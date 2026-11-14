import { useState } from 'react'
import { Link } from 'react-router'
import { Icono } from '../../compartido/componentes/Icono'
import { periodoCorto } from '../../compartido/formato/fechas'
import { SprintBacklog } from './SprintBacklog'
import type { HistoriaDeSprint, Sprint } from './sprint'

export function SprintFinalizado({ sprint }: { sprint: Sprint }) {
  const [historias, setHistorias] = useState<HistoriaDeSprint[]>()
  const completadas = historias?.filter((h) => h.estado === 'Completada') ?? []
  const puntos = completadas.reduce((total, h) => total + (h.estimacion_sp ?? 0), 0)

  return (
    <details className="sprint-terminado">
      <summary>
        <div>
          <h3>{`Sprint ${sprint.numero}`}</h3>
          {sprint.sprint_goal && <p className="sprint-goal">{sprint.sprint_goal}</p>}
          {sprint.fecha_inicio && sprint.fecha_fin && (
            <span className="sprint-fechas">
              <Icono nombre="calendario" />
              {periodoCorto(sprint.fecha_inicio, sprint.fecha_fin)}
            </span>
          )}
        </div>
        <span className="sprint-chevron">
          <Icono nombre="chevronAbajo" />
        </span>
        {historias && (
          <div className="sprint-chips">
            <span className={`sprint-chip sprint-chip--${completadas.length > 0 ? 'ok' : 'nueva'}`}>
              {`${completadas.length} completada${completadas.length === 1 ? '' : 's'}`}
            </span>
            <span className="sprint-chip sprint-chip--sp">{`${puntos} SP`}</span>
          </div>
        )}
      </summary>
      <SprintBacklog sprintId={sprint.id} finalizado alCargar={setHistorias} />
      <p className="sprint-enlace">
        <Link to={`/projects/${sprint.proyecto_id}/historial/${sprint.id}`}>
          Ver en el historial
          <Icono nombre="flecha" />
        </Link>
      </p>
    </details>
  )
}
