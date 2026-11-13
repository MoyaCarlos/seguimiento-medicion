import { render, screen, within } from '@testing-library/react'
import { createMemoryRouter } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { Lista } from './Lista'
import type { MetricasSprint, SprintHistorial } from './tipos'

const metricas: MetricasSprint = {
  story_points_planificados: 8,
  story_points_completados: 5,
  velocidad: 5,
  porcentaje_historias_completadas: 50,
  horas_estimadas: 32,
  horas_reales: 30.5,
  desviacion_horas: -1.5,
  desviacion_porcentaje: -4.69,
  defectos_detectados: 2,
  defectos_resueltos: 1,
}

const sprint = (numero: number, cambios: Partial<SprintHistorial> = {}): SprintHistorial => ({
  id: 70 + numero,
  proyecto_id: 7,
  numero,
  sprint_goal: `Objetivo del Sprint ${numero}`,
  fecha_inicio: '2026-10-20',
  fecha_fin: '2026-10-26',
  metricas,
  ...cambios,
})

function simularApi(responder: () => Response | Promise<Response>) {
  const simulado = vi.fn(() => Promise.resolve(responder()))
  vi.stubGlobal('fetch', simulado)
  return simulado
}

const respuesta = (status: number, cuerpo: unknown) => new Response(JSON.stringify(cuerpo), { status })

function abrir() {
  const router = createMemoryRouter(
    [{ path: '/projects/:id/historial', element: <Lista /> }],
    { initialEntries: ['/projects/7/historial'] },
  )
  render(<RouterProvider router={router} />)
}

const tarjeta = (nombre: string) => screen.findByRole('article', { name: nombre })

