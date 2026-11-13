/// <reference types="node" />
import { readFileSync } from 'node:fs'
import { expect, it } from 'vitest'

const css = readFileSync('src/features/backlog/backlog.css', 'utf8')
const token = (nombre: string) =>
  readFileSync('src/App.css', 'utf8').match(new RegExp(`${nombre}:\\s*(#[0-9a-f]{6})`))?.[1] ?? ''
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
const color = (bloque: string, propiedad: string) =>
  [...bloque.matchAll(new RegExp(`(?:^|[\\s;])${propiedad}:\\s*(#[0-9a-f]{6})`, 'g'))].at(-1)?.[1] ?? ''

it('Given la tarjeta de una historia, When se estiliza, Then tiene radio 14, relleno 18/20 y borde izquierdo de 4px', () => {
  const bloque = regla('.backlog-item')

  expect(bloque).toMatch(/border-radius:\s*14px/)
  expect(bloque).toMatch(/padding:\s*18px 20px/)
  expect(bloque).toMatch(/border-left:\s*4px solid/)
})

it('Given las cuatro prioridades, When se estilizan, Then cada una tiene su color de borde: Must tinta, Should marca, Could gris, Won\'t gris claro', () => {
  expect(regla('.backlog-item--must')).toMatch(/border-left-color:\s*var\(--ink\)/)
  expect(regla('.backlog-item--should')).toMatch(/border-left-color:\s*var\(--brand\)/)
  expect(regla('.backlog-item--could')).toMatch(/border-left-color:\s*#94a3b8/)
  expect(regla('.backlog-item--wont')).toMatch(/border-left-color:\s*#cbd5e1/)
})

it('Given las etiquetas, When se estilizan, Then son rectángulos de 7px de 12px de letra', () => {
  const bloque = regla('.backlog-chip')

  expect(bloque).toMatch(/border-radius:\s*7px/)
  expect(bloque).toMatch(/font-size:\s*12px/)
})

it('Given cada etiqueta con texto, When se estiliza, Then su texto cumple 4,5:1 contra su fondo', () => {
  const pares: [string, string][] = [
    ['.backlog-chip--must', '#0f1b2d'],
    ['.backlog-chip--should', ''],
    ['.backlog-chip--could', ''],
    ['.backlog-chip--ninguno', '#ffffff'],
    ['.backlog-chip--wont', '#ffffff'],
    ['.backlog-chip--sp', '#ffffff'],
    ['.backlog-chip--nueva', ''],
  ]
  for (const [selector, fondoFijo] of pares) {
    const bloque = regla(selector)
    const fondo = fondoFijo || color(bloque, 'background') || token('--brand-soft')
    expect(color(bloque, 'color'), selector).not.toBe('')
    expect(contraste(color(bloque, 'color'), fondo), selector).toBeGreaterThanOrEqual(4.5)
  }
})

it('Given la etiqueta de Won\'t have, When se estiliza, Then va tachada', () => {
  expect(regla('.backlog-chip--wont')).toMatch(/text-decoration:\s*line-through/)
})

it('Given las etiquetas "ninguno", When se estilizan, Then van en cursiva', () => {
  expect(regla('.backlog-chip--ninguno')).toMatch(/font-style:\s*italic/)
})

it('Given los botones de la tarjeta, When se estilizan, Then miden 32px de alto con letra de 13px', () => {
  const bloque = regla('.backlog-acciones .button')

  expect(bloque).toMatch(/height:\s*32px/)
  expect(bloque).toMatch(/font-size:\s*13px/)
})

it('Given "Votar", When se estiliza, Then es de contorno con borde de la marca y letra verde oscura', () => {
  const bloque = regla('.backlog-acciones .backlog-votar')

  expect(bloque).toMatch(/border-color:\s*var\(--brand\)/)
  expect(bloque).toMatch(/color:\s*var\(--brand-dark\)/)
})

it('Given los criterios de aceptación, When se estilizan, Then los separa una línea punteada y "sin criterios" va en cursiva gris', () => {
  expect(regla('.backlog-item__criterios')).toMatch(/border-top:\s*1px dashed/)
  expect(regla('.backlog-item__criterios .backlog-item__vacio')).toMatch(/font-style:\s*italic/)
})

it('Given un celular de 640px o menos, When se dibuja la tarjeta y la cabecera, Then ambas pasan a columna', () => {
  const movil = css.match(/@media \(max-width: 640px\)\s*\{([\s\S]*)\}\s*$/m)?.[1] ?? ''

  expect(movil).toMatch(/\.backlog-item__cabecera\s*\{[^}]*flex-direction:\s*column/)
  expect(movil).toMatch(/\.backlog-encabezado\s*\{[^}]*flex-direction:\s*column/)
})

it('Given la pantalla de Historias, When se estiliza, Then separa cabecera, aviso y lista 24px como el contenido del prototipo', () => {
  const bloque = regla('.backlog')

  expect(bloque).toMatch(/display:\s*grid/)
  expect(bloque).toMatch(/gap:\s*24px/)
})

const app = readFileSync('src/App.css', 'utf8')
const reglaApp = (selector: string) =>
  [...app.matchAll(/([^{}]+)\{([^}]*)\}/g)]
    .filter(([, lista]) => lista.split(',').some((parte) => parte.trim() === selector))
    .map(([, , cuerpo]) => cuerpo)
    .join('\n')

it('Given el cajón lateral, When se estiliza, Then es una columna de 480px pegada a la derecha, de alto completo y sin bordes redondeados', () => {
  expect(reglaApp('.ventana-fondo--lateral')).toMatch(/justify-items:\s*end/)
  expect(reglaApp('.ventana-fondo--lateral')).toMatch(/padding:\s*0/)
  const caja = reglaApp('.ventana--lateral')
  expect(caja).toMatch(/width:\s*min\(480px, 100%\)/)
  expect(caja).toMatch(/height:\s*100%/)
  expect(caja).toMatch(/border-radius:\s*0/)
  expect(caja).toMatch(/padding:\s*26px/)
})

it('Given el resumen «Historia creada», When se estiliza, Then es una lista de datos en grilla con los rótulos en mayúsculas de 12px', () => {
  expect(regla('.backlog-resumen')).toMatch(/display:\s*grid/)
  const rotulo = regla('.backlog-resumen dt')
  expect(rotulo).toMatch(/font-size:\s*12px/)
  expect(rotulo).toMatch(/text-transform:\s*uppercase/)
})

it('Given prioridad y valor de negocio, When se estilizan en el cajón, Then van en dos columnas', () => {
  expect(regla('.backlog-form__fila')).toMatch(/grid-template-columns:\s*1fr 1fr/)
})

it('Given los campos del cajón, When se estilizan, Then miden como los del prototipo: letra 14,5px, alto 42px, descripción de 92px y rótulo de 13px', () => {
  const campo = regla('.backlog-crear .field input')
  expect(campo).toMatch(/font-size:\s*14\.5px/)
  expect(campo).toMatch(/min-height:\s*42px/)
  expect(campo).toMatch(/padding:\s*10px 12px/)
  expect(regla('.backlog-crear .field textarea')).toMatch(/min-height:\s*92px/)
  expect(regla('.backlog-crear .field label')).toMatch(/font-size:\s*13px/)
})

it('Given el cajón, When se estiliza, Then el título es de 22px y los botones de 40px', () => {
  expect(regla('.backlog-crear__cabecera h2')).toMatch(/font-size:\s*22px/)
  expect(regla('.backlog-crear__acciones .button')).toMatch(/height:\s*40px/)
})

it('Given los criterios del resumen en Markdown, When se muestran, Then sus párrafos no suman margen', () => {
  expect(regla('.backlog-resumen .markdown p')).toMatch(/margin:\s*0/)
})

it('Given el cajón, When se estiliza, Then separa la cabecera del formulario 22px y «Cancelar» va en gris como en el prototipo', () => {
  expect(regla('.backlog-crear')).toMatch(/gap:\s*22px/)
  expect(regla('.backlog-crear__acciones .button--quiet')).toMatch(/color:\s*var\(--muted\)/)
})
