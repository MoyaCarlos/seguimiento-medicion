import { Button } from '../../compartido/componentes/Button'
import { NotaGestion } from '../../compartido/permisos/NotaGestion'
import { faltanPorVotar, puedeRevelar, type Votacion } from './votacion'

type Props = {
  votacion: Votacion
  deshabilitado: boolean
  alRevelar: () => void
  puedeGestionar?: boolean
}

export function AccionesScrumMaster({
  votacion,
  deshabilitado,
  alRevelar,
  puedeGestionar = true,
}: Props) {
  const faltan = faltanPorVotar(votacion)

  return (
    <div className="votacion-acciones">
      <Button
        type="button"
        disabled={deshabilitado || !puedeGestionar || !puedeRevelar(votacion)}
        aria-describedby={puedeGestionar ? undefined : 'votacion-nota-revelar'}
        onClick={alRevelar}
      >
        Revelar
      </Button>
      {!puedeGestionar && (
        <NotaGestion id="votacion-nota-revelar" texto="Solo el Scrum Master puede revelar los votos" />
      )}
      {faltan.length > 0 && <p>Falta votar: {faltan.map((voto) => voto.nombre).join(', ')}</p>}
    </div>
  )
}
