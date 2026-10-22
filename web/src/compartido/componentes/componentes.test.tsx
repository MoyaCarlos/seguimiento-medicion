import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { useState } from 'react'
import { describe, expect, it } from 'vitest'
import { ApiError } from '../api/client'
import { Button } from './Button'
import { ErrorMessage } from './ErrorMessage'
import { Field } from './Field'
import { LoadingIndicator } from './LoadingIndicator'

function FormularioDePrueba({ enviar }: { enviar: () => Promise<void> }) {
  const [enCurso, setEnCurso] = useState(false)
  const [error, setError] = useState<ApiError>()

  async function alEnviar(evento: React.FormEvent) {
    evento.preventDefault()
    setEnCurso(true)
    setError(undefined)
    try {
      await enviar()
    } catch (caught) {
      setError(caught as ApiError)
    } finally {
      setEnCurso(false)
    }
  }

  return (
    <form onSubmit={alEnviar} aria-label="Nuevo proyecto">
      {error && !error.field && <ErrorMessage message={error.message} />}
      <Field name="nombre" label="Nombre" error={error?.field === 'nombre' ? error.message : undefined} />
      {enCurso && <LoadingIndicator />}
      <Button type="submit" disabled={enCurso}>Crear</Button>
    </form>
  )
}

describe('Field', () => {
  it('asocia la etiqueta con el campo', async () => {
    render(<Field name="nombre" label="Nombre" />)

    await userEvent.click(screen.getByText('Nombre'))

    expect(screen.getByLabelText('Nombre')).toBe(document.activeElement)
  })

  it('muestra el error junto al campo y lo marca como inválido', () => {
    render(<Field name="nombre" label="Nombre" error="el nombre es obligatorio" />)

    const campo = screen.getByLabelText('Nombre')

    expect(campo.getAttribute('aria-invalid')).toBe('true')
    expect(campo.getAttribute('aria-describedby')).toBe(screen.getByRole('alert').id)
    expect(screen.getByRole('alert').textContent).toBe('el nombre es obligatorio')
  })

  it('no muestra error ni marca el campo como inválido si no hay error', () => {
    render(<Field name="nombre" label="Nombre" />)

    expect(screen.getByLabelText('Nombre').getAttribute('aria-invalid')).toBe('false')
    expect(screen.queryByRole('alert')).toBeNull()
  })
})

describe('Button', () => {
  it('queda deshabilitado cuando se lo indica', () => {
    render(<Button disabled>Crear</Button>)

    expect((screen.getByRole('button', { name: 'Crear' }) as HTMLButtonElement).disabled).toBe(true)
  })
})

describe('LoadingIndicator', () => {
  it('avisa que está cargando', () => {
    render(<LoadingIndicator />)

    expect(screen.getByRole('status').textContent).toBe('Cargando…')
  })
})

describe('formulario con una operación en curso', () => {
  it('muestra el indicador y deshabilita el botón mientras dura y los restablece al terminar', async () => {
    let terminar = () => {}
    const enviar = () => new Promise<void>((resolver) => { terminar = resolver })
    render(<FormularioDePrueba enviar={enviar} />)
    expect(screen.queryByRole('status')).toBeNull()

    await userEvent.click(screen.getByRole('button', { name: 'Crear' }))

    expect(await screen.findByRole('status')).toBeTruthy()
    expect((screen.getByRole('button', { name: 'Crear' }) as HTMLButtonElement).disabled).toBe(true)

    terminar()

    await waitFor(() => expect(screen.queryByRole('status')).toBeNull())
    expect((screen.getByRole('button', { name: 'Crear' }) as HTMLButtonElement).disabled).toBe(false)
  })

  it('muestra un error con campo junto al campo y no como aviso general', async () => {
    const enviar = () => Promise.reject(new ApiError(400, 'el nombre es obligatorio', 'nombre'))
    render(<FormularioDePrueba enviar={enviar} />)

    await userEvent.click(screen.getByRole('button', { name: 'Crear' }))

    const alerta = await screen.findByRole('alert')
    expect(alerta.textContent).toBe('el nombre es obligatorio')
    expect(screen.getByLabelText('Nombre').getAttribute('aria-describedby')).toBe(alerta.id)
  })

  it('muestra un error sin campo como aviso arriba del formulario', async () => {
    const enviar = () => Promise.reject(new ApiError(404, 'proyecto inexistente'))
    render(<FormularioDePrueba enviar={enviar} />)

    await userEvent.click(screen.getByRole('button', { name: 'Crear' }))

    const alerta = await screen.findByRole('alert')
    const formulario = screen.getByRole('form', { name: 'Nuevo proyecto' })
    expect(alerta.textContent).toBe('proyecto inexistente')
    expect(formulario.firstElementChild).toBe(alerta)
    expect(screen.getByLabelText('Nombre').getAttribute('aria-invalid')).toBe('false')
  })
})
