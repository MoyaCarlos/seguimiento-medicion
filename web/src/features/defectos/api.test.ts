import { afterEach, describe, expect, it, vi } from 'vitest'
import {
  listarDefectos,
  listarHistoriasDelProyecto,
  listarSprints,
  registrarDefecto,
  resolverDefecto,
  type SprintResumen,
} from './api'

function respuesta(valor: unknown, status = 200) {
  return new Response(JSON.stringify(valor), {
    status,
    headers: { 'Content-Type': 'application/json' },
  })
}

function prepararFetch(implementacion: (url: string, init?: RequestInit) => Response | Promise<Response>) {
  const fetchMock = vi.fn(implementacion)
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('API de defectos', () => {
  it('lista los defectos con GET /projects/{id}/defectos', async () => {
    const fetchMock = prepararFetch(() => respuesta([{ id: 7 }]))

    await expect(listarDefectos('1')).resolves.toEqual([{ id: 7 }])

    expect(fetchMock).toHaveBeenCalledWith('/projects/1/defectos', expect.objectContaining({ method: 'GET' }))
  })

  it('registra un defecto con POST y el cuerpo esperado', async () => {
    const fetchMock = prepararFetch(() => respuesta({ id: 8 }, 201))
    const cuerpo = { historia_id: 3, descripcion: 'No actualiza', severidad: 'Alta' as const }

    await registrarDefecto('1', cuerpo)

    expect(fetchMock).toHaveBeenCalledWith(
      '/projects/1/defectos',
      expect.objectContaining({ method: 'POST', body: JSON.stringify(cuerpo) }),
    )
  })

  it('resuelve un defecto con POST /defectos/{id}/resolver y sin cuerpo', async () => {
    const fetchMock = prepararFetch(() => respuesta({ id: 7, estado: 'Resuelto' }))

    await resolverDefecto(7)

    const [url, init] = fetchMock.mock.calls[0]
    expect(url).toBe('/defectos/7/resolver')
    expect(init?.method).toBe('POST')
    expect(init?.body).toBeUndefined()
  })

  it('lista los Sprints con GET /projects/{id}/sprints', async () => {
    const fetchMock = prepararFetch(() => respuesta([{ id: 2, numero: 1, estado: 'Activo' }]))

    await listarSprints('1')

    expect(fetchMock).toHaveBeenCalledWith('/projects/1/sprints', expect.objectContaining({ method: 'GET' }))
  })

  it('combina el Product Backlog y los backlogs de todos los Sprints, en orden', async () => {
    const sprints: SprintResumen[] = [
      { id: 2, numero: 1, estado: 'Finalizado' },
      { id: 3, numero: 2, estado: 'Activo' },
    ]
    prepararFetch((url) => {
      if (url === '/projects/1/backlog') return respuesta([{ id: 4, titulo: 'En backlog' }])
      if (url === '/sprints/2/backlog') return respuesta([{ id: 5, titulo: 'Vieja' }])
      if (url === '/sprints/3/backlog') return respuesta([{ id: 6, titulo: 'Activa' }])
      return respuesta([])
    })

    await expect(listarHistoriasDelProyecto('1', sprints)).resolves.toEqual([
      { id: 4, titulo: 'En backlog', origen: 'Product Backlog' },
      { id: 5, titulo: 'Vieja', origen: 'Sprint 1' },
      { id: 6, titulo: 'Activa', origen: 'Sprint 2' },
    ])
  })
})
