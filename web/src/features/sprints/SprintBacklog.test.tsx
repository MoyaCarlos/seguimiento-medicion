import { render, screen, within } from '@testing-library/react'
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

afterEach(() => vi.unstubAllGlobals())

describe('SprintBacklog', () => {
  it('US3.2: al montarse pide GET /sprints/{id}/backlog', async () => {
    const fetchSimulado = vi.fn().mockResolvedValue(respuesta(200, historias))
    vi.stubGlobal('fetch', fetchSimulado)

    render(<SprintBacklog sprintId={4} />)

    await screen.findByText('Recuperar clave')
    expect(fetchSimulado.mock.calls[0][0]).toBe('/sprints/4/backlog')
    expect(fetchSimulado.mock.calls[0][1].method).toBe('GET')
  })

  it('US3.2: muestra el indicador de carga mientras espera la respuesta', () => {
    vi.stubGlobal('fetch', vi.fn().mockReturnValue(new Promise(() => {})))

    render(<SprintBacklog sprintId={4} />)

    expect(screen.getByText('Cargando…')).toBeTruthy()
  })

  it('US3.2: lista las historias en el orden recibido', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(respuesta(200, historias)))

    render(<SprintBacklog sprintId={4} />)

    const items = await screen.findAllByRole('listitem')
    expect(items.map((item) => within(item).getByRole('heading').textContent)).toEqual([
      'Recuperar clave',
      'Iniciar sesión',
      'Exportar reporte',
    ])
  })

  it('US3.2: cada historia muestra prioridad con su etiqueta en español, estado y Story Points', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(respuesta(200, historias)))

    render(<SprintBacklog sprintId={4} />)

    const [primera, , tercera] = await screen.findAllByRole('listitem')
    expect(primera.textContent).toContain('Debería tener (S)')
    expect(primera.textContent).toContain('En progreso')
    expect(primera.textContent).toContain('3 SP')
    expect(tercera.textContent).toContain('No tendrá (W)')
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
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(respuesta(200, todas)))

    render(<SprintBacklog sprintId={4} />)

    const items = await screen.findAllByRole('listitem')
    expect(items.map((item) => item.textContent)).toEqual([
      expect.stringContaining('Debe tener (M)'),
      expect.stringContaining('Debería tener (S)'),
      expect.stringContaining('Podría tener (C)'),
      expect.stringContaining('No tendrá (W)'),
    ])
  })

  it('US3.2: una historia sin estimación muestra "Sin estimar"', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(respuesta(200, historias)))

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
