import { render, screen } from '@testing-library/react'
import { useState } from 'react'
import userEvent from '@testing-library/user-event'
import { createMemoryRouter } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { LayoutProyecto } from './compartido/router/LayoutProyecto'
import { rutas } from './rutas'

const PROYECTO = {
  id: 42,
  nombre: 'Proyecto simulado',
  descripcion: '',
  fecha_inicio: null,
  fecha_fin: null,
  creado_en: '2026-10-01T00:00:00Z',
}

let proyectosSimulados: { id: number; nombre: string }[] = []

function responderListado(url: string) {
  if (url === '/projects') return JSON.stringify(proyectosSimulados)
  if (/^\/projects\/\d+$/.test(url)) return JSON.stringify(PROYECTO)
  return '[]'
}

beforeEach(() => {
  proyectosSimulados = []
  vi.stubGlobal(
    'fetch',
    vi.fn((url: string) => Promise.resolve(new Response(responderListado(url), { status: 200 }))),
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

  it('la pantalla de una feature sigue mostrándose en sus sub-rutas', async () => {
    abrir('/projects/1/proyecto/settings')

    expect(await screen.findByRole('heading', { name: 'Panel del proyecto' })).toBeTruthy()
  })

  it('la pantalla principal permite crear un proyecto', async () => {
    abrir('/')

    expect(await screen.findByRole('heading', { name: 'Crear proyecto' })).toBeTruthy()
  })

  describe('inicio con lista de proyectos', () => {
    it('US3.1: el inicio muestra el formulario de creación y la lista de proyectos', async () => {
      proyectosSimulados = [{ id: 1, nombre: 'Alfa' }, { id: 2, nombre: 'Beta' }]
      abrir('/')

      expect(await screen.findByRole('heading', { name: 'Crear proyecto' })).toBeTruthy()
      expect(await screen.findByRole('link', { name: 'Alfa (#1)' })).toBeTruthy()
      expect(screen.getByRole('link', { name: 'Beta (#2)' })).toBeTruthy()
    })

    it('US3.2: un enlace de la lista abre el panel del proyecto', async () => {
      proyectosSimulados = [{ id: 1, nombre: 'Alfa' }]
      const usuario = userEvent.setup()
      const router = abrir('/')

      await usuario.click(await screen.findByRole('link', { name: 'Alfa (#1)' }))

      expect(router.state.location.pathname).toBe('/projects/1/proyecto')
      expect(await screen.findByRole('heading', { name: 'Panel del proyecto' })).toBeTruthy()
    })
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

  describe('selector de proyecto en el encabezado', () => {
    const DOS_PROYECTOS = [
      { id: 1, nombre: 'Alfa' },
      { id: 2, nombre: 'Beta' },
    ]

    function llamadasAlListado() {
      return vi.mocked(fetch).mock.calls.filter(([url]) => url === '/projects').length
    }

    it('US1.1: elegir otro proyecto desde el backlog abre el backlog del proyecto elegido', async () => {
      proyectosSimulados = DOS_PROYECTOS
      const usuario = userEvent.setup()
      const router = abrir('/projects/1/backlog')
      await screen.findByRole('option', { name: 'Beta (#2)' })

      await usuario.selectOptions(screen.getByRole('combobox', { name: 'Proyecto' }), '2')

      expect(router.state.location.pathname).toBe('/projects/2/backlog')
      expect(await screen.findByRole('heading', { name: 'Product Backlog' })).toBeTruthy()
      expect(screen.getByRole<HTMLSelectElement>('combobox', { name: 'Proyecto' }).value).toBe('2')
    })

    it('US2.1: "Nuevo proyecto" abre la pantalla de creación desde una pantalla de proyecto', async () => {
      const usuario = userEvent.setup()
      const router = abrir('/projects/1/backlog')

      await usuario.click(await screen.findByRole('link', { name: 'Nuevo proyecto' }))

      expect(router.state.location.pathname).toBe('/')
      expect(await screen.findByRole('heading', { name: 'Crear proyecto' })).toBeTruthy()
    })

    it('GET /projects se pide una sola vez al montar y no al cambiar de sección del mismo proyecto', async () => {
      proyectosSimulados = DOS_PROYECTOS
      const usuario = userEvent.setup()
      abrir('/projects/1/backlog')
      await screen.findByRole('option', { name: 'Beta (#2)' })

      await usuario.click(screen.getByRole('link', { name: 'Sprints' }))
      await screen.findByRole('heading', { name: 'Sprints' })

      expect(llamadasAlListado()).toBe(1)
    })

    it('US1.4: al cambiar de proyecto la pantalla no muestra datos del proyecto anterior', async () => {
      proyectosSimulados = DOS_PROYECTOS
      const historia = {
        id: 7,
        proyecto_id: 1,
        titulo: 'Historia del proyecto uno',
        descripcion: '',
        criterios_aceptacion: null,
        prioridad: 'M',
        estado: 'Nueva',
        valor_negocio: null,
        estimacion_sp: null,
        sprint_id: null,
      }
      const respuestaInicial = vi.mocked(fetch).getMockImplementation()!
      vi.stubGlobal(
        'fetch',
        vi.fn((url: string) => {
          if (url === '/projects/1/backlog') {
            return Promise.resolve(new Response(JSON.stringify([historia]), { status: 200 }))
          }
          if (url === '/projects/2/backlog') return new Promise<Response>(() => {})
          return respuestaInicial(url, undefined)
        }),
      )
      const usuario = userEvent.setup()
      abrir('/projects/1/backlog')
      await screen.findByText('Historia del proyecto uno')

      await usuario.selectOptions(screen.getByRole('combobox', { name: 'Proyecto' }), '2')

      expect(screen.queryByText('Historia del proyecto uno')).toBeNull()
    })
  })

  it('US1.4: el layout remonta cualquier pantalla al cambiar el :id, aunque no tenga su propia key', async () => {
    function PantallaConEstado() {
      const [marca] = useState(() => Math.random().toString(36).slice(2))
      return <span data-testid="marca">{marca}</span>
    }
    const router = createMemoryRouter(
      [{ path: 'projects/:id', element: <LayoutProyecto />, children: [{ path: 'pantalla', element: <PantallaConEstado /> }] }],
      { initialEntries: ['/projects/1/pantalla'] },
    )
    render(<RouterProvider router={router} />)
    const marcaDelProyecto1 = (await screen.findByTestId('marca')).textContent

    await router.navigate('/projects/2/pantalla')

    await screen.findByText('/projects/2/pantalla', { selector: 'code' })
    expect(screen.getByTestId('marca').textContent).not.toBe(marcaDelProyecto1)
  })

  it.each(['/no-existe', '/projects/1/otra'])('%s muestra "Página no encontrada"', async (url) => {
    abrir(url)

    expect(await screen.findByRole('heading', { name: 'Página no encontrada' })).toBeTruthy()
  })
})
