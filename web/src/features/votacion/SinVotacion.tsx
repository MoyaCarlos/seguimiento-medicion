import { Button } from '../../compartido/componentes/Button'
import { NotaGestion } from '../../compartido/permisos/NotaGestion'

type Props = {
  deshabilitado: boolean
  alIniciar: () => void
  puedeGestionar?: boolean
}

export function SinVotacion({
  deshabilitado,
  alIniciar,
  puedeGestionar = true,
}: Props) {
  return (
    <div className="votacion-acciones">
      <p>La historia todavía no tiene votación</p>
      <Button
        type="button"
        disabled={deshabilitado || !puedeGestionar}
        aria-describedby={puedeGestionar ? undefined : 'votacion-nota-iniciar'}
        onClick={alIniciar}
      >
        Iniciar votación
      </Button>
      {!puedeGestionar && (
        <NotaGestion id="votacion-nota-iniciar" texto="Solo el Scrum Master puede iniciar la votación" />
      )}
    </div>
  )
}
