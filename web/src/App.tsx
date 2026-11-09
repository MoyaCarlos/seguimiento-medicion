import { Link, Outlet } from 'react-router'
import { MenuUsuario } from './features/autenticacion/MenuUsuario'
import { SesionProvider } from './features/autenticacion/SesionContext'
import './App.css'

function App() {
  return (
    <SesionProvider>
      <div className="app-shell">
        <header className="app-header">
          <Link className="brand" to="/">Métrica</Link>
          <span className="header-label">Seguimiento de proyectos</span>
          <MenuUsuario />
        </header>
        <main className="app-main">
          <Outlet />
        </main>
      </div>
    </SesionProvider>
  )
}

export default App
