import { Link, useLocation } from 'react-router'
import { FormularioCredenciales } from './FormularioCredenciales'
import { useSesion } from './SesionContext'

export function Ingresar() {
  const { ingresar } = useSesion()
  const { state } = useLocation()
  const vencida = Boolean((state as { vencida?: boolean } | null)?.vencida)

  return (
    <FormularioCredenciales
      titulo="Ingresar"
      accion="Ingresar"
      aviso={vencida ? 'Tu sesión venció' : undefined}
      enviar={ingresar}
      pie={
        <>
          ¿No tenés cuenta? <Link to="/registro" state={state}>Registrate</Link>
        </>
      }
    />
  )
}
