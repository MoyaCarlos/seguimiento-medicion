export type Rol = 'scrum_master' | 'product_builder'

export type Carta = 1 | 2 | 3 | 5 | 8 | 13 | 21

export const CARTAS: Carta[] = [1, 2, 3, 5, 8, 13, 21]

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
