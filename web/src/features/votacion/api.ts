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

export const abrirNuevaRonda = (_historiaId: number, _usuarioId: number): Promise<Votacion> =>
  Promise.reject(new Error('sin implementar'))

export const registrarEstimacion = (
  _historiaId: number,
  _usuarioId: number,
  _estimacionSp: number,
): Promise<HistoriaResumen> => Promise.reject(new Error('sin implementar'))

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
