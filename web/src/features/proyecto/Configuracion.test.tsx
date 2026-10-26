import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Route, Routes } from 'react-router'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { Configuracion } from './Configuracion'

const integrantes = [
  {
    id: 1,
    proyecto_id: 3,
    nombre: 'Ana',
    rol: 'scrum_master',
    vinculado_en: '2026-09-01T00:00:00Z',
  },
]

function abrir() {
  render(
    <MemoryRouter initialEntries={['/projects/3/proyecto/settings']}>
      <Routes>
        <Route path="/projects/:id/proyecto/settings" element={<Configuracion />} />
      </Routes>
    </MemoryRouter>,
  )
}

function responderListaYAlta() {
  const fetchMock = vi.fn((_url: string, init?: RequestInit) => {
    if (init?.method === 'POST') {
      return Promise.resolve(
        new Response(
          JSON.stringify({
            id: 2,
            proyecto_id: 3,
            nombre: 'Jimena',
            rol: 'product_builder',
            vinculado_en: '2026-09-02T00:00:00Z',
          }),
          { status: 201 },
        ),
      )
    }
    return Promise.resolve(new Response(JSON.stringify(integrantes), { status: 200 }))
  })
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('pantalla de configuración del equipo', () => {
  it('lista los integrantes con su rol', async () => {
    responderListaYAlta()
    abrir()

    expect(await screen.findByText('Ana')).toBeTruthy()
    expect(screen.getByText('Ana').closest('li')?.textContent).toContain('Scrum Master')
  })

  it('valida nombre y rol antes de asignar', async () => {
    const fetchMock = responderListaYAlta()
    abrir()

    await screen.findByText('Ana')
    await userEvent.click(screen.getByRole('button', { name: 'Asignar integrante' }))

    expect(screen.getByText('El nombre del integrante es obligatorio')).toBeTruthy()
    expect(screen.getByText('Seleccioná un rol válido')).toBeTruthy()
    expect(fetchMock.mock.calls.filter(([, init]) => init?.method === 'POST')).toHaveLength(0)
  })

  it('asigna un integrante con rol y lo muestra en el equipo', async () => {
    const fetchMock = responderListaYAlta()
    abrir()

    await screen.findByText('Ana')
    await userEvent.type(screen.getByLabelText('Nombre del integrante'), 'Jimena')
    await userEvent.selectOptions(screen.getByLabelText('Rol'), 'product_builder')
    await userEvent.click(screen.getByRole('button', { name: 'Asignar integrante' }))

    expect(await screen.findByText('Jimena')).toBeTruthy()

    await waitFor(() =>
      expect(fetchMock.mock.calls.some(([, init]) => init?.method === 'POST')).toBe(true),
    )
    const llamada = fetchMock.mock.calls.find(([, init]) => init?.method === 'POST')
    expect(JSON.parse(String(llamada?.[1]?.body))).toEqual({
      nombre: 'Jimena',
      rol: 'product_builder',
    })
  })
})
