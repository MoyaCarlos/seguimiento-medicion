import { act, render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Route, Routes } from 'react-router'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { ApiError, ConnectionError } from '../../compartido/api/client'
import {
  consultarVotacion,
  iniciarVotacion,
  listarIntegrantes,
  obtenerHistoria,
  revelarVotacion,
  votar,
} from './api'
import { Component } from './index'
import type { HistoriaResumen, Integrante, Votacion } from './votacion'

vi.mock('./api')

const INTEGRANTES: Integrante[] = [
  { id: 1, nombre: 'Beta Uno', rol: 'scrum_master' },
  { id: 2, nombre: 'Beta Dos', rol: 'product_builder' },
  { id: 3, nombre: 'Beta Tres', rol: 'product_builder' },
]

const consultarMock = vi.mocked(consultarVotacion)
const votarMock = vi.mocked(votar)
const revelarMock = vi.mocked(revelarVotacion)
const iniciarMock = vi.mocked(iniciarVotacion)
const historiaMock = vi.mocked(obtenerHistoria)

const HISTORIA: HistoriaResumen = {
  id: 3,
  titulo: 'Pagar la tasa de la licencia en línea',
  descripcion: 'Como solicitante quiero pagar la tasa con tarjeta.',
  prioridad: 'M',
  estimacion_sp: null,
}
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
    promedio: servidor.revelada ? { exacto: 5.33, carta_sugerida: 5 } : null,
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

const fila = (nombre: string) =>
  within(screen.getByRole('list')).getByText(nombre).closest('li') as HTMLElement

async function elegirActor(nombre: string) {
  const selector = await screen.findByLabelText('Quién actúa')
  const opcion = within(selector).getByRole('option', { name: new RegExp(`^${nombre}`) })
  await userEvent.selectOptions(selector, opcion)
}

async function fallarConsultaAlActualizar() {
  consultarMock.mockRejectedValueOnce(new ConnectionError())
  await userEvent.click(screen.getByRole('button', { name: 'Actualizar' }))
  await screen.findByText('No se pudo conectar con el servidor')
}

