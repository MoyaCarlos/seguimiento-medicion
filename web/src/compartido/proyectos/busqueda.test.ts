import { describe, expect, it } from 'vitest'
import type { ProyectoResumen } from './api'
import { filtrarProyectos } from './busqueda'

function proyecto(id: number, nombre: string, descripcion = ''): ProyectoResumen {
  return { id, nombre, descripcion, fecha_inicio: null, fecha_fin: null }
}

const PROYECTOS = [
  proyecto(1, 'Gestión de Turnos'),
  proyecto(2, 'Licencia de conducir', 'Incluye gestión de turnos'),
  proyecto(3, 'Tienda en línea'),
  proyecto(4, 'Gestión de Stock'),
]

const nombres = (lista: ProyectoResumen[]) => lista.map((p) => p.nombre)

describe('filtrarProyectos', () => {
  it('deja los proyectos cuyo nombre contiene el texto', () => {
    expect(nombres(filtrarProyectos(PROYECTOS, 'turnos'))).toEqual(['Gestión de Turnos'])
  })

  it('ignora las mayúsculas', () => {
    expect(nombres(filtrarProyectos(PROYECTOS, 'TIENDA'))).toEqual(['Tienda en línea'])
  })

  it.each(['gestion', 'GESTIÓN', 'gestión'])('ignora los acentos buscando "%s"', (texto) => {
    expect(nombres(filtrarProyectos(PROYECTOS, texto))).toEqual(['Gestión de Turnos', 'Gestión de Stock'])
  })

  it('encuentra un nombre acentuado buscando sin acento y uno sin acento buscando con acento', () => {
    expect(nombres(filtrarProyectos(PROYECTOS, 'linea'))).toEqual(['Tienda en línea'])
    expect(nombres(filtrarProyectos([proyecto(9, 'Cafeteria')], 'cafetería'))).toEqual(['Cafeteria'])
  })

  it('recorta los espacios al inicio y al final', () => {
    expect(nombres(filtrarProyectos(PROYECTOS, '  tienda  '))).toEqual(['Tienda en línea'])
  })

  it.each(['', '   '])('con el texto %j devuelve todos los proyectos', (texto) => {
    expect(filtrarProyectos(PROYECTOS, texto)).toEqual(PROYECTOS)
  })

  it('filtra solo por nombre: no mira la descripción ni el identificador', () => {
    expect(filtrarProyectos(PROYECTOS, 'incluye')).toEqual([])
    expect(filtrarProyectos(PROYECTOS, '3')).toEqual([])
  })

  it('conserva el orden recibido', () => {
    expect(nombres(filtrarProyectos([...PROYECTOS].reverse(), 'gestion'))).toEqual(['Gestión de Stock', 'Gestión de Turnos'])
  })

  it('devuelve una lista vacía cuando ningún nombre coincide', () => {
    expect(filtrarProyectos(PROYECTOS, 'zzz')).toEqual([])
  })
})
