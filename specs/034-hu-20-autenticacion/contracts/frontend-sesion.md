# Contrato: sesión en el frontend

## Enganche del cliente compartido (`compartido/api/sesionHttp.ts`)

```ts
type EnganceSesion = {
  obtenerToken: () => string | null
  alNoAutenticado: () => void
}
export function configurarSesionHttp(enganche: EnganceSesion | null): void
```

- `client.ts` agrega `Authorization: Bearer <token>` cuando `obtenerToken()` devuelve un
  valor.
- Ante un 401 llama a `alNoAutenticado()` y relanza el `ApiError` solo si la petición llevó
  token y su ruta no es `/auth/ingreso` ni `/auth/registro`.
- Sin enganche configurado se comporta como hasta ahora.

## Contexto de sesión (`features/autenticacion`)

```ts
type EstadoSesion = 'cargando' | 'anonimo' | 'autenticado'

type Sesion = {
  estado: EstadoSesion
  usuario: { id: number; nombre: string } | null
  ingresar(nombre: string, contrasena: string): Promise<void>
  registrar(nombre: string, contrasena: string): Promise<void>
  salir(): Promise<void>
}
export function SesionProvider(props: { children: ReactNode }): JSX.Element
export function useSesion(): Sesion
```

- `ingresar` y `registrar` guardan el token, fijan el usuario y navegan a `state.desde`
  (o `/`). Relanzan el `ApiError` para que la pantalla muestre el mensaje y el campo.
- Al montar, si hay token guardado: `estado = 'cargando'` y se consulta `GET /auth/yo`.
  Éxito → `autenticado`. 401 → se descarta el token y se trata como sesión vencida.
- `alNoAutenticado`: borra el token, `estado = 'anonimo'`, navega a `/ingresar` con
  `state = { desde: <ruta actual>, vencida: true }`.
- `salir`: pide `POST /auth/salir` (un error se ignora), borra el token, `anonimo`, navega
  a `/ingresar`.

## Almacenamiento (`almacenamientoToken.ts`)

`leerToken()`, `guardarToken(t)`, `borrarToken()`: clave única, cada acceso dentro de
`try/catch`. Si falla, no lanza: devuelve `null` o no hace nada, y el token vive solo en
memoria.

## Componentes y rutas

| Pieza | Comportamiento |
|---|---|
| `Ingresar` (`/ingresar`) | formulario nombre + contraseña; muestra «Tu sesión venció» si `state.vencida`; error del servidor en el formulario sin borrar lo escrito; enlace a registro |
| `Registro` (`/registro`) | formulario nombre + contraseña; el error de campo (`nombre`, `contrasena`) se muestra junto a su campo; enlace a ingreso |
| `RutaProtegida` | `cargando` → indicador de carga; `anonimo` → `<Navigate to="/ingresar" state={{ desde }} replace />`; `autenticado` → `<Outlet />` |
| `MenuUsuario` | en el encabezado: con sesión, el nombre y «Cerrar sesión»; sin sesión, enlaces a ingresar y registrarse |

Rutas: `/ingresar` y `/registro` públicas; el resto de las rutas de la SPA bajo
`RutaProtegida` (decidido, [research.md](../research.md) D12).

## Convención de pruebas

Nombre de cada prueba con el escenario del spec (por ejemplo «Escenario 6: sesión vencida
lleva a ingresar y vuelve a la pantalla de origen»), con Vitest y Testing Library; `fetch`
simulado como en `client.test.ts`.
