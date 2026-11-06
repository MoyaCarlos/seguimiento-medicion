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

  it('la pantalla de una feature sigue mostrándose en sus sub-rutas', async () => {
    abrir('/projects/1/proyecto/settings')

    expect(await screen.findByRole('heading', { name: 'Panel del proyecto' })).toBeTruthy()
  })

  it('la pantalla principal permite crear un proyecto', async () => {
    abrir('/')

    expect(await screen.findByRole('heading', { name: 'Crear proyecto' })).toBeTruthy()
  })

  describe('inicio con datos de prueba', () => {
    it('el aviso de datos ya cargados queda junto a los botones, en el último bloque de la tarjeta del formulario', async () => {
      const respuestaBase = vi.mocked(fetch).getMockImplementation()!
      vi.stubGlobal(
        'fetch',
        vi.fn((url: string, init?: RequestInit) => {
          if (url === '/demo/estado') return Promise.resolve(new Response('{"habilitado":true}', { status: 200 }))
          if (url === '/demo/datos' && init?.method === 'POST') {
            return Promise.resolve(new Response(
              JSON.stringify({ campo: 'ya_cargado', mensaje: 'los datos de prueba ya están cargados' }),
              { status: 409 },
            ))
          }
          return Promise.resolve(respuestaBase(url, init))
        }),
      )
      const usuario = userEvent.setup()
      abrir('/')

      const botonDemo = await screen.findByRole('button', { name: 'Cargar datos de prueba' })
      await usuario.click(botonDemo)
      const aviso = await screen.findByText(/los datos de prueba ya están cargados/)

      const botonCrear = screen.getByRole('button', { name: 'Crear proyecto' })
      const tarjeta = botonCrear.closest('section')
      expect(tarjeta?.lastElementChild).toBe(botonDemo.parentElement)
      expect(botonDemo.parentElement?.contains(aviso)).toBe(true)
    })
  })

  describe('listado de proyectos caído', () => {
    function simularListadoCaido(extra?: (url: string, init?: RequestInit) => Response | undefined) {
      const respuestaBase = vi.mocked(fetch).getMockImplementation()!
      vi.stubGlobal(
        'fetch',
        vi.fn((url: string, init?: RequestInit) => {
          if (url === '/projects' && (!init?.method || init.method === 'GET')) {
            return Promise.resolve(new Response(JSON.stringify({ mensaje: 'Listado caído' }), { status: 500 }))
          }
          return Promise.resolve(extra?.(url, init) ?? respuestaBase(url, init))
        }),
      )
    }

    it('US4.2: el encabezado muestra el aviso y los enlaces de sección siguen navegando', async () => {
      simularListadoCaido()
      const usuario = userEvent.setup()
      const router = abrir('/projects/42/backlog')

      expect(await screen.findByText('Listado caído')).toBeTruthy()
      await usuario.click(screen.getByRole('link', { name: 'Sprints' }))

      expect(router.state.location.pathname).toBe('/projects/42/sprints')
      expect(screen.getByRole<HTMLSelectElement>('combobox', { name: 'Proyecto' }).disabled).toBe(true)
    })

    it('US3.4 y US4.4: en el inicio el aviso no impide crear un proyecto', async () => {
      simularListadoCaido((url, init) =>
        url === '/projects' && init?.method === 'POST'
          ? new Response(JSON.stringify({ ...PROYECTO, id: 77 }), { status: 201 })
          : undefined)
      const usuario = userEvent.setup()
      const router = abrir('/')

      expect(await screen.findByText('Listado caído')).toBeTruthy()
      await usuario.type(screen.getByLabelText('Nombre del proyecto'), 'Proyecto Alfa')
      await usuario.type(screen.getByLabelText('Tu nombre (Scrum Master)'), 'Ana')
      await usuario.click(screen.getByRole('button', { name: 'Crear proyecto' }))

      expect(await screen.findByRole('heading', { name: 'Panel del proyecto' })).toBeTruthy()
      expect(router.state.location.pathname).toBe('/projects/77/proyecto')
    })
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

  it('el menú de secciones de un proyecto se llama "Secciones del proyecto"', async () => {
    abrir('/projects/42/backlog')
    await screen.findByRole('heading', { name: 'Product Backlog' })

    expect(screen.getByRole('navigation', { name: 'Secciones del proyecto' })).toBeTruthy()
    expect(screen.queryByRole('navigation', { name: 'Historias del Sprint 2' })).toBeNull()
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
