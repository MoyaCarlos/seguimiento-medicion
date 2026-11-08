import { Route, Routes } from 'react-router'
import { NoEncontrada } from '../../compartido/router/NoEncontrada'
import { Detalle } from './Detalle'
import { Lista } from './Lista'
import './historial.css'

export function Component() {
  return (
    <Routes>
      <Route index element={<Lista />} />
      <Route path=":sprintId" element={<Detalle />} />
      <Route path="*" element={<NoEncontrada />} />
    </Routes>
  )
}
