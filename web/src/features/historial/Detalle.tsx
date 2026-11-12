import { useCallback } from 'react'
import { Link, useParams } from 'react-router'
import { BotonExportar } from '../../compartido/componentes/BotonExportar'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { descargarReporteSprint, obtenerDetalleSprint } from './api'
import { conUnidad, desviacion, numero, periodo } from './formato'
import { BarraAvance, Chip, Metrica } from './Piezas'
import type { DefectoDeSprint, DetalleSprint, HistoriaComprometida } from './tipos'
import { useConsulta } from './useConsulta'

export function Detalle() {
  const params = useParams()
  const proyectoId = params.id ?? ''
  const sprintId = params.sprintId ?? ''
  const consultar = useCallback(() => obtenerDetalleSprint(proyectoId, sprintId), [proyectoId, sprintId])
  const { datos: sprint, error } = useConsulta(consultar)

  return (
    <div className="historial-pantalla">
      <Link className="button button--quiet historial-volver" to={`/projects/${proyectoId}/historial`}>
        <span aria-hidden="true">←</span> Volver al historial
      </Link>
      {error && <ErrorMessage message={error} />}
      {!sprint && !error && <LoadingIndicator />}
      {sprint && <Contenido sprint={sprint} proyectoId={proyectoId} />}
    </div>
  )
}

function Contenido({ sprint, proyectoId }: { sprint: DetalleSprint; proyectoId: string }) {
  const { metricas } = sprint

  return (
    <>
      <header className="historial-encabezado">
        <div>
          <p className="historial-eyebrow">Historial de Sprints · Proyecto #{proyectoId}</p>
          <h1>Sprint {sprint.numero}</h1>
          <p className="historial-goal historial-goal--detalle">{sprint.sprint_goal}</p>
          <p className="historial-fechas">{periodo(sprint.fecha_inicio, sprint.fecha_fin)}</p>
        </div>
        <BotonExportar descargar={() => descargarReporteSprint(String(sprint.id))}>
          Exportar a PDF
        </BotonExportar>
      </header>

      <section className="historial-tarjeta historial-seccion" aria-labelledby="historial-metricas">
        <h2 id="historial-metricas">Métricas del Sprint</h2>
        <BarraAvance completados={metricas.story_points_completados} planificados={metricas.story_points_planificados} />
        <dl className="historial-metricas">
          <Metrica nombre="Story Points planificados" valor={numero(metricas.story_points_planificados)} />
          <Metrica nombre="Story Points completados" valor={numero(metricas.story_points_completados)} />
          <Metrica nombre="Velocidad" valor={numero(metricas.velocidad)} />
          <Metrica nombre="Historias completadas" valor={conUnidad(metricas.porcentaje_historias_completadas, '%')} />
          <Metrica nombre="Horas estimadas" valor={conUnidad(metricas.horas_estimadas, ' h')} />
          <Metrica nombre="Horas reales" valor={conUnidad(metricas.horas_reales, ' h')} />
          <Metrica nombre="Desviación" valor={desviacion(metricas.desviacion_horas, metricas.desviacion_porcentaje)} />
          <Metrica nombre="Defectos detectados" valor={numero(metricas.defectos_detectados)} />
          <Metrica nombre="Defectos resueltos" valor={numero(metricas.defectos_resueltos)} />
        </dl>
      </section>

      <div className="historial-columnas">
        <section className="historial-tarjeta historial-seccion" aria-labelledby="historial-historias">
          <h2 id="historial-historias">Historias comprometidas</h2>
          {sprint.historias.length === 0 ? (
            <p className="historial-vacio-corto">Sin historias comprometidas.</p>
          ) : (
            <ul className="historial-items">
              {sprint.historias.map((historia) => (
                <ItemHistoria key={historia.id} historia={historia} />
              ))}
            </ul>
          )}
          <dl className="historial-total">
            <Metrica nombre="Esfuerzo total del Sprint" valor={conUnidad(metricas.horas_reales, ' h')} />
          </dl>
        </section>

        <section className="historial-tarjeta historial-seccion" aria-labelledby="historial-defectos">
          <h2 id="historial-defectos">Defectos</h2>
          {sprint.defectos.length === 0 ? (
            <p className="historial-vacio-corto">Sin defectos en este Sprint.</p>
          ) : (
            <ul className="historial-items">
              {sprint.defectos.map((defecto) => (
                <ItemDefecto key={defecto.id} defecto={defecto} />
              ))}
            </ul>
          )}
        </section>
      </div>
    </>
  )
}

function ItemHistoria({ historia }: { historia: HistoriaComprometida }) {
  return (
    <li className="historial-item">
      <div className="historial-item__fila">
        <p className="historial-item__titulo">{historia.titulo}</p>
        <span className="historial-item__horas">{conUnidad(historia.horas, ' h')}</span>
      </div>
      <ul className="historial-chips">
        {historia.completada ? <Chip tono="ok">Completada</Chip> : <Chip tono="apagado">No completada</Chip>}
        {historia.estimacion_sp === null ? (
          <Chip tono="apagado">Sin estimar</Chip>
        ) : (
          <Chip>{historia.estimacion_sp} SP</Chip>
        )}
      </ul>
    </li>
  )
}

const tonoSeveridad = { Alta: 'alerta', Media: 'neutro', Baja: 'apagado' } as const

function ItemDefecto({ defecto }: { defecto: DefectoDeSprint }) {
  return (
    <li className="historial-item">
      <p className="historial-item__titulo">{defecto.descripcion}</p>
      <p className="historial-item__nota">Historia: {defecto.historia}</p>
      <ul className="historial-chips">
        <Chip tono={tonoSeveridad[defecto.severidad]}>{defecto.severidad}</Chip>
        <Chip tono={defecto.estado === 'Resuelto' ? 'ok' : 'apagado'}>{defecto.estado}</Chip>
        {defecto.detectado_en_sprint && <Chip tono="fuerte">Detectado en este Sprint</Chip>}
        {defecto.resuelto_en_sprint && <Chip tono="fuerte">Resuelto en este Sprint</Chip>}
      </ul>
    </li>
  )
}
