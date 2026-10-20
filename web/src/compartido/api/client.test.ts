import { describe, expect, it, vi } from 'vitest'
import { ApiError, apiFetch } from './client'

describe('apiFetch', () => {
  it('returns the decoded response for a successful request', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ id: 42 }), { status: 200 }),
    ))

    await expect(apiFetch<{ id: number }>('/projects/42')).resolves.toEqual({ id: 42 })
  })

  it('exposes field, message and status from API errors', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ campo: 'nombre', mensaje: 'el nombre es obligatorio' }), { status: 400 }),
    ))

    await expect(apiFetch('/projects')).rejects.toEqual(
      new ApiError(400, 'el nombre es obligatorio', 'nombre'),
    )
  })

  it('uses a general message when the API error has no field', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ mensaje: 'proyecto inexistente' }), { status: 404 }),
    ))

    await expect(apiFetch('/projects/99')).rejects.toMatchObject({
      status: 404,
      message: 'proyecto inexistente',
      field: undefined,
    })
  })
})
