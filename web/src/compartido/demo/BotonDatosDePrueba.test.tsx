import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { BotonDatosDePrueba } from './BotonDatosDePrueba'

const json = (status: number, cuerpo: unknown) => new Response(JSON.stringify(cuerpo), { status })
const sinCuerpo = (status: number) => new Response('404 page not found', { status })

const resumen = { proyecto_id: 25, usuarios: 4, integrantes: 4, sprints: 4, historias: 12, defectos: 5 }

function simularApi(estado: Response | Error, carga?: Response | Promise<Response>) {
  const fetchSimulado = vi.fn((ruta: string, _opciones?: RequestInit) => {
    if (ruta === '/demo/estado') return estado instanceof Error ? Promise.reject(estado) : Promise.resolve(estado)
    return carga ?? Promise.reject(new Error(`ruta inesperada ${ruta}`))
  })
  vi.stubGlobal('fetch', fetchSimulado)
  return fetchSimulado
}

function abrir(alUsar?: (nombre: string, contrasena: string) => void) {
  render(
    <MemoryRouter>
      <BotonDatosDePrueba alUsar={alUsar} />
    </MemoryRouter>,
  )
}

const tablaDeUsuarios = () => screen.findByRole('table', { name: 'Usuarios de la demostración' })

const boton = () => screen.findByRole('button', { name: 'Cargar datos de prueba' })

afterEach(() => vi.unstubAllGlobals())

