import { useState } from 'react'
import type { ChangeEvent, FormEvent } from 'react'
import { Button } from '../../compartido/componentes/Button'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { Field } from '../../compartido/componentes/Field'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { editarHistoria } from './api'
import {
  armarCuerpoEdicion,
  campoDeEdicion,
  campoDeError,
  formularioDe,
  mensajeDeError,
  validarEdicion,
} from './edicionHistoria'
import type { ErrorEdicion, FormularioEdicion, Historia } from './historia'
import { STORY_POINTS } from './opciones'

type Props = {
  historia: Historia
  onGuardada: (historia: Historia) => void
  onCancelar: () => void
}

export function Edicion({ historia, onGuardada, onCancelar }: Props) {
  const [formulario, setFormulario] = useState<FormularioEdicion>(() => formularioDe(historia))
  const [errores, setErrores] = useState<ErrorEdicion>({})
  const [errorGeneral, setErrorGeneral] = useState('')
  const [enviando, setEnviando] = useState(false)

  function actualizar<Campo extends keyof FormularioEdicion>(
    campo: Campo,
    valor: FormularioEdicion[Campo],
  ) {
    setFormulario((actual) => ({ ...actual, [campo]: valor }))
    setErrores((actual) => ({ ...actual, [campo]: undefined }))
    setErrorGeneral('')
  }

  async function guardar(evento: FormEvent<HTMLFormElement>) {
    evento.preventDefault()
    if (enviando) return

    setErrorGeneral('')
    const validacion = validarEdicion(formulario)
    setErrores(validacion)
    if (Object.keys(validacion).length > 0) return

    setEnviando(true)
    try {
      const actualizada = await editarHistoria(historia.id, armarCuerpoEdicion(formulario))
      onGuardada(actualizada)
    } catch (error) {
      const campo = campoDeError(error)
      const campoEditable = campo ? campoDeEdicion(campo) : undefined
      if (campoEditable) {
        setErrores({ [campoEditable]: mensajeDeError(error) })
      } else {
        setErrorGeneral(mensajeDeError(error))
      }
    } finally {
      setEnviando(false)
    }
  }

  function cambiarEstimacion(evento: ChangeEvent<HTMLSelectElement>) {
    actualizar('estimacion', evento.target.value as FormularioEdicion['estimacion'])
  }

  const id = (campo: string) => `edicion-${historia.id}-${campo}`

  return (
    <form
      className="backlog-edicion"
      aria-label={`Editar ${historia.titulo}`}
      onSubmit={guardar}
      noValidate
    >
      <h3>Editar historia</h3>
      {errorGeneral && <ErrorMessage message={errorGeneral} />}

      <Field
        id={id('titulo')}
        name="titulo"
        label="Título"
        value={formulario.titulo}
        onChange={(evento) => actualizar('titulo', evento.target.value)}
        disabled={enviando}
        error={errores.titulo}
      />

      <div className="field">
        <label htmlFor={id('descripcion')}>Descripción</label>
        <textarea
          id={id('descripcion')}
          name="descripcion"
          value={formulario.descripcion}
          onChange={(evento) => actualizar('descripcion', evento.target.value)}
          disabled={enviando}
          aria-invalid={Boolean(errores.descripcion)}
          aria-describedby={errores.descripcion ? `${id('descripcion')}-error` : undefined}
        />
        {errores.descripcion && (
          <ErrorMessage id={`${id('descripcion')}-error`} message={errores.descripcion} />
        )}
      </div>

      <div className="field">
        <label htmlFor={id('criterios')}>Criterios de aceptación</label>
        <textarea
          id={id('criterios')}
          name="criterios"
          value={formulario.criterios}
          onChange={(evento) => actualizar('criterios', evento.target.value)}
          disabled={enviando}
          aria-invalid={Boolean(errores.criterios)}
          aria-describedby={errores.criterios ? `${id('criterios')}-error` : undefined}
        />
        {errores.criterios && (
          <ErrorMessage id={`${id('criterios')}-error`} message={errores.criterios} />
        )}
      </div>

      <div className="field">
        <label htmlFor={id('estimacion')}>Story Points</label>
        <select
          id={id('estimacion')}
          name="estimacion"
          value={formulario.estimacion}
          onChange={cambiarEstimacion}
          disabled={enviando}
          aria-invalid={Boolean(errores.estimacion)}
          aria-describedby={errores.estimacion ? `${id('estimacion')}-error` : undefined}
        >
          <option value="">Sin estimar</option>
          {STORY_POINTS.map((storyPoint) => (
            <option key={storyPoint} value={String(storyPoint)}>
              {storyPoint}
            </option>
          ))}
        </select>
        {errores.estimacion && (
          <ErrorMessage id={`${id('estimacion')}-error`} message={errores.estimacion} />
        )}
      </div>

      {enviando && <LoadingIndicator />}

      <div className="backlog-acciones">
        <Button type="button" variant="secondary" onClick={onCancelar} disabled={enviando}>
          Cancelar
        </Button>
        <Button type="submit" disabled={enviando}>
          Guardar cambios
        </Button>
      </div>
    </form>
  )
}
