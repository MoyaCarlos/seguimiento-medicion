import { ApiError, ConnectionError } from '../../compartido/api/client'

const VOTACION_NO_INICIADA = 'la historia no tiene votación iniciada'

export function mensajeDeError(error: unknown): string {
  if (error instanceof ApiError || error instanceof ConnectionError) return error.message
  return 'Ocurrió un error inesperado'
}

export const esVotacionNoIniciada = (error: unknown): boolean =>
  error instanceof ApiError && error.status === 404 && error.message === VOTACION_NO_INICIADA
