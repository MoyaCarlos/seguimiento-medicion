import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { createMemoryRouter } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { Detalle } from './Detalle'
import type { DefectoDeSprint, DetalleSprint } from './tipos'

const defecto = (id: number, cambios: Partial<DefectoDeSprint> = {}): DefectoDeSprint => ({
  id,
  historia_id: 1,
  historia: 'Solicitar turno para el examen',
  descripcion: `Defecto ${id}`,
  severidad: 'Alta',
  estado: 'Pendiente',
  detectado_en_sprint: true,
  resuelto_en_sprint: false,
  ...cambios,
})

const detalle: DetalleSprint = {
  id: 71,
  proyecto_id: 7,
  numero: 1,
  sprint_goal: 'Gestión inicial del trámite',
  fecha_inicio: '2026-10-20',
  fecha_fin: '2026-10-26',
  metricas: {
    story_points_planificados: 8,
    story_points_completados: 5,
    velocidad: 5,
    porcentaje_historias_completadas: 50,
    horas_estimadas: 32,
    horas_reales: 4.33,
    desviacion_horas: -27.67,
    desviacion_porcentaje: -86.47,
    defectos_detectados: 2,
    defectos_resueltos: 1,
  },
  historias: [
    { id: 1, titulo: 'Solicitar turno para el examen', estimacion_sp: 5, completada: true, minutos: 260, horas: 4.33 },
    { id: 2, titulo: 'Validar el certificado psicofísico', estimacion_sp: 3, completada: false, minutos: 0, horas: 0 },
    { id: 3, titulo: 'Avisar el vencimiento', estimacion_sp: null, completada: false, minutos: 0, horas: 0 },
  ],
  defectos: [
    defecto(10),
    defecto(11, { severidad: 'Baja', estado: 'Resuelto', detectado_en_sprint: true, resuelto_en_sprint: true }),
    defecto(12, { severidad: 'Media', estado: 'Resuelto', detectado_en_sprint: false, resuelto_en_sprint: true }),
  ],
}

const respuesta = (status: number, cuerpo: unknown) => new Response(JSON.stringify(cuerpo), { status })

function respuestaPdf(nombre = 'reporte-proyecto-7-sprint-1.pdf') {
  return new Response('%PDF-1.4 contenido', {
    status: 200,
    headers: {
      'Content-Type': 'application/pdf',
      'Content-Disposition': `attachment; filename="${nombre}"`,
    },
  })
}

function simularApi(responder: (url: string) => Response | Promise<Response>) {
  const simulado = vi.fn((url: string) => Promise.resolve(responder(url)))
  vi.stubGlobal('fetch', simulado)
  return simulado
}

function simularDescarga() {
  const crear = vi.fn((_blob: Blob) => 'blob:reporte')
  const revocar = vi.fn()
  vi.stubGlobal('URL', { createObjectURL: crear, revokeObjectURL: revocar })

  let descargado: { nombre?: string; href?: string } = {}
  const click = vi.spyOn(HTMLAnchorElement.prototype, 'click').mockImplementation(function (
    this: HTMLAnchorElement,
  ) {
    descargado = { nombre: this.download, href: this.href }
  })

  return { crear, revocar, click, descargado: () => descargado }
}

function abrir() {
  const router = createMemoryRouter(
    [
      { path: '/projects/:id/historial', element: <h1>Lista del historial</h1> },
      { path: '/projects/:id/historial/:sprintId', element: <Detalle /> },
    ],
    { initialEntries: ['/projects/7/historial/71'] },
  )
  render(<RouterProvider router={router} />)
  return router
}

const seccion = (nombre: string) => screen.findByRole('region', { name: nombre })
const historia = (contenedor: HTMLElement, titulo: string) =>
  within(contenedor).getByText(titulo).closest('li') as HTMLElement
const defectoEn = (contenedor: HTMLElement, descripcion: string) =>
  within(contenedor).getByText(descripcion).closest('li') as HTMLElement

function metrica(contenedor: HTMLElement, nombre: string) {
  return within(contenedor).getByText(nombre, { selector: 'dt' }).nextElementSibling?.textContent
}

afterEach(() => {
  vi.unstubAllGlobals()
  vi.restoreAllMocks()
})

