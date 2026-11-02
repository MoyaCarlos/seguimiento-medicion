import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { FormularioHoras } from './FormularioHoras'
import { imputarHoras, listarIntegrantes } from './api'
import { fechaPropuesta } from './horas'

vi.mock('./api', () => ({ imputarHoras: vi.fn(), listarIntegrantes: vi.fn() }))

const imputarMock = vi.mocked(imputarHoras)
const integrantesMock = vi.mocked(listarIntegrantes)

const INTEGRANTES = [
  { id: 2, nombre: 'Ana' },
  { id: 4, nombre: 'Jimena' },
]

function abrir(fechaFin?: string) {
  const alRegistrar = vi.fn()
  render(<FormularioHoras historiaId={12} proyectoId={42} fechaFin={fechaFin} alRegistrar={alRegistrar} />)
  return { alRegistrar }
}

const lista = () => screen.findByLabelText('Integrante') as Promise<HTMLSelectElement>
const fecha = () => screen.getByLabelText('Fecha') as HTMLInputElement
const actividad = () => screen.getByLabelText('Actividad') as HTMLInputElement
const horas = () => screen.getByLabelText('Horas') as HTMLInputElement
const guardar = () => screen.getByRole('button', { name: 'Guardar horas' }) as HTMLButtonElement

async function completar(texto = '4') {
  const user = userEvent.setup()
  await user.selectOptions(await lista(), 'Ana')
  await user.type(actividad(), 'Revisión del PR')
  await user.type(horas(), texto)
  return user
}

beforeEach(() => {
  integrantesMock.mockResolvedValue(INTEGRANTES)
  imputarMock.mockResolvedValue({})
})

afterEach(() => {
  vi.resetAllMocks()
})

describe('FormularioHoras: registrar horas', () => {
  it('US1.1: consulta los integrantes del proyecto y muestra la carga mientras espera', () => {
    integrantesMock.mockReturnValue(new Promise(() => {}))
    abrir()

    expect(screen.getByText('Cargando…')).toBeTruthy()
    expect(integrantesMock).toHaveBeenCalledWith(42)
  })

  it('US1.1: muestra la lista sin nadie seleccionado y los campos de fecha, actividad y horas', async () => {
    abrir()

    const seleccion = await lista()
    expect(seleccion.value).toBe('')
    expect(within(seleccion).getAllByRole('option').map((o) => o.textContent)).toEqual([
      'Elegí un integrante',
      'Ana',
      'Jimena',
    ])
    expect(fecha().type).toBe('date')
    expect(actividad()).toBeTruthy()
    expect(horas().type).toBe('text')
  })

  it('US1.1: propone la fecha de hoy', async () => {
    abrir()

    await lista()
    expect(fecha().value).toBe(fechaPropuesta(new Date()))
  })

  it('US1.1: si el Sprint ya pasó su fecha de fin, propone esa fecha', async () => {
    abrir('2020-01-31')

    await lista()
    expect(fecha().value).toBe('2020-01-31')
  })

  it('US1.2: al guardar envía los datos elegidos, confirma y avisa', async () => {
    const { alRegistrar } = abrir()
    const user = await completar()

    await user.click(guardar())

    expect((await screen.findByRole('status')).textContent).toBe('Horas registradas')
    expect(imputarMock).toHaveBeenCalledWith(12, {
      integrante_id: 2,
      fecha: fechaPropuesta(new Date()),
      actividad: 'Revisión del PR',
      horas: '4',
    })
    expect(alRegistrar).toHaveBeenCalledTimes(1)
  })

  it('US1.2: al guardar vacía actividad y horas, y conserva integrante y fecha', async () => {
    abrir()
    const user = await completar()

    await user.click(guardar())

    await screen.findByText('Horas registradas')
    expect(actividad().value).toBe('')
    expect(horas().value).toBe('')
    expect((await lista()).value).toBe('2')
    expect(fecha().value).toBe(fechaPropuesta(new Date()))
  })

  it('US1.3: las horas escritas con coma se envían con punto', async () => {
    abrir()
    const user = await completar('1,5')

    await user.click(guardar())

    await screen.findByText('Horas registradas')
    expect(imputarMock.mock.calls[0][1].horas).toBe('1.5')
  })

  it('US1.5: mientras envía deshabilita el guardado, lo indica y no manda una segunda carga', async () => {
    imputarMock.mockReturnValue(new Promise(() => {}))
    abrir()
    const user = await completar()

    await user.click(guardar())
    await user.click(guardar())

    expect(guardar().disabled).toBe(true)
    expect(screen.getByText('Cargando…')).toBeTruthy()
    expect(imputarMock).toHaveBeenCalledTimes(1)
  })

  it('dos formularios abiertos no repiten los id de sus campos', async () => {
    render(
      <>
        <FormularioHoras historiaId={1} proyectoId={42} alRegistrar={vi.fn()} />
        <FormularioHoras historiaId={2} proyectoId={42} alRegistrar={vi.fn()} />
      </>,
    )

    await screen.findAllByLabelText('Integrante')
    const ids = [...document.querySelectorAll('input, select')].map((campo) => campo.id)
    expect(ids).toHaveLength(8)
    expect(new Set(ids).size).toBe(8)
  })
})
