import { useState } from 'react'
import { Link, useNavigate } from 'react-router'
import { AvatarProyecto } from '../componentes/AvatarProyecto'
import { Button } from '../componentes/Button'
import { ErrorMessage } from '../componentes/ErrorMessage'
import { Icono } from '../componentes/Icono'
import { LoadingIndicator } from '../componentes/LoadingIndicator'
import { periodoCorto } from '../formato/fechas'
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
        <AvatarProyecto id={proyecto.id} nombre={proyecto.nombre} />
        {proyecto.rol && <span className="tarjeta-proyecto__rol">{ETIQUETA_ROL[proyecto.rol]}</span>}
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
      {visibles.length > 0 && (
        <ul className="lista-proyectos__tarjetas">
          {visibles.map((proyecto) => (
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
    </section>
  )
}
