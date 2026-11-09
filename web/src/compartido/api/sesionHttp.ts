export type EnganceSesion = {
  obtenerToken: () => string | null
  alNoAutenticado: () => void
}

let enganche: EnganceSesion | null = null

export function configurarSesionHttp(nuevo: EnganceSesion | null): void {
  enganche = nuevo
}

export function tokenActual(): string | null {
  return enganche?.obtenerToken() ?? null
}

export function avisarNoAutenticado(): void {
  enganche?.alNoAutenticado()
}
