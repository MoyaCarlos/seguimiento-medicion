import { act, renderHook, waitFor } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { usePermisosProyecto } from './usePermisosProyecto'

function simularRespuesta(cuerpo: unknown, status = 200) {
  const fetchSimulado = vi.fn(() => Promise.resolve(new Response(JSON.stringify(cuerpo), { status })))
  vi.stubGlobal('fetch', fetchSimulado)
  return fetchSimulado
}

function respuesta(cuerpo: unknown) {
  return Promise.resolve(new Response(JSON.stringify(cuerpo), { status: 200 }))
}

afterEach(() => {
  vi.unstubAllGlobals()
  vi.restoreAllMocks()
})

const sinFunciones = (permisos: ReturnType<typeof usePermisosProyecto>) => {
  const { aplicarCierre, ...datos } = permisos
  void aplicarCierre
  return datos
}

const TODO_HABILITADO_SIN_DATOS = {
  rol: null, fechaCierre: null, cerrado: false, puedeGestionar: true, puedeEscribir: true,
}

describe('usePermisosProyecto', () => {
  it('Scrum Master con el proyecto abierto puede gestionar y escribir', async () => {
    const fetchSimulado = simularRespuesta({ id: 1, rol: 'scrum_master', fecha_cierre: null })

    const { result } = renderHook(() => usePermisosProyecto(1))

    await waitFor(() => expect(result.current.rol).toBe('scrum_master'))
    expect(sinFunciones(result.current)).toEqual({
      rol: 'scrum_master', fechaCierre: null, cerrado: false, puedeGestionar: true, puedeEscribir: true,
    })
    expect(fetchSimulado).toHaveBeenCalledTimes(1)
  })

  it('Product Builder con el proyecto abierto no puede gestionar pero sí escribir', async () => {
    const fetchSimulado = simularRespuesta({ id: 1, rol: 'product_builder', fecha_cierre: null })

    const { result } = renderHook(() => usePermisosProyecto(1))

    await waitFor(() => expect(result.current.rol).toBe('product_builder'))
    expect(result.current.puedeGestionar).toBe(false)
    expect(result.current.puedeEscribir).toBe(true)
    expect(fetchSimulado).toHaveBeenCalledTimes(1)
  })

  it('Scrum Master con el proyecto cerrado no puede gestionar ni escribir', async () => {
    const fetchSimulado = simularRespuesta({ id: 1, rol: 'scrum_master', fecha_cierre: '2026-11-11' })

    const { result } = renderHook(() => usePermisosProyecto(1))

    await waitFor(() => expect(result.current.cerrado).toBe(true))
    expect(result.current.fechaCierre).toBe('2026-11-11')
    expect(result.current.puedeGestionar).toBe(false)
    expect(result.current.puedeEscribir).toBe(false)
    expect(fetchSimulado).toHaveBeenCalledTimes(1)
  })

  it('Product Builder con el proyecto cerrado no puede gestionar ni escribir', async () => {
    const fetchSimulado = simularRespuesta({ id: 1, rol: 'product_builder', fecha_cierre: '2026-11-11' })

    const { result } = renderHook(() => usePermisosProyecto(1))

    await waitFor(() => expect(result.current.cerrado).toBe(true))
    expect(result.current.puedeGestionar).toBe(false)
    expect(result.current.puedeEscribir).toBe(false)
    expect(fetchSimulado).toHaveBeenCalledTimes(1)
  })

  it('mientras carga deja todo habilitado y sin rol ni cierre', () => {
    const fetchSimulado = vi.fn(() => new Promise<Response>(() => {}))
    vi.stubGlobal('fetch', fetchSimulado)

    const { result } = renderHook(() => usePermisosProyecto(1))

    expect(sinFunciones(result.current)).toEqual(TODO_HABILITADO_SIN_DATOS)
    expect(fetchSimulado).toHaveBeenCalledTimes(1)
  })

  it('ante un error de conexión deja todo habilitado y sin cierre', async () => {
    const fetchSimulado = vi.fn(() => Promise.reject(new TypeError('red caída')))
    vi.stubGlobal('fetch', fetchSimulado)

    const { result } = renderHook(() => usePermisosProyecto(1))

    await waitFor(() => expect(fetchSimulado).toHaveBeenCalledTimes(1))
    await new Promise((resolver) => setTimeout(resolver, 20))
    expect(sinFunciones(result.current)).toEqual(TODO_HABILITADO_SIN_DATOS)
    expect(fetchSimulado).toHaveBeenCalledTimes(1)
  })

  it('ante un 500 deja todo habilitado y sin cierre', async () => {
    const fetchSimulado = simularRespuesta({ mensaje: 'Falló' }, 500)

    const { result } = renderHook(() => usePermisosProyecto(1))

    await waitFor(() => expect(fetchSimulado).toHaveBeenCalledTimes(1))
    await new Promise((resolver) => setTimeout(resolver, 20))
    expect(sinFunciones(result.current)).toEqual(TODO_HABILITADO_SIN_DATOS)
    expect(fetchSimulado).toHaveBeenCalledTimes(1)
  })

  it('si no integra el proyecto no puede gestionar, pero puede escribir y no figura cerrado', async () => {
    const fetchSimulado = simularRespuesta({ mensaje: 'No sos integrante de este proyecto' }, 403)

    const { result } = renderHook(() => usePermisosProyecto(1))

    await waitFor(() => expect(result.current.puedeGestionar).toBe(false))
    expect(result.current.puedeEscribir).toBe(true)
    expect(result.current.cerrado).toBe(false)
    expect(result.current.rol).toBeNull()
    expect(fetchSimulado).toHaveBeenCalledTimes(1)
  })

  it.each([
    ['fecha_cierre null', { id: 1, rol: 'scrum_master', fecha_cierre: null }],
    ['sin el campo fecha_cierre', { id: 1, rol: 'scrum_master' }],
  ])('con el proyecto abierto (%s) no figura cerrado', async (_nombre, cuerpo) => {
    const fetchSimulado = simularRespuesta(cuerpo)

    const { result } = renderHook(() => usePermisosProyecto(1))

    await waitFor(() => expect(result.current.rol).toBe('scrum_master'))
    expect(result.current.cerrado).toBe(false)
    expect(result.current.fechaCierre).toBeNull()
    expect(fetchSimulado).toHaveBeenCalledTimes(1)
  })

  it('aplicarCierre con una fecha pasa a cerrado sin otro pedido y con null lo reabre', async () => {
    const fetchSimulado = simularRespuesta({ id: 1, rol: 'scrum_master', fecha_cierre: null })
    const { result } = renderHook(() => usePermisosProyecto(1))
    await waitFor(() => expect(result.current.rol).toBe('scrum_master'))

    act(() => result.current.aplicarCierre('2026-11-12'))

    expect(result.current.cerrado).toBe(true)
    expect(result.current.fechaCierre).toBe('2026-11-12')
    expect(result.current.puedeGestionar).toBe(false)
    expect(result.current.puedeEscribir).toBe(false)

    act(() => result.current.aplicarCierre(null))

    expect(result.current.cerrado).toBe(false)
    expect(result.current.fechaCierre).toBeNull()
    expect(result.current.puedeGestionar).toBe(true)
    expect(result.current.puedeEscribir).toBe(true)
    expect(fetchSimulado).toHaveBeenCalledTimes(1)
  })

  it('al cambiar de proyecto no arrastra el cierre aplicado al proyecto anterior', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn((ruta: string) =>
        ruta === '/projects/1'
          ? respuesta({ id: 1, rol: 'scrum_master', fecha_cierre: null })
          : new Promise<Response>(() => {}),
      ),
    )
    const { result, rerender } = renderHook(({ id }) => usePermisosProyecto(id), { initialProps: { id: 1 } })
    await waitFor(() => expect(result.current.rol).toBe('scrum_master'))
    act(() => result.current.aplicarCierre('2026-11-12'))
    expect(result.current.cerrado).toBe(true)

    rerender({ id: 2 })

    expect(result.current.cerrado).toBe(false)
    expect(result.current.fechaCierre).toBeNull()
  })

  it('al cambiar de proyecto no arrastra el cierre leído del proyecto anterior', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn((ruta: string) =>
        ruta === '/projects/1'
          ? respuesta({ id: 1, rol: 'scrum_master', fecha_cierre: '2026-11-11' })
          : new Promise<Response>(() => {}),
      ),
    )
    const { result, rerender } = renderHook(({ id }) => usePermisosProyecto(id), { initialProps: { id: 1 } })
    await waitFor(() => expect(result.current.cerrado).toBe(true))

    rerender({ id: 2 })

    expect(result.current.cerrado).toBe(false)
    expect(result.current.puedeEscribir).toBe(true)
  })
})
