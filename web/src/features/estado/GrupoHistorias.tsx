import { textoDeEstimacion, textoDeSprint, type GrupoAvance } from './avance'
import { Cifras } from './Cifras'

export function GrupoHistorias({ grupo }: { grupo: GrupoAvance }) {
  const titulo = `estado-grupo-${grupo.estado.toLowerCase().replace(/\s+/g, '-')}`

  return (
    <section className="estado-tarjeta" aria-labelledby={titulo}>
      <h2 id={titulo}>{grupo.estado}</h2>
      <Cifras historias={grupo.cantidad} storyPoints={grupo.story_points} />
      {grupo.historias.length === 0 ? (
        <p className="estado-grupo-vacio">No hay historias en este grupo.</p>
      ) : (
        <ul className="estado-lista">
          {grupo.historias.map((historia) => (
            <li key={historia.id} className="estado-historia">
              <h3>{historia.titulo}</h3>
              <p className="estado-historia__datos">
                <span>Story Points: {textoDeEstimacion(historia.estimacion_sp)}</span>
                <span>{textoDeSprint(historia.sprint)}</span>
              </p>
            </li>
          ))}
        </ul>
      )}
    </section>
  )
}
