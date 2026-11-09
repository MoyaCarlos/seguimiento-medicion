type CifrasProps = {
  historias: number
  storyPoints: number
}

export function Cifras({ historias, storyPoints }: CifrasProps) {
  return (
    <dl className="estado-cifras">
      <div>
        <dt>Historias</dt>
        <dd>{historias}</dd>
      </div>
      <div>
        <dt>Story Points</dt>
        <dd>{storyPoints}</dd>
      </div>
    </dl>
  )
}
