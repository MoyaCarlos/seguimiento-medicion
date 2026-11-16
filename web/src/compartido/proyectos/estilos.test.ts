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

it('Given una tarjeta de proyecto, When se muestra el pie, Then el #id queda a la derecha', () => {
  expect(bloque('.tarjeta-proyecto__pie')).toMatch(/justify-content:\s*flex-end/)
})

it('Given los botones de la pantalla Proyectos, When llevan ícono, Then el ícono y el texto van separados 8px', () => {
  expect(bloque('.lista-proyectos__acciones .button')).toMatch(/gap:\s*8px/)
})

it('Given el formulario de la ventana, When se muestra la descripción, Then el área de texto mide al menos 92px', () => {
  expect(bloque('.proyecto-crear textarea')).toMatch(/min-height:\s*92px/)
})
