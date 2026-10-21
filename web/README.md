# Frontend (web/)

Esta carpeta es la parte visual del proyecto: lo que se abre en el navegador.
Habla con el backend (la API que corre en el puerto 8080) para leer y guardar
datos.

Esta base deja todo listo para que las cuatro pantallas del Sprint 2 (HU-04,
HU-01, HU-05 y HU-13) se hagan **en paralelo**, cada una en su propia carpeta,
sin pisarse.

## Cómo correrlo

Necesitás Node.js 26 y el backend levantado. Todo el equipo usa la misma
versión de Node, y está fijada en el archivo `web/.nvmrc`. Con otra versión,
`npm install` se detiene con un error que indica la versión requerida
(`Required: {"node":">=26"}`), en lugar de fallar más adelante de forma confusa.

Para tener la versión correcta, desde esta carpeta (`web/`):

- Con **nvm**: `nvm install` y después `nvm use`.
- Con **fnm**: `fnm install` y después `fnm use`.

Los dos leen la versión del archivo `.nvmrc`. Podés comprobarlo con
`node --version`, que tiene que empezar con `v26`.

Con Node 26 listo:

1. Levantar el backend, desde la raíz del repositorio:

   ```bash
   go run ./cmd/api
   ```

   La terminal muestra `listening on :8080`. Para cargar un proyecto y un Sprint
   de ejemplo, usá los comandos de
   [`docs/curl-frontend-sprint-2.md`](../docs/curl-frontend-sprint-2.md) y
   anotá los `id` que devuelve.

2. Levantar el frontend, desde esta carpeta (`web/`):

   ```bash
   npm install
   npm run dev
   ```

   La terminal muestra la dirección, normalmente `http://localhost:5173`.

3. Abrir `http://localhost:5173/projects/<id>` con el `id` del proyecto que
   cargaste. La API todavía no tiene un listado de proyectos, así que hay que
   conocer el `id`. La página `/` solo explica esto.

### Otros comandos

| Comando | Para qué sirve |
|---|---|
| `npm run build` | Revisa los tipos y compila la aplicación. Termina sin errores si todo está bien. |
| `npm run lint` | Revisa el código con oxlint. No imprime nada si está todo bien. |
| `npm test` | Corre todas las pruebas automáticas (Vitest). |

Antes de abrir un Pull Request tienen que pasar los tres.

## Cómo se conecta con la API (el proxy)

El navegador no puede llamar directamente al backend (el backend no tiene CORS,
el permiso que necesita un navegador para llamar a otra dirección). Por eso, en
desarrollo, el servidor de Vite hace de intermediario. Está configurado en
`vite.config.ts` para las direcciones que empiezan con `/backlog`, `/projects` y
`/sprints`, y decide según lo que pide quien llama:

- Si la llamada pide **JSON** (`Accept: application/json`), la reenvía al backend
  en `http://localhost:8080`.
- Si la llamada viene del **navegador abriendo una página** (`Accept: text/html`),
  devuelve `index.html`. Por eso se puede recargar `/projects/42/backlog` sin
  recibir un 404 del backend, aunque esa dirección también exista en la API.

Consecuencia: **todas las llamadas a la API se hacen con el cliente de
`src/compartido/api/client.ts`** (ver más abajo), que siempre envía
`Accept: application/json`. Si hacés un `fetch` a mano sin ese encabezado,
recibirías `index.html` en lugar de datos.

Si el backend suma un grupo de direcciones nuevo (por ejemplo `/reportes`), hay
que agregarlo en `vite.config.ts`. Ese archivo solo lo modifica esta base.

## Cómo está organizado

