export const ESTADOS = ['Pendiente', 'Activo', 'Finalizado'] as const

export type Sprint = {
  id: number
  proyecto_id: number
  sprint_goal?: string | null
  fecha_inicio?: string | null
  fecha_fin?: string | null
  estado: (typeof ESTADOS)[number]
}
