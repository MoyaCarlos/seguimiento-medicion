import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { ResumenHistoria } from './ResumenHistoria'
import type { HistoriaResumen } from './votacion'

const HISTORIA: HistoriaResumen = {
  id: 4,
  titulo: 'Pagar la tasa de la licencia en línea',
  descripcion: 'Como solicitante quiero pagar la tasa con tarjeta.\n**sin formato** <b>ni html</b>',
  prioridad: 'M',
  estimacion_sp: 8,
}

describe('ResumenHistoria', () => {
  it('Given la historia encontrada, When se muestra el resumen, Then se ven título, descripción, prioridad y estimación', () => {
    render(<ResumenHistoria historiaId={4} historia={HISTORIA} />)

    expect(screen.getByRole('heading', { name: HISTORIA.titulo })).toBeTruthy()
    expect(screen.getByText(/Como solicitante quiero pagar la tasa con tarjeta/)).toBeTruthy()
    expect(screen.getByText('Must have')).toBeTruthy()
    expect(screen.getByText('8 SP')).toBeTruthy()
  })

  it('Given una descripción con símbolos de Markdown y HTML, When se muestra, Then se ve como texto plano', () => {
    const { container } = render(<ResumenHistoria historiaId={4} historia={HISTORIA} />)

    expect(container.textContent).toContain('**sin formato** <b>ni html</b>')
    expect(container.querySelector('strong, b')).toBeNull()
  })

  it.each([
    ['S', 'Should have'],
    ['C', 'Could have'],
    ['W', "Won't have"],
  ] as const)('Given prioridad %s, When se muestra el resumen, Then la etiqueta es %s', (prioridad, etiqueta) => {
    render(<ResumenHistoria historiaId={4} historia={{ ...HISTORIA, prioridad }} />)

    expect(screen.getByText(etiqueta)).toBeTruthy()
  })

  it('Given una historia sin estimación, When se muestra el resumen, Then figura "Sin estimar"', () => {
    render(<ResumenHistoria historiaId={4} historia={{ ...HISTORIA, estimacion_sp: null }} />)

    expect(screen.getByText('Sin estimar')).toBeTruthy()
  })

  it('Given que la historia no se encontró, When se muestra el resumen, Then solo figura "Historia #4"', () => {
    const { container } = render(<ResumenHistoria historiaId={4} historia={null} />)

    expect(container.textContent).toBe('Historia #4')
  })
})
