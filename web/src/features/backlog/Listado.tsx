import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router'
import { Button } from '../../compartido/componentes/Button'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { listarBacklog } from './api'
import { Edicion } from './Edicion'
import { mensajeDeError } from './edicionHistoria'
import type { Historia } from './historia'
import { Markdown } from './Markdown'
import { etiquetaPrioridad } from './opciones'
import './backlog.css'

export function Listado() {
  const { id = '' } = useParams()

  return <PantallaBacklog key={id} id={id} />
}

function PantallaBacklog({ id }: { id: string }) {
  const [historias, setHistorias] = useState<Historia[]>([])
  const [cargando, setCargando] = useState(true)
  const [error, setError] = useState('')
  const [editandoId, setEditandoId] = useState<number | null>(null)
  const [aviso, setAviso] = useState('')

  useEffect(() => {
    let vigente = true

    listarBacklog(id)
      .then((datos) => {
        if (vigente) setHistorias(datos)
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

  function guardar(actualizada: Historia) {
    setHistorias((actuales) =>
      actuales.map((historia) => (historia.id === actualizada.id ? actualizada : historia)),
    )
    setEditandoId(null)
    setAviso('Historia actualizada')
  }

  return (
    <section className="backlog" aria-labelledby="backlog-title">
      <div className="page-header backlog-encabezado">
        <div>
          <p className="eyebrow">Product Backlog</p>
          <h1 id="backlog-title">Product Backlog</h1>
          <p className="project-context">
            Proyecto actual: <strong>{id}</strong>
          </p>
        </div>
        <Link className="button button--secondary" to={`/projects/${id}/backlog/nueva`}>
          Nueva historia
        </Link>
      </div>

      {aviso && (
        <p className="message backlog-aviso" role="status">
          {aviso}
        </p>
      )}
      {cargando && <LoadingIndicator />}
      {error && <ErrorMessage message={error} />}

      {!cargando && !error && historias.length === 0 && (
        <p className="backlog-vacio">No hay historias en el backlog</p>
      )}

      <ul className="backlog-lista">
        {historias.map((historia) => (
          <li key={historia.id}>
            <article className="backlog-item" aria-labelledby={`historia-${historia.id}`}>
              <div className="backlog-item__cabecera">
                <div>
                  <p className="backlog-item__meta">Historia #{historia.id}</p>
                  <h2 id={`historia-${historia.id}`}>{historia.titulo}</h2>
                </div>
                <Button
                  type="button"
                  variant="secondary"
                  onClick={() => {
                    setAviso('')
                    setEditandoId(historia.id)
                  }}
                >
                  Editar
                </Button>
              </div>

              <p className="backlog-item__datos">
                <span>{etiquetaPrioridad(historia.prioridad)}</span>
                <span>{historia.estado}</span>
                <span>
                  {historia.estimacion_sp === null
                    ? 'Sin estimar'
                    : `${historia.estimacion_sp} Story Points`}
                </span>
                <span>
                  Valor de negocio:{' '}
                  {historia.valor_negocio === null ? 'Sin valor de negocio' : historia.valor_negocio}
                </span>
              </p>

              <div className="backlog-item__descripcion">
                <Markdown texto={historia.descripcion} />
              </div>

              <div className="backlog-item__criterios">
                <h3>Criterios de aceptación</h3>
                {historia.criterios_aceptacion === null ? (
                  <p>Sin criterios de aceptación</p>
                ) : (
                  <Markdown texto={historia.criterios_aceptacion} />
                )}
              </div>

              {editandoId === historia.id && (
                <Edicion
                  key={historia.id}
                  historia={historia}
                  onGuardada={guardar}
                  onCancelar={() => setEditandoId(null)}
                />
              )}
            </article>
          </li>
        ))}
      </ul>
    </section>
  )
}
