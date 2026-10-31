import { afterEach, describe, expect, it, vi } from 'vitest'
import { ApiError } from '../../compartido/api/client'
import { asignarASprint, listarSprintsDelProyecto } from './api'
import type { Historia } from './historia'

function respondWithJson(body: unknown, status: number) {
  const fetchMock = vi.fn().mockResolvedValue(new Response(JSON.stringify(body), { status }))
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

const historia: Historia = {
  id: 7,
  proyecto_id: 1,
  titulo: 'Iniciar sesión',
  descripcion: 'Como usuario quiero entrar',
  criterios_aceptacion: null,
  prioridad: 'M',
  estado: 'Nueva',
  valor_negocio: 5,
  estimacion_sp: 3,
  sprint_id: 4,
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('cliente de asignación a un Sprint', () => {
  it('US3.1: asignarASprint hace PUT /backlog/{id}/sprint con el sprint_id en el cuerpo', async () => {
    const fetchMock = respondWithJson(historia, 200)

    await asignarASprint(7, 4)

    expect(fetchMock).toHaveBeenCalledWith(
      '/backlog/7/sprint',
      expect.objectContaining({
        method: 'PUT',
        body: JSON.stringify({ sprint_id: 4 }),
      }),
    )
  })

  it('US3.1: asignarASprint devuelve la Historia con el Sprint asignado', async () => {
    respondWithJson(historia, 200)

    await expect(asignarASprint(7, 4)).resolves.toEqual(historia)
  })

  it('US3.1: asignarASprint rechaza con el mensaje de la API ante un 400', async () => {
    respondWithJson({ mensaje: 'la historia debe tener estimación', campo: 'estimacion_sp' }, 400)

    await expect(asignarASprint(7, 4)).rejects.toEqual(
      new ApiError(400, 'la historia debe tener estimación', 'estimacion_sp'),
    )
  })

  it('US3.1: listarSprintsDelProyecto hace GET /projects/{id}/sprints', async () => {
    const fetchMock = respondWithJson([], 200)

    await listarSprintsDelProyecto('3')

    expect(fetchMock).toHaveBeenCalledWith(
      '/projects/3/sprints',
      expect.objectContaining({ method: 'GET' }),
    )
  })

  it('US3.1: listarSprintsDelProyecto devuelve los Sprints con id, estado y objetivo', async () => {
    const sprints = [
      { id: 1, estado: 'Finalizado', sprint_goal: 'Login' },
      { id: 2, estado: 'Pendiente' },
    ]
    respondWithJson(sprints, 200)

    await expect(listarSprintsDelProyecto('3')).resolves.toEqual(sprints)
  })

  it('US3.1: listarSprintsDelProyecto rechaza con el mensaje de la API ante un 400', async () => {
    respondWithJson({ mensaje: 'identificador de proyecto inválido' }, 400)

    await expect(listarSprintsDelProyecto('abc')).rejects.toMatchObject({
      status: 400,
      message: 'identificador de proyecto inválido',
    })
  })
})
