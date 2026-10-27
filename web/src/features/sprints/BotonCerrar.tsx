import { useState } from 'react'
import { post } from '../../compartido/api/client'
import { Button } from '../../compartido/componentes/Button'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { clasificarError } from './errores'
import type { Sprint } from './sprint'

const cerrarSprint = (id: number) => post<Sprint>(`/sprints/${id}/cerrar`, undefined)

type Props = { sprint: Sprint; alActualizar: (sprint: Sprint) => void }

export function BotonCerrar({ sprint, alActualizar }: Props) {
  const [enviando, setEnviando] = useState(false)
  const [aviso, setAviso] = useState<string>()

  async function alHacerClic() {
    setAviso(undefined)
    setEnviando(true)
    try {
      alActualizar(await cerrarSprint(sprint.id))
    } catch (error) {
      const clasificado = clasificarError(error, [])
      setAviso('aviso' in clasificado ? clasificado.aviso : clasificado.mensaje)
    } finally {
      setEnviando(false)
    }
  }

  return (
    <div>
      {aviso && <ErrorMessage message={aviso} />}
      {enviando && <LoadingIndicator />}
      <Button type="button" onClick={alHacerClic} disabled={enviando}>Cerrar Sprint</Button>
    </div>
  )
}
