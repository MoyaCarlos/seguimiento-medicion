import { useState, type ReactNode } from 'react'
import { ApiError, ConnectionError, type Descarga } from '../api/client'
import { Button } from './Button'
import { ErrorMessage } from './ErrorMessage'
import './BotonExportar.css'

type BotonExportarProps = {
  children: ReactNode
  descargar: () => Promise<Descarga>
}

function mensajeDeError(causa: unknown): string {
  if (causa instanceof ApiError || causa instanceof ConnectionError) return causa.message
  return 'Ocurrió un error inesperado'
}

function guardarArchivo(blob: Blob, nombre: string) {
  const url = URL.createObjectURL(blob)
  const enlace = document.createElement('a')
  enlace.href = url
  enlace.download = nombre
  enlace.click()
  URL.revokeObjectURL(url)
}

export function BotonExportar({ children, descargar }: BotonExportarProps) {
  const [generando, setGenerando] = useState(false)
  const [error, setError] = useState<string>()

  async function exportar() {
    setError(undefined)
    setGenerando(true)
    try {
      const { blob, nombre } = await descargar()
      guardarArchivo(blob, nombre)
    } catch (causa) {
      setError(mensajeDeError(causa))
    } finally {
      setGenerando(false)
    }
  }

  return (
    <div className="reporte-exportar">
      <Button type="button" onClick={exportar} disabled={generando}>
        {generando ? 'Generando…' : children}
      </Button>
      {error && <ErrorMessage message={error} />}
    </div>
  )
}
