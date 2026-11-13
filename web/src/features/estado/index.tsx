import { useEffect, useState } from 'react'
import { useParams } from 'react-router'
import { ApiError, ConnectionError } from '../../compartido/api/client'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { consultarAvance, consultarEstadoDelProyecto, type EstadoDelProyecto } from './api'
import { esIdValido, type Avance } from './avance'
import { GrupoHistorias } from './GrupoHistorias'
import { Resumen } from './Resumen'
import './estado.css'

export function Component() {
  const { id = '' } = useParams()

  return (
    <div className="estado-pantalla">
      {esIdValido(id) && <p className="eyebrow estado-pretitulo">Proyecto #{id}</p>}
      <div className="estado-titulo">
        <h1>Estado del proyecto</h1>
        {esIdValido(id) && <InsigniaDeEstado key={id} id={id} />}
      </div>
      {esIdValido(id) ? (
        <>
          <p className="project-context">Proyecto actual: <strong>{id}</strong></p>
          <AvanceDelProyecto key={id} id={id} />
        </>
      ) : (
        <ErrorMessage message={`El proyecto "${id}" no es válido`} />
      )}
    </div>
  )
}

function InsigniaDeEstado({ id }: { id: string }) {
  const [estado, setEstado] = useState<EstadoDelProyecto>()

  useEffect(() => {
    let vigente = true
    consultarEstadoDelProyecto(id)
      .then((obtenido) => {
        if (vigente) setEstado(obtenido.estado)
      })
      .catch(() => {})
    return () => {
      vigente = false
    }
  }, [id])

  if (!estado) return null
  return <span className={`estado-insignia estado-insignia--${estado.toLowerCase()}`}>{estado}</span>
}

function mensajeDeError(causa: unknown): string {
  if (causa instanceof ApiError && causa.status === 404) return 'El proyecto no existe'
  if (causa instanceof ApiError || causa instanceof ConnectionError) return causa.message
  return 'Ocurrió un error inesperado'
}

function AvanceDelProyecto({ id }: { id: string }) {
  const [avance, setAvance] = useState<Avance>()
  const [error, setError] = useState<string>()

  useEffect(() => {
    let vigente = true

    consultarAvance(id)
      .then((obtenido) => {
        if (vigente) setAvance(obtenido)
      })
      .catch((causa: unknown) => {
        if (vigente) setError(mensajeDeError(causa))
      })

    return () => {
      vigente = false
    }
  }, [id])

  if (error) return <ErrorMessage message={error} />
  if (!avance) return <LoadingIndicator />

  return (
    <>
      {avance.total.historias === 0 && <p className="estado-vacio">El proyecto todavía no tiene historias.</p>}
      <Resumen avance={avance} />
      <div className="estado-tablero">
        {avance.grupos.map((grupo) => <GrupoHistorias key={grupo.estado} grupo={grupo} />)}
      </div>
    </>
  )
}
