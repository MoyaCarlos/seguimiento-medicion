import { useState } from 'react'
import { Link, NavLink, Outlet, useLocation, useParams } from 'react-router'
import { Icono } from './compartido/componentes/Icono'
import { ProveedorRecargaProyectos } from './compartido/proyectos/recargaProyectos'
import { SelectorProyecto } from './compartido/proyectos/SelectorProyecto'
import { MenuUsuario } from './features/autenticacion/MenuUsuario'
import { SesionProvider, useSesion } from './features/autenticacion/SesionContext'
import { navegacion } from './navegacion'
import './App.css'

const SECCION_DE_VOTACION = 'backlog'

const claseDeEnlace = ({ isActive }: { isActive: boolean }) =>
  `barra-lateral__enlace${isActive ? ' barra-lateral__enlace--activo' : ''}`

function Marca() {
  return (
    <Link className="marca" to="/">
      <span className="marca__logo"><Icono nombre="marca" /></span>
      Métrica
    </Link>
  )
}

function RutaDeNavegacion({ pathname, enProyecto }: { pathname: string; enProyecto: boolean }) {
  const tramo = pathname.split('/')[3]
  const seccion = enProyecto
    ? navegacion.find(({ ruta }) => ruta === (tramo === 'votacion' ? SECCION_DE_VOTACION : tramo))
    : undefined

  return (
    <nav className="ruta" aria-label="Ruta de navegación">
      {pathname === '/' ? <span aria-current="page">Proyectos</span> : <Link to="/">Proyectos</Link>}
      {seccion && (
        <>
          <span aria-hidden="true">›</span>
          <span aria-current="page">{seccion.texto}</span>
        </>
      )}
    </nav>
  )
}

function Armazon() {
  const { estado } = useSesion()
  const { id } = useParams()
  const { pathname } = useLocation()
  const [recargas, setRecargas] = useState(0)

  if (estado !== 'autenticado') {
    return (
      <div className="acceso">
        <header className="acceso__barra">
          <Marca />
          <MenuUsuario />
        </header>
        <main className="acceso__contenido">
          <Outlet />
        </main>
      </div>
    )
  }

  return (
    <div className="armazon">
      <aside className="barra-lateral" aria-label="Navegación principal">
        <Marca />
        <NavLink className={claseDeEnlace} to="/" end>
          <Icono nombre="carpeta" />
          Proyectos
        </NavLink>
        {id && (
          <>
            <SelectorProyecto key={recargas} />
            <p className="barra-lateral__rotulo" aria-hidden="true">Seguimiento</p>
            <nav className="barra-lateral__secciones" aria-label="Secciones del proyecto">
              {navegacion.map(({ ruta, texto, icono }) => (
                <NavLink key={ruta} className={claseDeEnlace} to={`/projects/${id}/${ruta}`}>
                  <Icono nombre={icono} />
                  {texto}
                </NavLink>
              ))}
            </nav>
          </>
        )}
      </aside>
      <div className="armazon__columna">
        <header className="barra-superior">
          <RutaDeNavegacion pathname={pathname} enProyecto={Boolean(id)} />
          <MenuUsuario />
        </header>
        <main className="contenido">
          <ProveedorRecargaProyectos value={() => setRecargas((actual) => actual + 1)}>
            <Outlet />
          </ProveedorRecargaProyectos>
        </main>
      </div>
    </div>
  )
}

function App() {
  return (
    <SesionProvider>
      <Armazon />
    </SesionProvider>
  )
}

export default App
