/// <reference types="node" />
import { readdirSync, readFileSync, statSync } from 'node:fs'
import { join } from 'node:path'
import { expect, it } from 'vitest'

function hojas(dir: string): string[] {
  return readdirSync(dir).flatMap((nombre) => {
    const ruta = join(dir, nombre)
    if (nombre === 'node_modules') return []
    return statSync(ruta).isDirectory() ? hojas(ruta) : ruta.endsWith('.css') ? [ruta] : []
  })
}

// Una regla sin su llave de cierre deja las siguientes anidadas dentro de ella: el CSS compila (lightningcss
// admite anidado) pero esas reglas dejan de aplicarse. Solo los bloques @… pueden contener otros bloques.
function reglasAnidadas(css: string): string[] {
  const sinComentarios = css.replace(/\/\*[\s\S]*?\*\//g, '')
  const pila: string[] = []
  const anidadas: string[] = []
  let cabecera = ''
  for (const caracter of sinComentarios) {
    if (caracter === '{') {
      const selector = cabecera.trim()
      if (pila.length > 0 && !pila[pila.length - 1].startsWith('@')) anidadas.push(selector)
      pila.push(selector)
      cabecera = ''
    } else if (caracter === '}') {
      pila.pop()
      cabecera = ''
    } else if (caracter === ';') {
      cabecera = ''
    } else {
      cabecera += caracter
    }
  }
  return pila.length === 0 ? anidadas : [...anidadas, 'llaves sin cerrar']
}

it.each(hojas('src'))('Given %s, When se lee, Then ninguna regla queda anidada por una llave de cierre faltante', (ruta) => {
  expect(reglasAnidadas(readFileSync(ruta, 'utf8'))).toEqual([])
})
