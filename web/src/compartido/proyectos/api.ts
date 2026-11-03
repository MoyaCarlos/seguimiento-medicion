import { get } from '../api/client'

export type ProyectoResumen = {
  id: number
  nombre: string
}

export const listarProyectos = () => get<ProyectoResumen[]>('/projects')
