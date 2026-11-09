import { cleanup, render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Route, Routes, useLocation } from 'react-router'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { borrarToken, leerToken } from './almacenamientoToken'
import { Ingresar } from './Ingresar'
import { Registro } from './Registro'
import { SesionProvider } from './SesionContext'

const RESPUESTA_SESION = {
  usuario: { id: 7, nombre: 'ana.perez' },
  token: 'token-nuevo',
  vence_en: '2026-11-10T12:00:00Z',
}

function json(cuerpo: unknown, status: number) {
  return new Response(JSON.stringify(cuerpo), { status })
}

function simularFetch(respuesta: () => Response) {
  const fetchMock = vi.fn(async () => respuesta())
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

function Ubicacion() {
  const { pathname } = useLocation()
  return <span data-testid="ubicacion">{pathname}</span>
}

function montar(entrada: string | { pathname: string; state: unknown } = '/registro') {
  return render(
    <MemoryRouter initialEntries={[entrada]}>
      <SesionProvider>
        <Routes>
          <Route path="/ingresar" element={<Ingresar />} />
          <Route path="/registro" element={<Registro />} />
          <Route path="*" element={<Ubicacion />} />
        </Routes>
      </SesionProvider>
    </MemoryRouter>,
  )
}

function descripcion(campo: HTMLElement) {
  const id = campo.getAttribute('aria-describedby') ?? ''
  return document.getElementById(id)?.textContent
}

async function completarYEnviar(nombre: string, contrasena: string) {
  const usuario = userEvent.setup()
  await usuario.type(screen.getByLabelText('Nombre de usuario'), nombre)
  await usuario.type(screen.getByLabelText('Contraseña'), contrasena)
  await usuario.click(screen.getByRole('button', { name: 'Crear cuenta' }))
}

beforeEach(() => {
  borrarToken()
  localStorage.clear()
})

afterEach(() => {
  cleanup()
  vi.unstubAllGlobals()
})

describe('pantalla de registro', () => {
  it('Escenario 1: registrarse con datos válidos deja la sesión iniciada y navega a la pantalla de origen', async () => {
    simularFetch(() => json(RESPUESTA_SESION, 201))
    montar({ pathname: '/registro', state: { desde: '/projects/3/backlog' } })

    await completarYEnviar('ana.perez', 'una-clave-larga')

    await waitFor(() => expect((screen.getByTestId('ubicacion')).textContent).toBe('/projects/3/backlog'))
    expect(leerToken()).toBe('token-nuevo')
  })

  it('Escenario 2: “Ese nombre de usuario ya existe” se muestra junto al campo nombre', async () => {
    simularFetch(() => json({ mensaje: 'Ese nombre de usuario ya existe' }, 409))
    montar()

    await completarYEnviar('ana.perez', 'una-clave-larga')

    await screen.findByText('Ese nombre de usuario ya existe')
    expect(descripcion(screen.getByLabelText('Nombre de usuario'))).toBe('Ese nombre de usuario ya existe')
    expect((screen.getByLabelText('Nombre de usuario') as HTMLInputElement).value).toBe('ana.perez')
  })

  it('Escenario 3: el error del campo contrasena se muestra junto a la contraseña', async () => {
    simularFetch(() =>
      json({ campo: 'contrasena', mensaje: 'La contraseña debe tener al menos 8 caracteres' }, 400),
    )
    montar()

    await completarYEnviar('ana.perez', 'corta')

    await screen.findByText('La contraseña debe tener al menos 8 caracteres')
    expect(descripcion(screen.getByLabelText('Contraseña'))).toBe('La contraseña debe tener al menos 8 caracteres')
    expect(screen.getByLabelText('Nombre de usuario').getAttribute('aria-describedby')).toBeNull()
  })

  it('Escenario 3: el error del campo nombre se muestra junto al nombre', async () => {
    simularFetch(() => json({ campo: 'nombre', mensaje: 'El nombre no es válido' }, 400))
    montar()

    await completarYEnviar('a b', 'una-clave-larga')

    await screen.findByText('El nombre no es válido')
    expect(descripcion(screen.getByLabelText('Nombre de usuario'))).toBe('El nombre no es válido')
  })

  it('la contraseña se escribe en un campo type="password"', () => {
    montar()

    expect(screen.getByLabelText('Contraseña').getAttribute('type')).toBe('password')
  })

  it('tiene un enlace a la pantalla de ingreso', async () => {
    montar()
    const usuario = userEvent.setup()

    await usuario.click(screen.getByRole('link', { name: 'Ingresá' }))

    expect(screen.getByRole('heading', { name: 'Ingresar' })).toBeTruthy()
  })
})
