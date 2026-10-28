import { describe, expect, it } from 'vitest'
import { ApiError, ConnectionError } from '../../compartido/api/client'
import { clasificarError } from './errores'

const CAMPOS = ['sprint_goal', 'fecha_inicio', 'fecha_fin']

describe('clasificarError', () => {
  it('asigna al campo un error de la API con un campo del formulario', () => {
    expect(clasificarError(new ApiError(400, 'm', 'sprint_goal'), CAMPOS)).toEqual({
      campo: 'sprint_goal',
      mensaje: 'm',
    })
  })

  it('manda al aviso un error con un campo que el formulario no tiene (estado)', () => {
    expect(clasificarError(new ApiError(400, 'm', 'estado'), CAMPOS)).toEqual({ aviso: 'm' })
  })

  it('manda al aviso un error de la API sin campo', () => {
    expect(clasificarError(new ApiError(404, 'm'), CAMPOS)).toEqual({ aviso: 'm' })
  })

  it('manda al aviso un error con campo proyecto_id', () => {
    expect(clasificarError(new ApiError(400, 'm', 'proyecto_id'), CAMPOS)).toEqual({ aviso: 'm' })
  })

  it('manda al aviso el error de conexión', () => {
    expect(clasificarError(new ConnectionError(), CAMPOS)).toEqual({
      aviso: 'No se pudo conectar con el servidor',
    })
  })

  it('manda al aviso el mensaje de un Error cualquiera', () => {
    expect(clasificarError(new Error('x'), CAMPOS)).toEqual({ aviso: 'x' })
  })

  it('manda al aviso un mensaje genérico si no es un Error', () => {
    expect(clasificarError('algo', CAMPOS)).toEqual({ aviso: 'Ocurrió un error inesperado' })
  })
})
