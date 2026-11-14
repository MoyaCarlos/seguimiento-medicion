import { useRolEnProyecto } from './useRolEnProyecto'

/**
 * Devuelve false SOLO si se sabe que la persona es Product Builder o que no integra el
 * proyecto. Mientras no se conoce el rol (cargando o error) devuelve true: decide el 403
 * del servidor.
 */
export function useEsScrumMaster(proyectoId: number): boolean {
  const rol = useRolEnProyecto(proyectoId)
  if (rol.estado === 'no_integrante') return false
  if (rol.estado === 'listo') return rol.rol === 'scrum_master'
  return true
}
