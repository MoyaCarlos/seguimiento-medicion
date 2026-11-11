import { afterEach, describe, expect, it, vi } from 'vitest'
import { ApiError } from '../../compartido/api/client'
import { listarHistorial, obtenerDetalleSprint } from './api'
import type { DetalleSprint, SprintHistorial } from './tipos'

const sprint: SprintHistorial = {
  id: 70,
  proyecto_id: 7,
  numero: 1,
  sprint_goal: 'Gestión inicial del trámite',
  fecha_inicio: '2026-10-20',
  fecha_fin: '2026-10-26',
  metricas: {
    story_points_planificados: 8,
    story_points_completados: 5,
    velocidad: 5,
    porcentaje_historias_completadas: 50,
    horas_estimadas: null,
    horas_reales: 3.5,
    desviacion_horas: null,
    desviacion_porcentaje: null,
    defectos_detectados: 1,
    defectos_resueltos: 0,
  },
}

const detalle: DetalleSprint = {
  ...sprint,
  historias: [{ id: 1, titulo: 'Solicitar turno', estimacion_sp: 5, completada: true, minutos: 210, horas: 3.5 }],
  defectos: [],
}

const responder = (status: number, cuerpo: unknown) => {
  const simulado = vi.fn().mockResolvedValue(new Response(JSON.stringify(cuerpo), { status }))
  vi.stubGlobal('fetch', simulado)
  return simulado
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('cliente del historial de Sprints', () => {
  it('listarHistorial hace GET /projects/{id}/historial y devuelve los Sprints finalizados', async () => {
    const api = responder(200, [sprint])

    await expect(listarHistorial('7')).resolves.toEqual([sprint])

    expect(api).toHaveBeenCalledWith('/projects/7/historial', expect.objectContaining({ method: 'GET' }))
  })

  it('obtenerDetalleSprint hace GET /projects/{id}/historial/{sprintId} y devuelve el detalle', async () => {
    const api = responder(200, detalle)

    await expect(obtenerDetalleSprint('7', '70')).resolves.toEqual(detalle)

    expect(api).toHaveBeenCalledWith('/projects/7/historial/70', expect.objectContaining({ method: 'GET' }))
  })

  it('Given un proyecto inexistente, When pido el historial, Then lanza con el mensaje del backend', async () => {
    responder(404, { mensaje: 'proyecto no encontrado' })

    await expect(listarHistorial('99')).rejects.toEqual(new ApiError(404, 'proyecto no encontrado'))
  })

  it('Given un Sprint no finalizado, When pido su detalle, Then lanza con el mensaje y el campo del backend', async () => {
    responder(400, { campo: 'estado', mensaje: 'el Sprint todavía no está finalizado' })

    const error = await obtenerDetalleSprint('7', '71').catch((causa: unknown) => causa)

    expect(error).toBeInstanceOf(ApiError)
    expect(error).toMatchObject({ status: 400, field: 'estado', message: 'el Sprint todavía no está finalizado' })
  })

  it('Given una respuesta sin lista, When pido el historial, Then lanza en vez de devolver datos parciales', async () => {
    responder(200, { sprints: [] })

    await expect(listarHistorial('7')).rejects.toThrow()
  })

  it('Given un detalle sin historias ni defectos, When lo pido, Then lanza en vez de devolver datos parciales', async () => {
    responder(200, { ...sprint })

    await expect(obtenerDetalleSprint('7', '70')).rejects.toThrow()
  })
})
