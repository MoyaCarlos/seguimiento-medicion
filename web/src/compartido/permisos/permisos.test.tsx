import { render, renderHook, screen, waitFor } from '@testing-library/react'
import { MemoryRouter } from 'react-router'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { AvisoSinPermiso } from './AvisoSinPermiso'
import { NoIntegrante } from './NoIntegrante'
import { useEsScrumMaster } from './useEsScrumMaster'
import { useRolEnProyecto } from './useRolEnProyecto'

function simularRespuesta(cuerpo: unknown, status = 200) {
  const fetchSimulado = vi.fn(() => Promise.resolve(new Response(JSON.stringify(cuerpo), { status })))
  vi.stubGlobal('fetch', fetchSimulado)
  return fetchSimulado
}

const NO_INTEGRANTE = { mensaje: 'No sos integrante de este proyecto' }

afterEach(() => {
  vi.unstubAllGlobals()
  vi.restoreAllMocks()
})

describe('NoIntegrante', () => {
  it('muestra «No sos integrante de este proyecto» con un enlace a «Mis proyectos»', () => {
    render(
      <MemoryRouter>
        <NoIntegrante />
      </MemoryRouter>,
    )

    expect(screen.getByText('No sos integrante de este proyecto')).toBeTruthy()
    const enlace = screen.getByRole('link', { name: 'Mis proyectos' })
    expect(enlace.getAttribute('href')).toBe('/')
  })
})

describe('AvisoSinPermiso', () => {
  it('muestra «No tenés permiso» y el mensaje del servidor', () => {
    render(<AvisoSinPermiso mensaje="Solo el Scrum Master puede realizar esta acción" />)

    expect(screen.getByText('No tenés permiso')).toBeTruthy()
    expect(screen.getByText('Solo el Scrum Master puede realizar esta acción')).toBeTruthy()
  })
})

