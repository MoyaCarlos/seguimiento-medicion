/// <reference types="node" />
import { readFileSync } from 'node:fs'
import { expect, it } from 'vitest'

const css = readFileSync('src/features/dashboard/dashboard.css', 'utf8')
const regla = (selector: string) =>
  css.match(new RegExp(`^${selector.replace(/[.:]/g, '\\$&')}\\s*\\{([^}]*)\\}`, 'm'))?.[1] ?? ''

it('Given las tarjetas de indicadores, When hay lugar, Then van en una grilla de columnas de al menos 230px', () => {
  expect(regla('.dashboard-indicadores')).toMatch(/repeat\(auto-fit, minmax\(230px, 1fr\)\)/)
})

it('Given el rótulo de un indicador, When se estiliza, Then va en mayúscula chica y la cifra en 26px verde oscuro', () => {
  const rotulo = regla('.dashboard-indicador h2')
  const valor = regla('.dashboard-indicador__valor')

  expect(rotulo).toMatch(/text-transform:\s*uppercase/)
  expect(rotulo).toMatch(/font-size:\s*12px/)
  expect(valor).toMatch(/font-size:\s*26px/)
  expect(valor).toMatch(/color:\s*var\(--brand-dark\)/)
})

it('Given una tarjeta, When se estiliza, Then tiene el radio de 14px y el relleno del prototipo', () => {
  expect(regla('.dashboard-tarjeta')).toMatch(/border-radius:\s*14px/)
})
