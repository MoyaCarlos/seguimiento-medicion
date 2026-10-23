import { describe, expect, it } from 'vitest'
import { ApiError, ConnectionError } from '../../compartido/api/client'
import { armarCuerpo, contarCaracteres, mapearError, parsearProyectoId, validarHistoria } from './historia'
import type { FormularioHistoria } from './historia'

const FORMULARIO: FormularioHistoria = {
  titulo: '  Como usuario quiero iniciar sesión  ',
  descripcion: '  Autenticación con email y contraseña  ',
  prioridad: 'M',
  valorNegocio: '',
}

describe('parsearProyectoId', () => {
  it('convierte un identificador válido en número', () => {
    expect(parsearProyectoId('42')).toBe(42)
  })

  it('acepta el máximo entero seguro', () => {
    expect(parsearProyectoId('9007199254740991')).toBe(9007199254740991)
  })

  it.each([
    ['undefined', undefined],
    ['vacío', ''],
    ['solo espacios', '   '],
    ['cero', '0'],
    ['negativo', '-1'],
    ['decimal', '1.5'],
    ['no numérico', 'abc'],
    ['con espacio al inicio', ' 1'],
    ['notación científica', '1e3'],
    ['fuera del rango seguro', '9007199254740993'],
    ['fuera del rango entero', '99999999999999999999'],
  ])('rechaza un identificador %s', (_caso, valor) => {
    expect(parsearProyectoId(valor)).toBeNull()
  })
})

describe('armarCuerpo', () => {
  it('arma el cuerpo con el proyecto numérico y los textos recortados', () => {
    const cuerpo = armarCuerpo(42, FORMULARIO)

    expect(cuerpo).toEqual({
      proyecto_id: 42,
      titulo: 'Como usuario quiero iniciar sesión',
      descripcion: 'Autenticación con email y contraseña',
      prioridad: 'M',
    })
  })

  it.each(['M', 'S', 'C', 'W'] as const)('envía la prioridad %s', (prioridad) => {
    const cuerpo = armarCuerpo(7, { ...FORMULARIO, prioridad })

    expect(cuerpo.prioridad).toBe(prioridad)
  })

  it('envía el valor de negocio como número cuando se informa', () => {
    const cuerpo = armarCuerpo(42, { ...FORMULARIO, valorNegocio: '13' })

    expect(cuerpo.valor_negocio).toBe(13)
  })

  it('omite el valor de negocio cuando la opción es "Sin valor"', () => {
    const cuerpo = armarCuerpo(42, { ...FORMULARIO, valorNegocio: '' })

    expect(cuerpo).not.toHaveProperty('valor_negocio')
  })
})

const VALIDO: FormularioHistoria = {
  titulo: 'Historia de usuario',
  descripcion: 'Descripción de la historia',
  prioridad: 'M',
  valorNegocio: '',
}

describe('contarCaracteres', () => {
  it('cuenta un emoji como un punto de código', () => {
    expect(contarCaracteres('😀')).toBe(1)
  })

  it('cuenta una letra acentuada como un punto de código', () => {
    expect(contarCaracteres('ñ')).toBe(1)
  })

  it('cuenta una letra con acento combinado como dos puntos de código', () => {
    expect(contarCaracteres('e\u0301')).toBe(2)
  })

  it('devuelve cero para un texto vacío', () => {
    expect(contarCaracteres('')).toBe(0)
  })
})

describe('validarHistoria', () => {
  it('acepta una historia válida', () => {
    expect(validarHistoria(VALIDO)).toEqual({})
  })

  it.each([
    ['vacío', ''],
    ['con solo espacios', '   '],
  ])('exige el título cuando está %s', (_caso, titulo) => {
    expect(validarHistoria({ ...VALIDO, titulo })).toEqual({ titulo: expect.any(String) })
  })

  it.each([
    ['vacía', ''],
    ['con solo espacios', '   '],
  ])('exige la descripción cuando está %s', (_caso, descripcion) => {
    expect(validarHistoria({ ...VALIDO, descripcion })).toEqual({
      descripcion: expect.any(String),
    })
  })

  it('rechaza un título de 201 puntos de código', () => {
    expect(validarHistoria({ ...VALIDO, titulo: 'a'.repeat(201) })).toEqual({
      titulo: expect.any(String),
    })
  })

  it('rechaza una descripción de 2001 puntos de código', () => {
    expect(validarHistoria({ ...VALIDO, descripcion: 'a'.repeat(2001) })).toEqual({
      descripcion: expect.any(String),
    })
  })

  it('acepta un título de exactamente 200 puntos de código', () => {
    expect(validarHistoria({ ...VALIDO, titulo: 'a'.repeat(200) })).toEqual({})
  })

  it('acepta una descripción de exactamente 2000 puntos de código', () => {
    expect(validarHistoria({ ...VALIDO, descripcion: 'a'.repeat(2000) })).toEqual({})
  })

  it('evalúa la longitud después de recortar los extremos', () => {
    expect(validarHistoria({ ...VALIDO, titulo: `  ${'a'.repeat(200)}  ` })).toEqual({})
  })

  it('cuenta puntos de código y no unidades UTF-16', () => {
    expect(validarHistoria({ ...VALIDO, titulo: '😀'.repeat(201) })).toEqual({
      titulo: expect.any(String),
    })
    expect(validarHistoria({ ...VALIDO, titulo: '😀'.repeat(200) })).toEqual({})
  })

  it('devuelve todos los errores de campo del intento', () => {
    expect(
      validarHistoria({ ...VALIDO, titulo: '', descripcion: 'a'.repeat(2001) }),
    ).toEqual({
      titulo: expect.any(String),
      descripcion: expect.any(String),
    })
  })
})

describe('mapearError', () => {
  it.each([
    ['titulo', 'el título es obligatorio'],
    ['descripcion', 'la descripción es obligatoria'],
    ['prioridad', 'la prioridad debe ser M, S, C o W'],
    ['valor_negocio', 'el valor de negocio debe ser uno de 1, 2, 3, 5, 8, 13, 21'],
  ] as const)('ubica el mensaje del campo %s junto a ese campo', (campo, mensaje) => {
    expect(mapearError(new ApiError(400, mensaje, campo))).toEqual({ campo, mensaje })
  })

  it('manda al aviso general un error sin campo', () => {
    const resultado = mapearError(new ApiError(400, 'datos inválidos'))

    expect(resultado.campo).toBeUndefined()
    expect(resultado.mensaje).toBe('datos inválidos')
  })

  it.each(['otro', 'proyecto_id'])(
    'manda al aviso general el error del campo %s',
    (campo) => {
      const resultado = mapearError(new ApiError(400, 'mensaje del servidor', campo))

      expect(resultado.campo).toBeUndefined()
      expect(resultado.mensaje).toBe('mensaje del servidor')
    },
  )

  it('manda al aviso general el proyecto no encontrado', () => {
    const resultado = mapearError(new ApiError(404, 'proyecto no encontrado'))

    expect(resultado.campo).toBeUndefined()
    expect(resultado.mensaje).toBe('proyecto no encontrado')
  })

  it('manda al aviso general un error de conexión', () => {
    const resultado = mapearError(new ConnectionError())

    expect(resultado.campo).toBeUndefined()
    expect(resultado.mensaje).toBe('No se pudo conectar con el servidor')
  })
})
