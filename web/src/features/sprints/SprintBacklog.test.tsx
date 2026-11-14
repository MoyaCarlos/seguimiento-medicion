import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { SprintBacklog } from './SprintBacklog'
import type { HistoriaDeSprint } from './sprint'

const respuesta = (status: number, cuerpo: unknown) =>
  new Response(JSON.stringify(cuerpo), { status })

const historias: HistoriaDeSprint[] = [
  { id: 2, titulo: 'Recuperar clave', prioridad: 'S', estado: 'En progreso', estimacion_sp: 3 },
  { id: 1, titulo: 'Iniciar sesión', prioridad: 'M', estado: 'Nueva', estimacion_sp: null },
  { id: 3, titulo: 'Exportar reporte', prioridad: 'W', estado: 'Completada', estimacion_sp: 8 },
]

const SIN_CARGAS = { total_horas: 0, registros: [] }

function simularBacklog(backlog: HistoriaDeSprint[]) {
  const simulado = vi.fn((ruta: string, _opciones?: RequestInit) =>
    Promise.resolve(ruta.endsWith('/horas') ? respuesta(200, SIN_CARGAS) : respuesta(200, backlog)),
  )
  vi.stubGlobal('fetch', simulado)
  return simulado
}

afterEach(() => vi.unstubAllGlobals())

describe('SprintBacklog', () => {
  it('US3.2: al montarse pide GET /sprints/{id}/backlog', async () => {
    const fetchSimulado = simularBacklog(historias)

    render(<SprintBacklog sprintId={4} />)

    await screen.findByText('Recuperar clave')
    expect(fetchSimulado.mock.calls[0][0]).toBe('/sprints/4/backlog')
    expect(fetchSimulado.mock.calls[0][1]?.method).toBe('GET')
  })

  it('US3.2: muestra el indicador de carga mientras espera la respuesta', () => {
    vi.stubGlobal('fetch', vi.fn().mockReturnValue(new Promise(() => {})))

    render(<SprintBacklog sprintId={4} />)

    expect(screen.getByText('Cargando…')).toBeTruthy()
  })

  it('US3.2: lista las historias en el orden recibido', async () => {
    simularBacklog(historias)

    render(<SprintBacklog sprintId={4} />)

    const items = await screen.findAllByRole('listitem')
    expect(items.map((item) => within(item).getByRole('heading').textContent)).toEqual([
      'Recuperar clave',
      'Iniciar sesión',
      'Exportar reporte',
    ])
  })

  it('US3.2: cada historia muestra prioridad con su etiqueta en español, estado y Story Points', async () => {
    simularBacklog(historias)

    render(<SprintBacklog sprintId={4} />)

    const [primera, , tercera] = await screen.findAllByRole('listitem')
    expect(primera.textContent).toContain('Should have (S)')
    expect(primera.textContent).toContain('En progreso')
    expect(primera.textContent).toContain('3 SP')
    expect(tercera.textContent).toContain("Won't have (W)")
    expect(tercera.textContent).toContain('Completada')
    expect(tercera.textContent).toContain('8 SP')
  })

  it('US3.2: etiqueta las cuatro prioridades MoSCoW en español', async () => {
    const todas: HistoriaDeSprint[] = (['M', 'S', 'C', 'W'] as const).map((prioridad, i) => ({
      id: i + 1,
      titulo: `Historia ${prioridad}`,
      prioridad,
      estado: 'Nueva',
      estimacion_sp: 1,
    }))
    simularBacklog(todas)

    render(<SprintBacklog sprintId={4} />)

    const items = await screen.findAllByRole('listitem')
    expect(items.map((item) => item.textContent)).toEqual([
      expect.stringContaining('Must have (M)'),
      expect.stringContaining('Should have (S)'),
      expect.stringContaining('Could have (C)'),
      expect.stringContaining("Won't have (W)"),
    ])
  })

  it('US3.2: una historia sin estimación muestra "Sin estimar"', async () => {
    simularBacklog(historias)

    render(<SprintBacklog sprintId={4} />)

    const items = await screen.findAllByRole('listitem')
    expect(items[1].textContent).toContain('Sin estimar')
    expect(items[1].textContent).not.toContain('SP')
  })

  it('US3.2: un Sprint sin historias muestra el mensaje vacío y no un error', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(respuesta(200, [])))

    render(<SprintBacklog sprintId={4} />)

    expect(await screen.findByText('Este Sprint no tiene historias')).toBeTruthy()
    expect(screen.queryByRole('alert')).toBeNull()
    expect(screen.queryByRole('listitem')).toBeNull()
  })

  it('US3.2: si la API falla muestra su mensaje y no la lista', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue(respuesta(400, { campo: 'id', mensaje: 'identificador de sprint inválido' })),
    )

    render(<SprintBacklog sprintId={4} />)

    expect((await screen.findByRole('alert')).textContent).toBe('identificador de sprint inválido')
    expect(screen.queryByText('Este Sprint no tiene historias')).toBeNull()
    expect(screen.queryByText('Cargando…')).toBeNull()
  })

  it('US3.2: ignora la respuesta de un Sprint anterior si sprintId cambia', async () => {
    let resolverViejo: (r: Response) => void = () => {}
    const fetchSimulado = vi
      .fn()
      .mockReturnValueOnce(new Promise<Response>((resolve) => (resolverViejo = resolve)))
      .mockResolvedValueOnce(
        respuesta(200, [{ id: 9, titulo: 'Del Sprint nuevo', prioridad: 'C', estado: 'Nueva', estimacion_sp: 2 }]),
      )
    vi.stubGlobal('fetch', fetchSimulado)

    const { rerender } = render(<SprintBacklog sprintId={4} />)
    rerender(<SprintBacklog sprintId={5} />)
    await screen.findByText('Del Sprint nuevo')
    resolverViejo(respuesta(200, historias))
    await new Promise((r) => setTimeout(r, 0))

    expect(screen.queryByText('Recuperar clave')).toBeNull()
    expect(screen.getByText('Del Sprint nuevo')).toBeTruthy()
  })
})

