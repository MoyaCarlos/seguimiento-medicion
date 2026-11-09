import { useEffect, useState } from 'react'
import { Link } from 'react-router'
import { ApiError } from '../api/client'
import { Button } from '../componentes/Button'
import { ErrorMessage } from '../componentes/ErrorMessage'
import { LoadingIndicator } from '../componentes/LoadingIndicator'
import { PROYECTO_DE_PRUEBA, cargarDatosDePrueba, datosDePruebaHabilitados } from './api'
import { CredencialesDemo } from './CredencialesDemo'

type Resultado = { tipo: 'aviso'; mensaje: string; proyecto: boolean } | { tipo: 'error'; mensaje: string }

type Props = { alUsar?: (nombre: string, contrasena: string) => void }

export function BotonDatosDePrueba({ alUsar }: Props) {
  const [habilitado, setHabilitado] = useState(false)
  const [enviando, setEnviando] = useState(false)
  const [resultado, setResultado] = useState<Resultado>()

  useEffect(() => {
    let vigente = true
    datosDePruebaHabilitados().then((respuesta) => {
      if (vigente) setHabilitado(respuesta)
    })
    return () => {
      vigente = false
    }
  }, [])

  async function cargar() {
    setResultado(undefined)
    setEnviando(true)
    try {
      await cargarDatosDePrueba()
      setResultado({ tipo: 'aviso', mensaje: 'Datos de prueba cargados.', proyecto: true })
    } catch (error) {
      setResultado(aResultado(error))
    } finally {
      setEnviando(false)
    }
  }

  if (!habilitado) return null

  return (
    <div>
      {resultado?.tipo === 'error' && <ErrorMessage message={resultado.mensaje} />}
      {resultado?.tipo === 'aviso' && (
        <p className="message" role="status">
          {resultado.mensaje}{' '}
          {resultado.proyecto && <Link to={`/projects/${PROYECTO_DE_PRUEBA}`}>Abrir el proyecto de prueba</Link>}
        </p>
      )}
      {resultado?.tipo === 'aviso' && resultado.proyecto && <CredencialesDemo alUsar={alUsar} />}
      {enviando && <LoadingIndicator />}
      <Button type="button" variant="secondary" onClick={cargar} disabled={enviando}>
        Cargar datos de prueba
      </Button>
    </div>
  )
}

function aResultado(error: unknown): Resultado {
  if (error instanceof ApiError && error.status === 409 && error.field === 'ya_cargado') {
    return { tipo: 'aviso', mensaje: error.message, proyecto: true }
  }
  if (error instanceof ApiError && error.status === 409 && error.field === 'usuario_existente') {
    return { tipo: 'aviso', mensaje: error.message, proyecto: false }
  }
  return { tipo: 'error', mensaje: error instanceof Error ? error.message : 'No se pudieron cargar los datos de prueba' }
}
