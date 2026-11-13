import { avisarNoAutenticado, tokenActual } from './sesionHttp'

export class ApiError extends Error {
  readonly status: number
  readonly field?: string
  readonly tipo?: 'no_integrante' | 'sin_permiso'

  constructor(
    status: number,
    message: string,
    field?: string,
  ) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.field = field
  }
}

export class ConnectionError extends Error {
  constructor() {
    super('No se pudo conectar con el servidor')
    this.name = 'ConnectionError'
  }
}

type Method = 'GET' | 'POST' | 'PUT'

export type Descarga = { blob: Blob; nombre: string }

export function get<T>(path: string): Promise<T> {
  return request<T>('GET', path)
}

export async function descargar(path: string): Promise<Descarga> {
  const token = tokenActual()
  const response = await send('GET', path, undefined, token, 'application/pdf')
  if (gatewayStatuses.includes(response.status)) throw new ConnectionError()
  if (response.status === 401 && token && !rutasDeAutenticacion.includes(path)) {
    avisarNoAutenticado()
  }
  if (!response.ok) throw errorFrom(response.status, parseJson(await response.text()))
  return { blob: await response.blob(), nombre: nombreDeArchivo(response) }
}

export function post<T>(path: string, body: unknown): Promise<T> {
  return request<T>('POST', path, body)
}

export function put<T>(path: string, body: unknown): Promise<T> {
  return request<T>('PUT', path, body)
}

async function request<T>(method: Method, path: string, body?: unknown): Promise<T> {
  const token = tokenActual()
  const response = await send(method, path, body, token)
  if (gatewayStatuses.includes(response.status)) throw new ConnectionError()
  const text = await response.text()
  const data = parseJson(text)

  if (response.status === 401 && token && !rutasDeAutenticacion.includes(path)) {
    avisarNoAutenticado()
  }
  if (!response.ok) throw errorFrom(response.status, data)
  if (data === invalid) throw genericError(response.status)
  return data as T
}

async function send(
  method: Method,
  path: string,
  body: unknown,
  token: string | null,
  accept = 'application/json',
): Promise<Response> {
  const headers: Record<string, string> = { Accept: accept }
  if (token) headers['Authorization'] = `Bearer ${token}`
  if (body !== undefined) headers['Content-Type'] = 'application/json'

  try {
    return await fetch(path, {
      method,
      headers,
      body: body === undefined ? undefined : JSON.stringify(body),
    })
  } catch {
    throw new ConnectionError()
  }
}

const gatewayStatuses = [502, 503, 504]

const rutasDeAutenticacion = ['/auth/ingreso', '/auth/registro']

const invalid = Symbol('cuerpo no interpretable')

function parseJson(text: string): unknown {
  if (!text) return undefined
  try {
    return JSON.parse(text)
  } catch {
    return invalid
  }
}

function errorFrom(status: number, data: unknown): ApiError {
  if (typeof data !== 'object' || data === null) return genericError(status)

  const { mensaje, campo } = data as { mensaje?: unknown; campo?: unknown }
  if (typeof mensaje !== 'string' || !mensaje) return genericError(status)
  return new ApiError(status, mensaje, typeof campo === 'string' && campo ? campo : undefined)
}

function genericError(status: number): ApiError {
  return new ApiError(status, `La API respondió con estado ${status}`)
}

function nombreDeArchivo(response: Response): string {
  const disposicion = response.headers.get('Content-Disposition') ?? ''
  const coincidencia = /filename\s*=\s*"?([^";]+)"?/i.exec(disposicion)
  return coincidencia?.[1] ?? 'reporte.pdf'
}
