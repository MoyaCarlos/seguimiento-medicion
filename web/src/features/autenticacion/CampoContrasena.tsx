import { useState } from 'react'
import type { InputHTMLAttributes } from 'react'
import { Field } from '../../compartido/componentes/Field'
import { Icono } from '../../compartido/componentes/Icono'

type Props = Omit<InputHTMLAttributes<HTMLInputElement>, 'type'> & {
  label: string
  error?: string
  ayuda?: string
  etiquetaVer?: string
}

export function CampoContrasena({ etiquetaVer = 'Mostrar contraseña', ...props }: Props) {
  const [visible, setVisible] = useState(false)

  return (
    <Field
      {...props}
      type={visible ? 'text' : 'password'}
      accion={
        <button
          type="button"
          className="campo-contrasena__ver"
          aria-label={etiquetaVer}
          aria-pressed={visible}
          onClick={() => setVisible((actual) => !actual)}
        >
          <Icono nombre={visible ? 'ojoTachado' : 'ojo'} />
        </button>
      }
    />
  )
}
