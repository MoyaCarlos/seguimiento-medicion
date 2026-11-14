import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Route, Routes } from 'react-router'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { Configuracion } from './Configuracion'

const integrantes = [
  {
    id: 1,
    proyecto_id: 3,
    nombre: 'Ana',
    rol: 'scrum_master',
    vinculado_en: '2026-09-01T00:00:00Z',
  },
]

function abrir() {
  render(
    <MemoryRouter initialEntries={['/projects/3/proyecto/settings']}>
      <Routes>
        <Route path="/projects/:id/proyecto/settings" element={<Configuracion />} />
      </Routes>
    </MemoryRouter>,
  )
}

const registrados = [{ id: 1, nombre: 'Ana' }, { id: 2, nombre: 'Jimena' }]

function responderListaYAlta() {
  const fetchMock = vi.fn((url: string, init?: RequestInit) => {
    if (init?.method === 'POST') {
      return Promise.resolve(
        new Response(
          JSON.stringify({
            id: 2,
            proyecto_id: 3,
            nombre: 'Jimena',
            rol: 'product_builder',
            vinculado_en: '2026-09-02T00:00:00Z',
          }),
          { status: 201 },
        ),
      )
    }
    if (url === '/users') {
      return Promise.resolve(new Response(JSON.stringify(registrados), { status: 200 }))
    }
    return Promise.resolve(new Response(JSON.stringify(integrantes), { status: 200 }))
  })
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('pantalla de configuración del equipo', () => {
  it('lista los integrantes con su rol', async () => {
    responderListaYAlta()
    abrir()

    expect(await screen.findByText('Ana')).toBeTruthy()
    expect(screen.getByText('Ana').closest('li')?.textContent).toContain('Scrum Master')
  })

  it('valida que se elija un usuario antes de asignar', async () => {
    const fetchMock = responderListaYAlta()
    abrir()

    await screen.findByText('Ana')
    await screen.findByLabelText('Usuario')
    await userEvent.click(screen.getByRole('button', { name: 'Asignar integrante' }))

    expect(screen.getByText('Elegí un usuario de la lista')).toBeTruthy()
    expect(fetchMock.mock.calls.filter(([, init]) => init?.method === 'POST')).toHaveLength(0)
  })
})

describe('configuración del equipo - elegir de la lista de usuarios (HU-21)', () => {
  const usuarios = [
    { id: 1, nombre: 'Ana' },
    { id: 2, nombre: 'Bruno' },
    { id: 3, nombre: 'Carla' },
  ]

  function responderUsuarios(alta: Response) {
    const fetchMock = vi.fn((url: string, init?: RequestInit) => {
      if (init?.method === 'POST') return Promise.resolve(alta)
      if (url === '/users') return Promise.resolve(new Response(JSON.stringify(usuarios)))
      if (url === '/projects/3') {
        return Promise.resolve(new Response(JSON.stringify({ id: 3, rol: 'scrum_master' })))
      }
      return Promise.resolve(new Response(JSON.stringify(integrantes)))
    })
    vi.stubGlobal('fetch', fetchMock)
    return fetchMock
  }

  it('Escenario 8: elijo de la lista un usuario registrado y entra como Product Builder', async () => {
    const fetchMock = responderUsuarios(
      new Response(
        JSON.stringify({
          id: 2,
          proyecto_id: 3,
          nombre: 'Bruno',
          rol: 'product_builder',
          vinculado_en: '2026-09-02T00:00:00Z',
        }),
        { status: 201 },
      ),
    )
    abrir()

    await screen.findByText('Ana')
    const selector = await screen.findByLabelText('Usuario')
    expect(fetchMock).toHaveBeenCalledWith('/users', expect.anything())
    expect(screen.queryByLabelText('Rol')).toBeNull()

    // Ana ya integra el proyecto: no se ofrece
    const opciones = Array.from(selector.querySelectorAll('option')).map((o) => o.textContent)
    expect(opciones).toContain('Bruno')
    expect(opciones).toContain('Carla')
    expect(opciones).not.toContain('Ana')

    await userEvent.selectOptions(selector, 'Bruno')
    await userEvent.click(screen.getByRole('button', { name: 'Asignar integrante' }))

    await waitFor(() =>
      expect(fetchMock.mock.calls.some(([, init]) => init?.method === 'POST')).toBe(true),
    )
    const llamada = fetchMock.mock.calls.find(([, init]) => init?.method === 'POST')
    expect(llamada?.[0]).toBe('/projects/3/members')
    expect(JSON.parse(String(llamada?.[1]?.body))).toEqual({ nombre: 'Bruno' })
    expect(await screen.findByText('Product Builder')).toBeTruthy()
  })

  it('muestra el error del servidor si el usuario no existe', async () => {
    responderUsuarios(
      new Response(JSON.stringify({ mensaje: 'usuario no encontrado' }), { status: 404 }),
    )
    abrir()

    await userEvent.selectOptions(await screen.findByLabelText('Usuario'), 'Bruno')
    await userEvent.click(screen.getByRole('button', { name: 'Asignar integrante' }))

    expect(await screen.findByText('usuario no encontrado')).toBeTruthy()
  })

  it('avisa cuando no hay usuarios registrados para agregar', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn((url: string) =>
        Promise.resolve(
          new Response(JSON.stringify(url === '/users' ? [] : integrantes)),
        ),
      ),
    )
    abrir()

    expect(await screen.findByText('No hay usuarios registrados para agregar')).toBeTruthy()
  })
})

