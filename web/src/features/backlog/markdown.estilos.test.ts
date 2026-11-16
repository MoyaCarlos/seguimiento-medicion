/// <reference types="node" />
import { readFileSync } from 'node:fs'
import { expect, it } from 'vitest'

const css = readFileSync('src/features/backlog/backlog.css', 'utf8')
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
const reglasDelEditor = [...css.matchAll(/([^{}]+)\{([^}]*)\}/g)]
  .filter(([, lista]) => lista.includes('.campo-markdown'))
  .map(([, , cuerpo]) => cuerpo)
  .join('\n')

it('Given el editor de Markdown, When se estiliza, Then usa la paleta teal del prototipo y no la lila', () => {
  expect(reglasDelEditor).not.toMatch(/--lilac/)
  expect(regla(".campo-markdown-pestanas button[aria-selected='true']")).toMatch(/border-color:\s*var\(--brand\)/)
  expect(regla('.campo-markdown-herramientas button:hover:not(:disabled)')).toMatch(/background:\s*var\(--brand-soft\)/)
})

it('Given las pestañas, When se estilizan, Then tienen una línea inferior y el foco por teclado es un anillo teal de 3px', () => {
  expect(regla('.campo-markdown-pestanas')).toMatch(/border-bottom:\s*1px solid var\(--line\)/)
  const foco = regla('.campo-markdown-pestanas button:focus-visible')
  expect(foco).toMatch(/outline:\s*3px solid rgb\(13 148 136 \/ 45%\)/)
  expect(regla('.campo-markdown-herramientas button:focus-visible')).toMatch(/outline:\s*3px solid rgb\(13 148 136 \/ 45%\)/)
})

it('Given la barra de formato, When se estiliza, Then sus botones miden 32px de alto con letra de 13px', () => {
  const boton = regla('.campo-markdown-herramientas button')
  expect(boton).toMatch(/min-height:\s*32px/)
  expect(boton).toMatch(/font-size:\s*13px/)
})

it('Given el marcador «Opcional», When se estiliza, Then es un texto gris de 12,5px junto a la etiqueta con contraste 4,5:1', () => {
  const marcador = regla('.campo-markdown-opcional')
  expect(marcador).toMatch(/font-size:\s*12\.5px/)
  expect(contraste(color(marcador, 'color'), '#ffffff')).toBeGreaterThanOrEqual(4.5)
  expect(regla('.campo-markdown-etiqueta')).toMatch(/align-items:\s*baseline/)
})

it('Given la vista previa, When se estiliza, Then los títulos van de 1,25 a 1,1 rem', () => {
  expect(regla('.campo-markdown-vista .markdown h1')).toMatch(/font-size:\s*1\.25rem/)
  expect(regla('.campo-markdown-vista .markdown h2')).toMatch(/font-size:\s*1\.15rem/)
  expect(regla('.campo-markdown-vista .markdown h3')).toMatch(/font-size:\s*1\.1rem/)
})

it('Given la vista previa, When tiene listas, código, citas, enlaces e imágenes, Then se ven como en react-markdown', () => {
  expect(regla('.campo-markdown-vista .markdown ul')).toMatch(/padding-left:\s*1\.4rem/)
  expect(regla('.campo-markdown-vista .markdown pre')).toMatch(/background:\s*#0f1b2d/)
  expect(contraste('#e6fffb', '#0f1b2d')).toBeGreaterThanOrEqual(4.5)
  expect(regla('.campo-markdown-vista .markdown blockquote')).toMatch(/border-left:\s*3px solid var\(--line-strong\)/)
  expect(regla('.campo-markdown-vista .markdown a')).toMatch(/color:\s*var\(--brand-dark\)/)
  expect(regla('.campo-markdown-vista .markdown img')).toMatch(/max-height:\s*20rem/)
  expect(regla('.campo-markdown-vista .markdown img')).toMatch(/max-width:\s*100%/)
})

it('Given una descripción deshabilitada, When se estiliza, Then el área de texto se ve gris y con cursor de no permitido', () => {
  const bloque = regla('.campo-markdown textarea:disabled')
  expect(bloque).toMatch(/background:\s*#eef2f6/)
  expect(bloque).toMatch(/cursor:\s*not-allowed/)
})
