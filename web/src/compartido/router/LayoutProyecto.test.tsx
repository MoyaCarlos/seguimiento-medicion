import { render, screen, waitFor } from '@testing-library/react'
import { createMemoryRouter } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { LayoutProyecto } from './LayoutProyecto'

function simular(status: number, cuerpo: unknown) {
  vi.stubGlobal('fetch', vi.fn(() => Promise.resolve(new Response(JSON.stringify(cuerpo), { status }))))
}

function abrir() {
  const router = createMemoryRouter(
    [
      {
        path: 'projects/:id',
        element: <LayoutProyecto />,
        children: [{ path: 'backlog', element: <p>Pantalla del proyecto</p> }],
      },
      { path: '/', element: <p>Mis proyectos (inicio)</p> },
    ],
    { initialEntries: ['/projects/5/backlog'] },
  )
  render(<RouterProvider router={router} />)
}

afterEach(() => vi.unstubAllGlobals())

describe('LayoutProyecto', () => {
  it('Escenario: quien no integra el proyecto ve «No sos integrante de este proyecto» en lugar de la pantalla', async () => {
    simular(403, { mensaje: 'No sos integrante de este proyecto' })

    abrir()

    expect(await screen.findByRole('heading', { name: 'No sos integrante de este proyecto' })).toBeTruthy()
    expect(screen.getByRole('link', { name: 'Mis proyectos' }).getAttribute('href')).toBe('/')
    expect(screen.queryByText('Pantalla del proyecto')).toBeNull()
  })

  it('Escenario: un integrante ve la pantalla del proyecto', async () => {
    simular(200, { id: 5, rol: 'product_builder', fecha_cierre: null })

    abrir()

    expect(await screen.findByText('Pantalla del proyecto')).toBeTruthy()
    await waitFor(() => expect(screen.queryByRole('heading', { name: 'No sos integrante de este proyecto' })).toBeNull())
  })

  it('Escenario: otro 403 no oculta la pantalla, decide el servidor en cada lectura', async () => {
    simular(403, { mensaje: 'Acceso denegado' })

    abrir()

    expect(await screen.findByText('Pantalla del proyecto')).toBeTruthy()
    expect(screen.queryByRole('heading', { name: 'No sos integrante de este proyecto' })).toBeNull()
  })
})
