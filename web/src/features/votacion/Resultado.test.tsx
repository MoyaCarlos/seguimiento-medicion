import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { Resultado } from './Resultado'

describe('Resultado', () => {
  it('Given una ronda revelada con promedio 5,33 y carta 5, When se muestra el resultado, Then se ven el promedio exacto y la carta sugerida', () => {
    render(<Resultado promedio={{ exacto: 5.33, carta_sugerida: 5 }} />)

    expect(screen.getByText('Promedio exacto').nextSibling?.textContent).toBe('5,33')
    expect(screen.getByText('Carta sugerida').nextSibling?.textContent).toBe('5')
  })

  it('Given que todos votaron la misma carta, When se muestra el resultado, Then el promedio no lleva decimales', () => {
    render(<Resultado promedio={{ exacto: 8, carta_sugerida: 8 }} />)

    expect(screen.getByText('Promedio exacto').nextSibling?.textContent).toBe('8')
  })

  it('Given una ronda sin promedio, When se muestra el resultado, Then no se muestra nada', () => {
    const { container } = render(<Resultado promedio={null} />)

    expect(container.textContent).toBe('')
  })
})
