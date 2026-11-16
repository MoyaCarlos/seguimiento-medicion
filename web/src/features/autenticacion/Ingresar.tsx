import { useState } from 'react'
import { Link, useLocation } from 'react-router'
import { BotonDatosDePrueba } from '../../compartido/demo/BotonDatosDePrueba'
import { FormularioCredenciales } from './FormularioCredenciales'
import { PantallaAcceso } from './PantallaAcceso'
import { useSesion } from './SesionContext'

export function Ingresar() {
  const { ingresar } = useSesion()
  const { state } = useLocation()
  const [relleno, setRelleno] = useState<{ nombre: string; contrasena: string }>()
  const vencida = Boolean((state as { vencida?: boolean } | null)?.vencida)

  return (
    <PantallaAcceso titulo="Iniciar sesión">
      <FormularioCredenciales
        titulo="Ingresar"
        sobretitulo="Bienvenido"
        accion="Ingresar"
        bajada="Ingresá con tu cuenta para ver tus proyectos."
        aviso={vencida ? 'Tu sesión venció. Volvé a ingresar para continuar.' : undefined}
        enviar={ingresar}
        relleno={relleno}
        pie={
          <>
            ¿No tenés cuenta? <Link to="/registro" state={state}>Registrate</Link>
          </>
        }
      />
      <BotonDatosDePrueba alUsar={(nombre, contrasena) => setRelleno({ nombre, contrasena })} />
    </PantallaAcceso>
  )
}
