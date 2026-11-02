import { useState } from 'react'
import type { FormEvent } from 'react'
import { useNavigate } from 'react-router'
import { Button } from '../../compartido/componentes/Button'
import { BotonDatosDePrueba } from '../../compartido/demo/BotonDatosDePrueba'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { Field } from '../../compartido/componentes/Field'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { aNuevoProyecto, crearProyecto } from './api'
import { CampoDescripcion } from './Campos'
import { campoDeError, mensajeDeError } from './errores'
import { validarCreacion } from './validacion'
import type { DatosProyecto, ErroresProyecto } from './validacion'
import './estilos.css'

const FORMULARIO_INICIAL: DatosProyecto = {
  nombre: '',
  descripcion: '',
  fecha_inicio: '',
  fecha_fin: '',
  creador: '',
}

export function Component() {
  const navegar = useNavigate()
  const [datos, setDatos] = useState(FORMULARIO_INICIAL)
  const [errores, setErrores] = useState<ErroresProyecto>({})
  const [errorGeneral, setErrorGeneral] = useState('')
  const [enviando, setEnviando] = useState(false)

  function actualizar(campo: keyof DatosProyecto, valor: string) {
    setDatos((actuales) => ({ ...actuales, [campo]: valor }))
    setErrores((actuales) => ({ ...actuales, [campo]: undefined }))
    setErrorGeneral('')
  }

  async function enviar(evento: FormEvent<HTMLFormElement>) {
    evento.preventDefault()
    setErrorGeneral('')

    const validacion = validarCreacion(datos)
    setErrores(validacion)
    if (Object.keys(validacion).length > 0) return

    setEnviando(true)
    try {
      const proyecto = await crearProyecto(aNuevoProyecto(datos))
      navegar(`/projects/${proyecto.id}/proyecto`, { replace: true })
    } catch (error) {
      const campo = campoDeError(error)
      if (campo && campo !== 'nombre' && campo in datos) {
        setErrores({ [campo]: mensajeDeError(error) })
      } else {
        setErrorGeneral(mensajeDeError(error))
      }
    } finally {
      setEnviando(false)
    }
  }

  return (
    <section className="hu04-tarjeta hu04-tarjeta--angosta" aria-labelledby="crear-title">
      <p className="eyebrow">Métrica · Scrum Master</p>
      <h1 id="crear-title">Crear proyecto</h1>
      <p className="hu04-intro">
        Definí los datos iniciales del proyecto para comenzar a hacerle seguimiento.
      </p>

      <form onSubmit={enviar} noValidate aria-label="Crear proyecto">
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

        <Field
          id="creador"
          name="creador"
          label="Tu nombre (Scrum Master)"
          value={datos.creador}
          onChange={(evento) => actualizar('creador', evento.target.value)}
          error={errores.creador}
        />

        {enviando && <LoadingIndicator />}

        <Button type="submit" disabled={enviando}>
          Crear proyecto
        </Button>
      </form>
      <BotonDatosDePrueba />
    </section>
  )
}
