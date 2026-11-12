import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { ApiError, ConnectionError } from '../../compartido/api/client'
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

function renderizarOtroControl() {
  render(<input aria-label="Otro control" />)
  return screen.getByLabelText('Otro control')
}

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

  it('formato h:mm: se muestra la ayuda y el campo de horas está asociado a ella', async () => {
    abrir()
    await lista()

    const ayuda = screen.getByText('Ej.: 1,5 o 1:30 (horas:minutos)')
    expect(ayuda.id).not.toBe('')
    expect(horas().getAttribute('aria-describedby')).toContain(ayuda.id)
  })

  it('formato h:mm: el campo de horas no usa el teclado decimal, que no tiene el signo ":"', async () => {
    abrir()
    await lista()

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

  it('foco: al cargar la lista de integrantes, el foco pasa a ella', async () => {
    abrir()

    const seleccion = await lista()
    await waitFor(() => expect(document.activeElement).toBe(seleccion))
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

  it('foco: si la persona se movió mientras cargaba la lista, no le saca el foco', async () => {
    let responderLista: (lista: typeof INTEGRANTES) => void = () => {}
    integrantesMock.mockReturnValue(new Promise((resolve) => (responderLista = resolve)))
    abrir()
    const otroControl = renderizarOtroControl()

    otroControl.focus()
    responderLista(INTEGRANTES)

    await lista()
    expect(document.activeElement).toBe(otroControl)
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

describe('FormularioHoras: errores', () => {
  const rechazar = (campo: string | undefined, mensaje: string) =>
    imputarMock.mockRejectedValueOnce(new ApiError(400, mensaje, campo))

  it('US2.1: sin integrante muestra el error junto a la lista y no envía la carga', async () => {
    const user = userEvent.setup()
    abrir()
    const seleccion = await lista()
    await user.type(actividad(), 'Revisión del PR')
    await user.type(horas(), '4')

    await user.click(guardar())

    const error = await screen.findByRole('alert')
    expect(error.textContent).toBe('Elegí un integrante')
    expect(seleccion.getAttribute('aria-invalid')).toBe('true')
    expect(seleccion.getAttribute('aria-describedby')).toBe(error.id)
    expect(imputarMock).not.toHaveBeenCalled()
  })

  it('US2.1: un proyecto sin integrantes no permite enviar la carga', async () => {
    integrantesMock.mockResolvedValue([])
    const user = userEvent.setup()
    abrir()
    await lista()

    await user.click(guardar())

    expect((await screen.findByRole('alert')).textContent).toBe('Elegí un integrante')
    expect(imputarMock).not.toHaveBeenCalled()
  })

  it.each([
    ['US2.2', 'horas', 'Horas', 'las horas deben ser un número'],
    ['US2.3', 'actividad', 'Actividad', 'la actividad es obligatoria'],
    ['US2.4', 'fecha', 'Fecha', 'la fecha no puede ser posterior a hoy'],
    ['US2.5', 'integrante_id', 'Integrante', 'el integrante no pertenece al proyecto de la historia'],
  ])('%s: el rechazo de "%s" se muestra junto a su campo', async (_escenario, campo, etiqueta, mensaje) => {
    rechazar(campo, mensaje)
    abrir()
    const user = await completar()

    await user.click(guardar())

    const error = await screen.findByRole('alert')
    expect(error.textContent).toBe(mensaje)
    const control = screen.getByLabelText(etiqueta)
    expect(control.getAttribute('aria-invalid')).toBe('true')
    expect(control.getAttribute('aria-describedby')).toBe(error.id)
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
    expect((await lista()).value).toBe('2')
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

  it('si falla la consulta de integrantes muestra el error y no permite guardar', async () => {
    integrantesMock.mockRejectedValue(new ApiError(404, 'proyecto no encontrado'))
    abrir()

    expect((await screen.findByRole('alert')).textContent).toBe('proyecto no encontrado')
    expect(screen.queryByRole('button', { name: 'Guardar horas' })).toBeNull()
    expect(screen.queryByText('Cargando…')).toBeNull()
  })
})
