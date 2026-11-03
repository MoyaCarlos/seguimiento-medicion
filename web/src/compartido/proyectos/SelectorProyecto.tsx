import type { ChangeEvent } from 'react'
import { Link, useLocation, useNavigate, useParams } from 'react-router'
import { etiquetaProyecto, rutaDeSeccion } from './ruta'
import { useProyectos } from './useProyectos'

export function SelectorProyecto() {
  const { id = '' } = useParams()
  const { pathname } = useLocation()
  const navegar = useNavigate()
  const proyectos = useProyectos()

  const listados = proyectos.estado === 'listo' ? proyectos.proyectos : []
  const actualListado = listados.some((proyecto) => String(proyecto.id) === id)

  function elegir(evento: ChangeEvent<HTMLSelectElement>) {
    if (evento.target.value === id) return
    navegar(rutaDeSeccion(pathname, evento.target.value))
  }

  return (
    <div className="selector-proyecto">
      <label className="field">
        <span>Proyecto</span>
        <select value={id} onChange={elegir}>
          {!actualListado && <option value={id}>{`Proyecto ${id}`}</option>}
          {listados.map((proyecto) => (
            <option key={proyecto.id} value={proyecto.id}>{etiquetaProyecto(proyecto)}</option>
          ))}
        </select>
      </label>
      <Link className="button button--secondary" to="/">Nuevo proyecto</Link>
    </div>
  )
}
