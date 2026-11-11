export const SIN_DATOS = 'sin datos'

export function numero(valor: number | null): string {
  if (valor === null) return SIN_DATOS
  return valor.toLocaleString('es-AR', { maximumFractionDigits: 2 })
}

export function conUnidad(valor: number | null, unidad: string): string {
  return valor === null ? SIN_DATOS : `${numero(valor)}${unidad}`
}

export function desviacion(horas: number | null, porcentaje: number | null): string {
  if (horas === null) return SIN_DATOS
  return `${conUnidad(horas, ' h')} (${conUnidad(porcentaje, '%')})`
}

export function fecha(iso: string): string {
  const [anio, mes, dia] = iso.split('-')
  return `${dia}/${mes}/${anio}`
}

export function periodo(inicio: string, fin: string): string {
  return `${fecha(inicio)} → ${fecha(fin)}`
}
