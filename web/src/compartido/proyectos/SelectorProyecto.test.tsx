import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { createMemoryRouter } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { SelectorProyecto } from './SelectorProyecto'

const PROYECTOS = [
  { id: 1, nombre: 'Alfa' },
  { id: 2, nombre: 'Beta' },
  { id: 3, nombre: 'Alfa' },
]

function abrir(url: string, proyectos = PROYECTOS) {
  vi.stubGlobal('fetch', vi.fn(() => Promise.resolve(new Response(JSON.stringify(proyectos), { status: 200 }))))
  const router = createMemoryRouter(
    [{ path: '/projects/:id/*', element: <SelectorProyecto /> }],
    { initialEntries: [url] },
  )
  render(<RouterProvider router={router} />)
  return router
}

afterEach(() => {
  vi.unstubAllGlobals()
})

const botonDelSelector = () => screen.getByRole('button', { name: /^Proyecto actual/ })
const lista = () => screen.getByRole('listbox', { name: 'Proyecto actual' })
const opciones = () => within(lista()).getAllByRole('option')
const nombresDeOpciones = () => opciones().map((opcion) => opcion.textContent)

async function abrirLista(usuario: ReturnType<typeof userEvent.setup>) {
  await usuario.click(botonDelSelector())
}

describe('SelectorProyecto', () => {
  it('US1.3: lista los proyectos como "<nombre> (#<id>)" en el orden de la API', async () => {
    const usuario = userEvent.setup()
    abrir('/projects/1/backlog')
    await screen.findByText('Alfa')

    await abrirLista(usuario)

    expect(nombresDeOpciones()).toEqual(['Alfa (#1)', 'Beta (#2)', 'Alfa (#3)'])
  })

  it('US1.3: marca como actual el proyecto del :id de la URL', async () => {
    const usuario = userEvent.setup()
    abrir('/projects/2/backlog')
    await screen.findByText('Beta')

    await abrirLista(usuario)

    expect(screen.getByRole('option', { name: 'Beta (#2)' }).getAttribute('aria-selected')).toBe('true')
    expect(within(botonDelSelector()).getByText('#2')).toBeTruthy()
  })

  it('Límite: con nombres repetidos marca el actual por su identificador y no por el nombre', async () => {
    const usuario = userEvent.setup()
    abrir('/projects/3/backlog')
    await screen.findByText('Alfa')

    await abrirLista(usuario)

    expect(screen.getByRole('option', { name: 'Alfa (#3)' }).getAttribute('aria-selected')).toBe('true')
    expect(screen.getByRole('option', { name: 'Alfa (#1)' }).getAttribute('aria-selected')).toBe('false')
  })

  it('Límite: si el proyecto actual no está en la lista muestra "Proyecto <id>" y no marca otro', async () => {
    const usuario = userEvent.setup()
    abrir('/projects/9/backlog')
    await screen.findByText('Proyecto 9', { selector: 'strong' })

    await abrirLista(usuario)

    expect(screen.getByRole('option', { name: 'Proyecto 9' }).getAttribute('aria-selected')).toBe('true')
    expect(screen.getByRole('option', { name: 'Beta (#2)' }).getAttribute('aria-selected')).toBe('false')
    expect(screen.getByRole('option', { name: 'Alfa (#1)' }).getAttribute('aria-selected')).toBe('false')
  })

  it('Límite: con la lista vacía el selector muestra solo "Proyecto <id>"', async () => {
    const usuario = userEvent.setup()
    abrir('/projects/5/backlog', [])
    await screen.findByText('Proyecto 5', { selector: 'strong' })

    await abrirLista(usuario)

    expect(nombresDeOpciones()).toEqual(['Proyecto 5'])
  })

  it('Accesibilidad: el botón se llama "Proyecto actual", se enfoca con Tab y se opera con el teclado (abrir, flechas, Enter)', async () => {
    const usuario = userEvent.setup()
    const router = abrir('/projects/1/backlog')
    await screen.findByText('Alfa')

    await usuario.tab()
    expect(document.activeElement).toBe(botonDelSelector())
    expect(botonDelSelector().getAttribute('aria-haspopup')).toBe('listbox')
    expect(botonDelSelector().getAttribute('aria-expanded')).toBe('false')

    await usuario.keyboard('{ArrowDown}')
    expect(botonDelSelector().getAttribute('aria-expanded')).toBe('true')
    await usuario.keyboard('{ArrowDown}{Enter}')

    expect(router.state.location.pathname).toBe('/projects/2/backlog')
  })

  it('US1.1: dado el proyecto 1 como actual, al elegir el 2 navega a la misma sección del 2', async () => {
    const usuario = userEvent.setup()
    const router = abrir('/projects/1/backlog')
    await screen.findByText('Alfa')
    await abrirLista(usuario)

    await usuario.click(screen.getByRole('option', { name: 'Beta (#2)' }))

    expect(router.state.location.pathname).toBe('/projects/2/backlog')
    expect(await screen.findByText('Beta')).toBeTruthy()
    expect(botonDelSelector().getAttribute('aria-expanded')).toBe('false')
  })

  it('US1.2: desde /proyecto/edit al elegir otro proyecto solo se conserva la sección', async () => {
    const usuario = userEvent.setup()
    const router = abrir('/projects/1/proyecto/edit')
    await screen.findByText('Alfa')
    await abrirLista(usuario)

    await usuario.click(screen.getByRole('option', { name: 'Beta (#2)' }))

    expect(router.state.location.pathname).toBe('/projects/2/proyecto')
  })

  it('US1.5: desde una sección que no es del menú al elegir otro proyecto cae en proyecto', async () => {
    const usuario = userEvent.setup()
    const router = abrir('/projects/1/otra')
    await screen.findByText('Alfa')
    await abrirLista(usuario)

    await usuario.click(screen.getByRole('option', { name: 'Beta (#2)' }))

    expect(router.state.location.pathname).toBe('/projects/2/proyecto')
  })

  it('Límite: con nombres repetidos abre el proyecto elegido por su identificador', async () => {
    const usuario = userEvent.setup()
    const router = abrir('/projects/1/sprints')
    await screen.findByText('Alfa')
    await abrirLista(usuario)

    await usuario.click(screen.getByRole('option', { name: 'Alfa (#3)' }))

    expect(router.state.location.pathname).toBe('/projects/3/sprints')
  })

  it('Límite: elegir el proyecto que ya está abierto no navega y cierra la lista', async () => {
    const usuario = userEvent.setup()
    const router = abrir('/projects/1/otra')
    await screen.findByText('Alfa')
    const historialAntes = router.state.historyAction
    await abrirLista(usuario)

    await usuario.click(screen.getByRole('option', { name: 'Alfa (#1)' }))

    expect(router.state.location.pathname).toBe('/projects/1/otra')
    expect(router.state.historyAction).toBe(historialAntes)
    expect(botonDelSelector().getAttribute('aria-expanded')).toBe('false')
  })

  it('US2.1: «Nuevo proyecto» es un botón aparte, fuera de la lista de opciones, que abre la ventana de creación', async () => {
    const usuario = userEvent.setup()
    const router = abrir('/projects/1/backlog')
    await screen.findByText('Alfa')
    await abrirLista(usuario)

    const nuevo = screen.getByRole('button', { name: 'Nuevo proyecto' })

    expect(lista().contains(nuevo)).toBe(false)
    expect(screen.queryByRole('option', { name: 'Nuevo proyecto' })).toBeNull()
    await usuario.click(nuevo)
    expect(router.state.location.pathname).toBe('/projects/new')
  })

  it('US2.2: «Nuevo proyecto» está disponible mientras el listado carga', async () => {
    const usuario = userEvent.setup()
    vi.stubGlobal('fetch', vi.fn(() => new Promise<Response>(() => {})))
    render(<RouterProvider router={createMemoryRouter(
      [{ path: '/projects/:id/*', element: <SelectorProyecto /> }],
      { initialEntries: ['/projects/1/backlog'] },
    )} />)

    await abrirLista(usuario)

    expect(screen.getByRole('button', { name: 'Nuevo proyecto' })).toBeTruthy()
  })

  it('US2.2: «Nuevo proyecto» está disponible cuando el listado falla', async () => {
    const usuario = userEvent.setup()
    vi.stubGlobal('fetch', vi.fn(() => Promise.reject(new TypeError('red caída'))))
    render(<RouterProvider router={createMemoryRouter(
      [{ path: '/projects/:id/*', element: <SelectorProyecto /> }],
      { initialEntries: ['/projects/1/backlog'] },
    )} />)
    await screen.findByRole('alert')

    await abrirLista(usuario)

    expect(screen.getByRole('button', { name: 'Nuevo proyecto' })).toBeTruthy()
  })

  describe('teclado y cierre', () => {
    it('Escape cierra la lista y devuelve el foco al botón', async () => {
      const usuario = userEvent.setup()
      abrir('/projects/1/backlog')
      await screen.findByText('Alfa')
      await abrirLista(usuario)

      await usuario.keyboard('{Escape}')

      expect(botonDelSelector().getAttribute('aria-expanded')).toBe('false')
      expect(document.activeElement).toBe(botonDelSelector())
    })

    it('Inicio y Fin saltan a la primera y a la última opción; Enter elige la resaltada', async () => {
      const usuario = userEvent.setup()
      const router = abrir('/projects/2/backlog')
      await screen.findByText('Beta')
      await usuario.tab()
      await usuario.keyboard('{Enter}')

      await usuario.keyboard('{End}{Enter}')

      expect(router.state.location.pathname).toBe('/projects/3/backlog')
    })

    it('al abrir, la opción resaltada es la del proyecto actual y las flechas no pasan de los extremos', async () => {
      const usuario = userEvent.setup()
      const router = abrir('/projects/3/backlog')
      await screen.findByText('Alfa')
      await usuario.tab()
      await usuario.keyboard('{Enter}')

      await usuario.keyboard('{ArrowDown}{ArrowDown}{Enter}')

      expect(router.state.location.pathname).toBe('/projects/3/backlog')
      await usuario.keyboard('{Enter}{Home}{ArrowUp}{Enter}')
      expect(router.state.location.pathname).toBe('/projects/1/backlog')
    })

    it('un clic fuera de la lista la cierra', async () => {
      const usuario = userEvent.setup()
      abrir('/projects/1/backlog')
      await screen.findByText('Alfa')
      await abrirLista(usuario)

      await usuario.click(document.body)

      expect(botonDelSelector().getAttribute('aria-expanded')).toBe('false')
    })

    it('el botón abre y cierra la lista al hacer clic', async () => {
      const usuario = userEvent.setup()
      abrir('/projects/1/backlog')
      await screen.findByText('Alfa')

      await abrirLista(usuario)
      expect(botonDelSelector().getAttribute('aria-expanded')).toBe('true')
      await abrirLista(usuario)

      expect(botonDelSelector().getAttribute('aria-expanded')).toBe('false')
    })
  })

  describe('estados del listado', () => {
    function abrirConFetch(implementacion: () => Promise<Response>, url = '/projects/1/backlog') {
      vi.stubGlobal('fetch', vi.fn(implementacion))
      const router = createMemoryRouter(
        [{ path: '/projects/:id/*', element: <SelectorProyecto /> }],
        { initialEntries: [url] },
      )
      render(<RouterProvider router={router} />)
    }

    it('US4.1: mientras carga muestra el indicador y la lista ofrece solo "Proyecto <id>"', async () => {
      const usuario = userEvent.setup()
      abrirConFetch(() => new Promise<Response>(() => {}))

      expect(screen.getByRole('status').textContent).toBe('Cargando…')
      await abrirLista(usuario)
      expect(nombresDeOpciones()).toEqual(['Proyecto 1'])
      expect(screen.getByRole('button', { name: 'Nuevo proyecto' })).toBeTruthy()
    })

    it('US4.2: si el listado falla muestra el mensaje de la API y la lista ofrece solo "Proyecto <id>"', async () => {
      const usuario = userEvent.setup()
      abrirConFetch(() => Promise.resolve(new Response(JSON.stringify({ mensaje: 'Falló el listado' }), { status: 500 })))

      expect((await screen.findByRole('alert')).textContent).toBe('Falló el listado')
      await abrirLista(usuario)
      expect(nombresDeOpciones()).toEqual(['Proyecto 1'])
      expect(screen.getByRole('button', { name: 'Nuevo proyecto' })).toBeTruthy()
      expect(screen.queryByRole('status')).toBeNull()
    })

    it('Cuando la carga se resuelve la lista ofrece todos los proyectos y desaparece el indicador', async () => {
      const usuario = userEvent.setup()
      abrirConFetch(() => Promise.resolve(new Response(JSON.stringify(PROYECTOS), { status: 200 })))
      await screen.findByText('Alfa')

      await abrirLista(usuario)

      expect(nombresDeOpciones()).toEqual(['Alfa (#1)', 'Beta (#2)', 'Alfa (#3)'])
      expect(screen.queryByRole('status')).toBeNull()
      expect(screen.queryByRole('alert')).toBeNull()
    })
  })
})
