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
