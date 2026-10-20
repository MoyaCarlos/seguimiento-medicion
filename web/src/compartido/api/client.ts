export type ApiErrorPayload = {
  campo?: string
  mensaje?: string
}

export class ApiError extends Error {
  readonly status: number
  readonly field?: string

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

export async function apiFetch<T>(
  path: string,
  options?: RequestInit,
): Promise<T> {
  const response = await fetch(path, {
    ...options,
    headers: {
      Accept: 'application/json',
      ...options?.headers,
    },
  })

  const body = await parseResponseBody(response)
  if (!response.ok) {
    const error = isApiErrorPayload(body) ? body : {}
    throw new ApiError(
      response.status,
      error.mensaje ?? `La API respondió con estado ${response.status}`,
      error.campo || undefined,
    )
  }

  return body as T
}

async function parseResponseBody(response: Response): Promise<unknown> {
  const text = await response.text()
  if (!text) return undefined

  try {
    return JSON.parse(text) as unknown
  } catch {
    throw new ApiError(response.status, 'La API devolvió una respuesta inválida')
  }
}

function isApiErrorPayload(value: unknown): value is ApiErrorPayload {
  return typeof value === 'object' && value !== null
}
