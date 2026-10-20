import { AppRoutes } from './compartido/router/AppRoutes'
import './App.css'

function App() {
  return (
    <div className="app-shell">
      <header className="app-header">
        <a className="brand" href="/projects/1">Métrica</a>
        <span className="header-label">Seguimiento de proyectos</span>
      </header>
      <main className="app-main">
        <AppRoutes />
      </main>
    </div>
  )
}

export default App
