import { describe, expect, it } from 'vitest'
import { formatoNumero, textoDeDesviacion } from './numeros'

describe('formatoNumero', () => {
  it.each([
    [-6.25, '-6,3'],
    [83.33, '83,3'],
    [-0.04, '0'],
    [4, '4'],
    [0, '0'],
    [1.5, '1,5'],
  ])('%d se muestra como %s', (valor, esperado) => {
    expect(formatoNumero(valor)).toBe(esperado)
  })

  it('null se muestra como sin datos', () => {
    expect(formatoNumero(null)).toBe('sin datos')
  })
})

describe('textoDeDesviacion', () => {
  it.each([
    [-32, -27.6, '-32 h (-27,6 %)'],
    [2, 5, '+2 h (+5 %)'],
    [0, 0, '0 h (0 %)'],
    [-0.04, -0.04, '0 h (0 %)'],
    [3, null, '+3 h (sin datos)'],
  ])('%d h y %s %% se muestran como %s', (horas, porcentaje, esperado) => {
    expect(textoDeDesviacion(horas, porcentaje)).toBe(esperado)
  })

  it('sin horas dice sin datos', () => {
    expect(textoDeDesviacion(null, null)).toBe('sin datos')
  })
})
