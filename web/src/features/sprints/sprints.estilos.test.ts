/// <reference types="node" />
import { readFileSync } from 'node:fs'
import { expect, it } from 'vitest'

const css = readFileSync('src/features/sprints/sprints.css', 'utf8')
const regla = (selector: string) =>
  [...css.matchAll(/([^{}]+)\{([^}]*)\}/g)]
    .filter(([, lista]) => lista.split(',').some((parte) => parte.trim() === selector))
    .map(([, , cuerpo]) => cuerpo)
    .join('\n')
const color = (bloque: string, propiedad: string) =>
  [...bloque.matchAll(new RegExp(`(?:^|[\\s;])${propiedad}:\\s*(#[0-9a-f]{6})`, 'g'))].at(-1)?.[1] ?? ''
function luminancia(hex: string): number {
  const [r, g, b] = [1, 3, 5].map((i) => parseInt(hex.slice(i, i + 2), 16) / 255)
  const lineal = (c: number) => (c <= 0.03928 ? c / 12.92 : ((c + 0.055) / 1.055) ** 2.4)
  return 0.2126 * lineal(r) + 0.7152 * lineal(g) + 0.0722 * lineal(b)
}
const contraste = (a: string, b: string) => {
  const [claro, oscuro] = [luminancia(a), luminancia(b)].sort((x, y) => y - x)
  return (claro + 0.05) / (oscuro + 0.05)
}

it('Given el tablero, When se estiliza, Then son tres columnas de igual ancho con 16px de hueco', () => {
  const bloque = regla('.sprints-tablero')

  expect(bloque).toMatch(/grid-template-columns:\s*repeat\(3, minmax\(0, 1fr\)\)/)
  expect(bloque).toMatch(/gap:\s*16px/)
  expect(bloque).toMatch(/align-items:\s*start/)
})

it('Given el ancho de tablet o menos, When se estiliza, Then las columnas se apilan', () => {
  expect(css).toMatch(/@media \(max-width: 1100px\)\s*\{[^@]*\.sprints-tablero\s*\{[^}]*grid-template-columns:\s*1fr/)
})

it('Given una columna, When se estiliza, Then es una tarjeta gris de radio 16 y relleno 12', () => {
  const bloque = regla('.sprints-carril')

  expect(bloque).toMatch(/border-radius:\s*16px/)
  expect(bloque).toMatch(/padding:\s*12px/)
  expect(bloque).toMatch(/background:\s*rgb\(15 27 45 \/ 3\.5%\)/)
})

it('Given el título de cada columna, When se estiliza, Then es una etiqueta de 12px en mayúsculas con un punto de color por estado', () => {
  const titulo = regla('.sprints-carril h2')

  expect(titulo).toMatch(/font-size:\s*12px/)
  expect(titulo).toMatch(/text-transform:\s*uppercase/)
  expect(regla('.sprints-carril h2::before')).toMatch(/border-radius:\s*50%/)
  expect(regla('.sprints-carril h2::before')).toMatch(/background:\s*#94a3b8/)
  expect(regla('.sprints-carril--activo h2::before')).toMatch(/background:\s*var\(--brand\)/)
  expect(regla('.sprints-carril--fin h2::before')).toMatch(/background:\s*#15803d/)
})

it('Given el título de cada columna, When se estiliza, Then su texto cumple 4,5:1 contra el fondo de la columna', () => {
  for (const [selector, fondo] of [
    ['.sprints-carril h2', '#f0f1f3'],
    ['.sprints-carril--activo h2', '#f0f1f3'],
    ['.sprints-carril--fin h2', '#f0f1f3'],
  ]) {
    const texto = color(regla(selector), 'color')

    expect(texto, selector).not.toBe('')
    expect(contraste(texto, fondo), selector).toBeGreaterThanOrEqual(4.5)
  }
})

it('Given la tarjeta de un Sprint, When se estiliza, Then tiene radio 14 y relleno 16', () => {
  const bloque = regla('.sprint-tarjeta')

  expect(bloque).toMatch(/border-radius:\s*14px/)
  expect(bloque).toMatch(/padding:\s*16px/)
})

it('Given las fechas del Sprint, When se estilizan, Then son una etiqueta de 12px con cifras tabulares', () => {
  const bloque = regla('.sprint-fechas')

  expect(bloque).toMatch(/border-radius:\s*7px/)
  expect(bloque).toMatch(/font-size:\s*12px/)
  expect(bloque).toMatch(/font-variant-numeric:\s*tabular-nums/)
  expect(regla('.sprint-fechas svg')).toMatch(/width:\s*13px/)
})

it('Given la columna vacía, When se estiliza, Then es un recuadro punteado centrado', () => {
  const bloque = regla('.sprints-vacio')

  expect(bloque).toMatch(/border:\s*1px dashed/)
  expect(bloque).toMatch(/text-align:\s*center/)
})

it('Given el plazo del proyecto, When se estiliza, Then es una línea de 13,5px', () => {
  expect(regla('.sprints-plazo')).toMatch(/font-size:\s*13\.5px/)
})
