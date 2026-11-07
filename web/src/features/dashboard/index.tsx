import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router'
import { ApiError, ConnectionError } from '../../compartido/api/client'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { consultarMetricas } from './api'
import { GraficoVelocidad } from './GraficoVelocidad'
import { Indicadores } from './Indicadores'
import type { MetricasProyecto } from './metricas'
import './dashboard.css'

export function Component() {
  const { id } = useParams()

  return <PantallaDashboard key={id} id={id ?? ''} />
}

function mensajeDeError(causa: unknown): string {
  if (causa instanceof ApiError || causa instanceof ConnectionError) return causa.message
  return 'Ocurrió un error inesperado'
}

function PantallaDashboard({ id }: { id: string }) {
  const [metricas, setMetricas] = useState<MetricasProyecto>()
  const [error, setError] = useState<string>()

  useEffect(() => {
    let vigente = true

    consultarMetricas(id)
      .then((obtenidas) => {
        if (vigente) setMetricas(obtenidas)
      })
      .catch((causa: unknown) => {
        if (vigente) setError(mensajeDeError(causa))
      })

    return () => {
      vigente = false
    }
  }, [id])

  return (
    <div className="dashboard-pantalla">
      <h1>Dashboard</h1>
      {error && <ErrorMessage message={error} />}
      {!metricas && !error && <LoadingIndicator />}
      {metricas && metricas.sprints.length === 0 && <EstadoVacio proyectoId={id} />}
      {metricas && metricas.sprints.length > 0 && (
        <>
          <GraficoVelocidad sprints={metricas.sprints} />
          <Indicadores
            metricas={metricas.proyecto}
            horasPorStoryPoint={metricas.horas_por_story_point}
            proyectoId={id}
          />
        </>
      )}
    </div>
  )
}

function EstadoVacio({ proyectoId }: { proyectoId: string }) {
  return (
    <section className="dashboard-tarjeta dashboard-vacio">
      <h2>Todavía no hay métricas</h2>
      <p>Las métricas aparecen cuando se cierra el primer Sprint del proyecto.</p>
      <Link className="button button--secondary" to={`/projects/${proyectoId}/sprints`}>
        Ir a Sprints
      </Link>
    </section>
  )
}
