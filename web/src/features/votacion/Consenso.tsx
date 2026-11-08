import { useState, type FormEvent } from 'react'
import { Button } from '../../compartido/componentes/Button'
import { CARTAS, type Votacion } from './votacion'

type Props = {
  votacion: Votacion
  deshabilitado: boolean
  alRegistrar: (valor: number) => void
}

export function Consenso({ votacion, deshabilitado, alRegistrar }: Props) {
  const [valor, setValor] = useState<number>(votacion.promedio?.carta_sugerida ?? CARTAS[0])

  function registrar(evento: FormEvent<HTMLFormElement>) {
    evento.preventDefault()
    alRegistrar(valor)
  }

  return (
    <form className="votacion-consenso" aria-labelledby="votacion-consenso-titulo" onSubmit={registrar}>
      <h2 id="votacion-consenso-titulo">Consenso del equipo</h2>
      <p className="votacion-consenso__ayuda">
        Registrá el valor que acordó el equipo. Puede ser cualquier carta de la baraja, aunque nadie la haya votado.
      </p>
      <div className="votacion-consenso__fila">
        <div className="field">
          <label htmlFor="votacion-estimacion">Estimación acordada (Story Points)</label>
          <select
            id="votacion-estimacion"
            value={valor}
            disabled={deshabilitado}
            onChange={(evento) => setValor(Number(evento.target.value))}
          >
            {CARTAS.map((carta) => (
              <option key={carta} value={carta}>
                {carta}
              </option>
            ))}
          </select>
        </div>
        <Button type="submit" disabled={deshabilitado}>
          Registrar estimación
        </Button>
      </div>
    </form>
  )
}
