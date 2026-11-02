import type { Integrante, Votacion } from './votacion'

export const consultarVotacion = (_historiaId: number, _usuarioId?: number): Promise<Votacion> =>
  Promise.reject(new Error('sin implementar'))

export const iniciarVotacion = (_historiaId: number, _usuarioId: number): Promise<Votacion> =>
  Promise.reject(new Error('sin implementar'))

export const votar = (_historiaId: number, _usuarioId: number, _valor: number): Promise<Votacion> =>
  Promise.reject(new Error('sin implementar'))

export const revelarVotacion = (_historiaId: number, _usuarioId: number): Promise<Votacion> =>
  Promise.reject(new Error('sin implementar'))

export const listarIntegrantes = (_proyectoId: string): Promise<Integrante[]> =>
  Promise.reject(new Error('sin implementar'))
