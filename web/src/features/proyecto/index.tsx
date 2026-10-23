import { Route, Routes } from 'react-router'
import { NoEncontrada } from '../../compartido/router/NoEncontrada'
import { Configuracion } from './Configuracion'
import { Editar } from './Editar'
import { Panel } from './Panel'
import './estilos.css'

export function Component() {
  return (
    <Routes>
      <Route index element={<Panel />} />
      <Route path="settings" element={<Configuracion />} />
      <Route path="edit" element={<Editar />} />
      <Route path="*" element={<NoEncontrada />} />
    </Routes>
  )
}
