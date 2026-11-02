import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Route, Routes } from 'react-router'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { ApiError, ConnectionError } from '../../compartido/api/client'
import { consultarVotacion, listarIntegrantes, votar } from './api'
import { Component } from './index'
import type { Integrante, Votacion } from './votacion'

vi.mock('./api')

const INTEGRANTES: Integrante[] = [
  { id: 1, nombre: 'Beta Uno', rol: 'scrum_master' },
  { id: 2, nombre: 'Beta Dos', rol: 'product_builder' },
  { id: 3, nombre: 'Beta Tres', rol: 'product_builder' },
]

const consultarMock = vi.mocked(consultarVotacion)
const votarMock = vi.mocked(votar)
const integrantesMock = vi.mocked(listarIntegrantes)

type Servidor = { votos: Record<number, number>; revelada: boolean }
let servidor: Servidor

function vista(consultante?: number): Votacion {
  return {
    historia_id: 3,
    ronda: 1,
    revelada_en: servidor.revelada ? '2026-11-02T10:00:00Z' : null,
    integrantes: INTEGRANTES.map((integrante) => {
      const voto = servidor.votos[integrante.id]
      const visible = servidor.revelada || integrante.id === consultante
      return {
        usuario_id: integrante.id,
        nombre: integrante.nombre,
        rol: integrante.rol,
        ya_voto: voto !== undefined,
        valor: visible && voto !== undefined ? voto : null,
      }
    }),
    promedio: null,
  }
}

function abrir() {
  render(
    <MemoryRouter initialEntries={['/projects/25/votacion/3']}>
      <Routes>
        <Route path="projects/:id/votacion/:historiaId" element={<Component />} />
      </Routes>
    </MemoryRouter>,
  )
}

const fila = (nombre: string) => screen.getByRole('listitem', { name: nombre })

async function elegirActor(nombre: string) {
  await userEvent.selectOptions(await screen.findByLabelText('Quién actúa'), nombre)
}

beforeEach(() => {
  servidor = { votos: {}, revelada: false }
  integrantesMock.mockResolvedValue(INTEGRANTES)
  consultarMock.mockImplementation(async (_historia, consultante) => vista(consultante))
  votarMock.mockImplementation(async (_historia, usuario, valor) => {
    servidor.votos[usuario] = valor
    return vista(usuario)
  })
})

afterEach(() => {
  vi.resetAllMocks()
})

