import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { ApiError } from '../../compartido/api/client'
import { EstadoHistoria } from './EstadoHistoria'
import { cambiarEstadoHistoria } from './api'
import type { HistoriaDeSprint } from './sprint'

vi.mock('./api', () => ({ cambiarEstadoHistoria: vi.fn() }))

const NUEVA: HistoriaDeSprint = {
  id: 12,
  titulo: 'Registrar horas',
  prioridad: 'M',
  estado: 'Nueva',
  estimacion_sp: 5,
}
const EN_PROGRESO: HistoriaDeSprint = { ...NUEVA, estado: 'En progreso' }
const COMPLETADA: HistoriaDeSprint = {
  ...NUEVA,
  estado: 'Completada',
  fecha_finalizacion: '2026-10-29',
}

const cambiarMock = vi.mocked(cambiarEstadoHistoria)

function abrir(historia: HistoriaDeSprint) {
  const alCambiar = vi.fn()
  render(<EstadoHistoria historia={historia} alCambiar={alCambiar} />)
  return { alCambiar }
}

afterEach(() => {
  vi.resetAllMocks()
})

describe('EstadoHistoria', () => {
  it('US3.1: con estado Nueva ofrece "Iniciar" y llama con "En progreso"', async () => {
    cambiarMock.mockResolvedValue(EN_PROGRESO)
    abrir(NUEVA)

    await userEvent.click(screen.getByRole('button', { name: 'Iniciar' }))

    expect(cambiarMock).toHaveBeenCalledWith(12, 'En progreso')
  })

  it('US3.2: con estado En progreso ofrece "Completar" y llama con "Completada"', async () => {
    cambiarMock.mockResolvedValue(COMPLETADA)
    abrir(EN_PROGRESO)

    await userEvent.click(screen.getByRole('button', { name: 'Completar' }))

    expect(cambiarMock).toHaveBeenCalledWith(12, 'Completada')
    expect(screen.queryByRole('button', { name: 'Iniciar' })).toBeNull()
  })

  it('US3.3: con estado Completada no ofrece ninguna acción', () => {
    abrir(COMPLETADA)

    expect(screen.queryByRole('button')).toBeNull()
  })

  it('US3.4: el botón queda deshabilitado mientras la petición está en curso y no envía dos veces', async () => {
    cambiarMock.mockReturnValue(new Promise(() => {}))
    abrir(NUEVA)

    await userEvent.click(screen.getByRole('button', { name: 'Iniciar' }))
    await userEvent.click(screen.getByRole('button', { name: 'Iniciar' }))

    expect((screen.getByRole('button', { name: 'Iniciar' }) as HTMLButtonElement).disabled).toBe(true)
    expect(cambiarMock).toHaveBeenCalledTimes(1)
  })

  it('US3.5: tras un 200 llama a alCambiar con la historia devuelta', async () => {
    cambiarMock.mockResolvedValue(EN_PROGRESO)
    const { alCambiar } = abrir(NUEVA)

    await userEvent.click(screen.getByRole('button', { name: 'Iniciar' }))

    await waitFor(() => expect(alCambiar).toHaveBeenCalledWith(EN_PROGRESO))
  })

  it('US3.6: ante un 400 muestra el mensaje de la API y no llama a alCambiar', async () => {
    cambiarMock.mockRejectedValue(
      new ApiError(400, 'solo se cambia el estado de historias de un Sprint activo', 'estado'),
    )
    const { alCambiar } = abrir(NUEVA)

    await userEvent.click(screen.getByRole('button', { name: 'Iniciar' }))

    expect((await screen.findByRole('alert')).textContent).toBe(
      'solo se cambia el estado de historias de un Sprint activo',
    )
    expect(alCambiar).not.toHaveBeenCalled()
    expect((screen.getByRole('button', { name: 'Iniciar' }) as HTMLButtonElement).disabled).toBe(false)
  })

  it('US3.7: ante un 409 muestra el mensaje de la API tal cual y no llama a alCambiar', async () => {
    const mensaje = 'la historia cambió de estado mientras tanto: recargá y volvé a intentar'
    cambiarMock.mockRejectedValue(new ApiError(409, mensaje))
    const { alCambiar } = abrir(EN_PROGRESO)

    await userEvent.click(screen.getByRole('button', { name: 'Completar' }))

    expect((await screen.findByRole('alert')).textContent).toBe(mensaje)
    expect(alCambiar).not.toHaveBeenCalled()
  })
})
