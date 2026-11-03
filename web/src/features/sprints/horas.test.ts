import { describe, expect, it } from 'vitest'
import { fechaPropuesta, formatearHoras, normalizarHoras } from './horas'

describe('normalizarHoras', () => {
  it.each([
    ['1,5', '1.5'],
    ['1.5', '1.5'],
    ['abc', 'abc'],
    ['  4 ', '4'],
    ['', ''],
  ])('convierte %j en %j', (texto, esperado) => {
    expect(normalizarHoras(texto)).toBe(esperado)
  })
})

describe('formatearHoras', () => {
  it.each([
    [5.5, '5,5 h'],
    [4, '4 h'],
    [0, '0 h'],
    [61 / 60, '1,02 h'],
  ])('muestra %d como %s', (horas, esperado) => {
    expect(formatearHoras(horas)).toBe(esperado)
  })
})

describe('fechaPropuesta', () => {
  it('propone hoy si el Sprint no pasó su fecha de fin', () => {
    expect(fechaPropuesta(new Date(2026, 10, 2, 12), '2026-11-09')).toBe('2026-11-02')
  })

  it('propone la fecha de fin si el Sprint ya la pasó', () => {
    expect(fechaPropuesta(new Date(2026, 10, 2, 12), '2026-10-30')).toBe('2026-10-30')
  })

  it('propone hoy si el Sprint no tiene fecha de fin', () => {
    expect(fechaPropuesta(new Date(2026, 10, 2, 12), null)).toBe('2026-11-02')
  })

  it('usa la fecha local aunque en UTC ya sea el día siguiente', () => {
    expect(fechaPropuesta(new Date(2026, 10, 2, 23, 30), null)).toBe('2026-11-02')
  })

  it('completa con ceros el mes y el día', () => {
    expect(fechaPropuesta(new Date(2026, 0, 5, 12), null)).toBe('2026-01-05')
  })
})
