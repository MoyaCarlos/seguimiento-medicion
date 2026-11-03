import { Outlet } from 'react-router'
import { ListaProyectos } from './ListaProyectos'

export function ComposicionInicio() {
  return (
    <>
      <Outlet />
      <ListaProyectos />
    </>
  )
}
