import { textoDeEstimacion, textoDeSprint, type GrupoAvance } from './avance'

const plural = (cantidad: number) => `${cantidad} ${cantidad === 1 ? 'historia' : 'historias'}`

export function GrupoHistorias({ grupo }: { grupo: GrupoAvance }) {
  const clave = grupo.estado.toLowerCase().replace(/\s+/g, '-')
  const titulo = `estado-grupo-${clave}`

  return (
    <section className={`estado-columna estado-columna--${clave}`} aria-labelledby={titulo}>
      <div className="estado-columna__cabecera">
        <h2 id={titulo}>{grupo.estado}</h2>
        <p className="estado-columna__chips">
          <span>{plural(grupo.cantidad)}</span>
          <span>{grupo.story_points} SP</span>
        </p>
      </div>
      {grupo.historias.length === 0 ? (
        <p className="estado-grupo-vacio">No hay historias en este grupo.</p>
      ) : (
        <ul className="estado-lista">
          {grupo.historias.map((historia) => (
            <li key={historia.id} className="estado-historia">
              <h3>{historia.titulo}</h3>
              <p className="estado-historia__chips">
                <span className={historia.estimacion_sp === null ? 'estado-historia__sin-estimar' : undefined}>
                  Story Points: {textoDeEstimacion(historia.estimacion_sp)}
                </span>
              </p>
              <span className={`estado-historia__sprint${historia.sprint ? '' : ' estado-historia__sprint--backlog'}`}>
                {textoDeSprint(historia.sprint)}
              </span>
            </li>
          ))}
        </ul>
      )}
    </section>
  )
}
