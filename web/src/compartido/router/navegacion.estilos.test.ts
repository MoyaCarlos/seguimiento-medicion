/// <reference types="node" />
import { readFileSync } from 'node:fs'
import { expect, it } from 'vitest'

const estilos = readFileSync('src/App.css', 'utf8')
const estilosAngostos = estilos.split('@media (max-width: 900px)')[1] ?? ''

it('Given una pantalla angosta, When se muestra el armazón, Then la barra lateral no ocupa ancho del contenido', () => {
  const armazon = estilosAngostos.match(/\.armazon\s*\{[^}]*\}/)?.[0] ?? ''
  const lateral = estilosAngostos.match(/\.barra-lateral\s*\{[^}]*\}/)?.[0] ?? ''

  expect(armazon).toMatch(/grid-template-columns:\s*minmax\(0,\s*1fr\)/)
  expect(lateral).toMatch(/position:\s*fixed/)
})
