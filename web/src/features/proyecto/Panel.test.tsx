import { render, screen, waitFor, within } from '@testing-library/react'
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

function abrir(hoy?: () => Date) {
  render(
    <MemoryRouter initialEntries={['/projects/3/proyecto']}>
      <Routes>
        <Route path="/projects/:id/proyecto" element={<Panel hoy={hoy} />} />
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
    if (url.endsWith('/members')) return Promise.resolve(new Response('[]', { status: 200 }))
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
    if (url.endsWith('/members')) return Promise.resolve(new Response('[]', { status: 200 }))
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
    if (url.endsWith('/members')) return Promise.resolve(new Response('[]', { status: 200 }))
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

      expect(await screen.findByRole('heading', { level: 1, name: /Proyecto Alfa/ })).toBeTruthy()
      expect(await screen.findAllByText(estado)).toHaveLength(2)
      expect(fetchMock.mock.calls.some(([url]) => url === '/projects/3/status')).toBe(true)
    },
  )

  it('muestra el error del estado sin ocultar los datos del proyecto', async () => {
    responderEstado(null)
    abrir()

    expect(await screen.findByRole('heading', { level: 1, name: /Proyecto Alfa/ })).toBeTruthy()
    expect(await screen.findByText('no se pudo calcular el estado')).toBeTruthy()
    expect(await screen.findByText('Sin datos')).toBeTruthy()
  })

  it('vuelve a consultar el estado al cambiar el id del proyecto', async () => {
    const fetchMock = responderCambio('En curso')
    const router = abrirNavegable()

    expect((await screen.findAllByText('Planificado')).length).toBeGreaterThan(0)
    await router.navigate('/projects/4/proyecto')

    expect(await screen.findByRole('heading', { level: 1, name: /Proyecto Beta/ })).toBeTruthy()
    expect((await screen.findAllByText('En curso')).length).toBeGreaterThan(0)
    expect(fetchMock.mock.calls.map(([url]) => url)).toContain('/projects/4/status')
  })

  it('muestra Consultando… mientras carga el estado del nuevo proyecto', async () => {
    responderCambio('pendiente')
    const router = abrirNavegable()

    expect((await screen.findAllByText('Planificado')).length).toBeGreaterThan(0)
    await router.navigate('/projects/4/proyecto')

    expect(await screen.findByRole('heading', { level: 1, name: /Proyecto Beta/ })).toBeTruthy()
    expect(screen.getByText('Consultando…')).toBeTruthy()
    expect(screen.queryByText('Planificado')).toBeNull()
  })

  it('muestra Sin datos si la consulta del estado falla', async () => {
    responderCambio('error')
    const router = abrirNavegable()

    expect((await screen.findAllByText('Planificado')).length).toBeGreaterThan(0)
    await router.navigate('/projects/4/proyecto')

    expect(await screen.findByRole('heading', { level: 1, name: /Proyecto Beta/ })).toBeTruthy()
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

  it('Escenario 4: el panel de un proyecto ajeno muestra No sos integrante', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn(() =>
        Promise.resolve(
          new Response(JSON.stringify({ mensaje: 'No sos integrante de este proyecto' }), {
            status: 403,
          }),
        ),
      ),
    )
    abrir()

    expect(await screen.findByRole('heading', { name: 'No sos integrante de este proyecto' })).toBeTruthy()
    expect(screen.queryByRole('alert')).toBeNull()
  })

  it('muestra tal cual un estado fuera de los tres esperados', async () => {
    responderEstado('Suspendido')
    abrir()

    expect(await screen.findByRole('heading', { level: 1, name: /Proyecto Alfa/ })).toBeTruthy()
    expect(await screen.findByText('Suspendido')).toBeTruthy()
    expect(document.querySelector('.badge')).toBeNull()
  })

  it('US5.2: Given el resumen de un proyecto, When presiono "Exportar reporte del proyecto", Then se descarga el PDF del proyecto', async () => {
    const usuario = userEvent.setup()
    const descarga = simularDescarga()
    const api = responderPanelYReporte(respuestaPdf())
    abrir()
    expect(await screen.findByRole('heading', { level: 1, name: /Proyecto Alfa/ })).toBeTruthy()

    await usuario.click(screen.getByRole('button', { name: 'Exportar reporte del proyecto' }))

    await waitFor(() => expect(descarga.click).toHaveBeenCalledTimes(1))
    expect(api.mock.calls.some(([url]) => url === '/projects/3/reporte')).toBe(true)
    expect(descarga.revocar).toHaveBeenCalledWith('blob:reporte')
    expect(descarga.descargado().nombre).toBe('reporte-proyecto-3.pdf')
    expect(await descarga.crear.mock.calls[0][0].text()).toBe('%PDF-1.4 contenido')
  })

  it('Given el resumen de un proyecto, When se muestra, Then "Exportar reporte del proyecto" es el botón compartido de exportar, con su flecha', async () => {
    responderPanelYReporte(respuestaPdf())
    abrir()
    expect(await screen.findByRole('heading', { level: 1, name: /Proyecto Alfa/ })).toBeTruthy()

    const boton = screen.getByRole('button', { name: 'Exportar reporte del proyecto' })

    expect(boton.className).toContain('reporte-exportar__boton')
    expect(boton.querySelector('svg')).toBeTruthy()
  })

  it('US5.3: Given que el backend informa un error de generación, When presiono "Exportar reporte del proyecto", Then veo el mensaje y no se descarga un archivo', async () => {
    const usuario = userEvent.setup()
    const descarga = simularDescarga()
    responderPanelYReporte(
      new Response(JSON.stringify({ mensaje: 'no se pudo generar el reporte' }), { status: 500 }),
    )
    abrir()
    expect(await screen.findByRole('heading', { level: 1, name: /Proyecto Alfa/ })).toBeTruthy()

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
            url.endsWith('/members')
              ? '[]'
              : JSON.stringify(
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
          new Response(
            JSON.stringify(
              url.endsWith('/members') ? [] : url.endsWith('/status') ? { estado: 'Planificado' } : proyecto,
            ),
            { status: 200 },
          ),
        ),
      ),
    )
  }

  it('Dado el panel principal, cuando se abre, entonces «Panel principal · #3» va encima del título con el nombre del proyecto', async () => {
    responderTodo()
    abrirNavegable()

    const titulo = await screen.findByRole('heading', { level: 1, name: /Proyecto Alfa/ })
    const pretitulo = screen.getByText('Panel principal · #3')

    expect(pretitulo.compareDocumentPosition(titulo) & Node.DOCUMENT_POSITION_FOLLOWING).toBeTruthy()
    expect(screen.getAllByText(/Panel principal/)).toHaveLength(1)
    expect(screen.queryByRole('heading', { name: 'Panel del proyecto' })).toBeNull()
  })

  it('Dado «Editar proyecto», cuando se abre, entonces no muestra «Panel principal»', async () => {
    responderTodo()
    const router = abrirNavegable()
    await router.navigate('/projects/3/proyecto/edit')

    await screen.findByRole('heading', { name: 'Editar proyecto' })

    expect(screen.queryByText(/Panel principal/)).toBeNull()
  })
})

