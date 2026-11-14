/// <reference types="node" />
import { readFileSync } from 'node:fs'
import { expect, it } from 'vitest'

const css = readFileSync('src/features/proyecto/estilos.css', 'utf8')
const token = (nombre: string) =>
  readFileSync('src/App.css', 'utf8').match(new RegExp(`${nombre}:\\s*(#[0-9a-f]{6})`))?.[1] ?? ''
const regla = (selector: string) =>
  [...css.matchAll(/([^{}]+)\{([^}]*)\}/g)]
    .filter(([, lista]) => lista.split(',').some((parte) => parte.trim() === selector))
    .map(([, , cuerpo]) => cuerpo)
    .join('\n')

function luminancia(hex: string): number {
  const [r, g, b] = [1, 3, 5].map((i) => parseInt(hex.slice(i, i + 2), 16) / 255)
  const lineal = (c: number) => (c <= 0.03928 ? c / 12.92 : ((c + 0.055) / 1.055) ** 2.4)
  return 0.2126 * lineal(r) + 0.7152 * lineal(g) + 0.0722 * lineal(b)
}
const contraste = (a: string, b: string) => {
  const [claro, oscuro] = [luminancia(a), luminancia(b)].sort((x, y) => y - x)
  return (claro + 0.05) / (oscuro + 0.05)
}
const color = (bloque: string, propiedad: string) =>
  [...bloque.matchAll(new RegExp(`(?:^|[\\s;])${propiedad}:\\s*(#[0-9a-f]{6})`, 'g'))].at(-1)?.[1] ?? ''

it('Given el título del panel, When se estiliza, Then el nombre y la insignia van en fila con 12px de separación', () => {
  const bloque = regla('.title-row')

  expect(bloque).toMatch(/display:\s*flex/)
  expect(bloque).toMatch(/flex-wrap:\s*wrap/)
  expect(bloque).toMatch(/gap:\s*12px/)
})

it('Given la insignia de estado, When se estiliza, Then es una píldora de 12,5px en negrita con un punto delante', () => {
  const bloque = regla('.badge')

  expect(bloque).toMatch(/border-radius:\s*999px/)
  expect(bloque).toMatch(/font-size:\s*12\.5px/)
  expect(bloque).toMatch(/font-weight:\s*700/)
  expect(regla('.badge::before')).toMatch(/border-radius:\s*50%/)
})

it('Given cada variante de la insignia, When se estiliza, Then su texto cumple 4,5:1 contra su fondo', () => {
  const pares: [string, string][] = [
    ['.badge', token('--brand-soft')],
    ['.badge--planificado', '#f1f5f9'],
    ['.badge--finalizado', '#dcfce7'],
  ]
  for (const [selector, fondo] of pares) {
    const bloque = regla(selector)
    const texto = color(bloque, 'color')
    const fondoReal = color(bloque, 'background') || fondo

    expect(texto, selector).not.toBe('')
    expect(contraste(texto, fondoReal), selector).toBeGreaterThanOrEqual(4.5)
  }
})

it('Given las tarjetas del panel, When se estilizan, Then la descripción ocupa 1,6 fr y el equipo 1 fr', () => {
  expect(regla('.panel-rejilla')).toMatch(/grid-template-columns:\s*minmax\(0, 1\.6fr\) minmax\(0, 1fr\)/)
  expect(regla('.panel-rejilla')).toMatch(/gap:\s*20px/)
})

it('Given una tarjeta del panel, When se estiliza, Then tiene radio 14 y relleno de 22px', () => {
  const bloque = regla('.panel-tarjeta')

  expect(bloque).toMatch(/border-radius:\s*var\(--radius\)/)
  expect(bloque).toMatch(/padding:\s*22px/)
})

it('Given el título de una tarjeta, When se estiliza, Then es una etiqueta de 12px en mayúsculas', () => {
  const bloque = regla('.panel-tarjeta__titulo')

  expect(bloque).toMatch(/font-size:\s*12px/)
  expect(bloque).toMatch(/text-transform:\s*uppercase/)
  expect(bloque).toMatch(/letter-spacing:\s*0\.08em/)
})

it('Given una descripción larga sin espacios, When se estiliza, Then se corta en varias líneas', () => {
  expect(regla('.panel-prosa')).toMatch(/overflow-wrap:\s*anywhere/)
})

it('Given las fechas del cronograma, When se estilizan, Then usan 17px con cifras tabulares', () => {
  const bloque = regla('.cronograma__fechas strong')

  expect(bloque).toMatch(/font-size:\s*17px/)
  expect(bloque).toMatch(/font-variant-numeric:\s*tabular-nums/)
})

it('Given el enlace «Ver estado del proyecto», When se estiliza, Then usa el color de marca oscuro, en negrita y con la flecha de 15px', () => {
  expect(regla('.enlace')).toMatch(/color:\s*var\(--brand-dark\)/)
  expect(regla('.enlace')).toMatch(/font-weight:\s*700/)
  expect(regla('.enlace svg')).toMatch(/width:\s*15px/)
})
