import { useEffect, useState } from 'react'
import { Button } from '../../compartido/componentes/Button'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { consultarEsfuerzo } from './api'
import { clasificarError } from './errores'
import { FormularioHoras } from './FormularioHoras'
import { formatearHoras, type EsfuerzoReal } from './horas'

export type Imputacion = { proyectoId: number; fechaFin?: string | null }

type Props = { historiaId: number; imputacion?: Imputacion }

type Estado =
  | { tipo: 'cargando' }
  | { tipo: 'listo'; esfuerzo: EsfuerzoReal }
  | { tipo: 'error'; mensaje: string }

export function EsfuerzoHistoria({ historiaId, imputacion }: Props) {
  const [estado, setEstado] = useState<Estado>({ tipo: 'cargando' })
  const [consulta, setConsulta] = useState(0)
  const [detalleAbierto, setDetalleAbierto] = useState(false)
  const [formularioAbierto, setFormularioAbierto] = useState(false)
  const [enviando, setEnviando] = useState(false)

  useEffect(() => {
    let vigente = true

    consultarEsfuerzo(historiaId)
      .then((esfuerzo) => {
        if (vigente) setEstado({ tipo: 'listo', esfuerzo })
      })
      .catch((error: unknown) => {
        if (!vigente) return
        const clasificado = clasificarError(error, [])
        setEstado({
          tipo: 'error',
          mensaje: 'aviso' in clasificado ? clasificado.aviso : clasificado.mensaje,
        })
      })

    return () => {
      vigente = false
    }
  }, [historiaId, consulta])

  return (
    <div className="esfuerzo-historia">
      {estado.tipo === 'cargando' && <LoadingIndicator />}
      {estado.tipo === 'error' && <ErrorMessage message={estado.mensaje} />}
      {estado.tipo === 'listo' && (
        <>
          <p>{`Esfuerzo real: ${formatearHoras(estado.esfuerzo.total_horas)}`}</p>
          {estado.esfuerzo.registros.length > 0 && (
            <Button
              type="button"
              variant="quiet"
              aria-expanded={detalleAbierto}
              onClick={() => setDetalleAbierto((actual) => !actual)}
            >
              Ver detalle
            </Button>
          )}
          {detalleAbierto && estado.esfuerzo.registros.length > 0 && (
            <ul className="esfuerzo-detalle">
              {estado.esfuerzo.registros.map((registro) => (
                <li key={registro.id}>
                  {`${registro.fecha} · ${registro.integrante} · ${formatearHoras(registro.horas)}`}
                  <br />
                  {registro.actividad}
                </li>
              ))}
            </ul>
          )}
        </>
      )}
      {imputacion && (
        <Button
          type="button"
          variant="secondary"
          aria-expanded={formularioAbierto}
          disabled={enviando}
          onClick={() => setFormularioAbierto((actual) => !actual)}
        >
          Registrar horas
        </Button>
      )}
      {imputacion && formularioAbierto && (
        <FormularioHoras
          historiaId={historiaId}
          fechaFin={imputacion.fechaFin}
          alRegistrar={() => setConsulta((actual) => actual + 1)}
          alCambiarEnvio={setEnviando}
        />
      )}
    </div>
  )
}
