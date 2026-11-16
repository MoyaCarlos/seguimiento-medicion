import { get } from '../../compartido/api/client'
import type { Avance } from './avance'

function tieneTotalNumerico(avance: Partial<Avance>): boolean {
  const total = avance.total
  return typeof total?.historias === 'number' && typeof total.story_points === 'number'
}

function tieneGruposConHistorias(avance: Partial<Avance>): boolean {
  return Array.isArray(avance.grupos) && avance.grupos.every((grupo) => Array.isArray(grupo?.historias))
}

function esAvanceCompleto(avance: Avance | undefined): avance is Avance {
  return !!avance && tieneGruposConHistorias(avance) && tieneTotalNumerico(avance)
}

export async function consultarAvance(proyectoId: string): Promise<Avance> {
  const avance = await get<Avance | undefined>(`/projects/${proyectoId}/avance`)
  if (!esAvanceCompleto(avance)) throw new Error('La respuesta no trae el avance')
  return avance
}

export type EstadoDelProyecto = 'Planificado' | 'En curso' | 'Finalizado'

export const consultarEstadoDelProyecto = (proyectoId: string) =>
  get<{ estado: EstadoDelProyecto }>(`/projects/${proyectoId}/status`)
