import { ApiError, ConnectionError } from '../../compartido/api/client'

export function mensajeDeError(error: unknown): string {
  if (error instanceof ApiError || error instanceof ConnectionError) {
    return error.message
  }
  return 'Ocurrió un error inesperado'
}

export function campoDeError(error: unknown): string | undefined {
  return error instanceof ApiError ? error.field : undefined
}
