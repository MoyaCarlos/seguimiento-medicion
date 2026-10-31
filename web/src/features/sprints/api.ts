import { get, post } from '../../compartido/api/client'
import type { HistoriaDeSprint, Sprint } from './sprint'

type DatosInicio = { sprint_goal: string; fecha_inicio: string; fecha_fin: string }

export const crearSprint = (proyectoId: number) =>
  post<Sprint>('/sprints', { proyecto_id: proyectoId })

export const iniciarSprint = (id: number, datos: DatosInicio) =>
  post<Sprint>(`/sprints/${id}/iniciar`, datos)

export const listarSprints = (proyectoId: string) =>
  get<Sprint[]>(`/projects/${proyectoId}/sprints`)

export const consultarSprintBacklog = (sprintId: number) => {
  void sprintId
  return Promise.resolve([] as HistoriaDeSprint[])
}
