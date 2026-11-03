import { NavLink, Outlet, useLocation, useParams } from 'react-router'
import { navegacion } from '../../navegacion'
import { SelectorProyecto } from '../proyectos/SelectorProyecto'

export function LayoutProyecto() {
  const { pathname } = useLocation()
  const { id } = useParams()

  return (
    <div className="page-stack">
      <div className="page-header">
        <div>
          <p className="eyebrow">Métrica · Seguimiento de proyectos</p>
          <p className="route-context">Ruta actual: <code>{pathname}</code></p>
          <SelectorProyecto />
        </div>
        <nav className="feature-navigation" aria-label="Historias del Sprint 2">
          {navegacion.map(({ ruta, texto }) => (
            <NavLink
              key={ruta}
              to={ruta}
              className={({ isActive }) =>
                `feature-navigation__item${isActive ? ' feature-navigation__item--active' : ''}`}
            >
              {texto}
            </NavLink>
          ))}
        </nav>
      </div>
      <Outlet key={id} />
    </div>
  )
}
