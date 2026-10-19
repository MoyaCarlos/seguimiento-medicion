export type ProjectFormData = {
  name: string
  description: string
  startDate: string
  endDate: string
  creator: string
}

export type ProjectFormErrors = Partial<Record<keyof ProjectFormData, string>>

export function validateProjectForm(data: ProjectFormData): ProjectFormErrors {
  const errors: ProjectFormErrors = {}

  if (!data.name.trim()) {
    errors.name = 'El nombre es obligatorio'
  }
  if (!data.creator.trim()) {
    errors.creator = 'El nombre del creador es obligatorio'
  }
  if (!data.startDate) {
    errors.startDate = 'La fecha de inicio es obligatoria'
  }
  if (!data.endDate) {
    errors.endDate = 'La fecha de fin es obligatoria'
  }
  if (data.startDate && data.endDate && data.endDate <= data.startDate) {
    errors.endDate = 'La fecha de fin debe ser posterior a la fecha de inicio'
  }

  return errors
}
