import { BotonCerrar } from './BotonCerrar'
import { FormularioInicio } from './FormularioInicio'
import { SprintBacklog } from './SprintBacklog'
import type { Sprint } from './sprint'

type Props = {
  estado: Sprint['estado']
  sprints: Sprint[]
  alActualizar: (sprint: Sprint) => void
}

export function SeccionSprints({ estado, sprints, alActualizar }: Props) {
  return (
    <section className="sprints-seccion">
      <h2>{estado}</h2>
      {sprints.length === 0 && <p>No hay Sprints</p>}
      {sprints.map((sprint) => (
        <article key={sprint.id} className="sprint-tarjeta">
          <h3>{`Sprint #${sprint.id}`}</h3>
          {sprint.sprint_goal && <p>{sprint.sprint_goal}</p>}
          {sprint.fecha_inicio && sprint.fecha_fin && (
            <p>{`${sprint.fecha_inicio} a ${sprint.fecha_fin}`}</p>
          )}
          <SprintBacklog sprintId={sprint.id} />
          {estado === 'Pendiente' && <FormularioInicio sprint={sprint} alActualizar={alActualizar} />}
          {estado === 'Activo' && <BotonCerrar sprint={sprint} alActualizar={alActualizar} />}
        </article>
      ))}
    </section>
  )
}
