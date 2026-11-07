import { afterEach, describe, expect, it, vi } from 'vitest'
import { ApiError } from '../../compartido/api/client'
import { consultarAvance } from './api'
import type { Avance } from './avance'

afterEach(() => {
  vi.unstubAllGlobals()
})

function responder(cuerpo: unknown, status = 200) {
  const fetchMock = vi.fn().mockResolvedValue(new Response(JSON.stringify(cuerpo), { status }))
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

const AVANCE: Avance = {
  total: { historias: 2, story_points: 8 },
  grupos: [
    {
      estado: 'Nueva',
      cantidad: 1,
      story_points: 3,
      historias: [{ id: 1, titulo: 'H1', estimacion_sp: 3, estado: 'Nueva', sprint: null }],
    },
    {
      estado: 'En progreso',
      cantidad: 1,
      story_points: 5,
      historias: [{ id: 2, titulo: 'H2', estimacion_sp: 5, estado: 'En progreso', sprint: { id: 20, numero: 2 } }],
    },
    { estado: 'Completada', cantidad: 0, story_points: 0, historias: [] },
  ],
}

describe('cliente del avance', () => {
  it('Dado un proyecto, cuando se consulta el avance, entonces hace GET /projects/{id}/avance y devuelve la respuesta', async () => {
    const fetchMock = responder(AVANCE)

    await expect(consultarAvance('7')).resolves.toEqual(AVANCE)

    expect(fetchMock).toHaveBeenCalledWith('/projects/7/avance', expect.objectContaining({ method: 'GET' }))
  })

  it('Dado una respuesta sin grupos, cuando se consulta el avance, entonces rechaza en vez de mostrar datos parciales', async () => {
    responder({ total: { historias: 0, story_points: 0 } })

    await expect(consultarAvance('7')).rejects.toThrow('La respuesta no trae el avance')
  })

  it.each([
    ['sin total', { grupos: AVANCE.grupos }],
    ['con total nulo', { total: null, grupos: AVANCE.grupos }],
    ['con total que no es un objeto', { total: 8, grupos: AVANCE.grupos }],
    ['con total sin historias ni story_points', { total: {}, grupos: AVANCE.grupos }],
    ['con total sin story_points', { total: { historias: 2 }, grupos: AVANCE.grupos }],
    ['con total sin historias', { total: { story_points: 8 }, grupos: AVANCE.grupos }],
    ['con historias que no es un número', { total: { historias: '2', story_points: 8 }, grupos: AVANCE.grupos }],
    ['con story_points nulo', { total: { historias: 2, story_points: null }, grupos: AVANCE.grupos }],
  ])(
    'Dado una respuesta %s, cuando se consulta el avance, entonces rechaza en vez de dejar la pantalla sin datos',
    async (_caso, cuerpo) => {
      responder(cuerpo)

      await expect(consultarAvance('7')).rejects.toThrow('La respuesta no trae el avance')
    },
  )

  it('Dado una respuesta vacía, cuando se consulta el avance, entonces rechaza', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response(null, { status: 200 })))

    await expect(consultarAvance('7')).rejects.toThrow('La respuesta no trae el avance')
  })

  it('Dado un proyecto inexistente, cuando se consulta el avance, entonces propaga el 404 con el mensaje de la API', async () => {
    responder({ mensaje: 'proyecto no encontrado' }, 404)

    const error = await consultarAvance('999').catch((causa: unknown) => causa)

    expect(error).toBeInstanceOf(ApiError)
    expect(error).toMatchObject({ status: 404, message: 'proyecto no encontrado' })
  })
})
