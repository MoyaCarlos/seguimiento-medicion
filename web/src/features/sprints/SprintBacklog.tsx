import { useEffect, useState } from 'react'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { consultarSprintBacklog } from './api'
import { clasificarError } from './errores'
import { EsfuerzoHistoria, type Imputacion } from './EsfuerzoHistoria'
import { EstadoHistoria } from './EstadoHistoria'
import { ETIQUETAS_PRIORIDAD, type HistoriaDeSprint } from './sprint'

type Props = { sprintId: number; permiteCambiarEstado?: boolean; imputacion?: Imputacion }

type Estado =
  | { tipo: 'cargando' }
  | { tipo: 'listo'; historias: HistoriaDeSprint[] }
  | { tipo: 'error'; mensaje: string }

export function SprintBacklog({ sprintId, permiteCambiarEstado = false, imputacion }: Props) {
  return (
    <Consulta
      key={sprintId}
      sprintId={sprintId}
      permiteCambiarEstado={permiteCambiarEstado}
      imputacion={imputacion}
    />
  )
}

function Consulta({ sprintId, permiteCambiarEstado, imputacion }: Props) {
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
          <h4>{historia.titulo}</h4>
          <p>{`${ETIQUETAS_PRIORIDAD[historia.prioridad]} (${historia.prioridad})`}</p>
          <p>
            {historia.estado === 'Completada' && historia.fecha_finalizacion
              ? `Completada el ${historia.fecha_finalizacion}`
              : historia.estado}
          </p>
          <p>{historia.estimacion_sp === null ? 'Sin estimar' : `${historia.estimacion_sp} SP`}</p>
          {permiteCambiarEstado && (
            <EstadoHistoria historia={historia} alCambiar={reemplazar} />
          )}
          <EsfuerzoHistoria historiaId={historia.id} imputacion={imputacion} />
        </li>
      ))}
    </ul>
  )
}
