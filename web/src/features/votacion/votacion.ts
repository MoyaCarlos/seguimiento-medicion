export type Rol = 'scrum_master' | 'product_builder'

export type Carta = 1 | 2 | 3 | 5 | 8 | 13 | 21

export const CARTAS: Carta[] = []

export type Integrante = { id: number; nombre: string; rol: Rol }

export type VotoDeIntegrante = {
  usuario_id: number
  nombre: string
  rol: Rol
  ya_voto: boolean
  valor: number | null
}

export type Votacion = {
  historia_id: number
  ronda: number
  revelada_en: string | null
  integrantes: VotoDeIntegrante[]
  promedio: { exacto: number; carta_sugerida: number } | null
}

export const esScrumMaster = (_integrante: Integrante | undefined): boolean => false

export const estaRevelada = (_votacion: Votacion): boolean => false

export const faltanPorVotar = (_votacion: Votacion): VotoDeIntegrante[] => []

export const puedeRevelar = (_votacion: Votacion, _actor: Integrante | undefined): boolean => false

export const textoDeVoto = (_voto: VotoDeIntegrante): string => ''
