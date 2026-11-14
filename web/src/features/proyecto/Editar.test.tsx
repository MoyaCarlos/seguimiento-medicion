import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Route, Routes } from 'react-router'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { Editar } from './Editar'

const proyecto = {
  id: 3,
  nombre: 'Proyecto Alfa',
  descripcion: 'Descripción original',
  fecha_inicio: '2026-10-01',
  fecha_fin: '2026-11-01',
  creado_en: '2026-09-01T00:00:00Z',
  horas_por_story_point: 4,
}

function abrir() {
  render(
    <MemoryRouter initialEntries={['/projects/3/proyecto/edit']}>
      <Routes>
        <Route path="/projects/:id/proyecto/edit" element={<Editar />} />
      </Routes>
    </MemoryRouter>,
  )
}

function responderProyecto() {
  const fetchMock = vi.fn().mockResolvedValue(new Response(JSON.stringify(proyecto), { status: 200 }))
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('pantalla de edición de proyecto', () => {
  it('carga los datos actuales del proyecto', async () => {
    responderProyecto()
    abrir()

    expect(await screen.findByDisplayValue('Proyecto Alfa')).toBeTruthy()
    expect(screen.getByDisplayValue('Descripción original')).toBeTruthy()
    expect(screen.getByDisplayValue('2026-10-01')).toBeTruthy()
  })

  it('muestra el error de nombre vacío y no guarda', async () => {
    const fetchMock = responderProyecto()
    abrir()

    const nombre = await screen.findByDisplayValue('Proyecto Alfa')
    await userEvent.clear(nombre)
    await userEvent.click(screen.getByRole('button', { name: 'Guardar cambios' }))

    expect(screen.getByText('El nombre es obligatorio')).toBeTruthy()
    expect(fetchMock.mock.calls.filter(([, init]) => init?.method === 'PUT')).toHaveLength(0)
  })

  it('muestra el error cuando la fecha de fin es igual a la de inicio', async () => {
    const fetchMock = responderProyecto()
    abrir()

    await screen.findByDisplayValue('Proyecto Alfa')
    fireEvent.change(screen.getByLabelText('Fecha de fin'), { target: { value: '2026-10-01' } })
    await userEvent.click(screen.getByRole('button', { name: 'Guardar cambios' }))

    expect(
      screen.getByText('La fecha de fin debe ser posterior a la fecha de inicio'),
    ).toBeTruthy()
    expect(fetchMock.mock.calls.filter(([, init]) => init?.method === 'PUT')).toHaveLength(0)
  })

  it('guarda los cambios válidos con PUT', async () => {
    const fetchMock = vi.fn((_url: string, init?: RequestInit) =>
      Promise.resolve(
        new Response(JSON.stringify(proyecto), {
          status: init?.method === 'PUT' ? 200 : 200,
        }),
      ),
    )
    vi.stubGlobal('fetch', fetchMock)
    abrir()

    const nombre = await screen.findByDisplayValue('Proyecto Alfa')
    await userEvent.clear(nombre)
    await userEvent.type(nombre, 'Proyecto Beta')
    await userEvent.click(screen.getByRole('button', { name: 'Guardar cambios' }))

    await waitFor(() =>
      expect(fetchMock.mock.calls.some(([, init]) => init?.method === 'PUT')).toBe(true),
    )

    const llamada = fetchMock.mock.calls.find(([, init]) => init?.method === 'PUT')
    expect(JSON.parse(String(llamada?.[1]?.body))).toEqual({
      nombre: 'Proyecto Beta',
      descripcion: 'Descripción original',
      fecha_inicio: '2026-10-01',
      fecha_fin: '2026-11-01',
      horas_por_story_point: 4,
    })
  })

  describe('horas por Story Point', () => {
    function responderEdicion(respuestaPut?: Response) {
      const fetchMock = vi.fn((_url: string, init?: RequestInit) =>
        Promise.resolve(
          init?.method === 'PUT' && respuestaPut
            ? respuestaPut
            : new Response(JSON.stringify(proyecto), { status: 200 }),
        ),
      )
      vi.stubGlobal('fetch', fetchMock)
      return fetchMock
    }

    async function guardarCon(valor: string, fetchMock: ReturnType<typeof responderEdicion>) {
      const campo = await screen.findByLabelText(/Horas por Story Point/)
      await userEvent.clear(campo)
      if (valor) await userEvent.type(campo, valor)
      await userEvent.click(screen.getByRole('button', { name: 'Guardar cambios' }))
      await waitFor(() =>
        expect(fetchMock.mock.calls.some(([, init]) => init?.method === 'PUT')).toBe(true),
      )
      const llamada = fetchMock.mock.calls.find(([, init]) => init?.method === 'PUT')
      return JSON.parse(String(llamada?.[1]?.body))
    }

    it('US4.1: dado un proyecto con el parámetro, al abrir la edición el campo muestra el valor actual', async () => {
      responderEdicion()
      abrir()

      const campo = await screen.findByLabelText<HTMLInputElement>(/Horas por Story Point/)
      expect(campo.value).toBe('4')
    })

    it('US4.1: dado un proyecto sin el parámetro, el campo aparece vacío', async () => {
      vi.stubGlobal(
        'fetch',
        vi.fn().mockResolvedValue(
          new Response(JSON.stringify({ ...proyecto, horas_por_story_point: null }), { status: 200 }),
        ),
      )
      abrir()

      const campo = await screen.findByLabelText<HTMLInputElement>(/Horas por Story Point/)
      expect(campo.value).toBe('')
    })

    it('US4.2: al cargar un valor válido y guardar, se envía el número', async () => {
      const fetchMock = responderEdicion()
      abrir()

      const cuerpo = await guardarCon('2,5', fetchMock)

      expect(cuerpo.horas_por_story_point).toBe(2.5)
    })

    it('US4.3: si el sistema rechaza el valor, la pantalla muestra el mensaje junto al campo', async () => {
      const mensaje = 'las horas por Story Point deben ser un número mayor que cero'
      const fetchMock = responderEdicion(
        new Response(JSON.stringify({ campo: 'horas_por_story_point', mensaje }), { status: 400 }),
      )
      abrir()

      await guardarCon('8', fetchMock)

      const campo = await screen.findByLabelText(/Horas por Story Point/)
      const error = await screen.findByText(mensaje)
      expect(campo.getAttribute('aria-describedby')).toBe(error.id)
    })

    it('US4.3: un valor inválido muestra el error y no se envía', async () => {
      const fetchMock = responderEdicion()
      abrir()

      const campo = await screen.findByLabelText(/Horas por Story Point/)
      await userEvent.clear(campo)
      await userEvent.type(campo, 'abc')
      await userEvent.click(screen.getByRole('button', { name: 'Guardar cambios' }))

      expect(
        screen.getByText('Las horas por Story Point deben ser un número mayor que cero'),
      ).toBeTruthy()
      expect(fetchMock.mock.calls.filter(([, init]) => init?.method === 'PUT')).toHaveLength(0)
    })

    it('US4.4: al dejar el campo vacío y guardar, se envía sin valor', async () => {
      const fetchMock = responderEdicion()
      abrir()

      const cuerpo = await guardarCon('', fetchMock)

      expect(cuerpo.horas_por_story_point).toBeNull()
    })
  })
})

describe('edición de proyecto - rol (HU-21)', () => {
  const nota = (boton: HTMLElement) =>
    document.getElementById(boton.getAttribute('aria-describedby') ?? '')?.textContent

  function responderConRol(rol: 'scrum_master' | 'product_builder', alGuardar?: () => Response) {
    const fetchMock = vi.fn((_url: string, init?: RequestInit) =>
      Promise.resolve(
        init?.method === 'PUT' && alGuardar
          ? alGuardar()
          : new Response(JSON.stringify({ ...proyecto, rol }), { status: 200 }),
      ),
    )
    vi.stubGlobal('fetch', fetchMock)
    return fetchMock
  }

  it('Escenario 2: como Product Builder veo «Guardar cambios» deshabilitado con la nota «Solo el Scrum Master puede guardar el proyecto»', async () => {
    responderConRol('product_builder')
    abrir()

    await screen.findByDisplayValue('Proyecto Alfa')
    const boton = screen.getByRole('button', { name: 'Guardar cambios' })
    await waitFor(() => expect(boton.hasAttribute('disabled')).toBe(true))
    expect(nota(boton)).toBe('Solo el Scrum Master puede guardar el proyecto')
  })

  it('Escenario 1: como Scrum Master «Guardar cambios» está habilitado', async () => {
    responderConRol('scrum_master')
    abrir()

    await screen.findByDisplayValue('Proyecto Alfa')
    const boton = screen.getByRole('button', { name: 'Guardar cambios' })
    expect(boton.hasAttribute('disabled')).toBe(false)
    expect(boton.hasAttribute('aria-describedby')).toBe(false)
  })

  it('un 403 de rol inesperado muestra «No tenés permiso» sin perder lo escrito', async () => {
    responderConRol('scrum_master', () =>
      new Response(JSON.stringify({ mensaje: 'Solo el Scrum Master puede realizar esta acción' }), { status: 403 }),
    )
    abrir()

    const nombre = await screen.findByDisplayValue('Proyecto Alfa')
    await userEvent.clear(nombre)
    await userEvent.type(nombre, 'Proyecto Beta')
    await userEvent.click(screen.getByRole('button', { name: 'Guardar cambios' }))

    expect(await screen.findByText('No tenés permiso')).toBeTruthy()
    expect(screen.getByText('Solo el Scrum Master puede realizar esta acción')).toBeTruthy()
    expect(screen.getByDisplayValue('Proyecto Beta')).toBeTruthy()
  })
})

describe('edición de proyecto - cierre (#312)', () => {
  const NOTA_CERRADO = 'El proyecto está cerrado'
  const NOTA_GESTION = 'Solo el Scrum Master puede guardar el proyecto'

  function responderCierre(opciones: {
    fechaCierre: string | null
    rol?: 'scrum_master' | 'product_builder'
    hook?: 'normal' | 'pendiente' | 'falla'
  }) {
    const { fechaCierre, rol = 'scrum_master', hook = 'normal' } = opciones
    let lecturas = 0
    const fetchMock = vi.fn((_url: string, init?: RequestInit) => {
      if (init?.method === 'POST') {
        return Promise.resolve(
          new Response(
            JSON.stringify({ ...proyecto, rol, fecha_cierre: '2026-11-11', estado: 'Finalizado' }),
            { status: 200 },
          ),
        )
      }
      lecturas += 1
      if (lecturas === 2 && hook === 'pendiente') return new Promise<Response>(() => {})
      if (lecturas === 2 && hook === 'falla') {
        return Promise.resolve(new Response(JSON.stringify({ mensaje: 'falló' }), { status: 500 }))
      }
      return Promise.resolve(
        new Response(JSON.stringify({ ...proyecto, rol, fecha_cierre: fechaCierre }), { status: 200 }),
      )
    })
    vi.stubGlobal('fetch', fetchMock)
    return fetchMock
  }

  const guardar = () => screen.getByRole('button', { name: 'Guardar cambios' })

  it('Escenario: Scrum Master abre Editar proyecto y aparece la zona «Cerrar proyecto»', async () => {
    responderCierre({ fechaCierre: null })
    abrir()

    await screen.findByDisplayValue('Proyecto Alfa')

    expect(await screen.findByRole('region', { name: 'Cerrar proyecto' })).toBeTruthy()
    expect(screen.getByRole('button', { name: 'Cerrar proyecto' })).toBeTruthy()
  })

  it('SC-001: cerrar y confirmar deja «Proyecto cerrado el AAAA-MM-DD» y «Finalizado» sin recargar, con «Guardar cambios» deshabilitado y la nota «El proyecto está cerrado»', async () => {
    responderCierre({ fechaCierre: null })
    abrir()

    await screen.findByDisplayValue('Proyecto Alfa')
    await userEvent.click(await screen.findByRole('button', { name: 'Cerrar proyecto' }))
    await userEvent.click(screen.getByRole('button', { name: 'Confirmar' }))

    expect(await screen.findByText('Proyecto cerrado el 2026-11-11')).toBeTruthy()
    expect(screen.getByText('Finalizado')).toBeTruthy()
    await waitFor(() => expect(guardar().hasAttribute('disabled')).toBe(true))
    expect(document.getElementById(guardar().getAttribute('aria-describedby') ?? '')?.textContent).toBe(
      NOTA_CERRADO,
    )
  })

  it('Escenario: con el proyecto ya cerrado al abrir se ve el aviso y la nota única, sin la nota de HU-21', async () => {
    responderCierre({ fechaCierre: '2026-11-11' })
    abrir()

    await screen.findByDisplayValue('Proyecto Alfa')

    expect(await screen.findByText('Proyecto cerrado el 2026-11-11')).toBeTruthy()
    expect(screen.getAllByText(NOTA_CERRADO)).toHaveLength(1)
    expect(screen.queryByText(NOTA_GESTION)).toBeNull()
    expect(guardar().hasAttribute('disabled')).toBe(true)
  })

  it('Escenario: la nota «El proyecto está cerrado» no aparece mientras el rol y el cierre se cargan y «Guardar cambios» sigue habilitado', async () => {
    responderCierre({ fechaCierre: '2026-11-11', hook: 'pendiente' })
    abrir()

    await screen.findByDisplayValue('Proyecto Alfa')

    expect(screen.queryByText(NOTA_CERRADO)).toBeNull()
    expect(guardar().hasAttribute('disabled')).toBe(false)
  })

  it('Escenario: la nota «El proyecto está cerrado» no aparece si el GET del hook falla y «Guardar cambios» sigue habilitado', async () => {
    const fetchMock = responderCierre({ fechaCierre: '2026-11-11', hook: 'falla' })
    abrir()

    await screen.findByDisplayValue('Proyecto Alfa')
    await waitFor(() => expect(fetchMock.mock.calls.length).toBeGreaterThanOrEqual(2))

    expect(screen.queryByText(NOTA_CERRADO)).toBeNull()
    expect(guardar().hasAttribute('disabled')).toBe(false)
  })

  it('Escenario: la nota «El proyecto está cerrado» no aparece en un proyecto abierto y «Guardar cambios» sigue habilitado', async () => {
    responderCierre({ fechaCierre: null })
    abrir()

    await screen.findByDisplayValue('Proyecto Alfa')

    expect(screen.queryByText(NOTA_CERRADO)).toBeNull()
    expect(screen.queryByText('Proyecto cerrado el 2026-11-11')).toBeNull()
    expect(guardar().hasAttribute('disabled')).toBe(false)
  })
})
