import { describe, expect, it } from 'vitest'
import { validarCreacion, validarEdicion, validarIntegrante } from './validacion'

describe('validarCreacion', () => {
  it('exige el nombre del proyecto y el del creador', () => {
    expect(
      validarCreacion({
        nombre: '   ',
        descripcion: '',
        fecha_inicio: '',
        fecha_fin: '',
        creador: '  ',
      }),
    ).toEqual({
      nombre: 'El nombre es obligatorio',
      creador: 'El nombre del creador es obligatorio',
    })
  })

  it('acepta fechas opcionales', () => {
    expect(
      validarCreacion({
        nombre: 'Proyecto Alfa',
        descripcion: '',
        fecha_inicio: '',
        fecha_fin: '',
        creador: 'Ana',
      }),
    ).toEqual({})
  })

  it('rechaza una fecha de fin que no es posterior a la de inicio', () => {
    expect(
      validarCreacion({
        nombre: 'Proyecto Alfa',
        descripcion: 'Descripción',
        fecha_inicio: '2026-10-15',
        fecha_fin: '2026-10-15',
        creador: 'Ana',
      }),
    ).toEqual({ fecha_fin: 'La fecha de fin debe ser posterior a la fecha de inicio' })
  })

  it('limita el nombre a 100 caracteres y la descripción a 2000', () => {
    expect(
      validarCreacion({
        nombre: 'a'.repeat(101),
        descripcion: 'b'.repeat(2001),
        fecha_inicio: '',
        fecha_fin: '',
        creador: 'Ana',
      }),
    ).toEqual({
      nombre: 'El nombre no puede superar los 100 caracteres',
      descripcion: 'La descripción no puede superar los 2000 caracteres',
    })
  })

  it('limita el nombre del creador a 200 caracteres contados por puntos de código', () => {
    const datosBase = {
      nombre: 'Proyecto Alfa',
      descripcion: '',
      fecha_inicio: '',
      fecha_fin: '',
    }

    expect(validarCreacion({ ...datosBase, creador: '😀'.repeat(200) })).toEqual({})
    expect(validarCreacion({ ...datosBase, creador: '😀'.repeat(201) })).toEqual({
      creador: 'El nombre del creador no puede superar los 200 caracteres',
    })
  })
})

describe('validarEdicion', () => {
  it('no exige un creador porque no se edita', () => {
    expect(
      validarEdicion({
        nombre: 'Proyecto Alfa',
        descripcion: '',
        fecha_inicio: '2026-10-15',
        fecha_fin: '2026-11-15',
        creador: '',
      }),
    ).toEqual({})
  })

  it('rechaza una fecha de fin anterior a la de inicio', () => {
    expect(
      validarEdicion({
        nombre: 'Proyecto Alfa',
        descripcion: '',
        fecha_inicio: '2026-10-15',
        fecha_fin: '2026-10-14',
        creador: '',
      }),
    ).toEqual({ fecha_fin: 'La fecha de fin debe ser posterior a la fecha de inicio' })
  })
})

describe('validarIntegrante', () => {
  it('exige nombre y un rol válido', () => {
    expect(validarIntegrante({ nombre: '', rol: '' })).toEqual({
      nombre: 'El nombre del integrante es obligatorio',
      rol: 'Seleccioná un rol válido',
    })
  })

  it('acepta los roles de la historia', () => {
    expect(validarIntegrante({ nombre: 'Jimena', rol: 'product_builder' })).toEqual({})
    expect(validarIntegrante({ nombre: 'Ana', rol: 'scrum_master' })).toEqual({})
  })
})
