import { Link } from 'react-router'
import { ErrorMessage } from '../componentes/ErrorMessage'
import { LoadingIndicator } from '../componentes/LoadingIndicator'
import { etiquetaProyecto } from './ruta'
import { useProyectos } from './useProyectos'

export function ListaProyectos() {
  const proyectos = useProyectos()

  if (proyectos.estado === 'cargando') return <LoadingIndicator />
  if (proyectos.estado === 'error') return <ErrorMessage message={proyectos.mensaje} />

  return (
    <section aria-labelledby="proyectos-existentes-title">
      <h2 id="proyectos-existentes-title">Proyectos existentes</h2>
      {proyectos.proyectos.length === 0 ? (
        <p>Todavía no hay proyectos. Creá el primero con el formulario.</p>
      ) : (
        <ul>
          {proyectos.proyectos.map((proyecto) => (
            <li key={proyecto.id}>
              <Link to={`/projects/${proyecto.id}/proyecto`}>{etiquetaProyecto(proyecto)}</Link>
            </li>
          ))}
        </ul>
      )}
    </section>
  )
}
