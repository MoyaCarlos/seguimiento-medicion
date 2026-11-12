import { Link, useLocation } from 'react-router'
import { FormularioCredenciales } from './FormularioCredenciales'
import { PantallaAcceso } from './PantallaAcceso'
import { useSesion } from './SesionContext'

export function Registro() {
  const { registrar } = useSesion()
  const { state } = useLocation()

  return (
    <PantallaAcceso titulo="Crear cuenta">
      <FormularioCredenciales
        titulo="Crear cuenta"
        sobretitulo="Nueva cuenta"
        accion="Crear cuenta"
        bajada="Elegí un nombre de usuario y una contraseña. Con tu cuenta vas a poder crear proyectos y sumarte a equipos."
        contrasenaNueva
        enviar={registrar}
        pie={
          <>
            ¿Ya tenés cuenta? <Link to="/ingresar" state={state}>Ingresá</Link>
          </>
        }
      />
    </PantallaAcceso>
  )
}
