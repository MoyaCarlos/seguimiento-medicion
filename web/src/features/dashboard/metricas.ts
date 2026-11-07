export type Metricas = {
  story_points_planificados: number
  story_points_completados: number
  velocidad: number
  porcentaje_historias_completadas: number
  horas_estimadas: number | null
  horas_reales: number
  desviacion_horas: number | null
  desviacion_porcentaje: number | null
  defectos_detectados: number
  defectos_resueltos: number
}

export type MetricasDeSprint = {
  sprint_id: number
  numero: number
  metricas: Metricas
}

export type MetricasProyecto = {
  proyecto_id: number
  horas_por_story_point: number | null
  proyecto: Metricas
  sprints: MetricasDeSprint[]
}

export type FilaDelGrafico = {
  sprint: string
  comprometidos: number
  completados: number
}

export function datosDelGrafico(sprints: MetricasDeSprint[]): FilaDelGrafico[] {
  return sprints.map(({ numero, metricas }) => ({
    sprint: `Sprint ${numero}`,
    comprometidos: metricas.story_points_planificados,
    completados: metricas.story_points_completados,
  }))
}

export function formatoNumero(valor: number | null): string {
  if (valor === null) return 'sin datos'
  return valor.toLocaleString('es-AR', { maximumFractionDigits: 2 })
}

export function textoDesviacion(horas: number | null): string | null {
  if (horas === null) return null
  if (horas > 0) return 'Se trabajó más de lo estimado'
  if (horas < 0) return 'Se trabajó menos de lo estimado'
  return 'Sin desviación: se trabajó lo estimado'
}

export function proporcion(valor: number | null, maximo: number): number {
  if (valor === null || maximo <= 0) return 0
  return Math.min((valor / maximo) * 100, 100)
}
