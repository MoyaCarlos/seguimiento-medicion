import { Link } from 'react-router'
import { ErrorMessage } from '../componentes/ErrorMessage'
import { LoadingIndicator } from '../componentes/LoadingIndicator'
import type { ProyectoResumen } from './api'
import { etiquetaProyecto } from './ruta'
import { useProyectos } from './useProyectos'

function TarjetaProyecto({ proyecto }: { proyecto: ProyectoResumen }) {
  const idDetalle = `proyecto-${proyecto.id}-detalle`

  return (
    <Link
      className="tarjeta-proyecto"
      to={`/projects/${proyecto.id}/proyecto`}
      aria-label={etiquetaProyecto(proyecto)}
      aria-describedby={idDetalle}
    >
      <span className="tarjeta-proyecto__cabecera">
        <span className="tarjeta-proyecto__nombre" title={proyecto.nombre}>{proyecto.nombre}</span>
        <span className="tarjeta-proyecto__id">#{proyecto.id}</span>
      </span>
      <span id={idDetalle}>
        <span className="tarjeta-proyecto__descripcion">{proyecto.descripcion || 'Sin descripción'}</span>
        <span className="tarjeta-proyecto__fechas">
          <span>Inicio: {proyecto.fecha_inicio ?? 'Sin definir'}</span>
          <span>Fin: {proyecto.fecha_fin ?? 'Sin definir'}</span>
        </span>
      </span>
    </Link>
  )
}

export function ListaProyectos() {
  const proyectos = useProyectos()

  return (
    <section className="lista-proyectos" aria-labelledby="proyectos-existentes-title">
      <h2 id="proyectos-existentes-title">
        Proyectos existentes
        {proyectos.estado === 'listo' && <> <span className="lista-proyectos__contador">{proyectos.proyectos.length}</span></>}
      </h2>
      {proyectos.estado === 'cargando' && <LoadingIndicator />}
      {proyectos.estado === 'error' && <ErrorMessage message={proyectos.mensaje} />}
      {proyectos.estado === 'listo' && proyectos.proyectos.length === 0 && (
        <p>Todavía no hay proyectos. Creá el primero con el formulario.</p>
      )}
      {proyectos.estado === 'listo' && proyectos.proyectos.length > 0 && (
        <ul className="lista-proyectos__tarjetas">
          {proyectos.proyectos.map((proyecto) => (
            <li key={proyecto.id}>
              <TarjetaProyecto proyecto={proyecto} />
            </li>
          ))}
        </ul>
      )}
    </section>
  )
}
