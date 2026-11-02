import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { BotonDatosDePrueba } from './BotonDatosDePrueba'

const json = (status: number, cuerpo: unknown) => new Response(JSON.stringify(cuerpo), { status })
const sinCuerpo = (status: number) => new Response('404 page not found', { status })

const resumen = { proyecto_id: 25, usuarios: 4, integrantes: 3, sprints: 3, historias: 10 }

function simularApi(estado: Response | Error, carga?: Response | Promise<Response>) {
  const fetchSimulado = vi.fn((ruta: string, _opciones?: RequestInit) => {
    if (ruta === '/demo/estado') return estado instanceof Error ? Promise.reject(estado) : Promise.resolve(estado)
    return carga ?? Promise.reject(new Error(`ruta inesperada ${ruta}`))
  })
  vi.stubGlobal('fetch', fetchSimulado)
  return fetchSimulado
}

function abrir() {
  render(
    <MemoryRouter>
      <BotonDatosDePrueba />
    </MemoryRouter>,
  )
}

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

  it('Dado el botón visible, cuando la carga responde 201, entonces muestra el aviso con el enlace al proyecto de prueba', async () => {
    const fetchSimulado = simularApi(json(200, { habilitado: true }), Promise.resolve(json(201, resumen)))
    abrir()

    await userEvent.click(await boton())

    const enlace = await screen.findByRole('link', { name: 'Abrir el proyecto de prueba' })
    expect(enlace.getAttribute('href')).toBe('/projects/25')
    expect(screen.getByRole('status').textContent).toContain('Datos de prueba cargados')
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
    expect(screen.getByRole('link', { name: 'Abrir el proyecto de prueba' }).getAttribute('href')).toBe('/projects/25')
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
})
