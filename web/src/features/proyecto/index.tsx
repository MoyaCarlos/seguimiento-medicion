import { Outlet, Route, Routes, useParams } from 'react-router'
import { NoEncontrada } from '../../compartido/router/NoEncontrada'
import { Configuracion } from './Configuracion'
import { Editar } from './Editar'
import { Panel } from './Panel'
import './estilos.css'

function PanelLayout() {
  return (
    <>
      <h1 id="panel-title">Panel del proyecto</h1>
      <Outlet />
    </>
  )
}

export function Component() {
  const { id } = useParams()

  return (
    <Routes>
      <Route element={<PanelLayout key={id} />}>
        <Route index element={<Panel />} />
        <Route path="settings" element={<Configuracion />} />
        <Route path="edit" element={<Editar />} />
        <Route path="*" element={<NoEncontrada />} />
      </Route>
    </Routes>
  )
}
