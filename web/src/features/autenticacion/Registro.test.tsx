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
  await usuario.type(screen.getByLabelText('Repetir contraseña'), contrasena)
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
    expect(descripcion(screen.getByLabelText('Nombre de usuario'))).not.toBe('La contraseña debe tener al menos 8 caracteres')
  })

  it('Given la pantalla de registro, When se muestra, Then cada campo describe la regla que aplica el servidor', () => {
    montar()

    expect(descripcion(screen.getByLabelText('Nombre de usuario'))).toBe(
      'De 3 a 30 caracteres: letras sin tildes, números, punto, guion y guion bajo.',
    )
    expect(descripcion(screen.getByLabelText('Contraseña'))).toBe('De 8 a 128 caracteres.')
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

  it('Given la pantalla de registro, When se muestra, Then el panel de presentación lleva la marca «Métrica» y el formulario el sobretítulo «Nueva cuenta»', () => {
    montar()

    expect(screen.queryByRole('link', { name: 'Métrica' })?.getAttribute('href')).toBe('/')
    expect(screen.queryByText('Nueva cuenta')).not.toBeNull()
  })

  it('tiene un enlace a la pantalla de ingreso', async () => {
    montar()
    const usuario = userEvent.setup()

    await usuario.click(screen.getByRole('link', { name: 'Ingresá' }))

    expect(screen.getByRole('heading', { name: 'Ingresar' })).toBeTruthy()
  })
})

