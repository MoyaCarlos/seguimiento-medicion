import { Outlet, useParams } from 'react-router'

export function LayoutProyecto() {
  const { id } = useParams()

  return <Outlet key={id} />
}
