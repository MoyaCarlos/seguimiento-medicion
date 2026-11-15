import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'
import { Link, useParams } from 'react-router'
import { AvatarProyecto } from '../../compartido/componentes/AvatarProyecto'
import { Button } from '../../compartido/componentes/Button'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { Icono } from '../../compartido/componentes/Icono'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { ApiError } from '../../compartido/api/client'
import { AvisoSinPermiso } from '../../compartido/permisos/AvisoSinPermiso'
import { AvisoProyectoCerrado } from '../../compartido/permisos/AvisoProyectoCerrado'
import { NotaCerrado } from '../../compartido/permisos/NotaCerrado'
import { NotaGestion } from '../../compartido/permisos/NotaGestion'
import { usePermisosProyecto } from '../../compartido/permisos/usePermisosProyecto'
import { ETIQUETAS_ROL, asignarIntegrante, listarIntegrantes, listarUsuarios, quitarIntegrante } from './api'
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
  const [confirmandoQuitarId, setConfirmandoQuitarId] = useState<number | null>(null)
  const [quitandoId, setQuitandoId] = useState<number | null>(null)
  const [errorQuitar, setErrorQuitar] = useState('')
  const [sinPermisoQuitar, setSinPermisoQuitar] = useState('')
  const { cerrado, fechaCierre, puedeGestionar } = usePermisosProyecto(Number(id))

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

    if (!puedeGestionar) return

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

  async function quitar(integrante: Integrante) {
    if (!puedeGestionar) return
    setErrorQuitar('')
    setSinPermisoQuitar('')
    setQuitandoId(integrante.id)
    try {
      await quitarIntegrante(id, integrante.id)
      setIntegrantes((actuales) => actuales.filter((actual) => actual.id !== integrante.id))
      setConfirmandoQuitarId(null)
    } catch (error) {
      setConfirmandoQuitarId(null)
      if (error instanceof ApiError && error.tipo === 'sin_permiso') {
        setSinPermisoQuitar(error.message)
      } else {
        setErrorQuitar(mensajeDeError(error))
      }
    } finally {
      setQuitandoId(null)
    }
  }

  return (
    <section className="panel-editar" aria-labelledby="configuracion-title">
      <header className="page-header">
        <div>
          <p className="eyebrow">Proyecto #{id}</p>
          <h1 id="configuracion-title">Configurar equipo</h1>
        </div>
        <Link className="button button--ghost" to={`/projects/${id}/proyecto`}>
          <Icono nombre="atras" />
          Volver al panel
        </Link>
      </header>

      <AvisoProyectoCerrado fechaCierre={fechaCierre} />

      <div className="equipo-rejilla">
        <article className="panel-tarjeta">
          <h2 className="equipo-titulo">Asignar integrante</h2>
          <form onSubmit={agregar} noValidate aria-label="Asignar integrante">
            {errorGeneral && <ErrorMessage message={errorGeneral} />}
            {sinPermiso && <AvisoSinPermiso mensaje={sinPermiso} />}

            {cargandoUsuarios && <LoadingIndicator />}
            {errorUsuarios && <ErrorMessage message={errorUsuarios} />}
            {sinUsuarios && (
              <p className="panel-ayuda">No hay usuarios registrados para agregar</p>
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
            {cerrado && <NotaCerrado id="configuracion-nota-agregar" />}
            {!cerrado && !puedeGestionar && (
              <NotaGestion
                id="configuracion-nota-agregar"
                texto="Solo el Scrum Master puede agregar integrantes"
              />
            )}
          </form>
        </article>

        <article className="panel-tarjeta">
          <h2 className="equipo-titulo">Equipo del proyecto</h2>
          {cargando && <LoadingIndicator />}
          {errorLista && <ErrorMessage message={errorLista} />}
          {errorQuitar && <ErrorMessage message={errorQuitar} />}
          {sinPermisoQuitar && <AvisoSinPermiso mensaje={sinPermisoQuitar} />}
          {cerrado && <NotaCerrado id="configuracion-nota-quitar" />}
          {!cerrado && !puedeGestionar && (
            <NotaGestion
              id="configuracion-nota-quitar"
              texto="Solo el Scrum Master puede quitar integrantes"
            />
          )}
          {!cargando && !errorLista && integrantes.length === 0 && (
            <p className="panel-ayuda">Todavía no hay integrantes asignados.</p>
          )}
          {integrantes.length > 0 && (
            <ul className="equipo-lista">
              {integrantes.map((integrante, indice) => (
                <li key={integrante.id}>
                  <AvatarProyecto id={indice + 1} nombre={integrante.nombre} />
                  <span className="equipo-persona">
                    {integrante.nombre}
                    <small>
                      {ETIQUETAS_ROL[integrante.rol]}
                      {integrante.rol === 'scrum_master' && <span> · fijo</span>}
                    </small>
                  </span>
                  {integrante.rol === 'product_builder' && (
                    <div className="equipo-acciones">
                      {confirmandoQuitarId === integrante.id ? (
                        <div className="hu04-confirmacion">
                          <p>¿Quitar a {integrante.nombre} del equipo?</p>
                          <Button
                            type="button"
                            variant="quiet"
                            className="equipo-quitar"
                            onClick={() => setConfirmandoQuitarId(null)}
                            disabled={quitandoId === integrante.id}
                          >
                            Cancelar
                          </Button>
                          <Button
                            type="button"
                            className="equipo-quitar"
                            onClick={() => quitar(integrante)}
                            disabled={quitandoId === integrante.id}
                          >
                            Sí, quitar
                          </Button>
                        </div>
                      ) : (
                        <Button
                          type="button"
                          variant="secondary"
                          className="equipo-quitar"
                          aria-label={`Quitar a ${integrante.nombre}`}
                          onClick={() => {
                            setErrorQuitar('')
                            setSinPermisoQuitar('')
                            setConfirmandoQuitarId(integrante.id)
                          }}
                          disabled={!puedeGestionar}
                          aria-describedby={
                            puedeGestionar ? undefined : 'configuracion-nota-quitar'
                          }
                        >
                          Quitar
                        </Button>
                      )}
                    </div>
                  )}
                </li>
              ))}
            </ul>
          )}
        </article>
      </div>
    </section>
  )
}