describe('SprintBacklog con registro de horas', () => {
  it('US1.1: con imputación ofrece "Registrar horas" en cada historia', async () => {
    simularBacklog(historias)

    render(<SprintBacklog sprintId={4} imputacion={{ proyectoId: 42, fechaFin: '2026-11-09' }} />)

    await screen.findByText('Recuperar clave')
    expect(screen.getAllByRole('button', { name: 'Registrar horas' })).toHaveLength(historias.length)
  })

  it('US1.4: sin imputación no ofrece "Registrar horas"', async () => {
    simularBacklog(historias)

    render(<SprintBacklog sprintId={4} />)

    await screen.findByText('Recuperar clave')
    expect(screen.queryByRole('button', { name: 'Registrar horas' })).toBeNull()
  })

  it('FR-010: si falla el esfuerzo de una historia, el aviso queda en ella y las demás siguen visibles', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn((ruta: string) =>
        Promise.resolve(
          ruta === '/backlog/1/horas'
            ? respuesta(404, { mensaje: 'historia no encontrada' })
            : ruta.endsWith('/horas')
              ? respuesta(200, SIN_CARGAS)
              : respuesta(200, historias),
        ),
      ),
    )

    render(<SprintBacklog sprintId={4} />)

    const [primera, segunda, tercera] = await screen.findAllByRole('listitem')
    expect((await within(segunda).findByRole('alert')).textContent).toBe('historia no encontrada')
    expect(segunda.textContent).toContain('Iniciar sesión')
    expect(await within(primera).findByText('Esfuerzo real: 0 h')).toBeTruthy()
    expect(await within(tercera).findByText('Esfuerzo real: 0 h')).toBeTruthy()
    expect(screen.getAllByRole('alert')).toHaveLength(1)
  })
})

