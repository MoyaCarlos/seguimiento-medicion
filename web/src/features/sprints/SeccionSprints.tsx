import { Icono } from '../../compartido/componentes/Icono'
import { periodoCorto } from '../../compartido/formato/fechas'
import { BotonCerrar } from './BotonCerrar'
import { FormularioInicio } from './FormularioInicio'
import { SprintBacklog } from './SprintBacklog'
import type { Sprint } from './sprint'

const VARIANTE = { Pendiente: '', Activo: ' sprints-carril--activo', Finalizado: ' sprints-carril--fin' } as const

type Props = {
  estado: Sprint['estado']
  sprints: Sprint[]
  alActualizar: (sprint: Sprint) => void
  puedeGestionar?: boolean
}

export function SeccionSprints({ estado, sprints, alActualizar, puedeGestionar = true }: Props) {
  return (
    <section className={`sprints-carril${VARIANTE[estado]}`}>
      <h2>{estado}</h2>
      {sprints.length === 0 && <p className="sprints-vacio">No hay Sprints</p>}
      {sprints.map((sprint) => (
        <article key={sprint.id} className="sprint-tarjeta">
          <h3>{`Sprint ${sprint.numero}`}</h3>
          {sprint.sprint_goal && <p className="sprint-goal">{sprint.sprint_goal}</p>}
          {sprint.fecha_inicio && sprint.fecha_fin && (
            <span className="sprint-fechas">
              <Icono nombre="calendario" />
              {periodoCorto(sprint.fecha_inicio, sprint.fecha_fin)}
            </span>
          )}
          <SprintBacklog
            sprintId={sprint.id}
            permiteCambiarEstado={estado === 'Activo'}
            imputacion={
              estado === 'Activo' ? { proyectoId: sprint.proyecto_id, fechaFin: sprint.fecha_fin } : undefined
            }
          />
          {estado === 'Pendiente' && (
            <FormularioInicio sprint={sprint} alActualizar={alActualizar} puedeGestionar={puedeGestionar} />
          )}
          {estado === 'Activo' && (
            <BotonCerrar sprint={sprint} alActualizar={alActualizar} puedeGestionar={puedeGestionar} />
          )}
        </article>
      ))}
    </section>
  )
}
