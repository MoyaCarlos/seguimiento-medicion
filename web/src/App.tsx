import { Link, Outlet } from 'react-router'
import './App.css'

function App() {
  return (
    <div className="app-shell">
      <header className="app-header">
        <Link className="brand" to="/">Métrica</Link>
        <span className="header-label">Seguimiento de proyectos</span>
      </header>
      <main className="app-main">
        <Outlet />
      </main>
    </div>
  )
}

export default App
