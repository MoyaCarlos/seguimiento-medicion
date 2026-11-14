import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'
import { Link, useParams } from 'react-router'
import { Button } from '../../compartido/componentes/Button'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { ApiError } from '../../compartido/api/client'
import { AvisoSinPermiso } from '../../compartido/permisos/AvisoSinPermiso'
import { NotaGestion } from '../../compartido/permisos/NotaGestion'
import { useEsScrumMaster } from '../../compartido/permisos/useEsScrumMaster'
import { ETIQUETAS_ROL, asignarIntegrante, listarIntegrantes, listarUsuarios } from './api'
import type { Integrante, UsuarioRegistrado } from './api'
import { CampoSelect } from './Campos'
import { mensajeDeError } from './errores'

export function Configuracion() {
  const { id = '' } = useParams()
  const [integrantes, setIntegrantes] = useState<Integrante[]>([])
  const [cargando, setCargando] = useState(true)
  const [errorLista, setErrorLista] = useState('')
  const [usuarios, setUsuarios] = useState<UsuarioRegistrado[]>([])
  const [cargandoUsuarios, setCargandoUsuarios] = useState(true)
  const [errorUsuarios, setErrorUsuarios] = useState('')
  const [nombre, setNombre] = useState('')
  const [errorNombre, setErrorNombre] = useState('')
  const [errorGeneral, setErrorGeneral] = useState('')
  const [enviando, setEnviando] = useState(false)
  const [sinPermiso, setSinPermiso] = useState('')
  const puedeGestionar = useEsScrumMaster(Number(id))

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

  useEffect(() => {
    let vigente = true

    listarUsuarios()
      .then((datos) => {
        if (vigente) setUsuarios(datos)
      })
      .catch((causa: unknown) => {
        if (vigente) setErrorUsuarios(mensajeDeError(causa))
      })
      .finally(() => {
        if (vigente) setCargandoUsuarios(false)
      })

    return () => {
      vigente = false
    }
  }, [])

  const yaIntegran = new Set(integrantes.map((integrante) => integrante.nombre))
  const disponibles = usuarios.filter((usuario) => !yaIntegran.has(usuario.nombre))
  const sinUsuarios = !cargandoUsuarios && !errorUsuarios && disponibles.length === 0

  async function agregar(evento: FormEvent<HTMLFormElement>) {
    evento.preventDefault()
    setErrorGeneral('')
    setSinPermiso('')

    if (!nombre) {
      setErrorNombre('Elegí un usuario de la lista')
      return
    }

    setEnviando(true)
    try {
      const integrante = await asignarIntegrante(id, { nombre })
      setIntegrantes((actuales) => [...actuales, integrante])
      setNombre('')
    } catch (error) {
      if (error instanceof ApiError && error.tipo === 'sin_permiso') {
        setSinPermiso(error.message)
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
          <h2 id="configuracion-title">Configurar equipo</h2>
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
            {sinPermiso && <AvisoSinPermiso mensaje={sinPermiso} />}

            {cargandoUsuarios && <LoadingIndicator />}
            {errorUsuarios && <ErrorMessage message={errorUsuarios} />}
            {sinUsuarios && (
              <p className="hu04-intro">No hay usuarios registrados para agregar</p>
            )}
            {!cargandoUsuarios && !errorUsuarios && !sinUsuarios && (
              <CampoSelect
                id="usuario"
                label="Usuario"
                placeholder="Elegí un usuario"
                value={nombre}
                opciones={disponibles.map((usuario) => ({
                  valor: usuario.nombre,
                  etiqueta: usuario.nombre,
                }))}
                onChange={(valor) => {
                  setNombre(valor)
                  setErrorNombre('')
                  setErrorGeneral('')
                }}
                error={errorNombre}
              />
            )}

            {enviando && <LoadingIndicator />}

            <Button
              type="submit"
              disabled={enviando || !puedeGestionar}
              aria-describedby={puedeGestionar ? undefined : 'configuracion-nota-agregar'}
            >
              Asignar integrante
            </Button>
            {!puedeGestionar && (
              <NotaGestion
                id="configuracion-nota-agregar"
                texto="Solo el Scrum Master puede agregar integrantes"
              />
            )}
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
