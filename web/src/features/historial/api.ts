import { get } from '../../compartido/api/client'
import type { DetalleSprint, SprintHistorial } from './tipos'

export async function listarHistorial(proyectoId: string): Promise<SprintHistorial[]> {
  const sprints = await get<unknown>(`/projects/${proyectoId}/historial`)
  if (!Array.isArray(sprints)) throw new Error('La respuesta no trae el historial')
  return sprints as SprintHistorial[]
}

export async function obtenerDetalleSprint(proyectoId: string, sprintId: string): Promise<DetalleSprint> {
  const detalle = await get<Partial<DetalleSprint> | undefined>(`/projects/${proyectoId}/historial/${sprintId}`)
  if (!Array.isArray(detalle?.historias) || !Array.isArray(detalle.defectos)) {
    throw new Error('La respuesta no trae el detalle del Sprint')
  }
  return detalle as DetalleSprint
}
