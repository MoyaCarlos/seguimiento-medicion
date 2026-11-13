import { describe, expect, it } from 'vitest'
import {
  datosDelGrafico,
  ejeDelGrafico,
  textoDeDesviacion,
  formatoNumero,
  proporcion,
  textoDesviacion,
  type Metricas,
  type MetricasDeSprint,
} from './metricas'

const metricas: Metricas = {
  story_points_planificados: 8,
  story_points_completados: 5,
  velocidad: 5,
  porcentaje_historias_completadas: 50,
  horas_estimadas: 32,
  horas_reales: 30,
  desviacion_horas: -2,
  desviacion_porcentaje: -6.25,
  defectos_detectados: 2,
  defectos_resueltos: 1,
}

function sprint(numero: number, planificados: number, completados: number): MetricasDeSprint {
  return {
    sprint_id: numero * 10,
    numero,
    metricas: {
      ...metricas,
      story_points_planificados: planificados,
      story_points_completados: completados,
    },
  }
}

describe('datosDelGrafico', () => {
  it('arma una fila por Sprint con comprometidos y completados', () => {
    expect(datosDelGrafico([sprint(1, 8, 5)])).toEqual([
      { sprint: 'Sprint 1', comprometidos: 8, completados: 5 },
    ])
  })

  it('conserva el orden recibido', () => {
    const filas = datosDelGrafico([sprint(1, 8, 5), sprint(2, 3, 3), sprint(3, 5, 1)])

    expect(filas.map((fila) => fila.sprint)).toEqual(['Sprint 1', 'Sprint 2', 'Sprint 3'])
  })

  it('incluye los Sprints con valores en cero', () => {
    expect(datosDelGrafico([sprint(4, 0, 0)])).toEqual([
      { sprint: 'Sprint 4', comprometidos: 0, completados: 0 },
    ])
  })

  it('sin Sprints devuelve una lista vacía', () => {
    expect(datosDelGrafico([])).toEqual([])
  })
})

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

describe('textoDesviacion', () => {
  it.each([
    [2, 'Se trabajó más de lo estimado'],
    [-2, 'Se trabajó menos de lo estimado'],
    [0, 'Sin desviación: se trabajó lo estimado'],
  ])('una desviación de %d horas dice "%s"', (horas, esperado) => {
    expect(textoDesviacion(horas)).toBe(esperado)
  })

  it('sin datos no da ninguna aclaración', () => {
    expect(textoDesviacion(null)).toBeNull()
  })
})

describe('proporcion', () => {
  it.each([
    [5, 10, 50],
    [10, 10, 100],
    [0, 10, 0],
  ])('%d sobre %d es %d%%', (valor, maximo, esperado) => {
    expect(proporcion(valor, maximo)).toBe(esperado)
  })

  it('con máximo cero es 0, sin proporciones indefinidas', () => {
    expect(proporcion(0, 0)).toBe(0)
  })

  it('sin datos es 0', () => {
    expect(proporcion(null, 10)).toBe(0)
  })

  it('nunca pasa de 100', () => {
    expect(proporcion(15, 10)).toBe(100)
  })
})

describe('ejeDelGrafico', () => {
  it.each([
    [16, [0, 8, 16]],
    [3, [0, 2, 4]],
    [40, [0, 20, 40]],
    [13, [0, 7, 14]],
    [0, [0, 1, 2]],
  ])('con un máximo de %d las marcas son %j', (maximo, esperado) => {
    expect(ejeDelGrafico(maximo)).toEqual(esperado)
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
