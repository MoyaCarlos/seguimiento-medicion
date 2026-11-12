import { afterEach, describe, expect, it, vi } from 'vitest'
import { ApiError, ConnectionError, descargar } from './client'
import { configurarSesionHttp } from './sesionHttp'

function responderDescarga(cuerpo: BodyInit, nombre?: string) {
  const headers: Record<string, string> = { 'Content-Type': 'application/pdf' }
  if (nombre) headers['Content-Disposition'] = `attachment; filename="${nombre}"`
  const fetchMock = vi.fn().mockResolvedValue(new Response(cuerpo, { status: 200, headers }))
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

afterEach(() => {
  vi.unstubAllGlobals()
  configurarSesionHttp(null)
})

describe('descarga de archivos', () => {
  it('US5.5: Given una descarga exitosa, When la solicito, Then recibo el archivo con el nombre de Content-Disposition', async () => {
    responderDescarga('%PDF-1.4 contenido', 'reporte-proyecto-3-sprint-1.pdf')

    const descarga = await descargar('/sprints/1/reporte')

    expect(await descarga.blob.text()).toBe('%PDF-1.4 contenido')
    expect(descarga.nombre).toBe('reporte-proyecto-3-sprint-1.pdf')
  })

  it('usa un nombre por defecto si la respuesta no trae Content-Disposition', async () => {
    responderDescarga('%PDF-1.4')

    await expect(descargar('/projects/1/reporte')).resolves.toMatchObject({ nombre: 'reporte.pdf' })
  })

  it('envía Accept: application/pdf y el token de la sesión', async () => {
    configurarSesionHttp({ obtenerToken: () => 'abc123', alNoAutenticado: () => {} })
    const fetchMock = responderDescarga('%PDF-1.4', 'reporte-proyecto-1.pdf')

    await descargar('/projects/1/reporte')

    expect(fetchMock).toHaveBeenCalledWith(
      '/projects/1/reporte',
      expect.objectContaining({
        method: 'GET',
        headers: expect.objectContaining({
          Accept: 'application/pdf',
          Authorization: 'Bearer abc123',
        }),
      }),
    )
  })

  it('US5.3: Given un error de generación del backend, When descargo, Then recibo su mensaje', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue(
        new Response(JSON.stringify({ mensaje: 'no se pudo generar el reporte' }), { status: 500 }),
      ),
    )

    await expect(descargar('/projects/1/reporte')).rejects.toEqual(
      new ApiError(500, 'no se pudo generar el reporte'),
    )
  })

  it('trata un gateway como error de conexión', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response(null, { status: 502 })))

    await expect(descargar('/projects/1/reporte')).rejects.toBeInstanceOf(ConnectionError)
  })

  it('lanza un error de conexión si no hay respuesta', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new TypeError('Failed to fetch')))

    await expect(descargar('/projects/1/reporte')).rejects.toBeInstanceOf(ConnectionError)
  })
})
