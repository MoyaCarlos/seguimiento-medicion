import { post } from '../../compartido/api/client'
import type { CuerpoCrearHistoria, HistoriaCreada } from './historia'

export function crearHistoria(cuerpo: CuerpoCrearHistoria): Promise<HistoriaCreada> {
  return post<HistoriaCreada>('/backlog', cuerpo)
}
