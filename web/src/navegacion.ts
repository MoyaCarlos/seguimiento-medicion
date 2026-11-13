import type { NombreIcono } from './compartido/componentes/Icono'

export const navegacion: { ruta: string; texto: string; icono: NombreIcono }[] = [
  { ruta: 'proyecto', texto: 'Proyecto', icono: 'panel' },
  { ruta: 'backlog', texto: 'Historias', icono: 'lista' },
  { ruta: 'sprints', texto: 'Sprints', icono: 'bandera' },
  { ruta: 'defectos', texto: 'Defectos', icono: 'defecto' },
  { ruta: 'status', texto: 'Estado', icono: 'estado' },
  { ruta: 'dashboard', texto: 'Dashboard', icono: 'barras' },
  { ruta: 'historial', texto: 'Historial de Sprints', icono: 'historial' },
]

export const ajustes: { ruta: string; texto: string; icono: NombreIcono }[] = [
  { ruta: 'proyecto/settings', texto: 'Equipo', icono: 'equipo' },
  { ruta: 'proyecto/edit', texto: 'Editar proyecto', icono: 'editar' },
]
