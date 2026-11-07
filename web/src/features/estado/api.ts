import { get } from '../../compartido/api/client'
import type { Avance } from './avance'

export async function consultarAvance(proyectoId: string): Promise<Avance> {
  const avance = await get<Avance | undefined>(`/projects/${proyectoId}/avance`)
  if (!Array.isArray(avance?.grupos)) throw new Error('La respuesta no trae el avance')
  return avance
}
