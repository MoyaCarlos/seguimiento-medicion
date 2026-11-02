import { useCallback, useEffect, useState } from 'react'
import { useParams } from 'react-router'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { consultarVotacion, listarIntegrantes, votar } from './api'
import { Cartas } from './Cartas'
import { mensajeDeError } from './errores'
import { Integrantes } from './Integrantes'
import { estaRevelada, type Integrante, type Votacion } from './votacion'
import './votacion.css'

export function Component() {
  const { id = '', historiaId = '' } = useParams()

  return <PantallaVotacion key={`${id}/${historiaId}`} proyectoId={id} historiaId={Number(historiaId)} />
}

type Props = { proyectoId: string; historiaId: number }

function PantallaVotacion({ proyectoId, historiaId }: Props) {
  const [integrantes, setIntegrantes] = useState<Integrante[] | null>(null)
  const [actorId, setActorId] = useState<number>()
  const [votacion, setVotacion] = useState<Votacion | null>(null)
  const [cargadaPara, setCargadaPara] = useState('')
  const [enviando, setEnviando] = useState(false)
  const [errorDeCarga, setErrorDeCarga] = useState('')
  const [errorDeAccion, setErrorDeAccion] = useState('')
  const [recarga, setRecarga] = useState(0)

  const consulta = `${historiaId}/${actorId}/${recarga}`
  const cargando = cargadaPara !== consulta

  const actor = integrantes?.find((integrante) => integrante.id === actorId)
  const recargar = useCallback(() => setRecarga((actual) => actual + 1), [])

  useEffect(() => {
    let vigente = true

    listarIntegrantes(proyectoId)
      .then((datos) => {
        if (vigente) setIntegrantes(datos)
      })
      .catch((causa: unknown) => {
        if (vigente) setErrorDeCarga(mensajeDeError(causa))
      })

    return () => {
      vigente = false
    }
  }, [proyectoId])

  useEffect(() => {
    let vigente = true

    consultarVotacion(historiaId, actorId)
      .then((datos) => {
        if (!vigente) return
        setVotacion(datos)
        setErrorDeCarga('')
      })
      .catch((causa: unknown) => {
        if (vigente) setErrorDeCarga(mensajeDeError(causa))
      })
      .finally(() => {
        if (vigente) setCargadaPara(consulta)
      })

    return () => {
      vigente = false
    }
  }, [historiaId, actorId, consulta])

  async function ejecutar(accion: () => Promise<Votacion>) {
    if (enviando) return
    setErrorDeAccion('')
    setEnviando(true)
    try {
      setVotacion(await accion())
    } catch (causa) {
      setErrorDeAccion(mensajeDeError(causa))
      recargar()
    } finally {
      setEnviando(false)
    }
  }

  const miVoto = votacion?.integrantes.find((voto) => voto.usuario_id === actorId)?.valor ?? null
  const puedeVotar = votacion !== null && !estaRevelada(votacion) && actor !== undefined

  return (
    <section className="votacion" aria-labelledby="votacion-title">
      <div className="page-header">
        <div>
          <p className="eyebrow">Planning Poker</p>
          <h1 id="votacion-title">Votación de la historia #{historiaId}</h1>
          <p className="project-context">
            Proyecto actual: <strong>{proyectoId}</strong>
          </p>
        </div>
      </div>

      {errorDeCarga && <ErrorMessage message={errorDeCarga} />}
      {errorDeAccion && <ErrorMessage message={errorDeAccion} />}
      {cargando && <LoadingIndicator />}

      {integrantes !== null && integrantes.length === 0 && <p>El proyecto no tiene integrantes</p>}

      {integrantes !== null && integrantes.length > 0 && (
        <div className="field">
          <label htmlFor="votacion-actor">Quién actúa</label>
          <select
            id="votacion-actor"
            value={actorId ?? ''}
            disabled={enviando}
            onChange={(evento) => {
              setErrorDeAccion('')
              setActorId(evento.target.value ? Number(evento.target.value) : undefined)
            }}
          >
            <option value="">Elegí un integrante</option>
            {integrantes.map((integrante) => (
              <option key={integrante.id} value={integrante.id}>
                {integrante.nombre}
              </option>
            ))}
          </select>
        </div>
      )}

      {votacion && (
        <>
          <p className="votacion-ronda">Ronda {votacion.ronda}</p>
          <Integrantes integrantes={votacion.integrantes} />
          {!estaRevelada(votacion) && actor === undefined && integrantes?.length ? (
            <p>Elegí quién actúa para votar</p>
          ) : null}
          {puedeVotar && (
            <Cartas
              seleccionada={miVoto}
              deshabilitada={enviando}
              alElegir={(valor) => ejecutar(() => votar(historiaId, actor.id, valor))}
            />
          )}
        </>
      )}
    </section>
  )
}
