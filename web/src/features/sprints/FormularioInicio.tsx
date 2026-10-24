import { useState, type FormEvent } from 'react'
import { Button } from '../../compartido/componentes/Button'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { Field } from '../../compartido/componentes/Field'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { iniciarSprint } from './api'
import { clasificarError } from './errores'
import type { Sprint } from './sprint'

const CAMPOS = ['sprint_goal', 'fecha_inicio', 'fecha_fin']

type Props = { sprint: Sprint; alActualizar: (sprint: Sprint) => void }

export function FormularioInicio({ sprint, alActualizar }: Props) {
  const [enviando, setEnviando] = useState(false)
  const [aviso, setAviso] = useState<string>()
  const [errores, setErrores] = useState<Record<string, string>>({})

  async function alEnviar(evento: FormEvent<HTMLFormElement>) {
    evento.preventDefault()
    const datos = new FormData(evento.currentTarget)
    setAviso(undefined)
    setErrores({})
    setEnviando(true)
    try {
      alActualizar(
        await iniciarSprint(sprint.id, {
          sprint_goal: String(datos.get('sprint_goal')),
          fecha_inicio: String(datos.get('fecha_inicio')),
          fecha_fin: String(datos.get('fecha_fin')),
        }),
      )
    } catch (error) {
      const clasificado = clasificarError(error, CAMPOS)
      if ('campo' in clasificado) setErrores({ [clasificado.campo]: clasificado.mensaje })
      else setAviso(clasificado.aviso)
    } finally {
      setEnviando(false)
    }
  }

  const id = (campo: string) => `sprint-${sprint.id}-${campo}`

  return (
    <form onSubmit={alEnviar}>
      {aviso && <ErrorMessage message={aviso} />}
      <Field id={id('sprint_goal')} name="sprint_goal" label="Sprint Goal" error={errores.sprint_goal} />
      <Field id={id('fecha_inicio')} name="fecha_inicio" type="date" label="Fecha de inicio" error={errores.fecha_inicio} />
      <Field id={id('fecha_fin')} name="fecha_fin" type="date" label="Fecha de fin" error={errores.fecha_fin} />
      {enviando && <LoadingIndicator />}
      <Button type="submit" disabled={enviando}>Iniciar Sprint</Button>
    </form>
  )
}
