import { useState } from 'react'
import { Button } from '../../compartido/componentes/Button'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { cambiarEstadoHistoria } from './api'
import type { HistoriaDeSprint } from './sprint'

type Props = {
  historia: HistoriaDeSprint
  alCambiar: (historia: HistoriaDeSprint) => void
}

const ACCIONES = {
  Nueva: { etiqueta: 'Iniciar', destino: 'En progreso' },
  'En progreso': { etiqueta: 'Completar', destino: 'Completada' },
} as const

export function EstadoHistoria({ historia, alCambiar }: Props) {
  const [enviando, setEnviando] = useState(false)
  const [error, setError] = useState('')

  if (historia.estado === 'Completada') return null
  const { etiqueta, destino } = ACCIONES[historia.estado]

  async function cambiar() {
    if (enviando) return
    setError('')
    setEnviando(true)
    try {
      alCambiar(await cambiarEstadoHistoria(historia.id, destino))
    } catch (causa) {
      setError(causa instanceof Error ? causa.message : 'Ocurrió un error inesperado')
    } finally {
      setEnviando(false)
    }
  }

  return (
    <div>
      {error && <ErrorMessage message={error} />}
      <Button type="button" variant="secondary" onClick={cambiar} disabled={enviando}>
        {etiqueta}
      </Button>
    </div>
  )
}