describe('pantalla de registro: campos vacíos', () => {
  it('Given los campos vacíos, When envío, Then cada campo muestra el mensaje de su regla y no se llama al servidor', async () => {
    const fetchMock = simularFetch(() => json(RESPUESTA_SESION, 201))
    montar()

    await userEvent.setup().click(screen.getByRole('button', { name: 'Crear cuenta' }))

    expect(descripcion(screen.getByLabelText('Nombre de usuario'))).toBe(
      'el nombre de usuario debe tener entre 3 y 30 caracteres',
    )
    expect(descripcion(screen.getByLabelText('Contraseña'))).toBe('la contraseña debe tener entre 8 y 128 caracteres')
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it('Given un nombre de solo espacios, When envío, Then cuenta como vacío', async () => {
    const fetchMock = simularFetch(() => json(RESPUESTA_SESION, 201))
    montar()
    const usuario = userEvent.setup()

    await usuario.type(screen.getByLabelText('Nombre de usuario'), '   ')
    await usuario.type(screen.getByLabelText('Contraseña'), 'una-clave-larga')
    await usuario.click(screen.getByRole('button', { name: 'Crear cuenta' }))

    expect(descripcion(screen.getByLabelText('Nombre de usuario'))).toBe(
      'el nombre de usuario debe tener entre 3 y 30 caracteres',
    )
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it('Given solo la contraseña vacía, When envío, Then el mensaje aparece solo junto a la contraseña', async () => {
    const fetchMock = simularFetch(() => json(RESPUESTA_SESION, 201))
    montar()
    const usuario = userEvent.setup()

    await usuario.type(screen.getByLabelText('Nombre de usuario'), 'ana.perez')
    await usuario.click(screen.getByRole('button', { name: 'Crear cuenta' }))

    expect(descripcion(screen.getByLabelText('Contraseña'))).toBe('la contraseña debe tener entre 8 y 128 caracteres')
    expect(descripcion(screen.getByLabelText('Nombre de usuario'))).not.toBe(
      'el nombre de usuario debe tener entre 3 y 30 caracteres',
    )
    expect(fetchMock).not.toHaveBeenCalled()
  })
})

describe('pantalla de registro: repetir contraseña', () => {
  async function enviarCon(nombre: string, contrasena: string, repeticion: string) {
    const usuario = userEvent.setup()
    await usuario.type(screen.getByLabelText('Nombre de usuario'), nombre)
    await usuario.type(screen.getByLabelText('Contraseña'), contrasena)
    if (repeticion) await usuario.type(screen.getByLabelText('Repetir contraseña'), repeticion)
    await usuario.click(screen.getByRole('button', { name: 'Crear cuenta' }))
  }

  it('Given la pantalla de registro, When se muestra, Then hay un campo «Repetir contraseña» oculto y sin texto de ayuda', () => {
    montar()

    const campo = screen.getByLabelText('Repetir contraseña')
    expect(campo.getAttribute('type')).toBe('password')
    expect(campo.getAttribute('aria-describedby')).toBeNull()
  })

  it('Given la pantalla de registro, When se muestra, Then no lleva ninguna etiqueta de propuesta del prototipo', () => {
    montar()

    expect(screen.queryByText(/Propuesta/)).toBeNull()
    expect(screen.queryByText(/La app real todavía/)).toBeNull()
  })

  it('Given la repetición oculta, When presiono «Mostrar la repetición de la contraseña», Then solo ese campo se ve como texto', async () => {
    montar()
    const usuario = userEvent.setup()
    const boton = screen.getByRole('button', { name: 'Mostrar la repetición de la contraseña' })
    expect(boton.getAttribute('aria-pressed')).toBe('false')

    await usuario.click(boton)

    expect(screen.getByLabelText('Repetir contraseña').getAttribute('type')).toBe('text')
    expect(screen.getByLabelText('Contraseña').getAttribute('type')).toBe('password')
    expect(boton.getAttribute('aria-pressed')).toBe('true')
  })

  it('Given los campos vacíos, When envío, Then también pide repetir la contraseña', async () => {
    const fetchMock = simularFetch(() => json(RESPUESTA_SESION, 201))
    montar()

    await userEvent.setup().click(screen.getByRole('button', { name: 'Crear cuenta' }))

    expect(descripcion(screen.getByLabelText('Repetir contraseña'))).toBe('Repetí la contraseña.')
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it('Given la contraseña escrita y la repetición vacía, When envío, Then pide repetirla y no se llama al servidor', async () => {
    const fetchMock = simularFetch(() => json(RESPUESTA_SESION, 201))
    montar()

    await enviarCon('ana.perez', 'una-clave-larga', '')

    expect(descripcion(screen.getByLabelText('Repetir contraseña'))).toBe('Repetí la contraseña.')
    expect(descripcion(screen.getByLabelText('Contraseña'))).not.toBe('Repetí la contraseña.')
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it('Given una repetición distinta, When envío, Then avisa que no coinciden y no se llama al servidor', async () => {
    const fetchMock = simularFetch(() => json(RESPUESTA_SESION, 201))
    montar()

    await enviarCon('ana.perez', 'una-clave-larga', 'una-clave-larga-distinta')

    expect(descripcion(screen.getByLabelText('Repetir contraseña'))).toBe('Las contraseñas no coinciden.')
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it('Given una repetición igual, When envío, Then viajan solo el nombre y la contraseña', async () => {
    const fetchMock = vi.fn(async (_ruta: string, _opciones?: RequestInit) => json(RESPUESTA_SESION, 201))
    vi.stubGlobal('fetch', fetchMock)
    montar()

    await enviarCon('ana.perez', 'una-clave-larga', 'una-clave-larga')

    await waitFor(() => expect(fetchMock).toHaveBeenCalledTimes(1))
    expect(JSON.parse(String(fetchMock.mock.calls[0][1]?.body))).toEqual({ nombre: 'ana.perez', contrasena: 'una-clave-larga' })
  })

  it('Given «no coinciden», When corrijo la repetición y el servidor rechaza el nombre, Then solo queda el error del servidor', async () => {
    simularFetch(() => json({ campo: 'nombre', mensaje: 'Ese nombre de usuario ya existe' }, 409))
    montar()
    const usuario = userEvent.setup()
    await enviarCon('ana.perez', 'una-clave-larga', 'otra')

    await usuario.clear(screen.getByLabelText('Repetir contraseña'))
    await usuario.type(screen.getByLabelText('Repetir contraseña'), 'una-clave-larga')
    await usuario.click(screen.getByRole('button', { name: 'Crear cuenta' }))

    await waitFor(() => expect(screen.getAllByRole('alert').map((aviso) => aviso.textContent)).toEqual(['Ese nombre de usuario ya existe']))
  })
})

describe('pantalla de registro: foco y título', () => {
  afterEach(() => {
    document.title = ''
  })

  it('Given la pantalla de registro, When se muestra, Then el foco queda en «Nombre de usuario»', () => {
    montar()

    expect(document.activeElement).toBe(screen.getByLabelText('Nombre de usuario'))
  })

  it('Given la pantalla de registro, When se muestra, Then la pestaña dice «Métrica · Crear cuenta» y al salir vuelve el título anterior', () => {
    document.title = 'Métrica · Seguimiento de proyectos'

    const { unmount } = montar()

    expect(document.title).toBe('Métrica · Crear cuenta')
    unmount()
    expect(document.title).toBe('Métrica · Seguimiento de proyectos')
  })
})
