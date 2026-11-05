/// <reference types="node" />
import { readFileSync } from 'node:fs'
import { expect, it } from 'vitest'

const estilos = readFileSync('src/App.css', 'utf8')
const estilosDeCelular = estilos.split('@media (max-width: 600px)')[1] ?? ''

it('Given una pantalla angosta, When las pestañas no entran en una línea, Then pasan a la siguiente sin desbordar la página', () => {
  const regla = estilosDeCelular.match(/\.feature-navigation\s*\{[^}]*\}/)?.[0] ?? ''

  expect(regla).toMatch(/flex-wrap:\s*wrap/)
})
