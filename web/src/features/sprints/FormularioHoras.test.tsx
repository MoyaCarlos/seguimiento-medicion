import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { ApiError, ConnectionError } from '../../compartido/api/client'
import { FormularioHoras } from './FormularioHoras'
import { imputarHoras } from './api'
import { fechaPropuesta } from './horas'

vi.mock('./api', () => ({ imputarHoras: vi.fn() }))

const imputarMock = vi.mocked(imputarHoras)

function abrir(fechaFin?: string) {
  const alRegistrar = vi.fn()
  render(<FormularioHoras historiaId={12} fechaFin={fechaFin} alRegistrar={alRegistrar} />)
  return { alRegistrar }
}

const fecha = () => screen.getByLabelText('Fecha') as HTMLInputElement
const actividad = () => screen.getByLabelText('Actividad') as HTMLInputElement
const horas = () => screen.getByLabelText('Horas') as HTMLInputElement
const guardar = () => screen.getByRole('button', { name: 'Guardar horas' }) as HTMLButtonElement

function renderizarOtroControl() {
  render(<input aria-label="Otro control" />)
  return screen.getByLabelText('Otro control')
}

async function completar(texto = '4') {
  const user = userEvent.setup()
  await user.type(actividad(), 'Revisión del PR')
  await user.type(horas(), texto)
  return user
}

beforeEach(() => {
  imputarMock.mockResolvedValue({})
})

afterEach(() => {
  vi.resetAllMocks()
})

