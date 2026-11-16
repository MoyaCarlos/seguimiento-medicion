import { useEffect, useState } from 'react'
import { get } from '../api/client'

/** Id de la persona con la sesión iniciada (GET /auth/yo), o null mientras no se conoce. Sirve para resaltar «mi» fila. */
export function useUsuarioActual(): number | null {
  const [id, setId] = useState<number | null>(null)

  useEffect(() => {
    let vigente = true

    get<{ id?: unknown }>('/auth/yo')
      .then((usuario) => {
        if (vigente && typeof usuario?.id === 'number') setId(usuario.id)
      })
      .catch(() => {})

    return () => {
      vigente = false
    }
  }, [])

  return id
}
