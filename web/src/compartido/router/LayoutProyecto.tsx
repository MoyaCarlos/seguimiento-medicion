import { NavLink, Outlet, useLocation, useParams } from 'react-router'
import { navegacion } from '../../navegacion'
import { SelectorProyecto } from '../proyectos/SelectorProyecto'

export function LayoutProyecto() {
  const { pathname } = useLocation()
  const { id } = useParams()

  return (
    <div className="page-stack">
      <div className="page-header project-header">
        <div>
          <p className="eyebrow">Métrica · Seguimiento de proyectos</p>
          <p className="route-context">Ruta actual: <code>{pathname}</code></p>
        </div>
        <div className="project-bar">
          <SelectorProyecto />
          <nav className="feature-navigation" aria-label="Secciones del proyecto">
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
      </div>
      <Outlet key={id} />
    </div>
  )
}
