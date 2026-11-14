import { renderHook, waitFor } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { useRolEnProyecto } from './useRolEnProyecto'

function simularRespuesta(cuerpo: unknown, status = 200) {
  const fetchSimulado = vi.fn(() => Promise.resolve(new Response(JSON.stringify(cuerpo), { status })))
  vi.stubGlobal('fetch', fetchSimulado)
  return fetchSimulado
}

afterEach(() => {
  vi.unstubAllGlobals()
  vi.restoreAllMocks()
})

describe('useRolEnProyecto: fecha de cierre', () => {
  it('devuelve listo con la fecha de cierre cuando el proyecto está cerrado', async () => {
    const fetchSimulado = simularRespuesta({ id: 1, rol: 'scrum_master', fecha_cierre: '2026-11-11' })

    const { result } = renderHook(() => useRolEnProyecto(1))

    await waitFor(() => expect(result.current.estado).toBe('listo'))
    expect(result.current).toEqual({ estado: 'listo', rol: 'scrum_master', fechaCierre: '2026-11-11' })
    expect(fetchSimulado).toHaveBeenCalledTimes(1)
  })

  it('devuelve fechaCierre null cuando el proyecto está abierto (fecha_cierre null)', async () => {
    simularRespuesta({ id: 1, rol: 'scrum_master', fecha_cierre: null })

    const { result } = renderHook(() => useRolEnProyecto(1))

    await waitFor(() => expect(result.current.estado).toBe('listo'))
    expect(result.current).toEqual({ estado: 'listo', rol: 'scrum_master', fechaCierre: null })
  })

  it('devuelve fechaCierre null cuando la respuesta no trae el campo fecha_cierre', async () => {
    simularRespuesta({ id: 1, rol: 'product_builder' })

    const { result } = renderHook(() => useRolEnProyecto(1))

    await waitFor(() => expect(result.current.estado).toBe('listo'))
    expect(result.current).toEqual({ estado: 'listo', rol: 'product_builder', fechaCierre: null })
  })

  it('devuelve error cuando fecha_cierre no es texto ni null', async () => {
    simularRespuesta({ id: 1, rol: 'scrum_master', fecha_cierre: 5 })

    const { result } = renderHook(() => useRolEnProyecto(1))

    await waitFor(() => expect(result.current.estado).toBe('error'))
    expect(result.current).toEqual({ estado: 'error', mensaje: 'No se pudo interpretar el rol en el proyecto' })
  })
})
