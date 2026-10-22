# Research: Base del frontend

## 1. Reenvío de llamadas (CORS)
- **Decision**: Conservar el proxy de Vite actual: `/backlog`, `/projects` y `/sprints` hacia `http://localhost:8080`, con `bypass` que devuelve `/index.html` cuando el navegador pide HTML (`Accept: text/html`).
- **Rationale**: Ya funciona, no toca Go y permite que `/projects/:id/...` sea a la vez una ruta de pantalla (navegación del navegador) y una ruta de la API (llamada del cliente). El cliente envía `Accept: application/json`, que es lo que activa el reenvío al backend.
- **Alternatives**: (a) Prefijo `/api` con reescritura: más claro pero cambia el cliente y el proxy que ya anda; descartado por el equipo. (b) Cabecera CORS en Go: cambio de backend, fuera de alcance.
- **Riesgos**: (1) Un `fetch` que no envíe `Accept: application/json` recibiría `index.html`: el cliente compartido siempre lo envía, y el README lo aclara. (2) Una ruta de pantalla nueva con un primer segmento que no sea `projects`, `backlog` ni `sprints` no pasa por el proxy y no necesita nada. (3) Para producción (Docker/Azure, probable Sprint 4) habrá que servir la SPA y la API de forma equivalente; fuera de alcance.

## 2. Ruteo
- **Decision**: `react-router` con `createBrowserRouter` (rutas como objetos) y `lazy` por feature. Aprobado en clarify.
- **Rationale**: Estándar; soporta `:id` en la URL, rutas anidadas, `NavLink` con pantalla activa y prioriza la ruta más específica. Una ruta = un objeto de una línea.
- **Alternatives**: Router propio (el actual): obliga a editar código compartido por cada pantalla; es lo que se reemplaza.
- **Por verificar al instalar**: versión de `react-router` compatible con React 19; el nombre del paquete (`react-router`) y la API (`createBrowserRouter`, `RouterProvider`, `NavLink`) se confirman con el build.

## 3. Pruebas
- **Decision**: Vitest (instalado) + `@testing-library/react` + `@testing-library/user-event` + `jsdom`. Entorno jsdom configurado en el bloque `test` de `vitest.config.ts` (la rama ya lo tiene; con los dos archivos presentes Vitest ignora el bloque `test` de `vite.config.ts`).
- **Rationale**: Reutiliza el plugin de React y la configuración de Vite; Testing Library prueba lo que ve la persona usuaria (etiqueta asociada, botón deshabilitado, aviso visible).
- **Alternatives**: Entorno jsdom por archivo con comentario `@vitest-environment`: más ruido. Si `Response`/`fetch` simulados fallan bajo jsdom en los tests del cliente, esos archivos pasan a entorno `node`.
- **Nota**: `fetch` se simula con `vi.stubGlobal`; sin librería de mocks de red.

## 4. Cliente de la API
- **Decision**: Una función base `request(path, método, cuerpo?)` y tres atajos `get`, `post`, `put`. Reglas: si `fetch` rechaza → `ConnectionError`; si la respuesta no es OK → intenta leer `{campo?, mensaje?}`; sin mensaje interpretable → `ApiError(estado, "La API respondió con estado N")`; si es OK y el cuerpo no es JSON → `ApiError(estado, mensaje genérico)`; cuerpo vacío → `undefined`.
- **Rationale**: Cubre FR-005 a FR-009b y deja los endpoints en las features. Hoy `parseResponseBody` lanza "respuesta inválida" aun en errores y oculta el estado real: se corrige leyendo el cuerpo de forma tolerante.
- **Alternatives**: Una función por endpoint en `compartido/` (acopla la base a la API, contradice FR-005); un tipo de error por código de estado (más clases sin necesidad).

## 5. Impedir que una feature importe a otra (FR-020)
- **Decision**: Probar primero una regla de oxlint (`no-restricted-imports` con `overrides` sobre `src/features/**`). Si oxlint no puede expresar "otra carpeta de `features/`" sin listar cada feature en la configuración, se usa un test de Vitest que recorre `src/features/**` y falla si un `import` apunta a otra carpeta de `features/`.
- **Rationale**: Un test corre en `npm test` y en CI sin dependencias nuevas, y no obliga a editar la configuración al sumar una feature. La regla de lint es preferible solo si no requiere mantenimiento por feature.
- **Alternatives**: Dejarlo a la revisión del PR (frágil, ya figura como respaldo en la spec); plugin de lint extra (otra dependencia, no aprobada).

## 6. Estilos y placeholders
- **Decision**: Se conserva el CSS actual (`index.css`, `App.css`). Cada placeholder usa las clases ya existentes de `placeholder`; el `FeaturePlaceholder` común deja `compartido/`.
- **Rationale**: Estilos elaborados están fuera de alcance; los placeholders son descartables cuando llegue cada pantalla.

## 7. Resueltos por clarify
Error sin `campo` → aviso arriba del formulario; librerías aprobadas; "prototipado" = solo código base.
