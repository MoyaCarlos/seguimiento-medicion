import { useState, type FormEvent } from 'react'
import { Button } from '../../compartido/componentes/Button'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { registrarDefecto, type HistoriaDeProyecto, type Severidad } from './api'
import { clasificarError } from './errores'

const SEVERIDADES: Severidad[] = ['Alta', 'Media', 'Baja']

type Props = {
  proyectoId: string
  historias: HistoriaDeProyecto[]
  alRegistrar: () => void
}

export function FormularioDefecto({ proyectoId, historias, alRegistrar }: Props) {
  const [descripcion, setDescripcion] = useState('')
  const [severidad, setSeveridad] = useState('')
  const [historiaId, setHistoriaId] = useState('')
  const [enviando, setEnviando] = useState(false)
  const [confirmado, setConfirmado] = useState(false)
  const [aviso, setAviso] = useState('')
  const [errores, setErrores] = useState<Record<string, string>>({})

  const grupos = agruparPorOrigen(historias)

  async function guardar(evento: FormEvent<HTMLFormElement>) {
    evento.preventDefault()
    if (enviando) return
    setConfirmado(false)
    setAviso('')

    const validacion: Record<string, string> = {}
    if (!descripcion.trim()) validacion.descripcion = 'La descripción es obligatoria'
    if (!severidad) validacion.severidad = 'La severidad es obligatoria'
    if (!historiaId) validacion.historia_id = 'La historia relacionada es obligatoria'
    setErrores(validacion)
    if (Object.keys(validacion).length > 0) return

    setEnviando(true)
    try {
      await registrarDefecto(proyectoId, {
        historia_id: Number(historiaId),
        descripcion: descripcion.trim(),
        severidad: severidad as Severidad,
      })
      setDescripcion('')
      setSeveridad('')
      setHistoriaId('')
      setConfirmado(true)
      alRegistrar()
    } catch (error) {
      const clasificado = clasificarError(error)
      if ('campo' in clasificado) setErrores({ [clasificado.campo]: clasificado.mensaje })
      else setAviso(clasificado.mensaje)
    } finally {
      setEnviando(false)
    }
  }

  function limpiar(campo: string) {
    setErrores((actuales) => ({ ...actuales, [campo]: '' }))
    setAviso('')
    setConfirmado(false)
  }

  return (
    <form onSubmit={guardar} noValidate aria-label="Registrar defecto">
      {aviso && <ErrorMessage message={aviso} />}
      {confirmado && <p className="message" role="status">Defecto registrado</p>}

      <div className="field">
        <label htmlFor="descripcion">Descripción</label>
        <textarea
          id="descripcion"
          name="descripcion"
          value={descripcion}
          maxLength={2000}
          onChange={(evento) => {
            setDescripcion(evento.target.value)
            limpiar('descripcion')
          }}
          disabled={enviando}
          aria-invalid={Boolean(errores.descripcion)}
          aria-describedby={errores.descripcion ? 'descripcion-error' : undefined}
        />
        {errores.descripcion && <ErrorMessage id="descripcion-error" message={errores.descripcion} />}
      </div>

      <div className="field">
        <label htmlFor="severidad">Severidad</label>
        <select
          id="severidad"
          name="severidad"
          value={severidad}
          onChange={(evento) => {
            setSeveridad(evento.target.value)
            limpiar('severidad')
          }}
          disabled={enviando}
          aria-invalid={Boolean(errores.severidad)}
          aria-describedby={errores.severidad ? 'severidad-error' : undefined}
        >
          <option value="">Elegí una severidad</option>
          {SEVERIDADES.map((valor) => (
            <option key={valor} value={valor}>{valor}</option>
          ))}
        </select>
        {errores.severidad && <ErrorMessage id="severidad-error" message={errores.severidad} />}
      </div>

      <div className="field">
        <label htmlFor="historia">Historia relacionada</label>
        <select
          id="historia"
          name="historia"
          value={historiaId}
          onChange={(evento) => {
            setHistoriaId(evento.target.value)
            limpiar('historia_id')
          }}
          disabled={enviando}
          aria-invalid={Boolean(errores.historia_id)}
          aria-describedby={errores.historia_id ? 'historia-error' : undefined}
        >
          <option value="">Elegí una historia</option>
          {grupos.map(([origen, items]) => (
            <optgroup key={origen} label={origen}>
              {items.map((historia) => (
                <option key={historia.id} value={String(historia.id)}>
                  {`${historia.titulo} (#${historia.id})`}
                </option>
              ))}
            </optgroup>
          ))}
        </select>
        {errores.historia_id && <ErrorMessage id="historia-error" message={errores.historia_id} />}
      </div>

      {enviando && <LoadingIndicator />}
      <Button type="submit" disabled={enviando}>Registrar defecto</Button>
    </form>
  )
}

function agruparPorOrigen(historias: HistoriaDeProyecto[]): Array<[string, HistoriaDeProyecto[]]> {
  const grupos = new Map<string, HistoriaDeProyecto[]>()
  for (const historia of historias) {
    const grupo = grupos.get(historia.origen) ?? []
    grupo.push(historia)
    grupos.set(historia.origen, grupo)
  }
  return [...grupos.entries()]
}
