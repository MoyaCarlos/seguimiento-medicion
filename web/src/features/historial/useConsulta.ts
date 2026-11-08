import { useEffect, useState } from 'react'
import { ApiError, ConnectionError } from '../../compartido/api/client'

type Consulta<T> = { datos?: T; error?: string }

function mensajeDeError(causa: unknown): string {
  if (causa instanceof ApiError || causa instanceof ConnectionError) return causa.message
  return 'Ocurrió un error inesperado'
}

export function useConsulta<T>(consultar: () => Promise<T>, clave: string): Consulta<T> {
  const [consulta, setConsulta] = useState<Consulta<T> & { clave?: string }>({})

  useEffect(() => {
    let vigente = true

    consultar()
      .then((datos) => {
        if (vigente) setConsulta({ datos, clave })
      })
      .catch((causa: unknown) => {
        if (vigente) setConsulta({ error: mensajeDeError(causa), clave })
      })

    return () => {
      vigente = false
    }
  }, [clave])

  return consulta.clave === clave ? consulta : {}
}