describe('cabecera y botones del panel como el prototipo', () => {
  function responder(estado: string, rol?: string, extra: object = {}) {
    vi.stubGlobal(
      'fetch',
      vi.fn((url: string) =>
        Promise.resolve(
          new Response(
            JSON.stringify(
              url.endsWith('/members') ? [] : url.endsWith('/status') ? { estado } : { ...proyecto, rol, ...extra },
            ),
            { status: 200 },
          ),
        ),
      ),
    )
  }

  it.each([
    ['Planificado', 'badge badge--planificado'],
    ['En curso', 'badge'],
    ['Finalizado', 'badge badge--finalizado'],
  ])('Given el estado %s, When abre el panel, Then el título lleva la insignia «%s» con su variante', async (estado, clases) => {
    responder(estado, 'scrum_master')
    abrir()

    const titulo = await screen.findByRole('heading', { level: 1, name: new RegExp(`Proyecto Alfa`) })
    const insignia = titulo.parentElement?.querySelector('.badge')

    expect(insignia?.textContent).toBe(estado)
    expect(insignia?.className).toBe(clases)
    expect(titulo.parentElement?.className).toBe('title-row')
  })

  it('Given que el estado todavía se consulta o falló, When abre el panel, Then el título no muestra insignia', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn((url: string) =>
        url.endsWith('/status')
          ? Promise.resolve(new Response(JSON.stringify({ mensaje: 'falló' }), { status: 500 }))
          : url.endsWith('/members')
            ? Promise.resolve(new Response('[]', { status: 200 }))
            : Promise.resolve(new Response(JSON.stringify({ ...proyecto, rol: 'scrum_master' }), { status: 200 })),
      ),
    )
    abrir()

    await screen.findByText('Sin datos')

    expect(document.querySelector('.title-row .badge')).toBeNull()
  })

  it('Given el Scrum Master, When abre el panel, Then los botones van bajo el título en orden Exportar, Configurar equipo (con ícono) y Editar proyecto (principal, con ícono)', async () => {
    responder('En curso', 'scrum_master')
    abrir()

    const equipo = await screen.findByRole('link', { name: 'Configurar equipo' })
    const editar = screen.getByRole('link', { name: 'Editar proyecto' })
    const exportar = screen.getByRole('button', { name: 'Exportar reporte del proyecto' })

    expect(equipo.getAttribute('href')).toBe('/projects/3/proyecto/settings')
    expect(editar.getAttribute('href')).toBe('/projects/3/proyecto/edit')
    expect(equipo.className).toBe('button button--secondary')
    expect(editar.className).toBe('button button--primary')
    expect(equipo.querySelector('svg')).toBeTruthy()
    expect(editar.querySelector('svg')).toBeTruthy()
    expect(exportar.compareDocumentPosition(equipo) & Node.DOCUMENT_POSITION_FOLLOWING).toBeTruthy()
    expect(equipo.compareDocumentPosition(editar) & Node.DOCUMENT_POSITION_FOLLOWING).toBeTruthy()
    expect(screen.queryByText(/Solo el Scrum Master puede editar/)).toBeNull()
  })

  it('Given un Product Builder, When abre el panel, Then «Configurar equipo» y «Editar proyecto» están deshabilitados y la nota explica por qué', async () => {
    responder('En curso', 'product_builder')
    abrir()

    const equipo = await screen.findByRole('button', { name: 'Configurar equipo' })
    const editar = screen.getByRole('button', { name: 'Editar proyecto' })
    const nota = screen.getByText('Solo el Scrum Master puede editar el proyecto y configurar el equipo.')

    expect((equipo as HTMLButtonElement).disabled).toBe(true)
    expect((editar as HTMLButtonElement).disabled).toBe(true)
    expect(equipo.getAttribute('aria-describedby')).toBe(nota.closest('[id]')?.id)
    expect(editar.getAttribute('aria-describedby')).toBe(nota.closest('[id]')?.id)
    expect(screen.queryByRole('link', { name: 'Editar proyecto' })).toBeNull()
    expect(nota.closest('[id]')?.querySelector('svg[aria-hidden="true"]')).toBeTruthy()
  })

  it('Given la descripción del proyecto, When abre el panel, Then va en una tarjeta «Descripción» con el enlace «Ver estado del proyecto» y su flecha', async () => {
    responder('En curso', 'scrum_master')
    abrir()

    const tarjeta = (await screen.findByRole('heading', { level: 2, name: 'Descripción' })).closest('article')
    const enlace = within(tarjeta as HTMLElement).getByRole('link', { name: 'Ver estado del proyecto' })

    expect(within(tarjeta as HTMLElement).getByText('Descripción', { selector: 'p' })).toBeTruthy()
    expect(enlace.getAttribute('href')).toBe('/projects/3/status')
    expect(enlace.className).toBe('enlace')
    expect(enlace.querySelector('svg')).toBeTruthy()
  })

  it('Given un proyecto sin descripción, When abre el panel, Then la tarjeta dice «Sin descripción»', async () => {
    responder('En curso', 'scrum_master', { descripcion: '' })
    abrir()

    expect(await screen.findByText('Sin descripción')).toBeTruthy()
  })

  it('Given las fechas del proyecto, When abre el panel, Then el cronograma las muestra como dd/mm/aaaa junto al estado', async () => {
    responder('En curso', 'scrum_master')
    abrir()

    const cronograma = (await screen.findByRole('heading', { level: 2, name: 'Cronograma' })).closest('section') as HTMLElement

    expect(within(cronograma).getByText('01/10/2026')).toBeTruthy()
    expect(within(cronograma).getByText('01/11/2026')).toBeTruthy()
    expect(within(cronograma).getByText('En curso')).toBeTruthy()
  })
})

