import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { createMemoryRouter } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import { afterEach, describe, expect, it, vi } from 'vitest'
import type { ProyectoResumen } from './api'
import { ListaProyectos } from './ListaProyectos'

function p(id: number, nombre: string, resto: Partial<ProyectoResumen> = {}): ProyectoResumen {
  return { id, nombre, descripcion: '', fecha_inicio: null, fecha_fin: null, ...resto }
}

function abrir(proyectos: ProyectoResumen[]) {
  vi.stubGlobal('fetch', vi.fn(() => Promise.resolve(new Response(JSON.stringify(proyectos), { status: 200 }))))
  const router = createMemoryRouter(
    [
      { path: '/', element: <ListaProyectos /> },
      { path: '/projects/:id/proyecto', element: <p>Panel abierto</p> },
    ],
    { initialEntries: ['/'] },
  )
  render(<RouterProvider router={router} />)
  return router
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('ListaProyectos', () => {
  it('US3.1: lista cada proyecto como "<nombre> (#<id>)" con un enlace a su panel, en el orden de la API', async () => {
    abrir([p(3, 'Gamma'), p(1, 'Alfa')])

    const enlaces = await screen.findAllByRole('link')

    expect(enlaces.map((enlace) => enlace.getAttribute('aria-label'))).toEqual(['Gamma (#3)', 'Alfa (#1)'])
    expect(enlaces.map((enlace) => enlace.getAttribute('href'))).toEqual([
      '/projects/3/proyecto',
      '/projects/1/proyecto',
    ])
  })

  it('US3.2: al activar el enlace de un proyecto se abre su panel', async () => {
    const usuario = userEvent.setup()
    const router = abrir([p(1, 'Alfa'), p(2, 'Beta')])

    await usuario.click(await screen.findByRole('link', { name: 'Beta (#2)' }))

    expect(router.state.location.pathname).toBe('/projects/2/proyecto')
  })

  it('Límite: con nombres repetidos cada enlace abre el proyecto de su identificador', async () => {
    abrir([p(1, 'Alfa'), p(2, 'Alfa')])

    const segundo = await screen.findByRole('link', { name: 'Alfa (#2)' })

    expect(segundo.getAttribute('href')).toBe('/projects/2/proyecto')
  })

  it('US3.3: sin proyectos invita a crear el primero y no muestra un error', async () => {
    abrir([])

    expect(await screen.findByText(/Todavía no hay proyectos/)).toBeTruthy()
    expect(screen.queryByRole('alert')).toBeNull()
    expect(screen.queryAllByRole('link')).toHaveLength(0)
  })

  describe('tarjetas y contador', () => {
    const COMPLETO = p(1, 'Licencia de conducir', {
      descripcion: 'Trámite de licencia',
      fecha_inicio: '2026-10-19',
      fecha_fin: '2026-12-18',
    })

    it('US3.5: la tarjeta muestra nombre, #id, descripción y fechas y es un único enlace al panel', async () => {
      abrir([COMPLETO])

      const tarjeta = await screen.findByRole('link', { name: 'Licencia de conducir (#1)' })

      expect(tarjeta.getAttribute('href')).toBe('/projects/1/proyecto')
      expect(within(tarjeta).getByText('Licencia de conducir')).toBeTruthy()
      expect(within(tarjeta).getByText('#1')).toBeTruthy()
      expect(within(tarjeta).getByText('Trámite de licencia')).toBeTruthy()
      expect(within(tarjeta).getByText('Inicio: 2026-10-19')).toBeTruthy()
      expect(within(tarjeta).getByText('Fin: 2026-12-18')).toBeTruthy()
      expect(screen.getAllByRole('link')).toHaveLength(1)
    })

    it('US3.6: sin descripción ni fechas la tarjeta muestra "Sin descripción" y "Sin definir"', async () => {
      abrir([p(2, 'Beta')])

      const tarjeta = await screen.findByRole('link', { name: 'Beta (#2)' })

      expect(within(tarjeta).getByText('Sin descripción')).toBeTruthy()
      expect(within(tarjeta).getByText('Inicio: Sin definir')).toBeTruthy()
      expect(within(tarjeta).getByText('Fin: Sin definir')).toBeTruthy()
    })

    it('US3.7: el nombre completo queda en el atributo title de la tarjeta', async () => {
      const nombreLargo = 'Sistema de turnos del hospital con un nombre bastante largo'
      abrir([p(3, nombreLargo)])

      const tarjeta = await screen.findByRole('link', { name: `${nombreLargo} (#3)` })

      expect(within(tarjeta).getByText(nombreLargo).getAttribute('title')).toBe(nombreLargo)
    })

    it('US3.8: el encabezado "Proyectos existentes" muestra la cantidad de proyectos', async () => {
      abrir([p(1, 'Alfa'), p(2, 'Beta'), p(3, 'Gamma')])

      expect(await screen.findByRole('heading', { name: 'Proyectos existentes 3' })).toBeTruthy()
    })

    it('Límite: con la lista vacía el contador muestra 0', async () => {
      abrir([])

      expect(await screen.findByRole('heading', { name: 'Proyectos existentes 0' })).toBeTruthy()
    })

    it('La tarjeta no muestra estado, Scrum Master ni barra de avance', async () => {
      abrir([COMPLETO])

      const tarjeta = await screen.findByRole('link', { name: 'Licencia de conducir (#1)' })

      expect(within(tarjeta).queryByText(/Por comenzar|En curso|Finalizado|Scrum Master/)).toBeNull()
      expect(within(tarjeta).queryByRole('progressbar')).toBeNull()
    })
  })

  describe('buscador', () => {
    const PROYECTOS = [
      p(1, 'Gestión de Turnos', { descripcion: 'Agenda médica' }),
      p(2, 'Licencia de conducir', { descripcion: 'Incluye gestión de turnos' }),
      p(3, 'Tienda en línea'),
      p(4, 'Gestión de Stock'),
    ]

    const nombresVisibles = () => screen.queryAllByRole('link').map((enlace) => enlace.getAttribute('aria-label'))

    async function abrirYBuscar(texto: string) {
      const usuario = userEvent.setup()
      abrir(PROYECTOS)
      const campo = await screen.findByRole('searchbox', { name: 'Buscar proyecto' })
      if (texto) await usuario.type(campo, texto)
      return { usuario, campo }
    }

    it('US5.1: al escribir muestra solo los proyectos cuyo nombre coincide y el contador "N de M"', async () => {
      await abrirYBuscar('tienda')

      expect(nombresVisibles()).toEqual(['Tienda en línea (#3)'])
      expect(screen.getByRole('heading', { name: 'Proyectos existentes 1 de 4' })).toBeTruthy()
    })

    it('US5.2: ignora mayúsculas y acentos', async () => {
      await abrirYBuscar('GESTION')

      expect(nombresVisibles()).toEqual(['Gestión de Turnos (#1)', 'Gestión de Stock (#4)'])
    })

    it('US5.2: buscar con acento encuentra un nombre sin acento', async () => {
      const usuario = userEvent.setup()
      abrir([p(7, 'Cafeteria del club')])
      await usuario.type(await screen.findByRole('searchbox', { name: 'Buscar proyecto' }), 'cafetería')

      expect(nombresVisibles()).toEqual(['Cafeteria del club (#7)'])
    })

    it('US5.3: recorta los espacios al inicio y al final', async () => {
      await abrirYBuscar('  tienda  ')

      expect(nombresVisibles()).toEqual(['Tienda en línea (#3)'])
    })

    it('US5.4: borrar la búsqueda o dejar solo espacios vuelve a mostrar todos con el contador "M"', async () => {
      const { usuario, campo } = await abrirYBuscar('tienda')

      await usuario.clear(campo)
      expect(nombresVisibles()).toHaveLength(4)
      expect(screen.getByRole('heading', { name: 'Proyectos existentes 4' })).toBeTruthy()

      await usuario.type(campo, '   ')
      expect(nombresVisibles()).toHaveLength(4)
      expect(screen.getByRole('heading', { name: 'Proyectos existentes 4' })).toBeTruthy()
    })

    it('US5.5: conserva el orden de la API entre los resultados', async () => {
      await abrirYBuscar('gestion')

      expect(nombresVisibles()).toEqual(['Gestión de Turnos (#1)', 'Gestión de Stock (#4)'])
    })

    it('US5.6: no filtra por descripción ni por identificador', async () => {
      const { usuario, campo } = await abrirYBuscar('agenda')
      expect(nombresVisibles()).toEqual([])

      await usuario.clear(campo)
      await usuario.type(campo, '3')
      expect(nombresVisibles()).toEqual([])
    })

    it('US5.7: sin coincidencias muestra el mensaje, "0 de M" y conserva el buscador', async () => {
      await abrirYBuscar('zzz')

      expect(screen.getByText('Ningún proyecto coincide con la búsqueda')).toBeTruthy()
      expect(screen.queryByText(/Todavía no hay proyectos/)).toBeNull()
      expect(screen.getByRole('heading', { name: 'Proyectos existentes 0 de 4' })).toBeTruthy()
      expect(screen.getByRole('searchbox', { name: 'Buscar proyecto' })).toBeTruthy()
    })

    it('US5.8: con la lista vacía no hay buscador y se invita a crear el primero', async () => {
      abrir([])

      expect(await screen.findByText(/Todavía no hay proyectos/)).toBeTruthy()
      expect(screen.queryByRole('searchbox')).toBeNull()
      expect(screen.queryByText('Ningún proyecto coincide con la búsqueda')).toBeNull()
    })

    it('US5.9: mientras carga y si el listado falla no hay buscador', async () => {
      vi.stubGlobal('fetch', vi.fn(() => new Promise<Response>(() => {})))
      const { unmount } = render(<RouterProvider router={createMemoryRouter(
        [{ path: '/', element: <ListaProyectos /> }],
        { initialEntries: ['/'] },
      )} />)
      expect(screen.queryByRole('searchbox')).toBeNull()
      unmount()

      vi.stubGlobal('fetch', vi.fn(() => Promise.resolve(new Response(JSON.stringify({ mensaje: 'Caído' }), { status: 500 }))))
      render(<RouterProvider router={createMemoryRouter(
        [{ path: '/', element: <ListaProyectos /> }],
        { initialEntries: ['/'] },
      )} />)
      await screen.findByRole('alert')
      expect(screen.queryByRole('searchbox')).toBeNull()
    })

    it('US5.10: filtrar y borrar la búsqueda no vuelve a llamar a GET /projects', async () => {
      const { usuario, campo } = await abrirYBuscar('tienda')
      await usuario.clear(campo)

      expect(vi.mocked(fetch)).toHaveBeenCalledTimes(1)
    })
  })

  describe('estados del listado', () => {
    function abrirConFetch(implementacion: () => Promise<Response>) {
      vi.stubGlobal('fetch', vi.fn(implementacion))
      render(<RouterProvider router={createMemoryRouter(
        [{ path: '/', element: <ListaProyectos /> }],
        { initialEntries: ['/'] },
      )} />)
    }

    it('Mientras carga muestra el indicador de carga', () => {
      abrirConFetch(() => new Promise<Response>(() => {}))

      expect(screen.getByRole('status').textContent).toBe('Cargando…')
    })

    it('Si el listado falla muestra el mensaje de la API y ninguna lista', async () => {
      abrirConFetch(() => Promise.resolve(new Response(JSON.stringify({ mensaje: 'Falló el listado' }), { status: 500 })))

      expect((await screen.findByRole('alert')).textContent).toBe('Falló el listado')
      expect(screen.queryByRole('list')).toBeNull()
      expect(screen.queryByText(/Todavía no hay proyectos/)).toBeNull()
    })
  })
})
