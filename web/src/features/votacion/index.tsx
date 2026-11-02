import { useCallback, useEffect, useRef, useState } from 'react'
import { Link, useParams } from 'react-router'
import { Button } from '../../compartido/componentes/Button'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { consultarVotacion, iniciarVotacion, listarIntegrantes, obtenerHistoria, revelarVotacion, votar } from './api'
import { Cartas } from './Cartas'
import { esVotacionNoIniciada, mensajeDeError } from './errores'
import { Integrantes } from './Integrantes'
import { Resultado } from './Resultado'
import { ResumenHistoria } from './ResumenHistoria'
import { estaRevelada, faltanPorVotar, puedeRevelar, esScrumMaster, type HistoriaResumen, type Integrante, type Votacion } from './votacion'
import './votacion.css'

export function Component() {
  const { id = '', historiaId = '' } = useParams()

  return <PantallaVotacion key={`${id}/${historiaId}`} proyectoId={id} historiaId={Number(historiaId)} />
}

type Props = { proyectoId: string; historiaId: number }

function PantallaVotacion({ proyectoId, historiaId }: Props) {
  const [integrantes, setIntegrantes] = useState<Integrante[] | null>(null)
  const [historia, setHistoria] = useState<HistoriaResumen | null>(null)
  const [actorId, setActorId] = useState<number>()
  const [votacion, setVotacion] = useState<Votacion | null>(null)
  const [sinVotacion, setSinVotacion] = useState(false)
  const [cargadaPara, setCargadaPara] = useState('')
  const [enviando, setEnviando] = useState(false)
  const [errorDeIntegrantes, setErrorDeIntegrantes] = useState('')
  const [errorDeVotacion, setErrorDeVotacion] = useState('')
  const [errorDeAccion, setErrorDeAccion] = useState('')
  const [recarga, setRecarga] = useState(0)
  const envioEnCurso = useRef(false)

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
        if (vigente) setErrorDeIntegrantes(mensajeDeError(causa))
      })

    return () => {
      vigente = false
    }
  }, [proyectoId])

  useEffect(() => {
    let vigente = true

    obtenerHistoria(proyectoId, historiaId)
      .then((datos) => {
        if (vigente) setHistoria(datos)
      })
      .catch(() => {})

    return () => {
      vigente = false
    }
  }, [proyectoId, historiaId])

  useEffect(() => {
    let vigente = true

    consultarVotacion(historiaId, actorId)
      .then((datos) => {
        if (!vigente) return
        setVotacion(datos)
        setSinVotacion(false)
        setErrorDeVotacion('')
      })
      .catch((causa: unknown) => {
        if (!vigente) return
        if (esVotacionNoIniciada(causa)) {
          setVotacion(null)
          setSinVotacion(true)
          setErrorDeVotacion('')
        } else {
          setErrorDeVotacion(mensajeDeError(causa))
        }
      })
      .finally(() => {
        if (vigente) setCargadaPara(consulta)
      })

    return () => {
      vigente = false
    }
  }, [historiaId, actorId, consulta])

  async function ejecutar(accion: () => Promise<Votacion>) {
    if (envioEnCurso.current) return
    envioEnCurso.current = true
    setErrorDeAccion('')
    setEnviando(true)
    try {
      setVotacion(await accion())
      setSinVotacion(false)
    } catch (causa) {
      setErrorDeAccion(mensajeDeError(causa))
      recargar()
    } finally {
      envioEnCurso.current = false
      setEnviando(false)
    }
  }

  const scrumMaster = actor !== undefined && esScrumMaster(actor) ? actor : undefined
  const miVoto = votacion?.integrantes.find((voto) => voto.usuario_id === actorId)?.valor ?? null
  const puedeVotar = votacion !== null && !estaRevelada(votacion) && actor !== undefined

  return (
    <section className="votacion" aria-labelledby="votacion-title">
      <div className="page-header votacion-encabezado">
        <div>
          <p className="eyebrow">Planning Poker</p>
          <h1 id="votacion-title">Votación de la historia #{historiaId}</h1>
          <p className="project-context">
            Proyecto actual: <strong>{proyectoId}</strong>
          </p>
        </div>
        <div className="votacion-encabezado__acciones">
          <Link className="button button--quiet" to={`/projects/${proyectoId}/backlog`}>
            Volver al Product Backlog
          </Link>
        <Button type="button" variant="secondary" onClick={recargar} disabled={enviando || cargando}>
          Actualizar
        </Button>
        </div>
      </div>

      <ResumenHistoria historiaId={historiaId} historia={historia} />

      {errorDeIntegrantes && <ErrorMessage message={errorDeIntegrantes} />}
      {errorDeVotacion && <ErrorMessage message={errorDeVotacion} />}
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

      {sinVotacion && (
        <div className="votacion-acciones">
          <p>La historia todavía no tiene votación</p>
          {scrumMaster && (
            <Button
              type="button"
              disabled={enviando}
              onClick={() => ejecutar(() => iniciarVotacion(historiaId, scrumMaster.id))}
            >
              Iniciar votación
            </Button>
          )}
          {actor && !scrumMaster && <p>Solo el Scrum Master puede iniciar la votación</p>}
          {!actor && integrantes?.length ? <p>Elegí quién actúa para iniciar la votación</p> : null}
        </div>
      )}

      {votacion && (
        <>
          <p className="votacion-ronda">Ronda {votacion.ronda}</p>
          <Integrantes integrantes={votacion.integrantes} />
          {!estaRevelada(votacion) && actor === undefined && integrantes?.length ? (
            <p>Elegí quién actúa para votar</p>
          ) : null}
          {estaRevelada(votacion) && (
            <>
              <Resultado promedio={votacion.promedio} />
              <p>Registrar la estimación acordada se hace en otra historia (HU-14).</p>
            </>
          )}
          {scrumMaster && !estaRevelada(votacion) && (
            <div className="votacion-acciones">
              <Button
                type="button"
                disabled={enviando || !puedeRevelar(votacion, scrumMaster)}
                onClick={() => ejecutar(() => revelarVotacion(historiaId, scrumMaster.id))}
              >
                Revelar
              </Button>
              {faltanPorVotar(votacion).length > 0 && (
                <p>Falta votar: {faltanPorVotar(votacion).map((voto) => voto.nombre).join(', ')}</p>
              )}
            </div>
          )}
          {puedeVotar && <p className="votacion-consigna">Elegí una carta: ¿cuántos Story Points creés que lleva esta historia?</p>}
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
