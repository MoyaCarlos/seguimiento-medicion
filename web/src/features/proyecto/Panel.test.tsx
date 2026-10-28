import { render, screen } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { Panel } from './Panel'

const proyecto = {
  id: 3,
  nombre: 'Proyecto Alfa',
  descripcion: 'Descripción',
  fecha_inicio: '2026-10-01',
  fecha_fin: '2026-11-01',
  creado_en: '2026-09-01T00:00:00Z',
}

function abrir() {
  render(
    <MemoryRouter initialEntries={['/projects/3/proyecto']}>
      <Routes>
        <Route path="/projects/:id/proyecto" element={<Panel />} />
      </Routes>
    </MemoryRouter>,
  )
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
  })
})
