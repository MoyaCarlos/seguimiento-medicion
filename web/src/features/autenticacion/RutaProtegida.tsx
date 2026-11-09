import { Navigate, Outlet, useLocation } from 'react-router'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { useSesion } from './SesionContext'

export function RutaProtegida() {
  const { estado } = useSesion()
  const location = useLocation()

  if (estado === 'cargando') return <LoadingIndicator />
  if (estado === 'anonimo') {
    return <Navigate to="/ingresar" replace state={{ desde: location.pathname + location.search }} />
  }
  return <Outlet />
}
