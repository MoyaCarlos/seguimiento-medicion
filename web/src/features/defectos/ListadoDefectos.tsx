import { useState } from 'react'
import { Button } from '../../compartido/componentes/Button'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { resolverDefecto, type Defecto, type SprintResumen } from './api'
import { mensajeDeError } from './errores'

type Props = {
  defectos: Defecto[]
  cargando: boolean
  error: string
  haySprintActivo: boolean
  sprints: SprintResumen[]
  alResolver: (mensaje: string) => void
}

export function ListadoDefectos({
  defectos,
  cargando,
  error,
  haySprintActivo,
  sprints,
  alResolver,
}: Props) {
  const [resolviendo, setResolviendo] = useState<number | null>(null)
  const [errorResolucion, setErrorResolucion] = useState('')

  async function resolver(id: number) {
    if (resolviendo !== null) return
    setErrorResolucion('')
    setResolviendo(id)
    try {
      const resuelto = await resolverDefecto(id)
      alResolver(`Defecto #${id} resuelto en el ${etiquetaSprint(resuelto.sprint_resolucion_id ?? 0, sprints)}`)
    } catch (causa) {
      setErrorResolucion(mensajeDeError(causa))
    } finally {
      setResolviendo(null)
    }
  }

  if (cargando) return <LoadingIndicator />
  if (error) return <ErrorMessage message={error} />
  if (defectos.length === 0) return <p className="defectos-vacio">Todavía no hay defectos registrados.</p>

  return (
    <div className="defectos-resultado">
      {errorResolucion && <ErrorMessage message={errorResolucion} />}
      <ul className="defectos-lista">
        {[...defectos].sort((a, b) => b.id - a.id).map((defecto) => (
          <li key={defecto.id} className="defecto">
            <div className="defecto__cabecera">
              <h3>Defecto #{defecto.id} · {defecto.historia}</h3>
              {defecto.estado === 'Pendiente' && haySprintActivo && (
                <Button
                  type="button"
                  variant="secondary"
                  className="defecto__resolver"
                  onClick={() => resolver(defecto.id)}
                  disabled={resolviendo !== null}
                >
                  {resolviendo === defecto.id ? 'Resolviendo…' : 'Resolver'}
                </Button>
              )}
            </div>
            <p className="defecto__descripcion">{defecto.descripcion}</p>
            <ul className="defecto__etiquetas">
              <li className={`defecto__etiqueta defecto__etiqueta--${defecto.severidad.toLowerCase()}`}>
                Severidad {defecto.severidad}
              </li>
              <li className={`defecto__etiqueta defecto__etiqueta--${defecto.estado.toLowerCase()}`}>
                {defecto.estado}
              </li>
            </ul>
            <p className="defecto__datos">
              {[
                `Detectado en ${etiquetaSprint(defecto.sprint_deteccion_id, sprints)}`,
                defecto.sprint_resolucion_id !== null && `Resuelto en ${etiquetaSprint(defecto.sprint_resolucion_id, sprints)}`,
                `Registrado: ${formatearFecha(defecto.creado_en)}`,
              ]
                .filter(Boolean)
                .join(' · ')}
            </p>
          </li>
        ))}
      </ul>
    </div>
  )
}

function etiquetaSprint(id: number, sprints: SprintResumen[]): string {
  const sprint = sprints.find((actual) => actual.id === id)
  return sprint ? `Sprint ${sprint.numero}` : `Sprint ${id}`
}

function formatearFecha(iso: string): string {
  const [anio, mes, dia] = iso.slice(0, 10).split('-')
  return `${dia}/${mes}/${anio}`
}
