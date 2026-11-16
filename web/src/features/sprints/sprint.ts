export const ESTADOS = ['Pendiente', 'Activo', 'Finalizado'] as const

export type Sprint = {
  id: number
  numero: number
  proyecto_id: number
  sprint_goal?: string | null
  fecha_inicio?: string | null
  fecha_fin?: string | null
  estado: (typeof ESTADOS)[number]
}

export function agruparPorEstado(sprints: Sprint[]): Record<Sprint['estado'], Sprint[]> {
  const ascendentes = [...sprints].sort((a, b) => a.numero - b.numero)
  return {
    Pendiente: ascendentes.filter((s) => s.estado === 'Pendiente'),
    Activo: ascendentes.filter((s) => s.estado === 'Activo'),
    Finalizado: ascendentes.filter((s) => s.estado === 'Finalizado').reverse(),
  }
}

export function guardarSprint(sprints: Sprint[], sprint: Sprint): Sprint[] {
  const existe = sprints.some((s) => s.id === sprint.id)
  return existe ? sprints.map((s) => (s.id === sprint.id ? sprint : s)) : [...sprints, sprint]
}

export type HistoriaDeSprint = {
  id: number
  titulo: string
  prioridad: 'M' | 'S' | 'C' | 'W'
  estado: 'Nueva' | 'En progreso' | 'Completada'
  estimacion_sp: number | null
  fecha_finalizacion?: string | null
}

export const ETIQUETAS_PRIORIDAD: Record<HistoriaDeSprint['prioridad'], string> = {
  M: 'Must have',
  S: 'Should have',
  C: 'Could have',
  W: "Won't have",
}
