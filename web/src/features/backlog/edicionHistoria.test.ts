import { describe, expect, it } from 'vitest'
import { armarCuerpoEdicion, formularioDe, validarEdicion } from './edicionHistoria'
import type { Historia } from './historia'

const HISTORIA: Historia = {
  id: 7,
  proyecto_id: 42,
  titulo: 'Como usuario quiero iniciar sesión',
  descripcion: 'Autenticación con email y contraseña',
  criterios_aceptacion: 'Dado un email válido, cuando ingreso, entonces veo el tablero',
  prioridad: 'M',
  estado: 'Nueva',
  valor_negocio: 13,
  estimacion_sp: 5,
  sprint_id: null,
}

const VALIDO = {
  titulo: 'Título válido',
  descripcion: 'Descripción válida',
  criterios: '',
  estimacion: '' as const,
}

describe('formularioDe', () => {
  it('carga los cuatro campos editable y deja vacíos los valores nulos', () => {
    expect(formularioDe(HISTORIA)).toEqual({
      titulo: HISTORIA.titulo,
      descripcion: HISTORIA.descripcion,
      criterios: HISTORIA.criterios_aceptacion,
      estimacion: '5',
    })

    expect(
      formularioDe({
        ...HISTORIA,
        criterios_aceptacion: null,
        estimacion_sp: null,
      }),
    ).toEqual({
      titulo: HISTORIA.titulo,
      descripcion: HISTORIA.descripcion,
      criterios: '',
      estimacion: '',
    })
  })
})

describe('armarCuerpoEdicion', () => {
  it('recorta los textos y conserva los Story Points elegidos', () => {
    expect(
      armarCuerpoEdicion({
        titulo: '  Título  ',
        descripcion: '  Descripción  ',
        criterios: '  Criterio  ',
        estimacion: '8',
      }),
    ).toEqual({
      titulo: 'Título',
      descripcion: 'Descripción',
      criterios_aceptacion: 'Criterio',
      estimacion_sp: 8,
    })
  })

  it('envía null cuando los criterios o los Story Points quedan vacíos', () => {
    expect(
      armarCuerpoEdicion({
        titulo: 'Título',
        descripcion: 'Descripción',
        criterios: '   ',
        estimacion: '',
      }),
    ).toEqual({
      titulo: 'Título',
      descripcion: 'Descripción',
      criterios_aceptacion: null,
      estimacion_sp: null,
    })
  })
})

describe('validarEdicion', () => {
  it('acepta una edición válida y los límites exactos', () => {
    expect(validarEdicion(VALIDO)).toEqual({})
    expect(
      validarEdicion({
        titulo: 'a'.repeat(200),
        descripcion: 'a'.repeat(2000),
        criterios: 'a'.repeat(2000),
        estimacion: '21',
      }),
    ).toEqual({})
  })

  it('exige el título y rechaza los que solo tienen espacios', () => {
    expect(validarEdicion({ ...VALIDO, titulo: '' }).titulo).toBe('El título es obligatorio')
    expect(validarEdicion({ ...VALIDO, titulo: '   ' }).titulo).toBe('El título es obligatorio')
  })

  it('rechaza un título de más de 200 caracteres', () => {
    expect(validarEdicion({ ...VALIDO, titulo: 'a'.repeat(201) }).titulo).toBe(
      'El título no puede superar los 200 caracteres',
    )
  })

  it('exige la descripción y rechaza las que solo tienen espacios', () => {
    expect(validarEdicion({ ...VALIDO, descripcion: '' }).descripcion).toBe(
      'La descripción es obligatoria',
    )
    expect(validarEdicion({ ...VALIDO, descripcion: '   ' }).descripcion).toBe(
      'La descripción es obligatoria',
    )
  })

  it('rechaza una descripción de más de 2000 caracteres', () => {
    expect(validarEdicion({ ...VALIDO, descripcion: 'a'.repeat(2001) }).descripcion).toBe(
      'La descripción no puede superar los 2000 caracteres',
    )
  })

  it('acepta criterios vacíos y rechaza los de más de 2000 caracteres', () => {
    expect(validarEdicion({ ...VALIDO, criterios: '   ' }).criterios).toBeUndefined()
    expect(validarEdicion({ ...VALIDO, criterios: 'a'.repeat(2001) }).criterios).toBe(
      'Los criterios de aceptación no pueden superar los 2000 caracteres',
    )
  })

  it('rechaza Story Points fuera de la escala Fibonacci', () => {
    for (const estimacion of ['0', '4', '22', '-1', '1.5']) {
      expect(validarEdicion({ ...VALIDO, estimacion: estimacion as never }).estimacion).toBe(
        'Los Story Points deben ser uno de 1, 2, 3, 5, 8, 13, 21',
      )
    }
  })

  it('informa todos los campos inválidos a la vez', () => {
    expect(
      validarEdicion({
        titulo: '',
        descripcion: '',
        criterios: 'a'.repeat(2001),
        estimacion: '4' as never,
      }),
    ).toEqual({
      titulo: 'El título es obligatorio',
      descripcion: 'La descripción es obligatoria',
      criterios: 'Los criterios de aceptación no pueden superar los 2000 caracteres',
      estimacion: 'Los Story Points deben ser uno de 1, 2, 3, 5, 8, 13, 21',
    })
  })
})
