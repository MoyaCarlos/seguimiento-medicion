import { Bar, BarChart, CartesianGrid, Legend, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts'
import { datosDelGrafico, type MetricasDeSprint } from './metricas'

const ALTO = 320
const ANCHO_POR_SPRINT = 88
const COLOR_COMPROMETIDOS = '#2a9d8f'
const COLOR_COMPLETADOS = '#604493'
const RAYADO = 'dashboard-comprometidos-rayado'

export function GraficoVelocidad({ sprints }: { sprints: MetricasDeSprint[] }) {
  const filas = datosDelGrafico(sprints)
  const anchoMinimo = filas.length * ANCHO_POR_SPRINT

  return (
    <figure className="dashboard-tarjeta dashboard-grafico" aria-labelledby="dashboard-grafico-titulo">
      <figcaption id="dashboard-grafico-titulo">
        <h2>Story Points comprometidos y completados por Sprint</h2>
      </figcaption>
      <div className="dashboard-grafico__desplazable">
        <div style={{ minWidth: anchoMinimo }}>
          <ResponsiveContainer width="100%" height={ALTO} initialDimension={{ width: anchoMinimo, height: ALTO }}>
            <BarChart data={filas} barGap={2} margin={{ top: 8, right: 8, bottom: 8, left: 0 }}>
              <defs>
                <pattern id={RAYADO} width="6" height="6" patternUnits="userSpaceOnUse" patternTransform="rotate(45)">
                  <rect width="6" height="6" fill="#ffffff" />
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
              <Tooltip cursor={{ fill: 'rgb(96 68 147 / 8%)' }} />
              <Legend />
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
