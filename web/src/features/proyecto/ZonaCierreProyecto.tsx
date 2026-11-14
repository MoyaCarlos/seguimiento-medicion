import { useEffect, useRef, useState } from 'react'
import type { KeyboardEvent } from 'react'
import { ApiError } from '../../compartido/api/client'
import { Button } from '../../compartido/componentes/Button'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { AvisoSinPermiso } from '../../compartido/permisos/AvisoSinPermiso'
import { NotaGestion } from '../../compartido/permisos/NotaGestion'
import { cerrarProyecto, reabrirProyecto } from './api'
import type { EstadoProyecto, Rol } from './api'
import { mensajeDeError } from './errores'

type Fase = 'reposo' | 'confirmando' | 'enviando'

type Props = {
  proyectoId: string
  rol: Rol | null
  fechaCierre: string | null
  onCambio: (fechaCierre: string | null) => void
}

const ID_NOTA = 'zona-cierre-nota'

/**
 * Cerrar pide confirmación en línea; reabrir no (research D5).
 * Decide por `rol` y no por `puedeGestionar`: con el proyecto cerrado este último es false
 * para todos, y el Scrum Master tiene que poder reabrir (research D3).
 */
export function ZonaCierreProyecto({ proyectoId, rol, fechaCierre, onCambio }: Props) {
  const bloqueado = rol === 'product_builder'
  const describedBy = bloqueado ? ID_NOTA : undefined
  const [fase, setFase] = useState<Fase>('reposo')
  const [estado, setEstado] = useState<EstadoProyecto | null>(null)
  const [sinPermiso, setSinPermiso] = useState('')
  const [error, setError] = useState('')
  const botonCancelar = useRef<HTMLButtonElement>(null)

  useEffect(() => {
    if (fase === 'confirmando') botonCancelar.current?.focus()
  }, [fase])

  function pedirConfirmacion() {
    setError('')
    setSinPermiso('')
    setFase('confirmando')
  }

  function alPresionarTecla(evento: KeyboardEvent<HTMLElement>) {
    if (evento.key === 'Escape' && fase === 'confirmando') setFase('reposo')
  }

  async function enviar(accion: typeof cerrarProyecto) {
    if (fase === 'enviando') return
    setError('')
    setSinPermiso('')
    setFase('enviando')
    try {
      const proyecto = await accion(proyectoId)
      setEstado(proyecto.estado)
      onCambio(proyecto.fecha_cierre)
    } catch (causa) {
      if (causa instanceof ApiError && causa.tipo === 'sin_permiso') {
        setSinPermiso(causa.message)
      } else {
        setError(mensajeDeError(causa))
      }
    } finally {
      setFase('reposo')
    }
  }

  return (
    <section className="zona-cierre" aria-labelledby="zona-cierre-titulo" onKeyDown={alPresionarTecla}>
      <h3 id="zona-cierre-titulo">Cerrar proyecto</h3>

      {error && <ErrorMessage message={error} />}
      {sinPermiso && <AvisoSinPermiso mensaje={sinPermiso} />}
      {estado && <p className="zona-cierre__estado">{estado}</p>}

      {fechaCierre !== null && (
        <>
          <p className="zona-cierre__fecha">Proyecto cerrado el {fechaCierre}</p>
          {fase === 'enviando' && <LoadingIndicator />}
          <Button
            onClick={() => enviar(reabrirProyecto)}
            disabled={bloqueado || fase === 'enviando'}
            aria-describedby={describedBy}
          >
            Reabrir proyecto
          </Button>
        </>
      )}

      {fechaCierre === null && fase === 'reposo' && (
        <Button onClick={pedirConfirmacion} disabled={bloqueado} aria-describedby={describedBy}>
          Cerrar proyecto
        </Button>
      )}

      {bloqueado && <NotaGestion id={ID_NOTA} texto="Solo el Scrum Master puede cerrar o reabrir el proyecto" />}

      {fechaCierre === null && fase !== 'reposo' && (
        <div className="zona-cierre__confirmacion">
          <p>¿Cerrar el proyecto?</p>
          {fase === 'enviando' && <LoadingIndicator />}
          <div className="hu04-acciones">
            <Button onClick={() => enviar(cerrarProyecto)} disabled={fase === 'enviando'}>
              Confirmar
            </Button>
            <button
              type="button"
              className="button button--quiet"
              ref={botonCancelar}
              onClick={() => setFase('reposo')}
              disabled={fase === 'enviando'}
            >
              Cancelar
            </button>
          </div>
        </div>
      )}
    </section>
  )
}
