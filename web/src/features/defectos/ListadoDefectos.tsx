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
  alResolver: () => void
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
      await resolverDefecto(id)
      alResolver()
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
        {defectos.map((defecto) => (
          <li key={defecto.id} className="defecto">
            <div className="defecto__cabecera">
              <div>
                <p className="defecto__meta">Defecto #{defecto.id} · {defecto.historia}</p>
                <p className="defecto__descripcion">{defecto.descripcion}</p>
              </div>
              <div className="defecto__etiquetas">
                <span className={`defecto__severidad defecto__severidad--${defecto.severidad.toLowerCase()}`}>
                  {defecto.severidad}
                </span>
                <span className={`defecto__estado defecto__estado--${defecto.estado.toLowerCase()}`}>
                  {defecto.estado}
                </span>
              </div>
            </div>
            <p className="defecto__datos">
              <span>Detectado en {etiquetaSprint(defecto.sprint_deteccion_id, sprints)}</span>
              {defecto.sprint_resolucion_id !== null && (
                <span>Resuelto en {etiquetaSprint(defecto.sprint_resolucion_id, sprints)}</span>
              )}
              <span>Registrado: {defecto.creado_en.slice(0, 10)}</span>
            </p>
            {defecto.estado === 'Pendiente' && haySprintActivo && (
              <Button
                type="button"
                variant="secondary"
                onClick={() => resolver(defecto.id)}
                disabled={resolviendo !== null}
              >
                {resolviendo === defecto.id ? 'Resolviendo…' : 'Resolver'}
              </Button>
            )}
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
