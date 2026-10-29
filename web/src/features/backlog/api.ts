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
