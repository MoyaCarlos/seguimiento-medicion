import type { Avance } from './avance'

const COLORES: Record<string, string> = { Nueva: '#cbd5e1', 'En progreso': '#0d9488', Completada: '#16a34a' }

export function Resumen({ avance }: { avance: Avance }) {
  const { historias, story_points: storyPoints } = avance.total
  const detalle = avance.grupos.map((grupo) => `${grupo.estado}: ${grupo.cantidad}`).join(', ')

  return (
    <section className="estado-resumen" aria-label="Total del proyecto">
      <dl className="estado-resumen__totales">
        <div>
          <dt>{historias === 1 ? 'Historia' : 'Historias'}</dt>
          <dd>{historias}</dd>
        </div>
        <div>
          <dt>Story Points</dt>
          <dd>{storyPoints}</dd>
        </div>
      </dl>
      <div className="estado-resumen__barra" role="img" aria-label={detalle}>
        {avance.grupos
          .filter((grupo) => grupo.cantidad > 0)
          .map((grupo) => (
            <i key={grupo.estado} style={{ flex: grupo.cantidad, background: COLORES[grupo.estado] }} />
          ))}
      </div>
    </section>
  )
}
