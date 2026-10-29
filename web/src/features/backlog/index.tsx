import { Route, Routes } from 'react-router'
import { NoEncontrada } from '../../compartido/router/NoEncontrada'
import { Component as Crear } from './crear'
import { Listado } from './Listado'

export function Component() {
  return (
    <Routes>
      <Route index element={<Listado />} />
      <Route path="nueva" element={<Crear />} />
      <Route path="*" element={<NoEncontrada />} />
    </Routes>
  )
}
