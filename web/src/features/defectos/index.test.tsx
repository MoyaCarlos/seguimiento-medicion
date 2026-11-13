/// <reference types="node" />
import { readFileSync } from 'node:fs'
import { fireEvent, render, screen, waitFor } from '@testing-library/react'
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
  demoraRegistro?: Promise<void>
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
  demoraRegistro,
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
      const responder = () => respuesta(nuevo, 201)
      return demoraRegistro ? demoraRegistro.then(responder) : Promise.resolve(responder())
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

function diferido() {
  let resolver!: () => void
  const promesa = new Promise<void>((resolve) => {
    resolver = resolve
  })
  return { promesa, resolver }
}

async function completarFormulario(usuario: ReturnType<typeof userEvent.setup>) {
  await screen.findByRole('heading', { name: 'Registrar defecto' })
  await usuario.type(screen.getByLabelText('Descripción'), 'El total no se actualiza')
  await usuario.selectOptions(screen.getByLabelText('Severidad'), 'Alta')
  await usuario.selectOptions(screen.getByLabelText('Historia relacionada'), '3')
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('pantalla de defectos', () => {
  it('US1.1 — Dado un Sprint activo y una historia, cuando registro un defecto, entonces queda registrado y aparece en el listado', async () => {
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

  it('US3.1 — Dado un proyecto con defectos pendientes y resueltos, cuando consulto el listado, entonces muestra cada defecto con su descripción, severidad, estado, historia, Sprints y fecha de registro', async () => {
    abrir({
      defectos: [DEFECTO_PENDIENTE, DEFECTO_RESUELTO],
      sprints: [SPRINT_ACTIVO, SPRINT_FINALIZADO],
      sprintsBacklog: { 2: [HISTORIA_SPRINT], 4: [] },
    })

    expect(await screen.findByText('El total no se actualiza')).toBeTruthy()
    expect(screen.getByText('El botón no responde')).toBeTruthy()
    expect(screen.getAllByText('Severidad Alta').length).toBeGreaterThan(0)
    expect(screen.getByText('Pendiente')).toBeTruthy()
    expect(screen.getAllByText('Resuelto').length).toBeGreaterThan(0)
    expect(screen.getAllByText(/Carrito de compras/).length).toBeGreaterThan(0)
    expect(screen.getAllByText(/Sprint 1/).length).toBeGreaterThan(0)
    expect(screen.getByText('Detectado en Sprint 1 · Registrado: 04/11/2026')).toBeTruthy()
  })

  it('US2.1 — Dado un defecto pendiente y un Sprint activo, cuando lo marco como resuelto, entonces pasa a Resuelto', async () => {
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

  it('US1.3 y US2.3 — Dado un proyecto sin Sprint activo, cuando abro la pantalla, entonces avisa y no permite registrar ni resolver', async () => {
    abrir({ defectos: [DEFECTO_PENDIENTE], sprints: [SPRINT_FINALIZADO], sprintsBacklog: { 4: [] } })

    expect(await screen.findByText(/No hay un Sprint activo/)).toBeTruthy()
    expect(screen.queryByRole('button', { name: 'Registrar defecto' })).toBeNull()
    expect(screen.queryByRole('button', { name: 'Resolver' })).toBeNull()
    expect(screen.getByText('El total no se actualiza')).toBeTruthy()
  })

  it('US1.1 — Dado un proyecto sin historias, cuando abro la pantalla, entonces no muestra el formulario de registro', async () => {
    abrir({ defectos: [], sprints: [SPRINT_ACTIVO], backlog: [], sprintsBacklog: { 2: [] } })

    expect(await screen.findByText(/No hay historias en el proyecto/)).toBeTruthy()
    expect(screen.queryByRole('button', { name: 'Registrar defecto' })).toBeNull()
  })

  it('US1.4 — Dado que el backend rechaza el registro con un campo inválido, cuando guardo, entonces muestra el error en ese campo', async () => {
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

  it('Dado que falla la consulta del listado, cuando abro la pantalla, entonces muestra el mensaje y el formulario sigue disponible', async () => {
    abrir({
      errorDefectos: { status: 500, cuerpo: { mensaje: 'Falló el listado' } },
    })

    expect(await screen.findByText('Falló el listado')).toBeTruthy()
    expect(await screen.findByRole('heading', { name: 'Registrar defecto' })).toBeTruthy()
  })

  it('US3.2 — Dado un proyecto sin defectos, cuando consulto el listado, entonces muestra el estado vacío', async () => {
    abrir({ defectos: [] })

    expect(await screen.findByText('Todavía no hay defectos registrados.')).toBeTruthy()
  })

  it('US1.2 — Dado un proyecto con historias en el Product Backlog y en Sprints finalizados, cuando abro el formulario, entonces el selector las incluye a todas', async () => {
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

  it('US1.4 — Dado un formulario sin los campos obligatorios, cuando intento guardar, entonces no envía el registro y muestra los errores', async () => {
    const usuario = userEvent.setup()
    const { fetchMock } = abrir({ defectos: [] })

    await screen.findByRole('heading', { name: 'Registrar defecto' })
    await usuario.click(screen.getByRole('button', { name: 'Registrar defecto' }))

    expect(screen.getByText('La descripción es obligatoria')).toBeTruthy()
    expect(screen.getByText('La severidad es obligatoria')).toBeTruthy()
    expect(screen.getByText('La historia relacionada es obligatoria')).toBeTruthy()
    expect(fetchMock.mock.calls.some(([, init]) => init?.method === 'POST')).toBe(false)
  })

  it('US2.3 — Dado que el backend rechaza la resolución por falta de Sprint activo, cuando marco el defecto como resuelto, entonces muestra el mensaje', async () => {
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

  it('US1.1 — Dado un envío en curso, cuando envío dos veces seguidas, entonces registra un solo defecto', async () => {
    const usuario = userEvent.setup()
    const demora = diferido()
    const { fetchMock } = abrir({ defectos: [], demoraRegistro: demora.promesa })

    await completarFormulario(usuario)

    const formulario = screen.getByRole('form', { name: 'Registrar defecto' })
    fireEvent.submit(formulario)
    await waitFor(() => {
      const boton = screen.getByRole('button', { name: 'Registrar defecto' }) as HTMLButtonElement
      expect(boton.disabled).toBe(true)
    })
    fireEvent.submit(formulario)

    const envios = fetchMock.mock.calls.filter(
      ([url, init]) => url === '/projects/1/defectos' && init?.method === 'POST',
    )
    expect(envios).toHaveLength(1)

    demora.resolver()
    expect(await screen.findByText('Defecto registrado')).toBeTruthy()
  })

  it('US1.4 — Dado una descripción con espacios al principio y al final, cuando registro el defecto, entonces se guarda sin esos espacios', async () => {
    const usuario = userEvent.setup()
    const { fetchMock } = abrir({ defectos: [] })

    await screen.findByRole('heading', { name: 'Registrar defecto' })
    await usuario.type(screen.getByLabelText('Descripción'), '  El total no se actualiza  ')
    await usuario.selectOptions(screen.getByLabelText('Severidad'), 'Alta')
    await usuario.selectOptions(screen.getByLabelText('Historia relacionada'), '3')
    await usuario.click(screen.getByRole('button', { name: 'Registrar defecto' }))

    expect(await screen.findByText('Defecto registrado')).toBeTruthy()
    const envio = fetchMock.mock.calls.find(
      ([url, init]) => url === '/projects/1/defectos' && init?.method === 'POST',
    )
    expect(JSON.parse(String(envio?.[1]?.body))).toMatchObject({
      descripcion: 'El total no se actualiza',
    })
  })

  it('US1.3 — Dado un Sprint que se cierra entre la carga y el envío, cuando registro el defecto, entonces muestra el rechazo por falta de Sprint activo', async () => {
    const usuario = userEvent.setup()
    abrir({
      defectos: [],
      errorRegistro: {
        status: 400,
        cuerpo: {
          mensaje: 'el proyecto no tiene un Sprint activo al que asociar el defecto como Sprint de detección',
        },
      },
    })

    await completarFormulario(usuario)
    await usuario.click(screen.getByRole('button', { name: 'Registrar defecto' }))

    expect(
      await screen.findByText(
        'el proyecto no tiene un Sprint activo al que asociar el defecto como Sprint de detección',
      ),
    ).toBeTruthy()
  })

  it('US1.6 — Dado que el proyecto dejó de existir, cuando registro el defecto, entonces muestra "proyecto no encontrado"', async () => {
    const usuario = userEvent.setup()
    abrir({
      defectos: [],
      errorRegistro: { status: 404, cuerpo: { mensaje: 'proyecto no encontrado' } },
    })

    await completarFormulario(usuario)
    await usuario.click(screen.getByRole('button', { name: 'Registrar defecto' }))

    expect(await screen.findByText('proyecto no encontrado')).toBeTruthy()
  })

  it('US2.5 — Dado que el defecto no existe, cuando intento resolverlo, entonces muestra "defecto no encontrado"', async () => {
    const usuario = userEvent.setup()
    abrir({
      defectos: [DEFECTO_PENDIENTE],
      errorResolucion: { status: 404, cuerpo: { mensaje: 'defecto no encontrado' } },
    })

    await screen.findByText('El total no se actualiza')
    await usuario.click(screen.getByRole('button', { name: 'Resolver' }))

    expect(await screen.findByText('defecto no encontrado')).toBeTruthy()
  })

  it('Dado un conflicto de integridad al registrar (409), cuando guardo, entonces muestra el mensaje del servidor', async () => {
    const usuario = userEvent.setup()
    abrir({
      defectos: [],
      errorRegistro: { status: 409, cuerpo: { mensaje: 'violación de integridad referencial' } },
    })

    await completarFormulario(usuario)
    await usuario.click(screen.getByRole('button', { name: 'Registrar defecto' }))

    expect(await screen.findByText('violación de integridad referencial')).toBeTruthy()
  })
})

describe('listado de defectos como el prototipo', () => {
  const ambos = { defectos: [DEFECTO_PENDIENTE, DEFECTO_RESUELTO], sprints: [SPRINT_ACTIVO, SPRINT_FINALIZADO], sprintsBacklog: { 2: [HISTORIA_SPRINT], 4: [] } }

  it('Dado un defecto, cuando se lista, entonces su título es «Defecto #N · historia» en un encabezado de nivel 3', async () => {
    abrir(ambos)

    expect(await screen.findByRole('heading', { level: 3, name: 'Defecto #7 · Carrito de compras' })).toBeTruthy()
  })

  it('Dados dos defectos, cuando se lista, entonces va primero el más nuevo', async () => {
    abrir(ambos)

    const titulos = (await screen.findAllByRole('heading', { level: 3 })).map((titulo) => titulo.textContent)
    expect(titulos).toEqual(['Defecto #8 · Carrito de compras', 'Defecto #7 · Carrito de compras'])
  })

  it('Dado un defecto pendiente, cuando se lista, entonces «Resolver» va antes de la descripción', async () => {
    abrir({ defectos: [DEFECTO_PENDIENTE] })

    const boton = await screen.findByRole('button', { name: 'Resolver' })
    const descripcion = screen.getByText('El total no se actualiza')
    expect(boton.compareDocumentPosition(descripcion) & Node.DOCUMENT_POSITION_FOLLOWING).toBeTruthy()
  })

  it('Dado un defecto resuelto, cuando se lista, entonces los sprints y la fecha van en una sola línea separada por «·»', async () => {
    abrir(ambos)

    expect(await screen.findByText('Detectado en Sprint 1 · Resuelto en Sprint 1 · Registrado: 03/11/2026')).toBeTruthy()
  })

  it('Dado un defecto, cuando se lista, entonces muestra «Severidad X» y el estado como etiquetas', async () => {
    abrir(ambos)

    expect(await screen.findByText('Severidad Alta')).toBeTruthy()
    expect(screen.getByText('Severidad Baja')).toBeTruthy()
    expect(screen.getByText('Pendiente')).toBeTruthy()
    expect(screen.getByText('Resuelto')).toBeTruthy()
  })
})

describe('avisos de éxito como el prototipo', () => {
  it('Dado un defecto pendiente, cuando lo resuelvo, entonces aparece el aviso «Defecto #7 resuelto en el Sprint 1»', async () => {
    const usuario = userEvent.setup()
    abrir({ defectos: [DEFECTO_PENDIENTE] })

    await screen.findByText('El total no se actualiza')
    await usuario.click(screen.getByRole('button', { name: 'Resolver' }))

    const aviso = await screen.findByText('Defecto #7 resuelto en el Sprint 1')
    expect(aviso.closest('.defectos-exito')?.getAttribute('role')).toBe('status')
  })

  it('Dado que registro un defecto, cuando se guarda, entonces el aviso «Defecto registrado» va arriba, antes del formulario', async () => {
    const usuario = userEvent.setup()
    abrir({ defectos: [] })
    await completarFormulario(usuario)

    await usuario.click(screen.getByRole('button', { name: 'Registrar defecto' }))

    const aviso = await screen.findByText('Defecto registrado')
    const formulario = screen.getByRole('heading', { name: 'Registrar defecto' })
    expect(aviso.compareDocumentPosition(formulario) & Node.DOCUMENT_POSITION_FOLLOWING).toBeTruthy()
    expect(aviso.closest('form')).toBeNull()
  })

  it('Dado un aviso de éxito, cuando se estiliza, entonces es una caja verde con texto que cumple 4,5:1', () => {
    const css = readFileSync('src/features/defectos/defectos.css', 'utf8')
    const regla = css.match(/^\.defectos-exito\s*\{([^}]*)\}/m)?.[1] ?? ''
    const fondo = regla.match(/background:\s*(#[0-9a-f]{6})/)?.[1] ?? ''
    const texto = regla.match(/\bcolor:\s*(#[0-9a-f]{6})/)?.[1] ?? ''
    const luz = (h: string) =>
      [1, 3, 5].map((i) => parseInt(h.slice(i, i + 2), 16) / 255).map((c) => (c <= 0.03928 ? c / 12.92 : ((c + 0.055) / 1.055) ** 2.4))
        .reduce((a, c, i) => a + c * [0.2126, 0.7152, 0.0722][i], 0)

    expect(fondo).toBe('#dcf5e4')
    const [a, b] = [luz(fondo), luz(texto)].sort((x, y) => y - x)
    expect((a + 0.05) / (b + 0.05)).toBeGreaterThanOrEqual(4.5)
  })
})
