/// <reference types="node" />
import { readFileSync } from 'node:fs'
import { act, render, screen, waitFor, within } from '@testing-library/react'
import { createMemoryRouter } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import { afterEach, describe, expect, it, vi } from 'vitest'
import type { Avance } from './avance'
import { Component } from './index'

const avance: Avance = {
  total: { historias: 4, story_points: 18 },
  grupos: [
    {
      estado: 'Nueva',
      cantidad: 1,
      story_points: 3,
      historias: [{ id: 1, titulo: 'Pagar la tasa', estimacion_sp: 3, estado: 'Nueva', sprint: null }],
    },
    {
      estado: 'En progreso',
      cantidad: 1,
      story_points: 5,
      historias: [{ id: 2, titulo: 'Rendir el examen', estimacion_sp: 5, estado: 'En progreso', sprint: { id: 20, numero: 2 } }],
    },
    {
      estado: 'Completada',
      cantidad: 2,
      story_points: 10,
      historias: [
        { id: 3, titulo: 'Solicitar turno', estimacion_sp: 8, estado: 'Completada', sprint: { id: 10, numero: 1 } },
        { id: 4, titulo: 'Registrar datos', estimacion_sp: 2, estado: 'Completada', sprint: { id: 10, numero: 1 } },
      ],
    },
  ],
}

const sinHistorias: Avance = {
  total: { historias: 0, story_points: 0 },
  grupos: (['Nueva', 'En progreso', 'Completada'] as const).map((estado) => ({
    estado,
    cantidad: 0,
    story_points: 0,
    historias: [],
  })),
}

const respuesta = (status: number, cuerpo: unknown) => new Response(JSON.stringify(cuerpo), { status })

function simularApi(...respuestas: (() => Response | Promise<Response>)[]) {
  let llamada = 0
  const simulado = vi.fn(() => {
    const responder = respuestas[Math.min(llamada, respuestas.length - 1)]
    llamada += 1
    return Promise.resolve(responder())
  })
  vi.stubGlobal('fetch', simulado)
  return simulado
}

function abrir(ruta = '/projects/7/status') {
  const router = createMemoryRouter([{ path: '/projects/:id/status/*', element: <Component /> }], {
    initialEntries: [ruta],
  })
  render(<RouterProvider router={router} />)
  return router
}

const region = (nombre: string) => screen.getByRole('region', { name: nombre })

