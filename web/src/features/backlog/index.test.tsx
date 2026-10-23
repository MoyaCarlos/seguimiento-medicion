import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { createMemoryRouter } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import { afterEach, describe, expect, it, vi } from 'vitest'
import type { HistoriaCreada } from './historia'
import { Component } from './index'

const HISTORIA: HistoriaCreada = {
  id: 7,
  proyecto_id: 42,
  titulo: 'Como usuario quiero iniciar sesión',
  descripcion: 'Autenticación con email y contraseña',
  prioridad: 'M',
  estado: 'Nueva',
  valor_negocio: null,
  estimacion_sp: null,
}

function responderCreada(overrides: Partial<HistoriaCreada> = {}) {
  const fetchMock = vi.fn().mockImplementation(() =>
    Promise.resolve(new Response(JSON.stringify({ ...HISTORIA, ...overrides }), { status: 201 })),
  )
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

function abrir(url = '/projects/42/backlog') {
  const router = createMemoryRouter(
    [{ path: '/projects/:id/backlog', element: <Component /> }],
    { initialEntries: [url] },
  )
  render(<RouterProvider router={router} />)
  return router
}

async function completarFormulario({
  titulo = 'Como usuario quiero iniciar sesión',
  descripcion = 'Autenticación con email y contraseña',
  prioridad,
}: {
  titulo?: string
  descripcion?: string
  prioridad?: string
} = {}) {
  const user = userEvent.setup()
  await user.type(screen.getByLabelText('Título'), titulo)
  await user.type(screen.getByLabelText('Descripción'), descripcion)
  if (prioridad) {
    await user.selectOptions(screen.getByLabelText('Prioridad'), prioridad)
  }
  return user
}

function cuerpoEnviado(fetchMock: ReturnType<typeof vi.fn>, llamada = 0) {
  const [, opciones] = fetchMock.mock.calls[llamada]
  return JSON.parse((opciones as RequestInit).body as string) as Record<string, unknown>
}

function botonGuardar() {
  return screen.getByRole('button', { name: 'Guardar' })
}

function mensajeDe(campo: HTMLElement) {
  const id = campo.getAttribute('aria-describedby')
  return id ? document.getElementById(id)?.textContent ?? null : null
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('creación de historias de US1', () => {
  it('US1.1: crea la historia y muestra la confirmación con estado y sin estimación', async () => {
    const fetchMock = responderCreada()
    abrir()
    const user = await completarFormulario({
      titulo: '  Como usuario quiero iniciar sesión  ',
      descripcion: '  Autenticación con email y contraseña  ',
      prioridad: 'M',
    })

    await user.click(botonGuardar())

    const confirmacion = await screen.findByRole('region', { name: 'Historia creada' })
    expect(fetchMock).toHaveBeenCalledTimes(1)
    expect(fetchMock.mock.calls[0][0]).toBe('/backlog')
    expect(cuerpoEnviado(fetchMock)).toEqual({
      proyecto_id: 42,
      titulo: 'Como usuario quiero iniciar sesión',
      descripcion: 'Autenticación con email y contraseña',
      prioridad: 'M',
    })
    expect(within(confirmacion).getByText('Como usuario quiero iniciar sesión')).toBeTruthy()
    expect(within(confirmacion).getByText('Must have')).toBeTruthy()
    expect(within(confirmacion).getByText('Nueva')).toBeTruthy()
    expect(within(confirmacion).getByText(/sin estimación/i)).toBeTruthy()
  })

  it('US1.2: restablece el formulario y conserva el proyecto para otra historia', async () => {
    const fetchMock = responderCreada()
    abrir()
    const user = await completarFormulario({ prioridad: 'S' })

    await user.click(botonGuardar())
    await screen.findByRole('region', { name: 'Historia creada' })

    expect((screen.getByLabelText('Título') as HTMLInputElement).value).toBe('')
    expect((screen.getByLabelText('Descripción') as HTMLTextAreaElement).value).toBe('')
    expect((screen.getByLabelText('Prioridad') as HTMLSelectElement).value).toBe('M')
    expect((screen.getByLabelText('Valor de negocio') as HTMLSelectElement).value).toBe('')
    expect(screen.getByRole('region', { name: 'Historia creada' })).toBeTruthy()

    await user.type(screen.getByLabelText('Título'), 'Otra historia')
    await user.type(screen.getByLabelText('Descripción'), 'Otra descripción')
    await user.click(botonGuardar())

    await waitFor(() => expect(fetchMock).toHaveBeenCalledTimes(2))
    expect(cuerpoEnviado(fetchMock, 1).proyecto_id).toBe(42)
  })

  it('US1.3: usa el proyecto de la URL al abrir o recargar la pantalla', async () => {
    const fetchMock = responderCreada({ proyecto_id: 7 })
    abrir('/projects/7/backlog')
    const user = await completarFormulario()

    await user.click(botonGuardar())

    await waitFor(() => expect(fetchMock).toHaveBeenCalledTimes(1))
    expect(cuerpoEnviado(fetchMock).proyecto_id).toBe(7)
  })
})

describe('casos límite de creación', () => {
  it.each(['abc', '0', '9007199254740993'])(
    'con un id inválido (%s) avisa y deshabilita el formulario sin llamar a la API',
    async (id) => {
      const fetchMock = vi.fn()
      vi.stubGlobal('fetch', fetchMock)

      abrir(`/projects/${id}/backlog`)

      expect(await screen.findByRole('alert')).toBeTruthy()
      expect((screen.getByLabelText('Título') as HTMLInputElement).disabled).toBe(true)
      expect((screen.getByLabelText('Descripción') as HTMLTextAreaElement).disabled).toBe(true)
      expect((screen.getByLabelText('Prioridad') as HTMLSelectElement).disabled).toBe(true)
      expect((screen.getByLabelText('Valor de negocio') as HTMLSelectElement).disabled).toBe(true)
      expect((botonGuardar() as HTMLButtonElement).disabled).toBe(true)
      expect(fetchMock).not.toHaveBeenCalled()
    },
  )

  it('FR-016: con una creación en curso deshabilita Guardar y evita un segundo envío', async () => {
    let resolver: (respuesta: Response) => void = () => {}
    const pendiente = new Promise<Response>((resolve) => {
      resolver = resolve
    })
    const fetchMock = vi.fn().mockReturnValue(pendiente)
    vi.stubGlobal('fetch', fetchMock)
    abrir()
    const user = await completarFormulario()

    await user.click(botonGuardar())

    expect(fetchMock).toHaveBeenCalledTimes(1)
    expect((botonGuardar() as HTMLButtonElement).disabled).toBe(true)
    expect(screen.getByRole('status')).toBeTruthy()

    const formulario = botonGuardar().closest('form')
    expect(formulario).not.toBeNull()
    fireEvent.submit(formulario as HTMLFormElement)
    expect(fetchMock).toHaveBeenCalledTimes(1)

    resolver(new Response(JSON.stringify(HISTORIA), { status: 201 }))
    await screen.findByRole('region', { name: 'Historia creada' })
  })

  it('permite crear dos historias con el mismo título sin tratarlo como error', async () => {
    const fetchMock = responderCreada({ titulo: 'Historia repetida' })
    abrir()
    const user = await completarFormulario({
      titulo: 'Historia repetida',
      descripcion: 'Descripción válida',
    })

    await user.click(botonGuardar())
    await screen.findByRole('region', { name: 'Historia creada' })

    await user.type(screen.getByLabelText('Título'), 'Historia repetida')
    await user.type(screen.getByLabelText('Descripción'), 'Otra descripción')
    await user.click(botonGuardar())

    await waitFor(() => expect(fetchMock).toHaveBeenCalledTimes(2))
    expect(cuerpoEnviado(fetchMock, 0).titulo).toBe('Historia repetida')
    expect(cuerpoEnviado(fetchMock, 1).titulo).toBe('Historia repetida')
    expect(screen.queryByRole('alert')).toBeNull()
  })
})

describe('validación de US2', () => {
  it('US2.1: sin título advierte junto al campo y no llama a la API', async () => {
    const fetchMock = vi.fn()
    vi.stubGlobal('fetch', fetchMock)
    abrir()
    const user = userEvent.setup()

    fireEvent.change(screen.getByLabelText('Descripción'), {
      target: { value: 'Descripción válida' },
    })
    await user.click(botonGuardar())

    const titulo = screen.getByLabelText('Título')
    expect(titulo.getAttribute('aria-invalid')).toBe('true')
    expect(mensajeDe(titulo)).toBeTruthy()
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it.each([
    ['título', 'Título'],
    ['descripción', 'Descripción'],
  ])('US2.2: con %s de solo espacios lo considera vacío', async (_caso, etiqueta) => {
    const fetchMock = vi.fn()
    vi.stubGlobal('fetch', fetchMock)
    abrir()
    const user = userEvent.setup()

    if (etiqueta === 'Título') {
      fireEvent.change(screen.getByLabelText('Título'), { target: { value: '   ' } })
      fireEvent.change(screen.getByLabelText('Descripción'), {
        target: { value: 'Descripción válida' },
      })
    } else {
      fireEvent.change(screen.getByLabelText('Título'), { target: { value: 'Título válido' } })
      fireEvent.change(screen.getByLabelText('Descripción'), { target: { value: '   ' } })
    }
    await user.click(botonGuardar())

    const campo = screen.getByLabelText(etiqueta)
    expect(campo.getAttribute('aria-invalid')).toBe('true')
    expect(mensajeDe(campo)).toBeTruthy()
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it.each([
    ['título de 201 caracteres', 'Título', 'a'.repeat(201)],
    ['descripción de 2001 caracteres', 'Descripción', 'a'.repeat(2001)],
  ])('US2.3: rechaza un %s sin llamar a la API', async (_caso, etiqueta, texto) => {
    const fetchMock = vi.fn()
    vi.stubGlobal('fetch', fetchMock)
    abrir()
    const user = userEvent.setup()

    if (etiqueta === 'Título') {
      fireEvent.change(screen.getByLabelText('Título'), { target: { value: texto } })
      fireEvent.change(screen.getByLabelText('Descripción'), {
        target: { value: 'Descripción válida' },
      })
    } else {
      fireEvent.change(screen.getByLabelText('Título'), { target: { value: 'Título válido' } })
      fireEvent.change(screen.getByLabelText('Descripción'), { target: { value: texto } })
    }
    await user.click(botonGuardar())

    const campo = screen.getByLabelText(etiqueta)
    expect(campo.getAttribute('aria-invalid')).toBe('true')
    expect(mensajeDe(campo)).toBeTruthy()
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it('US2.4: acepta exactamente 200 y 2000 caracteres', async () => {
    const fetchMock = responderCreada()
    abrir()
    const user = userEvent.setup()
    const titulo = 'a'.repeat(200)
    const descripcion = 'a'.repeat(2000)

    fireEvent.change(screen.getByLabelText('Título'), { target: { value: titulo } })
    fireEvent.change(screen.getByLabelText('Descripción'), { target: { value: descripcion } })
    await user.click(botonGuardar())

    await waitFor(() => expect(fetchMock).toHaveBeenCalledTimes(1))
    expect(cuerpoEnviado(fetchMock).titulo).toBe(titulo)
    expect(cuerpoEnviado(fetchMock).descripcion).toBe(descripcion)
  })

  it('US2.5: muestra todos los errores del intento a la vez', async () => {
    const fetchMock = vi.fn()
    vi.stubGlobal('fetch', fetchMock)
    abrir()
    const user = userEvent.setup()

    fireEvent.change(screen.getByLabelText('Descripción'), {
      target: { value: 'a'.repeat(2001) },
    })
    await user.click(botonGuardar())

    const titulo = screen.getByLabelText('Título')
    const descripcion = screen.getByLabelText('Descripción')
    expect(mensajeDe(titulo)).toBeTruthy()
    expect(mensajeDe(descripcion)).toBeTruthy()
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it('US2.6: un intento inválido reemplaza la confirmación anterior', async () => {
    const fetchMock = responderCreada()
    abrir()
    const user = await completarFormulario()

    await user.click(botonGuardar())
    await screen.findByRole('region', { name: 'Historia creada' })

    await user.click(botonGuardar())

    expect(screen.queryByRole('region', { name: 'Historia creada' })).toBeNull()
    expect(mensajeDe(screen.getByLabelText('Título'))).toBeTruthy()
    expect(fetchMock).toHaveBeenCalledTimes(1)
  })
})

describe('prioridad MoSCoW', () => {
  it('US4.1: muestra las cuatro etiquetas con "Must have" por defecto', () => {
    abrir()

    const prioridad = screen.getByLabelText('Prioridad') as HTMLSelectElement

    expect(prioridad.value).toBe('M')
    expect(Array.from(prioridad.options).map((opcion) => opcion.text)).toEqual([
      'Must have',
      'Should have',
      'Could have',
      "Won't have",
    ])
  })

  it.each([
    ['M', 'Must have'],
    ['S', 'Should have'],
    ['C', 'Could have'],
    ['W', "Won't have"],
  ] as const)('US4.2: envía la prioridad %s y muestra "%s" en la confirmación', async (codigo, etiqueta) => {
    const fetchMock = responderCreada({ prioridad: codigo })
    abrir()
    const user = await completarFormulario({ prioridad: codigo })

    await user.click(botonGuardar())

    const confirmacion = await screen.findByRole('region', { name: 'Historia creada' })
    expect(cuerpoEnviado(fetchMock).prioridad).toBe(codigo)
    expect(within(confirmacion).getByText(etiqueta)).toBeTruthy()
  })
})
