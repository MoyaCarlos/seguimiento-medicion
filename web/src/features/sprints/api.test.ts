import { afterEach, describe, expect, it, vi } from 'vitest'
import { ApiError } from '../../compartido/api/client'
import { cambiarEstadoHistoria, consultarSprintBacklog } from './api'
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

describe('cliente de cambio de estado de una historia', () => {
  const completada: HistoriaDeSprint = {
    id: 12,
    titulo: 'Registrar horas',
    prioridad: 'M',
    estado: 'Completada',
    estimacion_sp: 5,
    fecha_finalizacion: '2026-10-29',
  }

  it('US3: cambiarEstadoHistoria hace PUT /backlog/{id}/estado', async () => {
    const fetchMock = respondWithJson(completada, 200)

    await cambiarEstadoHistoria(12, 'Completada')

    expect(fetchMock).toHaveBeenCalledWith(
      '/backlog/12/estado',
      expect.objectContaining({ method: 'PUT' }),
    )
  })

  it('US3: cambiarEstadoHistoria envía el estado en el cuerpo', async () => {
    const fetchMock = respondWithJson(completada, 200)

    await cambiarEstadoHistoria(12, 'En progreso')

    const [, init] = fetchMock.mock.calls[0]
    expect(JSON.parse(init.body)).toEqual({ estado: 'En progreso' })
  })

  it('US3: cambiarEstadoHistoria devuelve la historia con su fecha de finalización', async () => {
    respondWithJson(completada, 200)

    await expect(cambiarEstadoHistoria(12, 'Completada')).resolves.toEqual(completada)
  })

  it('US3: cambiarEstadoHistoria devuelve fecha_finalizacion null al iniciar', async () => {
    const enProgreso = { ...completada, estado: 'En progreso', fecha_finalizacion: null }
    respondWithJson(enProgreso, 200)

    await expect(cambiarEstadoHistoria(12, 'En progreso')).resolves.toEqual(enProgreso)
  })

  it('US3: cambiarEstadoHistoria rechaza con el mensaje de la API ante un 400', async () => {
    respondWithJson({ campo: 'estado', mensaje: 'primero debe pasar por En progreso' }, 400)

    await expect(cambiarEstadoHistoria(12, 'Completada')).rejects.toEqual(
      new ApiError(400, 'primero debe pasar por En progreso', 'estado'),
    )
  })

  it('US3: cambiarEstadoHistoria rechaza con el mensaje de la API ante un 404', async () => {
    respondWithJson({ mensaje: 'historia no encontrada' }, 404)

    await expect(cambiarEstadoHistoria(99, 'En progreso')).rejects.toEqual(
      new ApiError(404, 'historia no encontrada'),
    )
  })

  it('US3: cambiarEstadoHistoria rechaza con el mensaje de la API ante un 409', async () => {
    const mensaje = 'la historia cambió de estado mientras tanto: recargá y volvé a intentar'
    respondWithJson({ mensaje }, 409)

    await expect(cambiarEstadoHistoria(12, 'En progreso')).rejects.toEqual(
      new ApiError(409, mensaje),
    )
  })
})
