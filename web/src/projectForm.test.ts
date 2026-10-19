import { describe, expect, it } from 'vitest'
import { validateProjectForm } from './projectForm'

describe('validateProjectForm', () => {
  it('reports required fields before submitting a project', () => {
    expect(
      validateProjectForm({
        name: '',
        description: '',
        startDate: '',
        endDate: '',
        creator: '',
      }),
    ).toEqual({
      name: 'El nombre es obligatorio',
      creator: 'El nombre del creador es obligatorio',
      startDate: 'La fecha de inicio es obligatoria',
      endDate: 'La fecha de fin es obligatoria',
    })
  })

  it('rejects an end date that is not after the start date', () => {
    expect(
      validateProjectForm({
        name: 'Proyecto Alfa',
        description: 'Descripción',
        startDate: '2026-10-15',
        endDate: '2026-10-15',
        creator: 'Ana',
      }),
    ).toEqual({ endDate: 'La fecha de fin debe ser posterior a la fecha de inicio' })
  })

  it('accepts a complete project form', () => {
    expect(
      validateProjectForm({
        name: 'Proyecto Alfa',
        description: 'Descripción',
        startDate: '2026-10-15',
        endDate: '2026-11-15',
        creator: 'Ana',
      }),
    ).toEqual({})
  })
})
