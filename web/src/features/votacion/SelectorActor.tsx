import { ETIQUETA_ROL, type Integrante } from './votacion'

type Props = {
  integrantes: Integrante[]
  actorId: number | undefined
  deshabilitado: boolean
  alElegir: (id: number | undefined) => void
}

export function SelectorActor({ integrantes, actorId, deshabilitado, alElegir }: Props) {
  return (
    <div className="field">
      <label htmlFor="votacion-actor">Quién actúa</label>
      <select
        id="votacion-actor"
        value={actorId ?? ''}
        disabled={deshabilitado}
        onChange={(evento) => alElegir(evento.target.value ? Number(evento.target.value) : undefined)}
      >
        <option value="">Elegí un integrante</option>
        {integrantes.map((integrante) => (
          <option key={integrante.id} value={integrante.id}>
            {integrante.nombre} ({ETIQUETA_ROL[integrante.rol]})
          </option>
        ))}
      </select>
    </div>
  )
}
