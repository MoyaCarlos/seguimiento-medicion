import { Button } from '../../compartido/componentes/Button'
import { NotaCerrado } from '../../compartido/permisos/NotaCerrado'
import { NotaGestion } from '../../compartido/permisos/NotaGestion'

type Props = {
  deshabilitado: boolean
  alIniciar: () => void
  puedeGestionar?: boolean
  cerrado?: boolean
}

const NOTA = 'votacion-nota-iniciar'

export function SinVotacion({
  deshabilitado,
  alIniciar,
  puedeGestionar = true,
  cerrado = false,
}: Props) {
  return (
    <div className="votacion-acciones">
      <p>La historia todavía no tiene votación</p>
      <Button
        type="button"
        disabled={deshabilitado || !puedeGestionar}
        aria-describedby={puedeGestionar ? undefined : NOTA}
        onClick={alIniciar}
      >
        Iniciar votación
      </Button>
      {/* Una sola nota por acción: el cierre gana a la de HU-21. */}
      {cerrado && <NotaCerrado id={NOTA} />}
      {!cerrado && !puedeGestionar && (
        <NotaGestion id={NOTA} texto="Solo el Scrum Master puede iniciar la votación" />
      )}
    </div>
  )
}
