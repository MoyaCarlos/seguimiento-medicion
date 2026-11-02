import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { ApiError } from '../../compartido/api/client'
import { EsfuerzoHistoria } from './EsfuerzoHistoria'
import { consultarEsfuerzo, imputarHoras, listarIntegrantes } from './api'
import type { EsfuerzoReal } from './horas'

vi.mock('./api', () => ({
  consultarEsfuerzo: vi.fn(),
  imputarHoras: vi.fn(),
  listarIntegrantes: vi.fn(),
}))

const imputacion = { proyectoId: 42, fechaFin: '2999-12-31' }

const SIN_CARGAS: EsfuerzoReal = { total_horas: 0, registros: [] }
const CARGA_ANA = { id: 10, integrante: 'Ana', fecha: '2026-10-27', actividad: 'Revisión del PR', horas: 1.5 }
const CARGA_JIMENA = { id: 11, integrante: 'Jimena', fecha: '2026-10-28', actividad: 'Implementar handler', horas: 4 }
const CON_CARGAS: EsfuerzoReal = { total_horas: 5.5, registros: [CARGA_ANA, CARGA_JIMENA] }

const esfuerzoMock = vi.mocked(consultarEsfuerzo)

const botonRegistrar = () => screen.getByRole('button', { name: 'Registrar horas' })
const botonDetalle = () => screen.getByRole('button', { name: 'Ver detalle' })

async function guardarCarga() {
  const user = userEvent.setup()
  await user.click(botonRegistrar())
  await user.selectOptions(await screen.findByLabelText('Integrante'), 'Ana')
  await user.type(screen.getByLabelText('Actividad'), 'Revisión')
  await user.type(screen.getByLabelText('Horas'), '4')
  await user.click(screen.getByRole('button', { name: 'Guardar horas' }))
  await screen.findByText('Horas registradas')
  return user
}

beforeEach(() => {
  esfuerzoMock.mockResolvedValue(SIN_CARGAS)
  vi.mocked(listarIntegrantes).mockResolvedValue([{ id: 2, nombre: 'Ana' }])
  vi.mocked(imputarHoras).mockResolvedValue({})
})

afterEach(() => {
  vi.resetAllMocks()
})

describe('EsfuerzoHistoria: registrar horas', () => {
  it('US1.1: "Registrar horas" abre el formulario dentro de la historia', async () => {
    render(<EsfuerzoHistoria historiaId={12} imputacion={imputacion} />)

    expect(screen.queryByLabelText('Actividad')).toBeNull()
    expect(botonRegistrar().getAttribute('aria-expanded')).toBe('false')

    await userEvent.click(botonRegistrar())

    expect(await screen.findByLabelText('Actividad')).toBeTruthy()
    expect(botonRegistrar().getAttribute('aria-expanded')).toBe('true')
    expect(listarIntegrantes).toHaveBeenCalledWith(42)
  })

  it('US1.1: al cerrar y volver a abrir, la lista vuelve sin nadie seleccionado', async () => {
    const user = userEvent.setup()
    render(<EsfuerzoHistoria historiaId={12} imputacion={imputacion} />)
    await user.click(botonRegistrar())
    await user.selectOptions(await screen.findByLabelText('Integrante'), 'Ana')

    await user.click(botonRegistrar())
    expect(screen.queryByLabelText('Integrante')).toBeNull()
    await user.click(botonRegistrar())

    expect(((await screen.findByLabelText('Integrante')) as HTMLSelectElement).value).toBe('')
  })

  it('US1.5: mientras envía no se puede cerrar el formulario para volver a guardar', async () => {
    let terminarEnvio: (valor: unknown) => void = () => {}
    vi.mocked(imputarHoras).mockReturnValue(new Promise((resolve) => (terminarEnvio = resolve)))
    const user = userEvent.setup()
    render(<EsfuerzoHistoria historiaId={12} imputacion={imputacion} />)
    await user.click(botonRegistrar())
    await user.selectOptions(await screen.findByLabelText('Integrante'), 'Ana')
    await user.type(screen.getByLabelText('Actividad'), 'Revisión')
    await user.type(screen.getByLabelText('Horas'), '4')

    await user.click(screen.getByRole('button', { name: 'Guardar horas' }))

    expect((botonRegistrar() as HTMLButtonElement).disabled).toBe(true)
    await user.click(botonRegistrar())
    expect(screen.getByLabelText('Actividad')).toBeTruthy()

    terminarEnvio({})
    await screen.findByText('Horas registradas')
    expect((botonRegistrar() as HTMLButtonElement).disabled).toBe(false)
  })

  it('US1.4: sin imputación no se ofrece registrar horas', () => {
    render(<EsfuerzoHistoria historiaId={12} />)

    expect(screen.queryByRole('button', { name: 'Registrar horas' })).toBeNull()
  })
})

