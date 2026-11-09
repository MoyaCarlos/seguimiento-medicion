import { useEffect, useState } from 'react'
import { useParams } from 'react-router'
import { ApiError, ConnectionError } from '../../compartido/api/client'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { consultarAvance } from './api'
import { esIdValido, type Avance } from './avance'
import { Cifras } from './Cifras'
import { GrupoHistorias } from './GrupoHistorias'
import './estado.css'

export function Component() {
  const { id = '' } = useParams()

  return (
    <div className="estado-pantalla">
      <h1>Estado del proyecto</h1>
      {esIdValido(id) ? (
        <>
          <p className="project-context">Proyecto actual: <strong>{id}</strong></p>
          <AvanceDelProyecto key={id} id={id} />
        </>
      ) : (
        <ErrorMessage message={`El proyecto "${id}" no es válido`} />
      )}
    </div>
  )
}

function mensajeDeError(causa: unknown): string {
  if (causa instanceof ApiError && causa.status === 404) return 'El proyecto no existe'
  if (causa instanceof ApiError || causa instanceof ConnectionError) return causa.message
  return 'Ocurrió un error inesperado'
}

function AvanceDelProyecto({ id }: { id: string }) {
  const [avance, setAvance] = useState<Avance>()
  const [error, setError] = useState<string>()

  useEffect(() => {
    let vigente = true

    consultarAvance(id)
      .then((obtenido) => {
        if (vigente) setAvance(obtenido)
      })
      .catch((causa: unknown) => {
        if (vigente) setError(mensajeDeError(causa))
      })

    return () => {
      vigente = false
    }
  }, [id])

  if (error) return <ErrorMessage message={error} />
  if (!avance) return <LoadingIndicator />

  return (
    <>
      {avance.total.historias === 0 && <p className="estado-vacio">El proyecto todavía no tiene historias.</p>}
      {avance.grupos.map((grupo) => <GrupoHistorias key={grupo.estado} grupo={grupo} />)}
      <section className="estado-tarjeta" aria-labelledby="estado-total">
        <h2 id="estado-total">Total del proyecto</h2>
        <Cifras historias={avance.total.historias} storyPoints={avance.total.story_points} />
      </section>
    </>
  )
}
