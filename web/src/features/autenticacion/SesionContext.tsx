import { createContext, useContext, useEffect, useMemo, useRef, useState } from 'react'
import type { ReactNode } from 'react'
import { useLocation, useNavigate } from 'react-router'
import { ApiError } from '../../compartido/api/client'
import { configurarSesionHttp } from '../../compartido/api/sesionHttp'
import * as api from './api'
import type { Usuario } from './api'
import { borrarToken, guardarToken, leerToken } from './almacenamientoToken'

export type EstadoSesion = 'cargando' | 'anonimo' | 'autenticado'

export type Sesion = {
  estado: EstadoSesion
  usuario: Usuario | null
  ingresar(nombre: string, contrasena: string): Promise<void>
  registrar(nombre: string, contrasena: string): Promise<void>
  salir(): Promise<void>
}

const SesionContext = createContext<Sesion | null>(null)

export function SesionProvider({ children }: { children: ReactNode }) {
  const navigate = useNavigate()
  const location = useLocation()
  const [estado, setEstado] = useState<EstadoSesion>(() => (leerToken() ? 'cargando' : 'anonimo'))
  const [usuario, setUsuario] = useState<Usuario | null>(null)

  const navegacion = useRef({ navigate, location })
  useEffect(() => {
    navegacion.current = { navigate, location }
  })

  const sesionVencida = useRef(() => {})
  useEffect(() => {
    sesionVencida.current = () => {
      const { navigate, location } = navegacion.current
      borrarToken()
      setUsuario(null)
      setEstado('anonimo')
      navigate('/ingresar', {
        state: { desde: location.pathname + location.search, vencida: true },
        replace: true,
      })
    }
  })

  useEffect(() => {
    configurarSesionHttp({
      obtenerToken: leerToken,
      alNoAutenticado: () => sesionVencida.current(),
    })
    return () => configurarSesionHttp(null)
  }, [])

  useEffect(() => {
    if (!leerToken()) return
    let activo = true
    api
      .yo()
      .then((actual) => {
        if (!activo) return
        setUsuario(actual)
        setEstado('autenticado')
      })
      .catch((error: unknown) => {
        if (!activo) return
        if (error instanceof ApiError && error.status === 401 && leerToken()) sesionVencida.current()
        else setEstado('anonimo')
      })
    return () => {
      activo = false
    }
  }, [])

  const sesion = useMemo<Sesion>(() => {
    async function iniciar(respuesta: api.RespuestaSesion) {
      guardarToken(respuesta.token)
      setUsuario(respuesta.usuario)
      setEstado('autenticado')
      const desde = (navegacion.current.location.state as { desde?: string } | null)?.desde
      navegacion.current.navigate(desde ?? '/', { replace: true })
    }

    return {
      estado,
      usuario,
      ingresar: async (nombre, contrasena) => iniciar(await api.ingresar(nombre, contrasena)),
      registrar: async (nombre, contrasena) => iniciar(await api.registrar(nombre, contrasena)),
      salir: async () => {
        try {
          await api.salir()
        } catch {
          // si el pedido falla igual se cierra la sesión en este navegador
        }
        borrarToken()
        setUsuario(null)
        setEstado('anonimo')
        navegacion.current.navigate('/ingresar', { replace: true })
      },
    }
  }, [estado, usuario])

  return <SesionContext value={sesion}>{children}</SesionContext>
}

// oxlint-disable-next-line react/only-export-components -- el contrato junta proveedor y hook
export function useSesion(): Sesion {
  const sesion = useContext(SesionContext)
  if (!sesion) throw new Error('useSesion debe usarse dentro de SesionProvider')
  return sesion
}
