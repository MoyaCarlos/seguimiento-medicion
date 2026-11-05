import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { Cartas } from './Cartas'

describe('Cartas', () => {
  it('Given una ronda abierta, When se muestra la baraja, Then hay un botón por cada carta 1, 2, 3, 5, 8, 13 y 21', () => {
    render(<Cartas alElegir={vi.fn()} />)

    const botones = screen.getAllByRole('button').map((boton) => boton.textContent)
    expect(botones).toEqual(['1', '2', '3', '5', '8', '13', '21'])
  })

  it('Given la baraja, When se pulsa la carta 8, Then se llama a alElegir con 8', async () => {
    const alElegir = vi.fn()
    render(<Cartas alElegir={alElegir} />)

    await userEvent.click(screen.getByRole('button', { name: '8' }))

    expect(alElegir).toHaveBeenCalledExactlyOnceWith(8)
  })

  it('Given la baraja deshabilitada, When se pulsa una carta, Then no se llama a alElegir', async () => {
    const alElegir = vi.fn()
    render(<Cartas deshabilitada alElegir={alElegir} />)

    await userEvent.click(screen.getByRole('button', { name: '5' }))

    expect(alElegir).not.toHaveBeenCalled()
  })

  it('Given que mi voto es 5, When se muestra la baraja, Then solo la carta 5 figura seleccionada', () => {
    render(<Cartas seleccionada={5} alElegir={vi.fn()} />)

    expect(screen.getByRole('button', { name: '5' }).getAttribute('aria-pressed')).toBe('true')
    expect(screen.getByRole('button', { name: '8' }).getAttribute('aria-pressed')).toBe('false')
  })
})
