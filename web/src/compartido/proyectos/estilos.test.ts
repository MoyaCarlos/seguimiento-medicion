/// <reference types="node" />
import { readFileSync } from 'node:fs'
import { expect, it } from 'vitest'

const css = readFileSync('src/App.css', 'utf8')
const bloque = (selector: string) =>
  css.match(new RegExp(`^${selector.replace(/[.:[\]]/g, '\\$&')}\\s*\\{([^}]*)\\}`, 'm'))?.[1] ?? ''

it('Given la grilla de proyectos, When hay lugar, Then arma columnas de al menos 300px con 18px de separación', () => {
  expect(bloque('.lista-proyectos__tarjetas')).toMatch(/repeat\(auto-fill, minmax\(300px, 1fr\)\)/)
  expect(bloque('.lista-proyectos__tarjetas')).toMatch(/gap:\s*18px/)
})

it('Given una tarjeta de proyecto, When se estiliza, Then es una tarjeta cuadrada de 16px de radio y 20px de relleno', () => {
  const tarjeta = bloque('.tarjeta-proyecto')

  expect(tarjeta).toMatch(/border-radius:\s*16px/)
  expect(tarjeta).toMatch(/padding:\s*20px/)
  expect(tarjeta).toMatch(/flex-direction:\s*column/)
})

it('Given la tarjeta «Nuevo proyecto», When se estiliza, Then tiene borde punteado de 2px y alto mínimo de 220px', () => {
  const nueva = bloque('.tarjeta-proyecto--nueva')

  expect(nueva).toMatch(/border:\s*2px dashed/)
  expect(nueva).toMatch(/min-height:\s*220px/)
})

it('Given una tarjeta con el foco del teclado, When se dibuja, Then el contorno es sólido y se distingue 3:1', () => {
  expect(bloque('.tarjeta-proyecto:focus-visible')).toMatch(/outline:\s*3px solid var\(--brand\)/)
})
