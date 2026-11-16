import { Icono } from '../../compartido/componentes/Icono'
import { useSesion } from './SesionContext'
import './estilos.css'

/** Pie de la barra lateral: quién tiene la sesión, su rol en el proyecto (si hay uno) y la salida. */
export function BloqueUsuario({ rol }: { rol?: string }) {
  const { estado, usuario, salir } = useSesion()

  if (estado !== 'autenticado' || !usuario) return null

  return (
    <div className="bloque-usuario">
      <span className="bloque-usuario__avatar" aria-hidden="true">
        {usuario.nombre?.charAt(0).toUpperCase()}
      </span>
      <div>
        <strong>{usuario.nombre}</strong>
        <small>{rol ?? 'Sesión iniciada'}</small>
      </div>
      <button
        type="button"
        className="bloque-usuario__salir"
        aria-label="Cerrar sesión"
        title="Cerrar sesión"
        onClick={() => void salir()}
      >
        <Icono nombre="salir" />
      </button>
    </div>
  )
}
