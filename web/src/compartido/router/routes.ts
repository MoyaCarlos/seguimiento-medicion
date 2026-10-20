export type FeatureName = 'hu04' | 'hu01' | 'hu05' | 'hu13' | 'not-found'
export type RouteMatch = {
  feature: FeatureName
  projectId?: string
  screen?: string
}

export const featureRoutes = [
  { feature: 'hu04' as const, paths: ['/projects/:id', '/projects/:id/settings', '/projects/:id/edit'] },
  { feature: 'hu01' as const, paths: ['/projects/:id/backlog'] },
  { feature: 'hu05' as const, paths: ['/projects/:id/sprints'] },
  { feature: 'hu13' as const, paths: ['/projects/:id/status'] },
]

export function matchRoute(pathname: string): RouteMatch {
  const segments = pathname.split('/').filter(Boolean)
  if (segments[0] !== 'projects' || !segments[1]) return { feature: 'not-found' }

  const projectId = segments[1]
  const suffix = segments[2]
  if (!suffix) return { feature: 'hu04', projectId, screen: 'dashboard' }
  if (suffix === 'settings' || suffix === 'edit') return { feature: 'hu04', projectId, screen: suffix }
  if (suffix === 'backlog') return { feature: 'hu01', projectId, screen: 'backlog' }
  if (suffix === 'sprints') return { feature: 'hu05', projectId, screen: 'sprints' }
  if (suffix === 'status') return { feature: 'hu13', projectId, screen: 'status' }
  return { feature: 'not-found' }
}
