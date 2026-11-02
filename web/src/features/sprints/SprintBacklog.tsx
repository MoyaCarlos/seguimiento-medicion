import { useEffect, useState } from 'react'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { consultarSprintBacklog } from './api'
import { clasificarError } from './errores'
import { ETIQUETAS_PRIORIDAD, type HistoriaDeSprint } from './sprint'

type Props = { sprintId: number }

type Estado =
  | { tipo: 'cargando' }
  | { tipo: 'listo'; historias: HistoriaDeSprint[] }
  | { tipo: 'error'; mensaje: string }

export function SprintBacklog({ sprintId }: Props) {
  return <Consulta key={sprintId} sprintId={sprintId} />
}

function Consulta({ sprintId }: Props) {
  const [estado, setEstado] = useState<Estado>({ tipo: 'cargando' })

  useEffect(() => {
    let vigente = true

    consultarSprintBacklog(sprintId)
      .then((historias) => {
        if (vigente) setEstado({ tipo: 'listo', historias })
      })
      .catch((error: unknown) => {
        if (!vigente) return
        const clasificado = clasificarError(error, [])
        setEstado({
          tipo: 'error',
          mensaje: 'aviso' in clasificado ? clasificado.aviso : clasificado.mensaje,
        })
      })

    return () => {
      vigente = false
    }
  }, [sprintId])

  if (estado.tipo === 'cargando') return <LoadingIndicator />
  if (estado.tipo === 'error') return <ErrorMessage message={estado.mensaje} />
  if (estado.historias.length === 0) return <p>Este Sprint no tiene historias</p>

  return (
    <ul className="sprint-backlog">
      {estado.historias.map((historia) => (
        <li key={historia.id}>
          <h4>{historia.titulo}</h4>
          <p>{`${ETIQUETAS_PRIORIDAD[historia.prioridad]} (${historia.prioridad})`}</p>
          <p>{historia.estado}</p>
          <p>{historia.estimacion_sp === null ? 'Sin estimar' : `${historia.estimacion_sp} SP`}</p>
        </li>
      ))}
    </ul>
  )
}
