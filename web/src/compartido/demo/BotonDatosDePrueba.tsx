import { useEffect, useState } from 'react'
import { ApiError } from '../api/client'
import { Button } from '../componentes/Button'
import { ErrorMessage } from '../componentes/ErrorMessage'
import { LoadingIndicator } from '../componentes/LoadingIndicator'
import { PROYECTO_DE_PRUEBA, cargarDatosDePrueba, datosDePruebaHabilitados } from './api'
import { CredencialesDemo } from './CredencialesDemo'

type Resultado = { tipo: 'aviso'; mensaje: string; conTabla: boolean; exito?: boolean } | { tipo: 'error'; mensaje: string }

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
      setResultado({
        tipo: 'aviso',
        mensaje: 'Datos de prueba cargados: ya podés ingresar con los usuarios de la lista',
        conTabla: true,
        exito: true,
      })
    } catch (error) {
      setResultado(aResultado(error))
    } finally {
      setEnviando(false)
    }
  }

  if (!habilitado) return null

  return (
    <div className="demo-tarjeta">
      <h2>Datos de prueba</h2>
      <p>
        Crea cuatro usuarios y el proyecto #{PROYECTO_DE_PRUEBA} «Licencia de conducir (datos de prueba)»: 4 Sprints
        (dos finalizados, uno activo y uno pendiente), 12 historias, 20 registros de horas, 2 votaciones y 5 defectos.
        Si algún usuario ya existe, no se carga nada.
      </p>
      {resultado?.tipo === 'error' && <ErrorMessage className="auth-alerta" message={resultado.mensaje} />}
      {resultado?.tipo === 'aviso' && (
        <p className={resultado.exito ? 'demo-exito' : 'auth-alerta'} role="status">{resultado.mensaje}</p>
      )}
      {resultado?.tipo === 'aviso' && resultado.conTabla && <CredencialesDemo alUsar={alUsar} />}
      {enviando && <LoadingIndicator />}
      <Button type="button" variant="secondary" onClick={cargar} disabled={enviando}>
        Cargar datos de prueba
      </Button>
    </div>
  )
}

function aResultado(error: unknown): Resultado {
  if (error instanceof ApiError && error.status === 409 && error.field === 'ya_cargado') {
    return { tipo: 'aviso', mensaje: error.message, conTabla: true }
  }
  if (error instanceof ApiError && error.status === 409 && error.field === 'usuario_existente') {
    return { tipo: 'aviso', mensaje: error.message, conTabla: false }
  }
  return { tipo: 'error', mensaje: error instanceof Error ? error.message : 'No se pudieron cargar los datos de prueba' }
}
