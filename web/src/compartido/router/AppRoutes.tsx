import { useEffect, useState } from 'react'
import { matchRoute, navigationRoutes } from './routes'

function FeaturePlaceholder({ title, projectId }: { title: string; projectId?: string }) {
  return (
    <section className="placeholder placeholder--feature" aria-labelledby="placeholder-title">
      <div className="placeholder-icon" aria-hidden="true">✦</div>
      <div>
        <p className="eyebrow">Base del frontend · Sprint 2</p>
        <h1 id="placeholder-title">{title}</h1>
        {projectId && <p className="project-context">Proyecto actual: <strong>{projectId}</strong></p>}
      </div>
      <p>Esta pantalla será implementada por su historia de usuario.</p>
    </section>
  )
}

export function AppRoutes() {
  const [pathname, setPathname] = useState(window.location.pathname)
  const route = matchRoute(pathname)

  useEffect(() => {
    const updatePath = () => setPathname(window.location.pathname)
    window.addEventListener('popstate', updatePath)
    return () => window.removeEventListener('popstate', updatePath)
  }, [])

  if (route.feature === 'not-found') {
    return (
      <div className="page-stack">
        <PageHeader pathname={pathname} projectId={route.projectId} />
        <FeaturePlaceholder title="Página no encontrada" />
      </div>
    )
  }

  const title = route.feature === 'hu04'
    ? route.screen === 'settings' ? 'Configuración del equipo' : route.screen === 'edit' ? 'Edición del proyecto' : 'Panel del proyecto'
    : route.feature === 'hu01' ? 'Product Backlog'
      : route.feature === 'hu05' ? 'Sprints'
        : 'Estado del proyecto'

  return (
    <div className="page-stack">
      <PageHeader pathname={pathname} projectId={route.projectId} />
      <FeaturePlaceholder title={title} projectId={route.projectId} />
    </div>
  )
}

function PageHeader({ pathname, projectId }: { pathname: string; projectId?: string }) {
  const currentProjectId = projectId ?? pathname.match(/^\/projects\/([^/]+)/)?.[1]

  function navigate(path: string) {
    if (!currentProjectId) return
    const nextPath = `/projects/${currentProjectId}${path}`
    window.history.pushState({}, '', nextPath)
    window.dispatchEvent(new PopStateEvent('popstate'))
  }

  return (
    <div className="page-header">
      <div>
        <p className="eyebrow">Métrica · Seguimiento de proyectos</p>
        <p className="route-context">Ruta actual: <code>{pathname}</code></p>
      </div>
      {currentProjectId && (
        <nav className="feature-navigation" aria-label="Historias del Sprint 2">
          {navigationRoutes.map((item) => {
            const active = matchRoute(pathname).feature === item.feature
            return (
              <button
                className={`feature-navigation__item${active ? ' feature-navigation__item--active' : ''}`}
                key={item.feature}
                onClick={() => navigate(item.path)}
                type="button"
                aria-current={active ? 'page' : undefined}
              >
                {item.label}
              </button>
            )
          })}
        </nav>
      )}
    </div>
  )
}
