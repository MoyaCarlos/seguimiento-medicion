import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { obtenerProyecto } from './api'
import type { Proyecto } from './api'
import { mensajeDeError } from './errores'

function mostrarFecha(fecha: string | null) {
  return fecha ?? 'Sin definir'
}

export function Panel() {
  const { id = '' } = useParams()
  const [proyecto, setProyecto] = useState<Proyecto>()
  const [error, setError] = useState('')
  const [cargando, setCargando] = useState(true)

  useEffect(() => {
    let vigente = true

    obtenerProyecto(id)
      .then((datos) => {
        if (vigente) setProyecto(datos)
      })
      .catch((causa: unknown) => {
        if (vigente) setError(mensajeDeError(causa))
      })
      .finally(() => {
        if (vigente) setCargando(false)
      })

    return () => {
      vigente = false
    }
  }, [id])

  return (
    <section className="hu04-panel" aria-labelledby="panel-title">
      <div className="hu04-encabezado">
        <div>
          <p className="eyebrow">Panel principal</p>
          <h1 id="panel-title">Panel del proyecto</h1>
          <p className="hu04-referencia">
            Proyecto <strong>{id}</strong>
          </p>
        </div>
        <div className="hu04-acciones">
          <Link className="button button--secondary" to="settings">
            Configurar equipo
          </Link>
          <Link className="button button--primary" to="edit">
            Editar proyecto
          </Link>
        </div>
      </div>

      {cargando && <LoadingIndicator />}
      {error && <ErrorMessage message={error} />}

      {proyecto && (
        <article className="hu04-tarjeta">
          <h2>{proyecto.nombre}</h2>
          <p>{proyecto.descripcion || 'Sin descripción'}</p>
          <dl className="hu04-datos">
            <div>
              <dt>Inicio</dt>
              <dd>{mostrarFecha(proyecto.fecha_inicio)}</dd>
            </div>
            <div>
              <dt>Fin</dt>
              <dd>{mostrarFecha(proyecto.fecha_fin)}</dd>
            </div>
          </dl>
        </article>
      )}

      <Link className="hu04-enlace" to={`/projects/${id}/status`}>
        Ver estado del proyecto
      </Link>
    </section>
  )
}
