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

function montar(entrada: string | { pathname: string; state: unknown } = '/ingresar') {
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

beforeEach(() => {
  borrarToken()
  localStorage.clear()
})

afterEach(() => {
  cleanup()
  vi.unstubAllGlobals()
})

describe('pantalla de ingreso', () => {
  it('Escenario 5: el error de credenciales se muestra en el formulario sin borrar lo escrito', async () => {
    simularFetch(() => json({ mensaje: 'Usuario o contraseña incorrectos' }, 401))
    montar()
    const usuario = userEvent.setup()

    await usuario.type(screen.getByLabelText('Nombre de usuario'), 'ana.perez')
    await usuario.type(screen.getByLabelText('Contraseña'), 'clave-equivocada')
    await usuario.click(screen.getByRole('button', { name: 'Ingresar' }))

    expect((await screen.findByRole('alert')).textContent).toBe('Usuario o contraseña incorrectos')
    expect((screen.getByLabelText('Nombre de usuario') as HTMLInputElement).value).toBe('ana.perez')
    expect((screen.getByLabelText('Contraseña') as HTMLInputElement).value).toBe('clave-equivocada')
  })

  it('Escenario 6: con state.vencida se muestra “Tu sesión venció”', () => {
    montar({ pathname: '/ingresar', state: { desde: '/backlog', vencida: true } })

    expect(screen.getByText('Tu sesión venció')).toBeTruthy()
  })

  it('Escenario 6: sin state.vencida no se muestra “Tu sesión venció”', () => {
    montar()

    expect(screen.queryByText('Tu sesión venció')).toBeNull()
  })

  it('Escenario 4: ingresar con datos válidos navega a la pantalla de origen', async () => {
    simularFetch(() => json(RESPUESTA_SESION, 200))
    montar({ pathname: '/ingresar', state: { desde: '/projects/3/backlog' } })
    const usuario = userEvent.setup()

    await usuario.type(screen.getByLabelText('Nombre de usuario'), 'ana.perez')
    await usuario.type(screen.getByLabelText('Contraseña'), 'una-clave-larga')
    await usuario.click(screen.getByRole('button', { name: 'Ingresar' }))

    await waitFor(() => expect((screen.getByTestId('ubicacion')).textContent).toBe('/projects/3/backlog'))
    expect(leerToken()).toBe('token-nuevo')
  })

  it('la contraseña se escribe en un campo type="password"', () => {
    montar()

    expect(screen.getByLabelText('Contraseña').getAttribute('type')).toBe('password')
  })

  it('tiene un enlace a la pantalla de registro', async () => {
    montar()
    const usuario = userEvent.setup()

    await usuario.click(screen.getByRole('link', { name: 'Registrate' }))

    expect(screen.getByRole('heading', { name: 'Crear cuenta' })).toBeTruthy()
  })
})