function valorDe(contenedor: HTMLElement, etiqueta: string): string | null {
  const termino = within(contenedor).getByText(etiqueta, { selector: 'dt' })
  return termino.nextElementSibling?.textContent ?? null
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('pantalla Estado', () => {
  it('Estado de carga: Given una consulta en curso, When ingreso a Estado, Then veo el título y un indicador de carga, sin datos ni errores', () => {
    simularApi(() => new Promise<Response>(() => {}))

    abrir()

    expect(screen.getByRole('heading', { level: 1, name: 'Estado del proyecto' })).toBeTruthy()
    expect(screen.getByRole('status').textContent).toMatch(/Cargando/)
    expect(screen.queryByRole('alert')).toBeNull()
    expect(screen.queryByRole('region')).toBeNull()
  })

  it('FR-001 y FR-007: When ingreso a Estado, Then hace una sola consulta de solo lectura al avance del proyecto', async () => {
    const api = simularApi(() => respuesta(200, avance))

    abrir()
    await screen.findByRole('region', { name: 'Nueva' })

    expect(api).toHaveBeenCalledTimes(1)
    expect(api).toHaveBeenCalledWith('/projects/7/avance', expect.objectContaining({ method: 'GET' }))
  })

  it('US1.1: Given historias en distintos estados, When ingreso a Estado, Then veo los tres grupos en orden con un encabezado cada uno', async () => {
    simularApi(() => respuesta(200, avance))

    abrir()

    await screen.findByRole('region', { name: 'Nueva' })
    const encabezados = screen.getAllByRole('heading', { level: 2 }).map((h) => h.textContent)
    expect(encabezados).toEqual(['Nueva', 'En progreso', 'Completada', 'Total del proyecto'])
  })

  it('US1.2: Given historias agrupadas, When ingreso a Estado, Then cada grupo informa su cantidad de historias y de Story Points', async () => {
    simularApi(() => respuesta(200, avance))

    abrir()

    await screen.findByRole('region', { name: 'Nueva' })
    expect([valorDe(region('Nueva'), 'Historias'), valorDe(region('Nueva'), 'Story Points')]).toEqual(['1', '3'])
    expect([valorDe(region('En progreso'), 'Historias'), valorDe(region('En progreso'), 'Story Points')]).toEqual(['1', '5'])
    expect([valorDe(region('Completada'), 'Historias'), valorDe(region('Completada'), 'Story Points')]).toEqual(['2', '10'])
  })

  it('US1.2: Given historias agrupadas, When ingreso a Estado, Then veo el total del proyecto', async () => {
    simularApi(() => respuesta(200, avance))

    abrir()

    await screen.findByRole('region', { name: 'Total del proyecto' })
    expect([valorDe(region('Total del proyecto'), 'Historias'), valorDe(region('Total del proyecto'), 'Story Points')]).toEqual(['4', '18'])
  })

  it('US1.3: Given una historia sin Story Points, When ingreso a Estado, Then figura "sin estimar" y no como 0', async () => {
    const conSinEstimar: Avance = {
      ...avance,
      grupos: avance.grupos.map((g) =>
        g.estado === 'Nueva'
          ? { ...g, cantidad: 2, historias: [...g.historias, { id: 9, titulo: 'Avisar por correo', estimacion_sp: null, estado: 'Nueva', sprint: null }] }
          : g,
      ),
    }
    simularApi(() => respuesta(200, conSinEstimar))

    abrir()

    const nueva = await screen.findByRole('region', { name: 'Nueva' })
    expect(valorDe(nueva, 'Historias')).toBe('2')
    expect(valorDe(nueva, 'Story Points')).toBe('3')
    const sinEstimar = within(nueva).getAllByRole('listitem')[1]
    expect(sinEstimar?.textContent).toContain('Story Points: sin estimar')
  })

  it('US1.4: Given un grupo sin historias, When ingreso a Estado, Then el grupo sigue presente con cantidad 0 y un texto de grupo vacío', async () => {
    const sinCompletadas: Avance = {
      total: { historias: 1, story_points: 3 },
      grupos: [avance.grupos[0]!, { ...sinHistorias.grupos[1]! }, { ...sinHistorias.grupos[2]! }],
    }
    simularApi(() => respuesta(200, sinCompletadas))

    abrir()

    const completada = await screen.findByRole('region', { name: 'Completada' })
    expect(valorDe(completada, 'Historias')).toBe('0')
    expect(valorDe(completada, 'Story Points')).toBe('0')
    expect(within(completada).getByText('No hay historias en este grupo.')).toBeTruthy()
    expect(screen.getByRole('region', { name: 'En progreso' })).toBeTruthy()
  })

  it('FR-007: Given las sumas del backend, When ingreso a Estado, Then las muestra tal cual, sin recalcular, y respeta el orden de la API', async () => {
    const desordenado: Avance = {
      total: { historias: 99, story_points: 77 },
      grupos: [
        {
          estado: 'Nueva',
          cantidad: 42,
          story_points: 55,
          historias: [
            { id: 9, titulo: 'Segunda por id', estimacion_sp: 1, estado: 'Nueva', sprint: null },
            { id: 2, titulo: 'Primera en la lista', estimacion_sp: 1, estado: 'Nueva', sprint: null },
          ],
        },
        sinHistorias.grupos[1]!,
        sinHistorias.grupos[2]!,
      ],
    }
    simularApi(() => respuesta(200, desordenado))

    abrir()

    const nueva = await screen.findByRole('region', { name: 'Nueva' })
    expect(valorDe(nueva, 'Historias')).toBe('42')
    expect(valorDe(nueva, 'Story Points')).toBe('55')
    expect(valorDe(region('Total del proyecto'), 'Historias')).toBe('99')
    expect(within(nueva).getAllByRole('heading', { level: 3 }).map((h) => h.textContent)).toEqual([
      'Segunda por id',
      'Primera en la lista',
    ])
  })

  it('US2.1: Given una historia en un Sprint, When ingreso a Estado, Then veo el Sprint al que pertenece', async () => {
    simularApi(() => respuesta(200, avance))

    abrir()

    const enProgreso = await screen.findByRole('region', { name: 'En progreso' })
    expect(within(enProgreso).getByRole('listitem').textContent).toContain('Sprint 2')
    const completada = region('Completada')
    expect(within(completada).getAllByRole('listitem').every((item) => item.textContent?.includes('Sprint 1'))).toBe(true)
  })

  it('US2.2: Given una historia sin Sprint, When ingreso a Estado, Then se indica que está en el Product Backlog', async () => {
    simularApi(() => respuesta(200, avance))

    abrir()

    const nueva = await screen.findByRole('region', { name: 'Nueva' })
    expect(within(nueva).getByRole('listitem').textContent).toContain('Product Backlog')
  })

  it('US3.1 y US1.4: Given un proyecto sin historias, When ingreso a Estado, Then veo el mensaje, los tres grupos con 0 historias y 0 Story Points y los totales en cero, sin errores', async () => {
    simularApi(() => respuesta(200, sinHistorias))

    abrir()

    expect(await screen.findByText('El proyecto todavía no tiene historias.')).toBeTruthy()
    for (const estado of ['Nueva', 'En progreso', 'Completada']) {
      const grupo = region(estado)
      expect([valorDe(grupo, 'Historias'), valorDe(grupo, 'Story Points')]).toEqual(['0', '0'])
      expect(within(grupo).getByText('No hay historias en este grupo.')).toBeTruthy()
    }
    const total = region('Total del proyecto')
    expect([valorDe(total, 'Historias'), valorDe(total, 'Story Points')]).toEqual(['0', '0'])
    expect(screen.queryByRole('alert')).toBeNull()
    expect(screen.queryByRole('status')).toBeNull()
  })

  it('US3.2: Given un proyecto inexistente, When ingreso a Estado, Then veo que el proyecto no existe, sin datos', async () => {
    simularApi(() => respuesta(404, { mensaje: 'proyecto no encontrado' }))

    abrir()

    expect((await screen.findByRole('alert')).textContent).toBe('El proyecto no existe')
    expect(screen.queryByRole('region')).toBeNull()
    expect(screen.queryByRole('status')).toBeNull()
  })

  it('US3.3: Given un id de ruta inválido, When ingreso a Estado, Then veo un mensaje claro y no consulto al servidor', () => {
    const api = simularApi(() => respuesta(200, avance))

    abrir('/projects/abc/status')

    expect(screen.getByRole('alert').textContent).toBe('El proyecto "abc" no es válido')
    expect(api).not.toHaveBeenCalled()
    expect(screen.queryByRole('status')).toBeNull()
  })

  it('US3.4: Given una falla de red, When ingreso a Estado, Then veo el error de conexión, sin datos parciales', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new TypeError('Failed to fetch')))

    abrir()

    expect((await screen.findByRole('alert')).textContent).toBe('No se pudo conectar con el servidor')
    expect(screen.queryByRole('region')).toBeNull()
  })

  it('Error de la API: Given un error 500, When ingreso a Estado, Then veo el mensaje de la API, sin datos parciales', async () => {
    simularApi(() => respuesta(500, { mensaje: 'error interno' }))

    abrir()

    expect((await screen.findByRole('alert')).textContent).toBe('error interno')
    expect(screen.queryByRole('region')).toBeNull()
  })

  it('Respuesta incompleta: Given una respuesta sin total, When ingreso a Estado, Then veo un mensaje de error y no una pantalla en blanco', async () => {
    simularApi(() => respuesta(200, { grupos: avance.grupos }))

    abrir()

    expect((await screen.findByRole('alert')).textContent).toBe('Ocurrió un error inesperado')
    expect(screen.getByRole('heading', { level: 1, name: 'Estado del proyecto' })).toBeTruthy()
    expect(screen.queryByRole('region')).toBeNull()
    expect(screen.queryByRole('status')).toBeNull()
  })

  it('Respuesta incompleta: Given un grupo sin su lista de historias, When ingreso a Estado, Then veo un mensaje de error y no una excepción', async () => {
    simularApi(() =>
      respuesta(200, { total: avance.total, grupos: [avance.grupos[0], { estado: 'En progreso', cantidad: 1, story_points: 5 }] }),
    )

    abrir()

    expect((await screen.findByRole('alert')).textContent).toBe('Ocurrió un error inesperado')
    expect(screen.getByRole('heading', { level: 1, name: 'Estado del proyecto' })).toBeTruthy()
    expect(screen.queryByRole('region')).toBeNull()
  })

  it('Doble carga: Given que cambio de proyecto con la pantalla montada, When llega tarde la respuesta del anterior, Then solo veo los datos del proyecto actual', async () => {
    let responderAnterior: (r: Response) => void = () => {}
    const delProyecto8: Avance = {
      ...avance,
      grupos: avance.grupos.map((g) =>
        g.estado === 'Nueva'
          ? { ...g, historias: [{ id: 8, titulo: 'Historia del proyecto 8', estimacion_sp: 3, estado: 'Nueva' as const, sprint: null }] }
          : g,
      ),
    }
    const delProyecto7: Avance = {
      ...avance,
      grupos: avance.grupos.map((g) =>
        g.estado === 'Nueva'
          ? { ...g, historias: [{ id: 7, titulo: 'Historia del proyecto 7', estimacion_sp: 3, estado: 'Nueva' as const, sprint: null }] }
          : g,
      ),
    }
    const api = simularApi(
      () => new Promise<Response>((resolver) => { responderAnterior = resolver }),
      () => respuesta(200, delProyecto8),
    )
    const router = abrir('/projects/7/status')

    await act(() => router.navigate('/projects/8/status'))
    await screen.findByText('Historia del proyecto 8')
    await act(async () => {
      responderAnterior(respuesta(200, delProyecto7))
    })

    expect(api).toHaveBeenCalledTimes(2)
    expect(api).toHaveBeenLastCalledWith('/projects/8/avance', expect.anything())
    expect(screen.getByText('Historia del proyecto 8')).toBeTruthy()
    expect(screen.queryByText('Historia del proyecto 7')).toBeNull()
  })

  it('Doble carga: Given que cambio de proyecto, When espero la respuesta del nuevo, Then no muestro los datos del anterior mientras carga', async () => {
    simularApi(
      () => respuesta(200, avance),
      () => new Promise<Response>(() => {}),
    )
    const router = abrir('/projects/7/status')
    await screen.findByText('Pagar la tasa')

    await act(() => router.navigate('/projects/8/status'))

    await waitFor(() => expect(screen.getByRole('status')).toBeTruthy())
    expect(screen.queryByText('Pagar la tasa')).toBeNull()
  })
})

describe('estilos de la pantalla Estado', () => {
  const estilos = readFileSync('src/features/estado/estado.css', 'utf8')

  it('Dado un título de historia largo y sin espacios, When se muestra, Then se parte en vez de desbordar', () => {
    expect(estilos).toMatch(/overflow-wrap:\s*anywhere/)
  })

  it('Dado un celular de 320 px, When se muestra la pantalla, Then las columnas pueden achicarse sin scroll horizontal', () => {
    expect(estilos).toMatch(/min-width:\s*0/)
    expect(estilos).toMatch(/minmax\(0,\s*1fr\)/)
  })
})