describe('FormularioHoras: registrar horas', () => {

  it('US1.1: muestra los campos de fecha, actividad y horas, sin elegir integrante', () => {
    abrir()

    expect(screen.queryByLabelText('Integrante')).toBeNull()
    expect(fecha().type).toBe('date')
    expect(actividad()).toBeTruthy()
    expect(horas().type).toBe('text')
  })

  it('US1.1: propone la fecha de hoy', () => {
    abrir()

    expect(fecha().value).toBe(fechaPropuesta(new Date()))
  })

  it('US1.1: si el Sprint ya pasó su fecha de fin, propone esa fecha', () => {
    abrir('2020-01-31')

    expect(fecha().value).toBe('2020-01-31')
  })

  it('US1.2: al guardar envía los datos elegidos, confirma y avisa', async () => {
    const { alRegistrar } = abrir()
    const user = await completar()

    await user.click(guardar())

    expect((await screen.findByRole('status')).textContent).toBe('Horas registradas')
    expect(imputarMock).toHaveBeenCalledWith(12, {
      fecha: fechaPropuesta(new Date()),
      actividad: 'Revisión del PR',
      horas: '4',
    })
    expect(alRegistrar).toHaveBeenCalledTimes(1)
  })

  it('US1.2: al guardar vacía actividad y horas, y conserva la fecha', async () => {
    abrir()
    const user = await completar()

    await user.click(guardar())

    await screen.findByText('Horas registradas')
    expect(actividad().value).toBe('')
    expect(horas().value).toBe('')
    expect(fecha().value).toBe(fechaPropuesta(new Date()))
  })

  it('US1.3: las horas escritas con coma se envían con punto', async () => {
    abrir()
    const user = await completar('1,5')

    await user.click(guardar())

    await screen.findByText('Horas registradas')
    expect(imputarMock.mock.calls[0][1].horas).toBe('1.5')
  })

  it('formato h:mm: se muestra la ayuda y el campo de horas está asociado a ella', async () => {
    abrir()

    const ayuda = screen.getByText('Ej.: 1,5 o 1:30 (horas:minutos)')
    expect(ayuda.id).not.toBe('')
    expect(horas().getAttribute('aria-describedby')).toContain(ayuda.id)
  })

  it('formato h:mm: el campo de horas no usa el teclado decimal, que no tiene el signo ":"', async () => {
    abrir()

    expect(horas().getAttribute('inputmode')).not.toBe('decimal')
    expect(horas().getAttribute('inputmode')).toBe('text')
  })

  it('formato h:mm: las horas escritas como 0:20 se envían sin cambios', async () => {
    abrir()
    const user = await completar('0:20')

    await user.click(guardar())

    await screen.findByText('Horas registradas')
    expect(imputarMock.mock.calls[0][1].horas).toBe('0:20')
  })

  it('formato h:mm: las horas escritas con coma siguen enviándose con punto', async () => {
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

  describe('foco: al terminar el envío vuelve a donde estaba', () => {
    function enviarConRespuestaPendiente() {
      let responder: (valor: unknown) => void = () => {}
      let rechazarEnvio: (error: unknown) => void = () => {}
      imputarMock.mockReturnValue(
        new Promise((resolve, reject) => {
          responder = resolve
          rechazarEnvio = reject
        }),
      )
      return { responder: (valor: unknown) => responder(valor), rechazar: (error: unknown) => rechazarEnvio(error) }
    }

    // jsdom no saca el foco de un control deshabilitado, ni con blur(); el navegador sí.
    function perderFocoComoElNavegador() {
      const temporal = document.body.appendChild(document.createElement('button'))
      temporal.focus()
      temporal.remove()
    }

    it('si se guardó con el botón, vuelve al botón', async () => {
      const envio = enviarConRespuestaPendiente()
      abrir()
      const user = await completar()

      await user.click(guardar())
      perderFocoComoElNavegador()
      envio.responder({})

      await screen.findByText('Horas registradas')
      await waitFor(() => expect(document.activeElement).toBe(guardar()))
    })

    it('si se guardó con Enter y se rechaza, vuelve al campo', async () => {
      const envio = enviarConRespuestaPendiente()
      abrir()
      const user = await completar()

      await user.type(horas(), '{Enter}')
      perderFocoComoElNavegador()
      envio.rechazar(new ApiError(400, 'las horas deben ser un número', 'horas'))

      await screen.findByRole('alert')
      await waitFor(() => expect(document.activeElement).toBe(horas()))
    })

    it('si la persona ya se movió a otro control, no le saca el foco', async () => {
      const envio = enviarConRespuestaPendiente()
      abrir()
      const user = await completar()
      const otroControl = renderizarOtroControl()

      await user.click(guardar())
      otroControl.focus()
      envio.responder({})

      await screen.findByText('Horas registradas')
      expect(document.activeElement).toBe(otroControl)
    })
  })

  it('dos formularios abiertos no repiten los id de sus campos', async () => {
    render(
      <>
        <FormularioHoras historiaId={1} alRegistrar={vi.fn()} />
        <FormularioHoras historiaId={2} alRegistrar={vi.fn()} />
      </>,
    )

    const ids = [...document.querySelectorAll('input, select')].map((campo) => campo.id)
    expect(ids).toHaveLength(6)
    expect(new Set(ids).size).toBe(6)
  })
})

describe('FormularioHoras: errores', () => {
  const rechazar = (campo: string | undefined, mensaje: string) =>
    imputarMock.mockRejectedValueOnce(new ApiError(400, mensaje, campo))

  it.each([
    ['US2.2', 'horas', 'Horas', 'las horas deben ser un número'],
    ['US2.3', 'actividad', 'Actividad', 'la actividad es obligatoria'],
    ['US2.4', 'fecha', 'Fecha', 'la fecha no puede ser posterior a hoy'],
  ])('%s: el rechazo de "%s" se muestra junto a su campo', async (_escenario, campo, etiqueta, mensaje) => {
    rechazar(campo, mensaje)
    abrir()
    const user = await completar()

    await user.click(guardar())

    const error = await screen.findByRole('alert')
    expect(error.textContent).toBe(mensaje)
    const control = screen.getByLabelText(etiqueta)
    expect(control.getAttribute('aria-invalid')).toBe('true')
    expect(control.getAttribute('aria-describedby')?.split(' ')).toContain(error.id)
  })

  it('US2.6: si el Sprint ya se cerró, el mensaje es un aviso general', async () => {
    rechazar('historia_id', 'solo se imputan horas a historias del Sprint activo')
    abrir()
    const user = await completar()

    await user.click(guardar())

    const error = await screen.findByRole('alert')
    expect(error.textContent).toBe('solo se imputan horas a historias del Sprint activo')
    expect(error.id).toBe('')
  })

  it('US2.7: ante un rechazo conserva lo escrito y no avisa que se registró', async () => {
    rechazar('horas', 'las horas deben ser un número')
    const { alRegistrar } = abrir()
    const user = await completar('abc')

    await user.click(guardar())

    await screen.findByRole('alert')
    expect(actividad().value).toBe('Revisión del PR')
    expect(horas().value).toBe('abc')
    expect(screen.queryByText('Horas registradas')).toBeNull()
    expect(alRegistrar).not.toHaveBeenCalled()
  })

  it('US2.8: el error desaparece cuando la carga se acepta', async () => {
    rechazar('horas', 'las horas deben ser un número')
    abrir()
    const user = await completar('abc')
    await user.click(guardar())
    await screen.findByRole('alert')

    await user.clear(horas())
    await user.type(horas(), '4')
    await user.click(guardar())

    await screen.findByText('Horas registradas')
    expect(screen.queryByRole('alert')).toBeNull()
    expect(horas().getAttribute('aria-invalid')).toBe('false')
  })

  it('US2.8: el aviso general desaparece cuando la carga se acepta', async () => {
    rechazar('historia_id', 'solo se imputan horas a historias del Sprint activo')
    abrir()
    const user = await completar()
    await user.click(guardar())
    await screen.findByRole('alert')

    await user.click(guardar())

    await screen.findByText('Horas registradas')
    expect(screen.queryByRole('alert')).toBeNull()
  })

  it('sin conexión al guardar muestra el aviso y conserva lo escrito', async () => {
    imputarMock.mockRejectedValueOnce(new ConnectionError())
    abrir()
    const user = await completar()

    await user.click(guardar())

    expect((await screen.findByRole('alert')).textContent).toBe('No se pudo conectar con el servidor')
    expect(horas().value).toBe('4')
  })
})

describe('FormularioHoras: la sesión es quien registra (HU-21 Escenario 7)', () => {
  it('Escenario 7: horas sin elegir integrante, sin selector «Integrante» y sin integrante_id en el pedido', async () => {
    abrir()
    const user = userEvent.setup()

    await screen.findByLabelText('Horas')
    expect(screen.queryByLabelText('Integrante')).toBeNull()

    await user.type(actividad(), 'Revisión del PR')
    await user.type(horas(), '1:30')
    await user.click(guardar())

    await waitFor(() =>
      expect(imputarMock).toHaveBeenCalledWith(12, {
        fecha: fechaPropuesta(new Date()),
        actividad: 'Revisión del PR',
        horas: '1:30',
      }),
    )
  })
})
