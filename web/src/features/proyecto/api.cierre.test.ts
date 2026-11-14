import { afterEach, describe, expect, it, vi } from 'vitest'
import { cerrarProyecto, reabrirProyecto } from './api'

const proyectoBase = {
  id: 7,
  nombre: 'Proyecto Alfa',
  descripcion: '',
  fecha_inicio: null,
  fecha_fin: null,
  creado_en: '2026-09-01T00:00:00Z',
  horas_por_story_point: null,
}

function responder(cuerpo: unknown) {
  const fetchMock = vi.fn().mockResolvedValue(new Response(JSON.stringify(cuerpo), { status: 200 }))
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('cliente de cierre de proyecto', () => {
  it('cerrarProyecto hace POST /projects/7/cerrar sin cuerpo y devuelve fecha_cierre y estado Finalizado', async () => {
    const fetchMock = responder({ ...proyectoBase, fecha_cierre: '2026-11-11', estado: 'Finalizado' })

    const respuesta = await cerrarProyecto('7')

    const [url, init] = fetchMock.mock.calls[0]
    expect(url).toBe('/projects/7/cerrar')
    expect(init.method).toBe('POST')
    expect(init.body).toBeUndefined()
    expect(respuesta.fecha_cierre).toBe('2026-11-11')
    expect(respuesta.estado).toBe('Finalizado')
  })

  it('reabrirProyecto hace POST /projects/7/reabrir sin cuerpo y devuelve fecha_cierre null', async () => {
    const fetchMock = responder({ ...proyectoBase, fecha_cierre: null, estado: 'En curso' })

    const respuesta = await reabrirProyecto('7')

    const [url, init] = fetchMock.mock.calls[0]
    expect(url).toBe('/projects/7/reabrir')
    expect(init.method).toBe('POST')
    expect(init.body).toBeUndefined()
    expect(respuesta.fecha_cierre).toBeNull()
    expect(respuesta.estado).toBe('En curso')
  })
})
