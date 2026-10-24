import { useState } from 'react'
import { useParams } from 'react-router'
import { Button } from '../../compartido/componentes/Button'
import { ErrorMessage } from '../../compartido/componentes/ErrorMessage'
import { LoadingIndicator } from '../../compartido/componentes/LoadingIndicator'
import { crearSprint } from './api'
import { clasificarError } from './errores'
import { SeccionSprints } from './SeccionSprints'
import { agruparPorEstado, cargarSprintsIniciales, ESTADOS, guardarSprint, type Sprint } from './sprint'
import './sprints.css'

export function Component() {
  const { id } = useParams()
  const [sprints, setSprints] = useState(cargarSprintsIniciales)
  const [creando, setCreando] = useState(false)
  const [aviso, setAviso] = useState<string>()

  const guardar = (sprint: Sprint) => setSprints((actuales) => guardarSprint(actuales, sprint))

  async function crear() {
    setAviso(undefined)
    setCreando(true)
    try {
      guardar(await crearSprint(Number(id)))
    } catch (error) {
      const clasificado = clasificarError(error, [])
      setAviso('aviso' in clasificado ? clasificado.aviso : clasificado.mensaje)
    } finally {
      setCreando(false)
    }
  }

  const grupos = agruparPorEstado(sprints)

  return (
    <div className="sprints-pantalla">
      <h1>Sprints</h1>
      <p>Proyecto actual: <strong>{id}</strong></p>
      <p>
        Los Sprints existentes no se listan todavía. Para consultarlos, usá los comandos de
        docs/curl-frontend-sprint-2.md.
      </p>
      {aviso && <ErrorMessage message={aviso} />}
      {creando && <LoadingIndicator />}
      <div>
        <Button type="button" onClick={crear} disabled={creando}>Crear Sprint</Button>
      </div>
      <div className="sprints-secciones">
        {ESTADOS.map((estado) => (
          <SeccionSprints key={estado} estado={estado} sprints={grupos[estado]} alActualizar={guardar} />
        ))}
      </div>
    </div>
  )
}
