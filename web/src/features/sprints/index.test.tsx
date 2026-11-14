import { act, render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { createMemoryRouter, MemoryRouter, Route, Routes } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { Component } from './index'
import type { Sprint } from './sprint'

const respuesta = (status: number, cuerpo: unknown) =>
  new Response(JSON.stringify(cuerpo), { status })

const pendiente = (id: number, numero = id): Sprint => ({ id, numero, proyecto_id: 42, estado: 'Pendiente' })

const RUTA_LISTADO = /^\/projects\/(\d+)\/sprints$/
const RUTA_BACKLOG = /^\/sprints\/(\d+)\/backlog$/
const RUTA_HORAS = /^\/backlog\/\d+\/horas$/

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
    if (RUTA_HORAS.test(ruta)) return Promise.resolve(respuesta(200, { total_horas: 1.5, registros: [] }))
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
  it('numeración por proyecto: la tarjeta muestra el número del Sprint y no su id', async () => {
    simularApi(conSprints(pendiente(7, 2)))
    abrir()

    expect(await screen.findByRole('heading', { name: 'Sprint 2' })).toBeTruthy()
    expect(screen.queryByRole('heading', { name: 'Sprint 7' })).toBeNull()
  })

  it('conserva el contrato con la base: un único h1 "Sprints" y el pretítulo «Proyecto #42» encima', () => {
    simularApi(conSprints())
    abrir()

    const titulo = screen.getByRole('heading', { level: 1, name: 'Sprints' })
    const pretitulo = screen.getByText('Proyecto #42')

    expect(screen.getAllByRole('heading', { name: 'Sprints' })).toHaveLength(1)
    expect(pretitulo.compareDocumentPosition(titulo) & Node.DOCUMENT_POSITION_FOLLOWING).toBeTruthy()
    expect(screen.queryByText(/Proyecto actual/)).toBeNull()
  })

  it('pide los Sprints del proyecto, muestra la carga y después los agrupa por estado', async () => {
    const activo: Sprint = { ...pendiente(2), estado: 'Activo', sprint_goal: 'Meta' }
    const finalizado: Sprint = { ...pendiente(1), estado: 'Finalizado', sprint_goal: 'Vieja' }
    const fetchSimulado = simularApi(conSprints(activo, pendiente(3), finalizado))
    abrir()

    expect(screen.getByText('Cargando…')).toBeTruthy()
    expect(await within(seccion('Pendiente')).findByRole('heading', { name: 'Sprint 3' })).toBeTruthy()
    expect(within(seccion('Activo')).getByRole('heading', { name: 'Sprint 2' })).toBeTruthy()
    expect(within(seccion('Finalizado')).getByRole('heading', { name: 'Sprint 1' })).toBeTruthy()
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

    await within(seccion('Pendiente')).findByRole('heading', { name: 'Sprint 3' })
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

    const tarjeta = (await within(seccion('Pendiente')).findByRole('heading', { name: 'Sprint 3' }))
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

    const tarjeta = await within(seccion('Pendiente')).findByRole('heading', { name: 'Sprint 5' })
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
    await screen.findByRole('heading', { name: 'Sprint 5' })
    await user.type(screen.getByLabelText('Sprint Goal'), 'Meta')
    await user.type(screen.getByLabelText('Fecha de inicio'), '2026-10-01')
    await user.type(screen.getByLabelText('Fecha de fin'), '2026-10-15')
    await user.click(screen.getByRole('button', { name: 'Iniciar Sprint' }))

    const enActivo = within(seccion('Activo'))
    expect(await enActivo.findByRole('heading', { name: 'Sprint 5' })).toBeTruthy()
    expect(enActivo.getByText('Meta')).toBeTruthy()
    expect(enActivo.getByText('01/10/2026 → 15/10/2026')).toBeTruthy()
    expect(within(seccion('Pendiente')).queryByRole('heading', { name: 'Sprint 5' })).toBeNull()
    expect(screen.queryByRole('button', { name: 'Iniciar Sprint' })).toBeNull()

    await user.click(screen.getByRole('button', { name: 'Cerrar Sprint' }))

    const enFinalizado = within(seccion('Finalizado'))
    expect(await enFinalizado.findByRole('heading', { name: 'Sprint 5' })).toBeTruthy()
    expect(within(seccion('Activo')).queryByRole('heading', { name: 'Sprint 5' })).toBeNull()
    expect(screen.queryByRole('button', { name: 'Cerrar Sprint' })).toBeNull()
    expect(screen.queryByRole('button', { name: 'Iniciar Sprint' })).toBeNull()
  })

  it('el error de un Sprint no aparece en otro', async () => {
    const respuestas = [respuesta(201, pendiente(1)), respuesta(201, pendiente(2))]
    simularApi(conSprints(), (ruta) =>
      ruta === '/projects/42'
        ? respuesta(200, { id: 42, rol: 'scrum_master' })
        : ruta === '/projects'
          ? respuesta(200, [])
          : ruta === '/sprints/1/iniciar'
          ? respuesta(400, { campo: 'estado', mensaje: 'error del Sprint 1' })
          : (respuestas.shift() as Response),
    )
    const user = userEvent.setup()
    abrir()

    await crear()
    await screen.findByRole('heading', { name: 'Sprint 1' })
    await crear()
    await screen.findByRole('heading', { name: 'Sprint 2' })
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
    await screen.findByRole('heading', { name: 'Sprint 5' })

    await act(() => router.navigate('/projects/2/sprints'))

    expect(screen.queryByRole('heading', { name: 'Sprint 5' })).toBeNull()
    await waitFor(() => expect(fetchSimulado).toHaveBeenCalledWith('/projects/2/sprints', expect.anything()))
    await waitFor(() => expect(screen.queryByText('Cargando…')).toBeNull())
    for (const nombre of ['Pendiente', 'Activo', 'Finalizado']) {
      expect(within(seccion(nombre)).getByText('No hay Sprints')).toBeTruthy()
    }
    expect(screen.getAllByRole('heading', { name: 'Sprints' })).toHaveLength(1)
    expect(screen.getByText('Proyecto #2')).toBeTruthy()
  })
})