function metrica(contenedor: HTMLElement, nombre: string) {
  const termino = within(contenedor).getByText(nombre, { selector: 'dt' })
  return termino.nextElementSibling?.textContent
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('lista del historial de Sprints', () => {
  it('Given una consulta en curso, When ingreso al historial, Then veo el título y un indicador de carga', () => {
    simularApi(() => new Promise<Response>(() => {}))

    abrir()

    expect(screen.getByRole('heading', { level: 1, name: 'Historial de Sprints' })).toBeTruthy()
    expect(screen.getByRole('status').textContent).toMatch(/Cargando/)
  })

  it('US4.1: Given Sprints cerrados, When ingreso al historial, Then veo cada uno como "Sprint <número>" con su goal y fechas', async () => {
    const api = simularApi(() => respuesta(200, [sprint(1), sprint(2, { sprint_goal: 'Examen teórico', fecha_inicio: '2026-10-27', fecha_fin: '2026-11-09' })]))

    abrir()

    const segundo = await tarjeta('Sprint 2')
    expect(await tarjeta('Sprint 1')).toBeTruthy()
    expect(within(segundo).getByText('Examen teórico')).toBeTruthy()
    expect(within(segundo).getByText('27/10/2026 → 09/11/2026')).toBeTruthy()
    expect(api).toHaveBeenCalledWith('/projects/7/historial', expect.objectContaining({ method: 'GET' }))
  })

  it('US4.1: Given la lista, When la leo, Then el Sprint más nuevo aparece primero', async () => {
    simularApi(() => respuesta(200, [sprint(1), sprint(2), sprint(3)]))

    abrir()

    await tarjeta('Sprint 3')
    const nombres = screen.getAllByRole('article').map((articulo) => articulo.getAttribute('aria-label'))
    expect(nombres).toEqual(['Sprint 3', 'Sprint 2', 'Sprint 1'])
  })

  it('US4.1 y FR-003: Given un Sprint cerrado, When lo veo en la lista, Then muestra sus métricas principales', async () => {
    simularApi(() => respuesta(200, [sprint(1)]))

    abrir()

    const primero = await tarjeta('Sprint 1')
    expect(metrica(primero, 'Story Points planificados')).toBe('8')
    expect(metrica(primero, 'Story Points completados')).toBe('5')
    expect(metrica(primero, 'Historias completadas')).toBe('50 %')
    expect(metrica(primero, 'Horas estimadas')).toBe('32 h')
    expect(metrica(primero, 'Horas reales')).toBe('30,5 h')
    expect(metrica(primero, 'Desviación')).toBe('-1,5 h (-4,7 %)')
    expect(metrica(primero, 'Defectos detectados')).toBe('2')
    expect(metrica(primero, 'Defectos resueltos')).toBe('1')
  })

  it('US4.2: Given la lista, When elijo un Sprint, Then su enlace lleva a su detalle', async () => {
    simularApi(() => respuesta(200, [sprint(1)]))

    abrir()

    const primero = await tarjeta('Sprint 1')
    expect(within(primero).getByRole('link', { name: 'Sprint 1' }).getAttribute('href')).toBe('/projects/7/historial/71')
  })

  it('US4.3: Given un proyecto sin Sprints cerrados, When ingreso al historial, Then veo que todavía no hay Sprints cerrados', async () => {
    simularApi(() => respuesta(200, []))

    abrir()

    expect(await screen.findByText('Todavía no hay Sprints cerrados')).toBeTruthy()
    expect(screen.getByRole('link', { name: 'Ir a Sprints' }).getAttribute('href')).toBe('/projects/7/sprints')
    expect(screen.queryByRole('article')).toBeNull()
    expect(screen.queryByRole('alert')).toBeNull()
  })

  it('US4.4: Given un proyecto sin horas por Story Point, When veo la lista, Then lo que falta dice sin datos y nunca cero', async () => {
    simularApi(() =>
      respuesta(200, [
        sprint(1, { metricas: { ...metricas, horas_estimadas: null, desviacion_horas: null, desviacion_porcentaje: null } }),
      ]),
    )

    abrir()

    const primero = await tarjeta('Sprint 1')
    expect(metrica(primero, 'Horas estimadas')).toBe('sin datos')
    expect(metrica(primero, 'Desviación')).toBe('sin datos')
    expect(metrica(primero, 'Horas reales')).toBe('30,5 h')
  })

  it('Caso límite: Given horas de desviación sin porcentaje, When veo la lista, Then el porcentaje dice sin datos', async () => {
    simularApi(() =>
      respuesta(200, [sprint(1, { metricas: { ...metricas, horas_estimadas: 0, desviacion_horas: 3, desviacion_porcentaje: null } })]),
    )

    abrir()

    expect(metrica(await tarjeta('Sprint 1'), 'Desviación')).toBe('+3 h (sin datos)')
  })

  it('US4.5: Given que la consulta falla, When ingreso al historial, Then veo el mensaje de error y ningún dato parcial', async () => {
    simularApi(() => respuesta(404, { mensaje: 'proyecto no encontrado' }))

    abrir()

    expect((await screen.findByRole('alert')).textContent).toBe('proyecto no encontrado')
    expect(screen.queryByRole('article')).toBeNull()
    expect(screen.queryByText('Todavía no hay Sprints cerrados')).toBeNull()
    expect(screen.queryByRole('status')).toBeNull()
  })

  it('US4.5: Given que no hay conexión, When ingreso al historial, Then veo que no se pudo conectar', async () => {
    vi.stubGlobal('fetch', vi.fn(() => Promise.reject(new TypeError('fallo de red'))))

    abrir()

    expect((await screen.findByRole('alert')).textContent).toBe('No se pudo conectar con el servidor')
  })

  it('Given un Sprint, When veo la lista, Then el pretítulo va encima del título y la bajada debajo', async () => {
    simularApi(() => respuesta(200, [sprint(1)]))

    abrir()

    await tarjeta('Sprint 1')
    const titulo = screen.getByRole('heading', { level: 1, name: 'Historial de Sprints' })
    const pretitulo = screen.getByText('Proyecto #7')
    const bajada = screen.getByText('Compará el desempeño del equipo en cada Sprint cerrado.')
    expect(pretitulo.compareDocumentPosition(titulo) & Node.DOCUMENT_POSITION_FOLLOWING).toBeTruthy()
    expect(titulo.compareDocumentPosition(bajada) & Node.DOCUMENT_POSITION_FOLLOWING).toBeTruthy()
    expect(titulo.parentElement).toBe(bajada.parentElement)
  })

  it('Given un Sprint con 5 de 8 Story Points, When veo la tarjeta, Then dice «5 / 8 SP» junto a la barra', async () => {
    simularApi(() => respuesta(200, [sprint(1)]))

    abrir()

    const primero = await tarjeta('Sprint 1')
    expect(within(primero).getByText('5 / 8 SP')).toBeTruthy()
  })

  it('Given una tarjeta, When la recorro con el teclado, Then tiene un solo enlace con el nombre del Sprint y ningún otro control', async () => {
    simularApi(() => respuesta(200, [sprint(2)]))

    abrir()

    const segundo = await tarjeta('Sprint 2')
    expect(within(segundo).getAllByRole('link').map((enlace) => enlace.textContent)).toEqual(['Sprint 2'])
    expect(within(segundo).queryByRole('button')).toBeNull()
    expect(segundo.querySelector('a a, a button, button a')).toBeNull()
  })

  it('Given las fechas de un Sprint, When veo la tarjeta, Then van en una etiqueta con ícono junto al título', async () => {
    simularApi(() => respuesta(200, [sprint(1)]))

    abrir()

    const fechas = within(await tarjeta('Sprint 1')).getByText('20/10/2026 → 26/10/2026')
    expect(fechas.querySelector('svg[aria-hidden="true"]')).not.toBeNull()
    expect(fechas.closest('.historial-sprint__cabecera')).not.toBeNull()
  })
})
