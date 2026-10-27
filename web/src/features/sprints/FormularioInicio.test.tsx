import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { FormularioInicio } from './FormularioInicio'
import type { Sprint } from './sprint'

const pendiente: Sprint = { id: 7, proyecto_id: 1, estado: 'Pendiente' }
const activo: Sprint = {
  ...pendiente,
  estado: 'Activo',
  sprint_goal: 'Meta',
  fecha_inicio: '2026-10-01',
  fecha_fin: '2026-10-15',
}

const respuesta = (status: number, cuerpo: unknown) =>
  new Response(JSON.stringify(cuerpo), { status })

async function completarYEnviar() {
  const user = userEvent.setup()
  await user.type(screen.getByLabelText('Sprint Goal'), 'Meta')
  await user.type(screen.getByLabelText('Fecha de inicio'), '2026-10-01')
  await user.type(screen.getByLabelText('Fecha de fin'), '2026-10-15')
  await user.click(screen.getByRole('button', { name: 'Iniciar Sprint' }))
}

afterEach(() => vi.unstubAllGlobals())

describe('FormularioInicio', () => {
  it('ofrece Goal y fechas asociados a su etiqueta y el botón de iniciar', () => {
    render(<FormularioInicio sprint={pendiente} alActualizar={vi.fn()} />)

    expect(screen.getByLabelText('Sprint Goal')).toBeTruthy()
    expect(screen.getByLabelText('Fecha de inicio').getAttribute('type')).toBe('date')
    expect(screen.getByLabelText('Fecha de fin').getAttribute('type')).toBe('date')
    expect(screen.getByRole('button', { name: 'Iniciar Sprint' })).toBeTruthy()
  })

  it('envía los datos a /sprints/{id}/iniciar y entrega el Sprint devuelto', async () => {
    const fetchSimulado = vi.fn().mockResolvedValue(respuesta(200, activo))
    vi.stubGlobal('fetch', fetchSimulado)
    const alActualizar = vi.fn()
    render(<FormularioInicio sprint={pendiente} alActualizar={alActualizar} />)

    await completarYEnviar()

    await waitFor(() => expect(alActualizar).toHaveBeenCalledWith(activo))
    const [ruta, opciones] = fetchSimulado.mock.calls[0]
    expect(ruta).toBe('/sprints/7/iniciar')
    expect(opciones.method).toBe('POST')
    expect(JSON.parse(opciones.body)).toEqual({
      sprint_goal: 'Meta',
      fecha_inicio: '2026-10-01',
      fecha_fin: '2026-10-15',
    })
  })

  it('deshabilita el botón y muestra la carga mientras envía', async () => {
    vi.stubGlobal('fetch', vi.fn().mockReturnValue(new Promise(() => {})))
    render(<FormularioInicio sprint={pendiente} alActualizar={vi.fn()} />)

    await completarYEnviar()

    expect(screen.getByText('Cargando…')).toBeTruthy()
    expect((screen.getByRole('button', { name: 'Iniciar Sprint' }) as HTMLButtonElement).disabled).toBe(true)
  })

  it('muestra junto al campo el error de un campo del formulario', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue(respuesta(400, { campo: 'fecha_fin', mensaje: 'fin inválido' })),
    )
    render(<FormularioInicio sprint={pendiente} alActualizar={vi.fn()} />)

    await completarYEnviar()

    const campo = await screen.findByLabelText('Fecha de fin')
    expect(campo.getAttribute('aria-invalid')).toBe('true')
    expect(document.getElementById(campo.getAttribute('aria-describedby')!)?.textContent).toBe('fin inválido')
  })

  it('muestra en un aviso un error con un campo que el formulario no tiene', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue(respuesta(400, { campo: 'estado', mensaje: 'ya hay un Sprint activo' })),
    )
    render(<FormularioInicio sprint={pendiente} alActualizar={vi.fn()} />)

    await completarYEnviar()

    expect((await screen.findByRole('alert')).textContent).toBe('ya hay un Sprint activo')
    expect(screen.getByLabelText('Fecha de fin').getAttribute('aria-invalid')).toBe('false')
  })

  it('muestra en el aviso el error de conexión', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new TypeError('red')))
    render(<FormularioInicio sprint={pendiente} alActualizar={vi.fn()} />)

    await completarYEnviar()

    expect((await screen.findByRole('alert')).textContent).toBe('No se pudo conectar con el servidor')
  })

  it('un nuevo envío borra los errores anteriores', async () => {
    const fetchSimulado = vi
      .fn()
      .mockResolvedValueOnce(respuesta(400, { campo: 'fecha_fin', mensaje: 'fin inválido' }))
      .mockResolvedValueOnce(respuesta(400, { campo: 'estado', mensaje: 'otro error' }))
    vi.stubGlobal('fetch', fetchSimulado)
    render(<FormularioInicio sprint={pendiente} alActualizar={vi.fn()} />)

    await completarYEnviar()
    await screen.findByText('fin inválido')
    await userEvent.click(screen.getByRole('button', { name: 'Iniciar Sprint' }))

    expect(await screen.findByText('otro error')).toBeTruthy()
    expect(screen.queryByText('fin inválido')).toBeNull()
  })
})