describe('pantalla de Sprints: cambio de estado de historias', () => {
  const historia = { id: 7, titulo: 'Exportar reporte', prioridad: 'M', estado: 'Nueva', estimacion_sp: 5 }

  it('US3.1: solo los Sprints Activos muestran las acciones de las historias', async () => {
    const activo: Sprint = { ...pendiente(2), estado: 'Activo' }
    const finalizado: Sprint = { ...pendiente(4), estado: 'Finalizado' }
    simularApi(conSprints(activo, pendiente(3), finalizado), undefined, () => respuesta(200, [historia]))

    abrir()
    await waitFor(() => expect(within(seccion('Activo')).getByText('Exportar reporte')).toBeTruthy())
    await waitFor(() => expect(within(seccion('Pendiente')).getByText('Exportar reporte')).toBeTruthy())
    await waitFor(() => expect(within(seccion('Finalizado')).getByText('Exportar reporte')).toBeTruthy())

    expect(within(seccion('Activo')).getByRole('button', { name: 'Iniciar' })).toBeTruthy()
    expect(within(seccion('Pendiente')).queryByRole('button', { name: 'Iniciar' })).toBeNull()
    expect(within(seccion('Finalizado')).queryByRole('button', { name: 'Iniciar' })).toBeNull()
  })
})

