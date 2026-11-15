import { useState } from 'react'
import type { FormEvent } from 'react'
import { Link } from 'react-router'
import { ApiError, ConnectionError } from '../../compartido/api/client'
import { Button } from '../../compartido/componentes/Button'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { Field } from '../../compartido/componentes/Field'
import { Icono } from '../../compartido/componentes/Icono'
import { useSesion } from './SesionContext'
import './estilos.css'

function mensajeDeError(error: unknown): string {
  if (error instanceof ApiError || error instanceof ConnectionError) return error.message
  return 'Ocurrió un error inesperado'
}

export function MenuUsuario() {
  const { estado, usuario, cambiarNombre, salir } = useSesion()
  const [cambiando, setCambiando] = useState(false)
  const [nombre, setNombre] = useState('')
  const [errorNombre, setErrorNombre] = useState('')
  const [errorGeneral, setErrorGeneral] = useState('')
  const [enviando, setEnviando] = useState(false)

  if (estado === 'autenticado' && usuario) {
    const nombreActual = usuario.nombre
    function abrir() {
      setNombre(nombreActual)
      setErrorNombre('')
      setErrorGeneral('')
      setCambiando(true)
    }

    async function guardar(evento: FormEvent<HTMLFormElement>) {
      evento.preventDefault()
      if (enviando) return
      setErrorNombre('')
      setErrorGeneral('')
      setEnviando(true)
      try {
        await cambiarNombre(nombre)
        setCambiando(false)
      } catch (error) {
        if (error instanceof ApiError && error.field === 'nombre') {
          setErrorNombre(mensajeDeError(error))
        } else {
          setErrorGeneral(mensajeDeError(error))
        }
      } finally {
        setEnviando(false)
      }
    }

    return (
      <div className="menu-usuario">
        <span className="menu-usuario-avatar" aria-hidden="true">
          {usuario.nombre?.charAt(0).toUpperCase()}
        </span>
        <span className="menu-usuario-nombre">{usuario.nombre}</span>
        <button type="button" onClick={abrir} aria-expanded={cambiando}>
          <Icono nombre="lapiz" />
          <span className="menu-usuario-accion">Cambiar nombre</span>
        </button>
        {cambiando && (
          <form
            className="menu-usuario-cambio"
            onSubmit={guardar}
            noValidate
            aria-label="Cambiar nombre de usuario"
          >
            {errorGeneral && <ErrorMessage message={errorGeneral} />}
            <Field
              id="menu-usuario-nombre"
              name="nombre"
              label="Nuevo nombre de usuario"
              autoComplete="username"
              value={nombre}
              onChange={(evento) => {
                setNombre(evento.target.value)
                setErrorNombre('')
                setErrorGeneral('')
              }}
              disabled={enviando}
              error={errorNombre}
            />
            <div className="menu-usuario-cambio-acciones">
              <Button
                type="button"
                variant="secondary"
                onClick={() => setCambiando(false)}
                disabled={enviando}
              >
                Cancelar
              </Button>
              <Button type="submit" disabled={enviando}>
                Guardar
              </Button>
            </div>
          </form>
        )}
        <button type="button" onClick={() => void salir()}>
          <Icono nombre="salir" />
          <span className="menu-usuario-accion">Cerrar sesión</span>
        </button>
      </div>
    )
  }
  if (estado === 'anonimo') {
    return (
      <nav className="menu-usuario">
        <Link to="/ingresar">Ingresar</Link>
        <Link to="/registro">Registrarse</Link>
      </nav>
    )
  }
  return null
}
