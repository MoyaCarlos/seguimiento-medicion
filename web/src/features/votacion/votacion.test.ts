import { describe, expect, it } from 'vitest'
import {
  CARTAS,
  esScrumMaster,
  estaRevelada,
  faltanPorVotar,
  puedeRevelar,
  textoDeVoto,
  type Integrante,
  type Votacion,
  type VotoDeIntegrante,
} from './votacion'

const MASTER: Integrante = { id: 1, nombre: 'Beta Uno', rol: 'scrum_master' }
const BUILDER: Integrante = { id: 2, nombre: 'Beta Dos', rol: 'product_builder' }

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

describe('esScrumMaster', () => {
  it('Given un scrum_master, Then es Scrum Master; un product_builder o nadie, no', () => {
    expect(esScrumMaster(MASTER)).toBe(true)
    expect(esScrumMaster(BUILDER)).toBe(false)
    expect(esScrumMaster(undefined)).toBe(false)
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

  it('Given votaron todos y el actor es Scrum Master, When la ronda está abierta, Then puede revelar', () => {
    expect(puedeRevelar(completa, MASTER)).toBe(true)
  })

  it('Given que falta un voto, Then no puede revelar', () => {
    expect(puedeRevelar(votacion([voto(1, true, null), voto(2, false, null)]), MASTER)).toBe(false)
  })

  it('Given un Product Builder o nadie elegido, Then no puede revelar', () => {
    expect(puedeRevelar(completa, BUILDER)).toBe(false)
    expect(puedeRevelar(completa, undefined)).toBe(false)
  })

  it('Given una ronda ya revelada, Then no puede revelar de nuevo', () => {
    const revelada = votacion([voto(1, true, 3)], '2026-11-02T10:00:00Z')

    expect(puedeRevelar(revelada, MASTER)).toBe(false)
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
