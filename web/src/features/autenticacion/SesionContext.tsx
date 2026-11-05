import type { ReactNode } from 'react'
import type { Usuario } from './api'

export type EstadoSesion = 'cargando' | 'anonimo' | 'autenticado'

export type Sesion = {
  estado: EstadoSesion
  usuario: Usuario | null
  ingresar(nombre: string, contrasena: string): Promise<void>
  registrar(nombre: string, contrasena: string): Promise<void>
  salir(): Promise<void>
}

export function SesionProvider({ children }: { children: ReactNode }) {
  return <>{children}</>
}

export function useSesion(): Sesion {
  return {
    estado: 'anonimo',
    usuario: null,
    ingresar: async () => {},
    registrar: async () => {},
    salir: async () => {},
  }
}
