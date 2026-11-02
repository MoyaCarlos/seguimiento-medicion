import { afterEach, describe, expect, it, vi } from 'vitest'
import { ApiError } from '../../compartido/api/client'
import { consultarSprintBacklog } from './api'
import type { HistoriaDeSprint } from './sprint'

function respondWithJson(body: unknown, status: number) {
  const fetchMock = vi.fn().mockResolvedValue(new Response(JSON.stringify(body), { status }))
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

const historias: HistoriaDeSprint[] = [
  { id: 1, titulo: 'Iniciar sesión', prioridad: 'M', estado: 'Nueva', estimacion_sp: 5 },
  { id: 2, titulo: 'Recuperar clave', prioridad: 'S', estado: 'En progreso', estimacion_sp: null },
]

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('cliente del Sprint Backlog', () => {
  it('US3.2: consultarSprintBacklog hace GET /sprints/{id}/backlog', async () => {
    const fetchMock = respondWithJson([], 200)

    await consultarSprintBacklog(4)

    expect(fetchMock).toHaveBeenCalledWith(
      '/sprints/4/backlog',
      expect.objectContaining({ method: 'GET' }),
    )
  })

  it('US3.2: consultarSprintBacklog devuelve las historias del Sprint', async () => {
    respondWithJson(historias, 200)

    await expect(consultarSprintBacklog(4)).resolves.toEqual(historias)
  })

  it('US3.2: consultarSprintBacklog devuelve una lista vacía si el Sprint no tiene historias', async () => {
    respondWithJson([], 200)

    await expect(consultarSprintBacklog(4)).resolves.toEqual([])
  })

  it('US3.2: consultarSprintBacklog rechaza con el mensaje de la API ante un 400', async () => {
    respondWithJson({ mensaje: 'identificador de sprint inválido', campo: 'id' }, 400)

    await expect(consultarSprintBacklog(0)).rejects.toEqual(
      new ApiError(400, 'identificador de sprint inválido', 'id'),
    )
  })
})
