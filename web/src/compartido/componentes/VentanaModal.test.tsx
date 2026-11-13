import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { useState } from 'react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { VentanaModal } from './VentanaModal'

function Escena({ alCerrar = vi.fn() }: { alCerrar?: () => void }) {
  const [abierta, setAbierta] = useState(false)

  return (
    <>
      <button type="button" onClick={() => setAbierta(true)}>Abrir</button>
      {abierta && (
        <VentanaModal
          alCerrar={() => {
            alCerrar()
            setAbierta(false)
          }}
        >
          <h2 id="titulo-de-la-ventana">Crear algo</h2>
          <input aria-label="Primero" />
          <button type="button">Último</button>
        </VentanaModal>
      )}
    </>
  )
}

async function abrir(alCerrar?: () => void) {
  const usuario = userEvent.setup()
  render(<Escena alCerrar={alCerrar} />)
  await usuario.click(screen.getByRole('button', { name: 'Abrir' }))
  return usuario
}

afterEach(() => {
  document.body.style.overflow = ''
})

describe('VentanaModal', () => {
  it('es un diálogo modal con el nombre de su título', async () => {
    await abrir()

    const dialogo = screen.getByRole('dialog', { name: 'Crear algo' })
    expect(dialogo.getAttribute('aria-modal')).toBe('true')
  })

  it('al abrirse lleva el foco al primer control del contenido, no al botón de cerrar', async () => {
    await abrir()

    expect(document.activeElement).toBe(screen.getByLabelText('Primero'))
  })

  it('Escape la cierra', async () => {
    const alCerrar = vi.fn()
    const usuario = await abrir(alCerrar)

    await usuario.keyboard('{Escape}')

    expect(alCerrar).toHaveBeenCalledTimes(1)
  })

  it('el botón «Cerrar» la cierra', async () => {
    const alCerrar = vi.fn()
    const usuario = await abrir(alCerrar)

    await usuario.click(screen.getByRole('button', { name: 'Cerrar' }))

    expect(alCerrar).toHaveBeenCalledTimes(1)
  })

  it('un clic en el fondo la cierra y uno dentro de la ventana no', async () => {
    const alCerrar = vi.fn()
    const usuario = await abrir(alCerrar)

    await usuario.click(screen.getByLabelText('Primero'))
    expect(alCerrar).not.toHaveBeenCalled()
    await usuario.click(screen.getByRole('dialog').parentElement as HTMLElement)

    expect(alCerrar).toHaveBeenCalledTimes(1)
  })

  it('el foco queda atrapado: Tab desde el último control vuelve al primero y Mayús+Tab desde el primero va al último', async () => {
    const usuario = await abrir()
    const botonCerrar = screen.getByRole('button', { name: 'Cerrar' })

    botonCerrar.focus()
    await usuario.tab()
    expect(document.activeElement).toBe(screen.getByLabelText('Primero'))

    await usuario.tab({ shift: true })
    expect(document.activeElement).toBe(botonCerrar)
  })

  it('al cerrarse devuelve el foco al botón que la abrió', async () => {
    const usuario = await abrir()

    await usuario.keyboard('{Escape}')

    expect(document.activeElement).toBe(screen.getByRole('button', { name: 'Abrir' }))
  })

  it('mientras está abierta bloquea el scroll de la página y al cerrarse lo restablece', async () => {
    const usuario = await abrir()
    expect(document.body.style.overflow).toBe('hidden')

    await usuario.keyboard('{Escape}')

    expect(document.body.style.overflow).toBe('')
  })
})
