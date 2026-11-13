import { describe, expect, it } from 'vitest'
import { seccionDeRuta } from './navegacion'

describe('seccionDeRuta', () => {
  it.each([
    ['/projects/42/proyecto', 'Proyecto'],
    ['/projects/42/backlog', 'Historias'],
    ['/projects/42/backlog/nueva', 'Nueva historia'],
    ['/projects/42/votacion/3', 'Planning Poker'],
    ['/projects/42/sprints', 'Sprints'],
    ['/projects/42/defectos', 'Defectos'],
    ['/projects/42/status', 'Estado'],
    ['/projects/42/dashboard', 'Dashboard'],
    ['/projects/42/historial', 'Historial de Sprints'],
    ['/projects/42/historial/7', 'Historial de Sprints'],
    ['/projects/42/proyecto/settings', 'Equipo'],
    ['/projects/42/proyecto/edit', 'Editar proyecto'],
  ])('%s → «%s»', (ruta, texto) => {
    expect(seccionDeRuta(ruta)).toBe(texto)
  })

  it.each([['/projects/42'], ['/projects/42/otra'], ['/']])('%s no tiene sección', (ruta) => {
    expect(seccionDeRuta(ruta)).toBeUndefined()
  })
})
