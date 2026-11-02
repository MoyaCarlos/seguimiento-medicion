import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { EsfuerzoHistoria } from './EsfuerzoHistoria'
import { imputarHoras, listarIntegrantes } from './api'

vi.mock('./api', () => ({
  consultarEsfuerzo: vi.fn(),
  imputarHoras: vi.fn(),
  listarIntegrantes: vi.fn(),
}))

const imputacion = { proyectoId: 42, fechaFin: '2999-12-31' }

const botonRegistrar = () => screen.getByRole('button', { name: 'Registrar horas' })

beforeEach(() => {
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

  it('US1.4: sin imputación no se ofrece registrar horas', () => {
    render(<EsfuerzoHistoria historiaId={12} />)

    expect(screen.queryByRole('button', { name: 'Registrar horas' })).toBeNull()
  })
})
