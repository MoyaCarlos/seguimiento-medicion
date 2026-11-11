import { renderHook, waitFor } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { useProyectos } from './useProyectos'

function simularRespuesta(cuerpo: string, status = 200) {
  const fetchSimulado = vi.fn(() => Promise.resolve(new Response(cuerpo, { status })))
  vi.stubGlobal('fetch', fetchSimulado)
  return fetchSimulado
}

afterEach(() => {
  vi.unstubAllGlobals()
  vi.restoreAllMocks()
})

describe('useProyectos', () => {
  it('arranca en estado cargando y pide GET /projects una sola vez', () => {
    const fetchSimulado = simularRespuesta('[]')

    const { result } = renderHook(() => useProyectos())

    expect(result.current).toEqual({ estado: 'cargando' })
    expect(fetchSimulado).toHaveBeenCalledTimes(1)
    expect(fetchSimulado).toHaveBeenCalledWith('/projects', expect.objectContaining({ method: 'GET' }))
  })

  it('pasa a listo conservando el orden de la API', async () => {
    simularRespuesta(JSON.stringify([{ id: 3, nombre: 'Gamma' }, { id: 1, nombre: 'Alfa' }]))

    const { result } = renderHook(() => useProyectos())

    await waitFor(() => expect(result.current.estado).toBe('listo'))
    expect(result.current).toEqual({
      estado: 'listo',
      proyectos: [{ id: 3, nombre: 'Gamma' }, { id: 1, nombre: 'Alfa' }],
    })
  })

  it('pasa a error con el mensaje de la API', async () => {
    simularRespuesta(JSON.stringify({ mensaje: 'Falló el listado' }), 500)

    const { result } = renderHook(() => useProyectos())

    await waitFor(() => expect(result.current.estado).toBe('error'))
    expect(result.current).toEqual({ estado: 'error', mensaje: 'Falló el listado' })
  })

  it('pasa a error cuando no hay conexión', async () => {
    vi.stubGlobal('fetch', vi.fn(() => Promise.reject(new TypeError('red caída'))))

    const { result } = renderHook(() => useProyectos())

    await waitFor(() => expect(result.current.estado).toBe('error'))
    expect(result.current).toEqual({ estado: 'error', mensaje: 'No se pudo conectar con el servidor' })
  })

  it('trata una respuesta 200 que no es un arreglo como error', async () => {
    simularRespuesta(JSON.stringify({ proyectos: [] }))

    const { result } = renderHook(() => useProyectos())

    await waitFor(() => expect(result.current.estado).toBe('error'))
    expect(result.current).toEqual({
      estado: 'error',
      mensaje: 'No se pudo interpretar la lista de proyectos',
    })
  })

  it('trata una respuesta 200 con HTML como lista no interpretable', async () => {
    simularRespuesta('<!doctype html><html>index</html>')

    const { result } = renderHook(() => useProyectos())

    await waitFor(() => expect(result.current.estado).toBe('error'))
    expect(result.current).toEqual({
      estado: 'error',
      mensaje: 'No se pudo interpretar la lista de proyectos',
    })
  })

  it('no falla si se desmonta antes de que responda la API', async () => {
    const consola = vi.spyOn(console, 'error').mockImplementation(() => {})
    let responder: (respuesta: Response) => void = () => {}
    vi.stubGlobal('fetch', vi.fn(() => new Promise<Response>((resolver) => { responder = resolver })))

    const { unmount } = renderHook(() => useProyectos())
    unmount()
    responder(new Response('[]', { status: 200 }))
    await Promise.resolve()

    expect(consola).not.toHaveBeenCalled()
  })
})