describe('configuración del equipo - rol (HU-21)', () => {
  const nota = (boton: HTMLElement) =>
    document.getElementById(boton.getAttribute('aria-describedby') ?? '')?.textContent

  function responderConRol(rol: 'scrum_master' | 'product_builder') {
    const fetchMock = vi.fn((url: string) =>
      Promise.resolve(
        new Response(JSON.stringify(url === '/projects/3' ? { id: 3, rol } : integrantes), { status: 200 }),
      ),
    )
    vi.stubGlobal('fetch', fetchMock)
    return fetchMock
  }

  it('Escenario 2: como Product Builder veo «Asignar integrante» deshabilitado con la nota «Solo el Scrum Master puede agregar integrantes»', async () => {
    responderConRol('product_builder')
    abrir()

    await screen.findByText('Ana')
    const boton = screen.getByRole('button', { name: 'Asignar integrante' })
    await waitFor(() => expect(boton.hasAttribute('disabled')).toBe(true))
    expect(nota(boton)).toBe('Solo el Scrum Master puede agregar integrantes')
  })

  it('Escenario 1: como Scrum Master «Asignar integrante» está habilitado', async () => {
    const fetchMock = responderConRol('scrum_master')
    abrir()

    await screen.findByText('Ana')
    await waitFor(() => expect(fetchMock).toHaveBeenCalledWith('/projects/3', expect.anything()))
    const boton = screen.getByRole('button', { name: 'Asignar integrante' })
    expect(boton.hasAttribute('disabled')).toBe(false)
    expect(boton.hasAttribute('aria-describedby')).toBe(false)
  })
})

