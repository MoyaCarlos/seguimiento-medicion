import { ApiError, ConnectionError } from '../../compartido/api/client'
import { contarCaracteres } from './historia'
import type {
  CuerpoEditarHistoria,
  ErrorEdicion,
  FormularioEdicion,
  Historia,
} from './historia'
import { STORY_POINTS } from './opciones'

const LONGITUD_MAXIMA_TITULO = 200
const LONGITUD_MAXIMA_DESCRIPCION = 2000
const LONGITUD_MAXIMA_CRITERIOS = 2000

export function formularioDe(historia: Historia): FormularioEdicion {
  const estimacion = historia.estimacion_sp === null ? '' : `${historia.estimacion_sp}`
  return {
    titulo: historia.titulo,
    descripcion: historia.descripcion,
    criterios: historia.criterios_aceptacion ?? '',
    estimacion: estimacion as FormularioEdicion['estimacion'],
  }
}

export function validarEdicion(formulario: FormularioEdicion): ErrorEdicion {
  const errores: ErrorEdicion = {}
  const titulo = formulario.titulo.trim()
  const descripcion = formulario.descripcion.trim()
  const criterios = formulario.criterios.trim()

  if (!titulo) {
    errores.titulo = 'El título es obligatorio'
  } else if (contarCaracteres(titulo) > LONGITUD_MAXIMA_TITULO) {
    errores.titulo = 'El título no puede superar los 200 caracteres'
  }

  if (!descripcion) {
    errores.descripcion = 'La descripción es obligatoria'
  } else if (contarCaracteres(descripcion) > LONGITUD_MAXIMA_DESCRIPCION) {
    errores.descripcion = 'La descripción no puede superar los 2000 caracteres'
  }

  if (contarCaracteres(criterios) > LONGITUD_MAXIMA_CRITERIOS) {
    errores.criterios = 'Los criterios de aceptación no pueden superar los 2000 caracteres'
  }

  if (formulario.estimacion !== '' && !esStoryPoint(formulario.estimacion)) {
    errores.estimacion = 'Los Story Points deben ser uno de 1, 2, 3, 5, 8, 13, 21'
  }

  return errores
}

export function armarCuerpoEdicion(formulario: FormularioEdicion): CuerpoEditarHistoria {
  const criterios = formulario.criterios.trim()
  return {
    titulo: formulario.titulo.trim(),
    descripcion: formulario.descripcion.trim(),
    criterios_aceptacion: criterios === '' ? null : criterios,
    estimacion_sp:
      formulario.estimacion === '' ? null : (Number(formulario.estimacion) as CuerpoEditarHistoria['estimacion_sp']),
  }
}

export function campoDeEdicion(campo: string): keyof ErrorEdicion | undefined {
  if (campo === 'criterios_aceptacion') return 'criterios'
  if (campo === 'estimacion_sp') return 'estimacion'
  if (campo === 'titulo' || campo === 'descripcion') return campo
  return undefined
}

export function mensajeDeError(error: unknown): string {
  if (error instanceof ApiError || error instanceof ConnectionError) return error.message
  return 'Ocurrió un error inesperado'
}

export function campoDeError(error: unknown): string | undefined {
  return error instanceof ApiError ? error.field : undefined
}

function esStoryPoint(valor: string): boolean {
  return STORY_POINTS.some((storyPoint) => `${storyPoint}` === valor)
}
