export type MetricasSprint = {
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

export type SprintHistorial = {
  id: number
  proyecto_id: number
  numero: number
  sprint_goal: string
  fecha_inicio: string
  fecha_fin: string
  metricas: MetricasSprint
}

export type HistoriaComprometida = {
  id: number
  titulo: string
  estimacion_sp: number | null
  completada: boolean
  minutos: number
  horas: number
}

export type DefectoDeSprint = {
  id: number
  historia_id: number
  historia: string
  descripcion: string
  severidad: 'Alta' | 'Media' | 'Baja'
  estado: 'Pendiente' | 'Resuelto'
  detectado_en_sprint: boolean
  resuelto_en_sprint: boolean
}

export type DetalleSprint = SprintHistorial & {
  historias: HistoriaComprometida[]
  defectos: DefectoDeSprint[]
}
