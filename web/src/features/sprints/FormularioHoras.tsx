import { useEffect, useState, type FormEvent } from 'react'
import { Button } from '../../compartido/componentes/Button'
import { Field } from '../../compartido/componentes/Field'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { imputarHoras, listarIntegrantes } from './api'
import { fechaPropuesta, normalizarHoras, type IntegranteDeProyecto } from './horas'

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

  useEffect(() => {
    let vigente = true

    listarIntegrantes(proyectoId).then((lista) => {
      if (vigente) setIntegrantes(lista)
    })

    return () => {
      vigente = false
    }
  }, [proyectoId])

  async function guardar(evento: FormEvent<HTMLFormElement>) {
    evento.preventDefault()
    if (enviando) return
    setConfirmado(false)
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
    } finally {
      setEnviando(false)
    }
  }

  if (!integrantes) return <LoadingIndicator />

  const id = (campo: string) => `historia-${historiaId}-${campo}`

  return (
    <form onSubmit={guardar} noValidate>
      {confirmado && <p className="message" role="status">Horas registradas</p>}
      <div className="field">
        <label htmlFor={id('integrante')}>Integrante</label>
        <select
          id={id('integrante')}
          name="integrante"
          value={integrante}
          onChange={(evento) => setIntegrante(evento.target.value)}
          disabled={enviando}
        >
          <option value="">Elegí un integrante</option>
          {integrantes.map((persona) => (
            <option key={persona.id} value={String(persona.id)}>
              {persona.nombre}
            </option>
          ))}
        </select>
      </div>
      <Field
        id={id('fecha')}
        name="fecha"
        type="date"
        label="Fecha"
        value={fecha}
        onChange={(evento) => setFecha(evento.target.value)}
        disabled={enviando}
      />
      <Field
        id={id('actividad')}
        name="actividad"
        label="Actividad"
        value={actividad}
        onChange={(evento) => setActividad(evento.target.value)}
        disabled={enviando}
      />
      <Field
        id={id('horas')}
        name="horas"
        inputMode="decimal"
        label="Horas"
        value={horas}
        onChange={(evento) => setHoras(evento.target.value)}
        disabled={enviando}
      />
      {enviando && <LoadingIndicator />}
      <Button type="submit" disabled={enviando}>Guardar horas</Button>
    </form>
  )
}
