import { describe, expect, it } from 'vitest'
import {
  CARTAS,
  estaRevelada,
  faltanPorVotar,
  puedeRevelar,
  textoDeDiferencias,
  textoDeVoto,
  type Votacion,
  type VotoDeIntegrante,
} from './votacion'

const voto = (usuario_id: number, ya_voto: boolean, valor: number | null): VotoDeIntegrante => ({
  usuario_id,
  nombre: `Integrante ${usuario_id}`,
  rol: 'product_builder',
  ya_voto,
  valor,
})

const votacion = (integrantes: VotoDeIntegrante[], revelada_en: string | null = null): Votacion => ({
  historia_id: 3,
  ronda: 1,
  revelada_en,
  integrantes,
  promedio: null,
  diferencias: null,
})

describe('CARTAS', () => {
  it('Given la baraja del proyecto, When se consulta, Then son exactamente 1, 2, 3, 5, 8, 13 y 21', () => {
    expect(CARTAS).toEqual([1, 2, 3, 5, 8, 13, 21])
  })
})

describe('estaRevelada', () => {
  it('Given revelada_en con fecha, Then está revelada; con null, no', () => {
    expect(estaRevelada(votacion([], '2026-11-02T10:00:00Z'))).toBe(true)
    expect(estaRevelada(votacion([]))).toBe(false)
  })
})

describe('faltanPorVotar', () => {
  it('Given dos votaron y uno no, Then lista solo al que falta', () => {
    const v = votacion([voto(1, true, null), voto(2, false, null), voto(3, true, 5)])

    expect(faltanPorVotar(v).map((i) => i.usuario_id)).toEqual([2])
  })
})

describe('puedeRevelar', () => {
  const completa = votacion([voto(1, true, null), voto(2, true, 5)])

  it('Given votaron todos, When la ronda está abierta, Then puede revelar', () => {
    expect(puedeRevelar(completa)).toBe(true)
  })

  it('Given que falta un voto, Then no puede revelar', () => {
    expect(puedeRevelar(votacion([voto(1, true, null), voto(2, false, null)]))).toBe(false)
  })

  it('Given una ronda ya revelada, Then no puede revelar de nuevo', () => {
    const revelada = votacion([voto(1, true, 3)], '2026-11-02T10:00:00Z')

    expect(puedeRevelar(revelada)).toBe(false)
  })
})

describe('textoDeVoto', () => {
  it('Given quien no votó, Then "Falta"', () => {
    expect(textoDeVoto(voto(1, false, null))).toBe('Falta')
  })

  it('Given quien votó y el valor viene oculto, Then "Oculto"', () => {
    expect(textoDeVoto(voto(1, true, null))).toBe('Oculto')
  })

  it('Given quien votó y el valor está disponible (propio o revelado), Then el valor', () => {
    expect(textoDeVoto(voto(1, true, 8))).toBe('8')
  })
})

describe('textoDeDiferencias', () => {
  it('Given votos 3, 5 y 8 revelados, Then indica que hay diferencias con el mínimo y el máximo', () => {
    expect(textoDeDiferencias({ hay_diferencias: true, minimo: 3, maximo: 8 })).toBe(
      'Hay diferencias: mínimo 3, máximo 8',
    )
  })

  it('Given votos que coinciden, Then indica que no hay diferencias', () => {
    expect(textoDeDiferencias({ hay_diferencias: false, minimo: null, maximo: null })).toBe(
      'No hay diferencias: todos votaron lo mismo',
    )
  })

  it('Given una ronda sin revelar (null), Then no dice nada', () => {
    expect(textoDeDiferencias(null)).toBe('')
  })
})
