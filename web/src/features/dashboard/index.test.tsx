import { act, render, screen, waitFor } from '@testing-library/react'
import { createMemoryRouter } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { Component } from './index'
import type { Metricas, MetricasProyecto } from './metricas'

const metricas: Metricas = {
  story_points_planificados: 8,
  story_points_completados: 5,
  velocidad: 5,
  porcentaje_historias_completadas: 50,
  horas_estimadas: 32,
  horas_reales: 30,
  desviacion_horas: -2,
  desviacion_porcentaje: -6.25,
  defectos_detectados: 2,
  defectos_resueltos: 1,
}

const conSprints = (completados = 5): MetricasProyecto => ({
  proyecto_id: 7,
  horas_por_story_point: 4,
  proyecto: metricas,
  sprints: [{ sprint_id: 70, numero: 1, metricas: { ...metricas, story_points_completados: completados } }],
})

const respuesta = (status: number, cuerpo: unknown) => new Response(JSON.stringify(cuerpo), { status })

function simularApi(...respuestas: (() => Response | Promise<Response>)[]) {
  let llamada = 0
  const simulado = vi.fn(() => {
    const responder = respuestas[Math.min(llamada, respuestas.length - 1)]
    llamada += 1
    return Promise.resolve(responder())
  })
  vi.stubGlobal('fetch', simulado)
  return simulado
}

function abrir() {
  const router = createMemoryRouter(
    [
      { path: '/projects/:id/dashboard', element: <Component /> },
      { path: '/projects/:id/sprints', element: <p>Pantalla de Sprints</p> },
    ],
    { initialEntries: ['/projects/7/dashboard'] },
  )
  render(<RouterProvider router={router} />)
  return router
}

const grafico = () => screen.queryByRole('figure', { name: /Story Points comprometidos y completados/ })

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('sección Dashboard', () => {
  it('US4.1: Given una consulta en curso, When ingreso al Dashboard, Then veo un indicador de carga', () => {
    simularApi(() => new Promise<Response>(() => {}))

    abrir()

    expect(screen.getByRole('status').textContent).toMatch(/Cargando/)
    expect(grafico()).toBeNull()
  })

  it('US1.1: Given Sprints finalizados, When ingreso al Dashboard, Then veo el título y el gráfico', async () => {
    simularApi(() => respuesta(200, conSprints()))

    abrir()

    expect(screen.getByRole('heading', { level: 1, name: 'Dashboard' })).toBeTruthy()
    await waitFor(() => expect(grafico()).not.toBeNull())
    expect(screen.queryByRole('status')).toBeNull()
  })

  it('FR-002 y FR-003: When ingreso al Dashboard, Then hace una sola consulta y solo de lectura', async () => {
    const api = simularApi(() => respuesta(200, conSprints()))

    abrir()
    await waitFor(() => expect(grafico()).not.toBeNull())

    expect(api).toHaveBeenCalledTimes(1)
    expect(api).toHaveBeenCalledWith('/projects/7/metricas', expect.objectContaining({ method: 'GET' }))
  })

  it('Caso límite: Given que cambiaron las métricas, When salgo y vuelvo al Dashboard, Then consulta de nuevo y veo los valores nuevos', async () => {
    const api = simularApi(
      () => respuesta(200, conSprints(5)),
      () => respuesta(200, conSprints(7)),
    )
    const router = abrir()
    await screen.findByRole('cell', { name: '5', hidden: true })

    await act(() => router.navigate('/projects/7/sprints'))
    await screen.findByText('Pantalla de Sprints')
    await act(() => router.navigate('/projects/7/dashboard'))

    expect(await screen.findByRole('cell', { name: '7', hidden: true })).toBeTruthy()
    expect(api).toHaveBeenCalledTimes(2)
  })
})
