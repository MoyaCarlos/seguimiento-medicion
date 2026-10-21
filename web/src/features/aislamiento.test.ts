import { describe, expect, it } from 'vitest'

const archivos = import.meta.glob<string>('./**/*.{ts,tsx}', {
  query: '?raw',
  import: 'default',
  eager: true,
})

const IMPORTACION = /(?:from|import)\s*\(?\s*['"]([^'"]+)['"]/g

function featureAjena(archivo: string, especificador: string): string | undefined {
  if (!especificador.startsWith('.')) return undefined

  const carpetas = archivo.replace('./', '').split('/').slice(0, -1)
  const destino = [...carpetas]
  for (const segmento of especificador.split('/')) {
    if (segmento === '..') {
      if (destino.length === 0) return undefined
      destino.pop()
    } else if (segmento !== '.') {
      destino.push(segmento)
    }
  }

  const feature = destino[0]
  return feature && feature !== carpetas[0] ? feature : undefined
}

describe('featureAjena', () => {
  it.each([
    ['./proyecto/index.tsx', '../backlog', 'backlog'],
    ['./proyecto/sub/pantalla.tsx', '../../backlog/index', 'backlog'],
    ['./proyecto/sub/pantalla.tsx', '../index', undefined],
    ['./proyecto/index.tsx', './sub/pantalla', undefined],
    ['./proyecto/index.tsx', '../../compartido/api/client', undefined],
    ['./proyecto/sub/pantalla.tsx', '../../../compartido/api/client', undefined],
    ['./proyecto/index.tsx', 'react-router', undefined],
  ])('%s importa %s -> %s', (archivo, especificador, esperado) => {
    expect(featureAjena(archivo, especificador)).toBe(esperado)
  })
})

describe('aislamiento entre features', () => {
  it('ninguna feature importa a otra', () => {
    const cruces = Object.entries(archivos)
      .filter(([archivo]) => archivo !== './aislamiento.test.ts')
      .flatMap(([archivo, codigo]) =>
        [...codigo.matchAll(IMPORTACION)]
          .map(([, especificador]) => ({ archivo, especificador, ajena: featureAjena(archivo, especificador) }))
          .filter(({ ajena }) => ajena !== undefined),
      )
      .map(({ archivo, especificador, ajena }) => `${archivo} importa ${especificador} (feature ${ajena})`)

    expect(cruces).toEqual([])
  })
})