```text
src/
  main.tsx              Arranque: monta el ruteo
  App.tsx               Marco de la página (encabezado y contenido)
  rutas.tsx             Registro de rutas: una línea por pantalla
  navegacion.ts         Registro del menú: una línea por pantalla
  compartido/           Código común. Solo lo modifica esta base
    api/                Cliente de la API (client.ts) y sus pruebas
    componentes/        Button, Field, ErrorMessage, LoadingIndicator
    router/             Marco del proyecto (menú), "página no encontrada" e inicio
  features/             Una carpeta por historia
    proyecto/           HU-04 (hoy solo un marcador de lugar)
    backlog/            HU-01
    sprints/            HU-05
    estado/             HU-13
  setupTests.ts         Limpia la pantalla entre pruebas
```

Hoy cada carpeta de `features/` tiene un `index.tsx` que solo muestra un
marcador de lugar con el título y el proyecto actual. Cada historia lo
reemplaza por su pantalla real.

### Las reglas que evitan los choques

- **Una carpeta por historia:** `src/features/<nombre>/`.
- **Ninguna pantalla importa a otra.** Lo que sea común va en
  `src/compartido/`. Hay una prueba (`src/features/aislamiento.test.ts`) que
  falla si una carpeta de `features/` importa algo de otra.
- **`src/compartido/` y `package.json` los modifica solo esta base.** Si una
  pantalla necesita algo común nuevo o una librería nueva, se pide y se agrega
  en la base, no en el PR de la pantalla.
- **El proyecto actual va en la dirección** (`/projects/:id/...`), no en una
  variable compartida entre pantallas. Una pantalla lo lee con `useParams`.
- **El panel del proyecto (`features/proyecto`) es el único punto de contacto
  entre HU-04 y HU-13.** HU-13 no usa nada de HU-04 ni al revés; si hay que
  juntarlas, se hace al final dentro de ese panel. Hoy el panel puede apuntar
  a HU-13 con un enlace a `/projects/<id>/status`. Si necesitara mostrar algo de
  HU-13 dentro de su propia pantalla, la prueba de aislamiento lo rechaza: hay
  que acordarlo con el equipo y ajustarla en la base.

## Cómo agregar una pantalla nueva

Ejemplo con una pantalla llamada `ejemplo`, que queda en
`/projects/<id>/ejemplo`.

**1. Crear la carpeta con su `index.tsx`**, que tiene que exportar una función
llamada `Component` (con ese nombre exacto):

```tsx
// src/features/ejemplo/index.tsx
import { useParams } from 'react-router'

export function Component() {
  const { id } = useParams()

  return <h1>Pantalla de ejemplo del proyecto {id}</h1>
}
```

**2. Agregar una línea en `src/rutas.tsx`**, junto a las otras pantallas, dentro
de la lista `children` de `projects/:id`:

```tsx
{ path: 'ejemplo/*', lazy: () => import('./features/ejemplo') },
```

**3. Agregar una línea en `src/navegacion.ts`** para que aparezca en el menú:

```ts
{ ruta: 'ejemplo', texto: 'Ejemplo' },
```

Eso es todo: fuera de tu carpeta cambian solo esas dos líneas. Si dos personas
agregan su línea en el mismo lugar y Git marca un conflicto, se resuelve
dejando **las dos líneas**.

El `/*` del final permite que tu pantalla tenga sub-pantallas propias
(por ejemplo `/projects/<id>/ejemplo/editar`), que declarás dentro de tu
carpeta con `Routes` y `Route` de `react-router` sin tocar el registro. Una
dirección que no corresponde a ninguna pantalla muestra "Página no
encontrada".

Después de crear tu pantalla corré `npm run build`, `npm run lint` y
`npm test`. Si tu pantalla necesita pruebas, guardalas junto a ella dentro de
tu carpeta.

## Cómo llamar a la API

El cliente está en `src/compartido/api/client.ts` y tiene tres funciones, que
envían y leen JSON:

| Función | Qué hace |
|---|---|
| `get<T>(ruta)` | Pide datos |
| `post<T>(ruta, cuerpo)` | Envía datos nuevos |
| `put<T>(ruta, cuerpo)` | Envía datos para actualizar |

