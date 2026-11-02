import { afterEach, describe, expect, it, vi } from 'vitest'
import {
  consultarVotacion,
  iniciarVotacion,
  listarIntegrantes,
  revelarVotacion,
  votar,
} from './api'
import type { Votacion } from './votacion'

function responder(cuerpo: unknown, status = 200) {
  const fetchMock = vi.fn().mockResolvedValue(new Response(JSON.stringify(cuerpo), { status }))
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

const VOTACION: Votacion = {
  historia_id: 3,
  ronda: 1,
  revelada_en: null,
  integrantes: [{ usuario_id: 1, nombre: 'Beta Uno', rol: 'scrum_master', ya_voto: false, valor: null }],
  promedio: null,
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('cliente de la votación', () => {
  it('consultarVotacion hace GET /backlog/{id}/votacion con usuario_id cuando hay quien consulta', async () => {
    const fetchMock = responder(VOTACION)

    await expect(consultarVotacion(3, 7)).resolves.toEqual(VOTACION)

    expect(fetchMock).toHaveBeenCalledWith(
      '/backlog/3/votacion?usuario_id=7',
      expect.objectContaining({ method: 'GET' }),
    )
  })

  it('consultarVotacion sin usuario no envía usuario_id', async () => {
    const fetchMock = responder(VOTACION)

    await consultarVotacion(3)

    expect(fetchMock).toHaveBeenCalledWith(
      '/backlog/3/votacion',
      expect.objectContaining({ method: 'GET' }),
    )
  })

  it('iniciarVotacion hace POST /backlog/{id}/votacion con usuario_id', async () => {
    const fetchMock = responder(VOTACION, 201)

    await expect(iniciarVotacion(3, 1)).resolves.toEqual(VOTACION)

    expect(fetchMock).toHaveBeenCalledWith(
      '/backlog/3/votacion',
      expect.objectContaining({ method: 'POST', body: JSON.stringify({ usuario_id: 1 }) }),
    )
  })

  it('votar hace POST /backlog/{id}/votacion/votos con usuario_id y valor', async () => {
    const fetchMock = responder(VOTACION)

    await votar(3, 2, 8)

    expect(fetchMock).toHaveBeenCalledWith(
      '/backlog/3/votacion/votos',
      expect.objectContaining({ method: 'POST', body: JSON.stringify({ usuario_id: 2, valor: 8 }) }),
    )
  })

  it('revelarVotacion hace POST /backlog/{id}/votacion/revelar con usuario_id', async () => {
    const fetchMock = responder(VOTACION)

    await revelarVotacion(3, 1)

    expect(fetchMock).toHaveBeenCalledWith(
      '/backlog/3/votacion/revelar',
      expect.objectContaining({ method: 'POST', body: JSON.stringify({ usuario_id: 1 }) }),
    )
  })

  it('listarIntegrantes hace GET /projects/{id}/members y devuelve id, nombre y rol', async () => {
    const fetchMock = responder([
      { id: 1, proyecto_id: 25, nombre: 'Beta Uno', rol: 'scrum_master', vinculado_en: '2026-11-01T00:00:00Z' },
    ])

    await expect(listarIntegrantes('25')).resolves.toEqual([
      { id: 1, nombre: 'Beta Uno', rol: 'scrum_master' },
    ])

    expect(fetchMock).toHaveBeenCalledWith(
      '/projects/25/members',
      expect.objectContaining({ method: 'GET' }),
    )
  })

  it('un 403 del servidor rechaza con ApiError y su mensaje', async () => {
    responder({ mensaje: 'solo el Scrum Master del proyecto puede hacerlo' }, 403)

    await expect(revelarVotacion(3, 2)).rejects.toMatchObject({
      status: 403,
      message: 'solo el Scrum Master del proyecto puede hacerlo',
    })
  })
})
