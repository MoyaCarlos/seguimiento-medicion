import { useState } from 'react'
import type { FormEvent } from 'react'
import { createProject } from './api'
import type { Project } from './api'
import { validateProjectForm } from './projectForm'
import type { ProjectFormData, ProjectFormErrors } from './projectForm'
import './App.css'

const initialForm: ProjectFormData = {
  name: '',
  description: '',
  startDate: '',
  endDate: '',
  creator: '',
}

function App() {
  const [form, setForm] = useState(initialForm)
  const [errors, setErrors] = useState<ProjectFormErrors>({})
  const [project, setProject] = useState<Project | null>(null)
  const [submitError, setSubmitError] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  function updateField(field: keyof ProjectFormData, value: string) {
    setForm((current) => ({ ...current, [field]: value }))
    setErrors((current) => ({ ...current, [field]: undefined }))
    setSubmitError('')
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const validationErrors = validateProjectForm(form)
    setErrors(validationErrors)
    setSubmitError('')
    if (Object.keys(validationErrors).length > 0) return

    setIsSubmitting(true)
    try {
      setProject(await createProject(form))
    } catch (error) {
      setSubmitError(error instanceof Error ? error.message : 'No se pudo crear el proyecto')
    } finally {
      setIsSubmitting(false)
    }
  }

  if (project) {
    return (
      <main className="project-shell">
        <header className="app-header">
          <span className="brand">Seguimiento</span>
          <span className="header-label">Proyecto creado</span>
        </header>
        <section className="project-panel" aria-labelledby="project-title">
          <p className="eyebrow">Panel principal</p>
          <h1 id="project-title">{project.nombre}</h1>
          <p className="project-description">{project.descripcion || 'Sin descripción'}</p>
          <div className="project-meta">
            <div><span>Inicio</span><strong>{project.fecha_inicio ?? 'Sin definir'}</strong></div>
            <div><span>Fin</span><strong>{project.fecha_fin ?? 'Sin definir'}</strong></div>
          </div>
          <button className="secondary-button" type="button" onClick={() => setProject(null)}>
            Crear otro proyecto
          </button>
        </section>
      </main>
    )
  }

  return (
    <main className="project-shell">
      <header className="app-header">
        <span className="brand">Seguimiento</span>
        <span className="header-label">Nuevo proyecto</span>
      </header>
      <section className="form-card" aria-labelledby="form-title">
        <p className="eyebrow">Sprint 2 · HU-04</p>
        <h1 id="form-title">Crear proyecto</h1>
        <p className="intro">Definí los datos iniciales del proyecto para comenzar a hacerle seguimiento.</p>
        <form onSubmit={handleSubmit} noValidate>
          <div className="field">
            <label htmlFor="name">Nombre del proyecto</label>
            <input id="name" value={form.name} onChange={(event) => updateField('name', event.target.value)} aria-invalid={Boolean(errors.name)} aria-describedby={errors.name ? 'name-error' : undefined} />
            {errors.name && <span className="field-error" id="name-error">{errors.name}</span>}
          </div>
          <div className="field">
            <label htmlFor="description">Descripción <span>(opcional)</span></label>
            <textarea id="description" rows={3} value={form.description} onChange={(event) => updateField('description', event.target.value)} />
          </div>
          <div className="field-row">
            <div className="field">
              <label htmlFor="startDate">Fecha de inicio</label>
              <input id="startDate" type="date" value={form.startDate} onChange={(event) => updateField('startDate', event.target.value)} aria-invalid={Boolean(errors.startDate)} aria-describedby={errors.startDate ? 'startDate-error' : undefined} />
              {errors.startDate && <span className="field-error" id="startDate-error">{errors.startDate}</span>}
            </div>
            <div className="field">
              <label htmlFor="endDate">Fecha de fin</label>
              <input id="endDate" type="date" value={form.endDate} onChange={(event) => updateField('endDate', event.target.value)} aria-invalid={Boolean(errors.endDate)} aria-describedby={errors.endDate ? 'endDate-error' : undefined} />
              {errors.endDate && <span className="field-error" id="endDate-error">{errors.endDate}</span>}
            </div>
          </div>
          <div className="field">
            <label htmlFor="creator">Tu nombre</label>
            <input id="creator" value={form.creator} onChange={(event) => updateField('creator', event.target.value)} aria-invalid={Boolean(errors.creator)} aria-describedby={errors.creator ? 'creator-error' : undefined} />
            {errors.creator && <span className="field-error" id="creator-error">{errors.creator}</span>}
          </div>
          {submitError && <p className="submit-error" role="alert">{submitError}</p>}
          <button className="primary-button" type="submit" disabled={isSubmitting}>
            {isSubmitting ? 'Creando…' : 'Crear proyecto'}
          </button>
        </form>
      </section>
    </main>
  )
}

export default App
