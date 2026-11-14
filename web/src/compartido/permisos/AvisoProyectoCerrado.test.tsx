import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { AvisoProyectoCerrado } from './AvisoProyectoCerrado'

describe('AvisoProyectoCerrado', () => {
  it('con fecha de cierre muestra «Proyecto cerrado el AAAA-MM-DD» tal como llega, con role status', () => {
    render(<AvisoProyectoCerrado fechaCierre="2026-11-11" />)

    const aviso = screen.getByRole('status')
    expect(aviso.textContent).toBe('Proyecto cerrado el 2026-11-11')
  })

  it('sin fecha de cierre (null) no renderiza nada', () => {
    const { container } = render(<AvisoProyectoCerrado fechaCierre={null} />)

    expect(container.innerHTML).toBe('')
    expect(screen.queryByRole('status')).toBeNull()
  })
})
