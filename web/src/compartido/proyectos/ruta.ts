import { navegacion } from '../../navegacion'

const SECCION_POR_DEFECTO = 'proyecto'

export function rutaDeSeccion(pathname: string, idNuevo: string | number): string {
  const [, , , primerTramo] = pathname.toLowerCase().split('/')
  const esDelMenu = navegacion.some(({ ruta }) => ruta === primerTramo)
  return `/projects/${idNuevo}/${esDelMenu ? primerTramo : SECCION_POR_DEFECTO}`
}

export function etiquetaProyecto({ id, nombre }: { id: number; nombre: string }): string {
  return `${nombre} (#${id})`
}