describe('pantalla de Sprints: registro de horas', () => {
  const historia = { id: 7, titulo: 'Exportar reporte', prioridad: 'M', estado: 'Nueva', estimacion_sp: 5 }

  it('FR-001: "Registrar horas" solo aparece en el Sprint Activo, y el total en todos', async () => {
    const activo: Sprint = { ...pendiente(2), estado: 'Activo', fecha_fin: '2026-11-09' }
    const finalizado: Sprint = { ...pendiente(4), estado: 'Finalizado' }
    simularApi(conSprints(activo, pendiente(3), finalizado), undefined, () => respuesta(200, [historia]))

    abrir()

    for (const nombre of ['Activo', 'Pendiente', 'Finalizado']) {
      expect(await within(seccion(nombre)).findByText('Esfuerzo real: 1,5 h')).toBeTruthy()
    }
    expect(within(seccion('Activo')).getByRole('button', { name: 'Registrar horas' })).toBeTruthy()
    expect(within(seccion('Pendiente')).queryByRole('button', { name: 'Registrar horas' })).toBeNull()
    expect(within(seccion('Finalizado')).queryByRole('button', { name: 'Registrar horas' })).toBeNull()
  })

  it('al cerrar el Sprint desaparece "Registrar horas" y sigue el total de la historia completada', async () => {
    const activo: Sprint = { ...pendiente(2), estado: 'Activo', fecha_fin: '2026-11-09' }
    simularApi(
      conSprints(activo),
      () => respuesta(200, { ...activo, estado: 'Finalizado' }),
      () => respuesta(200, [{ ...historia, estado: 'Completada' }]),
    )
    abrir()
    await within(seccion('Activo')).findByRole('button', { name: 'Registrar horas' })

    await userEvent.click(screen.getByRole('button', { name: 'Cerrar Sprint' }))

    expect(await within(seccion('Finalizado')).findByText('Esfuerzo real: 1,5 h')).toBeTruthy()
    expect(screen.queryByRole('button', { name: 'Registrar horas' })).toBeNull()
  })
})

describe('pantalla de Sprints - acciones de gestión según el rol (HU-21)', () => {
  const activo: Sprint = { ...pendiente(2), estado: 'Activo', sprint_goal: 'Meta' }

  const conRol = (rol: 'scrum_master' | 'product_builder') =>
    simularApi(conSprints(pendiente(1), activo), (ruta) =>
      ruta === '/projects/42' ? respuesta(200, { id: 42, rol }) : respuesta(404, { mensaje: 'no encontrado' }),
    )

  const nota = (boton: HTMLElement) =>
    document.getElementById(boton.getAttribute('aria-describedby') ?? '')?.textContent

  it('Escenario 2: como Product Builder veo crear, iniciar y cerrar Sprint deshabilitados con la nota «Solo el Scrum Master puede …»', async () => {
    conRol('product_builder')
    abrir()

    const crearBoton = await screen.findByRole('button', { name: 'Crear Sprint' })
    const iniciar = await screen.findByRole('button', { name: 'Iniciar Sprint' })
    const cerrar = await screen.findByRole('button', { name: 'Cerrar Sprint' })
    await waitFor(() => expect(crearBoton.hasAttribute('disabled')).toBe(true))
    expect(iniciar.hasAttribute('disabled')).toBe(true)
    expect(cerrar.hasAttribute('disabled')).toBe(true)
    expect(nota(crearBoton)).toBe('Solo el Scrum Master puede crear un Sprint')
    expect(nota(iniciar)).toBe('Solo el Scrum Master puede iniciar un Sprint')
    expect(nota(cerrar)).toBe('Solo el Scrum Master puede cerrar un Sprint')
  })

  it('Escenario 1: como Scrum Master crear, iniciar y cerrar Sprint están habilitados', async () => {
    const simulado = conRol('scrum_master')
    abrir()

    const crearBoton = await screen.findByRole('button', { name: 'Crear Sprint' })
    const iniciar = await screen.findByRole('button', { name: 'Iniciar Sprint' })
    const cerrar = await screen.findByRole('button', { name: 'Cerrar Sprint' })
    await waitFor(() => expect(simulado).toHaveBeenCalledWith('/projects/42', expect.anything()))
    expect(crearBoton.hasAttribute('disabled')).toBe(false)
    expect(iniciar.hasAttribute('disabled')).toBe(false)
    expect(cerrar.hasAttribute('disabled')).toBe(false)
    expect(crearBoton.hasAttribute('aria-describedby')).toBe(false)
  })
})

