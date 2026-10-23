import { get, post, put } from '../../compartido/api/client'
import type { DatosProyecto } from './validacion'

export type Rol = 'scrum_master' | 'product_builder'

export type Proyecto = {
  id: number
  nombre: string
  descripcion: string
  fecha_inicio: string | null
  fecha_fin: string | null
  creado_en: string
}

export type Integrante = {
  id: number
  proyecto_id: number
  nombre: string
  rol: Rol
  vinculado_en: string
}

export type NuevoProyecto = {
  nombre: string
  descripcion: string
  creador: string
  fecha_inicio: string | null
  fecha_fin: string | null
}

export type ProyectoEditado = {
  nombre: string
  descripcion: string
  fecha_inicio: string | null
  fecha_fin: string | null
}

export type NuevoIntegrante = {
  nombre: string
  rol: Rol
}

export const crearProyecto = (datos: NuevoProyecto) =>
  post<Proyecto>('/projects', datos)

export const obtenerProyecto = (id: string) =>
  get<Proyecto>(`/projects/${id}`)

export const editarProyecto = (id: string, datos: ProyectoEditado) =>
  put<Proyecto>(`/projects/${id}`, datos)

export const listarIntegrantes = (id: string) =>
  get<Integrante[]>(`/projects/${id}/members`)

export const asignarIntegrante = (id: string, datos: NuevoIntegrante) =>
  post<Integrante>(`/projects/${id}/members`, datos)

export const aNuevoProyecto = (datos: DatosProyecto): NuevoProyecto => ({
  nombre: datos.nombre.trim(),
  descripcion: datos.descripcion.trim(),
  creador: datos.creador.trim(),
  fecha_inicio: datos.fecha_inicio || null,
  fecha_fin: datos.fecha_fin || null,
})

export const aProyectoEditado = (datos: DatosProyecto): ProyectoEditado => ({
  nombre: datos.nombre.trim(),
  descripcion: datos.descripcion.trim(),
  fecha_inicio: datos.fecha_inicio || null,
  fecha_fin: datos.fecha_fin || null,
})

export const ETIQUETAS_ROL: Record<Rol, string> = {
  scrum_master: 'Scrum Master',
  product_builder: 'Product Builder',
}
