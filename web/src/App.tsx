import { useState } from 'react'
import { Link, NavLink, Outlet, useLocation, useParams } from 'react-router'
import { Icono } from './compartido/componentes/Icono'
import { Marca } from './compartido/componentes/Marca'
import { ProveedorRecargaProyectos } from './compartido/proyectos/recargaProyectos'
import { SelectorProyecto } from './compartido/proyectos/SelectorProyecto'
import { MenuUsuario } from './features/autenticacion/MenuUsuario'
import { SesionProvider, useSesion } from './features/autenticacion/SesionContext'
import { navegacion } from './navegacion'
import './App.css'

const SECCION_DE_VOTACION = 'backlog'

const claseDeEnlace = ({ isActive }: { isActive: boolean }) =>
  `barra-lateral__enlace${isActive ? ' barra-lateral__enlace--activo' : ''}`

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
  const [menuAbierto, setMenuAbierto] = useState(false)
  const cerrarMenu = () => setMenuAbierto(false)

  if (estado !== 'autenticado') {
    return (
      <main className="acceso">
        <Outlet />
      </main>
    )
  }

  return (
    <div
      className={`armazon${menuAbierto ? ' armazon--menu-abierto' : ''}`}
      onKeyDown={(evento) => {
        if (evento.key === 'Escape') cerrarMenu()
      }}
    >
      <aside
        id="barra-lateral"
        className="barra-lateral"
        aria-label="Navegación principal"
        onClick={(evento) => {
          if ((evento.target as HTMLElement).closest('a')) cerrarMenu()
        }}
      >
        <button type="button" className="barra-lateral__cerrar" aria-label="Cerrar menú" onClick={cerrarMenu}>
          <Icono nombre="cerrar" />
        </button>
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
      <div className="armazon__fondo" onClick={cerrarMenu} />
      <div className="armazon__columna">
        <header className="barra-superior">
          <button
            type="button"
            className="barra-superior__menu"
            aria-label="Abrir menú"
            aria-expanded={menuAbierto}
            aria-controls="barra-lateral"
            onClick={() => setMenuAbierto(true)}
          >
            <Icono nombre="menu" />
          </button>
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
