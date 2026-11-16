const SIN_DEFINIR = 'Sin definir'

export function fechaCorta(iso: string | null): string {
  if (!iso) return SIN_DEFINIR
  const [anio, mes, dia] = iso.slice(0, 10).split('-')
  return `${dia}/${mes}/${anio}`
}

export function periodoCorto(inicio: string | null, fin: string | null): string {
  return `${fechaCorta(inicio)} → ${fechaCorta(fin)}`
}
