import { render, screen } from '@testing-library/react'
import { createMemoryRouter, MemoryRouter, Route, Routes } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { Component } from './index'
import { Panel } from './Panel'

const proyecto = {
  id: 3,
  nombre: 'Proyecto Alfa',
  descripcion: 'Descripción',
  fecha_inicio: '2026-10-01',
  fecha_fin: '2026-11-01',
  creado_en: '2026-09-01T00:00:00Z',
}

const proyectoBeta = { ...proyecto, id: 4, nombre: 'Proyecto Beta' }

function abrir() {
  render(
    <MemoryRouter initialEntries={['/projects/3/proyecto']}>
      <Routes>
        <Route path="/projects/:id/proyecto" element={<Panel />} />
      </Routes>
    </MemoryRouter>,
  )
}

function abrirNavegable() {
  const router = createMemoryRouter([{ path: '/projects/:id/proyecto/*', element: <Component /> }], {
    initialEntries: ['/projects/3/proyecto'],
  })
  render(<RouterProvider router={router} />)
  return router
}

type EstadoDelNuevo = 'En curso' | 'pendiente' | 'error'

function responderCambio(estadoDelNuevo: EstadoDelNuevo) {
  const fetchMock = vi.fn((url: string) => {
    if (url === '/projects/3') {
      return Promise.resolve(new Response(JSON.stringify(proyecto), { status: 200 }))
    }
    if (url === '/projects/3/status') {
      return Promise.resolve(
        new Response(JSON.stringify({ estado: 'Planificado' }), { status: 200 }),
      )
    }
    if (url === '/projects/4') {
      return Promise.resolve(new Response(JSON.stringify(proyectoBeta), { status: 200 }))
    }
    if (url === '/projects/4/status') {
      if (estadoDelNuevo === 'pendiente') return new Promise<Response>(() => {})
      if (estadoDelNuevo === 'error') {
        return Promise.resolve(
          new Response(JSON.stringify({ mensaje: 'no se pudo calcular el estado' }), {
            status: 500,
          }),
        )
      }
      return Promise.resolve(
        new Response(JSON.stringify({ estado: estadoDelNuevo }), { status: 200 }),
      )
    }
    return Promise.reject(new Error(`URL inesperada: ${url}`))
  })
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

function responderEstado(estado: string | null) {
  const fetchMock = vi.fn((url: string) => {
    if (url === '/projects/3') {
      return Promise.resolve(new Response(JSON.stringify(proyecto), { status: 200 }))
    }
    if (url === '/projects/3/status') {
      if (estado === null) {
        return Promise.resolve(
          new Response(JSON.stringify({ mensaje: 'no se pudo calcular el estado' }), {
            status: 500,
          }),
        )
      }
      return Promise.resolve(new Response(JSON.stringify({ estado }), { status: 200 }))
    }
    return Promise.reject(new Error(`URL inesperada: ${url}`))
  })
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('panel del proyecto', () => {
  it.each(['Planificado', 'En curso', 'Finalizado'])(
    'muestra el estado %s que devuelve la API',
    async (estado) => {
      const fetchMock = responderEstado(estado)
      abrir()

      expect(await screen.findByText('Proyecto Alfa')).toBeTruthy()
      expect(await screen.findByText(estado)).toBeTruthy()
      expect(fetchMock.mock.calls.some(([url]) => url === '/projects/3/status')).toBe(true)
    },
  )

  it('muestra el error del estado sin ocultar los datos del proyecto', async () => {
    responderEstado(null)
    abrir()

    expect(await screen.findByText('Proyecto Alfa')).toBeTruthy()
    expect(await screen.findByText('no se pudo calcular el estado')).toBeTruthy()
    expect(await screen.findByText('Sin datos')).toBeTruthy()
  })

  it('vuelve a consultar el estado al cambiar el id del proyecto', async () => {
    const fetchMock = responderCambio('En curso')
    const router = abrirNavegable()

    expect(await screen.findByText('Planificado')).toBeTruthy()
    await router.navigate('/projects/4/proyecto')

    expect(await screen.findByText('Proyecto Beta')).toBeTruthy()
    expect(await screen.findByText('En curso')).toBeTruthy()
    expect(fetchMock.mock.calls.map(([url]) => url)).toContain('/projects/4/status')
  })

  it('muestra Consultando… mientras carga el estado del nuevo proyecto', async () => {
    responderCambio('pendiente')
    const router = abrirNavegable()

    expect(await screen.findByText('Planificado')).toBeTruthy()
    await router.navigate('/projects/4/proyecto')

    expect(await screen.findByText('Proyecto Beta')).toBeTruthy()
    expect(screen.getByText('Consultando…')).toBeTruthy()
    expect(screen.queryByText('Planificado')).toBeNull()
  })

  it('muestra Sin datos si la consulta del estado falla', async () => {
    responderCambio('error')
    const router = abrirNavegable()

    expect(await screen.findByText('Planificado')).toBeTruthy()
    await router.navigate('/projects/4/proyecto')

    expect(await screen.findByText('Proyecto Beta')).toBeTruthy()
    expect(await screen.findByText('Sin datos')).toBeTruthy()
    expect(screen.getByText('no se pudo calcular el estado')).toBeTruthy()
    expect(screen.queryByText('Planificado')).toBeNull()
  })

  it('muestra el error de un proyecto inexistente y no inventa un estado', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn(() =>
        Promise.resolve(
          new Response(JSON.stringify({ mensaje: 'no existe el proyecto' }), { status: 404 }),
        ),
      ),
    )
    abrir()

    expect(await screen.findByText('no existe el proyecto')).toBeTruthy()
    expect(screen.queryByText('Planificado')).toBeNull()
    expect(screen.queryByText('Sin datos')).toBeNull()
  })

  it('muestra tal cual un estado fuera de los tres esperados', async () => {
    responderEstado('Suspendido')
    abrir()

    expect(await screen.findByText('Proyecto Alfa')).toBeTruthy()
    expect(await screen.findByText('Suspendido')).toBeTruthy()
  })
})
