import { useEffect, useState } from 'react'
import { matchRoute } from './routes'

function FeaturePlaceholder({ title, projectId }: { title: string; projectId?: string }) {
  return (
    <section className="placeholder" aria-labelledby="placeholder-title">
      <p className="eyebrow">Base del frontend · Sprint 2</p>
      <h1 id="placeholder-title">{title}</h1>
      {projectId && <p>Proyecto actual: {projectId}</p>}
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
    return <FeaturePlaceholder title="Página no encontrada" />
  }

  const title = route.feature === 'hu04'
    ? route.screen === 'settings' ? 'Configuración del equipo' : route.screen === 'edit' ? 'Edición del proyecto' : 'Panel del proyecto'
    : route.feature === 'hu01' ? 'Product Backlog'
      : route.feature === 'hu05' ? 'Sprints'
        : 'Estado del proyecto'

  return <FeaturePlaceholder title={title} projectId={route.projectId} />
}
