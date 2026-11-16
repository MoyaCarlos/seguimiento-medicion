import { Link } from 'react-router'
import { Icono } from './Icono'

export function Marca() {
  return (
    <Link className="marca" to="/">
      <span className="marca__logo"><Icono nombre="marca" /></span>
      Métrica
    </Link>
  )
}
