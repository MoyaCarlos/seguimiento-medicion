/// <reference types="node" />
import { readFileSync } from 'node:fs'
import { expect, it } from 'vitest'

const css = readFileSync('src/features/defectos/defectos.css', 'utf8')
const regla = (selector: string) =>
  css.match(new RegExp(`^${selector.replace(/[.:]/g, '\\$&')}\\s*\\{([^}]*)\\}`, 'm'))?.[1] ?? ''
const color = (bloque: string, propiedad: string) =>
  bloque.match(new RegExp(`(?:^|[\\s;])${propiedad}:\\s*(#[0-9a-f]{6})`))?.[1] ?? ''

function luminancia(hex: string): number {
  const [r, g, b] = [1, 3, 5].map((i) => parseInt(hex.slice(i, i + 2), 16) / 255)
  const lineal = (c: number) => (c <= 0.03928 ? c / 12.92 : ((c + 0.055) / 1.055) ** 2.4)
  return 0.2126 * lineal(r) + 0.7152 * lineal(g) + 0.0722 * lineal(b)
}
const contraste = (a: string, b: string) => {
  const [claro, oscuro] = [luminancia(a), luminancia(b)].sort((x, y) => y - x)
  return (claro + 0.05) / (oscuro + 0.05)
}

it.each(['alta', 'media', 'baja', 'pendiente', 'resuelto'])(
  'Given la etiqueta %s, When se estiliza, Then su texto cumple 4,5:1 sobre su fondo',
  (nombre) => {
    const bloque = regla(`.defecto__etiqueta--${nombre}`)

    expect(contraste(color(bloque, 'color'), color(bloque, 'background'))).toBeGreaterThanOrEqual(4.5)
  },
)

it('Given una etiqueta, When se estiliza, Then es un rectángulo chico (radio de 7px) y no un óvalo alto', () => {
  const bloque = regla('.defecto__etiqueta')

  expect(bloque).toMatch(/border-radius:\s*7px/)
  expect(bloque).toMatch(/padding:\s*2px 9px/)
})

it('Given la pantalla, When hay lugar, Then el formulario y el listado van en dos columnas de al menos 300px', () => {
  expect(regla('.defectos-columnas')).toMatch(/grid-template-columns:\s*repeat\(auto-fit, minmax\(300px, 1fr\)\)/)
})
