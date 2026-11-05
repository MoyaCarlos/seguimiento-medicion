import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { createMemoryRouter } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { borrarToken, guardarToken } from './features/autenticacion/almacenamientoToken'
import { rutas } from './rutas'

const USUARIO = { id: 7, nombre: 'ana.perez' }

const VOTACION = { historia_id: 3, ronda: 1, revelada_en: null, integrantes: [], promedio: null }

function responder(url: string) {
  if (url.includes('/auth/yo')) return JSON.stringify(USUARIO)
  if (url.includes('/votacion')) return JSON.stringify(VOTACION)
  return '[]'
}

function simularFetch() {
  vi.stubGlobal(
    'fetch',
    vi.fn((url: string) => Promise.resolve(new Response(responder(url), { status: 200 }))),
  )
}

// Sesión simulada: token guardado y `GET /auth/yo` responde 200.
function simularSesion() {
  guardarToken('token-de-prueba')
  simularFetch()
}

beforeEach(() => {
  borrarToken()
  localStorage.clear()
  simularSesion()
})

afterEach(() => {
  borrarToken()
  localStorage.clear()
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
    abrir('/projects/42/votacion/3')

    expect(await screen.findByRole('heading', { name: 'Votación de la historia #3' })).toBeTruthy()
    expect(screen.getByText('42')).toBeTruthy()
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

  it('con sesión el encabezado muestra el menú de usuario', async () => {
    abrir('/')

    expect(await screen.findByText('ana.perez')).toBeTruthy()
    expect(screen.getByRole('button', { name: 'Cerrar sesión' })).toBeTruthy()
  })
})

describe('rutas protegidas', () => {
  // Todas las rutas que declara rutas.tsx, con una URL concreta de cada una.
  const protegidas = [
    '/',
    '/projects/42',
    '/projects/42/proyecto',
    '/projects/42/proyecto/settings',
    '/projects/42/backlog',
    '/projects/42/backlog/nueva',
    '/projects/42/votacion/3',
    '/projects/42/sprints',
    '/projects/42/status',
    '/projects/42/otra',
    '/no-existe',
  ]

  it.each(protegidas)('sin sesión %s redirige a /ingresar conservando la ruta en state.desde', async (url) => {
    borrarToken()
    const router = abrir(url)

    expect(await screen.findByRole('heading', { name: 'Ingresar' })).toBeTruthy()
    expect(router.state.location.pathname).toBe('/ingresar')
    expect((router.state.location.state as { desde?: string }).desde).toBe(url)
  })

  it('sin sesión conserva también la búsqueda de la ruta pedida', async () => {
    borrarToken()
    const router = abrir('/projects/42/backlog?orden=prioridad')

    expect(await screen.findByRole('heading', { name: 'Ingresar' })).toBeTruthy()
    expect((router.state.location.state as { desde?: string }).desde).toBe(
      '/projects/42/backlog?orden=prioridad',
    )
  })

  it('/ingresar abre sin sesión', async () => {
    borrarToken()
    const router = abrir('/ingresar')

    expect(await screen.findByRole('heading', { name: 'Ingresar' })).toBeTruthy()
    expect(router.state.location.pathname).toBe('/ingresar')
  })

  it('/registro abre sin sesión', async () => {
    borrarToken()
    const router = abrir('/registro')

    expect(await screen.findByRole('heading', { name: 'Crear cuenta' })).toBeTruthy()
    expect(router.state.location.pathname).toBe('/registro')
  })

  it('tras ingresar vuelve a la ruta pedida', async () => {
    borrarToken()
    vi.stubGlobal(
      'fetch',
      vi.fn((url: string) =>
        Promise.resolve(
          new Response(
            url.includes('/auth/ingreso')
              ? JSON.stringify({ usuario: USUARIO, token: 'token-nuevo', vence_en: '2026-11-10T12:00:00Z' })
              : '[]',
            { status: 200 },
          ),
        ),
      ),
    )
    const router = abrir('/projects/42/backlog')
    await screen.findByRole('heading', { name: 'Ingresar' })

    await userEvent.type(screen.getByLabelText('Nombre de usuario'), 'ana.perez')
    await userEvent.type(screen.getByLabelText('Contraseña'), 'secreta123')
    await userEvent.click(screen.getByRole('button', { name: 'Ingresar' }))

    expect(await screen.findByRole('heading', { name: 'Product Backlog' })).toBeTruthy()
    expect(router.state.location.pathname).toBe('/projects/42/backlog')
  })

  it('el encabezado ofrece ingresar y registrarse a quien no tiene sesión', async () => {
    borrarToken()
    abrir('/ingresar')
    await screen.findByRole('heading', { name: 'Ingresar' })

    expect(screen.getAllByRole('link', { name: 'Ingresar' }).length).toBeGreaterThan(0)
    expect(screen.getAllByRole('link', { name: 'Registrarse' }).length).toBeGreaterThan(0)
  })
})
