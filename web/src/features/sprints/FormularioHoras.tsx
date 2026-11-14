import { useRef, useEffect, useState, type FormEvent } from 'react'
import { Button } from '../../compartido/componentes/Button'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { Field } from '../../compartido/componentes/Field'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { imputarHoras } from './api'
import { clasificarError } from './errores'
import { fechaPropuesta, normalizarHoras } from './horas'

const CAMPOS = ['fecha', 'actividad', 'horas']

type Props = {
  historiaId: number
  fechaFin?: string | null
  alRegistrar: () => void
  alCambiarEnvio?: (enviando: boolean) => void
}

export function FormularioHoras({ historiaId, fechaFin, alRegistrar, alCambiarEnvio }: Props) {
  const [fecha, setFecha] = useState(() => fechaPropuesta(new Date(), fechaFin))
  const [actividad, setActividad] = useState('')
  const [horas, setHoras] = useState('')
  const [enviando, setEnviando] = useState(false)
  const [confirmado, setConfirmado] = useState(false)
  const [aviso, setAviso] = useState<string>()
  const [errores, setErrores] = useState<Record<string, string>>({})
  const focoAlEnviar = useRef<HTMLElement | null>(null)

  useEffect(() => {
    if (enviando) return
    if (document.activeElement === document.body) focoAlEnviar.current?.focus()
    focoAlEnviar.current = null
  }, [enviando])

  async function guardar(evento: FormEvent<HTMLFormElement>) {
    evento.preventDefault()
    if (enviando) return
    setConfirmado(false)
    setAviso(undefined)
    setErrores({})
    focoAlEnviar.current = document.activeElement as HTMLElement | null
    setEnviando(true)
    alCambiarEnvio?.(true)
    try {
      await imputarHoras(historiaId, {
        fecha,
        actividad,
        horas: normalizarHoras(horas),
      })
      setActividad('')
      setHoras('')
      setConfirmado(true)
      alRegistrar()
    } catch (error) {
      const clasificado = clasificarError(error, CAMPOS)
      if ('campo' in clasificado) setErrores({ [clasificado.campo]: clasificado.mensaje })
      else setAviso(clasificado.aviso)
    } finally {
      setEnviando(false)
      alCambiarEnvio?.(false)
    }
  }

  const id = (campo: string) => `historia-${historiaId}-${campo}`

  return (
    <form onSubmit={guardar} noValidate>
      {aviso && <ErrorMessage message={aviso} />}
      {confirmado && <p className="message" role="status">Horas registradas</p>}
      <Field
        id={id('fecha')}
        name="fecha"
        type="date"
        label="Fecha"
        value={fecha}
        onChange={(evento) => setFecha(evento.target.value)}
        disabled={enviando}
        error={errores.fecha}
      />
      <Field
        id={id('actividad')}
        name="actividad"
        label="Actividad"
        value={actividad}
        onChange={(evento) => setActividad(evento.target.value)}
        disabled={enviando}
        error={errores.actividad}
      />
      <Field
        id={id('horas')}
        name="horas"
        inputMode="text"
        label="Horas"
        help="Ej.: 1,5 o 1:30 (horas:minutos)"
        value={horas}
        onChange={(evento) => setHoras(evento.target.value)}
        disabled={enviando}
        error={errores.horas}
      />
      {enviando && <LoadingIndicator />}
      <Button type="submit" disabled={enviando}>Guardar horas</Button>
    </form>
  )
}
