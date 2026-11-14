import { useEffect, useState } from 'react'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { fechaCorta } from '../../compartido/formato/fechas'
import { consultarSprintBacklog } from './api'
import { clasificarError } from './errores'
import { EsfuerzoHistoria, type Imputacion } from './EsfuerzoHistoria'
import { EstadoHistoria } from './EstadoHistoria'
import { ETIQUETAS_PRIORIDAD, type HistoriaDeSprint } from './sprint'

const CLASE_PRIORIDAD = { M: 'must', S: 'should', C: 'could', W: 'wont' } as const
const CLASE_ESTADO = { Nueva: 'nueva', 'En progreso': 'progreso', Completada: 'ok' } as const

type Props = {
  sprintId: number
  permiteCambiarEstado?: boolean
  imputacion?: Imputacion
  finalizado?: boolean
  alCargar?: (historias: HistoriaDeSprint[]) => void
}

type Estado =
  | { tipo: 'cargando' }
  | { tipo: 'listo'; historias: HistoriaDeSprint[] }
  | { tipo: 'error'; mensaje: string }

export function SprintBacklog({ sprintId, permiteCambiarEstado = false, imputacion, finalizado, alCargar }: Props) {
  return (
    <Consulta
      key={sprintId}
      sprintId={sprintId}
      permiteCambiarEstado={permiteCambiarEstado}
      imputacion={imputacion}
      finalizado={finalizado}
      alCargar={alCargar}
    />
  )
}

function Consulta({ sprintId, permiteCambiarEstado, imputacion, finalizado, alCargar }: Props) {
  const [estado, setEstado] = useState<Estado>({ tipo: 'cargando' })

  useEffect(() => {
    let vigente = true

    consultarSprintBacklog(sprintId)
      .then((historias) => {
        if (!vigente) return
        setEstado({ tipo: 'listo', historias })
        alCargar?.(historias)
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

  function reemplazar(nueva: HistoriaDeSprint) {
    setEstado((actual) =>
      actual.tipo === 'listo'
        ? { tipo: 'listo', historias: actual.historias.map((h) => (h.id === nueva.id ? nueva : h)) }
        : actual,
    )
  }

  if (estado.tipo === 'cargando') return <LoadingIndicator />
  if (estado.tipo === 'error') return <ErrorMessage message={estado.mensaje} />
  if (estado.historias.length === 0) return <p>Este Sprint no tiene historias</p>

  return (
    <ul className="sprint-backlog">
      {estado.historias.map((historia) => (
        <li key={historia.id} className="sprint-historia">
          <h4 className="sprint-historia__titulo">{historia.titulo}</h4>
          <div className="sprint-chips">
            <span className={`sprint-chip sprint-chip--${CLASE_PRIORIDAD[historia.prioridad]}`}>
              {`${ETIQUETAS_PRIORIDAD[historia.prioridad]} (${historia.prioridad})`}
            </span>
            {!finalizado && (
              <span className={`sprint-chip sprint-chip--${CLASE_ESTADO[historia.estado]}`}>{historia.estado}</span>
            )}
            {historia.estimacion_sp === null ? (
              <span className="sprint-chip sprint-chip--ninguno">Sin estimar</span>
            ) : (
              <span className="sprint-chip sprint-chip--sp">{`${historia.estimacion_sp} SP`}</span>
            )}
            {finalizado && historia.estado !== 'Completada' && (
              <span className="sprint-chip sprint-chip--nueva">No completada</span>
            )}
          </div>
          {historia.estado === 'Completada' && historia.fecha_finalizacion && (
            <p className="sprint-historia__nota">{`Completada el ${fechaCorta(historia.fecha_finalizacion)}`}</p>
          )}
          <EsfuerzoHistoria
            historiaId={historia.id}
            imputacion={imputacion}
            acciones={
              permiteCambiarEstado ? <EstadoHistoria historia={historia} alCambiar={reemplazar} /> : undefined
            }
          />
        </li>
      ))}
    </ul>
  )
}
