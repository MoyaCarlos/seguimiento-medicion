import { useState } from 'react'
import { Link, useNavigate } from 'react-router'
import { AvatarProyecto } from '../componentes/AvatarProyecto'
import { Button } from '../componentes/Button'
import { ErrorMessage } from '../componentes/ErrorMessage'
import { Icono } from '../componentes/Icono'
import { LoadingIndicator } from '../componentes/LoadingIndicator'
import { fechaCorta, periodoCorto } from '../formato/fechas'
import { ETIQUETA_ROL } from './api'
import type { ProyectoResumen } from './api'
import { filtrarProyectos } from './busqueda'
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
        <AvatarProyecto id={proyecto.id} nombre={proyecto.nombre} />
        {(proyecto.rol || proyecto.fecha_cierre) && (
          <span className="tarjeta-proyecto__etiquetas">
            {proyecto.rol && <span className="tarjeta-proyecto__rol">{ETIQUETA_ROL[proyecto.rol]}</span>}
            {proyecto.fecha_cierre && (
              <span className="tarjeta-proyecto__cierre">{`Cerrado el ${fechaCorta(proyecto.fecha_cierre)}`}</span>
            )}
          </span>
        )}
      </span>
      <span className="tarjeta-proyecto__nombre" title={proyecto.nombre}>{proyecto.nombre}</span>
      <span id={idDetalle} className="tarjeta-proyecto__detalle">
        <span className="tarjeta-proyecto__descripcion">{proyecto.descripcion || 'Sin descripción'}</span>
        <span className="tarjeta-proyecto__fechas">
          <Icono nombre="calendario" />
          <span>{periodoCorto(proyecto.fecha_inicio, proyecto.fecha_fin)}</span>
        </span>
      </span>
      <span className="tarjeta-proyecto__pie">
        <span>#{proyecto.id}</span>
      </span>
    </Link>
  )
}

export function ListaProyectos() {
  const proyectos = useProyectos()
  const navegar = useNavigate()
  const [busqueda, setBusqueda] = useState('')

  const total = proyectos.estado === 'listo' ? proyectos.proyectos : []
  const visibles = [...filtrarProyectos(total, busqueda)].sort((a, b) => b.id - a.id)
  const hayFiltro = busqueda.trim() !== ''
  const abiertos = visibles.filter((proyecto) => !proyecto.fecha_cierre)
  const cerrados = visibles.filter((proyecto) => proyecto.fecha_cierre)
  const mostrarAbiertos = abiertos.length > 0 || (!hayFiltro && total.length > 0)
  const irANuevo = () => navegar('/projects/new')

  return (
    <section className="lista-proyectos" aria-labelledby="proyectos-title">
      <header className="lista-proyectos__encabezado">
        <div>
          <h1 id="proyectos-title">Proyectos</h1>
          <p className="project-context">Elegí un proyecto para ver su seguimiento o creá uno nuevo.</p>
        </div>
        <div className="lista-proyectos__acciones">
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
          <Button type="button" onClick={irANuevo}>
            <Icono nombre="mas" />
            Nuevo proyecto
          </Button>
        </div>
      </header>
      {proyectos.estado === 'cargando' && <LoadingIndicator />}
      {proyectos.estado === 'error' && <ErrorMessage message={proyectos.mensaje} />}
      {proyectos.estado === 'listo' && total.length === 0 && (
        <div className="lista-proyectos__vacio">
          <h2>Todavía no tenés proyectos</h2>
          <p>Creá tu primer proyecto para empezar el seguimiento.</p>
          <Button type="button" onClick={irANuevo}>
            <Icono nombre="mas" />
            Crear mi primer proyecto
          </Button>
        </div>
      )}
      {total.length > 0 && visibles.length === 0 && (
        <div className="lista-proyectos__vacio">No hay proyectos que coincidan con la búsqueda.</div>
      )}
      {mostrarAbiertos && (
        <ul className="lista-proyectos__tarjetas">
          {abiertos.map((proyecto) => (
            <li key={proyecto.id}>
              <TarjetaProyecto proyecto={proyecto} />
            </li>
          ))}
          {!hayFiltro && (
            <li>
              <button type="button" className="tarjeta-proyecto tarjeta-proyecto--nueva" onClick={irANuevo}>
                <Icono nombre="mas" />
                <span>Nuevo proyecto</span>
              </button>
            </li>
          )}
        </ul>
      )}
      {cerrados.length > 0 && (
        <section className="lista-proyectos__cerrados" aria-labelledby="proyectos-cerrados-titulo">
          <div className="lista-proyectos__cerrados-cabecera">
            <h2 id="proyectos-cerrados-titulo">Proyectos cerrados</h2>
            <span
              className="lista-proyectos__contador"
              aria-label={`${cerrados.length} ${cerrados.length === 1 ? 'proyecto cerrado' : 'proyectos cerrados'}`}
            >
              {cerrados.length}
            </span>
          </div>
          <ul className="lista-proyectos__tarjetas">
            {cerrados.map((proyecto) => (
              <li key={proyecto.id}>
                <TarjetaProyecto proyecto={proyecto} />
              </li>
            ))}
          </ul>
        </section>
      )}
    </section>
  )
}
