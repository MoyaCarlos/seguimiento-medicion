/// <reference types="node" />
import { existsSync, readFileSync } from 'node:fs'
import { expect, it } from 'vitest'

const estilos = readFileSync('src/App.css', 'utf8')

const bloque = (selector: string) =>
  estilos.match(new RegExp(`^${selector.replace(/[.:]/g, '\\$&')}\\s*\\{([^}]*)\\}`, 'm'))?.[1] ?? ''

const tokens = Object.fromEntries(
  [...bloque(':root').matchAll(/(--[\w-]+):\s*([^;]+);/g)].map(([, nombre, valor]) => [nombre, valor.trim()]),
)

function resolver(valor: string): string {
  const alias = valor.match(/^var\((--[\w-]+)\)$/)
  return alias ? resolver(tokens[alias[1]]) : valor
}

const token = (nombre: string) => resolver(tokens[nombre] ?? '')
const propiedad = (selector: string, nombre: string) =>
  resolver(bloque(selector).match(new RegExp(`${nombre}:\\s*([^;]+);`))?.[1].trim() ?? '')

function luminancia(hex: string): number {
  const [r, g, b] = [1, 3, 5].map((i) => parseInt(hex.slice(i, i + 2), 16) / 255)
  const lineal = (c: number) => (c <= 0.03928 ? c / 12.92 : ((c + 0.055) / 1.055) ** 2.4)
  return 0.2126 * lineal(r) + 0.7152 * lineal(g) + 0.0722 * lineal(b)
}

function contraste(a: string, b: string): number {
  const [claro, oscuro] = [luminancia(a), luminancia(b)].sort((x, y) => y - x)
  return (claro + 0.05) / (oscuro + 0.05)
}

const FONDO_PAGINA = '#f3f5f8'
const FONDO_SUAVE_DEL_PROTOTIPO = '#eef2f6'
const BLANCO = '#ffffff'

it('Given los tokens del prototipo, When se leen, Then el acento y el verde de los enlaces son los del prototipo', () => {
  expect(token('--brand')).toBe('#0d9488')
  expect(token('--brand-dark')).toBe('#0b7a70')
})

it.each([
  ['--muted', [FONDO_PAGINA, FONDO_SUAVE_DEL_PROTOTIPO, BLANCO]],
  ['--danger', [FONDO_PAGINA, FONDO_SUAVE_DEL_PROTOTIPO, BLANCO]],
  ['--brand-dark', [FONDO_PAGINA, BLANCO]],
])('Given el texto en %s, When se pone sobre %j, Then cumple el contraste AA de 4,5:1', (nombre, fondos) => {
  for (const fondo of fondos as string[]) expect(contraste(token(nombre), fondo)).toBeGreaterThanOrEqual(4.5)
})

it('Given el botón primario, When está en reposo, Then conserva el verde oscuro con texto blanco AA y no usa el acento', () => {
  const reposo = propiedad('.button--primary', 'background')

  expect(reposo).toBe('#0f766e')
  expect(reposo).not.toBe(token('--brand'))
  expect(contraste(BLANCO, reposo)).toBeGreaterThanOrEqual(4.5)
})

it('Given el botón primario, When se pasa el mouse, Then el fondo es más oscuro que en reposo', () => {
  const reposo = propiedad('.button--primary', 'background')
  const encima = propiedad('.button--primary:hover', 'background')

  expect(luminancia(encima)).toBeLessThan(luminancia(reposo))
})

it('Given los alias de la ola 2, When cambian los tokens, Then siguen con los valores anteriores', () => {
  expect(token('--lilac')).toBe('#0f766e')
  expect(token('--lilac-dark')).toBe('#115e59')
  expect(token('--complement')).toBe('#115e59')
})

it('Given el texto de la página, When se estiliza body, Then usa 15 px con interlineado 1,5 como el prototipo', () => {
  const regla = bloque('body')

  expect(regla).toMatch(/font-size:\s*15px/)
  expect(regla).toMatch(/line-height:\s*1\.5\b/)
})

it('Given la tipografía Inter, When se declara, Then sale de un woff2 propio con licencia y es la primera de la pila', () => {
  const fuente = estilos.match(/@font-face\s*\{([^}]*)\}/)?.[1] ?? ''
  const archivo = fuente.match(/url\(['"]?\.\/([^'")]+\.woff2)['"]?\)/)?.[1] ?? ''

  expect(fuente).toMatch(/font-family:\s*['"]Inter['"]/)
  expect(fuente).toMatch(/font-display:\s*swap/)
  expect(fuente).toMatch(/font-weight:\s*100 900/)
  expect(existsSync(`src/${archivo}`)).toBe(true)
  expect(existsSync('src/fuentes/LICENSE-Inter.txt')).toBe(true)
  expect(bloque(':root').match(/font-family:\s*([^;]+);/)?.[1]).toMatch(/^"Inter"/)
})

it('Given el foco visible, When se dibuja, Then el contorno es sólido y se distingue 3:1 de los fondos', () => {
  const contorno = resolver(bloque(':focus-visible').match(/outline:\s*3px solid\s*([^;]+);/)?.[1].trim() ?? '')

  expect(contorno).toMatch(/^#[0-9a-f]{6}$/)
  expect(contraste(contorno, FONDO_PAGINA)).toBeGreaterThanOrEqual(3)
  expect(contraste(contorno, BLANCO)).toBeGreaterThanOrEqual(3)
})
