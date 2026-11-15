import { cleanup, render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Route, Routes } from 'react-router'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { borrarToken, guardarToken } from './almacenamientoToken'
import { BloqueUsuario } from './BloqueUsuario'
import { SesionProvider } from './SesionContext'

const USUARIO = { id: 7, nombre: 'ana.perez' }
const json = (cuerpo: unknown, status: number) => new Response(JSON.stringify(cuerpo), { status })

function montar(rol?: string) {
  return render(
    <MemoryRouter initialEntries={['/']}>
      <SesionProvider>
        <BloqueUsuario rol={rol} />
        <Routes>
          <Route path="/ingresar" element={<h1>Pantalla de ingreso</h1>} />
          <Route path="*" element={<p>Inicio</p>} />
        </Routes>
      </SesionProvider>
    </MemoryRouter>,
  )
}

beforeEach(() => {
  borrarToken()
  localStorage.clear()
})

afterEach(() => {
  cleanup()
  vi.unstubAllGlobals()
})

describe('BloqueUsuario (pie de la barra lateral)', () => {
  it('Given una sesión iniciada, When se muestra, Then trae el avatar con la inicial, el nombre y «Sesión iniciada»', async () => {
    guardarToken('token')
    vi.stubGlobal('fetch', vi.fn(async () => json(USUARIO, 200)))
    montar()

    const bloque = (await screen.findByText('ana.perez')).closest('.bloque-usuario') as HTMLElement

    expect(bloque.querySelector('.bloque-usuario__avatar')?.textContent).toBe('A')
    expect(bloque.querySelector('.bloque-usuario__avatar')?.getAttribute('aria-hidden')).toBe('true')
    expect(bloque.querySelector('small')?.textContent).toBe('Sesión iniciada')
  })

  it('Given una pantalla de proyecto, When se muestra, Then el rol reemplaza a «Sesión iniciada»', async () => {
    guardarToken('token')
    vi.stubGlobal('fetch', vi.fn(async () => json(USUARIO, 200)))
    montar('Product Builder')

    const bloque = (await screen.findByText('ana.perez')).closest('.bloque-usuario') as HTMLElement

    expect(bloque.querySelector('small')?.textContent).toBe('Product Builder')
    expect(screen.queryByText('Sesión iniciada')).toBeNull()
  })

  it('Given el botón de salida, When se pulsa, Then cierra la sesión y lleva a /ingresar', async () => {
    guardarToken('token')
    const fetchSimulado = vi.fn(async (ruta: string) => (ruta === '/auth/salir' ? json({}, 200) : json(USUARIO, 200)))
    vi.stubGlobal('fetch', fetchSimulado)
    montar()

    await userEvent.click(await screen.findByRole('button', { name: 'Cerrar sesión' }))

    expect(await screen.findByRole('heading', { name: 'Pantalla de ingreso' })).toBeTruthy()
    expect(fetchSimulado.mock.calls.some(([ruta]) => ruta === '/auth/salir')).toBe(true)
  })

  it('Given que no hay sesión, When se muestra, Then no dibuja nada', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => json({}, 401)))
    const { container } = montar()

    expect(container.querySelector('.bloque-usuario')).toBeNull()
  })
})
