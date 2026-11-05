import { render, screen, within } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { GraficoVelocidad } from './GraficoVelocidad'
import type { MetricasDeSprint } from './metricas'

function sprint(numero: number, planificados: number, completados: number): MetricasDeSprint {
  return {
    sprint_id: numero * 10,
    numero,
    metricas: {
      story_points_planificados: planificados,
      story_points_completados: completados,
      velocidad: completados,
      porcentaje_historias_completadas: 0,
      horas_estimadas: null,
      horas_reales: 0,
      desviacion_horas: null,
      desviacion_porcentaje: null,
      defectos_detectados: 0,
      defectos_resueltos: 0,
    },
  }
}

function filasDeLaTabla() {
  const tabla = screen.getByRole('table', { name: 'Valores por Sprint', hidden: true })
  return within(tabla)
    .getAllByRole('row', { hidden: true })
    .slice(1)
    .map((fila) =>
      within(fila)
        .getAllByRole('cell', { hidden: true })
        .concat(within(fila).getAllByRole('rowheader', { hidden: true }))
        .map((celda) => celda.textContent),
    )
}

describe('GraficoVelocidad', () => {
  it('US1.1: Given Sprints finalizados, When veo el gráfico, Then hay un gráfico de barras con la leyenda de las dos series', () => {
    const { container } = render(<GraficoVelocidad sprints={[sprint(1, 8, 5)]} />)

    expect(screen.getByRole('figure', { name: /Story Points comprometidos y completados/ })).toBeTruthy()
    expect(container.querySelector('.recharts-surface')).not.toBeNull()
    const leyenda = container.querySelector<HTMLElement>('.recharts-legend-wrapper')
    expect(leyenda).not.toBeNull()
    expect(within(leyenda!).getByText('Comprometidos')).toBeTruthy()
    expect(within(leyenda!).getByText('Completados')).toBeTruthy()
  })

  it('US1.1: Given las dos series, When leo la leyenda, Then aparecen en el mismo orden que las barras', () => {
    const { container } = render(<GraficoVelocidad sprints={[sprint(1, 8, 5)]} />)

    const rotulos = [...container.querySelectorAll('.recharts-legend-item-text')].map((rotulo) => rotulo.textContent)
    expect(rotulos).toEqual(['Comprometidos', 'Completados'])
  })

  it('US1.2: Given un Sprint con 8 comprometidos y 5 completados, When leo sus valores, Then veo 8 y 5', () => {
    render(<GraficoVelocidad sprints={[sprint(1, 8, 5)]} />)

    expect(filasDeLaTabla()).toEqual([['8', '5', 'Sprint 1']])
  })

  it('US1.3: Given tres Sprints, When leo sus valores, Then aparecen en orden por número', () => {
    render(<GraficoVelocidad sprints={[sprint(1, 8, 5), sprint(2, 3, 3), sprint(3, 5, 1)]} />)

    expect(filasDeLaTabla().map((fila) => fila[2])).toEqual(['Sprint 1', 'Sprint 2', 'Sprint 3'])
  })

  it('US1.4: Given un Sprint con cero Story Points, When leo sus valores, Then sigue presente con sus valores en cero', () => {
    render(<GraficoVelocidad sprints={[sprint(1, 8, 5), sprint(2, 0, 0)]} />)

    expect(filasDeLaTabla()[1]).toEqual(['0', '0', 'Sprint 2'])
  })

  it('US1.5: Given los Sprints que informa HU-03, When los recibo, Then los muestro todos, sin filtrar', () => {
    render(<GraficoVelocidad sprints={[sprint(1, 8, 5), sprint(4, 2, 2)]} />)

    expect(filasDeLaTabla().map((fila) => fila[2])).toEqual(['Sprint 1', 'Sprint 4'])
  })
})
