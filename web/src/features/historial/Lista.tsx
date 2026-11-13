import { useCallback } from 'react'
import { Link, useParams } from 'react-router'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { listarHistorial } from './api'
import { formatoNumero as numero, textoDeDesviacion } from '../../compartido/formato/numeros'
import { conUnidad, periodo } from './formato'
import { BarraAvance, Metrica } from './Piezas'
import type { SprintHistorial } from './tipos'
import { useConsulta } from './useConsulta'

export function Lista() {
  const proyectoId = useParams().id ?? ''
  const consultar = useCallback(() => listarHistorial(proyectoId), [proyectoId])
  const { datos: sprints, error } = useConsulta(consultar)

  return (
    <div className="historial-pantalla">
      <header className="historial-encabezado">
        <p className="historial-eyebrow">Proyecto #{proyectoId}</p>
        <h1>Historial de Sprints</h1>
        <p className="historial-bajada">Compará el desempeño del equipo en cada Sprint cerrado.</p>
      </header>

      {error && <ErrorMessage message={error} />}
      {!sprints && !error && <LoadingIndicator />}
      {sprints?.length === 0 && <EstadoVacio proyectoId={proyectoId} />}
      {sprints && sprints.length > 0 && (
        <ol className="historial-lista">
          {[...sprints].sort((a, b) => b.numero - a.numero).map((sprint) => (
            <li key={sprint.id}>
              <TarjetaSprint sprint={sprint} proyectoId={proyectoId} />
            </li>
          ))}
        </ol>
      )}
    </div>
  )
}

function TarjetaSprint({ sprint, proyectoId }: { sprint: SprintHistorial; proyectoId: string }) {
  const { metricas } = sprint
  const nombre = `Sprint ${sprint.numero}`

  return (
    <article className="historial-tarjeta historial-sprint" aria-label={nombre}>
      <div className="historial-sprint__cabecera">
        <h2>
          <Link className="historial-sprint__enlace" to={`/projects/${proyectoId}/historial/${sprint.id}`}>
            {nombre}
          </Link>
        </h2>
        <p className="historial-fechas">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
            <rect x="3" y="4" width="18" height="18" rx="2" />
            <path d="M16 2v4M8 2v4M3 10h18" />
          </svg>
          {periodo(sprint.fecha_inicio, sprint.fecha_fin)}
        </p>
      </div>
      {sprint.sprint_goal && <p className="historial-goal">{sprint.sprint_goal}</p>}
      <div className="historial-avance">
        <BarraAvance completados={metricas.story_points_completados} planificados={metricas.story_points_planificados} />
        <span>
          {numero(metricas.story_points_completados)} / {numero(metricas.story_points_planificados)} SP
        </span>
      </div>
      <dl className="historial-metricas">
        <Metrica nombre="Story Points planificados" valor={numero(metricas.story_points_planificados)} />
        <Metrica nombre="Story Points completados" valor={numero(metricas.story_points_completados)} />
        <Metrica nombre="Historias completadas" valor={conUnidad(metricas.porcentaje_historias_completadas, ' %')} />
        <Metrica nombre="Horas estimadas" valor={conUnidad(metricas.horas_estimadas, ' h')} />
        <Metrica nombre="Horas reales" valor={conUnidad(metricas.horas_reales, ' h')} />
        <Metrica nombre="Desviación" valor={textoDeDesviacion(metricas.desviacion_horas, metricas.desviacion_porcentaje)} />
        <Metrica nombre="Defectos detectados" valor={numero(metricas.defectos_detectados)} />
        <Metrica nombre="Defectos resueltos" valor={numero(metricas.defectos_resueltos)} />
      </dl>
    </article>
  )
}

function EstadoVacio({ proyectoId }: { proyectoId: string }) {
  return (
    <section className="historial-vacio">
      <h2>Todavía no hay Sprints cerrados</h2>
      <p>Cada Sprint aparece en el historial cuando se cierra, con su Sprint Goal y sus métricas.</p>
      <Link className="button button--secondary" to={`/projects/${proyectoId}/sprints`}>
        Ir a Sprints
      </Link>
    </section>
  )
}
