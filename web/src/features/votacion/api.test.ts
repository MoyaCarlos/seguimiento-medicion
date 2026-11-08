import { afterEach, describe, expect, it, vi } from 'vitest'
import {
  abrirNuevaRonda,
  consultarVotacion,
  iniciarVotacion,
  listarIntegrantes,
  obtenerHistoria,
  registrarEstimacion,
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
  diferencias: null,
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

  it('abrirNuevaRonda hace POST /backlog/{id}/votacion/rondas con usuario_id y devuelve la nueva ronda', async () => {
    const nuevaRonda: Votacion = { ...VOTACION, ronda: 2 }
    const fetchMock = responder(nuevaRonda, 201)

    await expect(abrirNuevaRonda(3, 1)).resolves.toEqual(nuevaRonda)

    expect(fetchMock).toHaveBeenCalledWith(
      '/backlog/3/votacion/rondas',
      expect.objectContaining({ method: 'POST', body: JSON.stringify({ usuario_id: 1 }) }),
    )
  })

  it('registrarEstimacion hace POST /backlog/{id}/votacion/estimacion con usuario_id y estimacion_sp y devuelve la historia', async () => {
    const historia = { id: 3, titulo: 'Pagar la tasa', descripcion: 'x', prioridad: 'M', estimacion_sp: 5 }
    const fetchMock = responder(historia)

    await expect(registrarEstimacion(3, 1, 5)).resolves.toMatchObject(historia)

    expect(fetchMock).toHaveBeenCalledWith(
      '/backlog/3/votacion/estimacion',
      expect.objectContaining({ method: 'POST', body: JSON.stringify({ usuario_id: 1, estimacion_sp: 5 }) }),
    )
  })

  it('registrarEstimacion rechaza con el mensaje del servidor cuando los Story Points están congelados', async () => {
    const mensaje = 'los Story Points ya cargados no se pueden cambiar en un Sprint Activo o Finalizado'
    responder({ campo: 'estimacion_sp', mensaje }, 400)

    await expect(registrarEstimacion(3, 1, 8)).rejects.toMatchObject({ status: 400, message: mensaje })
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

  it('obtenerHistoria busca la historia en GET /projects/{id}/backlog', async () => {
    const fetchMock = responder([
      { id: 1, titulo: 'Otra', descripcion: 'x', prioridad: 'S', estimacion_sp: null },
      { id: 4, titulo: 'Pagar la tasa', descripcion: 'Pagar con tarjeta', prioridad: 'M', estimacion_sp: 8, valor_negocio: 13 },
    ])

    await expect(obtenerHistoria('25', 4)).resolves.toMatchObject({
      id: 4,
      titulo: 'Pagar la tasa',
      descripcion: 'Pagar con tarjeta',
      prioridad: 'M',
      estimacion_sp: 8,
    })

    expect(fetchMock).toHaveBeenCalledWith(
      '/projects/25/backlog',
      expect.objectContaining({ method: 'GET' }),
    )
  })

  it('obtenerHistoria devuelve null si la historia no figura en el listado', async () => {
    responder([{ id: 1, titulo: 'Otra', descripcion: 'x', prioridad: 'S', estimacion_sp: null }])

    await expect(obtenerHistoria('25', 4)).resolves.toBeNull()
  })
})
