import { ETIQUETA_ROL } from './api'
import { useProyectos } from './useProyectos'

/** Rol de la persona en el proyecto de la ruta, sacado de la lista compartida (sin pedidos nuevos). */
export function useRolEnRuta(proyectoId?: string): string | undefined {
  const proyectos = useProyectos()
  if (!proyectoId || proyectos.estado !== 'listo') return undefined
  const rol = proyectos.proyectos.find((proyecto) => String(proyecto.id) === proyectoId)?.rol
  return rol ? ETIQUETA_ROL[rol] : undefined
}
