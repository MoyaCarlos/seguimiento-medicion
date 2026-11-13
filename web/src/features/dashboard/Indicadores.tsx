import type { ReactNode } from 'react'
import { Link } from 'react-router'
import { formatoNumero, textoDeDesviacion } from '../../compartido/formato/numeros'
import { proporcion, textoDesviacion, type Metricas } from './metricas'

type IndicadoresProps = {
  metricas: Metricas
  horasPorStoryPoint: number | null
  proyectoId: string
}

const conUnidad = (valor: number | null, unidad: string) =>
  valor === null ? formatoNumero(valor) : `${formatoNumero(valor)}${unidad}`

export function Indicadores({ metricas, horasPorStoryPoint, proyectoId }: IndicadoresProps) {
  const aclaracion = textoDesviacion(metricas.desviacion_horas)

  return (
    <div className="dashboard-indicadores">
      <Indicador titulo="Velocidad">
        <p className="dashboard-indicador__valor">{formatoNumero(metricas.velocidad)}</p>
        <p className="dashboard-indicador__nota">Story Points por Sprint, en promedio</p>
      </Indicador>
      <Indicador titulo="Historias completadas">
        <p className="dashboard-indicador__valor">{conUnidad(metricas.porcentaje_historias_completadas, ' %')}</p>
        <Barra nombre="Historias completadas" valor={metricas.porcentaje_historias_completadas} maximo={100} />
      </Indicador>
      <Indicador titulo="Horas">
        <Pares
          pares={[
            ['Horas estimadas', conUnidad(metricas.horas_estimadas, ' h')],
            ['Horas reales', conUnidad(metricas.horas_reales, ' h')],
          ]}
        />
        {horasPorStoryPoint === null && (
          <p className="dashboard-indicador__nota">
            Las horas por Story Point no están definidas. Podés cargarlas al{' '}
            <Link to={`/projects/${proyectoId}/proyecto/edit`}>editar el proyecto</Link>.
          </p>
        )}
      </Indicador>
      <Indicador titulo="Desviación">
        <p className="dashboard-indicador__valor">
          {textoDeDesviacion(metricas.desviacion_horas, metricas.desviacion_porcentaje)}
        </p>
        {aclaracion && <p className="dashboard-indicador__nota">{aclaracion}</p>}
      </Indicador>
      <Indicador titulo="Defectos">
        <Pares
          pares={[
            ['Detectados', String(metricas.defectos_detectados)],
            ['Resueltos', String(metricas.defectos_resueltos)],
          ]}
        />
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

function Pares({ pares }: { pares: [string, string][] }) {
  return (
    <dl className="dashboard-pares">
      {pares.map(([rotulo, valor]) => (
        <div key={rotulo}>
          <dt>{rotulo}</dt>
          <dd>{valor}</dd>
        </div>
      ))}
    </dl>
  )
}

function Barra({ nombre, valor, maximo }: { nombre: string; valor: number; maximo: number }) {
  return (
    <div
      className="dashboard-barra"
      role="meter"
      aria-label={nombre}
      aria-valuemin={0}
      aria-valuemax={maximo > 0 ? maximo : 1}
      aria-valuenow={valor}
    >
      <span style={{ width: `${proporcion(valor, maximo)}%` }} />
    </div>
  )
}
