import { textoDeDiferencias, type Diferencias } from './votacion'

type Props = {
  promedio: { exacto: number; carta_sugerida: number } | null
  diferencias?: Diferencias | null
}

export function Resultado({ promedio, diferencias = null }: Props) {
  if (promedio === null) return null

  const texto = textoDeDiferencias(diferencias)
  const hayDiferencias = diferencias?.hay_diferencias === true

  return (
    <section className="votacion-resultado" aria-label="Resultado de la votación">
      {texto && (
        <p className={`votacion-diferencias${hayDiferencias ? ' votacion-diferencias--hay' : ''}`}>{texto}</p>
      )}
      <dl className="votacion-resultado__datos">
        <div>
          <dt>Promedio exacto</dt>
          <dd>{promedio.exacto.toLocaleString('es-AR', { maximumFractionDigits: 2 })}</dd>
        </div>
        <div>
          <dt>Carta sugerida</dt>
          <dd>{promedio.carta_sugerida}</dd>
        </div>
      </dl>
    </section>
  )
}
