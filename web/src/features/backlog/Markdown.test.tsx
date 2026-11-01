import { render } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { Markdown } from './Markdown'

function renderizar(texto: string) {
  return render(<Markdown texto={texto} />).container
}

describe('Markdown', () => {
  it('Escenario 1: muestra el texto plano sin romperlo', () => {
    const contenedor = renderizar('Autenticación con email y contraseña')

    expect(contenedor.textContent).toContain('Autenticación con email y contraseña')
  })

  it('Escenario 2: interpreta negritas, listas, código y enlaces comunes', () => {
    const contenedor = renderizar(
      [
        '**Importante:** validar el alta.',
        '',
        '- primer paso',
        '- segundo paso',
        '',
        '`npm test`',
        '',
        '[Documentación](https://example.com/docs)',
      ].join('\n'),
    )

    expect(contenedor.querySelector('strong')?.textContent).toBe('Importante:')
    expect(contenedor.querySelectorAll('li')).toHaveLength(2)
    expect(contenedor.querySelector('code')?.textContent).toBe('npm test')
    expect(contenedor.querySelector('a')?.getAttribute('href')).toBe('https://example.com/docs')
  })

  it('Escenario 3: no inserta ni ejecuta HTML crudo de un intento de XSS', () => {
    const alerta = vi.spyOn(window, 'alert').mockImplementation(() => {})
    const contenedor = renderizar(
      ['<script>alert(1)</script>', '', '<img src=x onerror="alert(2)">', '', 'Texto seguro'].join(
        '\n',
      ),
    )

    expect(contenedor.querySelector('script')).toBeNull()
    expect(contenedor.querySelectorAll('img')).toHaveLength(0)
    expect(contenedor.querySelector('[onerror]')).toBeNull()
    expect(contenedor.textContent).toContain('Texto seguro')
    expect(alerta).not.toHaveBeenCalled()

    alerta.mockRestore()
  })

  it('Escenario 6: no deja un enlace javascript: ejecutable', () => {
    const contenedor = renderizar('[peligroso](javascript:alert(1))')
    const enlace = contenedor.querySelector('a')

    expect(enlace?.getAttribute('href') ?? '').not.toMatch(/^javascript:/i)
  })
})
