import { useState } from 'react'
import { Button } from '../../compartido/componentes/Button'
import { FormularioHoras } from './FormularioHoras'

export type Imputacion = { proyectoId: number; fechaFin?: string | null }

type Props = { historiaId: number; imputacion?: Imputacion }

export function EsfuerzoHistoria({ historiaId, imputacion }: Props) {
  const [abierto, setAbierto] = useState(false)

  if (!imputacion) return null

  return (
    <div className="esfuerzo-historia">
      <Button
        type="button"
        variant="secondary"
        aria-expanded={abierto}
        onClick={() => setAbierto((actual) => !actual)}
      >
        Registrar horas
      </Button>
      {abierto && (
        <FormularioHoras
          historiaId={historiaId}
          proyectoId={imputacion.proyectoId}
          fechaFin={imputacion.fechaFin}
          alRegistrar={() => {}}
        />
      )}
    </div>
  )
}
