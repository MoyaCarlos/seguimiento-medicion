import { describe, expect, it } from 'vitest'
import { armarCuerpo, parsearProyectoId } from './historia'
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
