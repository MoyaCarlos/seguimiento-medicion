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
