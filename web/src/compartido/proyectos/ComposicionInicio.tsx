import { Outlet } from 'react-router'
import { ListaProyectos } from './ListaProyectos'

export function ComposicionInicio() {
  return (
    <div className="inicio-composicion">
      <Outlet />
      <ListaProyectos />
    </div>
  )
}
