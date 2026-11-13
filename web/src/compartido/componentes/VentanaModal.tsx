import { useEffect, useRef, useState } from 'react'
import type { KeyboardEvent, ReactNode } from 'react'
import { Icono } from './Icono'

type Props = { children: ReactNode; alCerrar(): void; lado?: 'derecha' }

const FOCALIZABLES =
  'a[href], button:not([disabled]), input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])'

export function VentanaModal({ children, alCerrar, lado }: Props) {
  const dialogo = useRef<HTMLDivElement>(null)
  const empezoEnElFondo = useRef(false)
  // Se toma al renderizar: los hijos con autoFocus ya habrían movido el foco cuando corre el efecto.
  const [abridor] = useState(() => document.activeElement as HTMLElement | null)

  const controles = () =>
    Array.from(dialogo.current?.querySelectorAll<HTMLElement>(FOCALIZABLES) ?? [])

  useEffect(() => {
    const anterior = document.body.style.overflow
    document.body.style.overflow = 'hidden'

    const preparar = () => {
      const caja = dialogo.current
      if (!caja) return
      const titulo = caja.querySelector<HTMLElement>('h1, h2, h3')
      if (titulo) {
        titulo.id ||= 'ventana-titulo'
        caja.setAttribute('aria-labelledby', titulo.id)
      }
      if (!caja.contains(document.activeElement)) {
        const [primero] = controles().filter((control) => !control.hasAttribute('data-cerrar'))
        ;(primero ?? caja).focus()
      }
    }
    preparar()
    const observador = new MutationObserver(preparar)
    if (dialogo.current) observador.observe(dialogo.current, { childList: true, subtree: true })

    return () => {
      observador.disconnect()
      document.body.style.overflow = anterior
      if (abridor?.isConnected) abridor.focus()
    }
  }, [abridor])

  function alTeclear(evento: KeyboardEvent<HTMLDivElement>) {
    if (evento.key === 'Escape') {
      evento.stopPropagation()
      alCerrar()
      return
    }
    if (evento.key !== 'Tab') return
    const todos = controles()
    if (todos.length === 0) return
    const primero = todos[0]
    const ultimo = todos[todos.length - 1]
    const activo = document.activeElement
    if (evento.shiftKey && (activo === primero || activo === dialogo.current)) {
      evento.preventDefault()
      ultimo.focus()
    } else if (!evento.shiftKey && activo === ultimo) {
      evento.preventDefault()
      primero.focus()
    }
  }

  return (
    <div
      className={`ventana-fondo${lado ? ' ventana-fondo--lateral' : ''}`}
      onMouseDown={(evento) => {
        empezoEnElFondo.current = evento.target === evento.currentTarget
      }}
      onClick={(evento) => {
        if (empezoEnElFondo.current && evento.target === evento.currentTarget) alCerrar()
      }}
      onKeyDown={alTeclear}
    >
      <div ref={dialogo} role="dialog" aria-modal="true" className={`ventana${lado ? ' ventana--lateral' : ''}`} tabIndex={-1}>
        {children}
        <button type="button" className="ventana__cerrar" aria-label="Cerrar" data-cerrar onClick={alCerrar}>
          <Icono nombre="cerrar" />
        </button>
      </div>
    </div>
  )
}
