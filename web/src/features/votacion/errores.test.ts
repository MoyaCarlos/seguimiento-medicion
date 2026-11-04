import { describe, expect, it } from 'vitest'
import { ApiError, ConnectionError } from '../../compartido/api/client'
import { esVotacionNoIniciada, mensajeDeError } from './errores'

describe('mensajeDeError', () => {
  it('Given un error de la API, Then devuelve el mensaje del servidor', () => {
    expect(mensajeDeError(new ApiError(403, 'solo el Scrum Master del proyecto puede hacerlo'))).toBe(
      'solo el Scrum Master del proyecto puede hacerlo',
    )
  })

  it('Given un error de conexión, Then devuelve el aviso de conexión', () => {
    expect(mensajeDeError(new ConnectionError())).toBe('No se pudo conectar con el servidor')
  })

  it('Given cualquier otro valor, Then devuelve el mensaje genérico', () => {
    expect(mensajeDeError(new Error('boom'))).toBe('Ocurrió un error inesperado')
    expect(mensajeDeError('texto')).toBe('Ocurrió un error inesperado')
  })
})

describe('esVotacionNoIniciada', () => {
  it('Given un 404 con el mensaje de votación no iniciada, Then es verdadero', () => {
    expect(esVotacionNoIniciada(new ApiError(404, 'la historia no tiene votación iniciada'))).toBe(true)
  })

  it('Given un 404 de historia no encontrada, Then es falso', () => {
    expect(esVotacionNoIniciada(new ApiError(404, 'historia no encontrada'))).toBe(false)
  })

  it('Given el mismo mensaje con otro estado u otro tipo de error, Then es falso', () => {
    expect(esVotacionNoIniciada(new ApiError(400, 'la historia no tiene votación iniciada'))).toBe(false)
    expect(esVotacionNoIniciada(new Error('la historia no tiene votación iniciada'))).toBe(false)
  })
})
