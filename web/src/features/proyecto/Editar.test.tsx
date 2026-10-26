import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Route, Routes } from 'react-router'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { Editar } from './Editar'

const proyecto = {
  id: 3,
  nombre: 'Proyecto Alfa',
  descripcion: 'Descripción original',
  fecha_inicio: '2026-10-01',
  fecha_fin: '2026-11-01',
  creado_en: '2026-09-01T00:00:00Z',
}

function abrir() {
  render(
    <MemoryRouter initialEntries={['/projects/3/proyecto/edit']}>
      <Routes>
        <Route path="/projects/:id/proyecto/edit" element={<Editar />} />
      </Routes>
    </MemoryRouter>,
  )
}

function responderProyecto() {
  const fetchMock = vi.fn().mockResolvedValue(new Response(JSON.stringify(proyecto), { status: 200 }))
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('pantalla de edición de proyecto', () => {
  it('carga los datos actuales del proyecto', async () => {
    responderProyecto()
    abrir()

    expect(await screen.findByDisplayValue('Proyecto Alfa')).toBeTruthy()
    expect(screen.getByDisplayValue('Descripción original')).toBeTruthy()
    expect(screen.getByDisplayValue('2026-10-01')).toBeTruthy()
  })

  it('muestra el error de nombre vacío y no guarda', async () => {
    const fetchMock = responderProyecto()
    abrir()

    const nombre = await screen.findByDisplayValue('Proyecto Alfa')
    await userEvent.clear(nombre)
    await userEvent.click(screen.getByRole('button', { name: 'Guardar cambios' }))

    expect(screen.getByText('El nombre es obligatorio')).toBeTruthy()
    expect(fetchMock.mock.calls.filter(([, init]) => init?.method === 'PUT')).toHaveLength(0)
  })

  it('muestra el error cuando la fecha de fin es igual a la de inicio', async () => {
    const fetchMock = responderProyecto()
    abrir()

    await screen.findByDisplayValue('Proyecto Alfa')
    fireEvent.change(screen.getByLabelText('Fecha de fin'), { target: { value: '2026-10-01' } })
    await userEvent.click(screen.getByRole('button', { name: 'Guardar cambios' }))

    expect(
      screen.getByText('La fecha de fin debe ser posterior a la fecha de inicio'),
    ).toBeTruthy()
    expect(fetchMock.mock.calls.filter(([, init]) => init?.method === 'PUT')).toHaveLength(0)
  })

  it('guarda los cambios válidos con PUT', async () => {
    const fetchMock = vi.fn((_url: string, init?: RequestInit) =>
      Promise.resolve(
        new Response(JSON.stringify(proyecto), {
          status: init?.method === 'PUT' ? 200 : 200,
        }),
      ),
    )
    vi.stubGlobal('fetch', fetchMock)
    abrir()

    const nombre = await screen.findByDisplayValue('Proyecto Alfa')
    await userEvent.clear(nombre)
    await userEvent.type(nombre, 'Proyecto Beta')
    await userEvent.click(screen.getByRole('button', { name: 'Guardar cambios' }))

    await waitFor(() =>
      expect(fetchMock.mock.calls.some(([, init]) => init?.method === 'PUT')).toBe(true),
    )

    const llamada = fetchMock.mock.calls.find(([, init]) => init?.method === 'PUT')
    expect(JSON.parse(String(llamada?.[1]?.body))).toEqual({
      nombre: 'Proyecto Beta',
      descripcion: 'Descripción original',
      fecha_inicio: '2026-10-01',
      fecha_fin: '2026-11-01',
    })
  })
})
