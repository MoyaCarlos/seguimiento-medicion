import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'

type CampoDescripcionProps = {
  value: string
  onChange: (value: string) => void
  error?: string
}

export function CampoDescripcion({ value, onChange, error }: CampoDescripcionProps) {
  return (
    <div className="field">
      <label htmlFor="descripcion">
        Descripción <span className="hu04-opcional">(opcional)</span>
      </label>
      <textarea
        id="descripcion"
        rows={3}
        value={value}
        onChange={(event) => onChange(event.target.value)}
        aria-invalid={Boolean(error)}
        aria-describedby={error ? 'descripcion-error' : undefined}
      />
      {error && <ErrorMessage id="descripcion-error" message={error} />}
    </div>
  )
}

type CampoSelectProps = {
  id: string
  label: string
  value: string
  opciones: { valor: string; etiqueta: string }[]
  onChange: (value: string) => void
  error?: string
}

export function CampoSelect({ id, label, value, opciones, onChange, error }: CampoSelectProps) {
  return (
    <div className="field">
      <label htmlFor={id}>{label}</label>
      <select
        id={id}
        value={value}
        onChange={(event) => onChange(event.target.value)}
        aria-invalid={Boolean(error)}
        aria-describedby={error ? `${id}-error` : undefined}
      >
        <option value="">Seleccioná un rol</option>
        {opciones.map((opcion) => (
          <option key={opcion.valor} value={opcion.valor}>
            {opcion.etiqueta}
          </option>
        ))}
      </select>
      {error && <ErrorMessage id={`${id}-error`} message={error} />}
    </div>
  )
}
