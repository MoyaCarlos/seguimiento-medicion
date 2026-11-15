import type { ChangeEvent } from 'react'

type Props = {
  id: string
  name: string
  label: string
  value: string
  onChange: (valor: string) => void
  error?: string
  disabled?: boolean
  rows?: number
}

// Stub del RED de #215; el GREEN agrega la barra, la ayuda y la vista previa.
export function CampoMarkdown({
  id,
  name,
  label,
  value,
  onChange,
  error,
  disabled = false,
  rows = 4,
}: Props) {
  return (
    <div className="field campo-markdown">
      <label htmlFor={id}>{label}</label>
      <textarea
        id={id}
        name={name}
        value={value}
        onChange={(evento: ChangeEvent<HTMLTextAreaElement>) => onChange(evento.target.value)}
        disabled={disabled}
        rows={rows}
        aria-invalid={Boolean(error)}
        aria-describedby={error ? `${id}-error` : undefined}
      />
      {error && (
        <p id={`${id}-error`} className="message message--error" role="alert">
          {error}
        </p>
      )}
    </div>
  )
}
