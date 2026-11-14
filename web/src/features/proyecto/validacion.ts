export type DatosProyecto = {
  nombre: string
  descripcion: string
  fecha_inicio: string
  fecha_fin: string
  horas_por_story_point: string
}

export type ErroresProyecto = Partial<Record<keyof DatosProyecto, string>>

const NOMBRE_MAX = 100
const DESCRIPCION_MAX = 2000

export const validarCreacion = (datos: DatosProyecto): ErroresProyecto => validarDatos(datos)

export function validarEdicion(datos: DatosProyecto): ErroresProyecto {
  const errores = validarDatos(datos)
  if (!horasPorStoryPointValidas(datos.horas_por_story_point)) {
    errores.horas_por_story_point = 'Las horas por Story Point deben ser un número mayor que cero'
  }
  return errores
}

const FORMATO_HORAS = /^\d*[.,]?\d+$/

function horasPorStoryPointValidas(texto: string): boolean {
  if (!texto.trim()) return true
  const horas = aHorasPorStoryPoint(texto)
  return FORMATO_HORAS.test(texto.trim()) && horas !== null && horas > 0 && Number.isFinite(horas)
}

export function aHorasPorStoryPoint(texto: string): number | null {
  const valor = texto.trim()
  return valor ? Number(valor.replace(',', '.')) : null
}

function validarDatos(datos: DatosProyecto): ErroresProyecto {
  const errores: ErroresProyecto = {}
  const nombre = datos.nombre.trim()
  const descripcion = datos.descripcion.trim()

  if (!nombre) {
    errores.nombre = 'El nombre es obligatorio'
  } else if ([...nombre].length > NOMBRE_MAX) {
    errores.nombre = `El nombre no puede superar los ${NOMBRE_MAX} caracteres`
  }

  if ([...descripcion].length > DESCRIPCION_MAX) {
    errores.descripcion = `La descripción no puede superar los ${DESCRIPCION_MAX} caracteres`
  }

  if (datos.fecha_inicio && datos.fecha_fin && datos.fecha_fin <= datos.fecha_inicio) {
    errores.fecha_fin = 'La fecha de fin debe ser posterior a la fecha de inicio'
  }

  return errores
}

export type DatosIntegrante = {
  nombre: string
  rol: string
}

export type ErroresIntegrante = Partial<Record<keyof DatosIntegrante, string>>

const ROLES_VALIDOS = ['scrum_master', 'product_builder']

export function validarIntegrante(datos: DatosIntegrante): ErroresIntegrante {
  const errores: ErroresIntegrante = {}
  const nombre = datos.nombre.trim()

  if (!nombre) {
    errores.nombre = 'El nombre del integrante es obligatorio'
  } else if ([...nombre].length > 200) {
    errores.nombre = 'El nombre del integrante no puede superar los 200 caracteres'
  }

  if (!ROLES_VALIDOS.includes(datos.rol)) {
    errores.rol = 'Seleccioná un rol válido'
  }

  return errores
}
