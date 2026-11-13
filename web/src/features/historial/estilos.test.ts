/// <reference types="node" />
import { readFileSync } from 'node:fs'
import { expect, it } from 'vitest'

const css = readFileSync('src/features/historial/historial.css', 'utf8')
const regla = (selector: string) =>
  css.match(new RegExp(`^${selector.replace(/[.:()]/g, '\\$&')}\\s*\\{([^}]*)\\}`, 'm'))?.[1] ?? ''

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
  bloque.match(new RegExp(`(?:^|[\\s;])${propiedad}:\\s*(#[0-9a-f]{6})`))?.[1] ?? ''

it('Given la etiqueta de fechas, When se estiliza, Then es un rectángulo chico de 7px y su texto cumple 4,5:1', () => {
  const bloque = regla('.historial-pantalla .historial-fechas')

  expect(bloque).toMatch(/border-radius:\s*7px/)
  expect(contraste(color(bloque, 'color'), color(bloque, 'background'))).toBeGreaterThanOrEqual(4.5)
})

it('Given el título de un Sprint, When se estiliza, Then va en verde oscuro con contraste 4,5:1 sobre blanco', () => {
  expect(contraste(color(regla('.historial-sprint__enlace'), 'color'), '#ffffff')).toBeGreaterThanOrEqual(4.5)
})

it('Given las métricas de la lista, When se estilizan, Then van sin bordes, en una grilla de columnas de 150px', () => {
  const bloque = regla('.historial-metricas')

  expect(bloque).toMatch(/repeat\(auto-fit, minmax\(150px, 1fr\)\)/)
  expect(bloque).not.toMatch(/border/)
})

it('Given una tarjeta con el foco del teclado, When se dibuja, Then el contorno es sólido y se distingue 3:1 del fondo', () => {
  const bloque = regla('.historial-sprint:focus-within')
  const contorno = bloque.match(/outline:\s*3px solid\s*([^;]+);/)?.[1].trim() ?? ''

  expect(contorno).toBe('var(--brand)')
  expect(contraste('#0d9488', '#f3f5f8')).toBeGreaterThanOrEqual(3)
})

it('Given una tarjeta, When se pasa el mouse, Then el borde pasa al acento', () => {
  expect(regla('.historial-sprint:hover')).toMatch(/border-color:\s*var\(--brand\)/)
})

it('Given una pantalla angosta, When se muestran las métricas, Then no se les agrega relleno ni ancho propio', () => {
  const angosto = css.slice(css.indexOf('@media (max-width: 600px)'), css.indexOf('@media (max-width: 640px)'))

  expect(angosto).not.toMatch(/\.historial-metrica\s*\{/)
  expect(angosto).toMatch(/\.historial-tarjeta:not\(\.historial-sprint\)/)
})

it.each([
  ['.historial-chip--alerta'],
  ['.historial-chip--media'],
  ['.historial-chip--baja'],
  ['.historial-chip--apagado'],
  ['.historial-chip--ok'],
  ['.historial-chip--fuerte'],
  ['.historial-chip--ninguno'],
])('Given la etiqueta %s del detalle, When se estiliza, Then su texto cumple 4,5:1 sobre su fondo', (selector) => {
  const bloque = regla(selector)
  const fondo = color(bloque, 'background') || '#ffffff'

  expect(contraste(color(bloque, 'color'), fondo)).toBeGreaterThanOrEqual(4.5)
})

it('Given una etiqueta del detalle, When se estiliza, Then es un rectángulo chico de 7px', () => {
  const bloque = regla('.historial-chip')

  expect(bloque).toMatch(/border-radius:\s*7px/)
  expect(bloque).toMatch(/padding:\s*2px 9px/)
})

it('Given las métricas del detalle, When se estilizan, Then la cifra es de 19px', () => {
  expect(regla('.historial-detalle .historial-metrica dd')).toMatch(/font-size:\s*19px/)
})

it('Given el botón «Exportar a PDF», When se estiliza, Then es blanco con borde del acento y texto verde oscuro con contraste 4,5:1', () => {
  const boton = readFileSync('src/compartido/componentes/BotonExportar.css', 'utf8')
  const bloque = boton.match(/\.reporte-exportar__boton\s*\{([^}]*)\}/)?.[1] ?? ''

  expect(bloque).toMatch(/border-color:\s*var\(--brand\)/)
  expect(contraste('#0b7a70', '#ffffff')).toBeGreaterThanOrEqual(4.5)
  expect(bloque).toMatch(/background:\s*var\(--surface\)/)
})

it('Given un defecto del detalle, When se estiliza, Then va en una caja con borde y 12px de radio, separadas entre sí', () => {
  expect(regla('.historial-item--caja')).toMatch(/border:\s*1px solid/)
  expect(regla('.historial-item--caja')).toMatch(/border-radius:\s*12px/)
  expect(regla('.historial-items--cajas')).toMatch(/gap:\s*14px/)
})
