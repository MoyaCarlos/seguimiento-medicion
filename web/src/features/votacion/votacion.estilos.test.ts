/// <reference types="node" />
import { readFileSync } from 'node:fs'
import { expect, it } from 'vitest'

const css = readFileSync('src/features/votacion/votacion.css', 'utf8')
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

it('Given mi fila, When se estiliza, Then lleva borde y halo de la marca', () => {
  const bloque = regla('.votacion-integrante--yo')

  expect(bloque).toMatch(/border-color:\s*var\(--brand\)/)
  expect(bloque).toMatch(/box-shadow:\s*0 0 0 1px var\(--brand\)/)
})

it('Given «Ya votó», When se estiliza, Then es verde, en negrita y cumple 4,5:1 contra blanco', () => {
  const bloque = regla('.votacion-integrante__estado')

  expect(bloque).toMatch(/font-weight:\s*700/)
  expect(contraste(color(bloque, 'color'), '#ffffff')).toBeGreaterThanOrEqual(4.5)
  expect(color(bloque, 'color')).toBe('#166534')
})

it('Given el texto de diferencias, When se estiliza, Then es texto en negrita sin tarjeta ni franja de color', () => {
  const bloque = regla('.votacion-diferencias')

  expect(bloque).toMatch(/font-weight:\s*700/)
  expect(bloque).not.toMatch(/background/)
  expect(bloque).not.toMatch(/border/)
  expect(css).not.toMatch(/votacion-diferencias--hay/)
})

it('Given el consenso, When se estiliza, Then el selector y el botón van en una fila alineada abajo', () => {
  const bloque = regla('.votacion-consenso__fila')

  expect(bloque).toMatch(/align-items:\s*flex-end/)
  expect(bloque).toMatch(/flex-wrap:\s*wrap/)
})

it('Given «Nueva ronda», When se estiliza, Then va debajo de su texto, sin línea separadora y alineada a la izquierda', () => {
  const bloque = regla('.votacion-consenso__ronda')

  expect(bloque).not.toMatch(/border-top/)
  expect(bloque).toMatch(/justify-items:\s*start/)
  expect(bloque).toMatch(/display:\s*grid/)
})

it('Given los párrafos del consenso, When se estilizan, Then usan el color de texto principal, como el prototipo, y cumplen 4,5:1', () => {
  const bloque = regla('.votacion-consenso p')

  expect(bloque).toMatch(/color:\s*#334155/)
  expect(contraste('#334155', '#ffffff')).toBeGreaterThanOrEqual(4.5)
})
