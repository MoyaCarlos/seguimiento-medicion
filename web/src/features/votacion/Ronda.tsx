import { Cartas } from './Cartas'
import { Integrantes } from './Integrantes'
import { estaRevelada, type Votacion } from './votacion'

type Props = {
  votacion: Votacion
  deshabilitado: boolean
  alVotar: (valor: number) => void
}

export function Ronda({ votacion, deshabilitado, alVotar }: Props) {
  const abierta = !estaRevelada(votacion)
  // Mientras la ronda está abierta el servidor solo muestra el voto de la sesión.
  const miVoto = votacion.integrantes.find((voto) => voto.valor !== null)?.valor ?? null

  return (
    <>
      <p className="votacion-ronda">Ronda {votacion.ronda}</p>
      <Integrantes integrantes={votacion.integrantes} />
      {abierta && (
        <>
          <p className="votacion-consigna">Elegí una carta: ¿cuántos Story Points creés que lleva esta historia?</p>
          <Cartas seleccionada={miVoto} deshabilitada={deshabilitado} alElegir={alVotar} />
        </>
      )}
    </>
  )
}
