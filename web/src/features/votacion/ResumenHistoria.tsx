import type { HistoriaResumen, Prioridad } from './votacion'

const ETIQUETA_PRIORIDAD: Record<Prioridad, string> = {
  M: 'Must have',
  S: 'Should have',
  C: 'Could have',
  W: "Won't have",
}

type Props = { historiaId: number; historia: HistoriaResumen | null }

export function ResumenHistoria({ historiaId, historia }: Props) {
  if (historia === null) return <p className="votacion-historia__numero">Historia #{historiaId}</p>

  return (
    <article className="votacion-historia" aria-labelledby="votacion-historia-titulo">
      <p className="votacion-historia__numero">Historia #{historia.id}</p>
      <h2 id="votacion-historia-titulo">{historia.titulo}</h2>
      <p className="votacion-historia__datos">
        <span>{ETIQUETA_PRIORIDAD[historia.prioridad]}</span>
        <span>{historia.estimacion_sp === null ? 'Sin estimar' : `${historia.estimacion_sp} SP`}</span>
      </p>
      <p className="votacion-historia__descripcion">{historia.descripcion}</p>
    </article>
  )
}
