import { useEffect, useState } from 'react'
import { Link, NavLink, Outlet, useLocation, useParams } from 'react-router'
import { Icono } from './compartido/componentes/Icono'
import { Marca } from './compartido/componentes/Marca'
import { ProveedorRecargaProyectos } from './compartido/proyectos/recargaProyectos'
import { SelectorProyecto } from './compartido/proyectos/SelectorProyecto'
import { ProveedorProyectos } from './compartido/proyectos/ProveedorProyectos'
import { useProyectos } from './compartido/proyectos/useProyectos'
import { MenuUsuario } from './features/autenticacion/MenuUsuario'
import { SesionProvider, useSesion } from './features/autenticacion/SesionContext'
import { ajustes, navegacion, seccionDeRuta } from './navegacion'
import './App.css'

const PANTALLA_ANGOSTA = '(max-width: 900px)'

function usePantallaAngosta(): boolean {
  const consulta = () => window.matchMedia?.(PANTALLA_ANGOSTA).matches ?? false
  const [angosta, setAngosta] = useState(consulta)

  useEffect(() => {
    const lista = window.matchMedia?.(PANTALLA_ANGOSTA)
    if (!lista) return
    const alCambiar = () => setAngosta(lista.matches)
    lista.addEventListener('change', alCambiar)
    return () => lista.removeEventListener('change', alCambiar)
  }, [])

  return angosta
}

const claseDeEnlace = ({ isActive }: { isActive: boolean }) =>
  `barra-lateral__enlace${isActive ? ' barra-lateral__enlace--activo' : ''}`

function Separador() {
  return (
    <span className="ruta__separador" aria-hidden="true">
      <Icono nombre="chevronDerecha" />
    </span>
  )
}

function RutaDeNavegacion({ pathname, id }: { pathname: string; id?: string }) {
  const proyectos = useProyectos()
  const seccion = id ? seccionDeRuta(pathname) : undefined
  const nombre =
    (proyectos.estado === 'listo' && proyectos.proyectos.find((proyecto) => String(proyecto.id) === id)?.nombre) ||
    `Proyecto ${id}`

  return (
    <nav className="ruta" aria-label="Ruta de navegación">
      {pathname === '/' ? <span aria-current="page">Proyectos</span> : <Link to="/">Proyectos</Link>}
      {pathname === '/projects/new' && (
        <>
          <Separador />
          <span aria-current="page">Nuevo proyecto</span>
        </>
      )}
      {id && (
        <>
          <Separador />
          <Link to={`/projects/${id}/proyecto`}>{nombre}</Link>
        </>
      )}
      {seccion && (
        <>
          <Separador />
          <span aria-current="page">{seccion}</span>
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
  const conPestanas = usePantallaAngosta() && Boolean(id)

  if (estado !== 'autenticado') {
    return (
      <main className="acceso">
        <Outlet />
      </main>
    )
  }

  return (
    <ProveedorProyectos recarga={recargas}>
      <div
        className={`armazon${menuAbierto ? ' armazon--menu-abierto' : ''}${conPestanas ? ' armazon--con-pestanas' : ''}`}
        onKeyDown={(evento) => {
          if (evento.key === 'Escape') cerrarMenu()
        }}
      >
        <aside
          id="barra-lateral"
          className="barra-lateral"
          aria-label="Navegación principal"
          onClick={(evento) => {
            if ((evento.target as HTMLElement).closest('a, [role="option"], .selector-proyecto__nuevo')) cerrarMenu()
          }}
        >
          <button type="button" className="barra-lateral__cerrar" aria-label="Cerrar menú" onClick={cerrarMenu}>
            <Icono nombre="cerrar" />
          </button>
          <Marca />
          <Link
            className={claseDeEnlace({ isActive: pathname === '/' || pathname === '/projects/new' })}
            to="/"
            aria-current={pathname === '/' || pathname === '/projects/new' ? 'page' : undefined}
          >
            <Icono nombre="carpeta" />
            Proyectos
          </Link>
          {id && (
            <>
              <SelectorProyecto />
              <p className="barra-lateral__rotulo" aria-hidden="true">Seguimiento</p>
              <nav className="barra-lateral__secciones" aria-label="Secciones del proyecto">
                {navegacion.map(({ ruta, texto, icono }) => (
                  <NavLink key={ruta} className={claseDeEnlace} to={`/projects/${id}/${ruta}`} end={ruta === 'proyecto'}>
                    <Icono nombre={icono} />
                    {texto}
                  </NavLink>
                ))}
              </nav>
              <p className="barra-lateral__rotulo" aria-hidden="true">Ajustes</p>
              <nav className="barra-lateral__secciones" aria-label="Ajustes del proyecto">
                {ajustes.map(({ ruta, texto, icono }) => (
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
            <RutaDeNavegacion pathname={pathname} id={id} />
            <MenuUsuario />
          </header>
          <main className="contenido">
            <ProveedorRecargaProyectos value={() => setRecargas((actual) => actual + 1)}>
              <Outlet />
            </ProveedorRecargaProyectos>
          </main>
        </div>
        {conPestanas && (
          <nav className="pestanas" aria-label="Pestañas del proyecto">
            {navegacion.map(({ ruta, texto, icono }) => (
              <NavLink
                key={ruta}
                className={({ isActive }) => `pestanas__enlace${isActive ? ' pestanas__enlace--activa' : ''}`}
                to={`/projects/${id}/${ruta}`}
                end={ruta === 'proyecto'}
                aria-label={texto}
              >
                <Icono nombre={icono} />
                <span>{texto === 'Historial de Sprints' ? 'Historial' : texto}</span>
              </NavLink>
            ))}
          </nav>
        )}
      </div>
    </ProveedorProyectos>
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
