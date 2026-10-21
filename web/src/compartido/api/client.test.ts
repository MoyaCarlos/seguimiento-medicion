import { afterEach, describe, expect, it, vi } from 'vitest'
import { ApiError, ConnectionError, get, post, put } from './client'

function respondWith(body: BodyInit | null, status: number) {
  const fetchMock = vi.fn().mockResolvedValue(new Response(body, { status }))
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

function respondWithJson(body: unknown, status: number) {
  return respondWith(JSON.stringify(body), status)
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('cliente de la API', () => {
  it('devuelve los datos de una respuesta exitosa', async () => {
    respondWithJson({ id: 42 }, 200)

    await expect(get<{ id: number }>('/projects/42')).resolves.toEqual({ id: 42 })
  })

  it('devuelve undefined cuando la respuesta exitosa no trae cuerpo', async () => {
    respondWith(null, 204)

    await expect(post('/sprints/1/cerrar', {})).resolves.toBeUndefined()
  })

  it('expone estado, campo y mensaje de un 400 con campo', async () => {
    respondWithJson({ campo: 'nombre', mensaje: 'el nombre es obligatorio' }, 400)

    await expect(post('/projects', {})).rejects.toEqual(
      new ApiError(400, 'el nombre es obligatorio', 'nombre'),
    )
  })

  it('expone estado y mensaje de un 400 sin campo, con el campo ausente', async () => {
    respondWithJson({ mensaje: 'datos inválidos' }, 400)

    await expect(post('/projects', {})).rejects.toMatchObject({
      status: 400,
      message: 'datos inválidos',
      field: undefined,
    })
  })

  it('expone estado y mensaje de un 404', async () => {
    respondWithJson({ mensaje: 'proyecto inexistente' }, 404)

    await expect(get('/projects/99')).rejects.toMatchObject({
      status: 404,
      message: 'proyecto inexistente',
      field: undefined,
    })
  })

  it('expone estado y mensaje de un 409', async () => {
    respondWithJson({ mensaje: 'el sprint ya está iniciado' }, 409)

    await expect(post('/sprints/1/iniciar', {})).rejects.toMatchObject({
      status: 409,
      message: 'el sprint ya está iniciado',
    })
  })

  it('devuelve un error con el estado y un mensaje genérico si el cuerpo del error no es JSON', async () => {
    respondWith('<html>Bad Gateway</html>', 500)

    await expect(get('/projects/1')).rejects.toEqual(
      new ApiError(500, 'La API respondió con estado 500'),
    )
  })

  it('devuelve un error con el estado y un mensaje genérico si una respuesta exitosa no es JSON', async () => {
    respondWith('<html>index</html>', 200)

    await expect(get('/projects/1')).rejects.toEqual(
      new ApiError(200, 'La API respondió con estado 200'),
    )
  })

  it('lanza un error de conexión, distinto de los de la API, si no hay respuesta', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new TypeError('Failed to fetch')))

    const error = await get('/projects/1').catch((caught: unknown) => caught)

    expect(error).toBeInstanceOf(ConnectionError)
    expect(error).not.toBeInstanceOf(ApiError)
  })

  it('envía Accept: application/json en todas las llamadas', async () => {
    const fetchMock = respondWithJson({}, 200)

    await get('/projects/1')

    expect(fetchMock).toHaveBeenCalledWith(
      '/projects/1',
      expect.objectContaining({
        method: 'GET',
        headers: expect.objectContaining({ Accept: 'application/json' }),
      }),
    )
  })

  it('post envía el cuerpo como JSON', async () => {
    const fetchMock = respondWithJson({ id: 1 }, 201)

    await post('/projects', { nombre: 'Demo' })

    expect(fetchMock).toHaveBeenCalledWith(
      '/projects',
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({ nombre: 'Demo' }),
        headers: expect.objectContaining({
          Accept: 'application/json',
          'Content-Type': 'application/json',
        }),
      }),
    )
  })

  it('put envía el cuerpo como JSON', async () => {
    const fetchMock = respondWithJson({ id: 1 }, 200)

    await put('/projects/1', { nombre: 'Nuevo' })

    expect(fetchMock).toHaveBeenCalledWith(
      '/projects/1',
      expect.objectContaining({
        method: 'PUT',
        body: JSON.stringify({ nombre: 'Nuevo' }),
        headers: expect.objectContaining({ 'Content-Type': 'application/json' }),
      }),
    )
  })
})
