import { Button } from '../../compartido/componentes/Button'
import { NotaGestion } from '../../compartido/permisos/NotaGestion'
import type { Integrante } from './votacion'

type Props = {
  actor: Integrante | undefined
  puedeIniciar: boolean
  hayIntegrantes: boolean
  deshabilitado: boolean
  alIniciar: () => void
  puedeGestionar?: boolean
}

export function SinVotacion({
  actor,
  puedeIniciar,
  hayIntegrantes,
  deshabilitado,
  alIniciar,
  puedeGestionar = true,
}: Props) {
  return (
    <div className="votacion-acciones">
      <p>La historia todavía no tiene votación</p>
      {(puedeIniciar || !puedeGestionar) && (
        <Button
          type="button"
          disabled={deshabilitado || !puedeGestionar}
          aria-describedby={puedeGestionar ? undefined : 'votacion-nota-iniciar'}
          onClick={alIniciar}
        >
          Iniciar votación
        </Button>
      )}
      {!puedeGestionar && (
        <NotaGestion id="votacion-nota-iniciar" texto="Solo el Scrum Master puede iniciar la votación" />
      )}
      {puedeGestionar && actor && !puedeIniciar && <p>Solo el Scrum Master puede iniciar la votación</p>}
      {puedeGestionar && !actor && hayIntegrantes && <p>Elegí quién actúa para iniciar la votación</p>}
    </div>
  )
}
