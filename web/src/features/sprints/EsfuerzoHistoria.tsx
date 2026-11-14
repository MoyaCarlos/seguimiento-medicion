import { useEffect, useState, type ReactNode } from 'react'
import { Button } from '../../compartido/componentes/Button'
import { Icono } from '../../compartido/componentes/Icono'
import { fechaCorta } from '../../compartido/formato/fechas'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { consultarEsfuerzo } from './api'
import { clasificarError } from './errores'
import { FormularioHoras } from './FormularioHoras'
import { formatearHoras, type EsfuerzoReal } from './horas'

export type Imputacion = { proyectoId: number; fechaFin?: string | null }

type Props = { historiaId: number; imputacion?: Imputacion; acciones?: ReactNode }

type Estado =
  | { tipo: 'cargando' }
  | { tipo: 'listo'; esfuerzo: EsfuerzoReal }
  | { tipo: 'error'; mensaje: string }

export function EsfuerzoHistoria({ historiaId, imputacion, acciones }: Props) {
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
          <p className="sprint-esfuerzo">{`Esfuerzo real: ${formatearHoras(estado.esfuerzo.total_horas)}`}</p>
          {estado.esfuerzo.registros.length > 0 && (
            <button
              type="button"
              className="sprint-detalle"
              aria-expanded={detalleAbierto}
              onClick={() => setDetalleAbierto((actual) => !actual)}
            >
              Ver detalle
            </button>
          )}
          {detalleAbierto && estado.esfuerzo.registros.length > 0 && (
            <ul className="esfuerzo-detalle">
              {estado.esfuerzo.registros.map((registro) => (
                <li key={registro.id}>
                  {`${fechaCorta(registro.fecha)} · ${registro.integrante} · ${formatearHoras(registro.horas)}`}
                  <br />
                  {registro.actividad}
                </li>
              ))}
            </ul>
          )}
        </>
      )}
      {(imputacion || acciones) && (
        <div className="sprint-historia__acciones">
          {imputacion && (
            <Button
              type="button"
              variant="secondary"
              className="sprint-boton"
              aria-expanded={formularioAbierto}
              disabled={enviando}
              onClick={() => setFormularioAbierto((actual) => !actual)}
            >
              <Icono nombre="reloj" />
              Registrar horas
            </Button>
          )}
          {acciones}
        </div>
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
