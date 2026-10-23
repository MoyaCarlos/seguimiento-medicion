import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'
import { Link, useParams } from 'react-router'
import { Button } from '../../compartido/componentes/Button'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { Field } from '../../compartido/componentes/Field'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { ETIQUETAS_ROL, asignarIntegrante, listarIntegrantes } from './api'
import type { Integrante, Rol } from './api'
import { CampoSelect } from './Campos'
import { campoDeError, mensajeDeError } from './errores'
import { validarIntegrante } from './validacion'
import type { ErroresIntegrante } from './validacion'

const OPCIONES_ROL = [
  { valor: 'scrum_master', etiqueta: ETIQUETAS_ROL.scrum_master },
  { valor: 'product_builder', etiqueta: ETIQUETAS_ROL.product_builder },
]

export function Configuracion() {
  const { id = '' } = useParams()
  const [integrantes, setIntegrantes] = useState<Integrante[]>([])
  const [cargando, setCargando] = useState(true)
  const [errorLista, setErrorLista] = useState('')
  const [nombre, setNombre] = useState('')
  const [rol, setRol] = useState('')
  const [errores, setErrores] = useState<ErroresIntegrante>({})
  const [errorGeneral, setErrorGeneral] = useState('')
  const [enviando, setEnviando] = useState(false)

  useEffect(() => {
    let vigente = true

    listarIntegrantes(id)
      .then((datos) => {
        if (vigente) setIntegrantes(datos)
      })
      .catch((causa: unknown) => {
        if (vigente) setErrorLista(mensajeDeError(causa))
      })
      .finally(() => {
        if (vigente) setCargando(false)
      })

    return () => {
      vigente = false
    }
  }, [id])

  async function agregar(evento: FormEvent<HTMLFormElement>) {
    evento.preventDefault()
    setErrorGeneral('')

    const validacion = validarIntegrante({ nombre, rol })
    setErrores(validacion)
    if (Object.keys(validacion).length > 0) return

    setEnviando(true)
    try {
      const integrante = await asignarIntegrante(id, {
        nombre: nombre.trim(),
        rol: rol as Rol,
      })
      setIntegrantes((actuales) => [...actuales, integrante])
      setNombre('')
      setRol('')
    } catch (error) {
      const campo = campoDeError(error)
      if (campo === 'nombre' || campo === 'rol') {
        setErrores({ [campo]: mensajeDeError(error) })
      } else {
        setErrorGeneral(mensajeDeError(error))
      }
    } finally {
      setEnviando(false)
    }
  }

  return (
    <section className="hu04-panel" aria-labelledby="configuracion-title">
      <div className="hu04-encabezado">
        <div>
          <p className="eyebrow">Configuración</p>
          <h1 id="configuracion-title">Configurar equipo</h1>
          <p className="hu04-referencia">
            Proyecto <strong>{id}</strong>
          </p>
        </div>
        <Link className="button button--quiet" to={`/projects/${id}/proyecto`}>
          Volver al panel
        </Link>
      </div>

      <div className="hu04-columnas">
        <article className="hu04-tarjeta">
          <h2>Asignar integrante</h2>
          <form onSubmit={agregar} noValidate aria-label="Asignar integrante">
            {errorGeneral && <ErrorMessage message={errorGeneral} />}

            <Field
              id="integrante"
              name="nombre"
              label="Nombre del integrante"
              value={nombre}
              onChange={(evento) => {
                setNombre(evento.target.value)
                setErrores((actuales) => ({ ...actuales, nombre: undefined }))
                setErrorGeneral('')
              }}
              error={errores.nombre}
            />

            <CampoSelect
              id="rol"
              label="Rol"
              value={rol}
              opciones={OPCIONES_ROL}
              onChange={(valor) => {
                setRol(valor)
                setErrores((actuales) => ({ ...actuales, rol: undefined }))
                setErrorGeneral('')
              }}
              error={errores.rol}
            />

            {enviando && <LoadingIndicator />}

            <Button type="submit" disabled={enviando}>
              Asignar integrante
            </Button>
          </form>
        </article>

        <article className="hu04-tarjeta">
          <h2>Equipo del proyecto</h2>
          {cargando && <LoadingIndicator />}
          {errorLista && <ErrorMessage message={errorLista} />}
          {!cargando && !errorLista && integrantes.length === 0 && (
            <p className="hu04-intro">Todavía no hay integrantes asignados.</p>
          )}
          {integrantes.length > 0 && (
            <ul className="hu04-equipo">
              {integrantes.map((integrante) => (
                <li key={integrante.id}>
                  <span>{integrante.nombre}</span>
                  <strong>{ETIQUETAS_ROL[integrante.rol]}</strong>
                </li>
              ))}
            </ul>
          )}
        </article>
      </div>
    </section>
  )
}