describe('SprintBacklog con cambio de estado', () => {
  const simularEstado = (respuestaPut: Response) => {
    const simulado = vi.fn((ruta: string, opciones?: RequestInit) =>
      Promise.resolve(
        opciones?.method === 'PUT'
          ? respuestaPut
          : ruta.endsWith('/horas')
            ? respuesta(200, SIN_CARGAS)
            : respuesta(200, historias),
      ),
    )
    vi.stubGlobal('fetch', simulado)
    return simulado
  }

  it('US3.1: con permiteCambiarEstado muestra las acciones de cada historia que las admite', async () => {
    simularEstado(respuesta(200, {}))

    render(<SprintBacklog sprintId={4} permiteCambiarEstado />)

    await screen.findByText('Recuperar clave')
    expect(screen.getAllByRole('button', { name: 'Iniciar' })).toHaveLength(1)
    expect(screen.getAllByRole('button', { name: 'Completar' })).toHaveLength(1)
  })

  it('US3.1: por defecto no muestra ninguna acción', async () => {
    simularEstado(respuesta(200, {}))

    render(<SprintBacklog sprintId={4} />)

    await screen.findByText('Recuperar clave')
    expect(screen.queryByRole('button')).toBeNull()
  })

  it('US3.1: iniciar una historia refleja el estado "En progreso" en la lista', async () => {
    const fetchSimulado = simularEstado(
      respuesta(200, { ...historias[1], estado: 'En progreso' }),
    )

    render(<SprintBacklog sprintId={4} permiteCambiarEstado />)
    await screen.findByText('Iniciar sesión')
    await userEvent.click(screen.getByRole('button', { name: 'Iniciar' }))

    const item = (await screen.findAllByRole('listitem'))[1]
    expect(await within(item).findByRole('button', { name: 'Completar' })).toBeTruthy()
    expect(item.textContent).toContain('En progreso')
    expect(fetchSimulado.mock.calls.find(([, o]) => o?.method === 'PUT')?.[0]).toBe('/backlog/1/estado')
  })

  it('US3.2: completar una historia muestra "Completada el dd/mm/aaaa" y quita la acción', async () => {
    simularEstado(
      respuesta(200, { ...historias[0], estado: 'Completada', fecha_finalizacion: '2026-10-29' }),
    )

    render(<SprintBacklog sprintId={4} permiteCambiarEstado />)
    await screen.findByText('Recuperar clave')
    await userEvent.click(screen.getByRole('button', { name: 'Completar' }))

    const item = (await screen.findAllByRole('listitem'))[0]
    expect(await within(item).findByText('Completada el 29/10/2026')).toBeTruthy()
    expect(within(item).queryByRole('button')).toBeNull()
  })

  it('US3.7: ante un 409 deja la lista como estaba y muestra el aviso de recargar', async () => {
    const mensaje = 'la historia cambió de estado mientras tanto: recargá y volvé a intentar'
    simularEstado(respuesta(409, { mensaje }))

    render(<SprintBacklog sprintId={4} permiteCambiarEstado />)
    await screen.findByText('Iniciar sesión')
    await userEvent.click(screen.getByRole('button', { name: 'Iniciar' }))

    expect((await screen.findByRole('alert')).textContent).toBe(mensaje)
    const item = screen.getAllByRole('listitem')[1]
    expect(item.textContent).toContain('Nueva')
    expect(within(item).getByRole('button', { name: 'Iniciar' })).toBeTruthy()
  })
})

