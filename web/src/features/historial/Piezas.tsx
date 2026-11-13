import type { ReactNode } from 'react'
const SIN_DATOS = 'sin datos'

export function Metrica({ nombre, valor }: { nombre: string; valor: string }) {
  return (
    <div className="historial-metrica">
      <dt>{nombre}</dt>
      <dd className={valor.startsWith(SIN_DATOS) ? 'historial-sin-datos' : undefined}>{valor}</dd>
    </div>
  )
}

export function BarraAvance({ completados, planificados }: { completados: number; planificados: number }) {
  const porcentaje = planificados > 0 ? Math.min((completados / planificados) * 100, 100) : 0

  return (
    <div
      className="historial-barra"
      role="meter"
      aria-label="Story Points completados sobre planificados"
      aria-valuemin={0}
      aria-valuemax={planificados > 0 ? planificados : 1}
      aria-valuenow={completados}
    >
      <span style={{ width: `${porcentaje}%` }} />
    </div>
  )
}

type Tono = 'neutro' | 'ok' | 'apagado' | 'fuerte' | 'alerta'

export function Chip({ tono = 'neutro', children }: { tono?: Tono; children: ReactNode }) {
  return <li className={`historial-chip historial-chip--${tono}`}>{children}</li>
}
