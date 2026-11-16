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

  it('Given diferencias entre los votos, When se muestra el resultado, Then el texto va debajo de la tarjeta y no dentro de ella', () => {
    render(
      <Resultado
        promedio={{ exacto: 5.25, carta_sugerida: 5 }}
        diferencias={{ hay_diferencias: true, minimo: 3, maximo: 8 }}
      />,
    )

    const tarjeta = screen.getByRole('region', { name: 'Resultado de la votación' })
    const texto = screen.getByText('Hay diferencias: mínimo 3, máximo 8')

    expect(tarjeta.contains(texto)).toBe(false)
    expect(tarjeta.compareDocumentPosition(texto) & Node.DOCUMENT_POSITION_FOLLOWING).toBeTruthy()
    expect(texto.className).toBe('votacion-diferencias')
  })

  it('Given votos iguales, When se muestra el resultado, Then dice que no hay diferencias con el mismo estilo', () => {
    render(
      <Resultado
        promedio={{ exacto: 5, carta_sugerida: 5 }}
        diferencias={{ hay_diferencias: false, minimo: null, maximo: null }}
      />,
    )

    expect(screen.getByText('No hay diferencias: todos votaron lo mismo').className).toBe('votacion-diferencias')
  })
})
