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

describe('CampoMarkdown según el prototipo (guía de estados #/guia/markdown)', () => {
  function EditorConOpcional({ opcional = false, error, inicial = '' }: { opcional?: boolean; error?: string; inicial?: string }) {
    const [valor, setValor] = useState(inicial)
    return (
      <CampoMarkdown
        id="criterios"
        name="criterios"
        label="Criterios de aceptación"
        value={valor}
        onChange={setValor}
        opcional={opcional}
        error={error}
      />
    )
  }

  it('Given las pestañas, When se muestran, Then cada una controla su panel y el panel está rotulado por su pestaña', () => {
    render(<Editor />)

    const escribir = screen.getByRole('tab', { name: 'Escribir' })
    const vista = screen.getByRole('tab', { name: 'Vista previa' })

    expect(document.getElementById(escribir.getAttribute('aria-controls') ?? '')?.getAttribute('role')).toBe('tabpanel')
    expect(document.getElementById(vista.getAttribute('aria-controls') ?? '')?.getAttribute('role')).toBe('tabpanel')
    expect(document.getElementById(escribir.getAttribute('aria-controls') ?? '')?.getAttribute('aria-labelledby')).toBe(escribir.id)
    expect(document.getElementById(vista.getAttribute('aria-controls') ?? '')?.getAttribute('aria-labelledby')).toBe(vista.id)
    expect(screen.getByRole('tabpanel', { name: 'Escribir' })).toBeTruthy()
  })

  it('Given la pestaña elegida, When se muestran las pestañas, Then solo ella entra en el orden del Tab (tabindex 0) y la otra queda en -1', async () => {
    const user = userEvent.setup()
    render(<Editor />)
    const escribir = screen.getByRole('tab', { name: 'Escribir' })
    const vista = screen.getByRole('tab', { name: 'Vista previa' })

    expect([escribir.tabIndex, vista.tabIndex]).toEqual([0, -1])
    await user.click(vista)

    expect([escribir.tabIndex, vista.tabIndex]).toEqual([-1, 0])
  })

  it('Given el foco en una pestaña, When pulso las flechas, Inicio o Fin, Then cambia de pestaña y el foco la sigue', async () => {
    const user = userEvent.setup()
    render(<Editor />)
    const escribir = screen.getByRole('tab', { name: 'Escribir' })
    const vista = screen.getByRole('tab', { name: 'Vista previa' })
    escribir.focus()

    await user.keyboard('{ArrowRight}')
    expect(vista.getAttribute('aria-selected')).toBe('true')
    expect(document.activeElement).toBe(vista)

    await user.keyboard('{ArrowLeft}')
    expect(escribir.getAttribute('aria-selected')).toBe('true')
    expect(document.activeElement).toBe(escribir)

    await user.keyboard('{End}')
    expect(vista.getAttribute('aria-selected')).toBe('true')
    expect(document.activeElement).toBe(vista)

    await user.keyboard('{Home}')
    expect(escribir.getAttribute('aria-selected')).toBe('true')
    expect(document.activeElement).toBe(escribir)
  })

  it('Given el foco en una pestaña, When pulso otra tecla, Then no cambia de pestaña', async () => {
    const user = userEvent.setup()
    render(<Editor />)
    const escribir = screen.getByRole('tab', { name: 'Escribir' })
    escribir.focus()

    await user.keyboard('{ArrowDown}a')

    expect(escribir.getAttribute('aria-selected')).toBe('true')
  })

  it('Given la vista previa, When se muestra, Then la etiqueta deja de apuntar al campo oculto y al volver a Escribir lo vuelve a rotular', async () => {
    const user = userEvent.setup()
    render(<Editor inicial="hola" />)
    const etiqueta = screen.getByText('Descripción', { selector: 'label' })
    expect(etiqueta.getAttribute('for')).toBe('descripcion')

    await user.click(screen.getByRole('tab', { name: 'Vista previa' }))
    expect(etiqueta.hasAttribute('for')).toBe(false)
    expect(screen.getByRole('tabpanel', { name: 'Escribir' }).hasAttribute('hidden')).toBe(true)

    await user.click(screen.getByRole('tab', { name: 'Escribir' }))
    expect(etiqueta.getAttribute('for')).toBe('descripcion')
    expect(campo().value).toBe('hola')
  })

  it('Given Criterios opcionales, When se muestra el campo, Then «Opcional» va junto a la etiqueta pero fuera de ella, y el nombre accesible no cambia', () => {
    render(<EditorConOpcional opcional />)

    const etiqueta = screen.getByText('Criterios de aceptación', { selector: 'label' })
    const marcador = screen.getByText('Opcional')
    const textarea = screen.getByLabelText('Criterios de aceptación')

    expect(etiqueta.textContent).toBe('Criterios de aceptación')
    expect(etiqueta.contains(marcador)).toBe(false)
    expect(etiqueta.parentElement).toBe(marcador.parentElement)
    expect(textarea.getAttribute('aria-describedby')).toBe(`${marcador.id} criterios-ayuda`)
  })

  it('Given Criterios opcionales con un error, When se muestra, Then el campo se describe con el marcador y con el error, y la ayuda sigue visible', () => {
    render(<EditorConOpcional opcional error="Los criterios de aceptación no pueden superar los 2000 caracteres" />)

    const textarea = screen.getByLabelText('Criterios de aceptación')
    const marcador = screen.getByText('Opcional')

    expect(textarea.getAttribute('aria-invalid')).toBe('true')
    expect(textarea.getAttribute('aria-describedby')).toBe(`${marcador.id} criterios-error`)
    expect(screen.getByText(AYUDA)).toBeTruthy()
  })

  it('Given un campo obligatorio, When se muestra, Then no lleva marcador «Opcional»', () => {
    render(<Editor />)

    expect(screen.queryByText('Opcional')).toBeNull()
  })

  it('Given un campo deshabilitado, When se muestra, Then las pestañas siguen activas', () => {
    render(<Editor disabled />)

    expect(screen.getByRole('tab', { name: 'Escribir' }).hasAttribute('disabled')).toBe(false)
    expect(screen.getByRole('tab', { name: 'Vista previa' }).hasAttribute('disabled')).toBe(false)
  })
})

