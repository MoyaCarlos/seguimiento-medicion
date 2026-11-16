import { act, renderHook, waitFor } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { useUsuarioActual } from './useUsuarioActual'

afterEach(() => vi.unstubAllGlobals())

describe('useUsuarioActual', () => {
  it('Given una sesión iniciada, When se monta, Then pide GET /auth/yo una sola vez y devuelve el id del usuario', async () => {
    const fetchSimulado = vi.fn(() => Promise.resolve(new Response(JSON.stringify({ id: 7, nombre: 'ana' }), { status: 200 })))
    vi.stubGlobal('fetch', fetchSimulado)

    const { result, rerender } = renderHook(() => useUsuarioActual())
    expect(result.current).toBeNull()
    await waitFor(() => expect(result.current).toBe(7))
    rerender()

    expect(fetchSimulado).toHaveBeenCalledTimes(1)
    expect(fetchSimulado).toHaveBeenCalledWith('/auth/yo', expect.objectContaining({ method: 'GET' }))
  })

  it.each([
    ['falla', () => new Response(JSON.stringify({ mensaje: 'no' }), { status: 401 })],
    ['no trae un id', () => new Response(JSON.stringify({ nombre: 'ana' }), { status: 200 })],
  ])('Given que el pedido %s, When se monta, Then no hay usuario y no se rompe', async (_caso, respuesta) => {
    const fetchSimulado = vi.fn(() => Promise.resolve(respuesta()))
    vi.stubGlobal('fetch', fetchSimulado)

    const { result } = renderHook(() => useUsuarioActual())
    await waitFor(() => expect(fetchSimulado).toHaveBeenCalled())
    await act(() => new Promise((resolver) => setTimeout(resolver, 0)))

    expect(result.current).toBeNull()
  })
})
