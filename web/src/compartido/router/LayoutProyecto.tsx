import { NavLink, Outlet, useLocation } from 'react-router'
import { navegacion } from '../../navegacion'

export function LayoutProyecto() {
  const { pathname } = useLocation()

  return (
    <div className="page-stack">
      <div className="page-header">
        <div>
          <p className="eyebrow">Métrica · Seguimiento de proyectos</p>
          <p className="route-context">Ruta actual: <code>{pathname}</code></p>
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
      <Outlet />
    </div>
  )
}
