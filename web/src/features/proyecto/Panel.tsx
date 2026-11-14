import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router'
import { ApiError } from '../../compartido/api/client'
import { NoIntegrante } from '../../compartido/permisos/NoIntegrante'
import { Button } from '../../compartido/componentes/Button'
import { BotonExportar } from '../../compartido/componentes/BotonExportar'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { Icono } from '../../compartido/componentes/Icono'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { fechaCorta } from '../../compartido/formato/fechas'
import { NotaGestion } from '../../compartido/permisos/NotaGestion'
import { AvatarProyecto } from '../../compartido/componentes/AvatarProyecto'
import { ETIQUETAS_ROL, descargarReporteProyecto, listarIntegrantes, obtenerEstadoProyecto, obtenerProyecto } from './api'
import type { EstadoProyecto, Integrante, Proyecto } from './api'
import { mensajeDeError } from './errores'

const ESTADOS_CON_INSIGNIA: string[] = ['Planificado', 'En curso', 'Finalizado']
const VARIANTE_INSIGNIA: Record<string, string> = {
  Planificado: ' badge--planificado',
  Finalizado: ' badge--finalizado',
}

export function Panel() {
  const { id = '' } = useParams()
  const [proyecto, setProyecto] = useState<Proyecto>()
  const [error, setError] = useState('')
  const [cargando, setCargando] = useState(true)
  const [estado, setEstado] = useState<EstadoProyecto>()
  const [errorEstado, setErrorEstado] = useState('')
  const [cargandoEstado, setCargandoEstado] = useState(true)
  const [noIntegrante, setNoIntegrante] = useState(false)
  const [integrantes, setIntegrantes] = useState<Integrante[]>([])
  const [errorEquipo, setErrorEquipo] = useState('')
  const [cargandoEquipo, setCargandoEquipo] = useState(true)

  useEffect(() => {
    let vigente = true

    obtenerProyecto(id)
      .then((datos) => {
        if (vigente) setProyecto(datos)
      })
      .catch((causa: unknown) => {
        if (!vigente) return
        if (causa instanceof ApiError && causa.tipo === 'no_integrante') setNoIntegrante(true)
        else setError(mensajeDeError(causa))
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

    obtenerEstadoProyecto(id)
      .then((respuesta) => {
        if (vigente) setEstado(respuesta.estado)
      })
      .catch((causa: unknown) => {
        if (vigente) setErrorEstado(mensajeDeError(causa))
      })
      .finally(() => {
        if (vigente) setCargandoEstado(false)
      })

    return () => {
      vigente = false
    }
  }, [id])

  useEffect(() => {
    let vigente = true

    listarIntegrantes(id)
      .then((datos) => {
        if (vigente) setIntegrantes(datos)
      })
      .catch((causa: unknown) => {
        if (vigente) setErrorEquipo(mensajeDeError(causa))
      })
      .finally(() => {
        if (vigente) setCargandoEquipo(false)
      })

    return () => {
      vigente = false
    }
  }, [id])

  if (noIntegrante) return <NoIntegrante />

  const soloLectura = proyecto?.rol === 'product_builder'

  return (
    <section className="hu04-panel" aria-labelledby="panel-title">
      <header className="page-header">
        <div>
          <p className="eyebrow">Panel principal · #{id}</p>
          {proyecto && (
            <div className="title-row">
              <h1 id="panel-title">{proyecto.nombre}</h1>
              {estado && ESTADOS_CON_INSIGNIA.includes(estado) && (
                <span className={`badge${VARIANTE_INSIGNIA[estado] ?? ''}`}>{estado}</span>
              )}
            </div>
          )}
        </div>
        <div className="hu04-acciones">
          <BotonExportar descargar={() => descargarReporteProyecto(id)}>
            Exportar reporte del proyecto
          </BotonExportar>
          {soloLectura ? (
            <>
              <Button variant="secondary" disabled aria-describedby="nota-rol">
                <Icono nombre="equipo" />
                Configurar equipo
              </Button>
              <Button disabled aria-describedby="nota-rol">
                <Icono nombre="editar" />
                Editar proyecto
              </Button>
            </>
          ) : (
            <>
              <Link className="button button--secondary" to="settings">
                <Icono nombre="equipo" />
                Configurar equipo
              </Link>
              <Link className="button button--primary" to="edit">
                <Icono nombre="editar" />
                Editar proyecto
              </Link>
            </>
          )}
        </div>
        {soloLectura && (
          <NotaGestion
            id="nota-rol"
            texto="Solo el Scrum Master puede editar el proyecto y configurar el equipo."
          />
        )}
      </header>

      {cargando && <LoadingIndicator />}
      {error && <ErrorMessage message={error} />}

      {proyecto && (
        <div className="panel-rejilla">
        <article className="panel-tarjeta">
          <h2 className="panel-tarjeta__titulo">Descripción</h2>
          <p className="panel-prosa">{proyecto.descripcion || 'Sin descripción'}</p>
          <hr className="panel-divisor" />
          <section aria-labelledby="cronograma-titulo">
            <h2 id="cronograma-titulo" className="panel-tarjeta__titulo">
              Cronograma
            </h2>
            <dl className="cronograma__fechas">
              <div>
                <dt>Inicio</dt>
                <dd>
                  <strong>{fechaCorta(proyecto.fecha_inicio)}</strong>
                </dd>
              </div>
              <div>
                <dt>Fin</dt>
                <dd>
                  <strong>{fechaCorta(proyecto.fecha_fin)}</strong>
                </dd>
              </div>
              <div>
                <dt>Estado</dt>
                <dd aria-live="polite">
                  <strong>{cargandoEstado ? 'Consultando…' : (estado ?? 'Sin datos')}</strong>
                </dd>
              </div>
            </dl>
            {errorEstado && <ErrorMessage message={errorEstado} />}
          </section>
          <p className="panel-tarjeta__pie">
            <Link className="enlace" to={`/projects/${id}/status`}>
              Ver estado del proyecto
              <Icono nombre="flecha" />
            </Link>
          </p>
        </article>
        <aside className="panel-tarjeta">
          <h2 className="panel-tarjeta__titulo">Equipo</h2>
          {cargandoEquipo && <LoadingIndicator />}
          {errorEquipo && <ErrorMessage message={errorEquipo} />}
          {!cargandoEquipo && !errorEquipo && integrantes.length === 0 && (
            <p className="panel-ayuda">Todavía no hay integrantes asignados.</p>
          )}
          {integrantes.length > 0 && (
            <ul className="panel-equipo">
              {integrantes.map((integrante, indice) => (
                <li key={integrante.id}>
                  <AvatarProyecto id={indice + 1} nombre={integrante.nombre} />
                  <span className="panel-equipo__nombre">
                    {integrante.nombre}
                    <small>{ETIQUETAS_ROL[integrante.rol]}</small>
                  </span>
                </li>
              ))}
            </ul>
          )}
          <p className="panel-tarjeta__pie">
            <Link className="enlace" to="settings">
              {soloLectura ? 'Ver equipo' : 'Configurar equipo'}
              <Icono nombre="flecha" />
            </Link>
          </p>
        </aside>
        </div>
      )}
    </section>
  )
}
