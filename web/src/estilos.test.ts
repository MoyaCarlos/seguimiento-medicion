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

it('Given el mensaje de datos de prueba cargados, When se estiliza, Then es una caja verde suave con texto que cumple 4,5:1', () => {
  const css = readFileSync('src/features/autenticacion/estilos.css', 'utf8')
  const regla = css.match(/^\.demo-exito\s*\{([^}]*)\}/m)?.[1] ?? ''
  const fondo = regla.match(/background:\s*(#[0-9a-f]{6})/)?.[1] ?? ''
  const texto = regla.match(/\bcolor:\s*(#[0-9a-f]{6})/)?.[1] ?? ''
  const luz = (h: string) =>
    [1, 3, 5].map((i) => parseInt(h.slice(i, i + 2), 16) / 255).map((c) => (c <= 0.03928 ? c / 12.92 : ((c + 0.055) / 1.055) ** 2.4))
      .reduce((a, c, i) => a + c * [0.2126, 0.7152, 0.0722][i], 0)

  expect(fondo).toBe('#dcf5e4')
  const [a, b] = [luz(fondo), luz(texto)].sort((x, y) => y - x)
  expect((a + 0.05) / (b + 0.05)).toBeGreaterThanOrEqual(4.5)
})
