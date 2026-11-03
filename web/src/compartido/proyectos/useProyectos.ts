import type { ProyectoResumen } from './api'

export type EstadoProyectos =
  | { estado: 'cargando' }
  | { estado: 'listo'; proyectos: ProyectoResumen[] }
  | { estado: 'error'; mensaje: string }

export function useProyectos(): EstadoProyectos {
  return { estado: 'cargando' }
}
