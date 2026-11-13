import { render, screen, within } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import type { GrupoAvance } from './avance'
import { GrupoHistorias } from './GrupoHistorias'

const grupo: GrupoAvance = {
  estado: 'En progreso',
  cantidad: 2,
  story_points: 5,
  historias: [
    { id: 2, titulo: 'Rendir el examen', estimacion_sp: 5, estado: 'En progreso', sprint: { id: 20, numero: 2 } },
    { id: 3, titulo: 'Emitir el comprobante', estimacion_sp: null, estado: 'En progreso', sprint: { id: 20, numero: 2 } },
  ],
}

const chipsDe = (region: HTMLElement) =>
  Array.from(region.querySelectorAll('.estado-columna__chips span')).map((chip) => chip.textContent)

describe('GrupoHistorias', () => {
  it('US1.1: Given un grupo, When se muestra, Then tiene un encabezado con su estado y cuenta como región', () => {
    render(<GrupoHistorias grupo={grupo} />)

    expect(screen.getByRole('heading', { level: 2, name: 'En progreso' })).toBeTruthy()
    expect(screen.getByRole('region', { name: 'En progreso' })).toBeTruthy()
  })

  it('US1.2: Given un grupo con historias, When se muestra, Then informa su cantidad y sus Story Points tal como vienen', () => {
    render(<GrupoHistorias grupo={{ ...grupo, cantidad: 7, story_points: 99 }} />)

    const region = screen.getByRole('region', { name: 'En progreso' })
    expect(chipsDe(region)).toEqual(['7 historias', '99 SP'])
  })

  it('US1: Given un grupo con historias, When se muestra, Then las lista en una lista semántica en el orden recibido', () => {
    render(<GrupoHistorias grupo={grupo} />)

    const items = within(screen.getByRole('list')).getAllByRole('listitem')
    expect(items.map((item) => item.querySelector('h3')?.textContent)).toEqual([
      'Rendir el examen',
      'Emitir el comprobante',
    ])
  })

  it('US1.3: Given una historia sin Story Points, When se muestra, Then dice "sin estimar" y no 0', () => {
    render(<GrupoHistorias grupo={grupo} />)

    const [conEstimacion, sinEstimar] = screen.getAllByRole('listitem')
    expect(conEstimacion?.textContent).toContain('Story Points: 5')
    expect(sinEstimar?.textContent).toContain('Story Points: sin estimar')
    expect(sinEstimar?.textContent).not.toMatch(/Story Points: 0/)
  })

  it('US2.1: Given una historia en un Sprint, When se muestra, Then dice "Sprint <número>"', () => {
    render(<GrupoHistorias grupo={grupo} />)

    expect(screen.getAllByRole('listitem')[0]?.textContent).toContain('Sprint 2')
  })

  it('US2.2: Given una historia sin Sprint, When se muestra, Then dice "Product Backlog"', () => {
    const nueva: GrupoAvance = {
      estado: 'Nueva',
      cantidad: 1,
      story_points: 3,
      historias: [{ id: 1, titulo: 'Pagar la tasa', estimacion_sp: 3, estado: 'Nueva', sprint: null }],
    }

    render(<GrupoHistorias grupo={nueva} />)

    expect(screen.getByRole('listitem').textContent).toContain('Product Backlog')
  })

  it('US1.4: Given un grupo vacío, When se muestra, Then sigue presente con cantidad 0 y un texto de grupo vacío, sin lista', () => {
    render(<GrupoHistorias grupo={{ estado: 'Completada', cantidad: 0, story_points: 0, historias: [] }} />)

    const region = screen.getByRole('region', { name: 'Completada' })
    expect(chipsDe(region)).toEqual(['0 historias', '0 SP'])
    expect(within(region).getByText('No hay historias en este grupo.')).toBeTruthy()
    expect(within(region).queryByRole('list')).toBeNull()
  })
})
