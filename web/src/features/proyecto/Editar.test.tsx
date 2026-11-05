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
  horas_por_story_point: 4,
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
      horas_por_story_point: 4,
    })
  })

  describe('horas por Story Point', () => {
    function responderEdicion(respuestaPut?: Response) {
      const fetchMock = vi.fn((_url: string, init?: RequestInit) =>
        Promise.resolve(
          init?.method === 'PUT' && respuestaPut
            ? respuestaPut
            : new Response(JSON.stringify(proyecto), { status: 200 }),
        ),
      )
      vi.stubGlobal('fetch', fetchMock)
      return fetchMock
    }

    async function guardarCon(valor: string, fetchMock: ReturnType<typeof responderEdicion>) {
      const campo = await screen.findByLabelText(/Horas por Story Point/)
      await userEvent.clear(campo)
      if (valor) await userEvent.type(campo, valor)
      await userEvent.click(screen.getByRole('button', { name: 'Guardar cambios' }))
      await waitFor(() =>
        expect(fetchMock.mock.calls.some(([, init]) => init?.method === 'PUT')).toBe(true),
      )
      const llamada = fetchMock.mock.calls.find(([, init]) => init?.method === 'PUT')
      return JSON.parse(String(llamada?.[1]?.body))
    }

    it('US4.1: dado un proyecto con el parámetro, al abrir la edición el campo muestra el valor actual', async () => {
      responderEdicion()
      abrir()

      const campo = await screen.findByLabelText<HTMLInputElement>(/Horas por Story Point/)
      expect(campo.value).toBe('4')
    })

    it('US4.1: dado un proyecto sin el parámetro, el campo aparece vacío', async () => {
      vi.stubGlobal(
        'fetch',
        vi.fn().mockResolvedValue(
          new Response(JSON.stringify({ ...proyecto, horas_por_story_point: null }), { status: 200 }),
        ),
      )
      abrir()

      const campo = await screen.findByLabelText<HTMLInputElement>(/Horas por Story Point/)
      expect(campo.value).toBe('')
    })

    it('US4.2: al cargar un valor válido y guardar, se envía el número', async () => {
      const fetchMock = responderEdicion()
      abrir()

      const cuerpo = await guardarCon('2,5', fetchMock)

      expect(cuerpo.horas_por_story_point).toBe(2.5)
    })

    it('US4.3: si el sistema rechaza el valor, la pantalla muestra el mensaje junto al campo', async () => {
      const mensaje = 'las horas por Story Point deben ser un número mayor que cero'
      const fetchMock = responderEdicion(
        new Response(JSON.stringify({ campo: 'horas_por_story_point', mensaje }), { status: 400 }),
      )
      abrir()

      await guardarCon('8', fetchMock)

      const campo = await screen.findByLabelText(/Horas por Story Point/)
      const error = await screen.findByText(mensaje)
      expect(campo.getAttribute('aria-describedby')).toBe(error.id)
    })

    it('US4.3: un valor inválido muestra el error y no se envía', async () => {
      const fetchMock = responderEdicion()
      abrir()

      const campo = await screen.findByLabelText(/Horas por Story Point/)
      await userEvent.clear(campo)
      await userEvent.type(campo, 'abc')
      await userEvent.click(screen.getByRole('button', { name: 'Guardar cambios' }))

      expect(
        screen.getByText('Las horas por Story Point deben ser un número mayor que cero'),
      ).toBeTruthy()
      expect(fetchMock.mock.calls.filter(([, init]) => init?.method === 'PUT')).toHaveLength(0)
    })

    it('US4.4: al dejar el campo vacío y guardar, se envía sin valor', async () => {
      const fetchMock = responderEdicion()
      abrir()

      const cuerpo = await guardarCon('', fetchMock)

      expect(cuerpo.horas_por_story_point).toBeNull()
    })
  })
})
