import { ApiError } from '../../compartido/api/client'
import type { ConnectionError } from '../../compartido/api/client'
import type { Prioridad, ValorNegocio } from './opciones'

export type FormularioHistoria = {
  titulo: string
  descripcion: string
  prioridad: Prioridad
  valorNegocio: '' | `${ValorNegocio}`
}

export type CuerpoCrearHistoria = {
  proyecto_id: number
  titulo: string
  descripcion: string
  prioridad: Prioridad
  valor_negocio?: ValorNegocio
}

export type HistoriaCreada = {
  id: number
  proyecto_id: number
  titulo: string
  descripcion: string
  prioridad: Prioridad
  estado: string
  valor_negocio: number | null
  estimacion_sp: number | null
}

export type EstadoHistoria = 'Nueva' | 'En progreso' | 'Completada'

export type Historia = {
  id: number
  proyecto_id: number
  titulo: string
  descripcion: string
  criterios_aceptacion: string | null
  prioridad: Prioridad
  estado: EstadoHistoria
  valor_negocio: number | null
  estimacion_sp: number | null
  sprint_id: number | null
}

export type FormularioEdicion = {
  titulo: string
  descripcion: string
  criterios: string
  estimacion: '' | `${ValorNegocio}`
}

export type CuerpoEditarHistoria = {
  titulo: string
  descripcion: string
  criterios_aceptacion: string | null
  estimacion_sp: ValorNegocio | null
}

export type ErrorEdicion = Partial<
  Record<'titulo' | 'descripcion' | 'criterios' | 'estimacion', string>
>

export type ErrorDeCampo = Partial<
  Record<'titulo' | 'descripcion' | 'prioridad' | 'valor_negocio', string>
>

export const ESTADO_INICIAL: FormularioHistoria = {
  titulo: '',
  descripcion: '',
  prioridad: 'M',
  valorNegocio: '',
}

export function parsearProyectoId(valor: string | undefined): number | null {
  if (!valor || !/^\d+$/.test(valor)) return null

  const id = Number(valor)
  return Number.isSafeInteger(id) && id > 0 ? id : null
}

export function armarCuerpo(
  proyectoId: number,
  formulario: FormularioHistoria,
): CuerpoCrearHistoria {
  const cuerpo: CuerpoCrearHistoria = {
    proyecto_id: proyectoId,
    titulo: formulario.titulo.trim(),
    descripcion: formulario.descripcion.trim(),
    prioridad: formulario.prioridad,
  }

  if (formulario.valorNegocio !== '') {
    cuerpo.valor_negocio = Number(formulario.valorNegocio) as ValorNegocio
  }

  return cuerpo
}

const LONGITUD_MAXIMA_TITULO = 200
const LONGITUD_MAXIMA_DESCRIPCION = 2000

export function contarCaracteres(texto: string): number {
  return [...texto].length
}

export function validarHistoria(formulario: FormularioHistoria): ErrorDeCampo {
  const errores: ErrorDeCampo = {}
  const titulo = formulario.titulo.trim()
  const descripcion = formulario.descripcion.trim()

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

  return errores
}

export type ErrorMapeado = {
  campo?: keyof ErrorDeCampo
  mensaje: string
}

const CAMPOS_CONOCIDOS: readonly (keyof ErrorDeCampo)[] = [
  'titulo',
  'descripcion',
  'prioridad',
  'valor_negocio',
]

function esCampoConocido(campo: string): campo is keyof ErrorDeCampo {
  return CAMPOS_CONOCIDOS.includes(campo as keyof ErrorDeCampo)
}

export function mapearError(error: ApiError | ConnectionError): ErrorMapeado {
  if (!(error instanceof ApiError)) {
    return { mensaje: error.message }
  }

  if (error.field && esCampoConocido(error.field)) {
    return { campo: error.field, mensaje: error.message }
  }

  return { mensaje: error.message }
}