describe('Pantalla de votación - votar en secreto (US1)', () => {
  it('US1.0: Given integrantes del proyecto y una votación abierta, When se abre la pantalla, Then se pide elegir quién actúa y no hay cartas', async () => {
    abrir()

    expect(await screen.findByText('Elegí quién actúa para votar')).toBeTruthy()
    expect(fila('Beta Uno')).toBeTruthy()
    expect(screen.queryByRole('button', { name: '8' })).toBeNull()
    expect(integrantesMock).toHaveBeenCalledWith('25')
    expect(consultarMock).toHaveBeenCalledWith(3, undefined)
  })

  it('US1.1: Given una votación abierta, When elijo a Beta Dos y la carta 8, Then se vota (3, 2, 8) y se ve mi voto y "Falta" en los demás', async () => {
    servidor.votos = { 3: 5 }
    abrir()
    await elegirActor('Beta Dos')

    await userEvent.click(await screen.findByRole('button', { name: '8' }))

    await waitFor(() => expect(within(fila('Beta Dos')).getByText('8')).toBeTruthy())
    expect(votarMock).toHaveBeenCalledWith(3, 2, 8)
    expect(within(fila('Beta Tres')).getByText('Oculto')).toBeTruthy()
    expect(within(fila('Beta Uno')).getByText('Falta')).toBeTruthy()
    expect(screen.getByRole('button', { name: '8' }).getAttribute('aria-pressed')).toBe('true')
  })

  it('US1.2: Given que ya voté 8, When elijo la carta 13, Then se muestra 13 como mi voto', async () => {
    servidor.votos = { 2: 8 }
    abrir()
    await elegirActor('Beta Dos')

    await userEvent.click(await screen.findByRole('button', { name: '13' }))

    await waitFor(() => expect(within(fila('Beta Dos')).getByText('13')).toBeTruthy())
    expect(votarMock).toHaveBeenCalledWith(3, 2, 13)
  })

  it('US1.3: Given que Beta Dos votó 8 y Beta Tres 5, When cambio de integrante, Then se consulta con el nuevo usuario y se ve su voto', async () => {
    servidor.votos = { 2: 8, 3: 5 }
    abrir()

    await elegirActor('Beta Dos')
    await waitFor(() => expect(within(fila('Beta Dos')).getByText('8')).toBeTruthy())
    expect(within(fila('Beta Tres')).getByText('Oculto')).toBeTruthy()

    await elegirActor('Beta Tres')
    await waitFor(() => expect(within(fila('Beta Tres')).getByText('5')).toBeTruthy())
    expect(consultarMock).toHaveBeenCalledWith(3, 3)
    expect(within(fila('Beta Dos')).getByText('Oculto')).toBeTruthy()
  })

  it('US1.4: Given un voto en curso, When pulso otra carta, Then se vota una sola vez y las cartas quedan deshabilitadas', async () => {
    votarMock.mockReturnValue(new Promise(() => {}))
    abrir()
    await elegirActor('Beta Dos')

    await userEvent.click(await screen.findByRole('button', { name: '8' }))
    await userEvent.click(screen.getByRole('button', { name: '13' }))

    expect(votarMock).toHaveBeenCalledTimes(1)
    expect((screen.getByRole('button', { name: '13' }) as HTMLButtonElement).disabled).toBe(true)
  })

  it('US1.5: Given que el servidor rechaza el voto, When ocurre el error, Then se muestra su mensaje, se vuelve a consultar y las cartas siguen operables', async () => {
    votarMock.mockRejectedValue(new ApiError(400, 'la ronda ya fue revelada', 'estado'))
    abrir()
    await elegirActor('Beta Dos')
    await screen.findByRole('button', { name: '8' })
    const consultasAntes = consultarMock.mock.calls.length

    await userEvent.click(screen.getByRole('button', { name: '8' }))

    expect((await screen.findByRole('alert')).textContent).toBe('la ronda ya fue revelada')
    await waitFor(() => expect(consultarMock.mock.calls.length).toBeGreaterThan(consultasAntes))
    expect((screen.getByRole('button', { name: '8' }) as HTMLButtonElement).disabled).toBe(false)
  })

  it('Given un proyecto sin integrantes, When se abre la pantalla, Then se informa y no hay selector ni cartas', async () => {
    integrantesMock.mockResolvedValue([])
    abrir()

    expect(await screen.findByText('El proyecto no tiene integrantes')).toBeTruthy()
    expect(screen.queryByLabelText('Quién actúa')).toBeNull()
  })

  it('US4.4: Given que el servidor no responde, When se abre la pantalla, Then se muestra el aviso de conexión', async () => {
    consultarMock.mockRejectedValue(new ConnectionError())
    abrir()

    expect((await screen.findByRole('alert')).textContent).toBe('No se pudo conectar con el servidor')
  })

  it('Given que la respuesta de un integrante anterior llega tarde, When ya elegí a otro, Then se descarta y se ve la del actual', async () => {
    servidor.votos = { 2: 8, 3: 5 }
    let resolverTardia: (v: Votacion) => void = () => {}
    abrir()
    await elegirActor('Beta Dos')
    await waitFor(() => expect(within(fila('Beta Dos')).getByText('8')).toBeTruthy())

    consultarMock.mockImplementationOnce(() => new Promise((resolver) => (resolverTardia = resolver)))
    await elegirActor('Beta Uno')
    await elegirActor('Beta Tres')
    await waitFor(() => expect(within(fila('Beta Tres')).getByText('5')).toBeTruthy())

    resolverTardia(vista(1))

    await waitFor(() => expect(within(fila('Beta Tres')).getByText('5')).toBeTruthy())
    expect(within(fila('Beta Dos')).getByText('Oculto')).toBeTruthy()
  })
})
