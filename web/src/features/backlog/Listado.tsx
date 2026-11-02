import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router'
import { Button } from '../../compartido/componentes/Button'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { listarBacklog, listarSprintsDelProyecto } from './api'
import type { SprintParaAsignar } from './api'
import { AsignarASprint } from './AsignarASprint'
import { Edicion } from './Edicion'
import { mensajeDeError } from './edicionHistoria'
import type { Historia } from './historia'
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
  const [asignandoId, setAsignandoId] = useState<number | null>(null)
  const [sprints, setSprints] = useState<SprintParaAsignar[] | null>(null)
  const [errorSprints, setErrorSprints] = useState('')
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

  useEffect(() => {
    let vigente = true

    listarSprintsDelProyecto(id)
      .then((datos) => {
        if (vigente) setSprints(datos)
      })
      .catch((causa: unknown) => {
        if (vigente) setErrorSprints(mensajeDeError(causa))
      })

    return () => {
      vigente = false
    }
  }, [id])

  function asignada(historia: Historia, sprintId: number) {
    setHistorias((actuales) => actuales.filter((actual) => actual.id !== historia.id))
    setAsignandoId(null)
    setAviso(`Historia asignada al Sprint #${sprintId}`)
  }

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
                <div className="backlog-acciones">
                  <Button
                    type="button"
                    variant="secondary"
                    onClick={() => {
                      setAviso('')
                      setAsignandoId(null)
                      setEditandoId(historia.id)
                    }}
                  >
                    Editar
                  </Button>
                  <Button
                    type="button"
                    variant="secondary"
                    onClick={() => {
                      setAviso('')
                      setEditandoId(null)
                      setAsignandoId(historia.id)
                    }}
                  >
                    Asignar a Sprint
                  </Button>
                </div>
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

              <p className="backlog-item__descripcion">{historia.descripcion}</p>

              <div className="backlog-item__criterios">
                <h3>Criterios de aceptación</h3>
                <p>{historia.criterios_aceptacion ?? 'Sin criterios de aceptación'}</p>
              </div>

              {editandoId === historia.id && (
                <Edicion
                  key={historia.id}
                  historia={historia}
                  onGuardada={guardar}
                  onCancelar={() => setEditandoId(null)}
                />
              )}

              {asignandoId === historia.id && (
                <AsignarASprint
                  key={historia.id}
                  historia={historia}
                  sprints={sprints}
                  errorSprints={errorSprints}
                  onAsignada={asignada}
                  onCancelar={() => setAsignandoId(null)}
                />
              )}
            </article>
          </li>
        ))}
      </ul>
    </section>
  )
}