El cliente no conoce los endpoints: cada pantalla arma sus llamadas en su
propia carpeta. Por ejemplo:

```ts
// src/features/backlog/api.ts
import { get, post } from '../../compartido/api/client'

type Historia = { id: number; titulo: string }

export const crearHistoria = (proyectoId: number, titulo: string) =>
  post<Historia>('/backlog', { proyecto_id: proyectoId, titulo, descripcion: '', prioridad: 'M' })

export const obtenerProyecto = (id: string) => get(`/projects/${id}`)
```

Las rutas son relativas, como `/projects/42` (sin `http://localhost:8080`); el
proxy se ocupa del resto. Los endpoints que existen hoy son:
`POST /backlog`, `POST /projects`, `GET` y `PUT /projects/{id}`,
`GET /projects/{id}/status`, `POST` y `GET /projects/{id}/members`,
`POST /sprints`, `POST /sprints/{id}/iniciar` y `POST /sprints/{id}/cerrar`.
**No existen** el listado de proyectos, el de Sprints ni el del backlog.

### Qué pasa cuando algo sale mal

Si todo va bien, la función devuelve los datos (o `undefined` si la respuesta
no trae cuerpo). Si no, **lanza un error**, que atrapás con `try/catch`:

| Qué pasó | Qué se lanza | Qué trae |
|---|---|---|
| El backend respondió con un error (400, 404, 409, etc.) | `ApiError` | `status`, `message` y, si el error es de un campo del formulario, `field` |
| El backend respondió algo que no se puede leer (por ejemplo un 500 con una página en vez de JSON) | `ApiError` | `status` y el mensaje "La API respondió con estado N" |
| No hubo respuesta (backend apagado o sin red) | `ConnectionError` | `message` |

Se distinguen con `instanceof`: `error instanceof ApiError` o
`error instanceof ConnectionError`.

### Cómo mostrar los errores y la carga en un formulario

Las cuatro piezas están en `src/compartido/componentes/`:

| Pieza | Para qué |
|---|---|
| `Field` | Campo con su etiqueta (`label`). Si le pasás `error`, lo muestra junto al campo |
| `Button` | Botón (`variant`: `primary`, `secondary` o `quiet`). Se deshabilita con `disabled` |
| `ErrorMessage` | Aviso de error general (`message`) |
| `LoadingIndicator` | Texto "Cargando…" |

La regla:

- Un error **con `field`** se muestra **junto a ese campo**: pasá `error.message`
  al `Field` cuyo `name` coincide con `error.field`.
- Un error **sin `field`** (y un `ConnectionError`) se muestra con
  `ErrorMessage` **arriba del formulario**, y queda ahí hasta el próximo envío.
- Mientras la operación está en curso, mostrá `LoadingIndicator` y deshabilitá
  el botón de envío.

El modelo para copiar es `FormularioDePrueba`, en
`src/compartido/componentes/componentes.test.tsx` (acá, simplificado):

```tsx
<form onSubmit={alEnviar}>
  {error && !error.field && <ErrorMessage message={error.message} />}
  <Field name="nombre" label="Nombre" error={error?.field === 'nombre' ? error.message : undefined} />
  {enCurso && <LoadingIndicator />}
  <Button type="submit" disabled={enCurso}>Crear</Button>
</form>
```

Ese ejemplo guarda el error en una variable de estado y la borra al empezar
cada envío. Está pensado para `ApiError`; si tu pantalla también quiere
mostrar un `ConnectionError`, agregá esa comprobación con `instanceof`.

## Qué queda fuera de esta base

Esta base no incluye las pantallas reales, autenticación ni estilos
elaborados. Tampoco las llamadas que hoy no existen en la API (listados de
proyectos, Sprints y backlog); cuando existan, se suman dentro de la carpeta de
la pantalla correspondiente, sin cambiar el código compartido.
