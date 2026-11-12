import { Button } from '../componentes/Button'
import './estilos.css'

export const CONTRASENA_DEMO = 'demo1234'

const USUARIOS_DEMO = [
  { nombre: 'carlos.moya', rol: 'Scrum Master' },
  { nombre: 'ana.iriarte', rol: 'Product Builder' },
  { nombre: 'candela.vulcano', rol: 'Product Builder' },
  { nombre: 'jimena.martinez', rol: 'Product Builder' },
]

type Props = { alUsar?: (nombre: string, contrasena: string) => void }

export function CredencialesDemo({ alUsar }: Props) {
  return (
    <table className="demo-credenciales">
      <caption>Usuarios de la demostración</caption>
      <thead>
        <tr>
          <th scope="col">Usuario</th>
          <th scope="col">Contraseña</th>
          <th scope="col">Rol</th>
        </tr>
      </thead>
      <tbody>
        {USUARIOS_DEMO.map(({ nombre, rol }) => (
          <tr key={nombre}>
            <td><code>{nombre}</code></td>
            <td><code>{CONTRASENA_DEMO}</code></td>
            <td>
              {rol}
              {alUsar && (
                <>
                  <br />
                  <Button
                    type="button"
                    variant="secondary"
                    aria-label={`Usar ${nombre}`}
                    onClick={() => alUsar(nombre, CONTRASENA_DEMO)}
                  >
                    Usar
                  </Button>
                </>
              )}
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}
