export type DatosProyecto = {
  nombre: string
  descripcion: string
  fecha_inicio: string
  fecha_fin: string
  creador: string
}

export type ErroresProyecto = Partial<Record<keyof DatosProyecto, string>>

const NOMBRE_MAX = 100
const DESCRIPCION_MAX = 2000

export function validarCreacion(datos: DatosProyecto): ErroresProyecto {
  const errores = validarDatos(datos)
  if (!datos.creador.trim()) {
    errores.creador = 'El nombre del creador es obligatorio'
  }
  return errores
}

export function validarEdicion(datos: DatosProyecto): ErroresProyecto {
  return validarDatos(datos)
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
