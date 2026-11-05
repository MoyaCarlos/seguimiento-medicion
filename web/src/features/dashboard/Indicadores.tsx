import type { ReactNode } from 'react'
import { Link } from 'react-router'
import { formatoNumero, proporcion, textoDesviacion, type Metricas } from './metricas'

type IndicadoresProps = {
  metricas: Metricas
  horasPorStoryPoint: number | null
  proyectoId: string
}

const conUnidad = (valor: number | null, unidad: string) =>
  valor === null ? formatoNumero(valor) : `${formatoNumero(valor)}${unidad}`

export function Indicadores({ metricas, horasPorStoryPoint, proyectoId }: IndicadoresProps) {
  const maximoHoras = Math.max(metricas.horas_estimadas ?? 0, metricas.horas_reales)
  const maximoDefectos = Math.max(metricas.defectos_detectados, metricas.defectos_resueltos)
  const aclaracion = textoDesviacion(metricas.desviacion_horas)

  return (
    <div className="dashboard-indicadores">
      <Indicador titulo="Velocidad">
        <p className="dashboard-indicador__valor">{formatoNumero(metricas.velocidad)}</p>
        <p className="dashboard-indicador__nota">Story Points por Sprint, en promedio</p>
      </Indicador>

      <Indicador titulo="Historias completadas">
        <p className="dashboard-indicador__valor">{conUnidad(metricas.porcentaje_historias_completadas, '%')}</p>
        <Barra nombre="Historias completadas" valor={metricas.porcentaje_historias_completadas} maximo={100} />
      </Indicador>

      <Indicador titulo="Horas">
        <Comparacion nombre="Horas estimadas" valor={metricas.horas_estimadas} maximo={maximoHoras} unidad=" h" />
        <Comparacion nombre="Horas reales" valor={metricas.horas_reales} maximo={maximoHoras} unidad=" h" />
        {horasPorStoryPoint === null && (
          <p className="dashboard-indicador__nota">
            Las horas por Story Point no están definidas. Podés cargarlas al{' '}
            <Link to={`/projects/${proyectoId}/proyecto/edit`}>editar el proyecto</Link>.
          </p>
        )}
      </Indicador>

      <Indicador titulo="Desviación">
        <p className="dashboard-indicador__valor">{conUnidad(metricas.desviacion_horas, ' h')}</p>
        <p className="dashboard-indicador__secundario">{conUnidad(metricas.desviacion_porcentaje, '%')}</p>
        {aclaracion && <p className="dashboard-indicador__nota">{aclaracion}</p>}
      </Indicador>

      <Indicador titulo="Defectos">
        <Comparacion nombre="Defectos detectados" valor={metricas.defectos_detectados} maximo={maximoDefectos} />
        <Comparacion nombre="Defectos resueltos" valor={metricas.defectos_resueltos} maximo={maximoDefectos} />
      </Indicador>
    </div>
  )
}

function Indicador({ titulo, children }: { titulo: string; children: ReactNode }) {
  return (
    <section className="dashboard-tarjeta dashboard-indicador" aria-label={titulo}>
      <h2>{titulo}</h2>
      {children}
    </section>
  )
}

type ComparacionProps = {
  nombre: string
  valor: number | null
  maximo: number
  unidad?: string
}

function Comparacion({ nombre, valor, maximo, unidad = '' }: ComparacionProps) {
  return (
    <div className="dashboard-comparacion">
      <div className="dashboard-comparacion__rotulo">
        <span>{nombre}</span>
        <strong>{conUnidad(valor, unidad)}</strong>
      </div>
      {valor !== null && <Barra nombre={nombre} valor={valor} maximo={maximo} />}
    </div>
  )
}

function Barra({ nombre, valor, maximo }: { nombre: string; valor: number; maximo: number }) {
  return (
    <div
      className="dashboard-barra"
      role="meter"
      aria-label={nombre}
      aria-valuemin={0}
      aria-valuemax={maximo}
      aria-valuenow={valor}
    >
      <span style={{ width: `${proporcion(valor, maximo)}%` }} />
    </div>
  )
}
