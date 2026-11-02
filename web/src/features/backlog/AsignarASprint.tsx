import { useState } from 'react'
import type { FormEvent } from 'react'
import { Button } from '../../compartido/componentes/Button'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { asignarASprint } from './api'
import type { SprintParaAsignar } from './api'
import { mensajeDeError } from './edicionHistoria'
import type { Historia } from './historia'

type Props = {
  historia: Historia
  sprints: SprintParaAsignar[] | null
  errorSprints: string
  onAsignada: (historia: Historia, sprintId: number) => void
  onCancelar: () => void
}

function etiquetaDe(sprint: SprintParaAsignar) {
  const base = `Sprint #${sprint.id} (${sprint.estado})`
  return sprint.sprint_goal ? `${base} - ${sprint.sprint_goal}` : base
}

export function AsignarASprint({ historia, sprints, errorSprints, onAsignada, onCancelar }: Props) {
  const [sprintId, setSprintId] = useState('')
  const [error, setError] = useState('')
  const [enviando, setEnviando] = useState(false)

  const disponibles = (sprints ?? []).filter((sprint) => sprint.estado !== 'Finalizado')
  const campo = `asignar-${historia.id}-sprint`

  async function asignar(evento: FormEvent<HTMLFormElement>) {
    evento.preventDefault()
    if (enviando) return

    if (!sprintId) {
      setError('Elegí un Sprint')
      return
    }

    setError('')
    setEnviando(true)
    try {
      const actualizada = await asignarASprint(historia.id, Number(sprintId))
      onAsignada(actualizada, Number(sprintId))
    } catch (causa) {
      setError(mensajeDeError(causa))
      setEnviando(false)
    }
  }

  function contenido() {
    if (errorSprints) return <ErrorMessage message={errorSprints} />
    if (sprints === null) return <LoadingIndicator />
    if (disponibles.length === 0) return <p>No hay Sprints disponibles para asignar</p>

    return (
      <div className="field">
        <label htmlFor={campo}>Sprint</label>
        <select
          id={campo}
          name="sprint"
          value={sprintId}
          onChange={(evento) => {
            setSprintId(evento.target.value)
            setError('')
          }}
          disabled={enviando}
        >
          <option value="">Elegí un Sprint</option>
          {disponibles.map((sprint) => (
            <option key={sprint.id} value={String(sprint.id)}>
              {etiquetaDe(sprint)}
            </option>
          ))}
        </select>
      </div>
    )
  }

  const puedeAsignar = !errorSprints && sprints !== null && disponibles.length > 0

  return (
    <form
      className="backlog-asignacion"
      aria-label={`Asignar ${historia.titulo} a un Sprint`}
      onSubmit={asignar}
      noValidate
    >
      <h3>Asignar a Sprint</h3>
      {contenido()}
      {error && <ErrorMessage message={error} />}
      {enviando && <LoadingIndicator />}
      <div className="backlog-acciones">
        <Button type="button" variant="secondary" onClick={onCancelar} disabled={enviando}>
          Cancelar
        </Button>
        {puedeAsignar && (
          <Button type="submit" disabled={enviando}>
            Confirmar asignación
          </Button>
        )}
      </div>
    </form>
  )
}
