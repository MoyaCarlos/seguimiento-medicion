import { describe, expect, it } from 'vitest'
import { agruparPorEstado, guardarSprint, type Sprint } from './sprint'

const sprint = (id: number, estado: Sprint['estado']): Sprint => ({ id, proyecto_id: 1, estado })

describe('agruparPorEstado', () => {
  it('devuelve las tres claves vacías si no hay Sprints', () => {
    expect(agruparPorEstado([])).toEqual({ Pendiente: [], Activo: [], Finalizado: [] })
  })

  it('pone cada Sprint solo en la lista de su estado', () => {
    const grupos = agruparPorEstado([sprint(1, 'Finalizado'), sprint(2, 'Activo'), sprint(3, 'Pendiente')])

    expect(grupos.Pendiente.map((s) => s.id)).toEqual([3])
    expect(grupos.Activo.map((s) => s.id)).toEqual([2])
    expect(grupos.Finalizado.map((s) => s.id)).toEqual([1])
  })

  it('ordena cada lista por id ascendente', () => {
    const grupos = agruparPorEstado([sprint(3, 'Finalizado'), sprint(1, 'Finalizado'), sprint(2, 'Finalizado')])

    expect(grupos.Finalizado.map((s) => s.id)).toEqual([1, 2, 3])
  })
})

describe('guardarSprint', () => {
  it('agrega un Sprint nuevo', () => {
    expect(guardarSprint([sprint(1, 'Pendiente')], sprint(2, 'Pendiente'))).toEqual([
      sprint(1, 'Pendiente'),
      sprint(2, 'Pendiente'),
    ])
  })

  it('reemplaza el Sprint de igual id', () => {
    expect(guardarSprint([sprint(1, 'Pendiente')], sprint(1, 'Activo'))).toEqual([sprint(1, 'Activo')])
  })

  it('no muta la lista original', () => {
    const lista = [sprint(1, 'Pendiente')]

    guardarSprint(lista, sprint(1, 'Activo'))
    guardarSprint(lista, sprint(2, 'Pendiente'))

    expect(lista).toEqual([sprint(1, 'Pendiente')])
  })
})
