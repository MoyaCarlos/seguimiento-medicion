import { afterEach, describe, expect, it, vi } from 'vitest'
import { ApiError, ConnectionError, get, post, put } from './client'
import { configurarSesionHttp } from './sesionHttp'

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
  configurarSesionHttp(null)
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

  it.each([502, 503, 504])(
    'lanza un error de conexión si el proxy responde %i con cuerpo vacío',
    async (status) => {
      respondWith(null, status)

      await expect(get('/projects/1')).rejects.toBeInstanceOf(ConnectionError)
    },
  )

  it('trata el 502 con cuerpo no JSON como error de conexión', async () => {
    respondWith('<html>Bad Gateway</html>', 502)

    await expect(post('/projects', {})).rejects.toBeInstanceOf(ConnectionError)
  })

  it('trata el 502 con cuerpo JSON como error de conexión', async () => {
    respondWithJson({ mensaje: 'x' }, 502)

    await expect(post('/projects', {})).rejects.toBeInstanceOf(ConnectionError)
  })

  it('no lee el cuerpo de una respuesta de gateway', async () => {
    const leerCuerpo = vi.fn(() => Promise.reject(new Error('no se debía leer el cuerpo')))
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({ status: 502, ok: false, text: leerCuerpo }),
    )

    await expect(get('/projects/1')).rejects.toBeInstanceOf(ConnectionError)
    expect(leerCuerpo).not.toHaveBeenCalled()
  })

  it('sigue dando un error de la API ante un 500 con cuerpo vacío', async () => {
    respondWith(null, 500)

    await expect(get('/projects/1')).rejects.toEqual(
      new ApiError(500, 'La API respondió con estado 500'),
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

describe('403 de pertenencia y de rol', () => {
  it('un 403 «No sos integrante de este proyecto» llega con tipo no_integrante', async () => {
    respondWithJson({ mensaje: 'No sos integrante de este proyecto' }, 403)

    await expect(get('/projects/1')).rejects.toMatchObject({
      status: 403,
      message: 'No sos integrante de este proyecto',
      tipo: 'no_integrante',
    })
  })

  it('un 403 de rol llega con tipo sin_permiso y conserva el mensaje del servidor', async () => {
    respondWithJson({ mensaje: 'Solo el Scrum Master puede realizar esta acción' }, 403)

    await expect(post('/projects/1/members', {})).rejects.toMatchObject({
      status: 403,
      message: 'Solo el Scrum Master puede realizar esta acción',
      tipo: 'sin_permiso',
    })
  })

  it('un 403 sin cuerpo interpretable llega con tipo sin_permiso', async () => {
    respondWith('<html>Forbidden</html>', 403)

    await expect(get('/projects/1')).rejects.toMatchObject({ status: 403, tipo: 'sin_permiso' })
  })

  it('un 404 y un 401 no llevan tipo', async () => {
    respondWithJson({ mensaje: 'proyecto inexistente' }, 404)
    await expect(get('/projects/1')).rejects.toMatchObject({ status: 404, tipo: undefined })

    respondWithJson({ mensaje: 'sesión vencida' }, 401)
    await expect(get('/projects/1')).rejects.toMatchObject({ status: 401, tipo: undefined })
  })
})

describe('sesión en el cliente de la API', () => {
  function configurar(token: string | null) {
    const alNoAutenticado = vi.fn()
    configurarSesionHttp({ obtenerToken: () => token, alNoAutenticado })
    return alNoAutenticado
  }

  it('con token envía Authorization: Bearer <token>', async () => {
    configurar('abc123')
    const fetchMock = respondWithJson({}, 200)

    await get('/projects/1')

    expect(fetchMock).toHaveBeenCalledWith(
      '/projects/1',
      expect.objectContaining({
        headers: expect.objectContaining({ Authorization: 'Bearer abc123' }),
      }),
    )
  })

  it('sin token no envía la cabecera Authorization', async () => {
    configurar(null)
    const fetchMock = respondWithJson({}, 200)

    await get('/projects/1')

    const { headers } = fetchMock.mock.calls[0][1] as { headers: Record<string, string> }
    expect(headers).not.toHaveProperty('Authorization')
  })

  it('un 401 de una petición con token llama a alNoAutenticado y relanza el ApiError', async () => {
    const alNoAutenticado = configurar('abc123')
    respondWithJson({ mensaje: 'sesión vencida' }, 401)

    await expect(get('/projects/1')).rejects.toEqual(new ApiError(401, 'sesión vencida'))
    expect(alNoAutenticado).toHaveBeenCalledTimes(1)
  })

  it('un 401 sin token no llama a alNoAutenticado', async () => {
    const alNoAutenticado = configurar(null)
    respondWithJson({ mensaje: 'no autenticado' }, 401)

    await expect(get('/projects/1')).rejects.toMatchObject({ status: 401 })
    expect(alNoAutenticado).not.toHaveBeenCalled()
  })

  it.each(['/auth/ingreso', '/auth/registro'])(
    'un 401 de %s aunque haya token no llama a alNoAutenticado y llega como ApiError 401',
    async (ruta) => {
      const alNoAutenticado = configurar('abc123')
      respondWithJson({ mensaje: 'credenciales inválidas' }, 401)

      await expect(post(ruta, {})).rejects.toEqual(new ApiError(401, 'credenciales inválidas'))
      expect(alNoAutenticado).not.toHaveBeenCalled()
    },
  )

  it('sin enganche configurado se comporta como antes', async () => {
    configurarSesionHttp(null)
    const fetchMock = respondWithJson({ mensaje: 'no autenticado' }, 401)

    await expect(get('/projects/1')).rejects.toEqual(new ApiError(401, 'no autenticado'))

    const { headers } = fetchMock.mock.calls[0][1] as { headers: Record<string, string> }
    expect(headers).not.toHaveProperty('Authorization')
  })
})
