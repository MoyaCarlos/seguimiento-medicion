import { afterEach, expect, it, vi } from 'vitest'
import { consultarMetricas } from './api'
import type { MetricasProyecto } from './metricas'

afterEach(() => {
  vi.unstubAllGlobals()
})

it('consultarMetricas hace GET /projects/{id}/metricas y devuelve la respuesta', async () => {
  const respuesta: MetricasProyecto = {
    proyecto_id: 7,
    horas_por_story_point: null,
    proyecto: {
      story_points_planificados: 0,
      story_points_completados: 0,
      velocidad: 0,
      porcentaje_historias_completadas: 0,
      horas_estimadas: null,
      horas_reales: 0,
      desviacion_horas: null,
      desviacion_porcentaje: null,
      defectos_detectados: 0,
      defectos_resueltos: 0,
    },
    sprints: [],
  }
  const fetchMock = vi.fn().mockResolvedValue(new Response(JSON.stringify(respuesta), { status: 200 }))
  vi.stubGlobal('fetch', fetchMock)

  await expect(consultarMetricas('7')).resolves.toEqual(respuesta)

  expect(fetchMock).toHaveBeenCalledWith('/projects/7/metricas', expect.objectContaining({ method: 'GET' }))
})
