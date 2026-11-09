import { get, post } from '../../compartido/api/client'

export type Usuario = { id: number; nombre: string }

export type RespuestaSesion = { usuario: Usuario; token: string; vence_en: string }

export function registrar(nombre: string, contrasena: string): Promise<RespuestaSesion> {
  return post<RespuestaSesion>('/auth/registro', { nombre, contrasena })
}

export function ingresar(nombre: string, contrasena: string): Promise<RespuestaSesion> {
  return post<RespuestaSesion>('/auth/ingreso', { nombre, contrasena })
}

export function yo(): Promise<Usuario> {
  return get<Usuario>('/auth/yo')
}

export function salir(): Promise<void> {
  return post<void>('/auth/salir', undefined)
}
