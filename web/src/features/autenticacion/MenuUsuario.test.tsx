import { cleanup, render, screen, within } from '@testing-library/react'
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
  it('con sesión muestra el nombre del usuario en el botón del menú, que empieza cerrado', async () => {
    guardarToken('token-guardado')
    vi.stubGlobal('fetch', vi.fn(async () => json(USUARIO, 200)))
    montar()

    expect(await screen.findByText('ana.perez')).toBeTruthy()
    const boton = screen.getByRole('button', { name: 'Menú de usuario' })
    expect(boton.getAttribute('aria-expanded')).toBe('false')
    expect(screen.queryByRole('button', { name: 'Cerrar sesión' })).toBeNull()
  })

  it('Escenario 7: elegir “Cerrar sesión” cierra la sesión y lleva a /ingresar', async () => {
    guardarToken('token-guardado')
    vi.stubGlobal(
      'fetch',
      vi.fn(async (url: string) => (String(url).endsWith('/auth/yo') ? json(USUARIO, 200) : new Response(null, { status: 204 }))),
    )
    montar()
    const usuario = userEvent.setup()

    await usuario.click(await screen.findByRole('button', { name: 'Menú de usuario' }))
    await usuario.click(screen.getByRole('button', { name: 'Cerrar sesión' }))

    expect(await screen.findByRole('heading', { name: 'Pantalla de ingreso' })).toBeTruthy()
    expect(leerToken()).toBeNull()
  })

  it('sin sesión ofrece enlaces a ingresar y registrarse', () => {
    montar()

    expect(screen.getByRole('link', { name: 'Ingresar' }).getAttribute('href')).toBe('/ingresar')
    expect(screen.getByRole('link', { name: 'Registrarse' }).getAttribute('href')).toBe('/registro')
  })

  async function abrirMenu() {
    guardarToken('token-guardado')
    vi.stubGlobal('fetch', vi.fn(async () => json(USUARIO, 200)))
    montar()
    const usuario = userEvent.setup()
    const boton = await screen.findByRole('button', { name: 'Menú de usuario' })
    await usuario.click(boton)
    return { usuario, boton }
  }

  it('al abrirlo muestra el nombre, «Sesión iniciada» y «Cerrar sesión», sin rol', async () => {
    const { boton } = await abrirMenu()

    expect(boton.getAttribute('aria-expanded')).toBe('true')
    const panel = screen.getByRole('group', { name: 'Cuenta' })
    expect(within(panel).getByText('ana.perez')).toBeTruthy()
    expect(within(panel).getByText('Sesión iniciada')).toBeTruthy()
    expect(within(panel).getByRole('button', { name: 'Cerrar sesión' })).toBeTruthy()
    expect(screen.queryByText(/Tu rol/)).toBeNull()
    expect(screen.queryByText(/Solo del prototipo/)).toBeNull()
  })

  it('el botón enlaza con su panel y vuelve a cerrarlo al pulsarlo de nuevo', async () => {
    const { usuario, boton } = await abrirMenu()

    expect(document.getElementById(boton.getAttribute('aria-controls') ?? '')).toBe(screen.getByRole('group', { name: 'Cuenta' }))
    await usuario.click(boton)

    expect(screen.queryByRole('group', { name: 'Cuenta' })).toBeNull()
    expect(boton.getAttribute('aria-expanded')).toBe('false')
  })

  it('Escape lo cierra y devuelve el foco al botón', async () => {
    const { usuario, boton } = await abrirMenu()

    await usuario.keyboard('{Escape}')

    expect(screen.queryByRole('group', { name: 'Cuenta' })).toBeNull()
    expect(document.activeElement).toBe(boton)
  })

  it('pulsar fuera lo cierra', async () => {
    const { usuario } = await abrirMenu()

    await usuario.click(screen.getByText('Inicio'))

    expect(screen.queryByRole('group', { name: 'Cuenta' })).toBeNull()
  })

  it('pulsar dentro del panel, fuera de «Cerrar sesión», no lo cierra', async () => {
    const { usuario } = await abrirMenu()

    await usuario.click(screen.getByText('Sesión iniciada'))

    expect(screen.getByRole('group', { name: 'Cuenta' })).toBeTruthy()
  })
})
