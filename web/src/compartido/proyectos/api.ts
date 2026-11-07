import { get } from '../api/client'

export type ProyectoResumen = {
  id: number
  nombre: string
  descripcion: string
  fecha_inicio: string | null
  fecha_fin: string | null
}

export const listarProyectos = () => get<ProyectoResumen[]>('/projects')