beforeEach(() => {
  servidor = { votos: {}, revelada: false }
  integrantesMock.mockResolvedValue(INTEGRANTES)
  historiaMock.mockResolvedValue(HISTORIA)
  consultarMock.mockImplementation(async (_historia, consultante) => vista(consultante))
  votarMock.mockImplementation(async (_historia, usuario, valor) => {
    servidor.votos[usuario] = valor
    return vista(usuario)
  })
  revelarMock.mockImplementation(async (_historia, usuario) => {
    servidor.revelada = true
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

  it('Given dos pulsaciones de carta en el mismo instante, When se procesan, Then se vota una sola vez', async () => {
    votarMock.mockReturnValue(new Promise(() => {}))
    abrir()
    await elegirActor('Beta Dos')
    const carta = await screen.findByRole('button', { name: '8' })

    act(() => {
      carta.click()
      carta.click()
    })

    expect(votarMock).toHaveBeenCalledTimes(1)
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

  it('Given que falla el listado de integrantes y la consulta responde bien, When llega la votación, Then el error de integrantes sigue visible', async () => {
    integrantesMock.mockRejectedValue(new ApiError(500, 'no se pudieron listar los integrantes'))
    abrir()

    await screen.findByRole('list')

    expect((await screen.findByRole('alert')).textContent).toBe('no se pudieron listar los integrantes')
  })

  it('Given que fallan el listado de integrantes y la consulta con mensajes distintos, When ocurren, Then se muestran los dos errores', async () => {
    integrantesMock.mockRejectedValue(new ApiError(500, 'no se pudieron listar los integrantes'))
    consultarMock.mockRejectedValue(new ApiError(500, 'no se pudo consultar la votación'))
    abrir()

    await waitFor(() => expect(screen.getAllByRole('alert')).toHaveLength(2))
    const mensajes = screen.getAllByRole('alert').map((alerta) => alerta.textContent)
    expect(mensajes).toContain('no se pudieron listar los integrantes')
    expect(mensajes).toContain('no se pudo consultar la votación')
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

describe('Pantalla de votación - recuperación y consistencia', () => {
  it('Given que falló el listado de integrantes, When pulso Actualizar y el servidor responde, Then aparece el selector con los integrantes y desaparece el error', async () => {
    integrantesMock.mockRejectedValueOnce(new ApiError(500, 'no se pudieron listar los integrantes'))
    abrir()
    await screen.findByText('no se pudieron listar los integrantes')
    expect(screen.queryByLabelText('Quién actúa')).toBeNull()
    expect(screen.queryByText('Elegí quién actúa para votar')).toBeNull()

    await userEvent.click(screen.getByRole('button', { name: 'Actualizar' }))

    expect(await screen.findByLabelText('Quién actúa')).toBeTruthy()
    expect(screen.queryByText('no se pudieron listar los integrantes')).toBeNull()
    expect(screen.getByText('Elegí quién actúa para votar')).toBeTruthy()
  })

  it('Given que no cargó la historia, When pulso Actualizar y el servidor responde, Then se ve su resumen', async () => {
    historiaMock.mockRejectedValueOnce(new ConnectionError())
    abrir()
    await screen.findByRole('list')
    expect(screen.queryByText(HISTORIA.titulo)).toBeNull()

    await userEvent.click(screen.getByRole('button', { name: 'Actualizar' }))

    expect(await screen.findByText(HISTORIA.titulo)).toBeTruthy()
  })

  it('US1.3: Given que Beta Dos votó 8 y la consulta de Beta Tres está en curso, When cambio de integrante, Then no se ve el 8 de Beta Dos', async () => {
    servidor.votos = { 2: 8 }
    abrir()
    await elegirActor('Beta Dos')
    await waitFor(() => expect(within(fila('Beta Dos')).getByText('8')).toBeTruthy())

    consultarMock.mockReturnValueOnce(new Promise(() => {}))
    await elegirActor('Beta Tres')

    expect(within(fila('Beta Dos')).queryByText('8')).toBeNull()
    expect(within(fila('Beta Dos')).getByText('Ya votó')).toBeTruthy()
  })

  it('US1.3: Given que Beta Dos votó 8 y la consulta de Beta Tres falla, When cambio de integrante, Then el 8 de Beta Dos no queda en pantalla', async () => {
    servidor.votos = { 2: 8 }
    abrir()
    await elegirActor('Beta Dos')
    await waitFor(() => expect(within(fila('Beta Dos')).getByText('8')).toBeTruthy())

    consultarMock.mockRejectedValueOnce(new ConnectionError())
    await elegirActor('Beta Tres')
    await screen.findByText('No se pudo conectar con el servidor')

    expect(within(fila('Beta Dos')).queryByText('8')).toBeNull()
  })

  it('FR-001: Given integrantes con rol, When se abre el selector, Then cada opción muestra "Nombre (Rol)"', async () => {
    abrir()

    const selector = await screen.findByLabelText('Quién actúa')

    const opciones = within(selector).getAllByRole('option').map((opcion) => opcion.textContent)
    expect(opciones).toEqual([
      'Elegí un integrante',
      'Beta Uno (Scrum Master)',
      'Beta Dos (Product Builder)',
      'Beta Tres (Product Builder)',
    ])
  })

  it('Given un voto emitido mientras Actualizar consulta, When la consulta responde después, Then el voto sigue visible', async () => {
    abrir()
    await elegirActor('Beta Dos')
    await screen.findByRole('list')
    let responder: (v: Votacion) => void = () => {}
    consultarMock.mockImplementationOnce(() => new Promise((resolver) => (responder = resolver)))
    const sinVoto = vista(2)

    await userEvent.click(screen.getByRole('button', { name: 'Actualizar' }))
    await userEvent.click(screen.getByRole('button', { name: '8' }))
    await waitFor(() => expect(within(fila('Beta Dos')).getByText('8')).toBeTruthy())
    await act(async () => responder(sinVoto))

    expect(within(fila('Beta Dos')).getByText('8')).toBeTruthy()
  })

  it('Given una acción rechazada, When pulso Actualizar y el servidor responde, Then el mensaje de la acción desaparece', async () => {
    abrir()
    await elegirActor('Beta Dos')
    await screen.findByRole('list')
    votarMock.mockRejectedValueOnce(new ApiError(409, 'la ronda ya fue revelada'))
    await userEvent.click(screen.getByRole('button', { name: '8' }))
    await screen.findByText('la ronda ya fue revelada')

    await userEvent.click(screen.getByRole('button', { name: 'Actualizar' }))

    await waitFor(() => expect(screen.queryByText('la ronda ya fue revelada')).toBeNull())
  })

  it('Given una acción rechazada, When se recarga sola la consulta, Then el mensaje de la acción sigue visible', async () => {
    abrir()
    await elegirActor('Beta Dos')
    await screen.findByRole('list')
    votarMock.mockRejectedValueOnce(new ApiError(409, 'la ronda ya fue revelada'))
    await userEvent.click(screen.getByRole('button', { name: '8' }))
    await screen.findByText('la ronda ya fue revelada')

    await waitFor(() => expect(consultarMock).toHaveBeenCalledTimes(3))
    expect(screen.getByText('la ronda ya fue revelada')).toBeTruthy()
  })

  it('Given que falló una consulta, When voto y el servidor responde bien, Then desaparece el aviso de conexión y queda el voto registrado', async () => {
    abrir()
    await elegirActor('Beta Dos')
    await screen.findByRole('button', { name: '8' })
    await fallarConsultaAlActualizar()

    await userEvent.click(screen.getByRole('button', { name: '8' }))

    await waitFor(() => expect(within(fila('Beta Dos')).getByText('8')).toBeTruthy())
    expect(screen.queryByRole('alert')).toBeNull()
  })

  it('Given que falló una consulta, When revelo y el servidor responde bien, Then desaparece el aviso de conexión y se ve el resultado', async () => {
    servidor.votos = { 1: 3, 2: 5, 3: 8 }
    abrir()
    await elegirActor('Beta Uno')
    await screen.findByRole('button', { name: 'Revelar' })
    await fallarConsultaAlActualizar()

    await userEvent.click(screen.getByRole('button', { name: 'Revelar' }))

    expect(await screen.findByText('Promedio exacto')).toBeTruthy()
    expect(screen.queryByRole('alert')).toBeNull()
  })

  it('Given un id de historia inválido en la ruta, When se abre la pantalla, Then se avisa y no se consulta al servidor', async () => {
    render(
      <MemoryRouter initialEntries={['/projects/25/votacion/abc']}>
        <Routes>
          <Route path="projects/:id/votacion/:historiaId" element={<Component />} />
        </Routes>
      </MemoryRouter>,
    )

    expect((await screen.findByRole('alert')).textContent).toBe('La historia "abc" no es válida')
    expect(screen.queryByText(/NaN/)).toBeNull()
    expect(consultarMock).not.toHaveBeenCalled()
    expect(integrantesMock).not.toHaveBeenCalled()
    expect(historiaMock).not.toHaveBeenCalled()
  })
})

describe('Pantalla de votación - revelar y ver el resultado (US3)', () => {
  it('US3.1: Given que votaron todos y soy Scrum Master, When pulso Revelar, Then se ven los valores, el promedio exacto y la carta sugerida, sin cartas', async () => {
    servidor.votos = { 1: 3, 2: 5, 3: 8 }
    abrir()
    await elegirActor('Beta Uno')

    await userEvent.click(await screen.findByRole('button', { name: 'Revelar' }))

    await waitFor(() => expect(within(fila('Beta Tres')).getByText('8')).toBeTruthy())
    expect(revelarMock).toHaveBeenCalledWith(3, 1)
    expect(within(fila('Beta Uno')).getByText('3')).toBeTruthy()
    expect(within(fila('Beta Dos')).getByText('5')).toBeTruthy()
    expect(screen.getByText('Promedio exacto').nextSibling?.textContent).toBe('5,33')
    expect(screen.getByText('Carta sugerida').nextSibling?.textContent).toBe('5')
    expect(screen.queryByRole('button', { name: '8' })).toBeNull()
    expect(screen.queryByRole('button', { name: 'Revelar' })).toBeNull()
  })

  it('US3.2: Given que falta votar a Beta Tres, When soy Scrum Master, Then Revelar está deshabilitado y se indica quién falta', async () => {
    servidor.votos = { 1: 3, 2: 5 }
    abrir()
    await elegirActor('Beta Uno')

    const revelar = (await screen.findByRole('button', { name: 'Revelar' })) as HTMLButtonElement

    expect(revelar.disabled).toBe(true)
    expect(screen.getByText('Falta votar: Beta Tres')).toBeTruthy()
  })

  it('US3.3: Given que soy Product Builder, When votaron todos, Then no hay botón Revelar', async () => {
    servidor.votos = { 1: 3, 2: 5, 3: 8 }
    abrir()
    await elegirActor('Beta Dos')

    await screen.findByRole('button', { name: '8' })

    expect(screen.queryByRole('button', { name: 'Revelar' })).toBeNull()
  })

  it('US3.4: Given una revelación en curso, When pulso Revelar otra vez, Then se revela una sola vez', async () => {
    servidor.votos = { 1: 3, 2: 5, 3: 8 }
    revelarMock.mockReturnValue(new Promise(() => {}))
    abrir()
    await elegirActor('Beta Uno')

    await userEvent.click(await screen.findByRole('button', { name: 'Revelar' }))
    await userEvent.click(screen.getByRole('button', { name: 'Revelar' }))

    expect(revelarMock).toHaveBeenCalledTimes(1)
    expect((screen.getByRole('button', { name: 'Revelar' }) as HTMLButtonElement).disabled).toBe(true)
  })

  it('US3.5: Given que el servidor rechaza revelar, When ocurre el error, Then se muestra su mensaje y la pantalla sigue operable', async () => {
    servidor.votos = { 1: 3, 2: 5, 3: 8 }
    revelarMock.mockRejectedValue(new ApiError(403, 'solo el Scrum Master del proyecto puede hacerlo'))
    abrir()
    await elegirActor('Beta Uno')

    await userEvent.click(await screen.findByRole('button', { name: 'Revelar' }))

    expect((await screen.findByRole('alert')).textContent).toBe('solo el Scrum Master del proyecto puede hacerlo')
    expect(fila('Beta Uno')).toBeTruthy()
  })

  it('Given una ronda revelada, When se muestra la pantalla, Then se aclara que registrar la estimación acordada se hace en otra historia', async () => {
    servidor.votos = { 1: 3, 2: 5, 3: 8 }
    servidor.revelada = true
    abrir()

    expect(await screen.findByText(/registrar la estimación acordada/i)).toBeTruthy()
    expect(screen.queryByText('Elegí quién actúa para votar')).toBeNull()
  })
})

describe('Pantalla de votación - iniciar la votación (US2)', () => {
  const SIN_VOTACION = new ApiError(404, 'la historia no tiene votación iniciada')
  let iniciada: boolean

  beforeEach(() => {
    iniciada = false
    consultarMock.mockImplementation(async (_historia, consultante) => {
      if (!iniciada) throw SIN_VOTACION
      return vista(consultante)
    })
    iniciarMock.mockImplementation(async (_historia, usuario) => {
      iniciada = true
      return vista(usuario)
    })
  })

  it('US2.0: Given una historia sin votación, When se abre la pantalla, Then se informa que no hay votación, sin cartas ni error', async () => {
    abrir()

    expect(await screen.findByText('La historia todavía no tiene votación')).toBeTruthy()
    expect(screen.queryByRole('button', { name: '8' })).toBeNull()
    expect(screen.queryByRole('alert')).toBeNull()
  })

  it('US2.1: Given el Scrum Master elegido, When pulso Iniciar votación, Then aparece la ronda 1 con todos los integrantes en "Falta"', async () => {
    abrir()
    await elegirActor('Beta Uno')

    await userEvent.click(await screen.findByRole('button', { name: 'Iniciar votación' }))

    expect(await screen.findByText('Ronda 1')).toBeTruthy()
    expect(iniciarMock).toHaveBeenCalledWith(3, 1)
    expect(within(fila('Beta Uno')).getByText('Falta')).toBeTruthy()
    expect(within(fila('Beta Dos')).getByText('Falta')).toBeTruthy()
    expect(within(fila('Beta Tres')).getByText('Falta')).toBeTruthy()
    expect(screen.queryByRole('button', { name: 'Iniciar votación' })).toBeNull()
  })

  it('US2.2: Given un Product Builder elegido, When la historia no tiene votación, Then no hay botón Iniciar y se indica que solo lo hace el Scrum Master', async () => {
    abrir()
    await elegirActor('Beta Dos')

    expect(await screen.findByText('Solo el Scrum Master puede iniciar la votación')).toBeTruthy()
    expect(screen.queryByRole('button', { name: 'Iniciar votación' })).toBeNull()
  })

  it('US2.4: Given un inicio en curso, When pulso Iniciar votación otra vez, Then se inicia una sola vez', async () => {
    iniciarMock.mockReturnValue(new Promise(() => {}))
    abrir()
    await elegirActor('Beta Uno')

    await userEvent.click(await screen.findByRole('button', { name: 'Iniciar votación' }))
    await userEvent.click(screen.getByRole('button', { name: 'Iniciar votación' }))

    expect(iniciarMock).toHaveBeenCalledTimes(1)
    expect((screen.getByRole('button', { name: 'Iniciar votación' }) as HTMLButtonElement).disabled).toBe(true)
  })

  it('US2.5: Given que el servidor rechaza el inicio, When ocurre el error, Then se muestra su mensaje', async () => {
    iniciarMock.mockRejectedValue(new ApiError(409, 'la historia ya tiene una votación'))
    abrir()
    await elegirActor('Beta Uno')

    await userEvent.click(await screen.findByRole('button', { name: 'Iniciar votación' }))

    expect((await screen.findByRole('alert')).textContent).toBe('la historia ya tiene una votación')
  })

  it('Given que falló una consulta, When inicio la votación y el servidor responde bien, Then desaparece el aviso de conexión y aparece la ronda 1', async () => {
    abrir()
    await elegirActor('Beta Uno')
    await screen.findByRole('button', { name: 'Iniciar votación' })
    await fallarConsultaAlActualizar()

    await userEvent.click(screen.getByRole('button', { name: 'Iniciar votación' }))

    expect(await screen.findByText('Ronda 1')).toBeTruthy()
    expect(screen.queryByRole('alert')).toBeNull()
  })

  it('Given una historia inexistente, When se abre la pantalla, Then se muestra "historia no encontrada" y no se ofrece iniciar', async () => {
    consultarMock.mockRejectedValue(new ApiError(404, 'historia no encontrada'))
    abrir()
    await elegirActor('Beta Uno')

    expect((await screen.findByRole('alert')).textContent).toBe('historia no encontrada')
    expect(screen.queryByRole('button', { name: 'Iniciar votación' })).toBeNull()
  })
})

describe('Pantalla de votación - mantenerla al día (US4)', () => {
  it('US4.2: Given que otro integrante votó, When pulso Actualizar, Then se vuelve a consultar y se refleja su voto', async () => {
    abrir()
    await elegirActor('Beta Uno')
    await waitFor(() => expect(within(fila('Beta Dos')).getByText('Falta')).toBeTruthy())

    servidor.votos = { 2: 8 }
    await userEvent.click(screen.getByRole('button', { name: 'Actualizar' }))

    await waitFor(() => expect(within(fila('Beta Dos')).getByText('Oculto')).toBeTruthy())
    expect(within(fila('Beta Dos')).getByText('Ya votó')).toBeTruthy()
  })

  it('US4.4: Given que la consulta falló por falta de conexión, When pulso Actualizar y el servidor responde, Then se muestra la votación y desaparece el aviso', async () => {
    consultarMock.mockRejectedValueOnce(new ConnectionError())
    abrir()
    expect((await screen.findByRole('alert')).textContent).toBe('No se pudo conectar con el servidor')

    await userEvent.click(screen.getByRole('button', { name: 'Actualizar' }))

    expect(await screen.findByRole('list')).toBeTruthy()
    expect(fila('Beta Uno')).toBeTruthy()
    expect(screen.queryByRole('alert')).toBeNull()
  })
})

describe('Pantalla de votación - contexto de la historia', () => {
  it('Given que la historia está en el listado, When se abre la pantalla, Then se ven su título, su descripción y la consulta con el proyecto y la historia', async () => {
    abrir()

    expect(await screen.findByRole('heading', { name: HISTORIA.titulo })).toBeTruthy()
    expect(screen.getByText(HISTORIA.descripcion)).toBeTruthy()
    expect(screen.getByText('Sin estimar')).toBeTruthy()
    expect(historiaMock).toHaveBeenCalledWith('25', 3)
  })

  it('Given que la historia no está en el listado, When se abre la pantalla, Then solo figura "Historia #3" y se puede votar', async () => {
    historiaMock.mockResolvedValue(null)
    abrir()
    await elegirActor('Beta Dos')

    expect(await screen.findByRole('button', { name: '8' })).toBeTruthy()
    expect(screen.getAllByText('Historia #3').length).toBeGreaterThan(0)
    expect(screen.queryByRole('heading', { name: HISTORIA.titulo })).toBeNull()
  })

  it('Given que falla la consulta de la historia, When se abre la pantalla, Then solo figura "Historia #3", sin error y se puede votar', async () => {
    historiaMock.mockRejectedValue(new ConnectionError())
    abrir()
    await elegirActor('Beta Dos')

    expect(await screen.findByRole('button', { name: '8' })).toBeTruthy()
    expect(screen.getAllByText('Historia #3').length).toBeGreaterThan(0)
    expect(screen.queryByRole('alert')).toBeNull()
  })

  it('Given una votación abierta y un integrante elegido, When se muestra la pantalla, Then se ve la consigna de la votación', async () => {
    abrir()
    await elegirActor('Beta Dos')

    expect(
      await screen.findByText('Elegí una carta: ¿cuántos Story Points creés que lleva esta historia?'),
    ).toBeTruthy()
  })

  it('Given una votación abierta sin integrante elegido, When se muestra la pantalla, Then la consigna es elegir quién actúa', async () => {
    abrir()

    expect(await screen.findByText('Elegí quién actúa para votar')).toBeTruthy()
    expect(screen.queryByText(/Elegí una carta/)).toBeNull()
  })

  it('Given una ronda revelada, When se muestra la pantalla, Then no se muestra la consigna de elegir carta', async () => {
    servidor.votos = { 1: 3, 2: 5, 3: 8 }
    servidor.revelada = true
    abrir()

    await screen.findByText(/registrar la estimación acordada/i)

    expect(screen.queryByText(/Elegí una carta/)).toBeNull()
  })

  it('Given la pantalla del proyecto 25, When se muestra, Then hay un enlace "Volver al Product Backlog" a /projects/25/backlog', async () => {
    abrir()

    const enlace = await screen.findByRole('link', { name: 'Volver al Product Backlog' })

    expect(enlace.getAttribute('href')).toBe('/projects/25/backlog')
  })
})

describe('Pantalla de votación - teclado', () => {
  it('Given una carta enfocada con Tab, When se pulsa Enter, Then se vota esa carta', async () => {
    abrir()
    await elegirActor('Beta Dos')
    const primera = await screen.findByRole('button', { name: '1' })

    await userEvent.tab()
    while (document.activeElement !== primera) await userEvent.tab()
    await userEvent.tab()
    await userEvent.keyboard('{Enter}')

    expect(votarMock).toHaveBeenCalledWith(3, 2, 2)
  })
})