describe('EsfuerzoHistoria: esfuerzo acumulado', () => {
  it('US3.1: una historia con cargas muestra el total en horas', async () => {
    esfuerzoMock.mockResolvedValue(CON_CARGAS)

    render(<EsfuerzoHistoria historiaId={12} />)

    expect(await screen.findByText('Esfuerzo real: 5,5 h')).toBeTruthy()
    expect(esfuerzoMock).toHaveBeenCalledWith(12)
  })

  it('US3.2: una historia sin cargas muestra 0 h, sin error', async () => {
    render(<EsfuerzoHistoria historiaId={12} />)

    expect(await screen.findByText('Esfuerzo real: 0 h')).toBeTruthy()
    expect(screen.queryByRole('alert')).toBeNull()
  })

  it('muestra la carga mientras consulta', () => {
    esfuerzoMock.mockReturnValue(new Promise(() => {}))

    render(<EsfuerzoHistoria historiaId={12} />)

    expect(screen.getByText('Cargando…')).toBeTruthy()
  })

  it('US3.3: al guardar una carga se actualiza el total sin recargar', async () => {
    esfuerzoMock.mockResolvedValueOnce(SIN_CARGAS).mockResolvedValue({ total_horas: 4, registros: [CARGA_JIMENA] })
    render(<EsfuerzoHistoria historiaId={12} imputacion={imputacion} />)
    await screen.findByText('Esfuerzo real: 0 h')

    await guardarCarga()

    expect(await screen.findByText('Esfuerzo real: 4 h')).toBeTruthy()
    expect(esfuerzoMock).toHaveBeenCalledTimes(2)
  })

  it('US3.3: mientras reconsulta conserva el total anterior', async () => {
    esfuerzoMock.mockResolvedValueOnce(CON_CARGAS).mockReturnValue(new Promise(() => {}))
    render(<EsfuerzoHistoria historiaId={12} imputacion={imputacion} />)
    await screen.findByText('Esfuerzo real: 5,5 h')

    await guardarCarga()

    expect(screen.getByText('Esfuerzo real: 5,5 h')).toBeTruthy()
    expect(screen.queryByText('Cargando…')).toBeNull()
  })

  it('US3.4: si falla la consulta muestra el aviso en lugar del total', async () => {
    esfuerzoMock.mockRejectedValue(new ApiError(404, 'historia no encontrada'))

    render(<EsfuerzoHistoria historiaId={12} />)

    expect((await screen.findByRole('alert')).textContent).toBe('historia no encontrada')
    expect(screen.queryByText(/Esfuerzo real/)).toBeNull()
  })

  it('US3.4: si falla la consulta se puede seguir registrando horas', async () => {
    esfuerzoMock.mockRejectedValue(new ApiError(500, 'error interno'))

    render(<EsfuerzoHistoria historiaId={12} imputacion={imputacion} />)

    await screen.findByRole('alert')
    expect(botonRegistrar()).toBeTruthy()
  })

  it('US3.4: si falla la reconsulta tras guardar, confirma la carga y avisa en lugar del total', async () => {
    esfuerzoMock.mockResolvedValueOnce(CON_CARGAS).mockRejectedValue(new ApiError(500, 'error interno'))
    render(<EsfuerzoHistoria historiaId={12} imputacion={imputacion} />)
    await screen.findByText('Esfuerzo real: 5,5 h')

    await guardarCarga()

    expect((await screen.findByRole('alert')).textContent).toBe('error interno')
    expect(screen.getByText('Horas registradas')).toBeTruthy()
    expect(screen.queryByText(/Esfuerzo real/)).toBeNull()
  })

  it('US3.5: el detalle empieza plegado y al desplegarlo muestra cada carga en el orden de la API', async () => {
    esfuerzoMock.mockResolvedValue(CON_CARGAS)
    render(<EsfuerzoHistoria historiaId={12} />)
    await screen.findByText('Esfuerzo real: 5,5 h')

    expect(screen.queryByRole('list')).toBeNull()
    expect(botonDetalle().getAttribute('aria-expanded')).toBe('false')

    await userEvent.click(botonDetalle())

    const cargas = within(screen.getByRole('list')).getAllByRole('listitem')
    expect(cargas).toHaveLength(2)
    expect(cargas[0].textContent).toContain('Ana')
    expect(cargas[0].textContent).toContain('2026-10-27')
    expect(cargas[0].textContent).toContain('Revisión del PR')
    expect(cargas[0].textContent).toContain('1,5 h')
    expect(cargas[1].textContent).toContain('Jimena')
    expect(botonDetalle().getAttribute('aria-expanded')).toBe('true')
  })

  it('US3.5: al volver a elegir "Ver detalle" se pliega', async () => {
    esfuerzoMock.mockResolvedValue(CON_CARGAS)
    const user = userEvent.setup()
    render(<EsfuerzoHistoria historiaId={12} />)
    await screen.findByText('Esfuerzo real: 5,5 h')

    await user.click(botonDetalle())
    await user.click(botonDetalle())

    expect(screen.queryByRole('list')).toBeNull()
  })

  it('US3.6: una historia sin cargas no ofrece "Ver detalle"', async () => {
    render(<EsfuerzoHistoria historiaId={12} />)

    await screen.findByText('Esfuerzo real: 0 h')
    expect(screen.queryByRole('button', { name: 'Ver detalle' })).toBeNull()
  })

  it('US3.7: con el detalle desplegado, la carga nueva aparece al guardar', async () => {
    esfuerzoMock
      .mockResolvedValueOnce({ total_horas: 1.5, registros: [CARGA_ANA] })
      .mockResolvedValue(CON_CARGAS)
    render(<EsfuerzoHistoria historiaId={12} imputacion={imputacion} />)
    await screen.findByText('Esfuerzo real: 1,5 h')
    await userEvent.click(botonDetalle())

    await guardarCarga()

    await screen.findByText('Esfuerzo real: 5,5 h')
    expect(screen.getByText(/Implementar handler/)).toBeTruthy()
  })
})
