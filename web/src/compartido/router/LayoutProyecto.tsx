import { Outlet, useParams } from 'react-router'
import { NoIntegrante } from '../permisos/NoIntegrante'
import { useRolEnProyecto } from '../permisos/useRolEnProyecto'

/**
 * Las pantallas de un proyecto leen sus datos del servidor, que responde 403 a quien no lo
 * integra. Una sola comprobación acá reemplaza todas por «No sos integrante de este
 * proyecto»; mientras no se sabe, o ante cualquier otro resultado, decide el servidor.
 */
export function LayoutProyecto() {
  const { id } = useParams()
  const rol = useRolEnProyecto(Number(id))

  if (rol.estado === 'no_integrante') return <NoIntegrante />
  return <Outlet key={id} />
}
