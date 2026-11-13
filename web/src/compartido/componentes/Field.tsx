import type { InputHTMLAttributes, ReactNode } from 'react'
import { ErrorMessage } from './ErrorMessage'

type FieldProps = InputHTMLAttributes<HTMLInputElement> & {
  label: string
  error?: string
  help?: string
  ayuda?: string
  accion?: ReactNode
}

export function Field({ label, error, help, ayuda, accion, id, ...props }: FieldProps) {
  const fieldId = id ?? props.name
  const errorId = error && fieldId ? `${fieldId}-error` : undefined
  const helpId = help && fieldId ? `${fieldId}-help` : undefined
  const ayudaId = !error && ayuda && fieldId ? `${fieldId}-ayuda` : undefined
  const describedBy = [helpId, errorId ?? ayudaId].filter(Boolean).join(' ') || undefined

  const input = (
    <input
      id={fieldId}
      aria-invalid={Boolean(error)}
      aria-describedby={describedBy}
      {...props}
    />
  )

  return (
    <div className="field">
      <label htmlFor={fieldId}>{label}</label>
      {accion ? <div className="field__control">{input}{accion}</div> : input}
      {help && <p id={helpId} className="message">{help}</p>}
      {ayudaId && <p id={ayudaId} className="field__ayuda">{ayuda}</p>}
      {error && <ErrorMessage id={errorId} message={error} />}
    </div>
  )
}
