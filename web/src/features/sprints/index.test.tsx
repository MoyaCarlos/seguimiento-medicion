import { act, render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { createMemoryRouter, MemoryRouter, Route, Routes } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { Component } from './index'
import type { Sprint } from './sprint'

const respuesta = (status: number, cuerpo: unknown) =>
  new Response(JSON.stringify(cuerpo), { status })

const pendiente = (id: number): Sprint => ({ id, proyecto_id: 42, estado: 'Pendiente' })

const RUTA_LISTADO = /^\/projects\/(\d+)\/sprints$/
const RUTA_BACKLOG = /^\/sprints\/(\d+)\/backlog$/

type Otras = (ruta: string, opciones?: RequestInit) => Response | Promise<Response>

type Backlog = (sprintId: string) => Response | Promise<Response>

function simularApi(
  listado: (proyecto: string) => Response | Promise<Response>,
  otras?: Otras,
  backlog: Backlog = () => respuesta(200, []),
) {
  const simulado = vi.fn((ruta: string, opciones?: RequestInit) => {
    const coincide = RUTA_LISTADO.exec(ruta)
    if (coincide) return Promise.resolve(listado(coincide[1]))
    const coincideBacklog = RUTA_BACKLOG.exec(ruta)
    if (coincideBacklog) return Promise.resolve(backlog(coincideBacklog[1]))
    return Promise.resolve(otras ? otras(ruta, opciones) : respuesta(404, { mensaje: 'no encontrado' }))
  })
  vi.stubGlobal('fetch', simulado)
  return simulado
}

const conSprints = (...sprints: Sprint[]) => () => respuesta(200, sprints)

function abrir(url = '/projects/42/sprints') {
  render(
    <MemoryRouter initialEntries={[url]}>
      <Routes>
        <Route path="/projects/:id/sprints" element={<Component />} />
      </Routes>
    </MemoryRouter>,
  )
}

const seccion = (nombre: string) =>
  screen.getByRole('heading', { name: nombre }).closest('section') as HTMLElement

const crear = () => userEvent.click(screen.getByRole('button', { name: 'Crear Sprint' }))

afterEach(() => vi.unstubAllGlobals())

describe('pantalla de Sprints', () => {
  it('conserva el contrato con la base: h1 "Sprints" y el id en un único <strong>', () => {
    simularApi(conSprints())
    abrir()

    expect(screen.getByRole('heading', { level: 1, name: 'Sprints' })).toBeTruthy()
    expect(screen.getAllByRole('heading', { name: 'Sprints' })).toHaveLength(1)
    expect(screen.getByText('42').tagName).toBe('STRONG')
    expect(screen.getByText('42').parentElement?.textContent).toBe('Proyecto actual: 42')
  })

  it('pide los Sprints del proyecto, muestra la carga y después los agrupa por estado', async () => {
    const activo: Sprint = { ...pendiente(2), estado: 'Activo', sprint_goal: 'Meta' }
    const finalizado: Sprint = { ...pendiente(1), estado: 'Finalizado', sprint_goal: 'Vieja' }
    const fetchSimulado = simularApi(conSprints(activo, pendiente(3), finalizado))
    abrir()

    expect(screen.getByText('Cargando…')).toBeTruthy()
    expect(await within(seccion('Pendiente')).findByRole('heading', { name: 'Sprint #3' })).toBeTruthy()
    expect(within(seccion('Activo')).getByRole('heading', { name: 'Sprint #2' })).toBeTruthy()
    expect(within(seccion('Finalizado')).getByRole('heading', { name: 'Sprint #1' })).toBeTruthy()
    await waitFor(() => expect(screen.queryAllByText('Cargando…')).toHaveLength(0))
    expect(fetchSimulado.mock.calls[0][0]).toBe('/projects/42/sprints')
  })

  it('US3.2: cada tarjeta de Sprint muestra su Sprint Backlog con las historias de ese Sprint', async () => {
    const activo: Sprint = { ...pendiente(2), estado: 'Activo', sprint_goal: 'Meta' }
    const historias: Record<string, unknown[]> = {
      '2': [{ id: 8, titulo: 'Historia del activo', prioridad: 'M', estado: 'En progreso', estimacion_sp: 5 }],
      '3': [],
    }
    const fetchSimulado = simularApi(
      conSprints(activo, pendiente(3)),
      undefined,
      (sprintId) => respuesta(200, historias[sprintId]),
    )
    abrir()

    await within(seccion('Pendiente')).findByRole('heading', { name: 'Sprint #3' })
    expect(await within(seccion('Activo')).findByText('Historia del activo')).toBeTruthy()
    expect(await within(seccion('Pendiente')).findByText('Este Sprint no tiene historias')).toBeTruthy()
    expect(within(seccion('Pendiente')).queryByText('Historia del activo')).toBeNull()
    const pedidos = fetchSimulado.mock.calls.map(([ruta]) => ruta)
    expect(pedidos).toContain('/sprints/2/backlog')
    expect(pedidos).toContain('/sprints/3/backlog')
  })

  it('US3.2: si falla el Sprint Backlog de una tarjeta el error queda en esa tarjeta y no en la pantalla', async () => {
    simularApi(conSprints(pendiente(3)), undefined, () =>
      respuesta(400, { mensaje: 'identificador de sprint inválido' }),
    )
    abrir()

    const tarjeta = (await within(seccion('Pendiente')).findByRole('heading', { name: 'Sprint #3' }))
      .closest('article') as HTMLElement
    expect((await within(tarjeta).findByRole('alert')).textContent).toBe('identificador de sprint inválido')
    expect(screen.getAllByRole('alert')).toHaveLength(1)
  })

  it('US3.2: un Sprint creado desde la pantalla muestra también su Sprint Backlog', async () => {
    const fetchSimulado = simularApi(conSprints(), () => respuesta(201, pendiente(5)))
    abrir()

    await crear()

    expect(await within(seccion('Pendiente')).findByText('Este Sprint no tiene historias')).toBeTruthy()
    expect(fetchSimulado.mock.calls.map(([ruta]) => ruta)).toContain('/sprints/5/backlog')
  })

  it('sin Sprints muestra las tres secciones vacías y ya no avisa de ninguna limitación', async () => {
    simularApi(conSprints())
    abrir()

    await waitFor(() => expect(screen.queryByText('Cargando…')).toBeNull())
    for (const nombre of ['Pendiente', 'Activo', 'Finalizado']) {
      expect(within(seccion(nombre)).getByText('No hay Sprints')).toBeTruthy()
    }
    expect(document.body.textContent).not.toContain('no se listan todavía')
  })

  it('si falla el listado muestra el mensaje de la API y deja la pantalla usable', async () => {
    simularApi(() => respuesta(404, { mensaje: 'proyecto no encontrado' }))
    abrir()

    expect((await screen.findByRole('alert')).textContent).toBe('proyecto no encontrado')
    expect(screen.queryByText('Cargando…')).toBeNull()
    expect(screen.getByRole('heading', { level: 1, name: 'Sprints' })).toBeTruthy()
    expect(within(seccion('Pendiente')).getByText('No hay Sprints')).toBeTruthy()
    expect((screen.getByRole('button', { name: 'Crear Sprint' }) as HTMLButtonElement).disabled).toBe(false)
  })

  it('crea un Sprint con el id del proyecto y lo muestra en Pendiente sin Goal ni fechas', async () => {
    const fetchSimulado = simularApi(conSprints(), () => respuesta(201, pendiente(5)))
    abrir()

    await crear()

    const tarjeta = await within(seccion('Pendiente')).findByRole('heading', { name: 'Sprint #5' })
    expect(tarjeta).toBeTruthy()
    const [ruta, opciones] = fetchSimulado.mock.calls.find(([r]) => r === '/sprints') as [string, RequestInit]
    expect(ruta).toBe('/sprints')
    expect(JSON.parse(opciones.body as string)).toEqual({ proyecto_id: 42 })
    expect(within(seccion('Pendiente')).queryByText(/\d{4}-\d{2}-\d{2}/)).toBeNull()
  })

  it('deshabilita "Crear Sprint" y muestra la carga mientras crea', async () => {
    simularApi(conSprints(), () => new Promise<Response>(() => {}))
    abrir()
    await waitFor(() => expect(screen.queryByText('Cargando…')).toBeNull())

    await crear()

    expect(screen.getByText('Cargando…')).toBeTruthy()
    expect((screen.getByRole('button', { name: 'Crear Sprint' }) as HTMLButtonElement).disabled).toBe(true)
  })

  it('muestra el error de crear en el aviso de la pantalla y no agrega ningún Sprint', async () => {
    simularApi(conSprints(), () => respuesta(400, { campo: 'estado', mensaje: 'ya hay un Sprint pendiente' }))
    abrir()

    await crear()

    expect((await screen.findByRole('alert')).textContent).toBe('ya hay un Sprint pendiente')
    expect(within(seccion('Pendiente')).getByText('No hay Sprints')).toBeTruthy()
  })

  it('recorre el ciclo: crear, iniciar y cerrar', async () => {
    const activo: Sprint = {
      ...pendiente(5),
      estado: 'Activo',
      sprint_goal: 'Meta',
      fecha_inicio: '2026-10-01',
      fecha_fin: '2026-10-15',
    }
    const respuestas: Record<string, Sprint> = {
      '/sprints': pendiente(5),
      '/sprints/5/iniciar': activo,
      '/sprints/5/cerrar': { ...activo, estado: 'Finalizado' },
    }
    simularApi(conSprints(), (ruta) => respuesta(200, respuestas[ruta]))
    const user = userEvent.setup()
    abrir()

    await crear()
    await screen.findByRole('heading', { name: 'Sprint #5' })
    await user.type(screen.getByLabelText('Sprint Goal'), 'Meta')
    await user.type(screen.getByLabelText('Fecha de inicio'), '2026-10-01')
    await user.type(screen.getByLabelText('Fecha de fin'), '2026-10-15')
    await user.click(screen.getByRole('button', { name: 'Iniciar Sprint' }))

    const enActivo = within(seccion('Activo'))
    expect(await enActivo.findByRole('heading', { name: 'Sprint #5' })).toBeTruthy()
    expect(enActivo.getByText('Meta')).toBeTruthy()
    expect(enActivo.getByText('2026-10-01 a 2026-10-15')).toBeTruthy()
    expect(within(seccion('Pendiente')).queryByRole('heading', { name: 'Sprint #5' })).toBeNull()
    expect(screen.queryByRole('button', { name: 'Iniciar Sprint' })).toBeNull()

    await user.click(screen.getByRole('button', { name: 'Cerrar Sprint' }))

    const enFinalizado = within(seccion('Finalizado'))
    expect(await enFinalizado.findByRole('heading', { name: 'Sprint #5' })).toBeTruthy()
    expect(within(seccion('Activo')).queryByRole('heading', { name: 'Sprint #5' })).toBeNull()
    expect(screen.queryByRole('button', { name: 'Cerrar Sprint' })).toBeNull()
    expect(screen.queryByRole('button', { name: 'Iniciar Sprint' })).toBeNull()
  })

  it('el error de un Sprint no aparece en otro', async () => {
    const respuestas = [respuesta(201, pendiente(1)), respuesta(201, pendiente(2))]
    simularApi(conSprints(), (ruta) =>
      ruta === '/sprints/1/iniciar'
        ? respuesta(400, { campo: 'estado', mensaje: 'error del Sprint 1' })
        : (respuestas.shift() as Response),
    )
    const user = userEvent.setup()
    abrir()

    await crear()
    await screen.findByRole('heading', { name: 'Sprint #1' })
    await crear()
    await screen.findByRole('heading', { name: 'Sprint #2' })
    const [primera, segunda] = screen.getAllByRole('article')
    await user.click(within(primera).getByRole('button', { name: 'Iniciar Sprint' }))

    expect((await within(primera).findByRole('alert')).textContent).toBe('error del Sprint 1')
    expect(within(segunda).queryByRole('alert')).toBeNull()
    await waitFor(() => expect(screen.getAllByRole('alert')).toHaveLength(1))
  })

  it('no valida el id de la URL: muestra el error que devuelve la API', async () => {
    const fetchSimulado = simularApi(conSprints(), () =>
      respuesta(400, { campo: 'proyecto_id', mensaje: 'proyecto_id inválido' }),
    )
    abrir('/projects/abc/sprints')

    await crear()

    expect((await screen.findByRole('alert')).textContent).toBe('proyecto_id inválido')
    expect(fetchSimulado.mock.calls.filter(([r]) => r === '/sprints')).toHaveLength(1)
  })

  it('al cambiar de proyecto recarga los Sprints y no deja a la vista los del anterior', async () => {
    const fetchSimulado = simularApi((proyecto) => respuesta(200, proyecto === '1' ? [pendiente(5)] : []))
    const router = createMemoryRouter(
      [{ path: '/projects/:id/sprints', element: <Component /> }],
      { initialEntries: ['/projects/1/sprints'] },
    )
    render(<RouterProvider router={router} />)
    await screen.findByRole('heading', { name: 'Sprint #5' })

    await act(() => router.navigate('/projects/2/sprints'))

    expect(screen.queryByRole('heading', { name: 'Sprint #5' })).toBeNull()
    await waitFor(() => expect(fetchSimulado).toHaveBeenCalledWith('/projects/2/sprints', expect.anything()))
    await waitFor(() => expect(screen.queryByText('Cargando…')).toBeNull())
    for (const nombre of ['Pendiente', 'Activo', 'Finalizado']) {
      expect(within(seccion(nombre)).getByText('No hay Sprints')).toBeTruthy()
    }
    expect(screen.getAllByRole('heading', { name: 'Sprints' })).toHaveLength(1)
    expect(screen.getByText('2').tagName).toBe('STRONG')
  })
})