describe('BotonDatosDePrueba', () => {
  it('Dado que /demo/estado responde 200, cuando se muestra el inicio, entonces aparece el botón', async () => {
    simularApi(json(200, { habilitado: true }))
    abrir()

    expect(await boton()).toBeTruthy()
  })

  it('Dado que /demo/estado responde 404, cuando se muestra el inicio, entonces no aparece el botón', async () => {
    const fetchSimulado = simularApi(sinCuerpo(404))
    abrir()

    await waitFor(() => expect(fetchSimulado).toHaveBeenCalledWith('/demo/estado', expect.anything()))
    await Promise.resolve()
    expect(screen.queryByRole('button', { name: 'Cargar datos de prueba' })).toBeNull()
  })

  it('Dado que no hay conexión con /demo/estado, cuando se muestra el inicio, entonces no aparece el botón', async () => {
    const fetchSimulado = simularApi(new TypeError('Failed to fetch'))
    abrir()

    await waitFor(() => expect(fetchSimulado).toHaveBeenCalled())
    await Promise.resolve()
    expect(screen.queryByRole('button', { name: 'Cargar datos de prueba' })).toBeNull()
  })

  it('Dado el botón visible, cuando la carga responde 201, entonces avisa en línea que ya se puede ingresar con los usuarios de la lista, sin enlace', async () => {
    const fetchSimulado = simularApi(json(200, { habilitado: true }), Promise.resolve(json(201, resumen)))
    abrir()

    await userEvent.click(await boton())

    expect((await screen.findByRole('status')).textContent).toBe(
      'Datos de prueba cargados: ya podés ingresar con los usuarios de la lista',
    )
    expect(screen.queryByRole('link')).toBeNull()
    const [, opciones] = fetchSimulado.mock.calls.find(([ruta]) => ruta === '/demo/datos')!
    expect(opciones?.method).toBe('POST')
  })

  it('Dado el botón visible, cuando la carga responde 409 ya cargado, entonces lo muestra como aviso y no como error', async () => {
    simularApi(
      json(200, { habilitado: true }),
      Promise.resolve(json(409, { campo: 'ya_cargado', mensaje: 'los datos de prueba ya están cargados' })),
    )
    abrir()

    await userEvent.click(await boton())

    expect((await screen.findByRole('status')).textContent).toContain('los datos de prueba ya están cargados')
    expect(screen.queryByRole('link')).toBeNull()
    expect(screen.queryByRole('alert')).toBeNull()
  })

  it('Dado el botón visible, cuando la carga responde 409 por id ocupado, entonces muestra un error', async () => {
    simularApi(
      json(200, { habilitado: true }),
      Promise.resolve(json(409, { campo: 'proyecto_id', mensaje: 'el id 25 ya lo usa otro proyecto' })),
    )
    abrir()

    await userEvent.click(await boton())

    expect((await screen.findByRole('alert')).textContent).toBe('el id 25 ya lo usa otro proyecto')
    expect(screen.queryByRole('link', { name: 'Abrir el proyecto de prueba' })).toBeNull()
  })

  it('Dado el botón visible, cuando la carga responde 500, entonces muestra un error', async () => {
    simularApi(
      json(200, { habilitado: true }),
      Promise.resolve(json(500, { mensaje: 'no se pudieron cargar los datos de prueba' })),
    )
    abrir()

    await userEvent.click(await boton())

    expect((await screen.findByRole('alert')).textContent).toBe('no se pudieron cargar los datos de prueba')
  })

  it('Dado el botón visible, cuando la carga está en curso, entonces el botón queda deshabilitado', async () => {
    simularApi(json(200, { habilitado: true }), new Promise(() => {}))
    abrir()

    await userEvent.click(await boton())

    expect((screen.getByRole('button', { name: 'Cargar datos de prueba' }) as HTMLButtonElement).disabled).toBe(true)
    expect(screen.getByText('Cargando…')).toBeTruthy()
  })

  it('Dado el botón visible, cuando la carga responde 201, entonces muestra la tabla con usuario, contraseña y rol de los cuatro usuarios', async () => {
    simularApi(json(200, { habilitado: true }), Promise.resolve(json(201, resumen)))
    abrir()

    await userEvent.click(await boton())

    const tabla = await tablaDeUsuarios()
    const filas = within(tabla).getAllByRole('row').slice(1)
    expect(filas.map((fila) => within(fila).getAllByRole('cell').map((celda) => celda.textContent))).toEqual([
      ['carlos.moya', 'demo1234', 'Scrum Master'],
      ['ana.iriarte', 'demo1234', 'Product Builder'],
      ['candela.vulcano', 'demo1234', 'Product Builder'],
      ['jimena.martinez', 'demo1234', 'Product Builder'],
    ])
    expect(within(tabla).getAllByRole('columnheader').map((c) => c.textContent)).toEqual(['Usuario', 'Contraseña', 'Rol'])
  })

  it('Dado el botón visible, cuando la carga responde 409 ya cargado, entonces también muestra la tabla de usuarios', async () => {
    simularApi(
      json(200, { habilitado: true }),
      Promise.resolve(json(409, { campo: 'ya_cargado', mensaje: 'los datos de prueba ya están cargados' })),
    )
    abrir()

    await userEvent.click(await boton())

    expect(await tablaDeUsuarios()).toBeTruthy()
  })

  it('Dado el botón visible, cuando la carga responde 409 por usuario existente, entonces avisa con el nombre del usuario, sin error ni tabla', async () => {
    simularApi(
      json(200, { habilitado: true }),
      Promise.resolve(json(409, { campo: 'usuario_existente', mensaje: 'ya existe el usuario ana.iriarte' })),
    )
    abrir()

    await userEvent.click(await boton())

    expect((await screen.findByRole('status')).textContent).toContain('ya existe el usuario ana.iriarte')
    expect(screen.queryByRole('alert')).toBeNull()
    expect(screen.queryByRole('table')).toBeNull()
    expect((screen.getByRole('button', { name: 'Cargar datos de prueba' }) as HTMLButtonElement).disabled).toBe(false)
  })

  it('Dado el botón visible, cuando la carga responde un error, entonces no muestra la tabla de usuarios', async () => {
    simularApi(
      json(200, { habilitado: true }),
      Promise.resolve(json(409, { campo: 'proyecto_id', mensaje: 'el id 25 ya lo usa otro proyecto' })),
    )
    abrir()

    await userEvent.click(await boton())

    await screen.findByRole('alert')
    expect(screen.queryByRole('table')).toBeNull()
  })

  it('Dado la tabla de usuarios y un destino para «Usar», cuando uso a ana.iriarte, entonces se entregan su usuario y su contraseña', async () => {
    simularApi(json(200, { habilitado: true }), Promise.resolve(json(201, resumen)))
    const alUsar = vi.fn()
    abrir(alUsar)

    await userEvent.click(await boton())
    const fila = within(await tablaDeUsuarios()).getByRole('row', { name: /ana\.iriarte/ })
    await userEvent.click(within(fila).getByRole('button', { name: 'Usar ana.iriarte' }))

    expect(alUsar).toHaveBeenCalledWith('ana.iriarte', 'demo1234')
  })

  it('Dado la tabla de usuarios y un destino para «Usar», entonces tiene tres columnas y «Usar» queda debajo del rol', async () => {
    simularApi(json(200, { habilitado: true }), Promise.resolve(json(201, resumen)))
    abrir(vi.fn())

    await userEvent.click(await boton())
    const tabla = await tablaDeUsuarios()

    expect(within(tabla).getAllByRole('columnheader').map((c) => c.textContent)).toEqual(['Usuario', 'Contraseña', 'Rol'])
    const filas = within(tabla).getAllByRole('row').slice(1)
    expect(filas.map((fila) => within(fila).getAllByRole('cell').length)).toEqual([3, 3, 3, 3])
    const celdaDelRol = within(filas[1]).getAllByRole('cell')[2]
    expect(celdaDelRol.textContent).toBe('Product BuilderUsar')
    expect(within(celdaDelRol).getByRole('button', { name: 'Usar ana.iriarte' })).toBeTruthy()
  })

  it('Dado la tabla de usuarios, entonces el usuario y la contraseña se muestran como código', async () => {
    simularApi(json(200, { habilitado: true }), Promise.resolve(json(201, resumen)))
    abrir()

    await userEvent.click(await boton())
    const [primera] = within(await tablaDeUsuarios()).getAllByRole('row').slice(1)

    const [usuario, contrasena] = within(primera).getAllByRole('cell')
    expect(within(usuario).getByText('carlos.moya').tagName).toBe('CODE')
    expect(within(contrasena).getByText('demo1234').tagName).toBe('CODE')
  })

  it('Dado la tabla de usuarios sin destino para «Usar», entonces no hay botones «Usar»', async () => {
    simularApi(json(200, { habilitado: true }), Promise.resolve(json(201, resumen)))
    abrir()

    await userEvent.click(await boton())
    await tablaDeUsuarios()

    expect(screen.queryByRole('button', { name: /^Usar / })).toBeNull()
  })
})
