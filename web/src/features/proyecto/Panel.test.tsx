import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { createMemoryRouter, MemoryRouter, Route, Routes } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import { afterEach, describe, expect, it, onTestFinished, vi } from 'vitest'
import { Component } from './index'
import estilos from './estilos.css?raw'
import { Panel } from './Panel'

const proyecto = {
  id: 3,
  nombre: 'Proyecto Alfa',
  descripcion: 'Descripción',
  fecha_inicio: '2026-10-01',
  fecha_fin: '2026-11-01',
  creado_en: '2026-09-01T00:00:00Z',
}

const proyectoBeta = { ...proyecto, id: 4, nombre: 'Proyecto Beta' }

function abrir() {
  render(
    <MemoryRouter initialEntries={['/projects/3/proyecto']}>
      <Routes>
        <Route path="/projects/:id/proyecto" element={<Panel />} />
      </Routes>
    </MemoryRouter>,
  )
}

function abrirNavegable() {
  const router = createMemoryRouter([{ path: '/projects/:id/proyecto/*', element: <Component /> }], {
    initialEntries: ['/projects/3/proyecto'],
  })
  render(<RouterProvider router={router} />)
  return router
}

type EstadoDelNuevo = 'En curso' | 'pendiente' | 'error'

function responderCambio(estadoDelNuevo: EstadoDelNuevo) {
  const fetchMock = vi.fn((url: string) => {
    if (url === '/projects/3') {
      return Promise.resolve(new Response(JSON.stringify(proyecto), { status: 200 }))
    }
    if (url === '/projects/3/status') {
      return Promise.resolve(
        new Response(JSON.stringify({ estado: 'Planificado' }), { status: 200 }),
      )
    }
    if (url === '/projects/4') {
      return Promise.resolve(new Response(JSON.stringify(proyectoBeta), { status: 200 }))
    }
    if (url === '/projects/4/status') {
      if (estadoDelNuevo === 'pendiente') return new Promise<Response>(() => {})
      if (estadoDelNuevo === 'error') {
        return Promise.resolve(
          new Response(JSON.stringify({ mensaje: 'no se pudo calcular el estado' }), {
            status: 500,
          }),
        )
      }
      return Promise.resolve(
        new Response(JSON.stringify({ estado: estadoDelNuevo }), { status: 200 }),
      )
    }
    return Promise.reject(new Error(`URL inesperada: ${url}`))
  })
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

function responderEstado(estado: string | null) {
  const fetchMock = vi.fn((url: string) => {
    if (url === '/projects/3') {
      return Promise.resolve(new Response(JSON.stringify(proyecto), { status: 200 }))
    }
    if (url === '/projects/3/status') {
      if (estado === null) {
        return Promise.resolve(
          new Response(JSON.stringify({ mensaje: 'no se pudo calcular el estado' }), {
            status: 500,
          }),
        )
      }
      return Promise.resolve(new Response(JSON.stringify({ estado }), { status: 200 }))
    }
    return Promise.reject(new Error(`URL inesperada: ${url}`))
  })
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

function respuestaPdf(nombre = 'reporte-proyecto-3.pdf') {
  return new Response('%PDF-1.4 contenido', {
    status: 200,
    headers: {
      'Content-Type': 'application/pdf',
      'Content-Disposition': `attachment; filename="${nombre}"`,
    },
  })
}

function simularDescarga() {
  const crear = vi.fn((_blob: Blob) => 'blob:reporte')
  const revocar = vi.fn()
  vi.stubGlobal('URL', { createObjectURL: crear, revokeObjectURL: revocar })

  let descargado: { nombre?: string; href?: string } = {}
  const click = vi.spyOn(HTMLAnchorElement.prototype, 'click').mockImplementation(function (
    this: HTMLAnchorElement,
  ) {
    descargado = { nombre: this.download, href: this.href }
  })

  return { crear, revocar, click, descargado: () => descargado }
}

function responderPanelYReporte(reporte: Response) {
  const fetchMock = vi.fn((url: string) => {
    if (url === '/projects/3') {
      return Promise.resolve(new Response(JSON.stringify(proyecto), { status: 200 }))
    }
    if (url === '/projects/3/status') {
      return Promise.resolve(
        new Response(JSON.stringify({ estado: 'Planificado' }), { status: 200 }),
      )
    }
    if (url === '/projects/3/reporte') {
      return Promise.resolve(reporte)
    }
    return Promise.reject(new Error(`URL inesperada: ${url}`))
  })
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

afterEach(() => {
  vi.unstubAllGlobals()
  vi.restoreAllMocks()
})

describe('panel del proyecto', () => {
  it.each(['Planificado', 'En curso', 'Finalizado'])(
    'muestra el estado %s que devuelve la API',
    async (estado) => {
      const fetchMock = responderEstado(estado)
      abrir()

      expect(await screen.findByText('Proyecto Alfa')).toBeTruthy()
      expect(await screen.findByText(estado)).toBeTruthy()
      expect(fetchMock.mock.calls.some(([url]) => url === '/projects/3/status')).toBe(true)
    },
  )

  it('muestra el error del estado sin ocultar los datos del proyecto', async () => {
    responderEstado(null)
    abrir()

    expect(await screen.findByText('Proyecto Alfa')).toBeTruthy()
    expect(await screen.findByText('no se pudo calcular el estado')).toBeTruthy()
    expect(await screen.findByText('Sin datos')).toBeTruthy()
  })

  it('vuelve a consultar el estado al cambiar el id del proyecto', async () => {
    const fetchMock = responderCambio('En curso')
    const router = abrirNavegable()

    expect(await screen.findByText('Planificado')).toBeTruthy()
    await router.navigate('/projects/4/proyecto')

    expect(await screen.findByText('Proyecto Beta')).toBeTruthy()
    expect(await screen.findByText('En curso')).toBeTruthy()
    expect(fetchMock.mock.calls.map(([url]) => url)).toContain('/projects/4/status')
  })

  it('muestra Consultando… mientras carga el estado del nuevo proyecto', async () => {
    responderCambio('pendiente')
    const router = abrirNavegable()

    expect(await screen.findByText('Planificado')).toBeTruthy()
    await router.navigate('/projects/4/proyecto')

    expect(await screen.findByText('Proyecto Beta')).toBeTruthy()
    expect(screen.getByText('Consultando…')).toBeTruthy()
    expect(screen.queryByText('Planificado')).toBeNull()
  })

  it('muestra Sin datos si la consulta del estado falla', async () => {
    responderCambio('error')
    const router = abrirNavegable()

    expect(await screen.findByText('Planificado')).toBeTruthy()
    await router.navigate('/projects/4/proyecto')

    expect(await screen.findByText('Proyecto Beta')).toBeTruthy()
    expect(await screen.findByText('Sin datos')).toBeTruthy()
    expect(screen.getByText('no se pudo calcular el estado')).toBeTruthy()
    expect(screen.queryByText('Planificado')).toBeNull()
  })

  it('muestra el error de un proyecto inexistente y no inventa un estado', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn(() =>
        Promise.resolve(
          new Response(JSON.stringify({ mensaje: 'no existe el proyecto' }), { status: 404 }),
        ),
      ),
    )
    abrir()

    expect(await screen.findByText('no existe el proyecto')).toBeTruthy()
    expect(screen.queryByText('Planificado')).toBeNull()
    expect(screen.queryByText('Sin datos')).toBeNull()
  })

  it('muestra tal cual un estado fuera de los tres esperados', async () => {
    responderEstado('Suspendido')
    abrir()

    expect(await screen.findByText('Proyecto Alfa')).toBeTruthy()
    expect(await screen.findByText('Suspendido')).toBeTruthy()
  })

  it('US5.2: Given el resumen de un proyecto, When presiono "Exportar reporte del proyecto", Then se descarga el PDF del proyecto', async () => {
    const usuario = userEvent.setup()
    const descarga = simularDescarga()
    const api = responderPanelYReporte(respuestaPdf())
    abrir()
    expect(await screen.findByText('Proyecto Alfa')).toBeTruthy()

    await usuario.click(screen.getByRole('button', { name: 'Exportar reporte del proyecto' }))

    await waitFor(() => expect(descarga.click).toHaveBeenCalledTimes(1))
    expect(api.mock.calls.some(([url]) => url === '/projects/3/reporte')).toBe(true)
    expect(descarga.revocar).toHaveBeenCalledWith('blob:reporte')
    expect(descarga.descargado().nombre).toBe('reporte-proyecto-3.pdf')
    expect(await descarga.crear.mock.calls[0][0].text()).toBe('%PDF-1.4 contenido')
  })

  it('US5.3: Given que el backend informa un error de generación, When presiono "Exportar reporte del proyecto", Then veo el mensaje y no se descarga un archivo', async () => {
    const usuario = userEvent.setup()
    const descarga = simularDescarga()
    responderPanelYReporte(
      new Response(JSON.stringify({ mensaje: 'no se pudo generar el reporte' }), { status: 500 }),
    )
    abrir()
    expect(await screen.findByText('Proyecto Alfa')).toBeTruthy()

    await usuario.click(screen.getByRole('button', { name: 'Exportar reporte del proyecto' }))

    expect((await screen.findByRole('alert')).textContent).toBe('no se pudo generar el reporte')
    expect(descarga.crear).not.toHaveBeenCalled()
    expect(descarga.click).not.toHaveBeenCalled()
  })

  it('dado una descripción larga sin espacios, permite cortarla en varias líneas', async () => {
    const descripcion = 'a'.repeat(300)
    vi.stubGlobal(
      'fetch',
      vi.fn((url: string) =>
        Promise.resolve(
          new Response(
            JSON.stringify(
              url === '/projects/3' ? { ...proyecto, descripcion } : { estado: 'Planificado' },
            ),
            { status: 200 },
          ),
        ),
      ),
    )
    const hojaDeEstilos = document.createElement('style')
    hojaDeEstilos.textContent = estilos
    document.head.append(hojaDeEstilos)
    onTestFinished(() => hojaDeEstilos.remove())
    abrir()

    const parrafo = await screen.findByText(descripcion)

    expect(getComputedStyle(parrafo).getPropertyValue('overflow-wrap')).toBe('anywhere')
  })
})

describe('pretítulo del panel', () => {
  function responderTodo() {
    vi.stubGlobal(
      'fetch',
      vi.fn((url: string) =>
        Promise.resolve(
          new Response(JSON.stringify(url.endsWith('/status') ? { estado: 'Planificado' } : proyecto), { status: 200 }),
        ),
      ),
    )
  }

  it('Dado el panel principal, cuando se abre, entonces «Panel principal» va encima del título', async () => {
    responderTodo()
    abrirNavegable()

    const titulo = await screen.findByRole('heading', { level: 1, name: 'Panel del proyecto' })
    const pretitulo = screen.getByText('Panel principal')

    expect(pretitulo.compareDocumentPosition(titulo) & Node.DOCUMENT_POSITION_FOLLOWING).toBeTruthy()
    expect(screen.getAllByText('Panel principal')).toHaveLength(1)
  })

  it('Dado «Editar proyecto», cuando se abre, entonces no muestra «Panel principal»', async () => {
    responderTodo()
    const router = abrirNavegable()
    await router.navigate('/projects/3/proyecto/edit')

    await screen.findByRole('heading', { name: 'Editar proyecto' })

    expect(screen.queryByText('Panel principal')).toBeNull()
  })
})
