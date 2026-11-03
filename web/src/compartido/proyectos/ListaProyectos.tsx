import { Link } from 'react-router'
import { etiquetaProyecto } from './ruta'
import { useProyectos } from './useProyectos'

export function ListaProyectos() {
  const proyectos = useProyectos()

  if (proyectos.estado !== 'listo') return null

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
