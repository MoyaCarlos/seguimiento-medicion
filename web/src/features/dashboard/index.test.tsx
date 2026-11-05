import { act, render, screen, waitFor, within } from '@testing-library/react'
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

  it('US2.1: Given valores por Sprint y del proyecto, When ingreso al Dashboard, Then los indicadores usan los del proyecto', async () => {
    simularApi(() =>
      respuesta(200, {
        ...conSprints(),
        sprints: [{ sprint_id: 70, numero: 1, metricas: { ...metricas, velocidad: 99 } }],
      }),
    )

    abrir()

    const velocidad = await screen.findByRole('region', { name: 'Velocidad' })
    expect(within(velocidad).getByText('5')).toBeTruthy()
    expect(within(velocidad).queryByText('99')).toBeNull()
  })

  it('US3.1: Given un proyecto sin Sprints finalizados, When ingreso al Dashboard, Then veo el mensaje, sin gráfico ni indicadores', async () => {
    simularApi(() => respuesta(200, { ...conSprints(), sprints: [] }))

    abrir()

    expect(await screen.findByText(/aparecen cuando se cierra el primer Sprint/)).toBeTruthy()
    expect(grafico()).toBeNull()
    expect(screen.queryByRole('region', { name: 'Velocidad' })).toBeNull()
    expect(screen.queryByRole('alert')).toBeNull()
  })

  it('US3.2: Given un único Sprint activo, When la consulta informa cero Sprints finalizados y valores en cero, Then veo el mismo estado vacío', async () => {
    const enCero: Metricas = {
      ...metricas,
      story_points_planificados: 0,
      story_points_completados: 0,
      velocidad: 0,
      porcentaje_historias_completadas: 0,
      horas_estimadas: 0,
      horas_reales: 0,
      desviacion_horas: 0,
      desviacion_porcentaje: null,
      defectos_detectados: 0,
      defectos_resueltos: 0,
    }
    simularApi(() => respuesta(200, { proyecto_id: 7, horas_por_story_point: 4, proyecto: enCero, sprints: [] }))

    abrir()

    expect(await screen.findByText(/aparecen cuando se cierra el primer Sprint/)).toBeTruthy()
    expect(screen.queryByRole('region', { name: 'Velocidad' })).toBeNull()
  })

  it('US3.3: Given el estado vacío, When lo leo, Then tengo un acceso a los Sprints del proyecto', async () => {
    simularApi(() => respuesta(200, { ...conSprints(), sprints: [] }))

    abrir()

    const enlace = await screen.findByRole('link', { name: /Sprints/ })
    expect(enlace.getAttribute('href')).toBe('/projects/7/sprints')
  })

  function sinContenido() {
    expect(grafico()).toBeNull()
    expect(screen.queryByRole('region', { name: 'Velocidad' })).toBeNull()
    expect(screen.queryByText(/aparecen cuando se cierra el primer Sprint/)).toBeNull()
    expect(screen.queryByRole('status')).toBeNull()
  }

  it('US4.2: Given un proyecto que no existe, When ingreso a su Dashboard, Then veo "proyecto no encontrado"', async () => {
    simularApi(() => respuesta(404, { mensaje: 'proyecto no encontrado' }))

    abrir()

    expect((await screen.findByRole('alert')).textContent).toBe('proyecto no encontrado')
    sinContenido()
  })

  it('US4.3: Given un servidor que no responde, When ingreso al Dashboard, Then veo el error de conexión', async () => {
    vi.stubGlobal('fetch', vi.fn(() => Promise.reject(new TypeError('Failed to fetch'))))

    abrir()

    expect((await screen.findByRole('alert')).textContent).toBe('No se pudo conectar con el servidor')
    sinContenido()
  })

  it('FR-016: Given un error interno, When ingreso al Dashboard, Then veo su mensaje, sin valores ni estado vacío', async () => {
    simularApi(() => respuesta(500, { mensaje: 'error interno' }))

    abrir()

    expect((await screen.findByRole('alert')).textContent).toBe('error interno')
    sinContenido()
  })

  it('FR-016: Given un error que no es de la API, When ingreso al Dashboard, Then no muestro su texto técnico', async () => {
    const cuerpoCortado = new Response('{}', { status: 200 })
    vi.spyOn(cuerpoCortado, 'text').mockRejectedValue(new TypeError('network error'))
    simularApi(() => cuerpoCortado)

    abrir()

    expect((await screen.findByRole('alert')).textContent).toBe('Ocurrió un error inesperado')
    sinContenido()
  })

  it.each([
    ['sin cuerpo', ''],
    ['sin Sprints en la respuesta', '{}'],
  ])('FR-016: Given una respuesta 200 %s, When ingreso al Dashboard, Then veo un error y no queda cargando', async (_caso, cuerpo) => {
    simularApi(() => new Response(cuerpo, { status: 200 }))

    abrir()

    expect((await screen.findByRole('alert')).textContent).toBe('Ocurrió un error inesperado')
    sinContenido()
  })
})
