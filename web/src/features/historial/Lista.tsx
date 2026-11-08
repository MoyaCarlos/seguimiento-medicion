import { Link, useParams } from 'react-router'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { listarHistorial } from './api'
import { conUnidad, desviacion, numero, periodo } from './formato'
import type { SprintHistorial } from './tipos'
import { useConsulta } from './useConsulta'

export function Lista() {
  const proyectoId = useParams().id ?? ''
  const { datos: sprints, error } = useConsulta(() => listarHistorial(proyectoId), proyectoId)

  return (
    <div className="historial-pantalla">
      <header className="historial-encabezado">
        <div>
          <p className="historial-eyebrow">Proyecto #{proyectoId}</p>
          <h1>Historial de Sprints</h1>
        </div>
        <p className="historial-bajada">Compará el desempeño del equipo en cada Sprint cerrado.</p>
      </header>

      {error && <ErrorMessage message={error} />}
      {!sprints && !error && <LoadingIndicator />}
      {sprints?.length === 0 && <EstadoVacio proyectoId={proyectoId} />}
      {sprints && sprints.length > 0 && (
        <ol className="historial-lista">
          {sprints.map((sprint) => (
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
        <div className="historial-sprint__titulos">
          <h2>
            <Link className="historial-sprint__enlace" to={`/projects/${proyectoId}/historial/${sprint.id}`}>
              {nombre}
            </Link>
          </h2>
          <p className="historial-goal">{sprint.sprint_goal}</p>
          <p className="historial-fechas">{periodo(sprint.fecha_inicio, sprint.fecha_fin)}</p>
        </div>
        <span className="historial-chevron" aria-hidden="true" />
      </div>
      <BarraAvance completados={metricas.story_points_completados} planificados={metricas.story_points_planificados} />
      <dl className="historial-metricas">
        <Metrica nombre="Story Points planificados" valor={numero(metricas.story_points_planificados)} />
        <Metrica nombre="Story Points completados" valor={numero(metricas.story_points_completados)} />
        <Metrica nombre="Historias completadas" valor={conUnidad(metricas.porcentaje_historias_completadas, '%')} />
        <Metrica nombre="Horas estimadas" valor={conUnidad(metricas.horas_estimadas, ' h')} />
        <Metrica nombre="Horas reales" valor={conUnidad(metricas.horas_reales, ' h')} />
        <Metrica nombre="Desviación" valor={desviacion(metricas.desviacion_horas, metricas.desviacion_porcentaje)} />
        <Metrica nombre="Defectos detectados" valor={numero(metricas.defectos_detectados)} />
        <Metrica nombre="Defectos resueltos" valor={numero(metricas.defectos_resueltos)} />
      </dl>
    </article>
  )
}

export function Metrica({ nombre, valor }: { nombre: string; valor: string }) {
  return (
    <div className="historial-metrica">
      <dt>{nombre}</dt>
      <dd className={valor.includes('sin datos') ? 'historial-sin-datos' : undefined}>{valor}</dd>
    </div>
  )
}

export function BarraAvance({ completados, planificados }: { completados: number; planificados: number }) {
  const porcentaje = planificados > 0 ? Math.min((completados / planificados) * 100, 100) : 0

  return (
    <div
      className="historial-barra"
      role="meter"
      aria-label="Story Points completados sobre planificados"
      aria-valuemin={0}
      aria-valuemax={planificados > 0 ? planificados : 1}
      aria-valuenow={completados}
    >
      <span style={{ width: `${porcentaje}%` }} />
    </div>
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
