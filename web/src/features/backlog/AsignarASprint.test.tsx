import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { ApiError } from '../../compartido/api/client'
import { AsignarASprint } from './AsignarASprint'
import { asignarASprint } from './api'
import type { SprintParaAsignar } from './api'
import type { Historia } from './historia'

vi.mock('./api', () => ({ asignarASprint: vi.fn() }))

const HISTORIA: Historia = {
  id: 7,
  proyecto_id: 42,
  titulo: 'Como usuario quiero iniciar sesión',
  descripcion: 'Autenticación con email y contraseña',
  criterios_aceptacion: null,
  prioridad: 'M',
  estado: 'Nueva',
  valor_negocio: null,
  estimacion_sp: null,
  sprint_id: null,
}

const SPRINTS: SprintParaAsignar[] = [
  { id: 1, estado: 'Finalizado', sprint_goal: 'Sprint viejo' },
  { id: 2, estado: 'Activo', sprint_goal: 'Sprint en curso' },
  { id: 3, estado: 'Pendiente', sprint_goal: null },
]

const asignarMock = vi.mocked(asignarASprint)

function abrir(sprints: SprintParaAsignar[] | null = SPRINTS) {
  const onAsignada = vi.fn()
  const onCancelar = vi.fn()
  render(
    <AsignarASprint
      historia={HISTORIA}
      sprints={sprints}
      errorSprints=""
      onAsignada={onAsignada}
      onCancelar={onCancelar}
    />,
  )
  return { onAsignada, onCancelar }
}

afterEach(() => {
  vi.resetAllMocks()
})

describe('AsignarASprint', () => {
  it('US3.1: el selector "Sprint" ofrece solo los Sprints Pendiente y Activo', () => {
    abrir()

    const opciones = screen.getAllByRole('option').map((opcion) => opcion.textContent)
    expect(opciones).toEqual([
      'Elegí un Sprint',
      'Sprint #2 (Activo) - Sprint en curso',
      'Sprint #3 (Pendiente)',
    ])
    expect(screen.getByLabelText('Sprint')).toBeTruthy()
  })

  it('US3.2: al confirmar llama a asignarASprint con la historia y el Sprint elegidos', async () => {
    asignarMock.mockResolvedValue({ ...HISTORIA, sprint_id: 3 })
    const user = userEvent.setup()
    abrir()

    await user.selectOptions(screen.getByLabelText('Sprint'), '3')
    await user.click(screen.getByRole('button', { name: 'Confirmar asignación' }))

    await waitFor(() => expect(asignarMock).toHaveBeenCalledTimes(1))
    expect(asignarMock).toHaveBeenCalledWith(7, 3)
  })

  it('US3.3: el botón de confirmar queda deshabilitado mientras la petición está en curso y no envía dos veces', async () => {
    asignarMock.mockReturnValue(new Promise(() => {}))
    const user = userEvent.setup()
    abrir()

    await user.selectOptions(screen.getByLabelText('Sprint'), '2')
    const boton = screen.getByRole('button', { name: 'Confirmar asignación' })
    await user.click(boton)

    await waitFor(() => expect((boton as HTMLButtonElement).disabled).toBe(true))
    await user.click(boton)
    expect(asignarMock).toHaveBeenCalledTimes(1)
  })

  it('US3.4: tras el 200 llama a onAsignada con la historia devuelta', async () => {
    const devuelta = { ...HISTORIA, estado: 'Pendiente' as const, sprint_id: 2 }
    asignarMock.mockResolvedValue(devuelta)
    const user = userEvent.setup()
    const { onAsignada } = abrir()

    await user.selectOptions(screen.getByLabelText('Sprint'), '2')
    await user.click(screen.getByRole('button', { name: 'Confirmar asignación' }))

    await waitFor(() => expect(onAsignada).toHaveBeenCalledTimes(1))
    expect(onAsignada).toHaveBeenCalledWith(devuelta, 2)
  })

  it('US3.5: ante un error muestra el mensaje de la API y no llama a onAsignada', async () => {
    asignarMock.mockRejectedValue(new ApiError(409, 'El Sprint ya finalizó'))
    const user = userEvent.setup()
    const { onAsignada } = abrir()

    await user.selectOptions(screen.getByLabelText('Sprint'), '2')
    await user.click(screen.getByRole('button', { name: 'Confirmar asignación' }))

    expect(await screen.findByText('El Sprint ya finalizó')).toBeTruthy()
    expect(onAsignada).not.toHaveBeenCalled()
  })

  it('US3.6: sin Sprints Pendiente/Activo muestra el aviso y no permite confirmar', () => {
    abrir([{ id: 1, estado: 'Finalizado', sprint_goal: null }])

    expect(screen.getByText('No hay Sprints disponibles para asignar')).toBeTruthy()
    expect(screen.queryByRole('button', { name: 'Confirmar asignación' })).toBeNull()
    expect(asignarMock).not.toHaveBeenCalled()
  })
})
