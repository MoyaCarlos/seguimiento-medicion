import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router'
import { Button } from '../../compartido/componentes/Button'
import { Icono } from '../../compartido/componentes/Icono'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { listarBacklog, listarSprintsDelProyecto } from './api'
import type { SprintParaAsignar } from './api'
import { AsignarASprint } from './AsignarASprint'
import { Edicion } from './Edicion'
import { mensajeDeError } from './edicionHistoria'
import type { Historia } from './historia'
import { Markdown } from './Markdown'
import { etiquetaPrioridad } from './opciones'
import type { Prioridad } from './opciones'
import './backlog.css'

const MARCADOR: Record<Prioridad, string> = { M: 'must', S: 'should', C: 'could', W: 'wont' }

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
    const numero = sprints?.find((sprint) => sprint.id === sprintId)?.numero ?? sprintId
    setAviso(`Historia asignada al Sprint ${numero}`)
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
          <p className="eyebrow">Proyecto #{id}</p>
          <h1 id="backlog-title">Product Backlog</h1>
          {!cargando && !error && (
            <p className="project-context">
              {historias.length} {historias.length === 1 ? 'historia' : 'historias'} sin Sprint
            </p>
          )}
        </div>
        <Link className="button button--primary" to={`/projects/${id}/backlog/nueva`}>
          <Icono nombre="mas" />
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
            <article
              className={`backlog-item backlog-item--${MARCADOR[historia.prioridad]}`}
              aria-labelledby={`historia-${historia.id}`}>
              <div className="backlog-item__cabecera">
                <div className="backlog-item__titulo">
                  <span className="backlog-item__id">#{historia.id}</span>
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
                    <Icono nombre="editar" />
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
                  <Link className="button button--secondary backlog-votar" to={`/projects/${id}/votacion/${historia.id}`}>
                    Votar
                  </Link>
                </div>
              </div>

              <ul className="backlog-chips">
                <li className={`backlog-chip backlog-chip--${MARCADOR[historia.prioridad]}`}>
                  {etiquetaPrioridad(historia.prioridad)}
                </li>
                <li className="backlog-chip backlog-chip--nueva">{historia.estado}</li>
                {historia.estimacion_sp === null ? (
                  <li className="backlog-chip backlog-chip--ninguno">Sin estimar</li>
                ) : (
                  <li className="backlog-chip backlog-chip--sp">{historia.estimacion_sp} Story Points</li>
                )}
                {historia.valor_negocio === null ? (
                  <li className="backlog-chip backlog-chip--ninguno">Sin valor de negocio</li>
                ) : (
                  <li className="backlog-chip backlog-chip--sp">Valor de negocio: {historia.valor_negocio}</li>
                )}
              </ul>

              <div className="backlog-item__descripcion">
                <Markdown texto={historia.descripcion} />
              </div>

              <div className="backlog-item__criterios">
                <h3>Criterios de aceptación</h3>
                {historia.criterios_aceptacion === null ? (
                  <p className="backlog-item__vacio">Sin criterios de aceptación</p>
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
