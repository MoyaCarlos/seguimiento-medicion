import { Button } from '../../compartido/componentes/Button'
import type { Integrante } from './votacion'

type Props = {
  actor: Integrante | undefined
  puedeIniciar: boolean
  hayIntegrantes: boolean
  deshabilitado: boolean
  alIniciar: () => void
}

export function SinVotacion({ actor, puedeIniciar, hayIntegrantes, deshabilitado, alIniciar }: Props) {
  return (
    <div className="votacion-acciones">
      <p>La historia todavía no tiene votación</p>
      {puedeIniciar && (
        <Button type="button" disabled={deshabilitado} onClick={alIniciar}>
          Iniciar votación
        </Button>
      )}
      {actor && !puedeIniciar && <p>Solo el Scrum Master puede iniciar la votación</p>}
      {!actor && hayIntegrantes && <p>Elegí quién actúa para iniciar la votación</p>}
    </div>
  )
}
