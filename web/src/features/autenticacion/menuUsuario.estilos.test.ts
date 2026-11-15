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

it('Given la línea «Tu rol en este proyecto», When se estiliza, Then es una franja de la marca de 12,5px con contraste 4,5:1', () => {
  const bloque = regla('.menu-usuario__rol')

  expect(bloque).toMatch(/background:\s*var\(--brand-soft\)/)
  expect(bloque).toMatch(/font-size:\s*12\.5px/)
  expect(bloque).toMatch(/border-radius:\s*8px/)
  expect(contraste(hex(bloque, 'color'), '#dcf3ef')).toBeGreaterThanOrEqual(4.5)
})

it('Given el bloque de usuario del pie de la barra lateral, When se estiliza, Then queda abajo, con línea superior y texto claro con contraste 4,5:1 sobre la barra', () => {
  const bloque = regla('.bloque-usuario')

  expect(bloque).toMatch(/margin-top:\s*auto/)
  expect(bloque).toMatch(/border-top:\s*1px solid/)
  expect(contraste(hex(regla('.bloque-usuario small'), 'color'), '#0b1626')).toBeGreaterThanOrEqual(4.5)
  expect(regla('.bloque-usuario strong')).toMatch(/color:\s*#fff/)
  expect(regla('.bloque-usuario__salir')).toMatch(/width:\s*34px/)
})
