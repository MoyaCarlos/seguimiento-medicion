import { act, renderHook } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { ApiError } from '../../compartido/api/client'
import { consultarVotacion, listarIntegrantes, obtenerHistoria, votar } from './api'
import { useVotacion } from './useVotacion'
import type { Votacion } from './votacion'

vi.mock('./api')

const consultarMock = vi.mocked(consultarVotacion)
const integrantesMock = vi.mocked(listarIntegrantes)
const historiaMock = vi.mocked(obtenerHistoria)
const votarMock = vi.mocked(votar)

const INTERVALO = 5000
let visibilidad: 'visible' | 'hidden' = 'visible'

function votacion(overrides: Partial<Votacion> = {}): Votacion {
  return {
    historia_id: 3,
    ronda: 1,
    revelada_en: null,
    integrantes: [
      { usuario_id: 1, nombre: 'Beta Uno', rol: 'scrum_master', ya_voto: true, valor: 5 },
      { usuario_id: 2, nombre: 'Beta Dos', rol: 'product_builder', ya_voto: false, valor: null },
    ],
    promedio: null,
    diferencias: null,
    ...overrides,
  }
}

function conSegundoVoto(): Votacion {
  return votacion({
    integrantes: [
      { usuario_id: 1, nombre: 'Beta Uno', rol: 'scrum_master', ya_voto: true, valor: 5 },
      { usuario_id: 2, nombre: 'Beta Dos', rol: 'product_builder', ya_voto: true, valor: null },
    ],
  })
}

async function flush() {
  await act(async () => {})
}

async function montar() {
  const hook = renderHook(() => useVotacion('25', 3))
  await flush()
  return hook
}

async function avanzar(ms: number) {
  await act(async () => {
    vi.advanceTimersByTime(ms)
  })
  await flush()
}

beforeEach(() => {
  vi.useFakeTimers()
  visibilidad = 'visible'
  Object.defineProperty(document, 'visibilityState', {
    configurable: true,
    get: () => visibilidad,
  })
  consultarMock.mockResolvedValue(votacion())
  integrantesMock.mockResolvedValue([])
  historiaMock.mockResolvedValue(null)
})

afterEach(() => {
  vi.useRealTimers()
  vi.resetAllMocks()
})

describe('useVotacion - actualización automática', () => {
  it('Escenario: con la ronda abierta consulta la votación cada 5 segundos sin pulsar Actualizar', async () => {
    const { result } = await montar()
    expect(consultarMock).toHaveBeenCalledTimes(1)

    consultarMock.mockResolvedValueOnce(conSegundoVoto())
    await avanzar(INTERVALO)

    expect(consultarMock).toHaveBeenCalledTimes(2)
    expect(result.current.votacion?.integrantes[1].ya_voto).toBe(true)
    expect(result.current.cargando).toBe(false)
  })

  it('Escenario: la consulta automática no muestra el indicador de carga', async () => {
    const { result } = await montar()
    consultarMock.mockReturnValueOnce(new Promise(() => {}))

    await avanzar(INTERVALO)

    expect(consultarMock).toHaveBeenCalledTimes(2)
    expect(result.current.cargando).toBe(false)
    expect(result.current.enviando).toBe(false)
  })

  it('Escenario: con la pestaña oculta no consulta y al volver a visible reanuda', async () => {
    await montar()
    visibilidad = 'hidden'

    await avanzar(INTERVALO * 3)
    expect(consultarMock).toHaveBeenCalledTimes(1)

    visibilidad = 'visible'
    await avanzar(INTERVALO)
    expect(consultarMock).toHaveBeenCalledTimes(2)
  })

  it('Escenario: cuando la consulta automática detecta la revelación deja de consultar', async () => {
    const { result } = await montar()
    const revelada = votacion({
      revelada_en: '2026-11-15T10:00:00Z',
      promedio: { exacto: 5, carta_sugerida: 5 },
      diferencias: { hay_diferencias: false, minimo: null, maximo: null },
    })
    consultarMock.mockResolvedValueOnce(revelada)

    await avanzar(INTERVALO)
    expect(result.current.votacion?.revelada_en).toBe('2026-11-15T10:00:00Z')
    expect(consultarMock).toHaveBeenCalledTimes(2)

    await avanzar(INTERVALO * 3)
    expect(consultarMock).toHaveBeenCalledTimes(2)
  })

  it('Escenario: sin votación iniciada no consulta sola', async () => {
    consultarMock.mockRejectedValue(new ApiError(404, 'la historia no tiene votación iniciada'))

    const { result } = await montar()
    expect(result.current.sinVotacion).toBe(true)

    await avanzar(INTERVALO * 3)
    expect(consultarMock).toHaveBeenCalledTimes(1)
  })

  it('Escenario: con una acción en curso no consulta sola hasta que termina', async () => {
    const { result } = await montar()
    let resolver: (valor: Votacion) => void = () => {}
    votarMock.mockReturnValue(new Promise((resolverVoto) => (resolver = resolverVoto)))

    act(() => {
      void result.current.votar(8)
    })
    await flush()
    expect(result.current.enviando).toBe(true)

    await avanzar(INTERVALO * 2)
    expect(consultarMock).toHaveBeenCalledTimes(1)

    await act(async () => resolver(conSegundoVoto()))
    await flush()
    expect(result.current.enviando).toBe(false)

    await avanzar(INTERVALO)
    expect(consultarMock).toHaveBeenCalledTimes(2)
  })

  it('Escenario: una consulta automática vieja no pisa el resultado de una acción', async () => {
    const { result } = await montar()
    let responderConsulta: (valor: Votacion) => void = () => {}
    consultarMock.mockReturnValueOnce(new Promise((resolver) => (responderConsulta = resolver)))
    await avanzar(INTERVALO)
    expect(consultarMock).toHaveBeenCalledTimes(2)

    votarMock.mockResolvedValue(conSegundoVoto())
    await act(async () => {
      await result.current.votar(8)
    })
    expect(result.current.votacion?.integrantes[1].ya_voto).toBe(true)

    await act(async () => responderConsulta(votacion()))
    await flush()

    expect(result.current.votacion?.integrantes[1].ya_voto).toBe(true)
  })
})
