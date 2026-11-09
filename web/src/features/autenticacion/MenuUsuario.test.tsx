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