describe('SprintBacklog como el prototipo', () => {
  const lista: HistoriaDeSprint[] = [
    { id: 1, titulo: 'Nueva con SP', prioridad: 'M', estado: 'Nueva', estimacion_sp: 5 },
    { id: 2, titulo: 'En curso', prioridad: 'S', estado: 'En progreso', estimacion_sp: null },
    { id: 3, titulo: 'Terminada', prioridad: 'C', estado: 'Completada', estimacion_sp: 8, fecha_finalizacion: '2026-10-29' },
    { id: 4, titulo: 'Descartada', prioridad: 'W', estado: 'Nueva', estimacion_sp: 1 },
  ]
  const imputacion = { proyectoId: 42, fechaFin: '2026-11-20' }

  async function items(props: object = {}) {
    simularBacklog(lista)
    render(<SprintBacklog sprintId={4} {...props} />)
    return screen.findAllByRole('listitem')
  }
  const chips = (item: HTMLElement) =>
    [...item.querySelectorAll('.sprint-chips > li')].map((c) => [c.textContent, c.className])

  it('Given las historias, When se listan, Then cada una es una tarjeta con su título', async () => {
    const [primera] = await items()

    expect(primera.className).toBe('sprint-historia')
    expect(within(primera).getByRole('heading').className).toBe('sprint-historia__titulo')
  })

  it('Given cada historia, When se muestran sus etiquetas, Then prioridad, estado y Story Points llevan su variante', async () => {
    const [nueva, enCurso, terminada, descartada] = await items()

    expect(chips(nueva)).toEqual([
      ['Must have (M)', 'sprint-chip sprint-chip--must'],
      ['Nueva', 'sprint-chip sprint-chip--nueva'],
      ['5 SP', 'sprint-chip sprint-chip--sp'],
    ])
    expect(chips(enCurso)).toEqual([
      ['Should have (S)', 'sprint-chip sprint-chip--should'],
      ['En progreso', 'sprint-chip sprint-chip--progreso'],
      ['Sin estimar', 'sprint-chip sprint-chip--ninguno'],
    ])
    expect(chips(terminada)).toEqual([
      ['Could have (C)', 'sprint-chip sprint-chip--could'],
      ['Completada', 'sprint-chip sprint-chip--ok'],
      ['8 SP', 'sprint-chip sprint-chip--sp'],
    ])
    expect(chips(descartada)[0]).toEqual(["Won't have (W)", 'sprint-chip sprint-chip--wont'])
  })

  it('Given una historia completada, When se muestra, Then la fecha va en una nota aparte', async () => {
    const [, , terminada] = await items()

    expect(within(terminada).getByText('Completada el 29/10/2026').className).toBe('sprint-historia__nota')
  })

  it('Given el esfuerzo real, When se muestra, Then va en verde de marca en su propia línea', async () => {
    const [primera] = await items()

    expect(within(primera).getByText('Esfuerzo real: 0 h').className).toBe('sprint-esfuerzo')
  })

  it('Given un Sprint activo, When se muestran las acciones, Then «Registrar horas» (con reloj) va antes que «Completar»/«Iniciar» (con ícono), alineadas en un pie', async () => {
    const [nueva, enCurso] = await items({ permiteCambiarEstado: true, imputacion })

    for (const [item, accion] of [[nueva, 'Iniciar'], [enCurso, 'Completar']] as const) {
      const pie = item.querySelector('.sprint-historia__acciones') as HTMLElement
      const botones = within(pie).getAllByRole('button')
      expect(botones.map((b) => b.textContent)).toEqual(['Registrar horas', accion])
      expect(botones[0].className).toBe('button button--secondary sprint-boton')
      expect(botones[1].className).toBe('button button--secondary sprint-boton sprint-boton--contorno')
      expect(botones.every((b) => b.querySelector('svg'))).toBe(true)
    }
  })

  it('Given un Sprint sin acciones, When se muestran las historias, Then no hay pie de acciones', async () => {
    const [primera] = await items()

    expect(primera.querySelector('.sprint-historia__acciones')).toBeNull()
  })

  it('Given una historia con cargas, When se muestra, Then «Ver detalle» es un control con triángulo y la lista usa dd/mm/aaaa', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn((ruta: string) =>
        Promise.resolve(
          ruta.endsWith('/horas')
            ? respuesta(200, { total_horas: 2, registros: [{ id: 1, integrante: 'Ana', fecha: '2026-10-27', actividad: 'Revisión', horas: 2 }] })
            : respuesta(200, [lista[0]]),
        ),
      ),
    )
    render(<SprintBacklog sprintId={4} imputacion={imputacion} />)

    const detalle = await screen.findByRole('button', { name: 'Ver detalle' })

    expect(detalle.className).toBe('sprint-detalle')
  })
})
