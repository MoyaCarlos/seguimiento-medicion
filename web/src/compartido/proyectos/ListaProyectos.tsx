import { useState } from 'react'
import { Link } from 'react-router'
import { ErrorMessage } from '../componentes/ErrorMessage'
import { LoadingIndicator } from '../componentes/LoadingIndicator'
import type { ProyectoResumen, RolEnProyecto } from './api'
import { filtrarProyectos } from './busqueda'
import { etiquetaProyecto } from './ruta'
import { useProyectos } from './useProyectos'

const ETIQUETA_ROL: Record<RolEnProyecto, string> = {
  scrum_master: 'Scrum Master',
  product_builder: 'Product Builder',
}

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
        {proyecto.rol && <span className="tarjeta-proyecto__rol">{ETIQUETA_ROL[proyecto.rol]}</span>}
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
  const [busqueda, setBusqueda] = useState('')

  const total = proyectos.estado === 'listo' ? proyectos.proyectos : []
  const visibles = filtrarProyectos(total, busqueda)
  const hayFiltro = busqueda.trim() !== ''

  return (
    <section className="lista-proyectos" aria-labelledby="proyectos-existentes-title">
      <div className="lista-proyectos__encabezado">
        <h2 id="proyectos-existentes-title">
          Proyectos existentes
          {proyectos.estado === 'listo' && (
            <> <span className="lista-proyectos__contador">
              {hayFiltro ? `${visibles.length} de ${total.length}` : total.length}
            </span></>
          )}
        </h2>
        {total.length > 0 && (
          <input
            className="lista-proyectos__buscador"
            type="search"
            aria-label="Buscar proyecto"
            placeholder="Buscar proyecto…"
            value={busqueda}
            onChange={(evento) => setBusqueda(evento.target.value)}
          />
        )}
      </div>
      {proyectos.estado === 'cargando' && <LoadingIndicator />}
      {proyectos.estado === 'error' && <ErrorMessage message={proyectos.mensaje} />}
      {proyectos.estado === 'listo' && total.length === 0 && (
        <p className="lista-proyectos__vacio">No integrás ningún proyecto. Creá uno con el formulario.</p>
      )}
      {total.length > 0 && visibles.length === 0 && <p className="lista-proyectos__vacio">Ningún proyecto coincide con la búsqueda</p>}
      {visibles.length > 0 && (
        <ul className="lista-proyectos__tarjetas">
          {visibles.map((proyecto) => (
            <li key={proyecto.id}>
              <TarjetaProyecto proyecto={proyecto} />
            </li>
          ))}
        </ul>
      )}
    </section>
  )
}
