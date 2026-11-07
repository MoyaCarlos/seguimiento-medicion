import { get } from '../../compartido/api/client'
import type { Avance } from './avance'

function tieneTotalNumerico(avance: Partial<Avance> | undefined): boolean {
  const total = avance?.total
  return typeof total?.historias === 'number' && typeof total.story_points === 'number'
}

export async function consultarAvance(proyectoId: string): Promise<Avance> {
  const avance = await get<Avance | undefined>(`/projects/${proyectoId}/avance`)
  if (!Array.isArray(avance?.grupos) || !tieneTotalNumerico(avance)) throw new Error('La respuesta no trae el avance')
  return avance
}
