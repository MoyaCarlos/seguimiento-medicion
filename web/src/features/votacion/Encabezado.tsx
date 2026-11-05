import { Link } from 'react-router'
import { Button } from '../../compartido/componentes/Button'

type Props = {
  proyectoId: string
  historiaId: number
  deshabilitado: boolean
  alActualizar: () => void
}

export function Encabezado({ proyectoId, historiaId, deshabilitado, alActualizar }: Props) {
  return (
    <div className="page-header votacion-encabezado">
      <div>
        <p className="eyebrow">Planning Poker</p>
        <h1 id="votacion-title">Votación de la historia #{historiaId}</h1>
        <p className="project-context">
          Proyecto actual: <strong>{proyectoId}</strong>
        </p>
      </div>
      <div className="votacion-encabezado__acciones">
        <Link className="button button--quiet" to={`/projects/${proyectoId}/backlog`}>
          Volver al Product Backlog
        </Link>
        <Button type="button" variant="secondary" onClick={alActualizar} disabled={deshabilitado}>
          Actualizar
        </Button>
      </div>
    </div>
  )
}
