import { formatoNumero } from '../../compartido/formato/numeros'

export function conUnidad(valor: number | null, unidad: string): string {
  return valor === null ? formatoNumero(valor) : `${formatoNumero(valor)}${unidad}`
}
