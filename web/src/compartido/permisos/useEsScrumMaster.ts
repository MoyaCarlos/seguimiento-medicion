import { usePermisosProyecto } from './usePermisosProyecto'

/**
 * Devuelve false SOLO si se sabe que la persona es Product Builder, que no integra el
 * proyecto o que el proyecto está cerrado. Mientras no se conoce el rol (cargando o error)
 * devuelve true: decide el 403 del servidor.
 */
export function useEsScrumMaster(proyectoId: number): boolean {
  return usePermisosProyecto(proyectoId).puedeGestionar
}
