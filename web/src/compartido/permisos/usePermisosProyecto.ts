import type { Rol } from './useRolEnProyecto'

export type PermisosProyecto = {
  rol: Rol | null
  fechaCierre: string | null
  cerrado: boolean
  puedeGestionar: boolean
  puedeEscribir: boolean
  aplicarCierre: (fechaCierre: string | null) => void
}

export function usePermisosProyecto(proyectoId: number): PermisosProyecto {
  void proyectoId
  return {
    rol: null,
    fechaCierre: null,
    cerrado: false,
    puedeGestionar: true,
    puedeEscribir: true,
    aplicarCierre: () => {},
  }
}
