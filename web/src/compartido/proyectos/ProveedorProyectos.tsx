import type { ReactNode } from 'react'
import { ContextoProyectos, useCargaDeProyectos } from './useProyectos'

export function ProveedorProyectos({ recarga, children }: { recarga: number; children: ReactNode }) {
  const estado = useCargaDeProyectos(true, recarga)

  return <ContextoProyectos.Provider value={estado}>{children}</ContextoProyectos.Provider>
}
