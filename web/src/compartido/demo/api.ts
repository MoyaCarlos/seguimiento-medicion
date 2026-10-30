import { get, post } from '../api/client'

export const PROYECTO_DE_PRUEBA = 25

export async function datosDePruebaHabilitados(): Promise<boolean> {
  try {
    await get('/demo/estado')
    return true
  } catch {
    return false
  }
}

export const cargarDatosDePrueba = () => post<{ proyecto_id: number }>('/demo/datos', undefined)
