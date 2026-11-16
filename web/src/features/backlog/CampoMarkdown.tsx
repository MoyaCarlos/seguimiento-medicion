import { useRef, useState } from 'react'
import type { KeyboardEvent } from 'react'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { Markdown } from './Markdown'

const AYUDA =
  'Admite Markdown: negrita, listas, enlaces e imágenes por URL con ![texto](https://...). No se pueden pegar imágenes.'

type Props = {
  id: string
  name: string
  label: string
  value: string
  onChange: (valor: string) => void
  error?: string
  disabled?: boolean
  rows?: number
  opcional?: boolean
}

type Vista = 'escribir' | 'preview'

export function CampoMarkdown({
  id,
  name,
  label,
  value,
  onChange,
  error,
  disabled = false,
  rows = 4,
  opcional = false,
}: Props) {
  const ref = useRef<HTMLTextAreaElement>(null)
  const [vista, setVista] = useState<Vista>('escribir')
  const ayudaId = `${id}-ayuda`
  const errorId = `${id}-error`
  const opcionalId = `${id}-opc`
  const pestanaEscribir = useRef<HTMLButtonElement>(null)
  const pestanaVista = useRef<HTMLButtonElement>(null)
  const descripcion = [opcional ? opcionalId : '', error ? errorId : ayudaId].filter(Boolean).join(' ')

  function elegir(nueva: Vista, enfocarPestana: boolean) {
    setVista(nueva)
    if (enfocarPestana) (nueva === 'escribir' ? pestanaEscribir : pestanaVista).current?.focus()
  }

  function alPulsarTecla(evento: KeyboardEvent<HTMLDivElement>) {
    const { key } = evento
    if (!['ArrowLeft', 'ArrowRight', 'Home', 'End'].includes(key)) return
    evento.preventDefault()
    const nueva: Vista = key === 'Home' ? 'escribir' : key === 'End' ? 'preview' : vista === 'escribir' ? 'preview' : 'escribir'
    elegir(nueva, true)
  }

  function enfocar(inicio: number, fin: number) {
    requestAnimationFrame(() => {
      const textarea = ref.current
      if (!textarea) return
      textarea.focus()
      textarea.setSelectionRange(inicio, fin)
    })
  }

  function envolver(antes: string, despues: string, porDefecto: string) {
    const textarea = ref.current
    if (!textarea) return
    const inicio = textarea.selectionStart
    const fin = textarea.selectionEnd
    const seleccion = value.slice(inicio, fin)
    const contenido = seleccion || porDefecto
    const nuevo = value.slice(0, inicio) + antes + contenido + despues + value.slice(fin)
    onChange(nuevo)
    enfocar(inicio + antes.length, inicio + antes.length + contenido.length)
  }

  function prefijar(prefijo: string) {
    const textarea = ref.current
    if (!textarea) return
    const inicio = textarea.selectionStart
    const fin = textarea.selectionEnd
    const inicioLinea = value.lastIndexOf('\n', inicio - 1) + 1
    const salto = value.indexOf('\n', fin)
    const finLinea = salto === -1 ? value.length : salto
    const bloque = value.slice(inicioLinea, finLinea)
    const nuevoBloque = bloque
      .split('\n')
      .map((linea) => (linea.startsWith(prefijo) ? linea : prefijo + linea))
      .join('\n')
    const nuevo = value.slice(0, inicioLinea) + nuevoBloque + value.slice(finLinea)
    onChange(nuevo)
    enfocar(inicioLinea + nuevoBloque.length, inicioLinea + nuevoBloque.length)
  }

  return (
    <div className="field campo-markdown">
      <div className="campo-markdown-etiqueta">
        <label id={`${id}-lbl`} htmlFor={vista === 'escribir' ? id : undefined}>
          {label}
        </label>
        {opcional && (
          <span id={opcionalId} className="campo-markdown-opcional">
            Opcional
          </span>
        )}
      </div>

      <div
        className="campo-markdown-pestanas"
        role="tablist"
        aria-label="Editor de Markdown"
        onKeyDown={alPulsarTecla}
      >
        <button
          ref={pestanaEscribir}
          type="button"
          role="tab"
          id={`${id}-tab-e`}
          aria-selected={vista === 'escribir'}
          aria-controls={`${id}-pe`}
          tabIndex={vista === 'escribir' ? 0 : -1}
          onClick={() => elegir('escribir', false)}
        >
          Escribir
        </button>
        <button
          ref={pestanaVista}
          type="button"
          role="tab"
          id={`${id}-tab-p`}
          aria-selected={vista === 'preview'}
          aria-controls={`${id}-pp`}
          tabIndex={vista === 'preview' ? 0 : -1}
          onClick={() => elegir('preview', false)}
        >
          Vista previa
        </button>
      </div>

      <div role="tabpanel" id={`${id}-pe`} aria-labelledby={`${id}-tab-e`} hidden={vista !== 'escribir'}>
        <div className="campo-markdown-herramientas" role="toolbar" aria-label="Formato Markdown">
          <button
            type="button"
            onMouseDown={(evento) => evento.preventDefault()}
            onClick={() => prefijar('# ')}
            disabled={disabled}
          >
            Título
          </button>
          <button
            type="button"
            onMouseDown={(evento) => evento.preventDefault()}
            onClick={() => envolver('**', '**', 'negrita')}
            disabled={disabled}
          >
            Negrita
          </button>
          <button
            type="button"
            onMouseDown={(evento) => evento.preventDefault()}
            onClick={() => envolver('*', '*', 'cursiva')}
            disabled={disabled}
          >
            Cursiva
          </button>
          <button
            type="button"
            onMouseDown={(evento) => evento.preventDefault()}
            onClick={() => prefijar('- ')}
            disabled={disabled}
          >
            Lista
          </button>
          <button
            type="button"
            onMouseDown={(evento) => evento.preventDefault()}
            onClick={() => envolver('[', '](https://ejemplo.com)', 'texto')}
            disabled={disabled}
          >
            Enlace
          </button>
          <button
            type="button"
            onMouseDown={(evento) => evento.preventDefault()}
            onClick={() => envolver('`', '`', 'código')}
            disabled={disabled}
          >
            Código
          </button>
          <button
            type="button"
            onMouseDown={(evento) => evento.preventDefault()}
            onClick={() => envolver('![', '](https://ejemplo.com/imagen.png)', 'texto')}
            disabled={disabled}
          >
            Imagen
          </button>
        </div>


        <textarea
          ref={ref}
          id={id}
          name={name}
          value={value}
          onChange={(evento) => onChange(evento.target.value)}
          disabled={disabled}
          rows={rows}
          aria-invalid={Boolean(error)}
          aria-describedby={descripcion}
        />
      </div>

      <div
        role="tabpanel"
        id={`${id}-pp`}
        className="campo-markdown-vista"
        aria-labelledby={`${id}-tab-p`}
        hidden={vista !== 'preview'}
      >
        {vista === 'preview' &&
          (value.trim() === '' ? (
            <p className="campo-markdown-vacio">Nada para previsualizar</p>
          ) : (
            <Markdown texto={value} />
          ))}
      </div>

      {error && <ErrorMessage id={errorId} message={error} />}
      <p id={ayudaId} className="campo-markdown-ayuda">
        {AYUDA}
      </p>
    </div>
  )
}
