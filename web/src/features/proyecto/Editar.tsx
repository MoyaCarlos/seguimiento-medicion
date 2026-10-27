import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'
import { Link, useNavigate, useParams } from 'react-router'
import { Button } from '../../compartido/componentes/Button'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { Field } from '../../compartido/componentes/Field'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { aProyectoEditado, editarProyecto, obtenerProyecto } from './api'
import { CampoDescripcion } from './Campos'
import { campoDeError, mensajeDeError } from './errores'
import { validarEdicion } from './validacion'
import type { DatosProyecto, ErroresProyecto } from './validacion'

const FORMULARIO_VACIO: DatosProyecto = {
  nombre: '',
  descripcion: '',
  fecha_inicio: '',
  fecha_fin: '',
  creador: '',
}

function esCampoDeProyecto(campo: string): campo is keyof DatosProyecto {
  return campo in FORMULARIO_VACIO
}

export function Editar() {
  const { id = '' } = useParams()
  const navegar = useNavigate()
  const [datos, setDatos] = useState(FORMULARIO_VACIO)
  const [errores, setErrores] = useState<ErroresProyecto>({})
  const [errorGeneral, setErrorGeneral] = useState('')
  const [errorCarga, setErrorCarga] = useState('')
  const [cargando, setCargando] = useState(true)
  const [enviando, setEnviando] = useState(false)

  useEffect(() => {
    let vigente = true

    obtenerProyecto(id)
      .then((proyecto) => {
        if (!vigente) return
        setErrorCarga('')
        setDatos({
          nombre: proyecto.nombre,
          descripcion: proyecto.descripcion,
          fecha_inicio: proyecto.fecha_inicio ?? '',
          fecha_fin: proyecto.fecha_fin ?? '',
          creador: '',
        })
      })
      .catch((causa: unknown) => {
        if (vigente) setErrorCarga(mensajeDeError(causa))
      })
      .finally(() => {
        if (vigente) setCargando(false)
      })

    return () => {
      vigente = false
    }
  }, [id])

  function actualizar(campo: keyof DatosProyecto, valor: string) {
    setDatos((actuales) => ({ ...actuales, [campo]: valor }))
    setErrores((actuales) => ({ ...actuales, [campo]: undefined }))
    setErrorGeneral('')
  }

  async function guardar(evento: FormEvent<HTMLFormElement>) {
    evento.preventDefault()
    setErrorGeneral('')

    const validacion = validarEdicion(datos)
    setErrores(validacion)
    if (Object.keys(validacion).length > 0) return

    setEnviando(true)
    try {
      await editarProyecto(id, aProyectoEditado(datos))
      navegar(`/projects/${id}/proyecto`)
    } catch (error) {
      const campo = campoDeError(error)
      if (campo && esCampoDeProyecto(campo)) {
        setErrores((actuales) => ({ ...actuales, [campo]: mensajeDeError(error) }))
      } else {
        setErrorGeneral(mensajeDeError(error))
      }
    } finally {
      setEnviando(false)
    }
  }

  return (
    <section className="hu04-tarjeta hu04-tarjeta--angosta" aria-labelledby="editar-title">
      <p className="eyebrow">Proyecto · Edición</p>
      <h2 id="editar-title">Editar proyecto</h2>
      <p className="hu04-intro">
        Corregí el nombre, la descripción o las fechas. Los cambios se guardan sobre el
        proyecto actual.
      </p>

      {cargando && <LoadingIndicator />}
      {errorCarga && <ErrorMessage message={errorCarga} />}

      {!cargando && !errorCarga && (
        <form onSubmit={guardar} noValidate aria-label="Editar proyecto">
          {errorGeneral && <ErrorMessage message={errorGeneral} />}

          <Field
            id="nombre"
            name="nombre"
            label="Nombre del proyecto"
            value={datos.nombre}
            onChange={(evento) => actualizar('nombre', evento.target.value)}
            error={errores.nombre}
          />

          <CampoDescripcion
            value={datos.descripcion}
            onChange={(valor) => actualizar('descripcion', valor)}
            error={errores.descripcion}
          />

          <div className="hu04-fila">
            <Field
              id="fecha_inicio"
              name="fecha_inicio"
              type="date"
              label="Fecha de inicio"
              value={datos.fecha_inicio}
              onChange={(evento) => actualizar('fecha_inicio', evento.target.value)}
              error={errores.fecha_inicio}
            />
            <Field
              id="fecha_fin"
              name="fecha_fin"
              type="date"
              label="Fecha de fin"
              value={datos.fecha_fin}
              onChange={(evento) => actualizar('fecha_fin', evento.target.value)}
              error={errores.fecha_fin}
            />
          </div>

          {enviando && <LoadingIndicator />}

          <div className="hu04-acciones">
            <Link className="button button--quiet" to={`/projects/${id}/proyecto`}>
              Cancelar
            </Link>
            <Button type="submit" disabled={enviando || cargando}>
              Guardar cambios
            </Button>
          </div>
        </form>
      )}
    </section>
  )
}
