import { get } from '../api/client'

export type RolEnProyecto = 'scrum_master' | 'product_builder'

export type ProyectoResumen = {
  rol?: RolEnProyecto
  id: number
  nombre: string
  descripcion: string
  fecha_inicio: string | null
  fecha_fin: string | null
}

export const listarProyectos = () => get<ProyectoResumen[]>('/projects')
