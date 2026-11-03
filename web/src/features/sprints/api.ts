import { get, post, put } from '../../compartido/api/client'
import type { DatosImputacion, EsfuerzoReal, IntegranteDeProyecto } from './horas'
import type { HistoriaDeSprint, Sprint } from './sprint'

type DatosInicio = { sprint_goal: string; fecha_inicio: string; fecha_fin: string }

export const crearSprint = (proyectoId: number) =>
  post<Sprint>('/sprints', { proyecto_id: proyectoId })

export const iniciarSprint = (id: number, datos: DatosInicio) =>
  post<Sprint>(`/sprints/${id}/iniciar`, datos)

export const listarSprints = (proyectoId: string) =>
  get<Sprint[]>(`/projects/${proyectoId}/sprints`)

export const consultarSprintBacklog = (sprintId: number) =>
  get<HistoriaDeSprint[]>(`/sprints/${sprintId}/backlog`)

export const cambiarEstadoHistoria = (historiaId: number, estado: 'En progreso' | 'Completada') =>
  put<HistoriaDeSprint>(`/backlog/${historiaId}/estado`, { estado })

export const consultarEsfuerzo = (historiaId: number) =>
  get<EsfuerzoReal>(`/backlog/${historiaId}/horas`)

export const imputarHoras = (historiaId: number, datos: DatosImputacion) =>
  post<unknown>(`/backlog/${historiaId}/horas`, datos)

export const listarIntegrantes = (proyectoId: number) =>
  get<IntegranteDeProyecto[]>(`/projects/${proyectoId}/members`)
