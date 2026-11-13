import { useCallback } from 'react'
import { Link, useParams } from 'react-router'
import { BotonExportar } from '../../compartido/componentes/BotonExportar'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { descargarReporteSprint, obtenerDetalleSprint } from './api'
import { formatoNumero as numero, textoDeDesviacion } from '../../compartido/formato/numeros'
import { periodoCorto } from '../../compartido/formato/fechas'
import { conUnidad } from './formato'
import { Chip, Metrica } from './Piezas'
import type { DefectoDeSprint, DetalleSprint, HistoriaComprometida } from './tipos'
import { useConsulta } from './useConsulta'

export function Detalle() {
  const params = useParams()
  const proyectoId = params.id ?? ''
  const sprintId = params.sprintId ?? ''
  const consultar = useCallback(() => obtenerDetalleSprint(proyectoId, sprintId), [proyectoId, sprintId])
  const { datos: sprint, error } = useConsulta(consultar)

  return (
    <div className="historial-pantalla historial-detalle">
      <Link className="historial-volver" to={`/projects/${proyectoId}/historial`}>
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
      <header className="historial-encabezado historial-encabezado--detalle">
        <div>
          <p className="historial-eyebrow">Historial de Sprints · Proyecto #{proyectoId}</p>
          <h1>Sprint {sprint.numero}</h1>
          <p className="historial-bajada">
            {sprint.sprint_goal ? `${sprint.sprint_goal} · ` : ''}
            {periodoCorto(sprint.fecha_inicio, sprint.fecha_fin)}
          </p>
        </div>
        <BotonExportar descargar={() => descargarReporteSprint(String(sprint.id))}>
          Exportar a PDF
        </BotonExportar>
      </header>

      <section className="historial-tarjeta historial-seccion" aria-labelledby="historial-metricas">
        <h2 id="historial-metricas">Métricas del Sprint</h2>
        <dl className="historial-metricas">
          <Metrica nombre="Velocidad" valor={numero(metricas.velocidad)} />
          <Metrica nombre="Story Points planificados" valor={numero(metricas.story_points_planificados)} />
          <Metrica nombre="Story Points completados" valor={numero(metricas.story_points_completados)} />
          <Metrica nombre="Historias completadas" valor={conUnidad(metricas.porcentaje_historias_completadas, ' %')} />
          <Metrica nombre="Horas estimadas" valor={conUnidad(metricas.horas_estimadas, ' h')} />
          <Metrica nombre="Horas reales" valor={conUnidad(metricas.horas_reales, ' h')} />
          <Metrica nombre="Desviación" valor={textoDeDesviacion(metricas.desviacion_horas, metricas.desviacion_porcentaje)} />
          <Metrica nombre="Defectos detectados" valor={numero(metricas.defectos_detectados)} />
          <Metrica nombre="Defectos resueltos" valor={numero(metricas.defectos_resueltos)} />
        </dl>
      </section>

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
        <p className="historial-total">
          Esfuerzo total del Sprint: <strong>{conUnidad(metricas.horas_reales, ' h')}</strong>
        </p>
      </section>

      <section className="historial-tarjeta historial-seccion" aria-labelledby="historial-defectos">
        <h2 id="historial-defectos">Defectos</h2>
        {sprint.defectos.length === 0 ? (
          <p className="historial-vacio-corto">Sin defectos en este Sprint.</p>
        ) : (
          <ul className="historial-items historial-items--cajas">
            {sprint.defectos.map((defecto) => (
              <ItemDefecto key={defecto.id} defecto={defecto} />
            ))}
          </ul>
        )}
      </section>
    </>
  )
}

function ItemHistoria({ historia }: { historia: HistoriaComprometida }) {
  return (
    <li className="historial-item">
      <p className="historial-item__titulo">{historia.titulo}</p>
      <ul className="historial-chips">
        {historia.completada ? <Chip tono="ok">Completada</Chip> : <Chip tono="apagado">No completada</Chip>}
        {historia.estimacion_sp === null ? (
          <Chip tono="ninguno">Sin estimar</Chip>
        ) : (
          <Chip tono="fuerte">{historia.estimacion_sp} Story Points</Chip>
        )}
        {historia.horas > 0 && <Chip tono="fuerte">{conUnidad(historia.horas, ' h')}</Chip>}
      </ul>
    </li>
  )
}

const tonoSeveridad = { Alta: 'alerta', Media: 'media', Baja: 'baja' } as const

function ItemDefecto({ defecto }: { defecto: DefectoDeSprint }) {
  return (
    <li className="historial-item historial-item--caja">
      <p className="historial-item__descripcion">{defecto.descripcion}</p>
      <p className="historial-item__nota">Historia: {defecto.historia}</p>
      <ul className="historial-chips">
        <Chip tono={tonoSeveridad[defecto.severidad]}>Severidad {defecto.severidad}</Chip>
        <Chip tono={defecto.estado === 'Resuelto' ? 'ok' : 'apagado'}>{defecto.estado}</Chip>
        {defecto.detectado_en_sprint && <Chip tono="fuerte">Detectado en este Sprint</Chip>}
        {defecto.resuelto_en_sprint && <Chip tono="fuerte">Resuelto en este Sprint</Chip>}
      </ul>
    </li>
  )
}
