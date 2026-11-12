import { useState } from 'react'
import type { FormEvent, ReactNode } from 'react'
import { ApiError, ConnectionError } from '../../compartido/api/client'
import { Button } from '../../compartido/componentes/Button'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { Field } from '../../compartido/componentes/Field'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { CampoContrasena } from './CampoContrasena'
import './estilos.css'

type Campo = 'nombre' | 'contrasena'

type Props = {
  titulo: string
  sobretitulo: string
  accion: string
  bajada?: string
  contrasenaNueva?: boolean
  aviso?: string
  pie: ReactNode
  relleno?: { nombre: string; contrasena: string }
  enviar(nombre: string, contrasena: string): Promise<void>
}

function mensajeDeError(error: unknown): string {
  if (error instanceof ApiError || error instanceof ConnectionError) return error.message
  return 'Ocurrió un error inesperado'
}

// Un 409 es «nombre en uso»; el servidor no manda el campo en ese caso.
function campoDeError(error: unknown): Campo | undefined {
  if (!(error instanceof ApiError)) return undefined
  if (error.field === 'nombre' || error.field === 'contrasena') return error.field
  return error.status === 409 ? 'nombre' : undefined
}

export function FormularioCredenciales({
  titulo,
  sobretitulo,
  accion,
  bajada,
  contrasenaNueva = false,
  aviso,
  pie,
  relleno,
  enviar,
}: Props) {
  const [nombre, setNombre] = useState('')
  const [contrasena, setContrasena] = useState('')
  const [errores, setErrores] = useState<Partial<Record<Campo, string>>>({})
  const [errorGeneral, setErrorGeneral] = useState('')
  const [enviando, setEnviando] = useState(false)
  const [rellenoAplicado, setRellenoAplicado] = useState(relleno)

  if (relleno !== rellenoAplicado) {
    setRellenoAplicado(relleno)
    if (relleno) {
      setNombre(relleno.nombre)
      setContrasena(relleno.contrasena)
    }
  }

  async function alEnviar(evento: FormEvent<HTMLFormElement>) {
    evento.preventDefault()
    setErrores({})
    setErrorGeneral('')
    setEnviando(true)
    try {
      await enviar(nombre, contrasena)
    } catch (error) {
      const campo = campoDeError(error)
      if (campo) setErrores({ [campo]: mensajeDeError(error) })
      else setErrorGeneral(mensajeDeError(error))
      setEnviando(false)
    }
  }

  return (
    <section className="auth-pantalla" aria-labelledby="credenciales-title">
      <p className="auth-sobretitulo">{sobretitulo}</p>
      <h1 id="credenciales-title">{titulo}</h1>
      {bajada && <p className="auth-bajada">{bajada}</p>}
      {aviso && <p className="message auth-aviso" role="status">{aviso}</p>}

      <form onSubmit={alEnviar} noValidate aria-label={titulo}>
        {errorGeneral && <ErrorMessage message={errorGeneral} />}

        <Field
          name="nombre"
          label="Nombre de usuario"
          autoComplete="username"
          ayuda={contrasenaNueva ? 'De 3 a 30 caracteres: letras sin tildes, números, punto, guion y guion bajo.' : undefined}
          value={nombre}
          onChange={(evento) => setNombre(evento.target.value)}
          error={errores.nombre}
        />
        <CampoContrasena
          name="contrasena"
          label="Contraseña"
          autoComplete={contrasenaNueva ? 'new-password' : 'current-password'}
          ayuda={contrasenaNueva ? 'De 8 a 128 caracteres.' : undefined}
          value={contrasena}
          onChange={(evento) => setContrasena(evento.target.value)}
          error={errores.contrasena}
        />

        {enviando && <LoadingIndicator />}

        <Button type="submit" disabled={enviando}>
          {accion}
        </Button>
      </form>

      <p>{pie}</p>
    </section>
  )
}
