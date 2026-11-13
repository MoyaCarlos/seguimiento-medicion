export type Rol = 'scrum_master' | 'product_builder'

export type EstadoRol =
  | { estado: 'cargando' }
  | { estado: 'error'; mensaje: string }
  | { estado: 'no_integrante' }
  | { estado: 'listo'; rol: Rol }

export function useRolEnProyecto(_proyectoId: number): EstadoRol {
  return { estado: 'cargando' }
}
