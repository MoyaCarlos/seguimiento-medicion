import type { ProjectFormData } from './projectForm'

export type Project = {
  id: number
  nombre: string
  descripcion: string
  fecha_inicio: string | null
  fecha_fin: string | null
  creado_en: string
}

export async function createProject(data: ProjectFormData): Promise<Project> {
  const response = await fetch('/projects', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      nombre: data.name.trim(),
      descripcion: data.description.trim(),
      creador: data.creator.trim(),
      fecha_inicio: data.startDate || null,
      fecha_fin: data.endDate || null,
    }),
  })

  const payload = (await response.json()) as Project | { mensaje?: string }
  if (!response.ok) {
    throw new Error('mensaje' in payload && payload.mensaje ? payload.mensaje : 'No se pudo crear el proyecto')
  }

  return payload as Project
}
