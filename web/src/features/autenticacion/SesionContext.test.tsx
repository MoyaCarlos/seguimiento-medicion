import { useEffect } from 'react'
import { act, cleanup, render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter, useLocation } from 'react-router'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { ApiError, get } from '../../compartido/api/client'
import { configurarSesionHttp } from '../../compartido/api/sesionHttp'
import { borrarToken, guardarToken, leerToken } from './almacenamientoToken'
import { SesionProvider, useSesion, type Sesion } from './SesionContext'

vi.mock('../../compartido/api/sesionHttp', async (importOriginal) => {
  const original = await importOriginal<typeof import('../../compartido/api/sesionHttp')>()
  return { ...original, configurarSesionHttp: vi.fn(original.configurarSesionHttp) }
})

const USUARIO = { id: 7, nombre: 'ana.perez' }
const RESPUESTA_SESION = { usuario: USUARIO, token: 'token-nuevo', vence_en: '2026-11-10T12:00:00Z' }

type Rutas = Record<string, () => Response | Promise<Response>>

function json(cuerpo: unknown, status: number) {
  return new Response(JSON.stringify(cuerpo), { status })
}

function simularFetch(rutas: Rutas) {
  const fetchMock = vi.fn(async (url: string, init?: RequestInit) => {
    const ruta = rutas[`${init?.method ?? 'GET'} ${url}`]
    if (!ruta) throw new Error(`petición no simulada: ${init?.method} ${url}`)
    return ruta()
  })
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

let sesion: Sesion

function Consumidor() {
  const actual = useSesion()
  useEffect(() => {
    sesion = actual
  })
  return <span data-testid="estado">{actual.estado}</span>
}

function Ubicacion() {
  const { pathname, state } = useLocation()
  return <span data-testid="ubicacion">{JSON.stringify({ pathname, state })}</span>
}

function ubicacion() {
  return JSON.parse(screen.getByTestId('ubicacion').textContent ?? '{}') as {
    pathname: string
    state: unknown
  }
}

function estado() {
  return screen.getByTestId('estado').textContent
}

function montar(entrada: string | { pathname: string; state: unknown } = '/') {
  return render(
    <MemoryRouter initialEntries={[entrada]}>
      <SesionProvider>
        <Consumidor />
        <Ubicacion />
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
  vi.mocked(configurarSesionHttp).mockClear()
})

describe('contexto de sesión', () => {
  it('Escenario 1: registrarse guarda el token, fija el usuario y navega a la pantalla de origen', async () => {
    simularFetch({ 'POST /auth/registro': () => json(RESPUESTA_SESION, 201) })
    montar({ pathname: '/registro', state: { desde: '/backlog' } })

    await act(() => sesion.registrar('ana.perez', 'una-clave-larga'))

    expect(leerToken()).toBe('token-nuevo')
    expect(sesion.estado).toBe('autenticado')
    expect(sesion.usuario).toEqual(USUARIO)
    expect(ubicacion().pathname).toBe('/backlog')
  })

  it('Escenario 1: registrarse sin pantalla de origen navega a /', async () => {
    simularFetch({ 'POST /auth/registro': () => json(RESPUESTA_SESION, 201) })
    montar('/registro')

    await act(() => sesion.registrar('ana.perez', 'una-clave-larga'))

    expect(ubicacion().pathname).toBe('/')
  })

  it('Escenario 4: ingresar guarda el token y deja la sesión autenticada', async () => {
    simularFetch({ 'POST /auth/ingreso': () => json(RESPUESTA_SESION, 200) })
    montar('/ingresar')

    await act(() => sesion.ingresar('ana.perez', 'una-clave-larga'))

    expect(leerToken()).toBe('token-nuevo')
    expect(sesion.estado).toBe('autenticado')
    expect(sesion.usuario).toEqual(USUARIO)
    expect(ubicacion().pathname).toBe('/')
  })

  it('Escenarios 2, 3 y 5: ante un error del servidor la promesa se rechaza con el ApiError y no se guarda token', async () => {
    simularFetch({
      'POST /auth/registro': () =>
        json({ campo: 'nombre', mensaje: 'Ese nombre de usuario ya existe' }, 400),
      'POST /auth/ingreso': () => json({ mensaje: 'Usuario o contraseña incorrectos' }, 401),
    })
    montar('/ingresar')

    await expect(act(() => sesion.registrar('ana', 'x'))).rejects.toEqual(
      new ApiError(400, 'Ese nombre de usuario ya existe', 'nombre'),
    )
    await expect(act(() => sesion.ingresar('ana', 'x'))).rejects.toEqual(
      new ApiError(401, 'Usuario o contraseña incorrectos'),
    )

    expect(leerToken()).toBeNull()
    expect(sesion.estado).toBe('anonimo')
    expect(sesion.usuario).toBeNull()
  })

  it('arranque con token válido (GET /auth/yo 200) deja la sesión autenticada', async () => {
    guardarToken('token-guardado')
    const fetchMock = simularFetch({ 'GET /auth/yo': () => json(USUARIO, 200) })
    montar('/backlog')

    expect(estado()).toBe('cargando')
    await waitFor(() => {
      expect(sesion.estado).toBe('autenticado')
      expect(sesion.usuario).toEqual(USUARIO)
    })

    expect(fetchMock).toHaveBeenCalledWith(
      '/auth/yo',
      expect.objectContaining({
        headers: expect.objectContaining({ Authorization: 'Bearer token-guardado' }),
      }),
    )
  })

  it('arranque sin token deja la sesión anónima sin consultar al servidor', () => {
    const fetchMock = simularFetch({})
    montar('/ingresar')

    expect(estado()).toBe('anonimo')
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it('Escenario 6: arranque con token que da 401 descarta el token, queda anónimo y navega a /ingresar con state.vencida', async () => {
    guardarToken('token-vencido')
    simularFetch({ 'GET /auth/yo': () => json({ mensaje: 'Tu sesión no es válida' }, 401) })
    montar('/backlog')

    await waitFor(() => expect(ubicacion().pathname).toBe('/ingresar'))

    expect(leerToken()).toBeNull()
    expect(estado()).toBe('anonimo')
    expect(ubicacion().state).toMatchObject({ vencida: true })
  })

  it('Escenario 8: arranque con token y GET /auth/yo caído (500) queda sin_conexion y conserva el token', async () => {
    guardarToken('token-guardado')
    simularFetch({ 'GET /auth/yo': () => json({ mensaje: 'Error interno' }, 500) })
    montar('/backlog')

    await waitFor(() => expect(estado()).toBe('sin_conexion'))

    expect(leerToken()).toBe('token-guardado')
    expect(sesion.usuario).toBeNull()
    expect(ubicacion().pathname).toBe('/backlog')
  })

  it('Escenario 8: arranque con token y la red caída queda sin_conexion y conserva el token', async () => {
    guardarToken('token-guardado')
    vi.stubGlobal('fetch', vi.fn(async () => { throw new TypeError('Failed to fetch') }))
    montar('/backlog')

    await waitFor(() => expect(estado()).toBe('sin_conexion'))

    expect(leerToken()).toBe('token-guardado')
  })

  it('Escenario 8: reintentar con el servidor de vuelta pasa a autenticado', async () => {
    guardarToken('token-guardado')
    const respuestas = [() => json({ mensaje: 'Error interno' }, 500), () => json(USUARIO, 200)]
    simularFetch({ 'GET /auth/yo': () => respuestas.shift()!() })
    montar('/backlog')
    await waitFor(() => expect(estado()).toBe('sin_conexion'))

    await act(() => sesion.reintentar())

    expect(estado()).toBe('autenticado')
    expect(sesion.usuario).toEqual(USUARIO)
    expect(leerToken()).toBe('token-guardado')
  })

  it('Escenario 8: reintentar que recibe 401 descarta el token y navega a /ingresar con state.vencida', async () => {
    guardarToken('token-guardado')
    const respuestas = [
      () => json({ mensaje: 'Error interno' }, 500),
      () => json({ mensaje: 'Tu sesión no es válida' }, 401),
    ]
    simularFetch({ 'GET /auth/yo': () => respuestas.shift()!() })
    montar('/backlog')
    await waitFor(() => expect(estado()).toBe('sin_conexion'))

    await act(() => sesion.reintentar())

    expect(leerToken()).toBeNull()
    expect(estado()).toBe('anonimo')
    expect(ubicacion().pathname).toBe('/ingresar')
    expect(ubicacion().state).toMatchObject({ vencida: true })
  })

  it('Escenario 7: salir pide POST /auth/salir, borra el token, queda anónimo y navega a /ingresar', async () => {
    guardarToken('token-guardado')
    const fetchMock = simularFetch({
      'GET /auth/yo': () => json(USUARIO, 200),
      'POST /auth/salir': () => new Response(null, { status: 204 }),
    })
    montar('/backlog')
    await waitFor(() => expect(estado()).toBe('autenticado'))

    await act(() => sesion.salir())

    expect(fetchMock).toHaveBeenCalledWith(
      '/auth/salir',
      expect.objectContaining({ method: 'POST' }),
    )
    expect(leerToken()).toBeNull()
    expect(sesion.estado).toBe('anonimo')
    expect(ubicacion().pathname).toBe('/ingresar')
  })

  it('Escenario 7: salir deja la sesión anónima incluso si el pedido falla', async () => {
    guardarToken('token-guardado')
    simularFetch({
      'GET /auth/yo': () => json(USUARIO, 200),
      'POST /auth/salir': () => json({ mensaje: 'error' }, 500),
    })
    montar('/backlog')
    await waitFor(() => expect(estado()).toBe('autenticado'))

    await act(() => sesion.salir())

    expect(leerToken()).toBeNull()
    expect(sesion.estado).toBe('anonimo')
    expect(ubicacion().pathname).toBe('/ingresar')
  })

  it('Escenario 6: con la sesión iniciada, una petición que recibe 401 descarta el token y navega a /ingresar con desde y vencida', async () => {
    guardarToken('token-guardado')
    simularFetch({
      'GET /auth/yo': () => json(USUARIO, 200),
      'GET /projects/1': () => json({ mensaje: 'Tu sesión no es válida' }, 401),
    })
    montar('/backlog')
    await waitFor(() => expect(estado()).toBe('autenticado'))

    await act(() => get('/projects/1').catch(() => undefined))

    expect(leerToken()).toBeNull()
    expect(sesion.estado).toBe('anonimo')
    expect(ubicacion()).toEqual({
      pathname: '/ingresar',
      state: { desde: '/backlog', vencida: true },
    })
  })

  it('registra el enganche del cliente al montar y lo quita al desmontar', () => {
    simularFetch({})
    const { unmount } = montar('/ingresar')

    expect(vi.mocked(configurarSesionHttp).mock.calls.at(-1)?.[0]).toEqual({
      obtenerToken: expect.any(Function),
      alNoAutenticado: expect.any(Function),
    })

    unmount()

    expect(vi.mocked(configurarSesionHttp).mock.calls.at(-1)?.[0]).toBeNull()
  })
})
