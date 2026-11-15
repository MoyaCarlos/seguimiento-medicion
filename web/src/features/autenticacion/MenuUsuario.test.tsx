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

function montar(rol?: string) {
  return render(
    <MemoryRouter initialEntries={['/']}>
      <SesionProvider>
        <MenuUsuario rol={rol} />
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
    await usuario.tab()
    await usuario.tab()
    expect(document.activeElement).toBe(screen.getByRole('button', { name: 'Cerrar sesión' }))

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

describe('MenuUsuario con el rol del proyecto', () => {
  async function abrirConRol(rol?: string) {
    guardarToken('token-guardado')
    vi.stubGlobal('fetch', vi.fn(async () => json(USUARIO, 200)))
    montar(rol)
    await userEvent.click(await screen.findByRole('button', { name: 'Menú de usuario' }))
    return screen.getByRole('group', { name: 'Cuenta' })
  }

  it('Given una pantalla de proyecto, When se abre el menú, Then dice «Tu rol en este proyecto: Scrum Master» y el rol reemplaza a «Sesión iniciada» bajo el nombre', async () => {
    const panel = await abrirConRol('Scrum Master')

    expect(within(panel).getByText('Tu rol en este proyecto: Scrum Master')).toBeTruthy()
    expect(within(panel).queryByText('Sesión iniciada')).toBeNull()
    expect(panel.querySelector('.menu-usuario__cabecera small')?.textContent).toBe('Scrum Master')
  })

  it('Given una pantalla fuera de un proyecto, When se abre el menú, Then no hay línea de rol y sigue «Sesión iniciada»', async () => {
    const panel = await abrirConRol(undefined)

    expect(within(panel).queryByText(/Tu rol/)).toBeNull()
    expect(within(panel).getByText('Sesión iniciada')).toBeTruthy()
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

  const botonDelMenu = () => screen.getByRole('button', { name: 'Menú de usuario' })

  async function abrirFormulario() {
    await userEvent.click(await screen.findByRole('button', { name: 'Menú de usuario' }))
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

    expect(await within(botonDelMenu()).findByText('ana.gomez')).toBeTruthy()
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
    expect(within(botonDelMenu()).getByText(USUARIO.nombre)).toBeTruthy()
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
    expect(within(botonDelMenu()).getByText(USUARIO.nombre)).toBeTruthy()
  })
})

describe('MenuUsuario - cambiar nombre dentro del menú rediseñado', () => {
  async function conMenuAbierto() {
    guardarToken('token-guardado')
    vi.stubGlobal('fetch', vi.fn(async () => json(USUARIO, 200)))
    montar()
    await userEvent.click(await screen.findByRole('button', { name: 'Menú de usuario' }))
    return screen.getByRole('group', { name: 'Cuenta' })
  }

  it('Given el menú abierto, When aún no pido el cambio, Then «Cambiar nombre» es un ítem del panel, con ícono, antes de «Cerrar sesión»', async () => {
    const panel = await conMenuAbierto()

    const cambiar = within(panel).getByRole('button', { name: 'Cambiar nombre' })
    const salir = within(panel).getByRole('button', { name: 'Cerrar sesión' })

    expect(cambiar.className).toBe('menu-usuario__item')
    expect(cambiar.querySelector('svg')).toBeTruthy()
    expect(cambiar.compareDocumentPosition(salir) & Node.DOCUMENT_POSITION_FOLLOWING).toBeTruthy()
    expect(within(panel).queryByLabelText('Nuevo nombre de usuario')).toBeNull()
  })

  it('Given el formulario de cambio abierto, When cierro el menú con Escape y lo vuelvo a abrir, Then el formulario no queda abierto', async () => {
    const panel = await conMenuAbierto()
    await userEvent.click(within(panel).getByRole('button', { name: 'Cambiar nombre' }))
    expect(screen.getByLabelText('Nuevo nombre de usuario')).toBeTruthy()

    await userEvent.keyboard('{Escape}')
    await userEvent.click(screen.getByRole('button', { name: 'Menú de usuario' }))

    expect(screen.queryByLabelText('Nuevo nombre de usuario')).toBeNull()
    expect(screen.getByRole('button', { name: 'Cambiar nombre' })).toBeTruthy()
  })

  it('Given un cambio guardado, When termina, Then el menú se cierra y el foco vuelve al botón del menú', async () => {
    guardarToken('token-guardado')
    vi.stubGlobal(
      'fetch',
      vi.fn(async (url: string, init?: RequestInit) =>
        url === '/auth/yo' && init?.method === 'PUT' ? json({ id: 7, nombre: 'ana.gomez' }, 200) : json(USUARIO, 200),
      ),
    )
    montar()
    await userEvent.click(await screen.findByRole('button', { name: 'Menú de usuario' }))
    await userEvent.click(screen.getByRole('button', { name: 'Cambiar nombre' }))
    const campo = screen.getByLabelText('Nuevo nombre de usuario')
    await userEvent.clear(campo)
    await userEvent.type(campo, 'ana.gomez')

    await userEvent.click(screen.getByRole('button', { name: 'Guardar' }))

    await within(screen.getByRole('button', { name: 'Menú de usuario' })).findByText('ana.gomez')
    expect(screen.queryByRole('group', { name: 'Cuenta' })).toBeNull()
    expect(document.activeElement).toBe(screen.getByRole('button', { name: 'Menú de usuario' }))
  })
})

