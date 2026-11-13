import { describe, expect, it } from 'vitest'
import { fechaCorta, periodoCorto } from './fechas'

describe('fechaCorta', () => {
  it.each([
    ['2026-10-14', '14/10/2026'],
    ['2026-12-28T00:00:00Z', '28/12/2026'],
  ])('%s se muestra como %s', (iso, esperado) => {
    expect(fechaCorta(iso)).toBe(esperado)
  })

  it.each([null, ''])('%j se muestra como «Sin definir»', (valor) => {
    expect(fechaCorta(valor)).toBe('Sin definir')
  })
})

describe('periodoCorto', () => {
  it('une inicio y fin con una flecha', () => {
    expect(periodoCorto('2026-10-14', '2026-12-28')).toBe('14/10/2026 → 28/12/2026')
  })

  it('un extremo faltante dice «Sin definir» sin ocultar el otro', () => {
    expect(periodoCorto(null, '2026-12-28')).toBe('Sin definir → 28/12/2026')
  })
})
