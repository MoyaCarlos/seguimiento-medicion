import { get } from '../api/client'

export type RolEnProyecto = 'scrum_master' | 'product_builder'

export type ProyectoResumen = {
  rol?: RolEnProyecto
  fecha_cierre?: string | null
  id: number
  nombre: string
  descripcion: string
  fecha_inicio: string | null
  fecha_fin: string | null
}

export const ETIQUETA_ROL: Record<RolEnProyecto, string> = {
  scrum_master: 'Scrum Master',
  product_builder: 'Product Builder',
}

export const listarProyectos = () => get<ProyectoResumen[]>('/projects')
