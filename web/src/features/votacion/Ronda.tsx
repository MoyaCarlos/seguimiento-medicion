import { Cartas } from './Cartas'
import { Integrantes } from './Integrantes'
import { estaRevelada, type Integrante, type Votacion } from './votacion'

type Props = {
  votacion: Votacion
  actor: Integrante | undefined
  hayIntegrantes: boolean
  deshabilitado: boolean
  alVotar: (valor: number) => void
}

export function Ronda({ votacion, actor, hayIntegrantes, deshabilitado, alVotar }: Props) {
  const abierta = !estaRevelada(votacion)
  const miVoto = votacion.integrantes.find((voto) => voto.usuario_id === actor?.id)?.valor ?? null

  return (
    <>
      <p className="votacion-ronda">Ronda {votacion.ronda}</p>
      <Integrantes integrantes={votacion.integrantes} />
      {abierta && !actor && hayIntegrantes && <p>Elegí quién actúa para votar</p>}
      {abierta && actor && (
        <>
          <p className="votacion-consigna">Elegí una carta: ¿cuántos Story Points creés que lleva esta historia?</p>
          <Cartas seleccionada={miVoto} deshabilitada={deshabilitado} alElegir={alVotar} />
        </>
      )}
    </>
  )
}
