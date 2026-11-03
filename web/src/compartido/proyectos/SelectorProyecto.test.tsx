import { render, screen } from '@testing-library/react'
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

describe('SelectorProyecto', () => {
  it('US1.3: lista los proyectos como "<nombre> (#<id>)" en el orden de la API', async () => {
    abrir('/projects/1/backlog')

    await screen.findByRole('option', { name: 'Beta (#2)' })
    expect(screen.getAllByRole('option').map((opcion) => opcion.textContent)).toEqual([
      'Alfa (#1)',
      'Beta (#2)',
      'Alfa (#3)',
    ])
  })

  it('US1.3: marca como actual el proyecto del :id de la URL', async () => {
    abrir('/projects/2/backlog')

    await screen.findByRole('option', { name: 'Beta (#2)' })
    expect(screen.getByRole<HTMLSelectElement>('combobox', { name: 'Proyecto' }).value).toBe('2')
  })

  it('Límite: con nombres repetidos marca el actual por su identificador y no por el nombre', async () => {
    abrir('/projects/3/backlog')

    await screen.findByRole('option', { name: 'Alfa (#3)' })
    expect(screen.getByRole<HTMLOptionElement>('option', { name: 'Alfa (#3)' }).selected).toBe(true)
    expect(screen.getByRole<HTMLOptionElement>('option', { name: 'Alfa (#1)' }).selected).toBe(false)
  })

  it('Límite: si el proyecto actual no está en la lista muestra "Proyecto <id>" y no marca otro', async () => {
    abrir('/projects/9/backlog')

    await screen.findByRole('option', { name: 'Beta (#2)' })
    const selector = screen.getByRole<HTMLSelectElement>('combobox', { name: 'Proyecto' })
    expect(selector.selectedOptions[0].textContent).toBe('Proyecto 9')
    expect(screen.getByRole<HTMLOptionElement>('option', { name: 'Alfa (#1)' }).selected).toBe(false)
  })

  it('Límite: con la lista vacía el selector muestra solo "Proyecto <id>"', async () => {
    abrir('/projects/5/backlog', [])

    expect(await screen.findByRole('option', { name: 'Proyecto 5' })).toBeTruthy()
    expect(screen.getAllByRole('option')).toHaveLength(1)
  })

  it('Accesibilidad: el selector se llama "Proyecto" y se enfoca y opera con el teclado', async () => {
    const usuario = userEvent.setup()
    const router = abrir('/projects/1/backlog')
    await screen.findByRole('option', { name: 'Beta (#2)' })

    await usuario.tab()
    const selector = screen.getByRole('combobox', { name: 'Proyecto' })
    expect(document.activeElement).toBe(selector)

    await usuario.selectOptions(selector, '2')
    expect(router.state.location.pathname).toBe('/projects/2/backlog')
  })

  it('US1.1: dado el proyecto 1 como actual, al elegir el 2 navega a la misma sección del 2', async () => {
    const usuario = userEvent.setup()
    const router = abrir('/projects/1/backlog')
    await screen.findByRole('option', { name: 'Beta (#2)' })

    await usuario.selectOptions(screen.getByRole('combobox', { name: 'Proyecto' }), 'Beta (#2)')

    expect(router.state.location.pathname).toBe('/projects/2/backlog')
    expect(screen.getByRole<HTMLSelectElement>('combobox', { name: 'Proyecto' }).value).toBe('2')
  })

  it('US1.2: desde /proyecto/edit al elegir otro proyecto solo se conserva la sección', async () => {
    const usuario = userEvent.setup()
    const router = abrir('/projects/1/proyecto/edit')
    await screen.findByRole('option', { name: 'Beta (#2)' })

    await usuario.selectOptions(screen.getByRole('combobox', { name: 'Proyecto' }), '2')

    expect(router.state.location.pathname).toBe('/projects/2/proyecto')
  })

  it('US1.5: desde una sección que no es del menú al elegir otro proyecto cae en proyecto', async () => {
    const usuario = userEvent.setup()
    const router = abrir('/projects/1/otra')
    await screen.findByRole('option', { name: 'Beta (#2)' })

    await usuario.selectOptions(screen.getByRole('combobox', { name: 'Proyecto' }), '2')

    expect(router.state.location.pathname).toBe('/projects/2/proyecto')
  })

  it('Límite: con nombres repetidos abre el proyecto elegido por su identificador', async () => {
    const usuario = userEvent.setup()
    const router = abrir('/projects/1/sprints')
    await screen.findByRole('option', { name: 'Alfa (#3)' })

    await usuario.selectOptions(screen.getByRole('combobox', { name: 'Proyecto' }), 'Alfa (#3)')

    expect(router.state.location.pathname).toBe('/projects/3/sprints')
  })

  it('Límite: elegir el proyecto que ya está abierto no navega', async () => {
    const usuario = userEvent.setup()
    const router = abrir('/projects/1/otra')
    await screen.findByRole('option', { name: 'Beta (#2)' })
    const historialAntes = router.state.historyAction

    await usuario.selectOptions(screen.getByRole('combobox', { name: 'Proyecto' }), 'Alfa (#1)')

    expect(router.state.location.pathname).toBe('/projects/1/otra')
    expect(router.state.historyAction).toBe(historialAntes)
  })
})
