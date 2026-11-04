import { useParams } from 'react-router'
import { PantallaDefectos } from './PantallaDefectos'

export function Component() {
  const { id = '' } = useParams()

  return <PantallaDefectos key={id} proyectoId={id} />
}
