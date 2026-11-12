import { cleanup, render, screen, waitFor, within } from '@testing-library/react'
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

type Rutas = { estado?: () => Response; datos?: () => Response }

// Enruta por URL: /demo/estado sin configurar responde 404 (función apagada).
function simularFetch(respuesta: () => Response, rutas: Rutas = {}) {
  const fetchMock = vi.fn(async (ruta: string, _opciones?: RequestInit) => {
    if (ruta === '/demo/estado') return rutas.estado ? rutas.estado() : new Response('404 page not found', { status: 404 })
    if (ruta === '/demo/datos') return rutas.datos ? rutas.datos() : new Response('404 page not found', { status: 404 })
    return respuesta()
  })
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

  it('Given la contraseña oculta, When presiono «Mostrar contraseña», Then se ve como texto y el botón queda presionado sin cambiar de nombre', async () => {
    montar()
    const usuario = userEvent.setup()
    const campo = screen.getByLabelText('Contraseña')
    const boton = screen.queryByRole('button', { name: 'Mostrar contraseña' })
    expect(boton).not.toBeNull()
    expect(boton!.getAttribute('aria-pressed')).toBe('false')

    await usuario.click(boton!)

    expect(campo.getAttribute('type')).toBe('text')
    expect(screen.getByRole('button', { name: 'Mostrar contraseña' }).getAttribute('aria-pressed')).toBe('true')
  })

  it('Given la contraseña visible, When presiono de nuevo «Mostrar contraseña», Then vuelve a ocultarse', async () => {
    montar()
    const usuario = userEvent.setup()
    const boton = screen.queryByRole('button', { name: 'Mostrar contraseña' })
    expect(boton).not.toBeNull()
    await usuario.click(boton!)
    expect(screen.getByLabelText('Contraseña').getAttribute('type')).toBe('text')

    await usuario.click(boton!)

    expect(screen.getByLabelText('Contraseña').getAttribute('type')).toBe('password')
    expect(boton!.getAttribute('aria-pressed')).toBe('false')
  })

  it('Given la pantalla de ingreso, When se muestra, Then el panel de presentación lleva la marca «Métrica» y el formulario el sobretítulo «Bienvenido»', () => {
    montar()

    expect(screen.queryByRole('link', { name: 'Métrica' })?.getAttribute('href')).toBe('/')
    expect(screen.queryByText('Bienvenido')).not.toBeNull()
  })

  it('tiene un enlace a la pantalla de registro', async () => {
    montar()
    const usuario = userEvent.setup()

    await usuario.click(screen.getByRole('link', { name: 'Registrate' }))

    expect(screen.getByRole('heading', { name: 'Crear cuenta' })).toBeTruthy()
  })

  it('Dado que /demo/estado responde 200, cuando se muestra /ingresar, entonces aparece el botón «Cargar datos de prueba»', async () => {
    simularFetch(() => json({}, 401), { estado: () => json({ habilitado: true }, 200) })
    montar()

    expect(await screen.findByRole('button', { name: 'Cargar datos de prueba' })).toBeTruthy()
  })

  it('Dado que /demo/estado responde 200, cuando se muestra /ingresar, entonces la tarjeta «Datos de prueba» explica qué se carga', async () => {
    simularFetch(() => json({}, 401), { estado: () => json({ habilitado: true }, 200) })
    montar()

    await screen.findByRole('button', { name: 'Cargar datos de prueba' })

    expect(screen.queryByRole('heading', { name: 'Datos de prueba' })).not.toBeNull()
    expect(screen.queryByText(/Crea cuatro usuarios y el proyecto #25/)).not.toBeNull()
  })

  it('Dado que /demo/estado responde 404, cuando se muestra /ingresar, entonces no aparece el botón ni la tabla', async () => {
    const fetchMock = simularFetch(() => json({}, 401))
    montar()

    await waitFor(() => expect(fetchMock).toHaveBeenCalledWith('/demo/estado', expect.anything()))
    await Promise.resolve()
    expect(screen.queryByRole('button', { name: 'Cargar datos de prueba' })).toBeNull()
    expect(screen.queryByRole('heading', { name: 'Datos de prueba' })).toBeNull()
    expect(screen.queryByRole('table')).toBeNull()
  })

  it('Dado que /demo/estado falla por conexión, cuando se muestra /ingresar, entonces no aparece el botón', async () => {
    const fetchMock = vi.fn(async () => {
      throw new TypeError('Failed to fetch')
    })
    vi.stubGlobal('fetch', fetchMock)
    montar()

    await waitFor(() => expect(fetchMock).toHaveBeenCalled())
    await Promise.resolve()
    expect(screen.queryByRole('button', { name: 'Cargar datos de prueba' })).toBeNull()
  })

  it('Dado el botón visible, cuando la carga responde 201, entonces /ingresar muestra la tabla de usuarios de la demo', async () => {
    simularFetch(() => json({}, 401), {
      estado: () => json({ habilitado: true }, 200),
      datos: () => json({ proyecto_id: 25 }, 201),
    })
    montar()
    const usuario = userEvent.setup()

    await usuario.click(await screen.findByRole('button', { name: 'Cargar datos de prueba' }))

    const tabla = await screen.findByRole('table', { name: 'Usuarios de la demostración' })
    expect(within(tabla).getByRole('cell', { name: 'carlos.moya' })).toBeTruthy()
  })

  it('Dado que la carga responde 409 por usuario existente, cuando la pido desde /ingresar, entonces el aviso nombra al usuario', async () => {
    simularFetch(() => json({}, 401), {
      estado: () => json({ habilitado: true }, 200),
      datos: () => json({ campo: 'usuario_existente', mensaje: 'ya existe el usuario ana.iriarte' }, 409),
    })
    montar()
    const usuario = userEvent.setup()

    await usuario.click(await screen.findByRole('button', { name: 'Cargar datos de prueba' }))

    expect((await screen.findByRole('status')).textContent).toContain('ya existe el usuario ana.iriarte')
    expect(screen.queryByRole('alert')).toBeNull()
  })

  it('Dado la tabla de la demo, cuando uso a ana.iriarte, entonces el formulario se rellena y al ingresar viajan esas credenciales', async () => {
    const fetchMock = simularFetch(() => json(RESPUESTA_SESION, 200), {
      estado: () => json({ habilitado: true }, 200),
      datos: () => json({ proyecto_id: 25 }, 201),
    })
    montar()
    const usuario = userEvent.setup()

    await usuario.click(await screen.findByRole('button', { name: 'Cargar datos de prueba' }))
    await usuario.click(await screen.findByRole('button', { name: 'Usar ana.iriarte' }))

    expect((screen.getByLabelText('Nombre de usuario') as HTMLInputElement).value).toBe('ana.iriarte')
    expect((screen.getByLabelText('Contraseña') as HTMLInputElement).value).toBe('demo1234')
    await usuario.click(screen.getByRole('button', { name: 'Ingresar' }))
    await waitFor(() => expect(fetchMock.mock.calls.some(([ruta]) => ruta === '/auth/ingreso')).toBe(true))
    const [, opciones] = fetchMock.mock.calls.find(([ruta]) => ruta === '/auth/ingreso')!
    expect(JSON.parse(String(opciones?.body))).toEqual({ nombre: 'ana.iriarte', contrasena: 'demo1234' })
  })
})
