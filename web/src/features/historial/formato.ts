import { formatoNumero } from '../../compartido/formato/numeros'

export function conUnidad(valor: number | null, unidad: string): string {
  return valor === null ? formatoNumero(valor) : `${formatoNumero(valor)}${unidad}`
}

export function fecha(iso: string): string {
  const [anio, mes, dia] = iso.split('-')
  return `${dia}/${mes}/${anio}`
}

export function periodo(inicio: string, fin: string): string {
  return `${fecha(inicio)} → ${fecha(fin)}`
}
