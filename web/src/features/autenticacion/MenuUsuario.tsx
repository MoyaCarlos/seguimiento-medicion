import { useEffect, useId, useRef, useState } from 'react'
import { Link } from 'react-router'
import { Icono } from '../../compartido/componentes/Icono'
import { useSesion } from './SesionContext'
import './estilos.css'

export function MenuUsuario() {
  const { estado, usuario, salir } = useSesion()
  const [abierto, setAbierto] = useState(false)
  const raiz = useRef<HTMLDivElement>(null)
  const boton = useRef<HTMLButtonElement>(null)
  const idPanel = useId()

  useEffect(() => {
    if (!abierto) return
    const alPulsar = (evento: PointerEvent) => {
      if (!raiz.current?.contains(evento.target as Node)) setAbierto(false)
    }
    const alTeclear = (evento: KeyboardEvent) => {
      if (evento.key !== 'Escape') return
      setAbierto(false)
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

    return (
      <div className="menu-usuario" ref={raiz}>
        <button
          ref={boton}
          type="button"
          className="menu-usuario__boton"
          aria-label="Menú de usuario"
          aria-expanded={abierto}
          aria-controls={idPanel}
          onClick={() => setAbierto((actual) => !actual)}
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
                <small>Sesión iniciada</small>
              </div>
            </div>
            <hr />
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
