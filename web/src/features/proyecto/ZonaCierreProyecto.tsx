import { useEffect, useRef, useState } from 'react'
import type { KeyboardEvent } from 'react'
import { ApiError } from '../../compartido/api/client'
import { Button } from '../../compartido/componentes/Button'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { AvisoSinPermiso } from '../../compartido/permisos/AvisoSinPermiso'
import { cerrarProyecto } from './api'
import type { EstadoProyecto, Rol } from './api'
import { mensajeDeError } from './errores'

type Fase = 'reposo' | 'confirmando' | 'enviando'

type Props = {
  proyectoId: string
  rol: Rol | null
  fechaCierre: string | null
  onCambio: (fechaCierre: string | null) => void
}

export function ZonaCierreProyecto({ proyectoId, onCambio }: Props) {
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

  async function confirmar() {
    if (fase === 'enviando') return
    setFase('enviando')
    try {
      const proyecto = await cerrarProyecto(proyectoId)
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

      {fase === 'reposo' && <Button onClick={pedirConfirmacion}>Cerrar proyecto</Button>}

      {fase !== 'reposo' && (
        <div className="zona-cierre__confirmacion">
          <p>¿Cerrar el proyecto?</p>
          {fase === 'enviando' && <LoadingIndicator />}
          <div className="hu04-acciones">
            <Button onClick={confirmar} disabled={fase === 'enviando'}>
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
