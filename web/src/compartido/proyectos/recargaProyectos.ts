import { createContext, useContext } from 'react'

const RecargaProyectos = createContext<() => void>(() => {})

export const ProveedorRecargaProyectos = RecargaProyectos.Provider

export function useRecargarProyectos() {
  return useContext(RecargaProyectos)
}
