import { act, fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { createMemoryRouter } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import { afterEach, describe, expect, it, vi } from 'vitest'
import type { Historia } from './historia'
import { Listado } from './Listado'

const respuesta = (status: number, cuerpo: unknown) =>
  new Response(JSON.stringify(cuerpo), { status })

const MUST: Historia = {
  id: 3,
  proyecto_id: 42,
  titulo: 'Como usuario quiero iniciar sesión',
  descripcion: 'Autenticación con email y contraseña',
  criterios_aceptacion: 'Dado un email válido, cuando ingreso, entonces veo el tablero',
  prioridad: 'M',
  estado: 'Nueva',
  valor_negocio: 13,
  estimacion_sp: 5,
  sprint_id: null,
}

const SHOULD: Historia = {
  id: 1,
  proyecto_id: 42,
  titulo: 'Como usuario quiero cerrar sesión',
  descripcion: 'Cerrar la sesión activa',
  criterios_aceptacion: null,
  prioridad: 'S',
  estado: 'Nueva',
  valor_negocio: null,
  estimacion_sp: null,
  sprint_id: null,
}

const RUTA_LISTADO = /^\/projects\/([^/]+)\/backlog$/
const RUTA_EDICION = /^\/backlog\/(\d+)$/

type ListadoApi = (proyecto: string) => Response | Promise<Response>
type EdicionApi = (id: number, opciones?: RequestInit) => Response | Promise<Response>

function simularApi(listado: ListadoApi, edicion?: EdicionApi) {
  const fetchSimulado = vi.fn((ruta: string, opciones?: RequestInit) => {
    const coincideListado = RUTA_LISTADO.exec(ruta)
    if (coincideListado) return Promise.resolve(listado(coincideListado[1]))

    const coincideEdicion = RUTA_EDICION.exec(ruta)
    if (coincideEdicion && edicion) return Promise.resolve(edicion(Number(coincideEdicion[1]), opciones))

    return Promise.resolve(respuesta(404, { mensaje: 'no encontrado' }))
  })
  vi.stubGlobal('fetch', fetchSimulado)
  return fetchSimulado
}

const conHistorias = (...historias: Historia[]) => () => respuesta(200, historias)

function abrir(url = '/projects/42/backlog') {
  const router = createMemoryRouter(
    [
      {
        path: '/projects/:id/backlog',
        children: [
          { index: true, element: <Listado /> },
          { path: 'nueva', element: <p>Nueva historia</p> },
        ],
      },
    ],
    { initialEntries: [url] },
  )
  render(<RouterProvider router={router} />)
  return router
}

const articulos = () => screen.queryAllByRole('article')

const articuloDe = (titulo: string) =>
  screen.getByRole('heading', { name: titulo }).closest('article') as HTMLElement

const abrirEdicion = (titulo: string) =>
  userEvent.click(within(articuloDe(titulo)).getByRole('button', { name: 'Editar' }))

const botonGuardar = () => screen.getByRole('button', { name: 'Guardar cambios' })

function cuerpoEnviado(fetchSimulado: ReturnType<typeof vi.fn>, indice = 0) {
  const llamada = fetchSimulado.mock.calls.filter(([ruta]) => RUTA_EDICION.test(ruta))[indice]
  return JSON.parse((llamada[1] as RequestInit).body as string) as Record<string, unknown>
}

function mensajeDe(campo: HTMLElement) {
  const id = campo.getAttribute('aria-describedby')
  return id ? document.getElementById(id)?.textContent ?? null : null
}

afterEach(() => vi.unstubAllGlobals())

describe('listado del Product Backlog (US1)', () => {
  it('US1.1: consulta el backlog y lo muestra en el orden que devuelve la API', async () => {
    const fetchSimulado = simularApi(conHistorias(MUST, SHOULD))
    abrir()

    expect(screen.getByText('Cargando…')).toBeTruthy()
    expect(await screen.findByRole('heading', { name: MUST.titulo })).toBeTruthy()
    expect(
      articulos().map((articulo) => within(articulo).getByRole('heading', { level: 2 }).textContent),
    ).toEqual([MUST.titulo, SHOULD.titulo])
    expect(fetchSimulado.mock.calls[0][0]).toBe('/projects/42/backlog')

    const must = articuloDe(MUST.titulo)
    expect(within(must).getByText('Historia #3')).toBeTruthy()
    expect(within(must).getByText('Must have')).toBeTruthy()
    expect(within(must).getByText('Nueva')).toBeTruthy()
    expect(within(must).getByText('5 Story Points')).toBeTruthy()
    expect(within(must).getByText('Valor de negocio: 13')).toBeTruthy()
    expect(within(must).getByText(MUST.descripcion)).toBeTruthy()
    expect(within(must).getByText(MUST.criterios_aceptacion as string)).toBeTruthy()
    expect(within(must).queryByText('Cargando…')).toBeNull()

    const should = articuloDe(SHOULD.titulo)
    expect(within(should).getByText('Sin estimar')).toBeTruthy()
    expect(within(should).getByText('Valor de negocio: Sin valor de negocio')).toBeTruthy()
    expect(within(should).getByText('Sin criterios de aceptación')).toBeTruthy()
  })

  it('US1.2: con el backlog vacío muestra un mensaje y no lo trata como error', async () => {
    simularApi(conHistorias())
    abrir()

    expect(await screen.findByText('No hay historias en el backlog')).toBeTruthy()
    expect(articulos()).toHaveLength(0)
    expect(screen.queryByRole('alert')).toBeNull()
    expect(screen.queryByText('Cargando…')).toBeNull()
  })

  it('US1.3: si el proyecto no existe muestra el mensaje de la API', async () => {
    simularApi(() => respuesta(404, { mensaje: 'proyecto no encontrado' }))
    abrir()

    expect((await screen.findByRole('alert')).textContent).toBe('proyecto no encontrado')
    expect(articulos()).toHaveLength(0)
    expect(screen.queryByText('No hay historias en el backlog')).toBeNull()
  })

  it('si se cae la conexión muestra el error de conexión y deja la pantalla usable', async () => {
    const fetchSimulado = vi.fn().mockRejectedValue(new TypeError('Failed to fetch'))
    vi.stubGlobal('fetch', fetchSimulado)
    abrir()

    expect((await screen.findByRole('alert')).textContent).toBe('No se pudo conectar con el servidor')
    expect(screen.getByRole('link', { name: 'Nueva historia' })).toBeTruthy()
  })

  it('permite ir a la pantalla de creación de HU-01', async () => {
    simularApi(conHistorias())
    const router = abrir()

    await screen.findByText('No hay historias en el backlog')
    await userEvent.click(screen.getByRole('link', { name: 'Nueva historia' }))

    expect(await screen.findByText('Nueva historia')).toBeTruthy()
    expect(router.state.location.pathname).toBe('/projects/42/backlog/nueva')
  })
})

describe('edición de una historia (US2)', () => {
  it('US2.1: al guardar envía los cuatro campos y actualiza el listado conservando id y estado', async () => {
    const actualizada: Historia = {
      ...MUST,
      titulo: 'Como usuario quiero iniciar sesión con Google',
      descripcion: 'Autenticación con OAuth',
      criterios_aceptacion: 'Dado que tengo cuenta de Google, cuando ingreso, entonces veo el tablero',
      estimacion_sp: 8,
    }
    const fetchSimulado = simularApi(conHistorias(MUST), () => respuesta(200, actualizada))
    const user = userEvent.setup()
    abrir()
    await screen.findByRole('heading', { name: MUST.titulo })

    await abrirEdicion(MUST.titulo)

    expect((screen.getByLabelText('Título') as HTMLInputElement).value).toBe(MUST.titulo)
    expect((screen.getByLabelText('Descripción') as HTMLTextAreaElement).value).toBe(MUST.descripcion)
    expect((screen.getByLabelText('Criterios de aceptación') as HTMLTextAreaElement).value).toBe(
      MUST.criterios_aceptacion,
    )
    expect((screen.getByLabelText('Story Points') as HTMLSelectElement).value).toBe('5')

    await user.clear(screen.getByLabelText('Título'))
    await user.type(screen.getByLabelText('Título'), '  Como usuario quiero iniciar sesión con Google  ')
    await user.clear(screen.getByLabelText('Descripción'))
    await user.type(screen.getByLabelText('Descripción'), '  Autenticación con OAuth  ')
    await user.clear(screen.getByLabelText('Criterios de aceptación'))
    await user.type(
      screen.getByLabelText('Criterios de aceptación'),
      '  Dado que tengo cuenta de Google, cuando ingreso, entonces veo el tablero  ',
    )
    await user.selectOptions(screen.getByLabelText('Story Points'), '8')
    await user.click(botonGuardar())

    expect(await screen.findByText('Historia actualizada')).toBeTruthy()
    expect(fetchSimulado.mock.calls.filter(([ruta]) => RUTA_EDICION.test(ruta))).toHaveLength(1)
    expect(cuerpoEnviado(fetchSimulado)).toEqual({
      titulo: 'Como usuario quiero iniciar sesión con Google',
      descripcion: 'Autenticación con OAuth',
      criterios_aceptacion:
        'Dado que tengo cuenta de Google, cuando ingreso, entonces veo el tablero',
      estimacion_sp: 8,
    })

    const actualizado = articuloDe(actualizada.titulo)
    expect(within(actualizado).getByText('Historia #3')).toBeTruthy()
    expect(within(actualizado).getByText('Nueva')).toBeTruthy()
    expect(within(actualizado).getByText('8 Story Points')).toBeTruthy()
    expect(screen.queryByLabelText('Título')).toBeNull()
  })

  it('US2.1b: al guardar sin estimación ni criterios envía null', async () => {
    const fetchSimulado = simularApi(conHistorias(MUST), () =>
      respuesta(200, { ...MUST, criterios_aceptacion: null, estimacion_sp: null }),
    )
    const user = userEvent.setup()
    abrir()
    await screen.findByRole('heading', { name: MUST.titulo })

    await abrirEdicion(MUST.titulo)
    await user.clear(screen.getByLabelText('Criterios de aceptación'))
    await user.selectOptions(screen.getByLabelText('Story Points'), '')
    await user.click(botonGuardar())

    await screen.findByText('Historia actualizada')
    expect(cuerpoEnviado(fetchSimulado).criterios_aceptacion).toBeNull()
    expect(cuerpoEnviado(fetchSimulado).estimacion_sp).toBeNull()
  })

  it('US2.2: sin título advierte junto al campo y no llama a la API', async () => {
    const fetchSimulado = simularApi(conHistorias(MUST))
    const user = userEvent.setup()
    abrir()
    await screen.findByRole('heading', { name: MUST.titulo })

    await abrirEdicion(MUST.titulo)
    await user.clear(screen.getByLabelText('Título'))
    await user.click(botonGuardar())

    const titulo = screen.getByLabelText('Título')
    expect(titulo.getAttribute('aria-invalid')).toBe('true')
    expect(mensajeDe(titulo)).toBe('El título es obligatorio')
    expect(fetchSimulado.mock.calls.filter(([ruta]) => RUTA_EDICION.test(ruta))).toHaveLength(0)
  })

  it('US2.3: sin descripción advierte junto al campo y no llama a la API', async () => {
    const fetchSimulado = simularApi(conHistorias(MUST))
    const user = userEvent.setup()
    abrir()
    await screen.findByRole('heading', { name: MUST.titulo })

    await abrirEdicion(MUST.titulo)
    await user.clear(screen.getByLabelText('Descripción'))
    await user.click(botonGuardar())

    const descripcion = screen.getByLabelText('Descripción')
    expect(descripcion.getAttribute('aria-invalid')).toBe('true')
    expect(mensajeDe(descripcion)).toBe('La descripción es obligatoria')
    expect(fetchSimulado.mock.calls.filter(([ruta]) => RUTA_EDICION.test(ruta))).toHaveLength(0)
  })

  it('US2.4: con criterios de más de 2000 caracteres advierte y no llama a la API', async () => {
    const fetchSimulado = simularApi(conHistorias(MUST))
    const user = userEvent.setup()
    abrir()
    await screen.findByRole('heading', { name: MUST.titulo })

    await abrirEdicion(MUST.titulo)
    const criterios = screen.getByLabelText('Criterios de aceptación')
    fireEvent.change(criterios, { target: { value: 'a'.repeat(2001) } })
    await user.click(botonGuardar())

    expect(criterios.getAttribute('aria-invalid')).toBe('true')
    expect(mensajeDe(criterios)).toBe('Los criterios de aceptación no pueden superar los 2000 caracteres')
    expect(fetchSimulado.mock.calls.filter(([ruta]) => RUTA_EDICION.test(ruta))).toHaveLength(0)
  })

  it('US2.5: si la API rechaza los Story Points muestra su mensaje junto al campo y conserva los datos', async () => {
    const mensaje =
      'los Story Points ya cargados no se pueden cambiar en un Sprint Activo o Finalizado'
    simularApi(conHistorias(MUST), () => respuesta(400, { campo: 'estimacion_sp', mensaje }))
    const user = userEvent.setup()
    abrir()
    await screen.findByRole('heading', { name: MUST.titulo })

    await abrirEdicion(MUST.titulo)
    await user.clear(screen.getByLabelText('Título'))
    await user.type(screen.getByLabelText('Título'), 'Título conservado')
    await user.selectOptions(screen.getByLabelText('Story Points'), '8')
    await user.click(botonGuardar())

    expect(await screen.findByText(mensaje)).toBeTruthy()
    const storyPoints = screen.getByLabelText('Story Points')
    expect(storyPoints.getAttribute('aria-invalid')).toBe('true')
    expect(mensajeDe(storyPoints)).toBe(mensaje)
    expect((screen.getByLabelText('Título') as HTMLInputElement).value).toBe('Título conservado')
    expect(screen.queryByText('Historia actualizada')).toBeNull()
  })

  it('US2.6: si la API responde 404 al guardar muestra el aviso general y no cierra el formulario', async () => {
    simularApi(conHistorias(MUST), () => respuesta(404, { mensaje: 'historia no encontrada' }))
    const user = userEvent.setup()
    abrir()
    await screen.findByRole('heading', { name: MUST.titulo })

    await abrirEdicion(MUST.titulo)
    await user.clear(screen.getByLabelText('Título'))
    await user.type(screen.getByLabelText('Título'), 'Título conservado')
    await user.click(botonGuardar())

    expect((await screen.findByRole('alert')).textContent).toBe('historia no encontrada')
    expect((screen.getByLabelText('Título') as HTMLInputElement).value).toBe('Título conservado')
    expect(botonGuardar()).toBeTruthy()
  })

  it('cancelar cierra el formulario sin llamar a la API', async () => {
    const fetchSimulado = simularApi(conHistorias(MUST))
    const user = userEvent.setup()
    abrir()
    await screen.findByRole('heading', { name: MUST.titulo })

    await abrirEdicion(MUST.titulo)
    await user.click(screen.getByRole('button', { name: 'Cancelar' }))

    expect(screen.queryByLabelText('Título')).toBeNull()
    expect(fetchSimulado.mock.calls.filter(([ruta]) => RUTA_EDICION.test(ruta))).toHaveLength(0)
  })

  it('no ofrece editar prioridad, valor de negocio, estado ni Sprint', async () => {
    simularApi(conHistorias(MUST))
    abrir()
    await screen.findByRole('heading', { name: MUST.titulo })

    await abrirEdicion(MUST.titulo)

    expect(screen.queryByLabelText('Prioridad')).toBeNull()
    expect(screen.queryByLabelText('Valor de negocio')).toBeNull()
    expect(screen.queryByLabelText('Estado')).toBeNull()
    expect(screen.queryByLabelText('Sprint')).toBeNull()
  })
})

describe('casos límite del listado y la edición', () => {
  it('al cambiar de proyecto descarta las historias y los avisos del anterior', async () => {
    const fetchSimulado = simularApi((proyecto) =>
      respuesta(200, proyecto === '42' ? [MUST] : [SHOULD]),
    )
    const router = abrir('/projects/42/backlog')
    await screen.findByRole('heading', { name: MUST.titulo })

    await act(() => router.navigate('/projects/7/backlog'))

    await waitFor(() => expect(screen.queryByRole('heading', { name: MUST.titulo })).toBeNull())
    expect(await screen.findByRole('heading', { name: SHOULD.titulo })).toBeTruthy()
    expect(fetchSimulado).toHaveBeenCalledWith('/projects/7/backlog', expect.anything())
  })

  it('un error de edición no se mezcla con otro al cambiar de historia', async () => {
    const fetchSimulado = simularApi(conHistorias(MUST, SHOULD), (id) =>
      id === MUST.id
        ? respuesta(400, { campo: 'titulo', mensaje: 'error de la primera historia' })
        : respuesta(200, { ...SHOULD, titulo: 'Segunda actualizada' }),
    )
    const user = userEvent.setup()
    abrir()
    await screen.findByRole('heading', { name: MUST.titulo })

    await abrirEdicion(MUST.titulo)
    await user.click(botonGuardar())
    expect(await screen.findByText('error de la primera historia')).toBeTruthy()
    await user.click(screen.getByRole('button', { name: 'Cancelar' }))

    await abrirEdicion(SHOULD.titulo)
    expect(screen.queryByText('error de la primera historia')).toBeNull()
    await user.click(botonGuardar())

    expect(await screen.findByText('Historia actualizada')).toBeTruthy()
    expect(fetchSimulado.mock.calls.filter(([ruta]) => RUTA_EDICION.test(ruta))).toHaveLength(2)
  })
})
