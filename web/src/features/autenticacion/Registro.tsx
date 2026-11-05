import { Link, useLocation } from 'react-router'
import { FormularioCredenciales } from './FormularioCredenciales'
import { useSesion } from './SesionContext'

export function Registro() {
  const { registrar } = useSesion()
  const { state } = useLocation()

  return (
    <FormularioCredenciales
      titulo="Crear cuenta"
      accion="Crear cuenta"
      contrasenaNueva
      enviar={registrar}
      pie={
        <>
          ¿Ya tenés cuenta? <Link to="/ingresar" state={state}>Ingresá</Link>
        </>
      }
    />
  )
}
