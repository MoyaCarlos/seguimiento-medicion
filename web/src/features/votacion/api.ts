import { get, post } from '../../compartido/api/client'
import type { HistoriaResumen, Integrante, Rol, Votacion } from './votacion'

type IntegranteDeApi = { id: number; nombre: string; rol: Rol }

// La identidad sale de la sesión (HU-21): ningún pedido lleva el usuario.
export const consultarVotacion = (historiaId: number) => get<Votacion>(`/backlog/${historiaId}/votacion`)

export const iniciarVotacion = (historiaId: number) => post<Votacion>(`/backlog/${historiaId}/votacion`, {})

export const votar = (historiaId: number, valor: number) =>
  post<Votacion>(`/backlog/${historiaId}/votacion/votos`, { valor })

export const revelarVotacion = (historiaId: number) =>
  post<Votacion>(`/backlog/${historiaId}/votacion/revelar`, {})

export const abrirNuevaRonda = (historiaId: number) =>
  post<Votacion>(`/backlog/${historiaId}/votacion/rondas`, {})

export const registrarEstimacion = (historiaId: number, estimacionSp: number) =>
  post<HistoriaResumen>(`/backlog/${historiaId}/votacion/estimacion`, { estimacion_sp: estimacionSp })

export async function listarIntegrantes(proyectoId: string): Promise<Integrante[]> {
  const integrantes = await get<IntegranteDeApi[]>(`/projects/${proyectoId}/members`)
  return integrantes.map(({ id, nombre, rol }) => ({ id, nombre, rol }))
}

// GET /backlog/{id} no existe todavía (issue #222): se busca en el listado del proyecto,
// que solo trae historias sin Sprint. Reemplazar solo esta función cuando exista.
export async function obtenerHistoria(proyectoId: string, historiaId: number): Promise<HistoriaResumen | null> {
  const historias = await get<HistoriaResumen[]>(`/projects/${proyectoId}/backlog`)
  return historias.find((historia) => historia.id === historiaId) ?? null
}
