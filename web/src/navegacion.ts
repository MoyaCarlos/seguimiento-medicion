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

const SECCION_DE_VOTACION = 'Planning Poker'

export function seccionDeRuta(pathname: string): string | undefined {
  const [, , , tramo, subtramo] = pathname.split('/')
  if (tramo === 'votacion') return SECCION_DE_VOTACION
  if (tramo === 'backlog' && subtramo === 'nueva') return 'Nueva historia'
  const ajuste = ajustes.find(({ ruta }) => ruta === `${tramo}/${subtramo}`)
  if (ajuste) return ajuste.texto
  return navegacion.find(({ ruta }) => ruta === tramo)?.texto
}