describe('detalle de un Sprint del historial', () => {
  it('Given una consulta en curso, When abro el detalle, Then veo un indicador de carga', () => {
    simularApi(() => new Promise<Response>(() => {}))

    abrir()

    expect(screen.getByRole('status').textContent).toMatch(/Cargando/)
  })

  it('US4.2: Given un Sprint cerrado, When lo selecciono, Then veo su número, Sprint Goal y fechas', async () => {
    const api = simularApi(() => respuesta(200, detalle))

    abrir()

    expect(await screen.findByRole('heading', { level: 1, name: 'Sprint 1' })).toBeTruthy()
    expect(screen.getByText('Gestión inicial del trámite')).toBeTruthy()
    expect(screen.getByText('20/10/2026 → 26/10/2026')).toBeTruthy()
    expect(api).toHaveBeenCalledWith('/projects/7/historial/71', expect.objectContaining({ method: 'GET' }))
  })

  it('US4.2: Given el detalle, When lo leo, Then muestra todas las métricas del Sprint', async () => {
    simularApi(() => respuesta(200, detalle))

    abrir()

    const metricas = await seccion('Métricas del Sprint')
    expect(metrica(metricas, 'Story Points planificados')).toBe('8')
    expect(metrica(metricas, 'Story Points completados')).toBe('5')
    expect(metrica(metricas, 'Velocidad')).toBe('5')
    expect(metrica(metricas, 'Historias completadas')).toBe('50%')
    expect(metrica(metricas, 'Horas estimadas')).toBe('32 h')
    expect(metrica(metricas, 'Horas reales')).toBe('4,33 h')
    expect(metrica(metricas, 'Desviación')).toBe('-27,67 h (-86,47%)')
    expect(metrica(metricas, 'Defectos detectados')).toBe('2')
    expect(metrica(metricas, 'Defectos resueltos')).toBe('1')
  })

  it('US4.2: Given historias comprometidas, When veo el detalle, Then cada una dice si se completó, sus Story Points y sus horas', async () => {
    simularApi(() => respuesta(200, detalle))

    abrir()

    const historias = await seccion('Historias comprometidas')
    const completada = historia(historias, 'Solicitar turno para el examen')
    expect(within(completada).getByText('Completada')).toBeTruthy()
    expect(within(completada).getByText('5 SP')).toBeTruthy()
    expect(within(completada).getByText('4,33 h')).toBeTruthy()

    const arrastrada = historia(historias, 'Validar el certificado psicofísico')
    expect(within(arrastrada).getByText('No completada')).toBeTruthy()
    expect(within(arrastrada).getByText('3 SP')).toBeTruthy()
    expect(within(arrastrada).getByText('0 h')).toBeTruthy()
  })

  it('US4.4: Given una historia sin Story Points, When veo el detalle, Then dice sin estimar y nunca 0 SP', async () => {
    simularApi(() => respuesta(200, detalle))

    abrir()

    const sinEstimar = historia(await seccion('Historias comprometidas'), 'Avisar el vencimiento')
    expect(within(sinEstimar).getByText(/sin estimar/i)).toBeTruthy()
    expect(within(sinEstimar).queryByText('0 SP')).toBeNull()
  })

  it('US4.2: Given horas imputadas, When veo el detalle, Then el esfuerzo total coincide con las horas reales del Sprint', async () => {
    simularApi(() => respuesta(200, detalle))

    abrir()

    const historias = await seccion('Historias comprometidas')
    expect(metrica(historias, 'Esfuerzo total del Sprint')).toBe('4,33 h')
  })

  it('US4.2: Given defectos del Sprint, When veo el detalle, Then cada uno indica si se detectó o resolvió acá y su estado actual', async () => {
    simularApi(() => respuesta(200, detalle))

    abrir()

    const defectos = await seccion('Defectos')
    const detectado = defectoEn(defectos, 'Defecto 10')
    expect(within(detectado).getByText('Alta')).toBeTruthy()
    expect(within(detectado).getByText('Pendiente')).toBeTruthy()
    expect(within(detectado).getByText('Detectado en este Sprint')).toBeTruthy()
    expect(within(detectado).queryByText('Resuelto en este Sprint')).toBeNull()
    expect(within(detectado).getByText(/Solicitar turno para el examen/)).toBeTruthy()

    const ambos = defectoEn(defectos, 'Defecto 11')
    expect(within(ambos).getByText('Detectado en este Sprint')).toBeTruthy()
    expect(within(ambos).getByText('Resuelto en este Sprint')).toBeTruthy()
    expect(within(ambos).getByText('Resuelto')).toBeTruthy()

    const resuelto = defectoEn(defectos, 'Defecto 12')
    expect(within(resuelto).queryByText('Detectado en este Sprint')).toBeNull()
    expect(within(resuelto).getByText('Resuelto en este Sprint')).toBeTruthy()
  })

  it('Caso límite: Given un Sprint sin historias ni defectos, When veo el detalle, Then lo indica en vez de listas vacías', async () => {
    simularApi(() => respuesta(200, { ...detalle, historias: [], defectos: [] }))

    abrir()

    expect(within(await seccion('Historias comprometidas')).getByText('Sin historias comprometidas.')).toBeTruthy()
    expect(within(await seccion('Defectos')).getByText('Sin defectos en este Sprint.')).toBeTruthy()
  })

  it('US4.4: Given un proyecto sin horas por Story Point, When veo el detalle, Then las métricas faltantes dicen sin datos', async () => {
    simularApi(() =>
      respuesta(200, {
        ...detalle,
        metricas: { ...detalle.metricas, horas_estimadas: null, desviacion_horas: null, desviacion_porcentaje: null },
      }),
    )

    abrir()

    const metricas = await seccion('Métricas del Sprint')
    expect(metrica(metricas, 'Horas estimadas')).toBe('sin datos')
    expect(metrica(metricas, 'Desviación')).toBe('sin datos')
  })

  it('US4.2: Given el detalle, When elijo volver, Then regreso a la lista del historial', async () => {
    simularApi(() => respuesta(200, detalle))
    const usuario = userEvent.setup()
    const router = abrir()
    await screen.findByRole('heading', { level: 1, name: 'Sprint 1' })

    await usuario.click(screen.getByRole('link', { name: /Volver al historial/ }))

    expect(router.state.location.pathname).toBe('/projects/7/historial')
  })

  it('US4.5: Given un Sprint todavía no finalizado, When abro su detalle, Then veo el mensaje y ningún dato parcial', async () => {
    simularApi(() => respuesta(400, { campo: 'estado', mensaje: 'el Sprint todavía no está finalizado' }))

    abrir()

    expect((await screen.findByRole('alert')).textContent).toBe('el Sprint todavía no está finalizado')
    expect(screen.queryByRole('region')).toBeNull()
    expect(screen.queryByRole('status')).toBeNull()
    expect(screen.getByRole('link', { name: /Volver al historial/ })).toBeTruthy()
  })

  it('US4.5: Given un Sprint inexistente, When abro su detalle, Then veo "sprint no encontrado"', async () => {
    simularApi(() => respuesta(404, { mensaje: 'sprint no encontrado' }))

    abrir()

    expect((await screen.findByRole('alert')).textContent).toBe('sprint no encontrado')
    expect(screen.queryByRole('heading', { level: 1, name: 'Sprint 1' })).toBeNull()
  })

  it('US5.1: Given el resumen de un Sprint finalizado, When presiono "Exportar a PDF", Then se descarga el PDF del Sprint', async () => {
    const usuario = userEvent.setup()
    const descarga = simularDescarga()
    const api = simularApi((url) =>
      url === '/sprints/71/reporte' ? respuestaPdf() : respuesta(200, detalle),
    )
    abrir()
    await screen.findByRole('heading', { level: 1, name: 'Sprint 1' })

    await usuario.click(screen.getByRole('button', { name: 'Exportar a PDF' }))

    await waitFor(() => expect(descarga.click).toHaveBeenCalledTimes(1))
    expect(api.mock.calls.some(([url]) => url === '/sprints/71/reporte')).toBe(true)
    expect(descarga.revocar).toHaveBeenCalledWith('blob:reporte')
    expect(descarga.descargado().nombre).toBe('reporte-proyecto-7-sprint-1.pdf')
    expect(await descarga.crear.mock.calls[0][0].text()).toBe('%PDF-1.4 contenido')
  })

  it('US5.3: Given que el backend informa un error de generación, When presiono "Exportar a PDF", Then veo el mensaje y no se descarga un archivo', async () => {
    const usuario = userEvent.setup()
    const descarga = simularDescarga()
    simularApi((url) =>
      url === '/sprints/71/reporte'
        ? respuesta(500, { mensaje: 'no se pudo generar el reporte' })
        : respuesta(200, detalle),
    )
    abrir()
    await screen.findByRole('heading', { level: 1, name: 'Sprint 1' })

    await usuario.click(screen.getByRole('button', { name: 'Exportar a PDF' }))

    expect((await screen.findByRole('alert')).textContent).toBe('no se pudo generar el reporte')
    expect(descarga.crear).not.toHaveBeenCalled()
    expect(descarga.click).not.toHaveBeenCalled()
  })

  it('US5.4: Given un Sprint todavía no finalizado, When abro su detalle, Then no ofrezco el botón de exportar', async () => {
    simularApi(() => respuesta(400, { campo: 'estado', mensaje: 'el Sprint todavía no está finalizado' }))

    abrir()

    expect((await screen.findByRole('alert')).textContent).toBe(
      'el Sprint todavía no está finalizado',
    )
    expect(screen.queryByRole('button', { name: 'Exportar a PDF' })).toBeNull()
  })
})
