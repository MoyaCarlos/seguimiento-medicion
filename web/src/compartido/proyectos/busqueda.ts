import type { ProyectoResumen } from './api'

function normalizar(texto: string): string {
  return texto.normalize('NFD').replace(/\p{M}/gu, '').toLowerCase()
}

export function filtrarProyectos(proyectos: ProyectoResumen[], texto: string): ProyectoResumen[] {
  const buscado = normalizar(texto.trim())
  if (!buscado) return proyectos
  return proyectos.filter((proyecto) => normalizar(proyecto.nombre).includes(buscado))
}
