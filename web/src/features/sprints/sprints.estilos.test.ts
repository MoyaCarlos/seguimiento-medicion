/// <reference types="node" />
import { readFileSync } from 'node:fs'
import { expect, it } from 'vitest'

const css = readFileSync('src/features/sprints/sprints.css', 'utf8')
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

it('Given el tablero, When se estiliza, Then son tres columnas de igual ancho con 16px de hueco', () => {
  const bloque = regla('.sprints-tablero')

  expect(bloque).toMatch(/grid-template-columns:\s*repeat\(3, minmax\(0, 1fr\)\)/)
  expect(bloque).toMatch(/gap:\s*16px/)
  expect(bloque).toMatch(/align-items:\s*start/)
})

it('Given el ancho de tablet o menos, When se estiliza, Then las columnas se apilan', () => {
  expect(css).toMatch(/@media \(max-width: 1100px\)\s*\{[^@]*\.sprints-tablero\s*\{[^}]*grid-template-columns:\s*1fr/)
})

it('Given una columna, When se estiliza, Then es una tarjeta gris de radio 16 y relleno 12', () => {
  const bloque = regla('.sprints-carril')

  expect(bloque).toMatch(/border-radius:\s*16px/)
  expect(bloque).toMatch(/padding:\s*12px/)
  expect(bloque).toMatch(/background:\s*rgb\(15 27 45 \/ 3\.5%\)/)
})

it('Given el título de cada columna, When se estiliza, Then es una etiqueta de 12px en mayúsculas con un punto de color por estado', () => {
  const titulo = regla('.sprints-carril h2')

  expect(titulo).toMatch(/font-size:\s*12px/)
  expect(titulo).toMatch(/text-transform:\s*uppercase/)
  expect(regla('.sprints-carril h2::before')).toMatch(/border-radius:\s*50%/)
  expect(regla('.sprints-carril h2::before')).toMatch(/background:\s*#94a3b8/)
  expect(regla('.sprints-carril--activo h2::before')).toMatch(/background:\s*var\(--brand\)/)
  expect(regla('.sprints-carril--fin h2::before')).toMatch(/background:\s*#15803d/)
})

it('Given el título de cada columna, When se estiliza, Then su texto cumple 4,5:1 contra el fondo de la columna', () => {
  for (const [selector, fondo] of [
    ['.sprints-carril h2', '#f0f1f3'],
    ['.sprints-carril--activo h2', '#f0f1f3'],
    ['.sprints-carril--fin h2', '#f0f1f3'],
  ]) {
    const texto = color(regla(selector), 'color')

    expect(texto, selector).not.toBe('')
    expect(contraste(texto, fondo), selector).toBeGreaterThanOrEqual(4.5)
  }
})

it('Given la tarjeta de un Sprint, When se estiliza, Then tiene radio 14 y relleno 16', () => {
  const bloque = regla('.sprint-tarjeta')

  expect(bloque).toMatch(/border-radius:\s*14px/)
  expect(bloque).toMatch(/padding:\s*16px/)
})

it('Given las fechas del Sprint, When se estilizan, Then son una etiqueta de 12px con cifras tabulares', () => {
  const bloque = regla('.sprint-fechas')

  expect(bloque).toMatch(/border-radius:\s*7px/)
  expect(bloque).toMatch(/font-size:\s*12px/)
  expect(bloque).toMatch(/font-variant-numeric:\s*tabular-nums/)
  expect(regla('.sprint-fechas svg')).toMatch(/width:\s*13px/)
})

it('Given la columna vacía, When se estiliza, Then es un recuadro punteado centrado', () => {
  const bloque = regla('.sprints-vacio')

  expect(bloque).toMatch(/border:\s*1px dashed/)
  expect(bloque).toMatch(/text-align:\s*center/)
})

it('Given el plazo del proyecto, When se estiliza, Then es una línea de 13,5px', () => {
  expect(regla('.sprints-plazo')).toMatch(/font-size:\s*13\.5px/)
})

it('Given una historia del Sprint, When se estiliza, Then es una tarjeta de radio 10 con relleno 10/12 y fondo casi blanco', () => {
  const bloque = regla('.sprint-historia')

  expect(bloque).toMatch(/border-radius:\s*10px/)
  expect(bloque).toMatch(/padding:\s*10px 12px/)
  expect(bloque).toMatch(/background:\s*#fbfcfd/)
})

it('Given el título de la historia, When se estiliza, Then es de 14px en negrita', () => {
  const bloque = regla('.sprint-historia__titulo')

  expect(bloque).toMatch(/font-size:\s*14px/)
  expect(bloque).toMatch(/font-weight:\s*700/)
})

it('Given las etiquetas de la historia, When se estilizan, Then son rectángulos de 7px de 12px de letra y su texto cumple 4,5:1', () => {
  expect(regla('.sprint-chip')).toMatch(/border-radius:\s*7px/)
  expect(regla('.sprint-chip')).toMatch(/font-size:\s*12px/)
  const pares: [string, string][] = [
    ['.sprint-chip--must', '#0f1b2d'],
    ['.sprint-chip--should', '#dcf3ef'],
    ['.sprint-chip--could', '#eef2f6'],
    ['.sprint-chip--nueva', '#eef2f6'],
    ['.sprint-chip--sp', '#ffffff'],
    ['.sprint-chip--ninguno', '#ffffff'],
    ['.sprint-chip--wont', '#ffffff'],
    ['.sprint-chip--progreso', '#fdf0d5'],
    ['.sprint-chip--ok', '#dcf5e4'],
  ]
  for (const [selector, fondoFijo] of pares) {
    const bloque = regla(selector)
    const texto = color(bloque, 'color')
    const fondo = color(bloque, 'background') || fondoFijo

    expect(texto, selector).not.toBe('')
    expect(contraste(texto, fondo), selector).toBeGreaterThanOrEqual(4.5)
  }
  expect(regla('.sprint-chip--ninguno')).toMatch(/font-style:\s*italic/)
  expect(regla('.sprint-chip--wont')).toMatch(/text-decoration:\s*line-through/)
})

it('Given el esfuerzo real, When se estiliza, Then es de 13px, negrita y verde de marca con cifras tabulares', () => {
  const bloque = regla('.sprint-esfuerzo')

  expect(bloque).toMatch(/font-size:\s*13px/)
  expect(bloque).toMatch(/font-weight:\s*700/)
  expect(bloque).toMatch(/font-variant-numeric:\s*tabular-nums/)
  expect(contraste(color(bloque, 'color'), '#fbfcfd')).toBeGreaterThanOrEqual(4.5)
})

it('Given el pie de acciones, When se estiliza, Then los botones van a la derecha y pueden pasar de renglón', () => {
  const bloque = regla('.sprint-historia__acciones')

  expect(bloque).toMatch(/display:\s*flex/)
  expect(bloque).toMatch(/justify-content:\s*flex-end/)
  expect(bloque).toMatch(/flex-wrap:\s*wrap/)
  expect(bloque).toMatch(/gap:\s*6px/)
})

it('Given los botones chicos, When se estilizan, Then miden 32px con letra de 13px; el de contorno usa el color de marca', () => {
  expect(regla('.sprint-boton')).toMatch(/min-height:\s*32px/)
  expect(regla('.sprint-boton')).toMatch(/font-size:\s*13px/)
  const contorno = regla('.sprint-boton--contorno')
  expect(contorno).toMatch(/border-color:\s*var\(--brand\)/)
  expect(contraste(color(contorno, 'color'), '#ffffff')).toBeGreaterThanOrEqual(4.5)
})

it('Given «Ver detalle», When se estiliza, Then es un control verde de marca, en negrita y con un triángulo delante', () => {
  expect(regla('.sprint-detalle')).toMatch(/font-weight:\s*700/)
  expect(regla('.sprint-detalle::before')).toMatch(/content:\s*'▸'/)
  expect(regla('.sprint-detalle[aria-expanded="true"]::before')).toMatch(/rotate\(90deg\)/)
})

it('Given el panel «Iniciar Sprint», When se estiliza, Then es un recuadro gris claro de radio 12 y relleno 18', () => {
  const bloque = regla('.sprint-panel')

  expect(bloque).toMatch(/background:\s*#f8fafc/)
  expect(bloque).toMatch(/border-radius:\s*12px/)
  expect(bloque).toMatch(/padding:\s*18px/)
  expect(regla('.sprint-boton-bloque')).toMatch(/width:\s*100%/)
})

it('Given el pie de un Sprint activo, When se estiliza, Then lo separa una línea y el botón queda a la derecha', () => {
  const bloque = regla('.sprint-pie')

  expect(bloque).toMatch(/border-top:\s*1px solid/)
  expect(bloque).toMatch(/justify-content:\s*flex-end/)
})

it('Given un Sprint finalizado, When se estiliza, Then es un desplegable de radio 14 con la línea inferior al abrir', () => {
  expect(regla('.sprint-terminado')).toMatch(/border-radius:\s*14px/)
  expect(regla('.sprint-terminado[open] summary')).toMatch(/border-bottom:\s*1px solid/)
  expect(regla('.sprint-terminado summary')).toMatch(/cursor:\s*pointer/)
})

it('Given el resumen del finalizado, When se abre, Then el chevrón gira 180 grados', () => {
  expect(regla('.sprint-terminado[open] .sprint-chevron')).toMatch(/rotate\(180deg\)/)
})

it('Given las historias dentro de un finalizado, When se estilizan, Then van en lista con separadores y sin tarjeta', () => {
  const bloque = regla('.sprint-terminado .sprint-historia')

  expect(bloque).toMatch(/border:\s*0/)
  expect(bloque).toMatch(/background:\s*none/)
  expect(regla('.sprint-terminado .sprint-backlog > li + li')).toMatch(/border-top:\s*1px solid/)
})

it('Given el encabezado, When lleva ícono, Then el ícono se separa 8px del texto', () => {
  expect(regla('.sprints-encabezado .button')).toMatch(/gap:\s*8px/)
})

it('Given una columna con formularios anchos, When se estiliza, Then su única columna puede encogerse y no empuja a la tarjeta fuera de la columna', () => {
  expect(regla('.sprints-carril')).toMatch(/grid-template-columns:\s*minmax\(0, 1fr\)/)
})

it('Given una pantalla angosta, When se estiliza, Then «Crear Sprint» no se estira a todo el ancho', () => {
  expect(regla('.sprints-encabezado .button')).toMatch(/align-self:\s*flex-start/)
})

it('Given una pantalla de 640px o menos, When se estiliza, Then «Completar»/«Iniciar» ocupa todo el ancho con 40px de alto', () => {
  const movil = css.slice(css.indexOf('@media (max-width: 640px)'))

  expect(movil).toMatch(/\.sprint-boton--contorno\s*\{[^}]*width:\s*100%/)
  expect(movil).toMatch(/\.sprint-boton--contorno\s*\{[^}]*min-height:\s*40px/)
  expect(movil).toMatch(/\.sprint-historia__acciones\s*\{[^}]*justify-content:\s*flex-start/)
})
