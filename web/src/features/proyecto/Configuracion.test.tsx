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

const registrados = [{ id: 1, nombre: 'Ana' }, { id: 2, nombre: 'Jimena' }]

function responderListaYAlta() {
  const fetchMock = vi.fn((url: string, init?: RequestInit) => {
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
    if (url === '/users') {
      return Promise.resolve(new Response(JSON.stringify(registrados), { status: 200 }))
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

  it('valida que se elija un usuario antes de asignar', async () => {
    const fetchMock = responderListaYAlta()
    abrir()

    await screen.findByText('Ana')
    await screen.findByLabelText('Usuario')
    await userEvent.click(screen.getByRole('button', { name: 'Asignar integrante' }))

    expect(screen.getByText('Elegí un usuario de la lista')).toBeTruthy()
    expect(fetchMock.mock.calls.filter(([, init]) => init?.method === 'POST')).toHaveLength(0)
  })
})

describe('configuración del equipo - elegir de la lista de usuarios (HU-21)', () => {
  const usuarios = [
    { id: 1, nombre: 'Ana' },
    { id: 2, nombre: 'Bruno' },
    { id: 3, nombre: 'Carla' },
  ]

  function responderUsuarios(alta: Response) {
    const fetchMock = vi.fn((url: string, init?: RequestInit) => {
      if (init?.method === 'POST') return Promise.resolve(alta)
      if (url === '/users') return Promise.resolve(new Response(JSON.stringify(usuarios)))
      if (url === '/projects/3') {
        return Promise.resolve(new Response(JSON.stringify({ id: 3, rol: 'scrum_master' })))
      }
      return Promise.resolve(new Response(JSON.stringify(integrantes)))
    })
    vi.stubGlobal('fetch', fetchMock)
    return fetchMock
  }

  it('Escenario 8: elijo de la lista un usuario registrado y entra como Product Builder', async () => {
    const fetchMock = responderUsuarios(
      new Response(
        JSON.stringify({
          id: 2,
          proyecto_id: 3,
          nombre: 'Bruno',
          rol: 'product_builder',
          vinculado_en: '2026-09-02T00:00:00Z',
        }),
        { status: 201 },
      ),
    )
    abrir()

    await screen.findByText('Ana')
    const selector = await screen.findByLabelText('Usuario')
    expect(fetchMock).toHaveBeenCalledWith('/users', expect.anything())
    expect(screen.queryByLabelText('Rol')).toBeNull()

    // Ana ya integra el proyecto: no se ofrece
    const opciones = Array.from(selector.querySelectorAll('option')).map((o) => o.textContent)
    expect(opciones).toContain('Bruno')
    expect(opciones).toContain('Carla')
    expect(opciones).not.toContain('Ana')

    await userEvent.selectOptions(selector, 'Bruno')
    await userEvent.click(screen.getByRole('button', { name: 'Asignar integrante' }))

    await waitFor(() =>
      expect(fetchMock.mock.calls.some(([, init]) => init?.method === 'POST')).toBe(true),
    )
    const llamada = fetchMock.mock.calls.find(([, init]) => init?.method === 'POST')
    expect(llamada?.[0]).toBe('/projects/3/members')
    expect(JSON.parse(String(llamada?.[1]?.body))).toEqual({ nombre: 'Bruno' })
    expect(await screen.findByText('Product Builder')).toBeTruthy()
  })

  it('muestra el error del servidor si el usuario no existe', async () => {
    responderUsuarios(
      new Response(JSON.stringify({ mensaje: 'usuario no encontrado' }), { status: 404 }),
    )
    abrir()

    await userEvent.selectOptions(await screen.findByLabelText('Usuario'), 'Bruno')
    await userEvent.click(screen.getByRole('button', { name: 'Asignar integrante' }))

    expect(await screen.findByText('usuario no encontrado')).toBeTruthy()
  })

  it('avisa cuando no hay usuarios registrados para agregar', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn((url: string) =>
        Promise.resolve(
          new Response(JSON.stringify(url === '/users' ? [] : integrantes)),
        ),
      ),
    )
    abrir()

    expect(await screen.findByText('No hay usuarios registrados para agregar')).toBeTruthy()
  })
})

describe('configuración del equipo - rol (HU-21)', () => {
  const nota = (boton: HTMLElement) =>
    document.getElementById(boton.getAttribute('aria-describedby') ?? '')?.textContent

  function responderConRol(rol: 'scrum_master' | 'product_builder') {
    const fetchMock = vi.fn((url: string) =>
      Promise.resolve(
        new Response(JSON.stringify(url === '/projects/3' ? { id: 3, rol } : integrantes), { status: 200 }),
      ),
    )
    vi.stubGlobal('fetch', fetchMock)
    return fetchMock
  }

  it('Escenario 2: como Product Builder veo «Asignar integrante» deshabilitado con la nota «Solo el Scrum Master puede agregar integrantes»', async () => {
    responderConRol('product_builder')
    abrir()

    await screen.findByText('Ana')
    const boton = screen.getByRole('button', { name: 'Asignar integrante' })
    await waitFor(() => expect(boton.hasAttribute('disabled')).toBe(true))
    expect(nota(boton)).toBe('Solo el Scrum Master puede agregar integrantes')
  })

  it('Escenario 1: como Scrum Master «Asignar integrante» está habilitado', async () => {
    const fetchMock = responderConRol('scrum_master')
    abrir()

    await screen.findByText('Ana')
    await waitFor(() => expect(fetchMock).toHaveBeenCalledWith('/projects/3', expect.anything()))
    const boton = screen.getByRole('button', { name: 'Asignar integrante' })
    expect(boton.hasAttribute('disabled')).toBe(false)
    expect(boton.hasAttribute('aria-describedby')).toBe(false)
  })
})
