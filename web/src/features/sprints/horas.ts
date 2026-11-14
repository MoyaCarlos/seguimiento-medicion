export type RegistroHoras = {
  id: number
  integrante: string
  fecha: string
  actividad: string
  horas: number
}

export type EsfuerzoReal = {
  total_horas: number
  registros: RegistroHoras[]
}

export type IntegranteDeProyecto = {
  id: number
  nombre: string
}

export type DatosImputacion = {
  fecha: string
  actividad: string
  horas: string
}

export const normalizarHoras = (texto: string) => texto.trim().replace(',', '.')

export const formatearHoras = (horas: number) =>
  `${horas.toLocaleString('es-AR', { maximumFractionDigits: 2 })} h`

export function fechaPropuesta(ahora: Date, fechaFin?: string | null): string {
  const hoy = fechaLocal(ahora)
  return fechaFin && fechaFin < hoy ? fechaFin : hoy
}

function fechaLocal(fecha: Date): string {
  const mes = String(fecha.getMonth() + 1).padStart(2, '0')
  const dia = String(fecha.getDate()).padStart(2, '0')
  return `${fecha.getFullYear()}-${mes}-${dia}`
}
