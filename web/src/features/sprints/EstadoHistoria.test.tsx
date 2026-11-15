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

async function pedir(etiqueta: 'Iniciar' | 'Completar') {
  await userEvent.click(screen.getByRole('button', { name: etiqueta }))
}

async function confirmar(etiqueta: 'Iniciar' | 'Completar') {
  await pedir(etiqueta)
  await userEvent.click(screen.getByRole('button', { name: 'Confirmar' }))
}

afterEach(() => {
  vi.resetAllMocks()
})

describe('EstadoHistoria', () => {
  it('US3.1: con estado Nueva, «Iniciar» pide confirmación y «Confirmar» llama con «En progreso»', async () => {
    cambiarMock.mockResolvedValue(EN_PROGRESO)
    abrir(NUEVA)

    await pedir('Iniciar')

    expect(screen.getByText('¿Iniciar la historia?')).toBeTruthy()
    expect(cambiarMock).not.toHaveBeenCalled()

    await userEvent.click(screen.getByRole('button', { name: 'Confirmar' }))

    expect(cambiarMock).toHaveBeenCalledWith(12, 'En progreso')
  })

  it('US3.2: con estado En progreso, «Completar» pide confirmación y «Confirmar» llama con «Completada»', async () => {
    cambiarMock.mockResolvedValue(COMPLETADA)
    abrir(EN_PROGRESO)

    await confirmar('Completar')

    expect(cambiarMock).toHaveBeenCalledWith(12, 'Completada')
    expect(screen.queryByRole('button', { name: 'Iniciar' })).toBeNull()
  })

  it('US3.3: con estado Completada no ofrece ninguna acción', () => {
    abrir(COMPLETADA)

    expect(screen.queryByRole('button')).toBeNull()
  })

  it('US3.4: con el pedido en curso los botones quedan deshabilitados y no se envía dos veces', async () => {
    cambiarMock.mockReturnValue(new Promise(() => {}))
    abrir(NUEVA)

    await pedir('Iniciar')
    await userEvent.click(screen.getByRole('button', { name: 'Confirmar' }))
    await userEvent.click(screen.getByRole('button', { name: 'Confirmar' }))

    expect((screen.getByRole('button', { name: 'Confirmar' }) as HTMLButtonElement).disabled).toBe(true)
    expect((screen.getByRole('button', { name: 'Cancelar' }) as HTMLButtonElement).disabled).toBe(true)
    expect(cambiarMock).toHaveBeenCalledTimes(1)
  })

  it('US3.5: tras un 200 llama a alCambiar con la historia devuelta', async () => {
    cambiarMock.mockResolvedValue(EN_PROGRESO)
    const { alCambiar } = abrir(NUEVA)

    await confirmar('Iniciar')

    await waitFor(() => expect(alCambiar).toHaveBeenCalledWith(EN_PROGRESO))
  })

  it('US3.6: ante un 400 muestra el mensaje de la API y no llama a alCambiar', async () => {
    cambiarMock.mockRejectedValue(
      new ApiError(400, 'solo se cambia el estado de historias de un Sprint activo', 'estado'),
    )
    const { alCambiar } = abrir(NUEVA)

    await confirmar('Iniciar')

    expect((await screen.findByRole('alert')).textContent).toBe(
      'solo se cambia el estado de historias de un Sprint activo',
    )
    expect(alCambiar).not.toHaveBeenCalled()
    expect(screen.queryByText('¿Iniciar la historia?')).toBeNull()
    expect((screen.getByRole('button', { name: 'Iniciar' }) as HTMLButtonElement).disabled).toBe(false)
  })

  it('US3.7: ante un 409 muestra el mensaje de la API tal cual y no llama a alCambiar', async () => {
    const mensaje = 'la historia cambió de estado mientras tanto: recargá y volvé a intentar'
    cambiarMock.mockRejectedValue(new ApiError(409, mensaje))
    const { alCambiar } = abrir(EN_PROGRESO)

    await confirmar('Completar')

    expect((await screen.findByRole('alert')).textContent).toBe(mensaje)
    expect(alCambiar).not.toHaveBeenCalled()
  })

  it('Defensa: un reintento exitoso tras un error borra el mensaje viejo', async () => {
    cambiarMock.mockRejectedValueOnce(new ApiError(409, 'la historia cambió de estado mientras tanto'))
    cambiarMock.mockResolvedValueOnce(EN_PROGRESO)
    const { alCambiar } = abrir(NUEVA)

    await confirmar('Iniciar')
    expect(await screen.findByRole('alert')).toBeTruthy()

    await confirmar('Iniciar')

    await waitFor(() => expect(alCambiar).toHaveBeenCalledWith(EN_PROGRESO))
    expect(screen.queryByRole('alert')).toBeNull()
  })

  it('US3.8: ante un 404 historia no encontrada muestra el mensaje de la API tal cual y no llama a alCambiar', async () => {
    const mensaje = 'historia no encontrada'
    cambiarMock.mockRejectedValue(new ApiError(404, mensaje))
    const { alCambiar } = abrir(NUEVA)

    await confirmar('Iniciar')

    expect((await screen.findByRole('alert')).textContent).toBe(mensaje)
    expect(alCambiar).not.toHaveBeenCalled()
  })

  it('Escenario: al pedir confirmación el foco queda en «Cancelar»', async () => {
    abrir(NUEVA)

    await pedir('Iniciar')

    expect(document.activeElement).toBe(screen.getByRole('button', { name: 'Cancelar' }))
  })

  it('Escenario: «Cancelar» cierra la confirmación sin enviar nada', async () => {
    abrir(NUEVA)

    await pedir('Iniciar')
    await userEvent.click(screen.getByRole('button', { name: 'Cancelar' }))

    expect(screen.queryByText('¿Iniciar la historia?')).toBeNull()
    expect(screen.getByRole('button', { name: 'Iniciar' })).toBeTruthy()
    expect(cambiarMock).not.toHaveBeenCalled()
  })

  it('Escenario: Escape cierra la confirmación sin enviar nada', async () => {
    abrir(NUEVA)

    await pedir('Iniciar')
    await userEvent.keyboard('{Escape}')

    expect(screen.queryByText('¿Iniciar la historia?')).toBeNull()
    expect(cambiarMock).not.toHaveBeenCalled()
  })

  it('Escenario: un clic fuera de la confirmación no envía nada', async () => {
    abrir(NUEVA)

    await pedir('Iniciar')
    await userEvent.click(document.body)

    expect(cambiarMock).not.toHaveBeenCalled()
  })
})

