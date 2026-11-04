import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { PantallaDefectos } from './PantallaDefectos'

type DefectoPrueba = {
  id: number
  proyecto_id: number
  historia_id: number
  historia: string
  descripcion: string
  severidad: 'Alta' | 'Media' | 'Baja'
  estado: 'Pendiente' | 'Resuelto'
  sprint_deteccion_id: number
  sprint_resolucion_id: number | null
  creado_en: string
}

type SprintPrueba = { id: number; numero: number; estado: 'Pendiente' | 'Activo' | 'Finalizado' }
type HistoriaPrueba = { id: number; titulo: string }

const SPRINT_ACTIVO: SprintPrueba = { id: 2, numero: 1, estado: 'Activo' }
const SPRINT_FINALIZADO: SprintPrueba = { id: 4, numero: 2, estado: 'Finalizado' }
const HISTORIA_BACKLOG: HistoriaPrueba = { id: 3, titulo: 'Carrito de compras' }
const HISTORIA_SPRINT: HistoriaPrueba = { id: 5, titulo: 'Pago' }

const DEFECTO_PENDIENTE: DefectoPrueba = {
  id: 7,
  proyecto_id: 1,
  historia_id: 3,
  historia: 'Carrito de compras',
  descripcion: 'El total no se actualiza',
  severidad: 'Alta',
  estado: 'Pendiente',
  sprint_deteccion_id: 2,
  sprint_resolucion_id: null,
  creado_en: '2026-11-04T10:00:00Z',
}

const DEFECTO_RESUELTO: DefectoPrueba = {
  ...DEFECTO_PENDIENTE,
  id: 8,
  descripcion: 'El botón no responde',
  severidad: 'Baja',
  estado: 'Resuelto',
  sprint_resolucion_id: 2,
  creado_en: '2026-11-03T10:00:00Z',
}

type Opciones = {
  defectos?: DefectoPrueba[]
  sprints?: SprintPrueba[]
  backlog?: HistoriaPrueba[]
  sprintsBacklog?: Record<number, HistoriaPrueba[]>
  errorDefectos?: { status: number; cuerpo: unknown }
  errorRegistro?: { status: number; cuerpo: unknown }
  errorResolucion?: { status: number; cuerpo: unknown }
}

function respuesta(valor: unknown, status = 200) {
  return new Response(JSON.stringify(valor), {
    status,
    headers: { 'Content-Type': 'application/json' },
  })
}

function montarFetch({
  defectos = [DEFECTO_PENDIENTE],
  sprints = [SPRINT_ACTIVO],
  backlog = [HISTORIA_BACKLOG],
  sprintsBacklog = { 2: [HISTORIA_SPRINT] },
  errorDefectos,
  errorRegistro,
  errorResolucion,
}: Opciones = {}) {
  const estado = {
    defectos: defectos.map((defecto) => ({ ...defecto })),
    sprints: sprints.map((sprint) => ({ ...sprint })),
  }
  const fetchMock = vi.fn((url: string, init?: RequestInit) => {
    const metodo = init?.method ?? 'GET'

    if (url === '/projects/1/defectos' && metodo === 'GET') {
      if (errorDefectos) return Promise.resolve(respuesta(errorDefectos.cuerpo, errorDefectos.status))
      return Promise.resolve(respuesta(estado.defectos))
    }
    if (url === '/projects/1/sprints' && metodo === 'GET') {
      return Promise.resolve(respuesta(estado.sprints))
    }
    if (url === '/projects/1/backlog' && metodo === 'GET') {
      return Promise.resolve(respuesta(backlog))
    }
    if (url.startsWith('/sprints/') && url.endsWith('/backlog') && metodo === 'GET') {
      const sprintId = Number(url.split('/')[2])
      return Promise.resolve(respuesta(sprintsBacklog[sprintId] ?? []))
    }
    if (url === '/projects/1/defectos' && metodo === 'POST') {
      if (errorRegistro) return Promise.resolve(respuesta(errorRegistro.cuerpo, errorRegistro.status))
      const cuerpo = JSON.parse(String(init?.body))
      const nuevo: DefectoPrueba = {
        ...DEFECTO_PENDIENTE,
        id: 9,
        ...cuerpo,
        historia: 'Carrito de compras',
        estado: 'Pendiente',
        sprint_resolucion_id: null,
      }
      estado.defectos.push(nuevo)
      return Promise.resolve(respuesta(nuevo, 201))
    }
    if (url.startsWith('/defectos/') && url.endsWith('/resolver') && metodo === 'POST') {
      if (errorResolucion) return Promise.resolve(respuesta(errorResolucion.cuerpo, errorResolucion.status))
      const id = Number(url.split('/')[2])
      const defecto = estado.defectos.find((actual) => actual.id === id)
      if (defecto) {
        defecto.estado = 'Resuelto'
        defecto.sprint_resolucion_id = estado.sprints.find((sprint) => sprint.estado === 'Activo')?.id ?? null
      }
      return Promise.resolve(respuesta(defecto))
    }
    return Promise.resolve(respuesta([]))
  })
  vi.stubGlobal('fetch', fetchMock)
  return { fetchMock, estado }
}

