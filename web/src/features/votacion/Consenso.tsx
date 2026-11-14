import { useState, type FormEvent } from 'react'
import { Button } from '../../compartido/componentes/Button'
import { NotaCerrado } from '../../compartido/permisos/NotaCerrado'
import { NotaGestion } from '../../compartido/permisos/NotaGestion'
import { CARTAS, type Votacion } from './votacion'

type Props = {
  votacion: Votacion
  deshabilitado: boolean
  alRegistrar: (valor: number) => void
  alAbrirRonda: () => void
  puedeGestionar?: boolean
  cerrado?: boolean
}

/** Una sola nota por acción: el cierre gana a la de HU-21. */
function Nota({ id, cerrado, texto }: { id: string; cerrado: boolean; texto: string }) {
  return cerrado ? <NotaCerrado id={id} /> : <NotaGestion id={id} texto={texto} />
}

export function Consenso({
  votacion,
  deshabilitado,
  alRegistrar,
  alAbrirRonda,
  puedeGestionar = true,
  cerrado = false,
}: Props) {
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
            disabled={deshabilitado || !puedeGestionar}
            onChange={(evento) => setValor(Number(evento.target.value))}
          >
            {CARTAS.map((carta) => (
              <option key={carta} value={carta}>
                {carta}
              </option>
            ))}
          </select>
        </div>
        <Button
          type="submit"
          disabled={deshabilitado || !puedeGestionar}
          aria-describedby={puedeGestionar ? undefined : 'votacion-nota-estimacion'}
        >
          Registrar estimación
        </Button>
      </div>
      {!puedeGestionar && (
        <Nota
          id="votacion-nota-estimacion"
          cerrado={cerrado}
          texto="Solo el Scrum Master puede registrar la estimación"
        />
      )}
      <div className="votacion-consenso__ronda">
        <p>¿Hace falta seguir conversando? Abrí otra ronda: los votos actuales quedan en el historial.</p>
        <Button
          type="button"
          variant="secondary"
          disabled={deshabilitado || !puedeGestionar}
          aria-describedby={puedeGestionar ? undefined : 'votacion-nota-nueva-ronda'}
          onClick={alAbrirRonda}
        >
          Nueva ronda
        </Button>
        {!puedeGestionar && (
          <Nota
            id="votacion-nota-nueva-ronda"
            cerrado={cerrado}
            texto="Solo el Scrum Master puede abrir una nueva ronda"
          />
        )}
      </div>
    </form>
  )
}
