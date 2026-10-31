export const ESTADOS = ['Pendiente', 'Activo', 'Finalizado'] as const

export type Sprint = {
  id: number
  proyecto_id: number
  sprint_goal?: string | null
  fecha_inicio?: string | null
  fecha_fin?: string | null
  estado: (typeof ESTADOS)[number]
}

export function agruparPorEstado(sprints: Sprint[]): Record<Sprint['estado'], Sprint[]> {
  const ordenados = [...sprints].sort((a, b) => a.id - b.id)
  return {
    Pendiente: ordenados.filter((s) => s.estado === 'Pendiente'),
    Activo: ordenados.filter((s) => s.estado === 'Activo'),
    Finalizado: ordenados.filter((s) => s.estado === 'Finalizado'),
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
}

export const ETIQUETAS_PRIORIDAD: Record<HistoriaDeSprint['prioridad'], string> = {
  M: 'Debe tener',
  S: 'Debería tener',
  C: 'Podría tener',
  W: 'No tendrá',
}
