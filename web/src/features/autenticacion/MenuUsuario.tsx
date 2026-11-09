import { Link } from 'react-router'
import { useSesion } from './SesionContext'
import './estilos.css'

export function MenuUsuario() {
  const { estado, usuario, salir } = useSesion()

  if (estado === 'autenticado' && usuario) {
    return (
      <div className="menu-usuario">
        <span className="menu-usuario-nombre">{usuario.nombre}</span>
        <button type="button" onClick={() => void salir()}>Cerrar sesión</button>
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
