export type EnganceSesion = {
  obtenerToken: () => string | null
  alNoAutenticado: () => void
}

export function configurarSesionHttp(_enganche: EnganceSesion | null): void {}
