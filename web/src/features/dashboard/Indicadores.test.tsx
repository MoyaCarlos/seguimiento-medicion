import { render, screen, within } from '@testing-library/react'
import { MemoryRouter } from 'react-router'
import { describe, expect, it } from 'vitest'
import { Indicadores } from './Indicadores'
import type { Metricas } from './metricas'

const base: Metricas = {
  story_points_planificados: 8,
  story_points_completados: 5,
  velocidad: 5.5,
  porcentaje_historias_completadas: 50,
  horas_estimadas: 32,
  horas_reales: 30,
  desviacion_horas: -2,
  desviacion_porcentaje: -6.25,
  defectos_detectados: 2,
  defectos_resueltos: 1,
}

function mostrar(cambios: Partial<Metricas> = {}, horasPorStoryPoint: number | null = 4) {
  render(
    <MemoryRouter>
      <Indicadores metricas={{ ...base, ...cambios }} horasPorStoryPoint={horasPorStoryPoint} proyectoId="7" />
    </MemoryRouter>,
  )
}

const indicador = (nombre: string) => screen.getByRole('region', { name: nombre })
const barra = (nombre: string) => screen.getByRole('meter', { name: nombre })

describe('Indicadores', () => {
  it('US2.1: Given valores del proyecto, When veo el panel, Then hay un indicador por cada métrica', () => {
    mostrar()

    expect(within(indicador('Velocidad')).getByText('5,5')).toBeTruthy()
    expect(within(indicador('Historias completadas')).getByText('50%')).toBeTruthy()
    expect(within(indicador('Horas')).getByText('Horas estimadas')).toBeTruthy()
    expect(within(indicador('Horas')).getByText('Horas reales')).toBeTruthy()
    expect(indicador('Desviación')).toBeTruthy()
    expect(within(indicador('Defectos')).getByText('Defectos detectados')).toBeTruthy()
    expect(within(indicador('Defectos')).getByText('Defectos resueltos')).toBeTruthy()
  })

  it('US2.2: Given 50% de historias completadas, When veo su indicador, Then se lee como proporción de un total', () => {
    mostrar()

    const proporcion = barra('Historias completadas')
    expect(proporcion.getAttribute('aria-valuenow')).toBe('50')
    expect(proporcion.getAttribute('aria-valuemax')).toBe('100')
    expect(proporcion.querySelector<HTMLElement>('span')?.style.width).toBe('50%')
  })

  it('US2.3: Given 32 horas estimadas y 30 reales, When las veo, Then se comparan en una misma escala', () => {
    mostrar()

    const estimadas = barra('Horas estimadas')
    const reales = barra('Horas reales')
    expect(estimadas.getAttribute('aria-valuemax')).toBe('32')
    expect(reales.getAttribute('aria-valuemax')).toBe('32')
    expect(estimadas.querySelector<HTMLElement>('span')?.style.width).toBe('100%')
    expect(reales.querySelector<HTMLElement>('span')?.style.width).toBe('93.75%')
    expect(within(indicador('Horas')).getByText('32 h')).toBeTruthy()
    expect(within(indicador('Horas')).getByText('30 h')).toBeTruthy()
  })

  it('US2.4: Given una desviación de -2 horas, When veo su indicador, Then veo ambos valores y que se trabajó menos', () => {
    mostrar()

    const desviacion = within(indicador('Desviación'))
    expect(desviacion.getByText('-2 h')).toBeTruthy()
    expect(desviacion.getByText('-6,25%')).toBeTruthy()
    expect(desviacion.getByText('Se trabajó menos de lo estimado')).toBeTruthy()
  })

  it('US2.5: Given 2 defectos detectados y 1 resuelto, When los veo, Then se comparan en una misma escala', () => {
    mostrar()

    expect(barra('Defectos detectados').getAttribute('aria-valuemax')).toBe('2')
    expect(barra('Defectos resueltos').getAttribute('aria-valuemax')).toBe('2')
    expect(barra('Defectos resueltos').querySelector<HTMLElement>('span')?.style.width).toBe('50%')
  })

  it('US2.6: Given un proyecto sin horas por Story Point, When veo los indicadores, Then lo que falta dice sin datos y se indica dónde definirlo', () => {
    mostrar({ horas_estimadas: null, desviacion_horas: null, desviacion_porcentaje: null }, null)

    const horas = within(indicador('Horas'))
    expect(horas.getByText('sin datos')).toBeTruthy()
    expect(screen.queryByRole('meter', { name: 'Horas estimadas' })).toBeNull()
    expect(horas.getByText('30 h')).toBeTruthy()
    expect(screen.queryByRole('meter', { name: 'Horas reales' })).toBeNull()
    expect(within(indicador('Desviación')).getAllByText('sin datos')).toHaveLength(2)
    expect(within(indicador('Desviación')).queryByText(/Se trabajó|Sin desviación/)).toBeNull()
    expect(screen.getByRole('link', { name: /editar el proyecto/i }).getAttribute('href')).toBe('/projects/7/proyecto/edit')
  })

  it('Given un proyecto con horas por Story Point, When veo los indicadores, Then no hay aviso', () => {
    mostrar()

    expect(screen.queryByRole('link', { name: /editar el proyecto/i })).toBeNull()
  })

  it('Caso límite: Given horas con valor y porcentaje sin datos, When veo la desviación, Then muestra las horas y sin datos', () => {
    mostrar({ horas_estimadas: 0, horas_reales: 3, desviacion_horas: 3, desviacion_porcentaje: null })

    const desviacion = within(indicador('Desviación'))
    expect(desviacion.getByText('3 h')).toBeTruthy()
    expect(desviacion.getByText('sin datos')).toBeTruthy()
    expect(desviacion.getByText('Se trabajó más de lo estimado')).toBeTruthy()
  })

  it('Caso límite: Given horas y defectos en cero, When los veo, Then muestran 0, sin barras rotas', () => {
    mostrar({
      horas_estimadas: 0,
      horas_reales: 0,
      desviacion_horas: 0,
      desviacion_porcentaje: null,
      defectos_detectados: 0,
      defectos_resueltos: 0,
    })

    expect(within(indicador('Horas')).getAllByText('0 h').length).toBeGreaterThan(0)
    expect(within(indicador('Defectos')).getAllByText('0')).toHaveLength(2)
    expect(barra('Horas reales').querySelector<HTMLElement>('span')?.style.width).toBe('0%')
    expect(barra('Defectos resueltos').querySelector<HTMLElement>('span')?.style.width).toBe('0%')
    expect(barra('Defectos resueltos').getAttribute('aria-valuemax')).toBe('1')
    expect(barra('Horas reales').getAttribute('aria-valuemax')).toBe('1')
    expect(within(indicador('Desviación')).getByText('Sin desviación: se trabajó lo estimado')).toBeTruthy()
  })
})
