import { Navigate, Outlet, useLocation } from 'react-router'
import { Button } from '../../compartido/componentes/Button'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { useSesion } from './SesionContext'

export function RutaProtegida() {
  const { estado, reintentar } = useSesion()
  const location = useLocation()

  if (estado === 'cargando') return <LoadingIndicator />
  if (estado === 'sin_conexion') {
    return (
      <section>
        <ErrorMessage message="No pudimos conectar con el servidor" />
        <Button onClick={() => void reintentar()}>Reintentar</Button>
      </section>
    )
  }
  if (estado === 'anonimo') {
    return <Navigate to="/ingresar" replace state={{ desde: location.pathname + location.search }} />
  }
  return <Outlet />
}
