// Las fechas AAAA-MM-DD se leen como locales: `new Date('2026-10-10')` las toma como UTC y corre el día en zonas al oeste.
function aFechaLocal(iso: string): number {
  const [anio, mes, dia] = iso.slice(0, 10).split('-').map(Number)
  return new Date(anio, mes - 1, dia).getTime()
}

/** Porcentaje (0–100) de tiempo transcurrido entre dos fechas. Solo presentación: no es una métrica del proyecto. */
export function porcentajeTranscurrido(inicio: string | null, fin: string | null, hoy: Date): number | null {
  if (!inicio || !fin) return null
  const desde = aFechaLocal(inicio)
  const hasta = aFechaLocal(fin)
  const ahora = hoy.getTime()
  if (ahora <= desde) return 0
  if (ahora >= hasta) return 100
  return Math.round(((ahora - desde) / (hasta - desde)) * 100)
}

