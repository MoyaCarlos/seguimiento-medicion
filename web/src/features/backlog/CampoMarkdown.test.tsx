import { useState } from 'react'
import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { CampoMarkdown } from './CampoMarkdown'

const AYUDA = /Admite Markdown: negrita, listas, enlaces e imágenes por URL/

function Editor({
  inicial = '',
  error,
  disabled,
}: {
  inicial?: string
  error?: string
  disabled?: boolean
}) {
  const [valor, setValor] = useState(inicial)
  return (
    <>
      <CampoMarkdown
        id="descripcion"
        name="descripcion"
        label="Descripción"
        value={valor}
        onChange={setValor}
        error={error}
        disabled={disabled}
      />
      <output data-testid="valor">{valor}</output>
    </>
  )
}

const campo = () => screen.getByLabelText('Descripción') as HTMLTextAreaElement
const valor = () => screen.getByTestId('valor').textContent

function seleccionar(inicio: number, fin: number) {
  const textarea = campo()
  textarea.focus()
  textarea.setSelectionRange(inicio, fin)
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('CampoMarkdown', () => {
  it('Given el editor, When se muestra, Then muestra la ayuda, la barra y las pestañas', () => {
    render(<Editor />)

    expect(screen.getByLabelText('Descripción')).toBeTruthy()
    expect(screen.getByText(AYUDA)).toBeTruthy()
    expect(screen.getByRole('toolbar', { name: 'Formato Markdown' })).toBeTruthy()
    expect(screen.getByRole('tab', { name: 'Escribir' })).toBeTruthy()
    expect(screen.getByRole('tab', { name: 'Vista previa' })).toBeTruthy()
  })

  it('Given un texto seleccionado, When pulso Negrita, Then la selección queda entre **', async () => {
    const user = userEvent.setup()
    render(<Editor inicial="texto" />)
    seleccionar(0, 5)

    await user.click(screen.getByRole('button', { name: 'Negrita' }))

    expect(campo().value).toBe('**texto**')
    expect(valor()).toBe('**texto**')
  })

  it.each([
    ['Cursiva', 'texto', '*texto*'],
    ['Código', 'texto', '`texto`'],
    ['Enlace', 'texto', '[texto](https://ejemplo.com)'],
    ['Imagen', 'texto', '![texto](https://ejemplo.com/imagen.png)'],
  ])('Given una selección, When pulso %s, Then inserta %s', async (boton, inicial, esperado) => {
    const user = userEvent.setup()
    render(<Editor inicial={inicial} />)
    seleccionar(0, inicial.length)

    await user.click(screen.getByRole('button', { name: boton }))

    expect(campo().value).toBe(esperado)
  })

  it('Given una línea, When pulso Título, Then la prefija con "# "', async () => {
    const user = userEvent.setup()
    render(<Editor inicial="linea" />)
    seleccionar(0, 5)

    await user.click(screen.getByRole('button', { name: 'Título' }))

    expect(campo().value).toBe('# linea')
  })

  it('Given dos líneas seleccionadas, When pulso Lista, Then prefija cada una con "- "', async () => {
    const user = userEvent.setup()
    render(<Editor inicial={'uno\ndos'} />)
    seleccionar(0, 7)

    await user.click(screen.getByRole('button', { name: 'Lista' }))

    expect(campo().value).toBe('- uno\n- dos')
  })

  it('Given Markdown escrito, When abro Vista previa, Then lo renderiza y conserva el texto', async () => {
    const user = userEvent.setup()
    render(<Editor inicial="# Título" />)

    await user.click(screen.getByRole('tab', { name: 'Vista previa' }))

    expect(screen.getByRole('heading', { level: 1, name: 'Título' })).toBeTruthy()
    expect(screen.queryByLabelText('Descripción')).toBeNull()

    await user.click(screen.getByRole('tab', { name: 'Escribir' }))

    expect(campo().value).toBe('# Título')
  })

  it('Given el editor vacío, When abro Vista previa, Then avisa que no hay contenido', async () => {
    const user = userEvent.setup()
    render(<Editor />)

    await user.click(screen.getByRole('tab', { name: 'Vista previa' }))

    expect(screen.getByText('Nada para previsualizar')).toBeTruthy()
  })

  it('Given un error, When se muestra, Then lo asocia al campo', () => {
    render(<Editor error="La descripción es obligatoria" />)

    const textarea = campo()
    expect(textarea.getAttribute('aria-invalid')).toBe('true')
    expect(textarea.getAttribute('aria-describedby')).toBe('descripcion-error')
    expect(screen.getByText('La descripción es obligatoria')).toBeTruthy()
  })

  it('Given un campo deshabilitado, When se muestra, Then el campo y la barra no permiten editar', () => {
    render(<Editor disabled />)

    expect(campo().disabled).toBe(true)
    const barra = screen.getByRole('toolbar', { name: 'Formato Markdown' })
    for (const boton of within(barra).getAllByRole('button')) {
      expect(boton.hasAttribute('disabled')).toBe(true)
    }
  })
})
