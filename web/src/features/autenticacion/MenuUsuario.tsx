import { useEffect, useId, useRef, useState } from 'react'
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

export function MenuUsuario({ rol }: { rol?: string }) {
  const { estado, usuario, cambiarNombre, salir } = useSesion()
  const [abierto, setAbierto] = useState(false)
  const [cambiando, setCambiando] = useState(false)
  const [nombre, setNombre] = useState('')
  const [errorNombre, setErrorNombre] = useState('')
  const [errorGeneral, setErrorGeneral] = useState('')
  const [enviando, setEnviando] = useState(false)
  const raiz = useRef<HTMLDivElement>(null)
  const boton = useRef<HTMLButtonElement>(null)
  const idPanel = useId()

  function cerrarMenu() {
    setAbierto(false)
    setCambiando(false)
  }

  useEffect(() => {
    if (!abierto) return
    const alPulsar = (evento: PointerEvent) => {
      if (!raiz.current?.contains(evento.target as Node)) cerrarMenu()
    }
    const alTeclear = (evento: KeyboardEvent) => {
      if (evento.key !== 'Escape') return
      cerrarMenu()
      boton.current?.focus()
    }
    document.addEventListener('pointerdown', alPulsar)
    document.addEventListener('keydown', alTeclear)
    return () => {
      document.removeEventListener('pointerdown', alPulsar)
      document.removeEventListener('keydown', alTeclear)
    }
  }, [abierto])

  if (estado === 'autenticado' && usuario) {
    const inicial = usuario.nombre?.charAt(0).toUpperCase()
    const nombreActual = usuario.nombre

    function abrirCambio() {
      setNombre(nombreActual)
      setErrorNombre('')
      setErrorGeneral('')
      setCambiando(true)
    }

    async function guardarNombre(evento: FormEvent<HTMLFormElement>) {
      evento.preventDefault()
      if (enviando) return
      setErrorNombre('')
      setErrorGeneral('')
      setEnviando(true)
      try {
        await cambiarNombre(nombre)
        setCambiando(false)
        setAbierto(false)
        boton.current?.focus()
      } catch (error) {
        if (error instanceof ApiError && error.field === 'nombre') setErrorNombre(mensajeDeError(error))
        else setErrorGeneral(mensajeDeError(error))
      } finally {
        setEnviando(false)
      }
    }

    return (
      <div className="menu-usuario" ref={raiz}>
        <button
          ref={boton}
          type="button"
          className="menu-usuario__boton"
          aria-label="Menú de usuario"
          aria-expanded={abierto}
          aria-controls={idPanel}
          onClick={() => (abierto ? cerrarMenu() : setAbierto(true))}
        >
          <span className="menu-usuario__avatar" aria-hidden="true">{inicial}</span>
          <span className="menu-usuario__nombre">{usuario.nombre}</span>
          <Icono nombre="chevronAbajo" />
        </button>
        {abierto && (
          <div id={idPanel} className="menu-usuario__panel" role="group" aria-label="Cuenta">
            <div className="menu-usuario__cabecera">
              <span className="menu-usuario__avatar" aria-hidden="true">{inicial}</span>
              <div>
                <strong>{usuario.nombre}</strong>
                <small>{rol ?? 'Sesión iniciada'}</small>
              </div>
            </div>
            {rol && <p className="menu-usuario__rol">{`Tu rol en este proyecto: ${rol}`}</p>}
            <hr />
            {cambiando ? (
              <form
                className="menu-usuario__cambio"
                onSubmit={guardarNombre}
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
                <div className="menu-usuario__cambio-acciones">
                  <Button type="button" variant="quiet" onClick={() => setCambiando(false)} disabled={enviando}>
                    Cancelar
                  </Button>
                  <Button type="submit" disabled={enviando}>
                    Guardar
                  </Button>
                </div>
              </form>
            ) : (
              <button type="button" className="menu-usuario__item" onClick={abrirCambio}>
                <Icono nombre="lapiz" />
                Cambiar nombre
              </button>
            )}
            <button
              type="button"
              className="menu-usuario__item menu-usuario__item--peligro"
              onClick={() => void salir()}
            >
              <Icono nombre="salir" />
              Cerrar sesión
            </button>
          </div>
        )}
      </div>
    )
  }
  if (estado === 'anonimo') {
    return (
      <nav className="menu-usuario menu-usuario--anonimo">
        <Link to="/ingresar">Ingresar</Link>
        <Link to="/registro">Registrarse</Link>
      </nav>
    )
  }
  return null
}
