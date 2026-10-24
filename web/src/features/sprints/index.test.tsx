import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Route, Routes } from 'react-router'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { Component } from './index'
import type { Sprint } from './sprint'

const respuesta = (status: number, cuerpo: unknown) =>
  new Response(JSON.stringify(cuerpo), { status })

const pendiente = (id: number): Sprint => ({ id, proyecto_id: 42, estado: 'Pendiente' })

function abrir(url = '/projects/42/sprints') {
  render(
    <MemoryRouter initialEntries={[url]}>
      <Routes>
        <Route path="/projects/:id/sprints" element={<Component />} />
      </Routes>
    </MemoryRouter>,
  )
}

const seccion = (nombre: string) =>
  screen.getByRole('heading', { name: nombre }).closest('section') as HTMLElement

const crear = () => userEvent.click(screen.getByRole('button', { name: 'Crear Sprint' }))

afterEach(() => vi.unstubAllGlobals())

describe('pantalla de Sprints', () => {
  it('conserva el contrato con la base: h1 "Sprints" y el id en un único <strong>', () => {
    abrir()

    expect(screen.getByRole('heading', { level: 1, name: 'Sprints' })).toBeTruthy()
    expect(screen.getAllByRole('heading', { name: 'Sprints' })).toHaveLength(1)
    expect(screen.getByText('42').tagName).toBe('STRONG')
    expect(screen.getByText('42').parentElement?.textContent).toBe('Proyecto actual: 42')
  })

  it('muestra las tres secciones vacías y el texto de la limitación', () => {
    abrir()

    for (const nombre of ['Pendiente', 'Activo', 'Finalizado']) {
      expect(within(seccion(nombre)).getByText('No hay Sprints')).toBeTruthy()
    }
    expect(screen.getByText(/Los Sprints existentes no se listan todavía/).textContent).toContain(
      'docs/curl-frontend-sprint-2.md',
    )
  })

  it('crea un Sprint con el id del proyecto y lo muestra en Pendiente sin Goal ni fechas', async () => {
    const fetchSimulado = vi.fn().mockResolvedValue(respuesta(201, pendiente(5)))
    vi.stubGlobal('fetch', fetchSimulado)
    abrir()

    await crear()

    const tarjeta = await within(seccion('Pendiente')).findByRole('heading', { name: 'Sprint #5' })
    expect(tarjeta).toBeTruthy()
    const [ruta, opciones] = fetchSimulado.mock.calls[0]
    expect(ruta).toBe('/sprints')
    expect(JSON.parse(opciones.body)).toEqual({ proyecto_id: 42 })
    expect(within(seccion('Pendiente')).queryByText(/\d{4}-\d{2}-\d{2}/)).toBeNull()
  })

  it('deshabilita "Crear Sprint" y muestra la carga mientras crea', async () => {
    vi.stubGlobal('fetch', vi.fn().mockReturnValue(new Promise(() => {})))
    abrir()

    await crear()

    expect(screen.getByText('Cargando…')).toBeTruthy()
    expect((screen.getByRole('button', { name: 'Crear Sprint' }) as HTMLButtonElement).disabled).toBe(true)
  })

  it('muestra el error de crear en el aviso de la pantalla y no agrega ningún Sprint', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue(respuesta(400, { campo: 'estado', mensaje: 'ya hay un Sprint pendiente' })),
    )
    abrir()

    await crear()

    expect((await screen.findByRole('alert')).textContent).toBe('ya hay un Sprint pendiente')
    expect(within(seccion('Pendiente')).getByText('No hay Sprints')).toBeTruthy()
  })

  it('recorre el ciclo: crear, iniciar y cerrar', async () => {
    const activo: Sprint = {
      ...pendiente(5),
      estado: 'Activo',
      sprint_goal: 'Meta',
      fecha_inicio: '2026-10-01',
      fecha_fin: '2026-10-15',
    }
    const respuestas: Record<string, Sprint> = {
      '/sprints': pendiente(5),
      '/sprints/5/iniciar': activo,
      '/sprints/5/cerrar': { ...activo, estado: 'Finalizado' },
    }
    vi.stubGlobal('fetch', vi.fn((ruta: string) => Promise.resolve(respuesta(200, respuestas[ruta]))))
    const user = userEvent.setup()
    abrir()

    await crear()
    await screen.findByRole('heading', { name: 'Sprint #5' })
    await user.type(screen.getByLabelText('Sprint Goal'), 'Meta')
    await user.type(screen.getByLabelText('Fecha de inicio'), '2026-10-01')
    await user.type(screen.getByLabelText('Fecha de fin'), '2026-10-15')
    await user.click(screen.getByRole('button', { name: 'Iniciar Sprint' }))

    const enActivo = within(seccion('Activo'))
    expect(await enActivo.findByRole('heading', { name: 'Sprint #5' })).toBeTruthy()
    expect(enActivo.getByText('Meta')).toBeTruthy()
    expect(enActivo.getByText('2026-10-01 a 2026-10-15')).toBeTruthy()
    expect(within(seccion('Pendiente')).queryByRole('heading', { name: 'Sprint #5' })).toBeNull()
    expect(screen.queryByRole('button', { name: 'Iniciar Sprint' })).toBeNull()

    await user.click(screen.getByRole('button', { name: 'Cerrar Sprint' }))

    const enFinalizado = within(seccion('Finalizado'))
    expect(await enFinalizado.findByRole('heading', { name: 'Sprint #5' })).toBeTruthy()
    expect(within(seccion('Activo')).queryByRole('heading', { name: 'Sprint #5' })).toBeNull()
    expect(screen.queryByRole('button', { name: 'Cerrar Sprint' })).toBeNull()
    expect(screen.queryByRole('button', { name: 'Iniciar Sprint' })).toBeNull()
  })

  it('el error de un Sprint no aparece en otro', async () => {
    const respuestas = [respuesta(201, pendiente(1)), respuesta(201, pendiente(2))]
    vi.stubGlobal(
      'fetch',
      vi.fn((ruta: string) =>
        Promise.resolve(
          ruta === '/sprints/1/iniciar'
            ? respuesta(400, { campo: 'estado', mensaje: 'error del Sprint 1' })
            : respuestas.shift(),
        ),
      ),
    )
    const user = userEvent.setup()
    abrir()

    await crear()
    await screen.findByRole('heading', { name: 'Sprint #1' })
    await crear()
    await screen.findByRole('heading', { name: 'Sprint #2' })
    const [primera, segunda] = screen.getAllByRole('article')
    await user.click(within(primera).getByRole('button', { name: 'Iniciar Sprint' }))

    expect((await within(primera).findByRole('alert')).textContent).toBe('error del Sprint 1')
    expect(within(segunda).queryByRole('alert')).toBeNull()
    await waitFor(() => expect(screen.getAllByRole('alert')).toHaveLength(1))
  })

  it('no valida el id de la URL: muestra el error que devuelve la API', async () => {
    const fetchSimulado = vi
      .fn()
      .mockResolvedValue(respuesta(400, { campo: 'proyecto_id', mensaje: 'proyecto_id inválido' }))
    vi.stubGlobal('fetch', fetchSimulado)
    abrir('/projects/abc/sprints')

    await crear()

    expect((await screen.findByRole('alert')).textContent).toBe('proyecto_id inválido')
    expect(fetchSimulado).toHaveBeenCalledTimes(1)
  })
})
