import { useRef, useState } from 'react'
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
}: Props) {
  const ref = useRef<HTMLTextAreaElement>(null)
  const [vista, setVista] = useState<Vista>('escribir')
  const ayudaId = `${id}-ayuda`
  const errorId = `${id}-error`

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

  const boton = (etiqueta: string, accion: () => void) => (
    <button
      type="button"
      onMouseDown={(evento) => evento.preventDefault()}
      onClick={accion}
      disabled={disabled}
    >
      {etiqueta}
    </button>
  )

  return (
    <div className="field campo-markdown">
      <label htmlFor={id}>{label}</label>

      <div className="campo-markdown-pestanas" role="tablist" aria-label="Editor de Markdown">
        <button
          type="button"
          role="tab"
          aria-selected={vista === 'escribir'}
          onClick={() => setVista('escribir')}
        >
          Escribir
        </button>
        <button
          type="button"
          role="tab"
          aria-selected={vista === 'preview'}
          onClick={() => setVista('preview')}
        >
          Vista previa
        </button>
      </div>

      {vista === 'escribir' ? (
        <>
          <div className="campo-markdown-herramientas" role="toolbar" aria-label="Formato Markdown">
            {boton('Título', () => prefijar('# '))}
            {boton('Negrita', () => envolver('**', '**', 'negrita'))}
            {boton('Cursiva', () => envolver('*', '*', 'cursiva'))}
            {boton('Lista', () => prefijar('- '))}
            {boton('Enlace', () => envolver('[', '](https://ejemplo.com)', 'texto'))}
            {boton('Código', () => envolver('`', '`', 'código'))}
            {boton('Imagen', () => envolver('![', '](https://ejemplo.com/imagen.png)', 'texto'))}
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
            aria-describedby={error ? errorId : ayudaId}
          />
        </>
      ) : (
        <div className="campo-markdown-vista">
          {value.trim() === '' ? (
            <p className="campo-markdown-vacio">Nada para previsualizar</p>
          ) : (
            <Markdown texto={value} />
          )}
        </div>
      )}

      {error && <ErrorMessage id={errorId} message={error} />}
      <p id={ayudaId} className="campo-markdown-ayuda">
        {AYUDA}
      </p>
    </div>
  )
}
