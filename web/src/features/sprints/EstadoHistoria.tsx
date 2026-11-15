import { useEffect, useRef, useState } from 'react'
import type { KeyboardEvent } from 'react'
import { Button } from '../../compartido/componentes/Button'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { NotaCerrado } from '../../compartido/permisos/NotaCerrado'
import { cambiarEstadoHistoria } from './api'
import type { HistoriaDeSprint } from './sprint'

type Props = {
  historia: HistoriaDeSprint
  alCambiar: (historia: HistoriaDeSprint) => void
  puedeEscribir?: boolean
}

type Fase = 'reposo' | 'confirmando'

const ACCIONES = {
  Nueva: { etiqueta: 'Iniciar', destino: 'En progreso' },
  'En progreso': { etiqueta: 'Completar', destino: 'Completada' },
} as const

export function EstadoHistoria({ historia, alCambiar, puedeEscribir = true }: Props) {
  const [fase, setFase] = useState<Fase>('reposo')
  const [enviando, setEnviando] = useState(false)
  const [error, setError] = useState('')
  const botonCancelar = useRef<HTMLButtonElement>(null)

  useEffect(() => {
    if (fase === 'confirmando') botonCancelar.current?.focus()
  }, [fase])

  if (historia.estado === 'Completada') return null
  const { etiqueta, destino } = ACCIONES[historia.estado]
  const idNota = `historia-${historia.id}-nota-estado`

  function pedirConfirmacion() {
    if (enviando || !puedeEscribir) return
    setError('')
    setFase('confirmando')
  }

  function cancelar() {
    if (!enviando) setFase('reposo')
  }

  async function confirmar() {
    if (enviando || !puedeEscribir) return
    setError('')
    setEnviando(true)
    try {
      alCambiar(await cambiarEstadoHistoria(historia.id, destino))
    } catch (causa) {
      setError(causa instanceof Error ? causa.message : 'Ocurrió un error inesperado')
    } finally {
      setEnviando(false)
      setFase('reposo')
    }
  }

  function alPresionarTecla(evento: KeyboardEvent<HTMLDivElement>) {
    if (evento.key === 'Escape' && fase === 'confirmando') cancelar()
  }

  return (
    <div className="estado-historia" onKeyDown={alPresionarTecla}>
      {error && <ErrorMessage message={error} />}
      {fase === 'reposo' && (
        <Button
          type="button"
          variant="secondary"
          onClick={pedirConfirmacion}
          disabled={enviando || !puedeEscribir}
          aria-describedby={puedeEscribir ? undefined : idNota}
        >
          {etiqueta}
        </Button>
      )}
      {fase === 'confirmando' && (
        <div className="estado-historia__confirmacion">
          <p>{`¿${etiqueta} la historia?`}</p>
          <div className="estado-historia__acciones">
            <Button type="button" onClick={confirmar} disabled={enviando || !puedeEscribir}>
              Confirmar
            </Button>
            <button
              type="button"
              className="button button--quiet"
              ref={botonCancelar}
              onClick={cancelar}
              disabled={enviando}
            >
              Cancelar
            </button>
          </div>
        </div>
      )}
      {!puedeEscribir && <NotaCerrado id={idNota} />}
    </div>
  )
}
