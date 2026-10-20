import { describe, expect, it } from 'vitest'
import { featureRoutes, matchRoute, navigationRoutes } from './routes'

describe('frontend route registry', () => {
  it('extracts the project id from the project route', () => {
    expect(matchRoute('/projects/42')).toEqual({
      feature: 'hu04',
      projectId: '42',
      screen: 'dashboard',
    })
  })

  it('registers one route entry for each Sprint 2 feature', () => {
    expect(featureRoutes.map((route) => route.feature)).toEqual(['hu04', 'hu01', 'hu05', 'hu13'])
  })

  it('registers one navigation entry for each Sprint 2 feature', () => {
    expect(navigationRoutes.map((route) => route.feature)).toEqual(['hu04', 'hu01', 'hu05', 'hu13'])
  })

  it('returns a fallback for an unknown route', () => {
    expect(matchRoute('/unknown')).toEqual({ feature: 'not-found' })
  })
})
