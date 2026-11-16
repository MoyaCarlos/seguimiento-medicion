import { render, screen, within } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { Integrantes } from './Integrantes'
import type { VotoDeIntegrante } from './votacion'

const voto = (usuario_id: number, nombre: string, ya_voto: boolean, valor: number | null): VotoDeIntegrante => ({
  usuario_id,
  nombre,
  rol: usuario_id === 1 ? 'scrum_master' : 'product_builder',
  ya_voto,
  valor,
})

const fila = (nombre: string) =>
  within(screen.getByRole('list')).getByText(nombre).closest('li') as HTMLElement

describe('Integrantes', () => {
  it('Given una ronda sin revelar, When se muestran los integrantes, Then quien votó figura "Ya votó" con "Oculto" y quien no, "Falta"', () => {
    render(<Integrantes integrantes={[voto(1, 'Beta Uno', true, null), voto(2, 'Beta Dos', false, null)]} />)

    expect(within(fila('Beta Uno')).getByText('Ya votó')).toBeTruthy()
    expect(within(fila('Beta Uno')).getByText('Oculto')).toBeTruthy()
    expect(within(fila('Beta Dos')).getByText('Falta')).toBeTruthy()
    expect(within(fila('Beta Dos')).queryByText('Oculto')).toBeNull()
  })

  it('Given que el integrante consultante votó 8, When se muestran los integrantes, Then se ve su 8 y el de otro sigue "Oculto"', () => {
    render(<Integrantes integrantes={[voto(1, 'Beta Uno', true, 8), voto(2, 'Beta Dos', true, null)]} />)

    expect(within(fila('Beta Uno')).getByText('8')).toBeTruthy()
    expect(within(fila('Beta Dos')).getByText('Oculto')).toBeTruthy()
  })

  it('Given una ronda revelada, When se muestran los integrantes, Then se ve el valor de cada uno', () => {
    render(<Integrantes integrantes={[voto(1, 'Beta Uno', true, 3), voto(2, 'Beta Dos', true, 13)]} />)

    expect(within(fila('Beta Uno')).getByText('3')).toBeTruthy()
    expect(within(fila('Beta Dos')).getByText('13')).toBeTruthy()
  })

  it('Given integrantes con rol, When se muestran, Then cada uno indica si es Scrum Master o Product Builder', () => {
    render(<Integrantes integrantes={[voto(1, 'Beta Uno', false, null), voto(2, 'Beta Dos', false, null)]} />)

    expect(within(fila('Beta Uno')).getByText('Scrum Master')).toBeTruthy()
    expect(within(fila('Beta Dos')).getByText('Product Builder')).toBeTruthy()
  })

  it('Given la lista de integrantes, When la lee un lector de pantalla, Then ninguna fila lleva aria-label y cada una contiene nombre, rol, estado y valor', () => {
    render(<Integrantes integrantes={[voto(2, 'Beta Dos', true, 8)]} />)

    const item = fila('Beta Dos')
    expect(item.hasAttribute('aria-label')).toBe(false)
    expect(item.textContent).toContain('Beta Dos')
    expect(item.textContent).toContain('Product Builder')
    expect(item.textContent).toContain('Ya votó')
    expect(item.textContent).toContain('Voto: 8')
  })
})

describe('Integrantes: mi fila', () => {
  const lista = [voto(1, 'Beta Uno', true, 5), voto(2, 'Beta Dos', true, 8)]

  it('Given el id del usuario actual, When se muestran los integrantes, Then su fila queda resaltada y marcada como actual', () => {
    render(<Integrantes integrantes={lista} usuarioActualId={2} />)

    expect(fila('Beta Dos').getAttribute('aria-current')).toBe('true')
    expect(fila('Beta Dos').className).toBe('votacion-integrante votacion-integrante--yo')
    expect(fila('Beta Uno').hasAttribute('aria-current')).toBe(false)
    expect(fila('Beta Uno').className).toBe('votacion-integrante')
  })

  it('Given que no se conoce el usuario actual, When se muestran los integrantes, Then ninguna fila se resalta', () => {
    render(<Integrantes integrantes={lista} usuarioActualId={null} />)

    expect(document.querySelector('[aria-current]')).toBeNull()
  })
})
