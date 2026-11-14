import { useEffect, useState } from 'react'
import { ApiError, ConnectionError, get } from '../api/client'

export type Rol = 'scrum_master' | 'product_builder'

export type EstadoRol =
  | { estado: 'cargando' }
  | { estado: 'error'; mensaje: string }
  | { estado: 'no_integrante' }
  | { estado: 'listo'; rol: Rol; fechaCierre: string | null }

const MENSAJE_NO_INTERPRETABLE = 'No se pudo interpretar el rol en el proyecto'

function esRol(valor: unknown): valor is Rol {
  return valor === 'scrum_master' || valor === 'product_builder'
}

function estadoDeError(error: unknown): EstadoRol {
  if (error instanceof ApiError && error.tipo === 'no_integrante') return { estado: 'no_integrante' }
  if (error instanceof ApiError || error instanceof ConnectionError) {
    return { estado: 'error', mensaje: error.message }
  }
  return { estado: 'error', mensaje: MENSAJE_NO_INTERPRETABLE }
}

export function useRolEnProyecto(proyectoId: number): EstadoRol {
  const [consulta, setConsulta] = useState<{ proyectoId: number; resultado: EstadoRol } | null>(null)

  useEffect(() => {
    let vigente = true

    get<{ rol?: unknown; fecha_cierre?: unknown }>(`/projects/${proyectoId}`)
      .then((proyecto) => {
        if (!vigente) return
        const rol = proyecto?.rol
        const fechaCierre = proyecto?.fecha_cierre ?? null
        setConsulta({
          proyectoId,
          resultado:
            esRol(rol) && (fechaCierre === null || typeof fechaCierre === 'string')
              ? { estado: 'listo', rol, fechaCierre }
              : { estado: 'error', mensaje: MENSAJE_NO_INTERPRETABLE },
        })
      })
      .catch((error: unknown) => {
        if (vigente) setConsulta({ proyectoId, resultado: estadoDeError(error) })
      })

    return () => {
      vigente = false
    }
  }, [proyectoId])

  return consulta?.proyectoId === proyectoId ? consulta.resultado : { estado: 'cargando' }
}
