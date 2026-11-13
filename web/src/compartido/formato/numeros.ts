export function formatoNumero(valor: number | null): string {
  if (valor === null) return 'sin datos'
  const texto = valor.toLocaleString('es-AR', { maximumFractionDigits: 1 })
  return texto === '-0' ? '0' : texto
}

function conSigno(valor: number): string {
  const texto = formatoNumero(valor)
  return valor > 0 && texto !== '0' ? `+${texto}` : texto
}

export function textoDeDesviacion(horas: number | null, porcentaje: number | null): string {
  if (horas === null) return 'sin datos'
  const relativa = porcentaje === null ? 'sin datos' : `${conSigno(porcentaje)} %`
  return `${conSigno(horas)} h (${relativa})`
}
