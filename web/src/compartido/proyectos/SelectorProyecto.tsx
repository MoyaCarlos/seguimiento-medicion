import { useEffect, useId, useRef, useState } from 'react'
import type { KeyboardEvent } from 'react'
import { useLocation, useNavigate, useParams } from 'react-router'
import { ErrorMessage } from '../componentes/ErrorMessage'
import { Icono } from '../componentes/Icono'
import { LoadingIndicator } from '../componentes/LoadingIndicator'
import { etiquetaProyecto, rutaDeSeccion } from './ruta'
import { useProyectos } from './useProyectos'

type Entrada = { valor: string; etiqueta: string }

const inicial = (nombre: string) => (nombre.trim().charAt(0) || '?').toUpperCase()

export function SelectorProyecto() {
  const { id = '' } = useParams()
  const { pathname } = useLocation()
  const navegar = useNavigate()
  const proyectos = useProyectos()
  const [abierto, setAbierto] = useState(false)
  const [resaltada, setResaltada] = useState(0)
  const raiz = useRef<HTMLDivElement>(null)
  const boton = useRef<HTMLButtonElement>(null)
  const lista = useRef<HTMLUListElement>(null)
  const base = useId()

  const listados = proyectos.estado === 'listo' ? proyectos.proyectos : []
  const actual = listados.find((proyecto) => String(proyecto.id) === id)
  const nombreActual = actual ? actual.nombre : `Proyecto ${id}`
  const entradas: Entrada[] = [
    ...(actual ? [] : [{ valor: id, etiqueta: `Proyecto ${id}` }]),
    ...listados.map((proyecto) => ({ valor: String(proyecto.id), etiqueta: etiquetaProyecto(proyecto) })),
  ]
  const idDeOpcion = (indice: number) => `${base}-opcion-${indice}`

  useEffect(() => {
    if (abierto) lista.current?.focus()
  }, [abierto])

  useEffect(() => {
    if (!abierto) return
    const alPulsar = (evento: PointerEvent) => {
      if (!raiz.current?.contains(evento.target as Node)) setAbierto(false)
    }
    document.addEventListener('pointerdown', alPulsar)
    return () => document.removeEventListener('pointerdown', alPulsar)
  }, [abierto])

  function abrir() {
    setResaltada(Math.max(0, entradas.findIndex((entrada) => entrada.valor === id)))
    setAbierto(true)
  }

  function cerrar(devolverFoco = true) {
    setAbierto(false)
    if (devolverFoco) boton.current?.focus()
  }

  function elegir(entrada: Entrada) {
    cerrar()
    if (entrada.valor !== id) navegar(rutaDeSeccion(pathname, entrada.valor))
  }

  function irANuevo() {
    cerrar()
    navegar('/')
  }

  function alTeclearEnLaLista(evento: KeyboardEvent<HTMLUListElement>) {
    const ultima = entradas.length - 1
    const movimientos: Record<string, number> = {
      ArrowDown: Math.min(resaltada + 1, ultima),
      ArrowUp: Math.max(resaltada - 1, 0),
      Home: 0,
      End: ultima,
    }
    if (evento.key in movimientos) {
      evento.preventDefault()
      setResaltada(movimientos[evento.key])
    } else if (evento.key === 'Enter' || evento.key === ' ') {
      evento.preventDefault()
      elegir(entradas[resaltada])
    }
  }

  return (
    <div
      ref={raiz}
      className={`selector-proyecto${abierto ? ' selector-proyecto--abierto' : ''}`}
      onKeyDown={(evento) => {
        if (evento.key === 'Escape' && abierto) {
          evento.stopPropagation()
          cerrar()
        }
      }}
      onBlur={(evento) => {
        const destino = evento.relatedTarget as Node | null
        if (abierto && destino && !raiz.current?.contains(destino)) cerrar(false)
      }}
    >
      <p id={`${base}-rotulo`} className="barra-lateral__rotulo">Proyecto actual</p>
      <button
        ref={boton}
        id={`${base}-boton`}
        type="button"
        className="selector-proyecto__boton"
        aria-haspopup="listbox"
        aria-expanded={abierto}
        aria-controls={`${base}-popup`}
        aria-labelledby={`${base}-rotulo ${base}-boton`}
        onClick={() => (abierto ? cerrar(false) : abrir())}
        onKeyDown={(evento) => {
          if (!abierto && (evento.key === 'ArrowDown' || evento.key === 'ArrowUp')) {
            evento.preventDefault()
            abrir()
          }
        }}
      >
        <span
          className={`selector-proyecto__avatar selector-proyecto__avatar--${Number(id) % 5}`}
          data-inicial={inicial(nombreActual)}
          aria-hidden="true"
        />
        <span className="selector-proyecto__texto">
          <strong>{nombreActual}</strong>
          <small>#{id}</small>
        </span>
        <Icono nombre="chevronAbajo" />
      </button>
      <div id={`${base}-popup`} className="selector-proyecto__popup" hidden={!abierto}>
        <ul
          ref={lista}
          role="listbox"
          aria-labelledby={`${base}-rotulo`}
          aria-activedescendant={idDeOpcion(resaltada)}
          tabIndex={-1}
          className="selector-proyecto__lista"
          onKeyDown={alTeclearEnLaLista}
        >
          {entradas.map((entrada, indice) => (
            <li
              key={entrada.valor}
              id={idDeOpcion(indice)}
              role="option"
              aria-selected={entrada.valor === id}
              className={`selector-proyecto__opcion${indice === resaltada ? ' selector-proyecto__opcion--resaltada' : ''}`}
              onClick={() => elegir(entrada)}
              onMouseEnter={() => setResaltada(indice)}
            >
              <span
                className={`selector-proyecto__avatar selector-proyecto__avatar--chico selector-proyecto__avatar--${Number(entrada.valor) % 5}`}
                data-inicial={inicial(entrada.etiqueta)}
                aria-hidden="true"
              />
              <span className="selector-proyecto__nombre">{entrada.etiqueta}</span>
            </li>
          ))}
        </ul>
        <button type="button" className="selector-proyecto__nuevo" onClick={irANuevo}>
          <Icono nombre="mas" />
          Nuevo proyecto
        </button>
      </div>
      {proyectos.estado === 'cargando' && <LoadingIndicator />}
      {proyectos.estado === 'error' && <ErrorMessage message={proyectos.mensaje} />}
    </div>
  )
}
