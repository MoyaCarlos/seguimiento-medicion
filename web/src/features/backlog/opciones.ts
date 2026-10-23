export type Prioridad = 'M' | 'S' | 'C' | 'W'

export type ValorNegocio = 1 | 2 | 3 | 5 | 8 | 13 | 21

export const PRIORIDADES: { codigo: Prioridad; etiqueta: string }[] = [
  { codigo: 'M', etiqueta: 'Must have' },
  { codigo: 'S', etiqueta: 'Should have' },
  { codigo: 'C', etiqueta: 'Could have' },
  { codigo: 'W', etiqueta: "Won't have" },
]

export const VALORES_NEGOCIO: ValorNegocio[] = [1, 2, 3, 5, 8, 13, 21]

export function etiquetaPrioridad(codigo: Prioridad): string {
  return PRIORIDADES.find((prioridad) => prioridad.codigo === codigo)?.etiqueta ?? codigo
}
