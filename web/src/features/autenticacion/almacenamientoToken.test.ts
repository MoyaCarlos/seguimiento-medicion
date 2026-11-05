import { afterEach, describe, expect, it, vi } from 'vitest'
import { borrarToken, guardarToken, leerToken } from './almacenamientoToken'

afterEach(() => {
  vi.restoreAllMocks()
  borrarToken()
  localStorage.clear()
})

describe('almacenamiento del token', () => {
  it('guarda y lee el token', () => {
    guardarToken('abc123')

    expect(leerToken()).toBe('abc123')
  })

  it('devuelve null si no hay token guardado', () => {
    expect(leerToken()).toBeNull()
  })

  it('borra el token guardado', () => {
    guardarToken('abc123')

    borrarToken()

    expect(leerToken()).toBeNull()
  })

  it('guardarToken no lanza si setItem lanza', () => {
    vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => {
      throw new Error('almacenamiento no disponible')
    })

    expect(() => guardarToken('abc123')).not.toThrow()
  })

  it('leerToken devuelve null si getItem lanza', () => {
    vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => {
      throw new Error('almacenamiento no disponible')
    })

    expect(leerToken()).toBeNull()
  })

  it('borrarToken no lanza si removeItem lanza', () => {
    vi.spyOn(Storage.prototype, 'removeItem').mockImplementation(() => {
      throw new Error('almacenamiento no disponible')
    })

    expect(() => borrarToken()).not.toThrow()
  })
})
