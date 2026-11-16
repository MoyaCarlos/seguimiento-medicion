import { useEffect } from 'react'
import type { ReactNode } from 'react'
import { Icono } from '../../compartido/componentes/Icono'
import { Marca } from '../../compartido/componentes/Marca'
import type { NombreIcono } from '../../compartido/componentes/Icono'
import './estilos.css'

const PUNTOS: { icono: NombreIcono; texto: string }[] = [
  { icono: 'lista', texto: 'Historias priorizadas con Story Points' },
  { icono: 'bandera', texto: 'Sprints pendientes, activos y finalizados' },
  { icono: 'equipo', texto: 'Equipo con roles: Scrum Master y Product Builder' },
]

export function PantallaAcceso({ titulo, children }: { titulo: string; children: ReactNode }) {
  useEffect(() => {
    const anterior = document.title
    document.title = `Métrica · ${titulo}`
    return () => {
      document.title = anterior
    }
  }, [titulo])

  return (
    <div className="pantalla-acceso">
      <section className="pantalla-acceso__presentacion" aria-labelledby="presentacion-title">
        <Marca />
        <div>
          <h2 id="presentacion-title">Seguimiento Scrum, claro y al día.</h2>
          <p>Product Backlog, Sprints y estado del proyecto en un solo lugar para todo el equipo.</p>
        </div>
        <ul>
          {PUNTOS.map(({ icono, texto }) => (
            <li key={texto}>
              <span><Icono nombre={icono} /></span>
              {texto}
            </li>
          ))}
        </ul>
        <small>UTN FRSR · Ingeniería y Calidad de Software</small>
      </section>
      <div className="pantalla-acceso__formulario">
        <div className="pantalla-acceso__contenido">{children}</div>
      </div>
    </div>
  )
}