describe('quitar integrante desde Configurar equipo (HU-16, #182)', () => {
  const equipo = [
    {
      id: 1,
      proyecto_id: 3,
      nombre: 'Ana',
      rol: 'scrum_master',
      vinculado_en: '2026-09-01T00:00:00Z',
    },
    {
      id: 2,
      proyecto_id: 3,
      nombre: 'Jimena',
      rol: 'product_builder',
      vinculado_en: '2026-09-02T00:00:00Z',
    },
  ]
  const usuarios = [
    { id: 1, nombre: 'Ana' },
    { id: 2, nombre: 'Jimena' },
  ]

  function responderEquipo(quitar?: (id: number) => Response, fechaCierre: string | null = null) {
    const fetchMock = vi.fn((url: string, init?: RequestInit) => {
      if (url === '/projects/3' && (!init?.method || init.method === 'GET')) {
        return Promise.resolve(
          new Response(
            JSON.stringify({ id: 3, rol: 'scrum_master', fecha_cierre: fechaCierre }),
          ),
        )
      }
      if (url === '/users') {
        return Promise.resolve(new Response(JSON.stringify(usuarios)))
      }
      if (url.startsWith('/projects/3/members/') && init?.method === 'DELETE') {
        const id = Number(url.split('/').pop())
        return Promise.resolve(quitar?.(id) ?? new Response(null, { status: 204 }))
      }
      if (url === '/projects/3/members' && init?.method === 'POST') {
        return Promise.resolve(
          new Response(
            JSON.stringify({
              id: 2,
              proyecto_id: 3,
              nombre: 'Jimena',
              rol: 'product_builder',
              vinculado_en: '2026-09-02T00:00:00Z',
            }),
            { status: 201 },
          ),
        )
      }
      return Promise.resolve(new Response(JSON.stringify(equipo)))
    })
    vi.stubGlobal('fetch', fetchMock)
    return fetchMock
  }

  const filaDe = (nombre: string) => screen.getByText(nombre).closest('li') as HTMLElement
  const articuloEquipo = () =>
    screen.getByRole('heading', { name: 'Equipo del proyecto' }).closest('article') as HTMLElement

  const pedirQuitar = (nombre: string) =>
    userEvent.click(within(filaDe(nombre)).getByRole('button', { name: 'Quitar' }))

  const confirmarQuitar = (nombre: string) =>
    userEvent.click(within(filaDe(nombre)).getByRole('button', { name: 'Sí, quitar' }))

  it('Given un equipo con Scrum Master y Product Builder, When se muestra, Then solo el Product Builder tiene «Quitar»', async () => {
    responderEquipo()
    abrir()

    await screen.findByText('Ana')
    expect(within(filaDe('Jimena')).getByRole('button', { name: 'Quitar' })).toBeTruthy()
    expect(within(filaDe('Ana')).queryByRole('button', { name: 'Quitar' })).toBeNull()
  })

  it('Given un Product Builder, When pulso «Quitar» y cancelo, Then no se llama a DELETE y la fila sigue', async () => {
    const fetchMock = responderEquipo()
    abrir()
    await screen.findByText('Jimena')

    await pedirQuitar('Jimena')
    expect(within(filaDe('Jimena')).getByText('¿Quitar a Jimena del equipo?')).toBeTruthy()
    expect(fetchMock.mock.calls.filter(([, init]) => init?.method === 'DELETE')).toHaveLength(0)

    await userEvent.click(within(filaDe('Jimena')).getByRole('button', { name: 'Cancelar' }))

    expect(fetchMock.mock.calls.filter(([, init]) => init?.method === 'DELETE')).toHaveLength(0)
    expect(screen.getByText('Jimena')).toBeTruthy()
    expect(within(filaDe('Jimena')).queryByRole('button', { name: 'Sí, quitar' })).toBeNull()
  })

  it('Given un Product Builder, When confirmo «Quitar», Then se llama a DELETE, la fila desaparece y vuelve a estar disponible para el alta', async () => {
    const fetchMock = responderEquipo()
    abrir()
    await screen.findByText('Jimena')

    await pedirQuitar('Jimena')
    await confirmarQuitar('Jimena')

    await waitFor(() => expect(within(articuloEquipo()).queryByText('Jimena')).toBeNull())
    expect(fetchMock).toHaveBeenCalledWith(
      '/projects/3/members/2',
      expect.objectContaining({ method: 'DELETE' }),
    )

    const selector = await screen.findByLabelText('Usuario')
    await userEvent.selectOptions(selector, 'Jimena')
    await userEvent.click(screen.getByRole('button', { name: 'Asignar integrante' }))

    await waitFor(() => expect(within(articuloEquipo()).getByText('Jimena')).toBeTruthy())
    expect(within(articuloEquipo()).getByText('Jimena').closest('li')?.textContent).toContain(
      'Product Builder',
    )
    expect(fetchMock).toHaveBeenCalledWith(
      '/projects/3/members',
      expect.objectContaining({ method: 'POST' }),
    )
  })

  it('Given un proyecto cerrado con un Product Builder, When se muestra el equipo, Then «Quitar» queda deshabilitado con la nota y no envía', async () => {
    const fetchMock = responderEquipo(undefined, '2026-11-10')
    abrir()
    await screen.findByText('Jimena')

    const boton = within(filaDe('Jimena')).getByRole('button', {
      name: 'Quitar',
    }) as HTMLButtonElement
    await waitFor(() => expect(boton.disabled).toBe(true))
    expect(boton.getAttribute('aria-describedby')).toBe('configuracion-nota-cerrado-equipo')
    expect(document.getElementById('configuracion-nota-cerrado-equipo')?.textContent).toBe(
      'El proyecto está cerrado',
    )
    expect(fetchMock.mock.calls.filter(([, init]) => init?.method === 'DELETE')).toHaveLength(0)
  })

  it.each([
    [409, 'No se puede quitar al Scrum Master del proyecto'],
    [404, 'integrante no encontrado'],
    [403, 'Solo el Scrum Master puede realizar esta acción'],
  ])(
    'Given el servidor responde %i, When confirmo «Quitar», Then se muestra «%s» y la fila sigue',
    async (status, mensaje) => {
      responderEquipo(() => new Response(JSON.stringify({ mensaje }), { status }))
      abrir()
      await screen.findByText('Jimena')

      await pedirQuitar('Jimena')
      await confirmarQuitar('Jimena')

      expect(await screen.findByText(mensaje)).toBeTruthy()
      expect(screen.getByText('Jimena')).toBeTruthy()
    },
  )
})