describe('tarjeta «Equipo» del panel como el prototipo', () => {
  const integrantes = [
    { id: 1, proyecto_id: 3, nombre: 'carlos.moya', rol: 'scrum_master', vinculado_en: '2026-10-01T00:00:00Z' },
    { id: 2, proyecto_id: 3, nombre: 'ana.iriarte', rol: 'product_builder', vinculado_en: '2026-10-01T00:00:00Z' },
  ]

  function responder(rol: string, respuestaEquipo: () => Promise<Response>) {
    const fetchMock = vi.fn((url: string) => {
      if (url === '/projects/3/members') return respuestaEquipo()
      if (url === '/projects/3/status') {
        return Promise.resolve(new Response(JSON.stringify({ estado: 'En curso' }), { status: 200 }))
      }
      return Promise.resolve(new Response(JSON.stringify({ ...proyecto, rol }), { status: 200 }))
    })
    vi.stubGlobal('fetch', fetchMock)
    return fetchMock
  }
  const conEquipo = () => Promise.resolve(new Response(JSON.stringify(integrantes), { status: 200 }))

  async function tarjetaEquipo() {
    return (await screen.findByRole('heading', { level: 2, name: 'Equipo' })).closest('aside') as HTMLElement
  }

  it('Given un proyecto con integrantes, When abre el panel, Then la tarjeta «Equipo» lista nombre y rol de cada uno', async () => {
    const api = responder('scrum_master', conEquipo)
    abrir()

    const tarjeta = await tarjetaEquipo()
    const filas = await within(tarjeta).findAllByRole('listitem')

    expect(filas).toHaveLength(2)
    expect(within(filas[0]).getByText('carlos.moya')).toBeTruthy()
    expect(within(filas[0]).getByText('Scrum Master')).toBeTruthy()
    expect(within(filas[1]).getByText('ana.iriarte')).toBeTruthy()
    expect(within(filas[1]).getByText('Product Builder')).toBeTruthy()
    expect(api.mock.calls.filter(([url]) => url === '/projects/3/members')).toHaveLength(1)
  })

  it('Given los integrantes, When se listan, Then cada uno lleva un avatar con su inicial y colores que rotan desde el segundo', async () => {
    responder('scrum_master', conEquipo)
    abrir()

    const filas = await within(await tarjetaEquipo()).findAllByRole('listitem')
    const avatares = filas.map((fila) => fila.querySelector('.avatar-proyecto'))

    expect(avatares.map((a) => a?.getAttribute('data-inicial'))).toEqual(['C', 'A'])
    expect(avatares.map((a) => a?.className)).toEqual([
      'avatar-proyecto avatar-proyecto--1',
      'avatar-proyecto avatar-proyecto--2',
    ])
  })

  it('Given el Scrum Master, When mira la tarjeta, Then el pie dice «Configurar equipo» y lleva a la configuración', async () => {
    responder('scrum_master', conEquipo)
    abrir()

    const enlace = within(await tarjetaEquipo()).getByRole('link', { name: 'Configurar equipo' })

    expect(enlace.getAttribute('href')).toBe('/projects/3/proyecto/settings')
    expect(enlace.className).toBe('enlace')
    expect(enlace.querySelector('svg')).toBeTruthy()
  })

  it('Given un Product Builder, When mira la tarjeta, Then el pie dice «Ver equipo» y sigue llevando a la configuración', async () => {
    responder('product_builder', conEquipo)
    abrir()

    const enlace = within(await tarjetaEquipo()).getByRole('link', { name: 'Ver equipo' })

    expect(enlace.getAttribute('href')).toBe('/projects/3/proyecto/settings')
    expect(within(await tarjetaEquipo()).queryByRole('link', { name: 'Configurar equipo' })).toBeNull()
  })

  it('Given un proyecto sin integrantes, When abre el panel, Then la tarjeta lo dice', async () => {
    responder('scrum_master', () => Promise.resolve(new Response('[]', { status: 200 })))
    abrir()

    const tarjeta = await tarjetaEquipo()

    expect(await within(tarjeta).findByText('Todavía no hay integrantes asignados.')).toBeTruthy()
    expect(within(tarjeta).queryByRole('list')).toBeNull()
  })

  it('Given que el equipo falla, When abre el panel, Then la tarjeta muestra el error y el resto del panel sigue', async () => {
    responder('scrum_master', () =>
      Promise.resolve(new Response(JSON.stringify({ mensaje: 'no se pudo leer el equipo' }), { status: 500 })),
    )
    abrir()

    const tarjeta = await tarjetaEquipo()

    expect((await within(tarjeta).findByRole('alert')).textContent).toBe('no se pudo leer el equipo')
    expect(within(tarjeta).queryByText('Todavía no hay integrantes asignados.')).toBeNull()
    expect(screen.getByRole('heading', { level: 1, name: /Proyecto Alfa/ })).toBeTruthy()
  })

  it('Given que el equipo está cargando, When abre el panel, Then la tarjeta muestra el indicador de carga y no «Todavía no hay integrantes»', async () => {
    responder('scrum_master', () => new Promise<Response>(() => {}))
    abrir()

    const tarjeta = await tarjetaEquipo()

    expect(within(tarjeta).getByRole('status')).toBeTruthy()
    expect(within(tarjeta).queryByText('Todavía no hay integrantes asignados.')).toBeNull()
  })

  it('Given las dos tarjetas, When abre el panel, Then comparten la rejilla con la descripción', async () => {
    responder('scrum_master', conEquipo)
    abrir()

    const tarjeta = await tarjetaEquipo()

    expect(tarjeta.parentElement?.className).toBe('panel-rejilla')
    expect(tarjeta.className).toBe('panel-tarjeta')
    expect(tarjeta.previousElementSibling?.tagName).toBe('ARTICLE')
  })
})

