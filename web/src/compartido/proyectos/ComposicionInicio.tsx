import { Outlet, useMatch, useNavigate } from 'react-router'
import { VentanaModal } from '../componentes/VentanaModal'
import { ListaProyectos } from './ListaProyectos'

export function ComposicionInicio() {
  const navegar = useNavigate()
  const ventanaAbierta = Boolean(useMatch('/projects/new'))

  return (
    <>
      <ListaProyectos />
      {ventanaAbierta && (
        <VentanaModal alCerrar={() => navegar('/')}>
          <Outlet />
        </VentanaModal>
      )}
    </>
  )
}
