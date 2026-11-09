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

export function textoDeSprint(sprint: SprintRef | null): string {
  return sprint ? `Sprint ${sprint.numero}` : 'Product Backlog'
}

export function textoDeEstimacion(storyPoints: number | null): string {
  return storyPoints === null ? 'sin estimar' : String(storyPoints)
}

export function esIdValido(id: string): boolean {
  return /^[1-9]\d*$/.test(id)
}
