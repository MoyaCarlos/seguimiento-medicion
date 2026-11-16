import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { AvisoProyectoCerrado } from './AvisoProyectoCerrado'

describe('AvisoProyectoCerrado', () => {
  it('con fecha de cierre muestra «Proyecto cerrado el dd/mm/aaaa» (la fecha ISO que llega, en formato local), con role status', () => {
    render(<AvisoProyectoCerrado fechaCierre="2026-11-11" />)

    const aviso = screen.getByRole('status')
    expect(aviso.textContent).toBe('Proyecto cerrado el 11/11/2026')
  })

  it('sin fecha de cierre (null) no renderiza nada', () => {
    const { container } = render(<AvisoProyectoCerrado fechaCierre={null} />)

    expect(container.innerHTML).toBe('')
    expect(screen.queryByRole('status')).toBeNull()
  })

  it.each([
    ['día y mes de un solo dígito conservan el cero', '2026-01-05', '05/01/2026'],
    ['último día del año', '2026-12-31', '31/12/2026'],
    ['primer día del año', '2027-01-01', '01/01/2027'],
    ['una marca de tiempo completa solo muestra la fecha', '2026-03-09T14:30:00Z', '09/03/2026'],
  ])('Given una fecha de cierre con %s, When se muestra el aviso, Then dice «Proyecto cerrado el %s»', (_caso, iso, esperado) => {
    render(<AvisoProyectoCerrado fechaCierre={iso} />)

    expect(screen.getByRole('status').textContent).toBe(`Proyecto cerrado el ${esperado}`)
  })
})

