type Props = { promedio: { exacto: number; carta_sugerida: number } | null }

export function Resultado({ promedio }: Props) {
  if (promedio === null) return null

  return (
    <dl className="votacion-resultado" aria-label="Resultado de la votación">
      <dt>Promedio exacto</dt>
      <dd>{promedio.exacto.toLocaleString('es-AR', { maximumFractionDigits: 2 })}</dd>
      <dt>Carta sugerida</dt>
      <dd>{promedio.carta_sugerida}</dd>
    </dl>
  )
}
