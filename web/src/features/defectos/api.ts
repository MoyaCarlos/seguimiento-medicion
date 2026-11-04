import { get, post } from '../../compartido/api/client'

export type Severidad = 'Alta' | 'Media' | 'Baja'
export type EstadoDefecto = 'Pendiente' | 'Resuelto'
export type EstadoSprint = 'Pendiente' | 'Activo' | 'Finalizado'

export type Defecto = {
  id: number
  proyecto_id: number
  historia_id: number
  historia: string
  descripcion: string
  severidad: Severidad
  estado: EstadoDefecto
  sprint_deteccion_id: number
  sprint_resolucion_id: number | null
  creado_en: string
}

export type SprintResumen = {
  id: number
  numero: number
  estado: EstadoSprint
}

export type HistoriaDeProyecto = {
  id: number
  titulo: string
  origen: string
}

export type NuevoDefecto = {
  historia_id: number
  descripcion: string
  severidad: Severidad
}

type HistoriaResumen = {
  id: number
  titulo: string
}

export const listarDefectos = (proyectoId: string) =>
  get<Defecto[]>(`/projects/${proyectoId}/defectos`)

export const registrarDefecto = (proyectoId: string, cuerpo: NuevoDefecto) =>
  post<Defecto>(`/projects/${proyectoId}/defectos`, cuerpo)

export const resolverDefecto = (id: number) =>
  post<Defecto>(`/defectos/${id}/resolver`, undefined)

export const listarSprints = (proyectoId: string) =>
  get<SprintResumen[]>(`/projects/${proyectoId}/sprints`)

export async function listarHistoriasDelProyecto(
  proyectoId: string,
  sprints: SprintResumen[],
): Promise<HistoriaDeProyecto[]> {
  const [backlog, ...porSprint] = await Promise.all([
    get<HistoriaResumen[]>(`/projects/${proyectoId}/backlog`),
    ...sprints.map((sprint) => get<HistoriaResumen[]>(`/sprints/${sprint.id}/backlog`)),
  ])

  return [
    ...backlog.map((historia) => ({ ...historia, origen: 'Product Backlog' })),
    ...porSprint.flatMap((historias, indice) =>
      historias.map((historia) => ({ ...historia, origen: `Sprint ${sprints[indice].numero}` })),
    ),
  ]
}