function abrir(opciones: Opciones = {}) {
  const api = montarFetch(opciones)
  render(<PantallaDefectos proyectoId="1" />)
  return api
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('pantalla de defectos', () => {
  it('Escenario 1: con Sprint activo e historias, registrar un defecto lo muestra en el listado', async () => {
    const usuario = userEvent.setup()
    const { fetchMock } = abrir({ defectos: [] })

    await screen.findByRole('heading', { name: 'Registrar defecto' })
    await usuario.type(screen.getByLabelText('Descripción'), 'El total no se actualiza')
    await usuario.selectOptions(screen.getByLabelText('Severidad'), 'Alta')
    await usuario.selectOptions(screen.getByLabelText('Historia relacionada'), '3')
    await usuario.click(screen.getByRole('button', { name: 'Registrar defecto' }))

    expect(await screen.findByText('Defecto registrado')).toBeTruthy()
    expect(await screen.findByText('El total no se actualiza')).toBeTruthy()

    const llamada = fetchMock.mock.calls.find(
      ([url, init]) => url === '/projects/1/defectos' && init?.method === 'POST',
    )
    expect(JSON.parse(String(llamada?.[1]?.body))).toEqual({
      historia_id: 3,
      descripcion: 'El total no se actualiza',
      severidad: 'Alta',
    })
  })

  it('Escenario 2: el listado muestra descripción, severidad, estado, historia y Sprints', async () => {
    abrir({
      defectos: [DEFECTO_PENDIENTE, DEFECTO_RESUELTO],
      sprints: [SPRINT_ACTIVO, SPRINT_FINALIZADO],
      sprintsBacklog: { 2: [HISTORIA_SPRINT], 4: [] },
    })

    expect(await screen.findByText('El total no se actualiza')).toBeTruthy()
    expect(screen.getByText('El botón no responde')).toBeTruthy()
    expect(screen.getAllByText('Alta').length).toBeGreaterThan(0)
    expect(screen.getByText('Pendiente')).toBeTruthy()
    expect(screen.getAllByText('Resuelto').length).toBeGreaterThan(0)
    expect(screen.getAllByText(/Carrito de compras/).length).toBeGreaterThan(0)
    expect(screen.getAllByText(/Sprint 1/).length).toBeGreaterThan(0)
    expect(screen.getByText(/Registrado: 2026-11-04/)).toBeTruthy()
  })

  it('Escenario 3: resolver un defecto pendiente lo pasa a Resuelto', async () => {
    const usuario = userEvent.setup()
    const { fetchMock } = abrir({ defectos: [DEFECTO_PENDIENTE] })

    await screen.findByText('El total no se actualiza')
    await usuario.click(screen.getByRole('button', { name: 'Resolver' }))

    expect(await screen.findByText('Resuelto')).toBeTruthy()
    expect(screen.queryByRole('button', { name: 'Resolver' })).toBeNull()
    expect(
      fetchMock.mock.calls.some(([url, init]) => url === '/defectos/7/resolver' && init?.method === 'POST'),
    ).toBe(true)
  })

  it('Escenario 4: sin Sprint activo muestra el mensaje y no permite registrar ni resolver', async () => {
    abrir({ defectos: [DEFECTO_PENDIENTE], sprints: [SPRINT_FINALIZADO], sprintsBacklog: { 4: [] } })

    expect(await screen.findByText(/No hay un Sprint activo/)).toBeTruthy()
    expect(screen.queryByRole('button', { name: 'Registrar defecto' })).toBeNull()
    expect(screen.queryByRole('button', { name: 'Resolver' })).toBeNull()
    expect(screen.getByText('El total no se actualiza')).toBeTruthy()
  })

  it('Escenario alternativo: sin historias no se muestra el formulario', async () => {
    abrir({ defectos: [], sprints: [SPRINT_ACTIVO], backlog: [], sprintsBacklog: { 2: [] } })

    expect(await screen.findByText(/No hay historias en el proyecto/)).toBeTruthy()
    expect(screen.queryByRole('button', { name: 'Registrar defecto' })).toBeNull()
  })

  it('Escenario de error: el backend rechaza el registro y muestra el campo', async () => {
    const usuario = userEvent.setup()
    abrir({
      defectos: [],
      errorRegistro: {
        status: 400,
        cuerpo: { campo: 'descripcion', mensaje: 'la descripción es obligatoria' },
      },
    })

    await screen.findByRole('heading', { name: 'Registrar defecto' })
    await usuario.type(screen.getByLabelText('Descripción'), 'Algo')
    await usuario.selectOptions(screen.getByLabelText('Severidad'), 'Alta')
    await usuario.selectOptions(screen.getByLabelText('Historia relacionada'), '3')
    await usuario.click(screen.getByRole('button', { name: 'Registrar defecto' }))

    expect(await screen.findByText('la descripción es obligatoria')).toBeTruthy()
    expect(screen.getByLabelText('Descripción').getAttribute('aria-invalid')).toBe('true')
  })

  it('Escenario de error: si falla el listado, muestra el mensaje y el formulario sigue disponible', async () => {
    abrir({
      errorDefectos: { status: 500, cuerpo: { mensaje: 'Falló el listado' } },
    })

    expect(await screen.findByText('Falló el listado')).toBeTruthy()
    expect(await screen.findByRole('heading', { name: 'Registrar defecto' })).toBeTruthy()
  })

  it('Escenario alternativo: sin defectos muestra el estado vacío', async () => {
    abrir({ defectos: [] })

    expect(await screen.findByText('Todavía no hay defectos registrados.')).toBeTruthy()
  })

  it('Escenario límite: el selector incluye historias del Product Backlog y de Sprints finalizados', async () => {
    abrir({
      defectos: [],
      sprints: [SPRINT_FINALIZADO, SPRINT_ACTIVO],
      backlog: [HISTORIA_BACKLOG],
      sprintsBacklog: {
        4: [{ id: 6, titulo: 'Historia vieja' }],
        2: [HISTORIA_SPRINT],
      },
    })

    await screen.findByRole('heading', { name: 'Registrar defecto' })

    expect(screen.getByRole('option', { name: 'Carrito de compras (#3)' })).toBeTruthy()
    expect(screen.getByRole('option', { name: 'Historia vieja (#6)' })).toBeTruthy()
    expect(screen.getByRole('option', { name: 'Pago (#5)' })).toBeTruthy()
    expect(screen.getByRole('group', { name: 'Sprint 2' })).toBeTruthy()
    expect(screen.getByRole('group', { name: 'Sprint 1' })).toBeTruthy()
  })

  it('Escenario de error: sin los campos obligatorios no envía el registro', async () => {
    const usuario = userEvent.setup()
    const { fetchMock } = abrir({ defectos: [] })

    await screen.findByRole('heading', { name: 'Registrar defecto' })
    await usuario.click(screen.getByRole('button', { name: 'Registrar defecto' }))

    expect(screen.getByText('La descripción es obligatoria')).toBeTruthy()
    expect(screen.getByText('La severidad es obligatoria')).toBeTruthy()
    expect(screen.getByText('La historia relacionada es obligatoria')).toBeTruthy()
    expect(fetchMock.mock.calls.some(([, init]) => init?.method === 'POST')).toBe(false)
  })

  it('Escenario de error: si el backend rechaza la resolución, muestra el mensaje', async () => {
    const usuario = userEvent.setup()
    abrir({
      defectos: [DEFECTO_PENDIENTE],
      errorResolucion: {
        status: 400,
        cuerpo: { campo: 'sprint_resolucion_id', mensaje: 'el proyecto no tiene un Sprint activo' },
      },
    })

    await screen.findByText('El total no se actualiza')
    await usuario.click(screen.getByRole('button', { name: 'Resolver' }))

    expect(await screen.findByText('el proyecto no tiene un Sprint activo')).toBeTruthy()
  })
})
