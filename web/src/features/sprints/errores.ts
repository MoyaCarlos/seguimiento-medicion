import { ApiError } from '../../compartido/api/client'

export type ErrorClasificado = { campo: string; mensaje: string } | { aviso: string }

export function clasificarError(error: unknown, campos: readonly string[]): ErrorClasificado {
  if (error instanceof ApiError && error.field && campos.includes(error.field)) {
    return { campo: error.field, mensaje: error.message }
  }
  if (error instanceof Error) return { aviso: error.message }
  return { aviso: 'Ocurrió un error inesperado' }
}