describe('pantalla de Sprints como el prototipo: cabecera, plazo y tablero', () => {
  const activo: Sprint = { ...pendiente(2), estado: 'Activo', sprint_goal: 'Meta del Sprint', fecha_inicio: '2026-11-07', fecha_fin: '2026-11-20' }
  const proyectos = [
    { id: 41, nombre: 'Otro', descripcion: '', fecha_inicio: '2026-01-01', fecha_fin: '2026-02-02' },
    { id: 42, nombre: 'Alfa', descripcion: '', fecha_inicio: '2026-10-15', fecha_fin: '2026-12-29' },
  ]
  const conProyectos = (lista: unknown) =>
    simularApi(conSprints(pendiente(1), activo), (ruta) =>
      ruta === '/projects' ? respuesta(200, lista) : respuesta(404, { mensaje: 'no encontrado' }),
    )

  it('Given la pantalla, When abre, Then «Crear Sprint» va en la cabecera como botón principal con el ícono «+»', async () => {
    conProyectos(proyectos)
    abrir()

    const boton = await screen.findByRole('button', { name: 'Crear Sprint' })

    expect(boton.className).toBe('button button--primary')
    expect(boton.querySelector('svg')).toBeTruthy()
    expect(boton.closest('header')?.querySelector('h1')?.textContent).toBe('Sprints')
  })

  it('Given un proyecto con fechas, When abre, Then muestra el plazo del proyecto y la regla de los Sprints', async () => {
    conProyectos(proyectos)
    abrir()

    const plazo = (await screen.findByText('Plazo del proyecto:')).parentElement as HTMLElement

    expect(plazo.textContent).toBe(
      'Plazo del proyecto: 15/10/2026 → 29/12/2026. Los Sprints tienen que quedar dentro de este rango.',
    )
    expect(plazo.querySelector('strong')?.textContent).toBe('Plazo del proyecto:')
  })

  it.each([
    ['sin fechas', [{ ...proyectos[1], fecha_inicio: null, fecha_fin: null }]],
    ['sin ese proyecto en la lista', [proyectos[0]]],
    ['con la lista caída', null],
  ])('Given un proyecto %s, When abre, Then no muestra el plazo', async (_caso, lista) => {
    if (lista === null) {
      simularApi(conSprints(pendiente(1)), () => respuesta(500, { mensaje: 'caído' }))
    } else conProyectos(lista)
    abrir()

    await screen.findByRole('heading', { name: 'Sprint 1' })

    expect(screen.queryByText('Plazo del proyecto:')).toBeNull()
  })

  it('Given los Sprints, When abre, Then las tres columnas están en el tablero con sus variantes', async () => {
    conProyectos(proyectos)
    abrir()

    await screen.findByRole('heading', { name: 'Sprint 1' })
    const columnas = ['Pendiente', 'Activo', 'Finalizado'].map(seccion)

    expect(columnas.map((c) => c.className)).toEqual([
      'sprints-carril',
      'sprints-carril sprints-carril--activo',
      'sprints-carril sprints-carril--fin',
    ])
    expect(new Set(columnas.map((c) => c.parentElement))).toHaveProperty('size', 1)
    expect(columnas[0].parentElement?.className).toBe('sprints-tablero')
  })

  it('Given una columna sin Sprints, When abre, Then dice «No hay Sprints» en un recuadro punteado', async () => {
    conProyectos(proyectos)
    abrir()

    await screen.findByRole('heading', { name: 'Sprint 1' })

    expect(within(seccion('Finalizado')).getByText('No hay Sprints').className).toBe('sprints-vacio')
  })

  it('Given un Sprint con Goal y fechas, When abre, Then el Goal va en gris y las fechas en una etiqueta con calendario', async () => {
    conProyectos(proyectos)
    abrir()

    const tarjeta = (await screen.findByRole('heading', { name: 'Sprint 2' })).closest('article') as HTMLElement
    const fechas = within(tarjeta).getByText('07/11/2026 → 20/11/2026')

    expect(within(tarjeta).getByText('Meta del Sprint').className).toBe('sprint-goal')
    expect(fechas.className).toBe('sprint-fechas')
    expect(fechas.querySelector('svg')).toBeTruthy()
  })

  it('Given un Sprint sin Goal ni fechas, When abre, Then la tarjeta no dibuja etiquetas vacías', async () => {
    conProyectos(proyectos)
    abrir()

    const tarjeta = (await screen.findByRole('heading', { name: 'Sprint 1' })).closest('article') as HTMLElement

    expect(tarjeta.querySelector('.sprint-goal')).toBeNull()
    expect(tarjeta.querySelector('.sprint-fechas')).toBeNull()
  })
})
