import { Bar, BarChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts'
import { datosDelGrafico, type MetricasDeSprint } from './metricas'

const ALTO = 320
const ANCHO_POR_SPRINT = 72
const ANCHO_DEL_EJE_Y = 56
const COLOR_COMPROMETIDOS = 'var(--dashboard-comprometidos)'
const COLOR_COMPLETADOS = 'var(--lilac-dark)'
const RAYADO = 'dashboard-comprometidos-rayado'

export function GraficoVelocidad({ sprints }: { sprints: MetricasDeSprint[] }) {
  const filas = datosDelGrafico(sprints)
  const anchoMinimo = ANCHO_DEL_EJE_Y + filas.length * ANCHO_POR_SPRINT

  return (
    <figure className="dashboard-tarjeta dashboard-grafico" aria-labelledby="dashboard-grafico-titulo">
      <figcaption id="dashboard-grafico-titulo">
        <h2>Story Points comprometidos y completados por Sprint</h2>
      </figcaption>
      <ul className="dashboard-leyenda" aria-label="Series del gráfico">
        <li>
          <span className="dashboard-leyenda__marca dashboard-leyenda__marca--comprometidos" aria-hidden="true" />
          Comprometidos
        </li>
        <li>
          <span className="dashboard-leyenda__marca dashboard-leyenda__marca--completados" aria-hidden="true" />
          Completados
        </li>
      </ul>
      <div
        className="dashboard-grafico__desplazable"
        role="region"
        aria-label="Gráfico por Sprint, desplazable"
        tabIndex={0}
      >
        <div style={{ minWidth: anchoMinimo }}>
          <ResponsiveContainer width="100%" height={ALTO} initialDimension={{ width: anchoMinimo, height: ALTO }}>
            <BarChart data={filas} barGap={2} margin={{ top: 8, right: 8, bottom: 8, left: 0 }}>
              <defs>
                <pattern id={RAYADO} width="6" height="6" patternUnits="userSpaceOnUse" patternTransform="rotate(45)">
                  <rect width="6" height="6" fill="var(--surface)" />
                  <rect width="3" height="6" fill={COLOR_COMPROMETIDOS} />
                </pattern>
              </defs>
              <CartesianGrid vertical={false} stroke="var(--line)" />
              <XAxis dataKey="sprint" interval={0} tickLine={false} stroke="var(--muted)" />
              <YAxis
                allowDecimals={false}
                width={48}
                tickLine={false}
                axisLine={false}
                stroke="var(--muted)"
                label={{ value: 'Story Points', angle: -90, position: 'insideLeft', fill: 'var(--muted)' }}
              />
              <Tooltip
                itemSorter={(serie) => (serie.dataKey === 'comprometidos' ? 0 : 1)}
                itemStyle={{ color: 'var(--ink)' }}
                cursor={{ fill: 'rgb(96 68 147 / 8%)' }}
              />
              <Bar
                dataKey="comprometidos"
                name="Comprometidos"
                fill={`url(#${RAYADO})`}
                stroke={COLOR_COMPROMETIDOS}
                radius={[4, 4, 0, 0]}
                maxBarSize={32}
                isAnimationActive={false}
              />
              <Bar
                dataKey="completados"
                name="Completados"
                fill={COLOR_COMPLETADOS}
                radius={[4, 4, 0, 0]}
                maxBarSize={32}
                isAnimationActive={false}
              />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </div>
      <details className="dashboard-valores">
        <summary>Ver valores por Sprint</summary>
        <div className="dashboard-grafico__desplazable">
          <table aria-label="Valores por Sprint">
            <thead>
              <tr>
                <th scope="col">Sprint</th>
                <th scope="col">Comprometidos</th>
                <th scope="col">Completados</th>
              </tr>
            </thead>
            <tbody>
              {filas.map((fila) => (
                <tr key={fila.sprint}>
                  <th scope="row">{fila.sprint}</th>
                  <td>{fila.comprometidos}</td>
                  <td>{fila.completados}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </details>
    </figure>
  )
}