describe('barra de avance del cronograma', () => {
  function responder(extra: object = {}) {
    vi.stubGlobal(
      'fetch',
      vi.fn((url: string) =>
        Promise.resolve(
          new Response(
            JSON.stringify(
              url.endsWith('/members') ? [] : url.endsWith('/status') ? { estado: 'En curso' } : { ...proyecto, rol: 'scrum_master', ...extra },
            ),
            { status: 200 },
          ),
        ),
      ),
    )
  }
  const pista = () => screen.findByRole('img', { name: /^Avance / })

  it('Given que hoy es anterior al inicio, When abre el panel, Then la barra marca 0%', async () => {
    responder()
    abrir(() => new Date(2026, 8, 15))

    const barra = await pista()

    expect(barra.getAttribute('aria-label')).toBe('Avance 0%')
    expect((barra.querySelector('i') as HTMLElement).style.width).toBe('0%')
    expect((barra.querySelector('b') as HTMLElement).style.left).toBe('0%')
  })

  it('Given que hoy está a mitad del período, When abre el panel, Then la barra marca el porcentaje transcurrido', async () => {
    responder({ fecha_inicio: '2026-10-01', fecha_fin: '2026-10-31' })
    abrir(() => new Date(2026, 9, 16))

    const barra = await pista()

    expect(barra.getAttribute('aria-label')).toBe('Avance 50%')
    expect((barra.querySelector('i') as HTMLElement).style.width).toBe('50%')
    expect((barra.querySelector('b') as HTMLElement).style.left).toBe('50%')
  })

  it('Given que hoy es posterior al fin, When abre el panel, Then la barra marca 100%', async () => {
    responder()
    abrir(() => new Date(2026, 11, 1))

    const barra = await pista()

    expect(barra.getAttribute('aria-label')).toBe('Avance 100%')
    expect((barra.querySelector('i') as HTMLElement).style.width).toBe('100%')
  })

  it('Given un proyecto sin fechas, When abre el panel, Then no se dibuja la barra', async () => {
    responder({ fecha_inicio: null, fecha_fin: null })
    abrir(() => new Date(2026, 9, 16))

    await screen.findByRole('heading', { level: 2, name: 'Cronograma' })

    expect(screen.queryByRole('img', { name: /^Avance / })).toBeNull()
  })
})