describe('useRolEnProyecto', () => {
  it('arranca en cargando y pide GET /projects/{id} una sola vez', () => {
    const fetchSimulado = simularRespuesta({ id: 7, rol: 'scrum_master' })

    const { result } = renderHook(() => useRolEnProyecto(7))

    expect(result.current).toEqual({ estado: 'cargando' })
    expect(fetchSimulado).toHaveBeenCalledTimes(1)
    expect(fetchSimulado).toHaveBeenCalledWith('/projects/7', expect.objectContaining({ method: 'GET' }))
  })

  it.each(['scrum_master', 'product_builder'] as const)(
    'devuelve listo con el rol %s cuando el servidor responde con ese rol',
    async (rol) => {
      simularRespuesta({ id: 7, rol })

      const { result } = renderHook(() => useRolEnProyecto(7))

      await waitFor(() => expect(result.current.estado).toBe('listo'))
      expect(result.current).toEqual({ estado: 'listo', rol })
    },
  )

  it('devuelve no_integrante ante el 403 «No sos integrante de este proyecto»', async () => {
    simularRespuesta(NO_INTEGRANTE, 403)

    const { result } = renderHook(() => useRolEnProyecto(7))

    await waitFor(() => expect(result.current.estado).toBe('no_integrante'))
    expect(result.current).toEqual({ estado: 'no_integrante' })
  })

  it('devuelve error con el mensaje del servidor ante otro 403', async () => {
    simularRespuesta({ mensaje: 'Acceso denegado' }, 403)

    const { result } = renderHook(() => useRolEnProyecto(7))

    await waitFor(() => expect(result.current.estado).toBe('error'))
    expect(result.current).toEqual({ estado: 'error', mensaje: 'Acceso denegado' })
  })

  it('devuelve error con el mensaje de conexión cuando no hay respuesta', async () => {
    vi.stubGlobal('fetch', vi.fn(() => Promise.reject(new TypeError('red caída'))))

    const { result } = renderHook(() => useRolEnProyecto(7))

    await waitFor(() => expect(result.current.estado).toBe('error'))
    expect(result.current).toEqual({ estado: 'error', mensaje: 'No se pudo conectar con el servidor' })
  })

  it('devuelve error si el rol de la respuesta no se reconoce', async () => {
    simularRespuesta({ id: 7, rol: 'visitante' })

    const { result } = renderHook(() => useRolEnProyecto(7))

    await waitFor(() => expect(result.current.estado).toBe('error'))
  })

  it('ignora la respuesta vieja de un proyecto anterior al cambiar de proyecto', async () => {
    const respuestas: Record<string, (respuesta: Response) => void> = {}
    vi.stubGlobal(
      'fetch',
      vi.fn((ruta: string) => new Promise<Response>((resolver) => { respuestas[ruta] = resolver })),
    )

    const { result, rerender } = renderHook(({ id }) => useRolEnProyecto(id), {
      initialProps: { id: 1 },
    })
    rerender({ id: 2 })
    respuestas['/projects/2'](new Response(JSON.stringify({ id: 2, rol: 'product_builder' }), { status: 200 }))
    await waitFor(() => expect(result.current).toEqual({ estado: 'listo', rol: 'product_builder' }))
    respuestas['/projects/1'](new Response(JSON.stringify({ id: 1, rol: 'scrum_master' }), { status: 200 }))
    await Promise.resolve()
    await Promise.resolve()

    expect(result.current).toEqual({ estado: 'listo', rol: 'product_builder' })
  })

  it('al cambiar de proyecto vuelve a cargando hasta que responde el nuevo', async () => {
    simularRespuesta({ id: 1, rol: 'scrum_master' })
    const { result, rerender } = renderHook(({ id }) => useRolEnProyecto(id), {
      initialProps: { id: 1 },
    })
    await waitFor(() => expect(result.current.estado).toBe('listo'))

    vi.stubGlobal('fetch', vi.fn(() => new Promise<Response>(() => {})))
    rerender({ id: 2 })

    await waitFor(() => expect(result.current).toEqual({ estado: 'cargando' }))
  })

  it('no falla si se desmonta antes de que responda el servidor', async () => {
    const consola = vi.spyOn(console, 'error').mockImplementation(() => {})
    let responder: (respuesta: Response) => void = () => {}
    vi.stubGlobal('fetch', vi.fn(() => new Promise<Response>((resolver) => { responder = resolver })))

    const { unmount } = renderHook(() => useRolEnProyecto(7))
    unmount()
    responder(new Response(JSON.stringify({ id: 7, rol: 'scrum_master' }), { status: 200 }))
    await Promise.resolve()

    expect(consola).not.toHaveBeenCalled()
  })
})

describe('useEsScrumMaster', () => {
  it('devuelve true mientras se carga el rol', () => {
    vi.stubGlobal('fetch', vi.fn(() => new Promise<Response>(() => {})))

    const { result } = renderHook(() => useEsScrumMaster(7))

    expect(result.current).toBe(true)
  })

  it('devuelve true para el Scrum Master', async () => {
    const fetchSimulado = simularRespuesta({ id: 7, rol: 'scrum_master' })

    const { result } = renderHook(() => useEsScrumMaster(7))

    await waitFor(() => expect(fetchSimulado).toHaveBeenCalled())
    await Promise.resolve()
    expect(result.current).toBe(true)
  })

  it('devuelve false para el Product Builder', async () => {
    simularRespuesta({ id: 7, rol: 'product_builder' })

    const { result } = renderHook(() => useEsScrumMaster(7))

    await waitFor(() => expect(result.current).toBe(false))
  })

  it('devuelve false cuando la persona no integra el proyecto', async () => {
    simularRespuesta(NO_INTEGRANTE, 403)

    const { result } = renderHook(() => useEsScrumMaster(7))

    await waitFor(() => expect(result.current).toBe(false))
  })

  it('devuelve true cuando falla la consulta del rol, porque decide el 403 del servidor', async () => {
    simularRespuesta({ mensaje: 'Falló' }, 500)

    const { result } = renderHook(() => useEsScrumMaster(7))

    await new Promise((resolver) => setTimeout(resolver, 20))
    expect(result.current).toBe(true)
  })
})
