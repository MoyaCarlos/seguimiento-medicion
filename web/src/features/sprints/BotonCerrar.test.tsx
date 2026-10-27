import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { BotonCerrar } from './BotonCerrar'
import type { Sprint } from './sprint'

const activo: Sprint = { id: 7, proyecto_id: 1, estado: 'Activo' }
const finalizado: Sprint = { ...activo, estado: 'Finalizado' }

const respuesta = (status: number, cuerpo: unknown) =>
  new Response(JSON.stringify(cuerpo), { status })

const cerrar = () => userEvent.click(screen.getByRole('button', { name: 'Cerrar Sprint' }))

afterEach(() => vi.unstubAllGlobals())

describe('BotonCerrar', () => {
  it('hace POST /sprints/{id}/cerrar sin cuerpo y entrega el Sprint devuelto', async () => {
    const fetchSimulado = vi.fn().mockResolvedValue(respuesta(200, finalizado))
    vi.stubGlobal('fetch', fetchSimulado)
    const alActualizar = vi.fn()
    render(<BotonCerrar sprint={activo} alActualizar={alActualizar} />)

    await cerrar()

    await waitFor(() => expect(alActualizar).toHaveBeenCalledWith(finalizado))
    const [ruta, opciones] = fetchSimulado.mock.calls[0]
    expect(ruta).toBe('/sprints/7/cerrar')
    expect(opciones.method).toBe('POST')
    expect(opciones.body).toBeUndefined()
  })

  it('deshabilita el botón y muestra la carga mientras envía', async () => {
    vi.stubGlobal('fetch', vi.fn().mockReturnValue(new Promise(() => {})))
    render(<BotonCerrar sprint={activo} alActualizar={vi.fn()} />)

    await cerrar()

    expect(screen.getByText('Cargando…')).toBeTruthy()
    expect((screen.getByRole('button', { name: 'Cerrar Sprint' }) as HTMLButtonElement).disabled).toBe(true)
  })

  it('muestra el error de la API en un aviso y vuelve a habilitar el botón', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue(respuesta(400, { campo: 'estado', mensaje: 'solo se cierran Sprints activos' })),
    )
    render(<BotonCerrar sprint={activo} alActualizar={vi.fn()} />)

    await cerrar()

    expect((await screen.findByRole('alert')).textContent).toBe('solo se cierran Sprints activos')
    expect((screen.getByRole('button', { name: 'Cerrar Sprint' }) as HTMLButtonElement).disabled).toBe(false)
  })

  it('muestra el error de conexión en el aviso', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new TypeError('red')))
    render(<BotonCerrar sprint={activo} alActualizar={vi.fn()} />)

    await cerrar()

    expect((await screen.findByRole('alert')).textContent).toBe('No se pudo conectar con el servidor')
  })

  it('un nuevo clic borra el aviso anterior', async () => {
    vi.stubGlobal(
      'fetch',
      vi
        .fn()
        .mockResolvedValueOnce(respuesta(400, { mensaje: 'primer error' }))
        .mockResolvedValueOnce(respuesta(400, { mensaje: 'segundo error' })),
    )
    render(<BotonCerrar sprint={activo} alActualizar={vi.fn()} />)

    await cerrar()
    await screen.findByText('primer error')
    await cerrar()

    expect(await screen.findByText('segundo error')).toBeTruthy()
    expect(screen.queryByText('primer error')).toBeNull()
  })
})
