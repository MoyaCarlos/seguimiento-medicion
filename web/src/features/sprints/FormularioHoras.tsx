import { useEffect, useState, type FormEvent } from 'react'
import { Button } from '../../compartido/componentes/Button'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { Field } from '../../compartido/componentes/Field'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { imputarHoras, listarIntegrantes } from './api'
import { clasificarError } from './errores'
import { fechaPropuesta, normalizarHoras, type IntegranteDeProyecto } from './horas'

const CAMPOS = ['integrante_id', 'fecha', 'actividad', 'horas']

type Props = {
  historiaId: number
  proyectoId: number
  fechaFin?: string | null
  alRegistrar: () => void
}

export function FormularioHoras({ historiaId, proyectoId, fechaFin, alRegistrar }: Props) {
  const [integrantes, setIntegrantes] = useState<IntegranteDeProyecto[]>()
  const [integrante, setIntegrante] = useState('')
  const [fecha, setFecha] = useState(() => fechaPropuesta(new Date(), fechaFin))
  const [actividad, setActividad] = useState('')
  const [horas, setHoras] = useState('')
  const [enviando, setEnviando] = useState(false)
  const [confirmado, setConfirmado] = useState(false)
  const [errorIntegrantes, setErrorIntegrantes] = useState<string>()
  const [aviso, setAviso] = useState<string>()
  const [errores, setErrores] = useState<Record<string, string>>({})

  useEffect(() => {
    let vigente = true

    listarIntegrantes(proyectoId)
      .then((lista) => {
        if (vigente) setIntegrantes(lista)
      })
      .catch((error: unknown) => {
        if (!vigente) return
        const clasificado = clasificarError(error, [])
        setErrorIntegrantes('aviso' in clasificado ? clasificado.aviso : clasificado.mensaje)
      })

    return () => {
      vigente = false
    }
  }, [proyectoId])

  async function guardar(evento: FormEvent<HTMLFormElement>) {
    evento.preventDefault()
    if (enviando) return
    setConfirmado(false)
    setAviso(undefined)
    if (!integrante) {
      setErrores({ integrante_id: 'Elegí un integrante' })
      return
    }
    setErrores({})
    setEnviando(true)
    try {
      await imputarHoras(historiaId, {
        integrante_id: Number(integrante),
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
    }
  }

  if (errorIntegrantes) return <ErrorMessage message={errorIntegrantes} />
  if (!integrantes) return <LoadingIndicator />

  const id = (campo: string) => `historia-${historiaId}-${campo}`

  return (
    <form onSubmit={guardar} noValidate>
      {aviso && <ErrorMessage message={aviso} />}
      {confirmado && <p className="message" role="status">Horas registradas</p>}
      <div className="field">
        <label htmlFor={id('integrante')}>Integrante</label>
        <select
          id={id('integrante')}
          name="integrante"
          value={integrante}
          onChange={(evento) => setIntegrante(evento.target.value)}
          disabled={enviando}
          aria-invalid={Boolean(errores.integrante_id)}
          aria-describedby={errores.integrante_id ? `${id('integrante')}-error` : undefined}
        >
          <option value="">Elegí un integrante</option>
          {integrantes.map((persona) => (
            <option key={persona.id} value={String(persona.id)}>
              {persona.nombre}
            </option>
          ))}
        </select>
        {errores.integrante_id && (
          <ErrorMessage id={`${id('integrante')}-error`} message={errores.integrante_id} />
        )}
      </div>
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
        inputMode="decimal"
        label="Horas"
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
