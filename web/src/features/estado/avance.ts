export type EstadoHistoria = 'Nueva' | 'En progreso' | 'Completada'

export type SprintRef = {
  id: number
  numero: number
}

export type HistoriaAvance = {
  id: number
  titulo: string
  estimacion_sp: number | null
  estado: EstadoHistoria
  sprint: SprintRef | null
}

export type GrupoAvance = {
  estado: EstadoHistoria
  cantidad: number
  story_points: number
  historias: HistoriaAvance[]
}

export type Avance = {
  total: { historias: number; story_points: number }
  grupos: GrupoAvance[]
}

export function textoDeSprint(_sprint: SprintRef | null): string {
  return ''
}

export function textoDeEstimacion(_storyPoints: number | null): string {
  return ''
}

export function esIdValido(_id: string): boolean {
  return false
}
