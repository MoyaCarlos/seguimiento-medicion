import { get } from '../../compartido/api/client'
import type { MetricasProyecto } from './metricas'

export const consultarMetricas = (proyectoId: string) =>
  get<MetricasProyecto>(`/projects/${proyectoId}/metricas`)
