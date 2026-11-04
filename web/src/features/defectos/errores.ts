import { ApiError } from '../../compartido/api/client'

const CAMPOS = ['descripcion', 'severidad', 'historia_id']

export type ErrorClasificado = { campo: string; mensaje: string } | { mensaje: string }

export function clasificarError(error: unknown): ErrorClasificado {
  if (error instanceof ApiError && error.field && CAMPOS.includes(error.field)) {
    return { campo: error.field, mensaje: error.message }
  }
  return { mensaje: mensajeDeError(error) }
}

export function mensajeDeError(error: unknown): string {
  return error instanceof Error ? error.message : 'Ocurrió un error inesperado'
}
