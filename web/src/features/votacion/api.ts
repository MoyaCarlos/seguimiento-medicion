import { get, post } from '../../compartido/api/client'
import type { HistoriaResumen, Integrante, Rol, Votacion } from './votacion'

type IntegranteDeApi = { id: number; nombre: string; rol: Rol }

export const consultarVotacion = (historiaId: number, usuarioId?: number) => {
  const consulta = usuarioId === undefined ? '' : `?usuario_id=${usuarioId}`
  return get<Votacion>(`/backlog/${historiaId}/votacion${consulta}`)
}

export const iniciarVotacion = (historiaId: number, usuarioId: number) =>
  post<Votacion>(`/backlog/${historiaId}/votacion`, { usuario_id: usuarioId })

export const votar = (historiaId: number, usuarioId: number, valor: number) =>
  post<Votacion>(`/backlog/${historiaId}/votacion/votos`, { usuario_id: usuarioId, valor })

export const revelarVotacion = (historiaId: number, usuarioId: number) =>
  post<Votacion>(`/backlog/${historiaId}/votacion/revelar`, { usuario_id: usuarioId })

export async function listarIntegrantes(proyectoId: string): Promise<Integrante[]> {
  const integrantes = await get<IntegranteDeApi[]>(`/projects/${proyectoId}/members`)
  return integrantes.map(({ id, nombre, rol }) => ({ id, nombre, rol }))
}

export const obtenerHistoria = (_proyectoId: string, _historiaId: number): Promise<HistoriaResumen | null> =>
  Promise.resolve(null)
