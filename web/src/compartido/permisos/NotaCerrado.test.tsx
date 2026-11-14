import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { NotaCerrado } from './NotaCerrado'

describe('NotaCerrado', () => {
  it('muestra «El proyecto está cerrado» con el id recibido para aria-describedby', () => {
    render(<NotaCerrado id="nota-cierre" />)

    const nota = screen.getByText('El proyecto está cerrado')
    expect(nota.id).toBe('nota-cierre')
  })
})
