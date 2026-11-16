import { render } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { AvatarProyecto } from './AvatarProyecto'

const avatar = (props: Parameters<typeof AvatarProyecto>[0]) =>
  render(<AvatarProyecto {...props} />).container.firstElementChild as HTMLElement

describe('AvatarProyecto', () => {
  it('muestra la inicial del nombre en mayúscula y se oculta a los lectores de pantalla', () => {
    const el = avatar({ id: 1, nombre: ' beta' })

    expect(el.getAttribute('data-inicial')).toBe('B')
    expect(el.getAttribute('aria-hidden')).toBe('true')
  })

  it('sin nombre usa «?»', () => {
    expect(avatar({ id: 1, nombre: '   ' }).getAttribute('data-inicial')).toBe('?')
  })

  it.each([[0, 0], [1, 1], [4, 4], [5, 0], [7, 2], ['12', 2]])('el proyecto %s usa la paleta %s', (id, paleta) => {
    expect(avatar({ id, nombre: 'A' }).className).toContain(`avatar-proyecto--${paleta}`)
  })

  it('chico agrega la variante pequeña', () => {
    expect(avatar({ id: 1, nombre: 'A', chico: true }).className).toContain('avatar-proyecto--chico')
    expect(avatar({ id: 1, nombre: 'A' }).className).not.toContain('avatar-proyecto--chico')
  })
})
