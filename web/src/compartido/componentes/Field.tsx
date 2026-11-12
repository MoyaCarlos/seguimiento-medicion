import type { InputHTMLAttributes, ReactNode } from 'react'
import { ErrorMessage } from './ErrorMessage'

type FieldProps = InputHTMLAttributes<HTMLInputElement> & {
  label: string
  error?: string
  ayuda?: string
  accion?: ReactNode
}

export function Field({ label, error, ayuda, accion, id, ...props }: FieldProps) {
  const fieldId = id ?? props.name
  const errorId = error && fieldId ? `${fieldId}-error` : undefined
  const ayudaId = !error && ayuda && fieldId ? `${fieldId}-ayuda` : undefined

  const input = (
    <input
      id={fieldId}
      aria-invalid={Boolean(error)}
      aria-describedby={errorId ?? ayudaId}
      {...props}
    />
  )

  return (
    <div className="field">
      <label htmlFor={fieldId}>{label}</label>
      {accion ? <div className="field__control">{input}{accion}</div> : input}
      {ayudaId && <p id={ayudaId} className="field__ayuda">{ayuda}</p>}
      {error && <ErrorMessage id={errorId} message={error} />}
    </div>
  )
}
