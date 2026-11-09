const CLAVE = 'seguimiento-medicion.token'

let copiaEnMemoria: string | null = null

export function leerToken(): string | null {
  try {
    return localStorage.getItem(CLAVE) ?? copiaEnMemoria
  } catch {
    return copiaEnMemoria
  }
}

export function guardarToken(token: string): void {
  copiaEnMemoria = token
  try {
    localStorage.setItem(CLAVE, token)
  } catch {
    // el token queda solo en memoria
  }
}

export function borrarToken(): void {
  copiaEnMemoria = null
  try {
    localStorage.removeItem(CLAVE)
  } catch {
    // sin almacenamiento no hay nada más que borrar
  }
}
