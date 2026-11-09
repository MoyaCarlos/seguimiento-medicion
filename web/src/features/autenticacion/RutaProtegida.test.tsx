import { cleanup, render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Route, Routes, useLocation } from 'react-router'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { borrarToken, guardarToken, leerToken } from './almacenamientoToken'
import { Ingresar } from './Ingresar'
import { RutaProtegida } from './RutaProtegida'
import { SesionProvider } from './SesionContext'

const USUARIO = { id: 7, nombre: 'ana.perez' }
const RESPUESTA_SESION = { usuario: USUARIO, token: 'token-nuevo', vence_en: '2026-11-10T12:00:00Z' }

function json(cuerpo: unknown, status: number) {
  return new Response(JSON.stringify(cuerpo), { status })
}

function Ubicacion() {
  const { pathname, state } = useLocation()
  return (
    <span data-testid="ubicacion">
      {pathname}|{(state as { desde?: string } | null)?.desde ?? ''}
    </span>
  )
}

function montar(entrada: string) {
  return render(
    <MemoryRouter initialEntries={[entrada]}>
      <SesionProvider>
        <Routes>
          <Route path="/ingresar" element={<><Ubicacion /><Ingresar /></>} />
          <Route element={<RutaProtegida />}>
            <Route path="/backlog" element={<p>Contenido protegido</p>} />
          </Route>
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

describe('RutaProtegida', () => {
  it('Escenario 6: sin sesión redirige a /ingresar con state.desde', async () => {
    montar('/backlog')

    expect((await screen.findByTestId('ubicacion')).textContent).toBe('/ingresar|/backlog')
    expect(screen.queryByText('Contenido protegido')).toBeNull()
  })

  it('Escenario 6: tras ingresar vuelve a la ruta pedida', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => json(RESPUESTA_SESION, 200)))
    montar('/backlog')
    const usuario = userEvent.setup()

    await usuario.type(await screen.findByLabelText('Nombre de usuario'), 'ana.perez')
    await usuario.type(screen.getByLabelText('Contraseña'), 'una-clave-larga')
    await usuario.click(screen.getByRole('button', { name: 'Ingresar' }))

    await waitFor(() => expect(screen.getByText('Contenido protegido')).toBeTruthy())
  })

  it('mientras está cargando muestra un indicador y no el contenido', () => {
    guardarToken('token-guardado')
    vi.stubGlobal('fetch', vi.fn(() => new Promise<Response>(() => {})))
    montar('/backlog')

    expect(screen.getByRole('status')).toBeTruthy()
    expect(screen.queryByText('Contenido protegido')).toBeNull()
  })

  it('con sesión muestra el contenido', async () => {
    guardarToken('token-guardado')
    vi.stubGlobal('fetch', vi.fn(async () => json(USUARIO, 200)))
    montar('/backlog')

    expect(await screen.findByText('Contenido protegido')).toBeTruthy()
  })

  it('Escenario 8: con el servidor caído al arrancar y token guardado muestra el aviso en lugar de redirigir', async () => {
    guardarToken('token-guardado')
    vi.stubGlobal('fetch', vi.fn(async () => json({ mensaje: 'Error interno' }, 500)))
    montar('/backlog')

    expect(await screen.findByText('No pudimos conectar con el servidor')).toBeTruthy()
    expect(screen.getByRole('button', { name: 'Reintentar' })).toBeTruthy()
    expect(screen.queryByTestId('ubicacion')).toBeNull()
    expect(screen.queryByText('Contenido protegido')).toBeNull()
  })

  it('Escenario 8: Reintentar con el servidor de vuelta muestra el contenido protegido', async () => {
    guardarToken('token-guardado')
    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce(json({ mensaje: 'Error interno' }, 500))
      .mockResolvedValueOnce(json(USUARIO, 200))
    vi.stubGlobal('fetch', fetchMock)
    montar('/backlog')
    const usuario = userEvent.setup()

    await usuario.click(await screen.findByRole('button', { name: 'Reintentar' }))

    expect(await screen.findByText('Contenido protegido')).toBeTruthy()
    expect(screen.queryByText('No pudimos conectar con el servidor')).toBeNull()
  })

  it('Escenario 8: Reintentar que recibe 401 redirige a /ingresar con state.vencida', async () => {
    guardarToken('token-guardado')
    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce(json({ mensaje: 'Error interno' }, 500))
      .mockResolvedValueOnce(json({ mensaje: 'Tu sesión no es válida' }, 401))
    vi.stubGlobal('fetch', fetchMock)
    montar('/backlog')
    const usuario = userEvent.setup()

    await usuario.click(await screen.findByRole('button', { name: 'Reintentar' }))

    expect((await screen.findByTestId('ubicacion')).textContent).toBe('/ingresar|/backlog')
    expect(leerToken()).toBeNull()
  })
})
