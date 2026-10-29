import { useState } from 'react'
import type { ChangeEvent, FormEvent } from 'react'
import { Link, useParams } from 'react-router'
import { ApiError, ConnectionError } from '../../compartido/api/client'
import { Button } from '../../compartido/componentes/Button'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { Field } from '../../compartido/componentes/Field'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { crearHistoria } from './api'
import { ESTADO_INICIAL, armarCuerpo, mapearError, parsearProyectoId, validarHistoria } from './historia'
import type { ErrorDeCampo, FormularioHistoria, HistoriaCreada } from './historia'
import { PRIORIDADES, VALORES_NEGOCIO, etiquetaPrioridad } from './opciones'
import './backlog.css'

export function Component() {
  const { id } = useParams()
  const proyectoId = parsearProyectoId(id)
  const deshabilitado = proyectoId === null

  const [formulario, setFormulario] = useState<FormularioHistoria>(ESTADO_INICIAL)
  const [enviando, setEnviando] = useState(false)
  const [confirmacion, setConfirmacion] = useState<HistoriaCreada | null>(null)
  const [errorGeneral, setErrorGeneral] = useState<string>()
  const [errores, setErrores] = useState<ErrorDeCampo>({})

  function cambiarTitulo(evento: ChangeEvent<HTMLInputElement>) {
    setFormulario((actual) => ({ ...actual, titulo: evento.target.value }))
  }

  function cambiarDescripcion(evento: ChangeEvent<HTMLTextAreaElement>) {
    setFormulario((actual) => ({ ...actual, descripcion: evento.target.value }))
  }

  function cambiarPrioridad(evento: ChangeEvent<HTMLSelectElement>) {
    setFormulario((actual) => ({
      ...actual,
      prioridad: evento.target.value as FormularioHistoria['prioridad'],
    }))
  }

  function cambiarValorNegocio(evento: ChangeEvent<HTMLSelectElement>) {
    setFormulario((actual) => ({
      ...actual,
      valorNegocio: evento.target.value as FormularioHistoria['valorNegocio'],
    }))
  }

  async function enviar(evento: FormEvent<HTMLFormElement>) {
    evento.preventDefault()
    if (enviando || proyectoId === null) return

    setConfirmacion(null)
    setErrorGeneral(undefined)
    setErrores({})

    const erroresDelIntento = validarHistoria(formulario)
    if (Object.keys(erroresDelIntento).length > 0) {
      setErrores(erroresDelIntento)
      return
    }

    setEnviando(true)

    try {
      const historia = await crearHistoria(armarCuerpo(proyectoId, formulario))
      setConfirmacion(historia)
      setFormulario(ESTADO_INICIAL)
    } catch (error) {
      if (error instanceof ApiError || error instanceof ConnectionError) {
        const mapeado = mapearError(error)

        if (mapeado.campo) {
          setErrores({ [mapeado.campo]: mapeado.mensaje } as ErrorDeCampo)
        } else {
          setErrorGeneral(mapeado.mensaje)
        }
      } else {
        setErrorGeneral(error instanceof Error ? error.message : 'No se pudo crear la historia')
      }
    } finally {
      setEnviando(false)
    }
  }

  return (
    <section className="backlog" aria-labelledby="backlog-title">
      <div className="page-header">
        <div>
          <p className="eyebrow">Product Backlog</p>
          <h1 id="backlog-title">Nueva historia</h1>
          <p className="project-context">
            Proyecto actual: <strong>{id}</strong>
          </p>
        </div>
        <Link className="button button--quiet" to={`/projects/${id}/backlog`}>
          Volver al backlog
        </Link>
      </div>

      {deshabilitado && (
        <ErrorMessage message="El identificador del proyecto de la dirección no es válido." />
      )}
      {errorGeneral && <ErrorMessage message={errorGeneral} />}

      <form className="backlog-form" onSubmit={enviar}>
        <fieldset disabled={deshabilitado}>
          <Field
            label="Título"
            name="titulo"
            value={formulario.titulo}
            onChange={cambiarTitulo}
            disabled={deshabilitado}
            error={errores.titulo}
          />

          <div className="field">
            <label htmlFor="descripcion">Descripción</label>
            <textarea
              id="descripcion"
              name="descripcion"
              value={formulario.descripcion}
              onChange={cambiarDescripcion}
              disabled={deshabilitado}
              aria-invalid={Boolean(errores.descripcion)}
              aria-describedby={errores.descripcion ? 'descripcion-error' : undefined}
            />
            {errores.descripcion && (
              <ErrorMessage id="descripcion-error" message={errores.descripcion} />
            )}
          </div>

          <div className="field">
            <label htmlFor="prioridad">Prioridad</label>
            <select
              id="prioridad"
              name="prioridad"
              value={formulario.prioridad}
              onChange={cambiarPrioridad}
              disabled={deshabilitado}
              aria-invalid={Boolean(errores.prioridad)}
              aria-describedby={errores.prioridad ? 'prioridad-error' : undefined}
            >
              {PRIORIDADES.map((prioridad) => (
                <option key={prioridad.codigo} value={prioridad.codigo}>
                  {prioridad.etiqueta}
                </option>
              ))}
            </select>
            {errores.prioridad && (
              <ErrorMessage id="prioridad-error" message={errores.prioridad} />
            )}
          </div>

          <div className="field">
            <label htmlFor="valor-negocio">Valor de negocio</label>
            <select
              id="valor-negocio"
              name="valorNegocio"
              value={formulario.valorNegocio}
              onChange={cambiarValorNegocio}
              disabled={deshabilitado}
              aria-invalid={Boolean(errores.valor_negocio)}
              aria-describedby={errores.valor_negocio ? 'valor-negocio-error' : undefined}
            >
              <option value="">Sin valor</option>
              {VALORES_NEGOCIO.map((valor) => (
                <option key={valor} value={String(valor)}>
                  {valor}
                </option>
              ))}
            </select>
            {errores.valor_negocio && (
              <ErrorMessage id="valor-negocio-error" message={errores.valor_negocio} />
            )}
          </div>

          {enviando && <LoadingIndicator />}

          <Button type="submit" disabled={enviando || deshabilitado}>
            Guardar
          </Button>
        </fieldset>
      </form>

      {confirmacion && (
        <section className="backlog-confirmacion" aria-label="Historia creada">
          <h2>Historia creada</h2>
          <p>{confirmacion.titulo}</p>
          <dl>
            <dt>Prioridad</dt>
            <dd>{etiquetaPrioridad(confirmacion.prioridad)}</dd>
            <dt>Estado</dt>
            <dd>{confirmacion.estado}</dd>
            <dt>Estimación</dt>
            <dd>
              {confirmacion.estimacion_sp === null
                ? 'sin estimación'
                : `${confirmacion.estimacion_sp} Story Points`}
            </dd>
            {confirmacion.valor_negocio !== null && (
              <>
                <dt>Valor de negocio</dt>
                <dd>{confirmacion.valor_negocio}</dd>
              </>
            )}
          </dl>
        </section>
      )}
    </section>
  )
}
