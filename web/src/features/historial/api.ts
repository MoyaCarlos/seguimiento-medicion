import type { DetalleSprint, SprintHistorial } from './tipos'

export async function listarHistorial(_proyectoId: string): Promise<SprintHistorial[]> {
  throw new Error('sin implementar')
}

export async function obtenerDetalleSprint(_proyectoId: string, _sprintId: string): Promise<DetalleSprint> {
  throw new Error('sin implementar')
}
