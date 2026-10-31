import { get, post, put } from '../../compartido/api/client'
import type { CuerpoCrearHistoria, CuerpoEditarHistoria, Historia, HistoriaCreada } from './historia'

export function crearHistoria(cuerpo: CuerpoCrearHistoria): Promise<HistoriaCreada> {
  return post<HistoriaCreada>('/backlog', cuerpo)
}

export function listarBacklog(proyectoId: string): Promise<Historia[]> {
  return get<Historia[]>(`/projects/${proyectoId}/backlog`)
}

export function editarHistoria(
  id: number,
  cuerpo: CuerpoEditarHistoria,
): Promise<Historia> {
  return put<Historia>(`/backlog/${id}`, cuerpo)
}

export type SprintParaAsignar = {
  id: number
  estado: 'Pendiente' | 'Activo' | 'Finalizado'
  sprint_goal?: string | null
}

export function asignarASprint(historiaId: number, sprintId: number): Promise<Historia> {
  void historiaId
  void sprintId
  return Promise.resolve(undefined as unknown as Historia)
}

export function listarSprintsDelProyecto(proyectoId: string): Promise<SprintParaAsignar[]> {
  void proyectoId
  return Promise.resolve([])
}
