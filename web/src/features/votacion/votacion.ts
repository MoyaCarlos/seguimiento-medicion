export type Rol = 'scrum_master' | 'product_builder'

export const ETIQUETA_ROL: Record<Rol, string> = {
  scrum_master: 'Scrum Master',
  product_builder: 'Product Builder',
}

export type Carta = 1 | 2 | 3 | 5 | 8 | 13 | 21

export const CARTAS: Carta[] = [1, 2, 3, 5, 8, 13, 21]

export type Prioridad = 'M' | 'S' | 'C' | 'W'

export type HistoriaResumen = {
  id: number
  titulo: string
  descripcion: string
  prioridad: Prioridad
  estimacion_sp: number | null
}

export type Integrante = { id: number; nombre: string; rol: Rol }

export type VotoDeIntegrante = {
  usuario_id: number
  nombre: string
  rol: Rol
  ya_voto: boolean
  valor: number | null
}

export type Diferencias = {
  hay_diferencias: boolean
  minimo: number | null
  maximo: number | null
}

export type Votacion = {
  historia_id: number
  ronda: number
  revelada_en: string | null
  integrantes: VotoDeIntegrante[]
  promedio: { exacto: number; carta_sugerida: number } | null
  diferencias: Diferencias | null
}

export const esScrumMaster = (integrante: Integrante | undefined): boolean =>
  integrante?.rol === 'scrum_master'

export const estaRevelada = (votacion: Votacion): boolean => votacion.revelada_en !== null

export const faltanPorVotar = (votacion: Votacion): VotoDeIntegrante[] =>
  votacion.integrantes.filter((integrante) => !integrante.ya_voto)

export const puedeRevelar = (votacion: Votacion, actor: Integrante | undefined): boolean =>
  esScrumMaster(actor) && !estaRevelada(votacion) && faltanPorVotar(votacion).length === 0

export function textoDeVoto(voto: VotoDeIntegrante): string {
  if (!voto.ya_voto) return 'Falta'
  return voto.valor === null ? 'Oculto' : String(voto.valor)
}

export function textoDeDiferencias(diferencias: Diferencias | null): string {
  if (diferencias === null) return ''
  if (!diferencias.hay_diferencias) return 'No hay diferencias: todos votaron lo mismo'
  return `Hay diferencias: mínimo ${diferencias.minimo}, máximo ${diferencias.maximo}`
}

export const sinValores = (votacion: Votacion): Votacion => ({
  ...votacion,
  integrantes: votacion.integrantes.map((integrante) => ({ ...integrante, valor: null })),
})
