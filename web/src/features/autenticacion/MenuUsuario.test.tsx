import { cleanup, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Route, Routes } from 'react-router'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { borrarToken, guardarToken, leerToken } from './almacenamientoToken'
import { MenuUsuario } from './MenuUsuario'
import { SesionProvider } from './SesionContext'

const USUARIO = { id: 7, nombre: 'ana.perez' }

function json(cuerpo: unknown, status: number) {
  return new Response(JSON.stringify(cuerpo), { status })
}

function montar() {
  return render(
    <MemoryRouter initialEntries={['/']}>
      <SesionProvider>
        <MenuUsuario />
        <Routes>
          <Route path="/ingresar" element={<h1>Pantalla de ingreso</h1>} />
          <Route path="*" element={<p>Inicio</p>} />
        </Routes>
      </SesionProvider>
    </MemoryRouter>,
  )
}

beforeEach(() => {
  borrarToken()
  localStorage.clear()
})

afterEach(() => {
  cleanup()
  vi.unstubAllGlobals()
})

describe('MenuUsuario', () => {
  it('con sesión muestra el nombre del usuario y “Cerrar sesión”', async () => {
    guardarToken('token-guardado')
    vi.stubGlobal('fetch', vi.fn(async () => json(USUARIO, 200)))
    montar()

    expect(await screen.findByText('ana.perez')).toBeTruthy()
    expect(screen.getByRole('button', { name: 'Cerrar sesión' })).toBeTruthy()
  })

  it('Escenario 7: elegir “Cerrar sesión” cierra la sesión y lleva a /ingresar', async () => {
    guardarToken('token-guardado')
    vi.stubGlobal(
      'fetch',
      vi.fn(async (url: string) => (String(url).endsWith('/auth/yo') ? json(USUARIO, 200) : new Response(null, { status: 204 }))),
    )
    montar()
    const usuario = userEvent.setup()

    await usuario.click(await screen.findByRole('button', { name: 'Cerrar sesión' }))

    expect(await screen.findByRole('heading', { name: 'Pantalla de ingreso' })).toBeTruthy()
    expect(leerToken()).toBeNull()
  })

  it('sin sesión ofrece enlaces a ingresar y registrarse', () => {
    montar()

    expect(screen.getByRole('link', { name: 'Ingresar' }).getAttribute('href')).toBe('/ingresar')
    expect(screen.getByRole('link', { name: 'Registrarse' }).getAttribute('href')).toBe('/registro')
  })
})

describe('MenuUsuario - cambiar nombre (HU-17)', () => {
  function preparar(actualizado = USUARIO, putStatus = 200, putCuerpo: unknown = actualizado) {
    guardarToken('token-guardado')
    const fetchMock = vi.fn(async (url: string, init?: RequestInit) => {
      if (url === '/auth/yo' && init?.method === 'PUT') return json(putCuerpo, putStatus)
      return json(USUARIO, 200)
    })
    vi.stubGlobal('fetch', fetchMock)
    montar()
    return fetchMock
  }

  async function abrirFormulario() {
    await screen.findByText(USUARIO.nombre)
    await userEvent.click(screen.getByRole('button', { name: 'Cambiar nombre' }))
    return screen.getByLabelText('Nuevo nombre de usuario') as HTMLInputElement
  }

  it('Given una sesión, When elijo Cambiar nombre, Then veo el formulario con el nombre actual', async () => {
    preparar()

    const input = await abrirFormulario()

    expect(input.value).toBe(USUARIO.nombre)
  })

  it('Given un nombre válido, When guardo, Then se envía PUT /auth/yo y el menú muestra el nombre nuevo', async () => {
    const actualizado = { id: 7, nombre: 'ana.gomez' }
    const fetchMock = preparar(actualizado)
    const input = await abrirFormulario()

    await userEvent.clear(input)
    await userEvent.type(input, 'ana.gomez')
    await userEvent.click(screen.getByRole('button', { name: 'Guardar' }))

    expect(await screen.findByText('ana.gomez')).toBeTruthy()
    expect(fetchMock).toHaveBeenCalledWith(
      '/auth/yo',
      expect.objectContaining({
        method: 'PUT',
        body: JSON.stringify({ nombre: 'ana.gomez' }),
      }),
    )
    expect(screen.queryByLabelText('Nuevo nombre de usuario')).toBeNull()
  })

  it('Given un nombre en uso, When guardo, Then veo el mensaje y el menú conserva el nombre anterior', async () => {
    preparar(USUARIO, 409, { mensaje: 'Ese nombre de usuario ya existe' })
    const input = await abrirFormulario()

    await userEvent.clear(input)
    await userEvent.type(input, 'bruno')
    await userEvent.click(screen.getByRole('button', { name: 'Guardar' }))

    expect(await screen.findByText('Ese nombre de usuario ya existe')).toBeTruthy()
    expect(screen.getByText(USUARIO.nombre)).toBeTruthy()
    expect(screen.getByLabelText('Nuevo nombre de usuario')).toBeTruthy()
  })

  it('Given el servidor rechaza el nombre, When guardo, Then veo el mensaje junto al campo y conservo lo escrito', async () => {
    const mensaje = 'el nombre de usuario debe tener entre 3 y 30 caracteres'
    preparar(USUARIO, 400, { campo: 'nombre', mensaje })
    const input = await abrirFormulario()

    await userEvent.clear(input)
    await userEvent.type(input, 'ab')
    await userEvent.click(screen.getByRole('button', { name: 'Guardar' }))

    expect(await screen.findByText(mensaje)).toBeTruthy()
    expect((screen.getByLabelText('Nuevo nombre de usuario') as HTMLInputElement).value).toBe('ab')
  })

  it('Given el formulario abierto, When cancelo, Then no se llama a la API y el menú sigue igual', async () => {
    const fetchMock = preparar()
    await abrirFormulario()

    await userEvent.click(screen.getByRole('button', { name: 'Cancelar' }))

    expect(screen.queryByLabelText('Nuevo nombre de usuario')).toBeNull()
    expect(fetchMock.mock.calls.filter(([, init]) => init?.method === 'PUT')).toHaveLength(0)
    expect(screen.getByText(USUARIO.nombre)).toBeTruthy()
  })
})
