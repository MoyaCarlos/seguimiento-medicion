import { useEffect, useState } from 'react'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import {
  listarDefectos,
  listarHistoriasDelProyecto,
  listarSprints,
  type Defecto,
  type HistoriaDeProyecto,
  type SprintResumen,
} from './api'
import { FormularioDefecto } from './FormularioDefecto'
import { ListadoDefectos } from './ListadoDefectos'
import { mensajeDeError } from './errores'
import './defectos.css'

type Props = {
  proyectoId: string
}

export function PantallaDefectos({ proyectoId }: Props) {
  const [defectos, setDefectos] = useState<Defecto[]>([])
  const [cargandoDefectos, setCargandoDefectos] = useState(true)
  const [errorDefectos, setErrorDefectos] = useState('')
  const [sprints, setSprints] = useState<SprintResumen[] | null>(null)
  const [historias, setHistorias] = useState<HistoriaDeProyecto[] | null>(null)
  const [cargandoDatos, setCargandoDatos] = useState(true)
  const [errorDatos, setErrorDatos] = useState('')
  const [exito, setExito] = useState('')

  function recargarDefectos() {
    setCargandoDefectos(true)
    setErrorDefectos('')
    listarDefectos(proyectoId)
      .then(setDefectos)
      .catch((causa: unknown) => setErrorDefectos(mensajeDeError(causa)))
      .finally(() => setCargandoDefectos(false))
  }

  useEffect(() => {
    let vigente = true
    listarDefectos(proyectoId)
      .then((datos) => {
        if (vigente) setDefectos(datos)
      })
      .catch((causa: unknown) => {
        if (vigente) setErrorDefectos(mensajeDeError(causa))
      })
      .finally(() => {
        if (vigente) setCargandoDefectos(false)
      })

    return () => {
      vigente = false
    }
  }, [proyectoId])

  useEffect(() => {
    let vigente = true
    ;(async () => {
      const listaSprints = await listarSprints(proyectoId)
      const listaHistorias = await listarHistoriasDelProyecto(proyectoId, listaSprints)
      if (vigente) {
        setSprints(listaSprints)
        setHistorias(listaHistorias)
      }
    })()
      .catch((causa: unknown) => {
        if (vigente) setErrorDatos(mensajeDeError(causa))
      })
      .finally(() => {
        if (vigente) setCargandoDatos(false)
      })

    return () => {
      vigente = false
    }
  }, [proyectoId])

  const activo = sprints?.find((sprint) => sprint.estado === 'Activo')
  const sinSprintActivo = sprints !== null && !activo
  const sinHistorias =
    sprints !== null && activo !== undefined && historias !== null && historias.length === 0

  return (
    <section className="defectos" aria-labelledby="defectos-title">
      <div className="page-header">
        <div>
          <p className="eyebrow">Calidad</p>
          <h1 id="defectos-title">Defectos</h1>
          <p className="project-context">Proyecto actual: <strong>{proyectoId}</strong></p>
        </div>
      </div>

      {exito && (
        <p className="defectos-exito" role="status">
          {exito}
        </p>
      )}
      {cargandoDatos && <LoadingIndicator />}
      {errorDatos && <ErrorMessage message={errorDatos} />}
      {sinSprintActivo && (
        <p className="message defectos-aviso" role="status">
          No hay un Sprint activo. Para registrar o resolver defectos, iniciá un Sprint.
        </p>
      )}
      {sinHistorias && (
        <p className="message defectos-aviso" role="status">
          No hay historias en el proyecto. Creá una antes de registrar un defecto.
        </p>
      )}

      <div className="defectos-columnas">
        {sprints && activo && historias && historias.length > 0 && (
          <section className="defectos-registro" aria-labelledby="registrar-defecto-title">
            <h2 id="registrar-defecto-title">Registrar defecto</h2>
            <FormularioDefecto
              proyectoId={proyectoId}
              historias={historias}
              alRegistrar={() => {
              recargarDefectos()
              setExito('Defecto registrado')
            }}
            />
          </section>
        )}

        <section className="defectos-listado" aria-labelledby="defectos-listado-title">
          <h2 id="defectos-listado-title">Defectos registrados</h2>
          <ListadoDefectos
            defectos={defectos}
            cargando={cargandoDefectos}
            error={errorDefectos}
            haySprintActivo={Boolean(activo)}
            sprints={sprints ?? []}
            alResolver={(mensaje) => {
            recargarDefectos()
            setExito(mensaje)
          }}
          />
        </section>
      </div>
    </section>
  )
}
