import { createContext, useContext, useEffect, useState } from 'react'
import { ApiError, ConnectionError } from '../api/client'
import { listarProyectos } from './api'
import type { ProyectoResumen } from './api'

export type EstadoProyectos =
  | { estado: 'cargando' }
  | { estado: 'listo'; proyectos: ProyectoResumen[] }
  | { estado: 'error'; mensaje: string }

const MENSAJE_NO_INTERPRETABLE = 'No se pudo interpretar la lista de proyectos'

function mensajeDeError(error: unknown): string {
  if (error instanceof ApiError && error.status >= 300) return error.message
  if (error instanceof ConnectionError) return error.message
  return MENSAJE_NO_INTERPRETABLE
}

export const ContextoProyectos = createContext<EstadoProyectos | null>(null)

export function useCargaDeProyectos(habilitado: boolean, recarga = 0): EstadoProyectos {
  const [estado, setEstado] = useState<EstadoProyectos>({ estado: 'cargando' })

  useEffect(() => {
    if (!habilitado) return
    let vigente = true

    listarProyectos()
      .then((proyectos) => {
        if (!vigente) return
        setEstado(
          Array.isArray(proyectos)
            ? { estado: 'listo', proyectos }
            : { estado: 'error', mensaje: MENSAJE_NO_INTERPRETABLE },
        )
      })
      .catch((error: unknown) => {
        if (vigente) setEstado({ estado: 'error', mensaje: mensajeDeError(error) })
      })

    return () => {
      vigente = false
    }
  }, [habilitado, recarga])

  return estado
}

export function useProyectos(): EstadoProyectos {
  const compartido = useContext(ContextoProyectos)
  const propio = useCargaDeProyectos(compartido === null)

  return compartido ?? propio
}
