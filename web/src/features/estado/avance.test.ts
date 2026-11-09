import { describe, expect, it } from 'vitest'
import { esIdValido, textoDeEstimacion, textoDeSprint } from './avance'

describe('textoDeSprint', () => {
  it('Dado una historia en un Sprint, cuando se muestra, entonces dice "Sprint <número>"', () => {
    expect(textoDeSprint({ id: 20, numero: 2 })).toBe('Sprint 2')
  })

  it('Dado una historia sin Sprint, cuando se muestra, entonces dice "Product Backlog"', () => {
    expect(textoDeSprint(null)).toBe('Product Backlog')
  })
})

describe('textoDeEstimacion', () => {
  it('Dado una historia estimada, cuando se muestra, entonces dice sus Story Points', () => {
    expect(textoDeEstimacion(8)).toBe('8')
  })

  it('Dado una historia sin Story Points, cuando se muestra, entonces dice "sin estimar" y nunca 0', () => {
    expect(textoDeEstimacion(null)).toBe('sin estimar')
    expect(textoDeEstimacion(null)).not.toBe('0')
  })
})

describe('esIdValido', () => {
  it.each(['1', '7', '25', '1200'])(
    'Dado el id de ruta "%s", cuando se valida, entonces es un entero positivo válido',
    (id) => {
      expect(esIdValido(id)).toBe(true)
    },
  )

  it.each(['', 'abc', '0', '-1', '1.5', '007', '1e3', ' 7', '7 '])(
    'Dado el id de ruta "%s", cuando se valida, entonces no es válido',
    (id) => {
      expect(esIdValido(id)).toBe(false)
    },
  )
})
