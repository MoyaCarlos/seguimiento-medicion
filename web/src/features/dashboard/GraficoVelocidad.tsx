import { Bar, BarChart, CartesianGrid, LabelList, ResponsiveContainer, XAxis, YAxis } from 'recharts'
import { datosDelGrafico, ejeDelGrafico, type MetricasDeSprint } from './metricas'

const ALTO = 230
const ANCHO_MINIMO = 520
const ANCHO_POR_SPRINT = 120
const COLOR_COMPROMETIDOS = '#64748b'
const COLOR_COMPLETADOS = '#0d9488'
const RAYADO = 'dashboard-comprometidos-rayado'
const TEXTO_DEL_EJE = { fill: '#64748b', fontSize: 11.5 }

export function GraficoVelocidad({ sprints }: { sprints: MetricasDeSprint[] }) {
  const filas = datosDelGrafico(sprints)
  const anchoMinimo = Math.max(ANCHO_MINIMO, filas.length * ANCHO_POR_SPRINT)
  const marcas = ejeDelGrafico(Math.max(0, ...filas.flatMap((fila) => [fila.comprometidos, fila.completados])))

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
        <div style={{ width: anchoMinimo }}>
          <ResponsiveContainer width="100%" height={ALTO} initialDimension={{ width: anchoMinimo, height: ALTO }}>
            <BarChart data={filas} barGap={8} margin={{ top: 26, right: 10, bottom: 8, left: 0 }}>
              <defs>
                <pattern id={RAYADO} width="6" height="6" patternUnits="userSpaceOnUse" patternTransform="rotate(45)">
                  <rect width="6" height="6" fill="#eef2f6" />
                  <rect width="2.5" height="6" fill="#94a3b8" />
                </pattern>
              </defs>
              <CartesianGrid vertical={false} stroke="#e2e8f0" />
              <XAxis dataKey="sprint" interval={0} tickLine={false} axisLine={false} tick={TEXTO_DEL_EJE} />
              <YAxis
                ticks={marcas}
                width={44}
                tickLine={false}
                axisLine={false}
                tick={TEXTO_DEL_EJE}
              />
              <Bar
                dataKey="comprometidos"
                name="Comprometidos"
                fill={`url(#${RAYADO})`}
                stroke={COLOR_COMPROMETIDOS}
                strokeWidth={1.2}
                radius={3}
                barSize={34}
                isAnimationActive={false}
              >
                <LabelList dataKey="comprometidos" position="top" fill="#0f1b2d" fontSize={12} fontWeight={700} />
              </Bar>
              <Bar
                dataKey="completados"
                name="Completados"
                fill={COLOR_COMPLETADOS}
                radius={3}
                barSize={34}
                isAnimationActive={false}
              >
                <LabelList dataKey="completados" position="top" fill="#0f1b2d" fontSize={12} fontWeight={700} />
              </Bar>
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
