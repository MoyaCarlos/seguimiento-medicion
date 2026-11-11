import { useState } from 'react'
import { Link, useLocation } from 'react-router'
import { BotonDatosDePrueba } from '../../compartido/demo/BotonDatosDePrueba'
import { FormularioCredenciales } from './FormularioCredenciales'
import { useSesion } from './SesionContext'

export function Ingresar() {
  const { ingresar } = useSesion()
  const { state } = useLocation()
  const [relleno, setRelleno] = useState<{ nombre: string; contrasena: string }>()
  const vencida = Boolean((state as { vencida?: boolean } | null)?.vencida)

  return (
    <>
      <FormularioCredenciales
        titulo="Ingresar"
        accion="Ingresar"
        aviso={vencida ? 'Tu sesión venció' : undefined}
        enviar={ingresar}
        relleno={relleno}
        pie={
          <>
            ¿No tenés cuenta? <Link to="/registro" state={state}>Registrate</Link>
          </>
        }
      />
      <div className="auth-demo">
        <BotonDatosDePrueba alUsar={(nombre, contrasena) => setRelleno({ nombre, contrasena })} />
      </div>
    </>
  )
}
