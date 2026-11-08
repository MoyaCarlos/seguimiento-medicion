import { describe, expect, it } from 'vitest'
import { etiquetaProyecto, rutaDeSeccion } from './ruta'

describe('rutaDeSeccion', () => {
  it.each([
    ['conserva la sección del menú', '/projects/1/backlog', 2, '/projects/2/backlog'],
    ['descarta las sub-rutas', '/projects/1/proyecto/edit', 2, '/projects/2/proyecto'],
    ['usa proyecto cuando no hay sección', '/projects/1', 2, '/projects/2/proyecto'],
    ['ignora la barra final', '/projects/1/backlog/', '2', '/projects/2/backlog'],
    ['cae en proyecto si la sección no es del menú', '/projects/1/otra', 2, '/projects/2/proyecto'],
    ['no traslada el id de una historia', '/projects/1/votacion/3', 2, '/projects/2/proyecto'],
  ])('%s', (_caso, pathname, idNuevo, esperada) => {
    expect(rutaDeSeccion(pathname, idNuevo)).toBe(esperada)
  })
})

describe('etiquetaProyecto', () => {
  it('muestra nombre e identificador', () => {
    expect(etiquetaProyecto({ id: 7, nombre: 'Alfa' })).toBe('Alfa (#7)')
  })

  it('distingue nombres repetidos por el identificador', () => {
    expect(etiquetaProyecto({ id: 8, nombre: 'Alfa' })).toBe('Alfa (#8)')
  })
})
