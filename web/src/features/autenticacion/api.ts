export type Usuario = { id: number; nombre: string }

export type RespuestaSesion = { usuario: Usuario; token: string; vence_en: string }

export async function registrar(_nombre: string, _contrasena: string): Promise<RespuestaSesion> {
  return undefined as never
}

export async function ingresar(_nombre: string, _contrasena: string): Promise<RespuestaSesion> {
  return undefined as never
}

export async function yo(): Promise<Usuario> {
  return undefined as never
}

export async function salir(): Promise<void> {}
