import { useCallback, useState } from 'react'
import { useRolEnProyecto, type Rol } from './useRolEnProyecto'

export type PermisosProyecto = {
  rol: Rol | null
  fechaCierre: string | null
  cerrado: boolean
  puedeGestionar: boolean
  puedeEscribir: boolean
  aplicarCierre: (fechaCierre: string | null) => void
}

/**
 * Un solo GET /projects/{id}. Mientras el rol o el cierre no se conocen (cargando o error)
 * todo queda habilitado: decide el servidor.
 */
export function usePermisosProyecto(proyectoId: number): PermisosProyecto {
  const estadoRol = useRolEnProyecto(proyectoId)
  const [reemplazo, setReemplazo] = useState<{ proyectoId: number; fechaCierre: string | null } | null>(null)

  const aplicarCierre = useCallback(
    (fechaCierre: string | null) => setReemplazo({ proyectoId, fechaCierre }),
    [proyectoId],
  )

  const listo = estadoRol.estado === 'listo'
  const rol = listo ? estadoRol.rol : null
  const reemplazoVigente = reemplazo?.proyectoId === proyectoId ? reemplazo : null
  const fechaCierre = listo ? (reemplazoVigente ? reemplazoVigente.fechaCierre : estadoRol.fechaCierre) : null
  const cerrado = fechaCierre !== null

  return {
    rol,
    fechaCierre,
    cerrado,
    puedeGestionar: !cerrado && rol !== 'product_builder' && estadoRol.estado !== 'no_integrante',
    puedeEscribir: !cerrado,
    aplicarCierre,
  }
}
