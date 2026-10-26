import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { Component } from './crear'

function abrir() {
  render(
    <MemoryRouter>
      <Component />
    </MemoryRouter>,
  )
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('pantalla de creación de proyecto', () => {
  it('muestra los errores de validación y no llama a la API', async () => {
    const fetchMock = vi.fn()
    vi.stubGlobal('fetch', fetchMock)
    abrir()

    await userEvent.click(screen.getByRole('button', { name: 'Crear proyecto' }))

    expect(screen.getByText('El nombre es obligatorio')).toBeTruthy()
    expect(screen.getByText('El nombre del creador es obligatorio')).toBeTruthy()
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it('envía las fechas opcionales vacías como null y navega al panel', async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      new Response(
        JSON.stringify({
          id: 7,
          nombre: 'Proyecto Alfa',
          descripcion: '',
          fecha_inicio: null,
          fecha_fin: null,
          creado_en: '2026-10-23T00:00:00Z',
        }),
        { status: 201 },
      ),
    )
    vi.stubGlobal('fetch', fetchMock)
    abrir()

    await userEvent.type(screen.getByLabelText('Nombre del proyecto'), 'Proyecto Alfa')
    await userEvent.type(screen.getByLabelText('Tu nombre (Scrum Master)'), 'Ana')
    await userEvent.click(screen.getByRole('button', { name: 'Crear proyecto' }))

    await waitFor(() => expect(fetchMock).toHaveBeenCalledTimes(1))
    expect(fetchMock).toHaveBeenCalledWith(
      '/projects',
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({
          nombre: 'Proyecto Alfa',
          descripcion: '',
          creador: 'Ana',
          fecha_inicio: null,
          fecha_fin: null,
        }),
      }),
    )
  })

  it('muestra como error general el campo "nombre" que devuelve la API al crear', async () => {
    const mensaje = 'el nombre del integrante no puede superar los 200 caracteres'
    const fetchMock = vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ campo: 'nombre', mensaje }), { status: 400 }),
    )
    vi.stubGlobal('fetch', fetchMock)
    abrir()

    await userEvent.type(screen.getByLabelText('Nombre del proyecto'), 'Proyecto Alfa')
    await userEvent.type(screen.getByLabelText('Tu nombre (Scrum Master)'), 'Ana')
    await userEvent.click(screen.getByRole('button', { name: 'Crear proyecto' }))

    expect(await screen.findByText(mensaje)).toBeTruthy()
    const nombre = screen.getByLabelText('Nombre del proyecto')
    expect(nombre.getAttribute('aria-describedby')).toBeNull()
  })
})
