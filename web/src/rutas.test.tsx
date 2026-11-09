import { render, screen } from '@testing-library/react'
import { useState } from 'react'
import userEvent from '@testing-library/user-event'
import { createMemoryRouter } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { LayoutProyecto } from './compartido/router/LayoutProyecto'
import { borrarToken, guardarToken } from './features/autenticacion/almacenamientoToken'
import { rutas } from './rutas'

const USUARIO = { id: 7, nombre: 'ana.perez' }

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
  if (url.includes('/auth/yo')) return JSON.stringify(USUARIO)
  if (url === '/projects') return JSON.stringify(proyectosSimulados)
  if (/^\/projects\/\d+$/.test(url)) return JSON.stringify(PROYECTO)
  return '[]'
}

// Sesión simulada: token guardado y `GET /auth/yo` responde 200.
beforeEach(() => {
  borrarToken()
  localStorage.clear()
  guardarToken('token-de-prueba')
  proyectosSimulados = []
  vi.stubGlobal(
    'fetch',
    vi.fn((url: string) => Promise.resolve(new Response(responderListado(url), { status: 200 }))),
  )
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
    vi.stubGlobal(
      'fetch',
      vi.fn((url: string) =>
        Promise.resolve(
          new Response(
            url.includes('/auth/yo')
              ? JSON.stringify(USUARIO)
              : url.includes('/votacion')
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
      vi.fn((url: string) =>
        Promise.resolve(
          new Response(
            JSON.stringify(
              url.includes('/auth/yo')
                ? USUARIO
                : { proyecto_id: 42, horas_por_story_point: null, proyecto: metricas, sprints: [] },
            ),
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

    describe('renombrar el proyecto abierto', () => {
      async function abrirEdicionDeBeta(usuario: ReturnType<typeof userEvent.setup>) {
        proyectosSimulados = [{ id: 27, nombre: 'Beta' }]
        const router = abrir('/projects/27/proyecto/edit')
        await screen.findByRole('option', { name: 'Beta (#27)' })
        const nombre = await screen.findByLabelText('Nombre del proyecto')
        await usuario.clear(nombre)
        await usuario.type(nombre, 'Beta Renombrado')
        return router
      }

      it('al guardar el nombre nuevo, el selector lo muestra sin recargar la página', async () => {
        const usuario = userEvent.setup()
        const router = await abrirEdicionDeBeta(usuario)
        proyectosSimulados = [{ id: 27, nombre: 'Beta Renombrado' }]

        await usuario.click(screen.getByRole('button', { name: 'Guardar cambios' }))

        expect(await screen.findByRole('option', { name: 'Beta Renombrado (#27)' })).toBeTruthy()
        expect(screen.queryByRole('option', { name: 'Beta (#27)' })).toBeNull()
        expect(screen.getByRole<HTMLSelectElement>('combobox', { name: 'Proyecto' }).value).toBe('27')
        expect(router.state.location.pathname).toBe('/projects/27/proyecto')
      })

      it('al cancelar la edición no se vuelve a pedir el listado', async () => {
        const usuario = userEvent.setup()
        const router = await abrirEdicionDeBeta(usuario)

        await usuario.click(screen.getByRole('link', { name: 'Cancelar' }))

        expect(router.state.location.pathname).toBe('/projects/27/proyecto')
        expect(await screen.findByRole('option', { name: 'Beta (#27)' })).toBeTruthy()
        expect(llamadasAlListado()).toBe(1)
      })

      it('si el guardado falla no se vuelve a pedir el listado', async () => {
        const usuario = userEvent.setup()
        await abrirEdicionDeBeta(usuario)
        const respuestaBase = vi.mocked(fetch).getMockImplementation()!
        const llamadasPrevias = llamadasAlListado()
        vi.stubGlobal(
          'fetch',
          vi.fn((url: string, init?: RequestInit) =>
            init?.method === 'PUT'
              ? Promise.resolve(new Response(JSON.stringify({ mensaje: 'No se pudo guardar' }), { status: 500 }))
              : respuestaBase(url, init)),
        )

        await usuario.click(screen.getByRole('button', { name: 'Guardar cambios' }))

        expect(await screen.findByText('No se pudo guardar')).toBeTruthy()
        expect(screen.getByRole('option', { name: 'Beta (#27)' })).toBeTruthy()
        expect(llamadasPrevias).toBe(1)
        expect(llamadasAlListado()).toBe(0)
      })
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
    '/projects/42/dashboard',
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
