import { Button } from '../../compartido/componentes/Button'
import { faltanPorVotar, puedeRevelar, type Integrante, type Votacion } from './votacion'

type Props = {
  votacion: Votacion
  scrumMaster: Integrante
  deshabilitado: boolean
  alRevelar: () => void
}

export function AccionesScrumMaster({ votacion, scrumMaster, deshabilitado, alRevelar }: Props) {
  const faltan = faltanPorVotar(votacion)

  return (
    <div className="votacion-acciones">
      <Button type="button" disabled={deshabilitado || !puedeRevelar(votacion, scrumMaster)} onClick={alRevelar}>
        Revelar
      </Button>
      {faltan.length > 0 && <p>Falta votar: {faltan.map((voto) => voto.nombre).join(', ')}</p>}
    </div>
  )
}
