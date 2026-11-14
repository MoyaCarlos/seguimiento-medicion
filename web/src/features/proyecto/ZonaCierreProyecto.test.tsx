import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { ZonaCierreProyecto } from './ZonaCierreProyecto'

const respuestaCierre = {
  id: 7,
  nombre: 'Proyecto Alfa',
  descripcion: '',
  fecha_inicio: null,
  fecha_fin: null,
  creado_en: '2026-09-01T00:00:00Z',
  horas_por_story_point: null,
  fecha_cierre: '2026-11-11',
  estado: 'Finalizado',
}

function json(cuerpo: unknown, status = 200) {
  return new Response(JSON.stringify(cuerpo), { status })
}

function responder(respuesta: Response | (() => Promise<Response>)) {
  const fetchMock = vi.fn(() => (typeof respuesta === 'function' ? respuesta() : Promise.resolve(respuesta)))
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

function abrir(onCambio = vi.fn()) {
  render(<ZonaCierreProyecto proyectoId="7" rol="scrum_master" fechaCierre={null} onCambio={onCambio} />)
  return onCambio
}

async function pedirConfirmacion() {
  await userEvent.click(screen.getByRole('button', { name: 'Cerrar proyecto' }))
}

async function confirmar() {
  await pedirConfirmacion()
  await userEvent.click(screen.getByRole('button', { name: 'Confirmar' }))
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('zona de cierre del proyecto - cerrar', () => {
  it('Escenario: abierto y Scrum Master, «Cerrar proyecto» está habilitado y no se envía ningún pedido', () => {
    const fetchMock = responder(json(respuestaCierre))
    abrir()

    expect(screen.getByRole('button', { name: 'Cerrar proyecto' }).hasAttribute('disabled')).toBe(false)
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it('Escenario: pulsar «Cerrar proyecto» pide confirmación con «Confirmar» y «Cancelar» y el foco queda en «Cancelar»', async () => {
    responder(json(respuestaCierre))
    abrir()

    await pedirConfirmacion()

    expect(screen.getByText('¿Cerrar el proyecto?')).toBeTruthy()
    expect(screen.getByRole('button', { name: 'Confirmar' })).toBeTruthy()
    expect(document.activeElement).toBe(screen.getByRole('button', { name: 'Cancelar' }))
  })

  it('SC-002: «Cancelar» vuelve al reposo sin enviar nada', async () => {
    const fetchMock = responder(json(respuestaCierre))
    const onCambio = abrir()

    await pedirConfirmacion()
    await userEvent.click(screen.getByRole('button', { name: 'Cancelar' }))

    expect(screen.queryByText('¿Cerrar el proyecto?')).toBeNull()
    expect(screen.getByRole('button', { name: 'Cerrar proyecto' })).toBeTruthy()
    expect(fetchMock).not.toHaveBeenCalled()
    expect(onCambio).not.toHaveBeenCalled()
  })

  it('Escenario: Escape cancela la confirmación sin enviar nada', async () => {
    const fetchMock = responder(json(respuestaCierre))
    abrir()

    await pedirConfirmacion()
    await userEvent.keyboard('{Escape}')

    expect(screen.queryByText('¿Cerrar el proyecto?')).toBeNull()
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it('Escenario: «Confirmar» hace un POST /projects/7/cerrar, llama onCambio con la fecha y muestra el estado «Finalizado»', async () => {
    const fetchMock = responder(json(respuestaCierre))
    const onCambio = abrir()

    await confirmar()

    expect(await screen.findByText('Finalizado')).toBeTruthy()
    expect(fetchMock).toHaveBeenCalledTimes(1)
    const [url, init] = fetchMock.mock.calls[0] as unknown as [string, RequestInit]
    expect(url).toBe('/projects/7/cerrar')
    expect(init.method).toBe('POST')
    expect(onCambio).toHaveBeenCalledWith('2026-11-11')
  })

  it('Escenario: con el pedido en curso los botones están deshabilitados, un segundo clic no envía otro pedido y Escape no cancela', async () => {
    let resolver: (r: Response) => void = () => {}
    const fetchMock = responder(() => new Promise<Response>((resuelto) => { resolver = resuelto }))
    abrir()

    await confirmar()

    const confirmarBoton = screen.getByRole('button', { name: 'Confirmar' })
    const cancelarBoton = screen.getByRole('button', { name: 'Cancelar' })
    expect(confirmarBoton.hasAttribute('disabled')).toBe(true)
    expect(cancelarBoton.hasAttribute('disabled')).toBe(true)

    await userEvent.click(confirmarBoton)
    await userEvent.keyboard('{Escape}')

    expect(fetchMock).toHaveBeenCalledTimes(1)
    expect(screen.getByText('¿Cerrar el proyecto?')).toBeTruthy()

    resolver(json(respuestaCierre))
    expect(await screen.findByText('Finalizado')).toBeTruthy()
  })

  it('Escenario: 409 «Cerrá el Sprint activo antes de cerrar el proyecto» cierra la confirmación, muestra el mensaje tal cual y no llama onCambio', async () => {
    const mensaje = 'Cerrá el Sprint activo antes de cerrar el proyecto'
    responder(json({ mensaje }, 409))
    const onCambio = abrir()

    await confirmar()

    expect((await screen.findByRole('alert')).textContent).toBe(mensaje)
    expect(screen.queryByText('¿Cerrar el proyecto?')).toBeNull()
    expect(screen.getByRole('button', { name: 'Cerrar proyecto' })).toBeTruthy()
    expect(onCambio).not.toHaveBeenCalled()
  })

  it('Escenario: 409 «El proyecto ya está cerrado» cierra la confirmación, muestra el mensaje tal cual y no llama onCambio', async () => {
    const mensaje = 'El proyecto ya está cerrado'
    responder(json({ mensaje }, 409))
    const onCambio = abrir()

    await confirmar()

    expect((await screen.findByRole('alert')).textContent).toBe(mensaje)
    expect(screen.queryByText('¿Cerrar el proyecto?')).toBeNull()
    expect(onCambio).not.toHaveBeenCalled()
  })

  it('Escenario: 403 muestra el aviso «No tenés permiso» con el mensaje del servidor', async () => {
    responder(json({ mensaje: 'Solo el Scrum Master puede realizar esta acción' }, 403))
    const onCambio = abrir()

    await confirmar()

    expect(await screen.findByText('No tenés permiso')).toBeTruthy()
    expect(screen.getByText('Solo el Scrum Master puede realizar esta acción')).toBeTruthy()
    expect(onCambio).not.toHaveBeenCalled()
  })

  it('Escenario: 404 muestra el aviso de proyecto inexistente', async () => {
    responder(json({ mensaje: 'no existe el proyecto' }, 404))
    const onCambio = abrir()

    await confirmar()

    expect((await screen.findByRole('alert')).textContent).toBe('no existe el proyecto')
    expect(onCambio).not.toHaveBeenCalled()
  })

  it('Escenario: conexión caída muestra el error de conexión y el estado no cambia', async () => {
    const fetchMock = vi.fn().mockRejectedValue(new TypeError('Failed to fetch'))
    vi.stubGlobal('fetch', fetchMock)
    const onCambio = abrir()

    await confirmar()

    expect((await screen.findByRole('alert')).textContent).toBe('No se pudo conectar con el servidor')
    await waitFor(() => expect(screen.getByRole('button', { name: 'Cerrar proyecto' })).toBeTruthy())
    expect(screen.queryByText('Finalizado')).toBeNull()
    expect(onCambio).not.toHaveBeenCalled()
  })
})

const respuestaReapertura = { ...respuestaCierre, fecha_cierre: null, estado: 'En curso' }

function abrirCerrado(onCambio = vi.fn()) {
  render(<ZonaCierreProyecto proyectoId="7" rol="scrum_master" fechaCierre="2026-11-11" onCambio={onCambio} />)
  return onCambio
}

async function reabrir() {
  await userEvent.click(screen.getByRole('button', { name: 'Reabrir proyecto' }))
}

describe('zona de cierre del proyecto - reabrir', () => {
  it('Escenario: cerrado y Scrum Master, se ve «Proyecto cerrado el 2026-11-11» y «Reabrir proyecto» habilitado, sin «Cerrar proyecto»', () => {
    const fetchMock = responder(json(respuestaReapertura))
    abrirCerrado()

    expect(screen.getByText('Proyecto cerrado el 2026-11-11')).toBeTruthy()
    expect(screen.getByRole('button', { name: 'Reabrir proyecto' }).hasAttribute('disabled')).toBe(false)
    expect(screen.queryByRole('button', { name: 'Cerrar proyecto' })).toBeNull()
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it('US3.1: pulsar «Reabrir proyecto» hace un POST /projects/7/reabrir sin confirmación, llama onCambio(null) y muestra el estado devuelto', async () => {
    const fetchMock = responder(json(respuestaReapertura))
    const onCambio = abrirCerrado()

    await reabrir()

    expect(await screen.findByText('En curso')).toBeTruthy()
    expect(screen.queryByText('¿Cerrar el proyecto?')).toBeNull()
    expect(screen.queryByRole('button', { name: 'Confirmar' })).toBeNull()
    expect(fetchMock).toHaveBeenCalledTimes(1)
    const [url, init] = fetchMock.mock.calls[0] as unknown as [string, RequestInit]
    expect(url).toBe('/projects/7/reabrir')
    expect(init.method).toBe('POST')
    expect(onCambio).toHaveBeenCalledWith(null)
  })

  it('Escenario: con la reapertura en curso «Reabrir proyecto» está deshabilitado y un segundo clic no envía otro pedido', async () => {
    let resolver: (r: Response) => void = () => {}
    const fetchMock = responder(() => new Promise<Response>((resuelto) => { resolver = resuelto }))
    abrirCerrado()

    await reabrir()

    const boton = screen.getByRole('button', { name: 'Reabrir proyecto' })
    expect(boton.hasAttribute('disabled')).toBe(true)
    await userEvent.click(boton)
    expect(fetchMock).toHaveBeenCalledTimes(1)

    resolver(json(respuestaReapertura))
    expect(await screen.findByText('En curso')).toBeTruthy()
  })

  it('US3.2: 409 «El proyecto no está cerrado» muestra el mensaje tal cual y no llama onCambio', async () => {
    const mensaje = 'El proyecto no está cerrado'
    responder(json({ mensaje }, 409))
    const onCambio = abrirCerrado()

    await reabrir()

    expect((await screen.findByRole('alert')).textContent).toBe(mensaje)
    expect(screen.getByRole('button', { name: 'Reabrir proyecto' }).hasAttribute('disabled')).toBe(false)
    expect(onCambio).not.toHaveBeenCalled()
  })

  it('US4.2: 403 al reabrir muestra el aviso «No tenés permiso» con el mensaje del servidor', async () => {
    responder(json({ mensaje: 'Solo el Scrum Master puede realizar esta acción' }, 403))
    const onCambio = abrirCerrado()

    await reabrir()

    expect(await screen.findByText('No tenés permiso')).toBeTruthy()
    expect(screen.getByText('Solo el Scrum Master puede realizar esta acción')).toBeTruthy()
    expect(onCambio).not.toHaveBeenCalled()
  })

  it('US4.3: 404 al reabrir muestra el aviso de proyecto inexistente', async () => {
    responder(json({ mensaje: 'no existe el proyecto' }, 404))
    const onCambio = abrirCerrado()

    await reabrir()

    expect((await screen.findByRole('alert')).textContent).toBe('no existe el proyecto')
    expect(onCambio).not.toHaveBeenCalled()
  })

  it('Escenario: conexión caída al reabrir muestra el error de conexión y el estado no cambia', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new TypeError('Failed to fetch')))
    const onCambio = abrirCerrado()

    await reabrir()

    expect((await screen.findByRole('alert')).textContent).toBe('No se pudo conectar con el servidor')
    expect(screen.getByText('Proyecto cerrado el 2026-11-11')).toBeTruthy()
    expect(screen.queryByText('En curso')).toBeNull()
    expect(onCambio).not.toHaveBeenCalled()
  })
})