describe('EstadoHistoria - proyecto cerrado (cerrar y reabrir, US2)', () => {
  const boton = () => screen.getByRole('button', { name: 'Iniciar' }) as HTMLButtonElement
  const nota = () => document.getElementById(boton().getAttribute('aria-describedby') ?? '')?.textContent

  it('Escenario: proyecto cerrado: con puedeEscribir en false «Iniciar» queda deshabilitado con la nota «El proyecto está cerrado» y no envía', async () => {
    cambiarMock.mockResolvedValue(EN_PROGRESO)
    const alCambiar = vi.fn()
    render(<EstadoHistoria historia={NUEVA} alCambiar={alCambiar} puedeEscribir={false} />)

    await userEvent.click(boton())

    expect(boton().disabled).toBe(true)
    expect(nota()).toBe('El proyecto está cerrado')
    expect(screen.queryByText('¿Iniciar la historia?')).toBeNull()
    expect(cambiarMock).not.toHaveBeenCalled()
    expect(alCambiar).not.toHaveBeenCalled()
  })

  it('Escenario: proyecto abierto: con el valor por defecto «Iniciar» está habilitado y sin nota', () => {
    abrir(NUEVA)

    expect(boton().disabled).toBe(false)
    expect(boton().hasAttribute('aria-describedby')).toBe(false)
    expect(screen.queryByText('El proyecto está cerrado')).toBeNull()
  })
})
