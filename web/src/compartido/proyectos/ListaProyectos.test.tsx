import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { createMemoryRouter } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { ListaProyectos } from './ListaProyectos'

function abrir(proyectos: { id: number; nombre: string }[]) {
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
    abrir([{ id: 3, nombre: 'Gamma' }, { id: 1, nombre: 'Alfa' }])

    const enlaces = await screen.findAllByRole('link')

    expect(enlaces.map((enlace) => enlace.textContent)).toEqual(['Gamma (#3)', 'Alfa (#1)'])
    expect(enlaces.map((enlace) => enlace.getAttribute('href'))).toEqual([
      '/projects/3/proyecto',
      '/projects/1/proyecto',
    ])
  })

  it('US3.2: al activar el enlace de un proyecto se abre su panel', async () => {
    const usuario = userEvent.setup()
    const router = abrir([{ id: 1, nombre: 'Alfa' }, { id: 2, nombre: 'Beta' }])

    await usuario.click(await screen.findByRole('link', { name: 'Beta (#2)' }))

    expect(router.state.location.pathname).toBe('/projects/2/proyecto')
  })

  it('Límite: con nombres repetidos cada enlace abre el proyecto de su identificador', async () => {
    abrir([{ id: 1, nombre: 'Alfa' }, { id: 2, nombre: 'Alfa' }])

    const segundo = await screen.findByRole('link', { name: 'Alfa (#2)' })

    expect(segundo.getAttribute('href')).toBe('/projects/2/proyecto')
  })

  it('US3.3: sin proyectos invita a crear el primero y no muestra un error', async () => {
    abrir([])

    expect(await screen.findByText(/Todavía no hay proyectos/)).toBeTruthy()
    expect(screen.queryByRole('alert')).toBeNull()
    expect(screen.queryAllByRole('link')).toHaveLength(0)
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
