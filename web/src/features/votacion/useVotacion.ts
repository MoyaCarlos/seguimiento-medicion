import { useCallback, useEffect, useRef, useState } from 'react'
import { consultarVotacion, iniciarVotacion, listarIntegrantes, obtenerHistoria, revelarVotacion, votar } from './api'
import { esVotacionNoIniciada, mensajeDeError } from './errores'
import { esScrumMaster, estaRevelada, sinValores, type HistoriaResumen, type Integrante, type Votacion } from './votacion'

export function useVotacion(proyectoId: string, historiaId: number) {
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
  const [actualizacion, setActualizacion] = useState(0)
  const envioEnCurso = useRef(false)
  const accionesIniciadas = useRef(0)

  const consulta = `${historiaId}/${actorId}/${recarga}`
  const cargando = cargadaPara !== consulta
  const actor = integrantes?.find((integrante) => integrante.id === actorId)
  const scrumMaster = actor !== undefined && esScrumMaster(actor) ? actor : undefined
  const reconsultar = useCallback(() => setRecarga((actual) => actual + 1), [])
  const recargar = useCallback(() => {
    setErrorDeAccion('')
    setActualizacion((actual) => actual + 1)
    reconsultar()
  }, [reconsultar])

  useEffect(() => {
    let vigente = true

    listarIntegrantes(proyectoId)
      .then((datos) => {
        if (!vigente) return
        setIntegrantes(datos)
        setErrorDeIntegrantes('')
      })
      .catch((causa: unknown) => {
        if (vigente) setErrorDeIntegrantes(mensajeDeError(causa))
      })

    return () => {
      vigente = false
    }
  }, [proyectoId, actualizacion])

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
  }, [proyectoId, historiaId, actualizacion])

  useEffect(() => {
    let vigente = true
    const accionesAlConsultar = accionesIniciadas.current
    const sigueVigente = () => vigente && accionesAlConsultar === accionesIniciadas.current

    consultarVotacion(historiaId, actorId)
      .then((datos) => {
        if (!sigueVigente()) return
        setVotacion(datos)
        setSinVotacion(false)
        setErrorDeVotacion('')
      })
      .catch((causa: unknown) => {
        if (!sigueVigente()) return
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
    accionesIniciadas.current += 1
    setErrorDeAccion('')
    setEnviando(true)
    try {
      setVotacion(await accion())
      setSinVotacion(false)
    } catch (causa) {
      setErrorDeAccion(mensajeDeError(causa))
      reconsultar()
    } finally {
      envioEnCurso.current = false
      setEnviando(false)
    }
  }

  function elegirActor(id: number | undefined) {
    setErrorDeAccion('')
    setVotacion((actual) => (actual && !estaRevelada(actual) ? sinValores(actual) : actual))
    setActorId(id)
  }

  return {
    historia,
    integrantes,
    actor,
    scrumMaster,
    votacion,
    sinVotacion,
    cargando,
    enviando,
    errorDeIntegrantes,
    errorDeVotacion,
    errorDeAccion,
    recargar,
    elegirActor,
    iniciar: () => (scrumMaster ? ejecutar(() => iniciarVotacion(historiaId, scrumMaster.id)) : undefined),
    revelar: () => (scrumMaster ? ejecutar(() => revelarVotacion(historiaId, scrumMaster.id)) : undefined),
    votar: (valor: number) => (actor ? ejecutar(() => votar(historiaId, actor.id, valor)) : undefined),
  }
}
