/// <reference types="node" />
import { readFileSync } from 'node:fs'
import { expect, it } from 'vitest'

const estilos = readFileSync('src/App.css', 'utf8')

it('Given un enlace con clase button, When se estiliza, Then no se subraya igual que un botón', () => {
  const regla = estilos.match(/^\.button\s*\{[^}]*\}/m)?.[0] ?? ''

  expect(regla).toMatch(/text-decoration:\s*none/)
})

it('Given una pantalla de 320 px con barra de scroll que ocupa lugar, When la página tiene scroll vertical, Then body no fija un ancho mínimo que genere scroll horizontal', () => {
  const regla = estilos.match(/^body\s*\{[^}]*\}/m)?.[0] ?? ''

  expect(regla).not.toBe('')
  expect(regla).not.toMatch(/min-width/)
})
