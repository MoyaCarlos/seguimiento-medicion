import { render, screen } from '@testing-library/react'
import { createMemoryRouter } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { rutas } from './rutas'

beforeEach(() => {
  vi.stubGlobal(
    'fetch',
    vi.fn(() => Promise.resolve(new Response('[]', { status: 200 }))),
  )
})

afterEach(() => {
  vi.unstubAllGlobals()
})

function abrir(url: string) {
  const router = createMemoryRouter(rutas, { initialEntries: [url] })
  render(<RouterProvider router={router} />)
  return router
}

describe('registro de rutas', () => {
  it.each([
    ['/projects/42/proyecto', 'Panel del proyecto'],
    ['/projects/42/backlog', 'Product Backlog'],
    ['/projects/42/backlog/nueva', 'Nueva historia'],
    ['/projects/42/sprints', 'Sprints'],
    ['/projects/42/status', 'Estado del proyecto'],
  ])('%s muestra su pantalla con el proyecto de la URL', async (url, titulo) => {
    abrir(url)

    expect(await screen.findByRole('heading', { name: titulo })).toBeTruthy()
    expect(screen.getByText('42')).toBeTruthy()
  })

  it('/projects/:id/votacion/:historiaId muestra la votación de la historia', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn((url: string) =>
        Promise.resolve(
          new Response(
            url.includes('/votacion')
              ? JSON.stringify({ historia_id: 3, ronda: 1, revelada_en: null, integrantes: [], promedio: null })
              : '[]',
            { status: 200 },
          ),
        ),
      ),
    )
    abrir('/projects/42/votacion/3')

    expect(await screen.findByRole('heading', { name: 'Votación de la historia #3' })).toBeTruthy()
    expect(screen.getByText('42')).toBeTruthy()
  })

  it('/projects/:id/dashboard muestra el Dashboard y la navegación lo ofrece como activo', async () => {
    const metricas = {
      story_points_planificados: 0,
      story_points_completados: 0,
      velocidad: 0,
      porcentaje_historias_completadas: 0,
      horas_estimadas: null,
      horas_reales: 0,
      desviacion_horas: null,
      desviacion_porcentaje: null,
      defectos_detectados: 0,
      defectos_resueltos: 0,
    }
    vi.stubGlobal(
      'fetch',
      vi.fn(() =>
        Promise.resolve(
          new Response(
            JSON.stringify({ proyecto_id: 42, horas_por_story_point: null, proyecto: metricas, sprints: [] }),
            { status: 200 },
          ),
        ),
      ),
    )
    abrir('/projects/42/dashboard')

    expect(await screen.findByRole('heading', { name: 'Dashboard' })).toBeTruthy()
    expect(screen.getByRole('link', { name: 'Dashboard' }).getAttribute('aria-current')).toBe('page')
  })

  it('la pantalla de una feature sigue mostrándose en sus sub-rutas', async () => {
    abrir('/projects/1/proyecto/settings')

    expect(await screen.findByRole('heading', { name: 'Panel del proyecto' })).toBeTruthy()
  })

  it('la pantalla principal permite crear un proyecto', async () => {
    abrir('/')

    expect(await screen.findByRole('heading', { name: 'Crear proyecto' })).toBeTruthy()
  })

  it('/projects/:id redirige al panel del proyecto', async () => {
    const router = abrir('/projects/1')

    expect(await screen.findByRole('heading', { name: 'Panel del proyecto' })).toBeTruthy()
    expect(router.state.location.pathname).toBe('/projects/1/proyecto')
  })

  it('marca como activa solo la pantalla actual en la navegación', async () => {
    abrir('/projects/42/backlog')
    await screen.findByRole('heading', { name: 'Product Backlog' })

    expect(screen.getByRole('link', { name: 'Historias' }).getAttribute('aria-current')).toBe('page')
    expect(screen.getByRole('link', { name: 'Sprints' }).getAttribute('aria-current')).toBeNull()
  })

  it('mantiene activo el enlace de la feature en sus sub-rutas', async () => {
    abrir('/projects/1/proyecto/edit')
    await screen.findByRole('heading', { name: 'Panel del proyecto' })

    expect(screen.getByRole('link', { name: 'Proyecto' }).getAttribute('aria-current')).toBe('page')
  })

  it('cambiar el id de la URL cambia el proyecto mostrado, sin estado compartido', async () => {
    const router = abrir('/projects/1/backlog')
    expect(await screen.findByText('1', { selector: 'strong' })).toBeTruthy()

    await router.navigate('/projects/2/backlog')

    expect(await screen.findByText('2', { selector: 'strong' })).toBeTruthy()
  })

  it.each(['/no-existe', '/projects/1/otra'])('%s muestra "Página no encontrada"', async (url) => {
    abrir(url)

    expect(await screen.findByRole('heading', { name: 'Página no encontrada' })).toBeTruthy()
  })
})
