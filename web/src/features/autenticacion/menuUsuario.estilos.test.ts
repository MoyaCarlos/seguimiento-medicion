/// <reference types="node" />
import { readFileSync } from 'node:fs'
import { expect, it } from 'vitest'

const css = readFileSync('src/features/autenticacion/estilos.css', 'utf8')
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
const hex = (bloque: string, propiedad: string) =>
  bloque.match(new RegExp(`(?:^|[\\s;])${propiedad}:\\s*(#[0-9a-f]{6})`))?.[1] ?? ''

it('Given el botón del menú de usuario, When se estiliza, Then es una píldora de 40px con borde', () => {
  const bloque = regla('.menu-usuario__boton')

  expect(bloque).toMatch(/height:\s*40px/)
  expect(bloque).toMatch(/border-radius:\s*999px/)
  expect(bloque).toMatch(/border:\s*1px solid/)
})

it('Given el nombre del botón, When se estiliza, Then se corta a 14ch y desaparece en 640px o menos', () => {
  expect(regla('.menu-usuario__nombre')).toMatch(/max-width:\s*14ch/)
  const movil = css.match(/@media \(max-width: 640px\)\s*\{([\s\S]*?)\n\}/g)?.join('\n') ?? ''

  expect(movil).toMatch(/\.menu-usuario__nombre\s*\{[^}]*display:\s*none/)
})

it('Given el panel del menú, When se estiliza, Then flota a la derecha del botón, de 268px, con radio 14', () => {
  const bloque = regla('.menu-usuario__panel')

  expect(bloque).toMatch(/position:\s*absolute/)
  expect(bloque).toMatch(/right:\s*0/)
  expect(bloque).toMatch(/top:\s*calc\(100% \+ 8px\)/)
  expect(bloque).toMatch(/width:\s*268px/)
  expect(bloque).toMatch(/border-radius:\s*14px/)
})

it('Given «Cerrar sesión», When se estiliza, Then va en rojo con contraste 4,5:1 en reposo y al pasar el mouse', () => {
  const reposo = regla('.menu-usuario__item--peligro')
  const encima = regla('.menu-usuario__item--peligro:hover')

  expect(contraste(hex(reposo, 'color'), '#ffffff')).toBeGreaterThanOrEqual(4.5)
  expect(contraste(hex(reposo, 'color'), hex(encima, 'background'))).toBeGreaterThanOrEqual(4.5)
})

it('Given la «Sesión iniciada» de la cabecera, When se estiliza, Then cumple 4,5:1 sobre blanco', () => {
  expect(contraste(hex(regla('.menu-usuario__cabecera small'), 'color'), '#ffffff')).toBeGreaterThanOrEqual(4.5)
  expect(regla('.menu-usuario__cabecera small')).toMatch(/color:\s*#606f86/)
})
