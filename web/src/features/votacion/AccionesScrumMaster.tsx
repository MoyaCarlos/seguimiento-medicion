import { Button } from '../../compartido/componentes/Button'
import { NotaCerrado } from '../../compartido/permisos/NotaCerrado'
import { NotaGestion } from '../../compartido/permisos/NotaGestion'
import { faltanPorVotar, puedeRevelar, type Votacion } from './votacion'

type Props = {
  votacion: Votacion
  deshabilitado: boolean
  alRevelar: () => void
  puedeGestionar?: boolean
  cerrado?: boolean
}

const NOTA = 'votacion-nota-revelar'

export function AccionesScrumMaster({
  votacion,
  deshabilitado,
  alRevelar,
  puedeGestionar = true,
  cerrado = false,
}: Props) {
  const faltan = faltanPorVotar(votacion)

  return (
    <div className="votacion-acciones">
      <Button
        type="button"
        disabled={deshabilitado || !puedeGestionar || !puedeRevelar(votacion)}
        aria-describedby={puedeGestionar ? undefined : NOTA}
        onClick={alRevelar}
      >
        Revelar
      </Button>
      {/* Una sola nota por acción: el cierre gana a la de HU-21. */}
      {cerrado && <NotaCerrado id={NOTA} />}
      {!cerrado && !puedeGestionar && (
        <NotaGestion id={NOTA} texto="Solo el Scrum Master puede revelar los votos" />
      )}
      {faltan.length > 0 && <p>Falta votar: {faltan.map((voto) => voto.nombre).join(', ')}</p>}
    </div>
  )
}
