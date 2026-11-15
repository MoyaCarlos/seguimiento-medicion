import { useEffect, useState } from 'react'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { consultarSprintBacklog } from './api'
import { clasificarError } from './errores'
import { EsfuerzoHistoria, type Imputacion } from './EsfuerzoHistoria'
import { EstadoHistoria } from './EstadoHistoria'
import { ETIQUETAS_PRIORIDAD, type HistoriaDeSprint } from './sprint'

type Props = {
  sprintId: number
  permiteCambiarEstado?: boolean
  imputacion?: Imputacion
  puedeEscribir?: boolean
}

type Estado =
  | { tipo: 'cargando' }
  | { tipo: 'listo'; historias: HistoriaDeSprint[] }
  | { tipo: 'error'; mensaje: string }

export function SprintBacklog({ sprintId, permiteCambiarEstado = false, imputacion, puedeEscribir = true }: Props) {
  return (
    <Consulta
      key={sprintId}
      sprintId={sprintId}
      permiteCambiarEstado={permiteCambiarEstado}
      imputacion={imputacion}
      puedeEscribir={puedeEscribir}
    />
  )
}

function Consulta({ sprintId, permiteCambiarEstado, imputacion, puedeEscribir }: Props) {
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
        <li key={historia.id}>
          <details className="historia-tarjeta">
            <summary className="historia-tarjeta__resumen">
              <h4>{historia.titulo}</h4>
              <span className="historia-marca" data-estado={historia.estado}>
                {historia.estado}
              </span>
              <span className="historia-tarjeta__sp">
                {historia.estimacion_sp === null ? 'Sin estimar' : `${historia.estimacion_sp} SP`}
              </span>
            </summary>
            <div className="historia-tarjeta__detalle">
              <p>{`${ETIQUETAS_PRIORIDAD[historia.prioridad]} (${historia.prioridad})`}</p>
              {historia.estado === 'Completada' && historia.fecha_finalizacion && (
                <p>{`Completada el ${historia.fecha_finalizacion}`}</p>
              )}
              {permiteCambiarEstado && (
                <EstadoHistoria historia={historia} alCambiar={reemplazar} puedeEscribir={puedeEscribir} />
              )}
              <EsfuerzoHistoria historiaId={historia.id} imputacion={imputacion} puedeEscribir={puedeEscribir} />
            </div>
          </details>
        </li>
      ))}
    </ul>
  )
}
