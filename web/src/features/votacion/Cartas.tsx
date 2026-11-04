import { CARTAS } from './votacion'

type Props = {
  seleccionada?: number | null
  deshabilitada?: boolean
  alElegir: (valor: number) => void
}

export function Cartas({ seleccionada = null, deshabilitada = false, alElegir }: Props) {
  return (
    <div className="votacion-cartas" role="group" aria-label="Cartas de estimación">
      {CARTAS.map((carta) => (
        <button
          key={carta}
          type="button"
          className="votacion-carta"
          aria-pressed={seleccionada === carta}
          disabled={deshabilitada}
          onClick={() => alElegir(carta)}
        >
          {carta}
        </button>
      ))}
    </div>
  )
}