describe('configuración del equipo - proyecto cerrado (cerrar y reabrir, US2)', () => {
  const NOTA_CERRADO = 'El proyecto está cerrado'
  const NOTA_HU21 = 'Solo el Scrum Master puede agregar integrantes'
  const usuarios = [
    { id: 1, nombre: 'Ana' },
    { id: 2, nombre: 'Bruno' },
  ]

  type Proyecto = { status: number; cuerpo: unknown } | 'pendiente'

  const proyectoCon = (rol: string, fechaCierre: string | null): Proyecto => ({
    status: 200,
    cuerpo: { id: 3, rol, fecha_cierre: fechaCierre },
  })

  function responder(proyecto: Proyecto, alta?: Response) {
    const fetchMock = vi.fn((url: string, init?: RequestInit) => {
      if (init?.method === 'POST') {
        return Promise.resolve(alta ?? new Response(JSON.stringify({}), { status: 201 }))
      }
      if (url === '/projects/3') {
        if (proyecto === 'pendiente') return new Promise<Response>(() => {})
        return Promise.resolve(
          new Response(JSON.stringify(proyecto.cuerpo), { status: proyecto.status }),
        )
      }
      if (url === '/users') return Promise.resolve(new Response(JSON.stringify(usuarios)))
      return Promise.resolve(new Response(JSON.stringify(integrantes)))
    })
    vi.stubGlobal('fetch', fetchMock)
    return fetchMock
  }

  function llamadasAlProyecto(fetchMock: ReturnType<typeof responder>) {
    return fetchMock.mock.calls.filter(
      ([url, init]) => url === '/projects/3' && (init?.method ?? 'GET') === 'GET',
    )
  }

  function descripcionDe(elemento: HTMLElement) {
    const id = elemento.getAttribute('aria-describedby')
    return id ? document.getElementById(id)?.textContent : undefined
  }

  it('Escenario: proyecto cerrado: muestra el aviso y la lista de integrantes se sigue viendo', async () => {
    responder(proyectoCon('scrum_master', '2026-11-10'))
    abrir()

    expect(await screen.findByText('Proyecto cerrado el 2026-11-10')).toBeTruthy()
    expect(await screen.findByText('Ana')).toBeTruthy()
    expect(screen.getByText('Ana').closest('li')?.textContent).toContain('Scrum Master')
  })

  it('Escenario: proyecto cerrado: «Asignar integrante» queda deshabilitado con una sola nota «El proyecto está cerrado» y no envía', async () => {
    const fetchMock = responder(proyectoCon('scrum_master', '2026-11-10'))
    abrir()

    await screen.findByText('Proyecto cerrado el 2026-11-10')
    const selector = await screen.findByLabelText('Usuario')
    const boton = screen.getByRole('button', { name: 'Asignar integrante' }) as HTMLButtonElement
    expect(boton.disabled).toBe(true)
    expect(descripcionDe(boton)).toBe(NOTA_CERRADO)

    const formulario = screen.getByRole('form', { name: 'Asignar integrante' })
    expect(within(formulario).getAllByText(NOTA_CERRADO)).toHaveLength(1)

    fireEvent.change(selector, { target: { value: 'Bruno' } })
    fireEvent.submit(formulario)
    expect(fetchMock.mock.calls.some(([, init]) => init?.method === 'POST')).toBe(false)
  })

  it('Escenario: proyecto cerrado y Product Builder: la nota de cierre gana y no aparece la de HU-21', async () => {
    responder(proyectoCon('product_builder', '2026-11-10'))
    abrir()

    await screen.findByText('Proyecto cerrado el 2026-11-10')
    const boton = screen.getByRole('button', { name: 'Asignar integrante' }) as HTMLButtonElement
    expect(boton.disabled).toBe(true)
    expect(descripcionDe(boton)).toBe(NOTA_CERRADO)
    const formulario = screen.getByRole('form', { name: 'Asignar integrante' })
    expect(within(formulario).getAllByText(NOTA_CERRADO)).toHaveLength(1)
    expect(screen.queryByText(NOTA_HU21)).toBeNull()
  })

  it('Escenario: proyecto abierto y Scrum Master: «Asignar integrante» habilitado, sin aviso y sin nota', async () => {
    const fetchMock = responder(proyectoCon('scrum_master', null))
    abrir()

    await screen.findByText('Ana')
    await waitFor(() => expect(llamadasAlProyecto(fetchMock)).toHaveLength(1))
    const boton = screen.getByRole('button', { name: 'Asignar integrante' }) as HTMLButtonElement
    expect(boton.disabled).toBe(false)
    expect(boton.hasAttribute('aria-describedby')).toBe(false)
    expect(screen.queryByText(NOTA_CERRADO)).toBeNull()
    expect(screen.queryByText(/Proyecto cerrado el/)).toBeNull()
  })

  it('Escenario: proyecto abierto y Product Builder: se mantiene la nota de HU-21 y no la de cierre', async () => {
    responder(proyectoCon('product_builder', null))
    abrir()

    await screen.findByText('Ana')
    const boton = screen.getByRole('button', { name: 'Asignar integrante' }) as HTMLButtonElement
    await waitFor(() => expect(boton.disabled).toBe(true))
    expect(descripcionDe(boton)).toBe(NOTA_HU21)
    expect(screen.queryByText(NOTA_CERRADO)).toBeNull()
    expect(screen.queryByText(/Proyecto cerrado el/)).toBeNull()
  })

  it('Escenario: cierre desconocido (cargando): todo habilitado, sin aviso y sin nota', async () => {
    responder('pendiente')
    abrir()

    await screen.findByText('Ana')
    await screen.findByLabelText('Usuario')
    const boton = screen.getByRole('button', { name: 'Asignar integrante' }) as HTMLButtonElement
    expect(boton.disabled).toBe(false)
    expect(boton.hasAttribute('aria-describedby')).toBe(false)
    expect(screen.queryByText(NOTA_CERRADO)).toBeNull()
    expect(screen.queryByText(NOTA_HU21)).toBeNull()
    expect(screen.queryByText(/Proyecto cerrado el/)).toBeNull()
  })

  it('Escenario: cierre desconocido (error): todo habilitado, sin aviso y sin nota', async () => {
    const fetchMock = responder({ status: 500, cuerpo: { mensaje: 'falla' } })
    abrir()

    await screen.findByText('Ana')
    await waitFor(() => expect(llamadasAlProyecto(fetchMock)).toHaveLength(1))
    const boton = screen.getByRole('button', { name: 'Asignar integrante' }) as HTMLButtonElement
    expect(boton.disabled).toBe(false)
    expect(boton.hasAttribute('aria-describedby')).toBe(false)
    expect(screen.queryByText(NOTA_CERRADO)).toBeNull()
    expect(screen.queryByText(/Proyecto cerrado el/)).toBeNull()
  })

  it('Escenario: 409 «El proyecto está cerrado» al asignar: se muestra el mensaje tal cual y una sola vez', async () => {
    responder(
      proyectoCon('scrum_master', null),
      new Response(JSON.stringify({ mensaje: NOTA_CERRADO }), { status: 409 }),
    )
    abrir()

    await userEvent.selectOptions(await screen.findByLabelText('Usuario'), 'Bruno')
    await userEvent.click(screen.getByRole('button', { name: 'Asignar integrante' }))

    expect(await screen.findByText(NOTA_CERRADO)).toBeTruthy()
    expect(screen.getAllByText(NOTA_CERRADO)).toHaveLength(1)
  })

  it('Escenario: una sola llamada a GET /projects/{id} por pantalla', async () => {
    const fetchMock = responder(proyectoCon('scrum_master', '2026-11-10'))
    abrir()

    await screen.findByText('Proyecto cerrado el 2026-11-10')
    await screen.findByText('Ana')
    await screen.findByLabelText('Usuario')

    expect(llamadasAlProyecto(fetchMock)).toHaveLength(1)
  })
})
