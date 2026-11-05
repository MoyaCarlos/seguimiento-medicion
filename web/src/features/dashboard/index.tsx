import { useEffect, useState } from 'react'
import { useParams } from 'react-router'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { consultarMetricas } from './api'
import { GraficoVelocidad } from './GraficoVelocidad'
import type { MetricasProyecto } from './metricas'
import './dashboard.css'

export function Component() {
  const { id } = useParams()

  return <PantallaDashboard key={id} id={id ?? ''} />
}

function PantallaDashboard({ id }: { id: string }) {
  const [metricas, setMetricas] = useState<MetricasProyecto>()

  useEffect(() => {
    let vigente = true

    consultarMetricas(id).then((obtenidas) => {
      if (vigente) setMetricas(obtenidas)
    })

    return () => {
      vigente = false
    }
  }, [id])

  return (
    <div className="dashboard-pantalla">
      <h1>Dashboard</h1>
      {!metricas && <LoadingIndicator />}
      {metricas && <GraficoVelocidad sprints={metricas.sprints} />}
    </div>
  )
}
