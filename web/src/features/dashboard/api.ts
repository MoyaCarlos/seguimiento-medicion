import { get } from '../../compartido/api/client'
import type { MetricasProyecto } from './metricas'

export async function consultarMetricas(proyectoId: string): Promise<MetricasProyecto> {
  const metricas = await get<MetricasProyecto | undefined>(`/projects/${proyectoId}/metricas`)
  if (!Array.isArray(metricas?.sprints)) throw new Error('La respuesta no trae métricas')
  return metricas
}
