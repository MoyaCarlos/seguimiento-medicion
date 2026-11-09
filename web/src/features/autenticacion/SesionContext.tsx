import { createContext, useContext, useEffect, useMemo, useRef, useState } from 'react'
import type { ReactNode } from 'react'
import { useLocation, useNavigate } from 'react-router'
import { ApiError } from '../../compartido/api/client'
import { configurarSesionHttp } from '../../compartido/api/sesionHttp'
import * as api from './api'
import type { Usuario } from './api'
import { borrarToken, guardarToken, leerToken } from './almacenamientoToken'

export type EstadoSesion = 'cargando' | 'anonimo' | 'autenticado' | 'sin_conexion'

export type Sesion = {
  estado: EstadoSesion
  usuario: Usuario | null
  ingresar(nombre: string, contrasena: string): Promise<void>
  registrar(nombre: string, contrasena: string): Promise<void>
  salir(): Promise<void>
  reintentar(): Promise<void>
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

  const consultarSesion = useRef<(sigueActivo: () => boolean) => Promise<void>>(async () => {})
  useEffect(() => {
    consultarSesion.current = async (sigueActivo) => {
      try {
        const actual = await api.yo()
        if (!sigueActivo()) return
        setUsuario(actual)
        setEstado('autenticado')
      } catch (error: unknown) {
        if (!sigueActivo()) return
        if (!(error instanceof ApiError && error.status === 401)) setEstado('sin_conexion')
        else if (leerToken()) sesionVencida.current()
        else setEstado('anonimo')
      }
    }
  })

  useEffect(() => {
    if (!leerToken()) return
    let activo = true
    void consultarSesion.current(() => activo)
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
      reintentar: async () => {
        setEstado('cargando')
        await consultarSesion.current(() => true)
      },
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
