# Implementation Plan: Navegar entre proyectos desde la interfaz

**Branch**: `feature/148-navegacion-proyectos` | **Date**: 2026-10-27 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/022-navegar-entre-proyectos/spec.md`

**Note**: Plan de implementación de la historia #148 (cierra #141). Solo frontend; la API y el
esquema de la base no cambian.

## Summary

Agregar navegación entre proyectos reutilizando el listado existente `GET /projects`: un selector
en el encabezado de cada pantalla de proyecto que conserva la sección actual al cambiar de
proyecto, un enlace permanente a "Nuevo proyecto" y la lista de proyectos en la pantalla de
inicio junto al formulario de creación. La lógica pura de ruta y etiqueta, el acceso a la API y
el estado de la lista viven en `web/src/compartido/proyectos/`; las pantallas consumen piezas
compartidas y no se toca `features/proyecto/crear.tsx` ni su prueba. Se elimina el `Inicio.tsx`
sin uso y se cierra #141.

## Technical Context

**Language/Version**: TypeScript `~6.0.2`, React `19.2.8`, React Router `8.5.0`, Vite `8.3.0`,
Node.js `>=24.15.0` (fijado en `web/.nvmrc` y `web/package.json`).

**Primary Dependencies**: `react-router` (navegación y `Outlet`), cliente existente
`web/src/compartido/api/client.ts`, `LoadingIndicator` y `ErrorMessage` compartidos. No se
agregan dependencias.

**Storage**: N/A. Solo lectura del endpoint existente `GET /projects`; no hay persistencia nueva
ni cambios de esquema.

**Testing**: Vitest `4.x` + Testing Library `16.x` + jsdom `30.x`; `fetch` simulado con
`vi.stubGlobal`; `npm test`, `npm run lint` y `npm run build` desde `web/`. La suite existente de
`web/src/rutas.test.tsx` se mantiene en verde; `web/src/features/proyecto/crear.test.tsx` no se
modifica.

**Target Platform**: SPA en el navegador, servida por Vite; en desarrollo el proxy de Vite lleva
`/projects` al backend en `http://localhost:8080`.

**Project Type**: Frontend web (React + Vite) dentro de un repositorio con backend Go. Esta
historia no modifica el backend.

**Performance Goals**: Un `GET /projects` por consumidor montado (selector o lista; en desarrollo
React StrictMode puede repetir el efecto); funciones puras de costo lineal en la cantidad de
tramos o proyectos; sin polling ni recargas automáticas. La lista se renderiza en una sola pasada.

**Constraints**: Sin dependencias nuevas ni cambios de API/esquema. Archivos permitidos:
`LayoutProyecto.tsx`, `rutas.tsx`, `rutas.test.tsx`, archivos nuevos en
`web/src/compartido/proyectos/` y eliminación de `Inicio.tsx`; no se tocan
`features/proyecto/crear.tsx` ni `crear.test.tsx`. Respetar RED/GREEN/REFACTOR y el orden de
commits pedido. Los escenarios BDD se cubren por equivalencia con pruebas de componente.

**Scale/Scope**: SPA pequeña, N proyectos sin paginación (se muestran en el orden de la API), una
pantalla de inicio, un encabezado de proyecto, cinco archivos de producción nuevos y cuatro
archivos de prueba nuevos. No hay roles ni permisos en esta historia.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principio | Resultado | Evidencia |
|---|---|---|
| I. Test-First estricto (NON-NEGOTIABLE) | PASS | Cada comportamiento tiene su prueba antes de la implementación; un commit por fase RED/GREEN y REFACTOR solo si hay limpieza. El primer commit es infraestructura de prueba (mock de `GET /projects`) para que la suite existente siga verde antes y después del cambio. Merge con merge commit, nunca squash. |
| II. Specification-Driven Development | PASS | Spec de la historia en `specs/022-navegar-entre-proyectos/spec.md`; este plan y los artefactos de Fase 0/1 son su salida. No se implementa nada fuera de la spec. |
| III. Behavior-Driven Development | PASS (equivalencia documentada) | La spec 022 establece que los escenarios de aceptación Given-When-Then se automatizan como pruebas de componente con Vitest + Testing Library. No hay comportamiento de backend nuevo que requiera Godog. |
| IV. Arquitectura Clean/Hexagonal | PASS | El backend y `internal/domain` no se tocan. El frontend separa lógica pura (`ruta.ts`), acceso a datos (`api.ts`), estado de UI (`useProyectos.ts`) y presentación (selector y lista). |
| V. Clean Code, SOLID, KISS, YAGNI, DRY | PASS | Un hook y un adaptador de API reutilizados por selector y lista; una sola función de ruta y una sola de etiqueta; sin contexto global, store, polling ni retry no pedidos. |
| VI. Stack tecnológico obligatorio | PASS | React + Vite + Vitest/Testing Library ya presentes; sin dependencias nuevas; backend Go intacto. |

**Post-Design Re-check**: PASS. El diseño de Fase 1 no agrega dependencias, no modifica contratos
de API ni el esquema, y conserva la separación entre lógica pura, estado y presentación.

## Project Structure

### Documentation (this feature)

```text
specs/022-navegar-entre-proyectos/
├── plan.md              # Este archivo
├── research.md          # Fase 0: decisiones y alternativas
├── data-model.md        # Fase 1: entidades y estados de UI
├── quickstart.md        # Fase 1: guía de validación
├── contracts/
│   └── ui.md            # Fase 1: contrato de UI y consumo de GET /projects
└── checklists/
    └── requirements.md  # Validación de la spec
```

### Source Code (repository root)

```text
web/src/
├── compartido/
│   ├── api/
│   │   └── client.ts                    # Reutilizado, sin cambios
│   ├── componentes/                     # LoadingIndicator y ErrorMessage reutilizados
│   ├── proyectos/                       # Nuevo módulo compartido
│   │   ├── api.ts                       # ProyectoResumen + listarProyectos()
│   │   ├── ruta.ts                      # rutaDeSeccion() + etiquetaProyecto()
│   │   ├── ruta.test.ts
│   │   ├── useProyectos.ts              # Hook de estados cargando/listo/error
│   │   ├── useProyectos.test.tsx
│   │   ├── SelectorProyecto.tsx
│   │   ├── SelectorProyecto.test.tsx
│   │   ├── ListaProyectos.tsx
│   │   └── ListaProyectos.test.tsx
│   └── router/
│       ├── LayoutProyecto.tsx           # Modificado: monta el selector y reinicia por :id
│       └── Inicio.tsx                   # Eliminado (sin uso, #141)
├── features/
│   └── proyecto/
│       ├── crear.tsx                    # SIN CAMBIOS
│       └── crear.test.tsx               # SIN CAMBIOS
├── rutas.tsx                            # Modificado: composición del inicio
└── rutas.test.tsx                       # Modificado: fetch simulado + pruebas de integración
```

**Structure Decision**: El selector y la lista son navegación transversal, no una pantalla de una
feature, así que viven en `compartido/proyectos/`. `rutas.tsx` compone el inicio con un layout
local que renderiza el formulario (hijo index lazy existente) y `ListaProyectos`; así no se toca
`features/proyecto/crear.tsx` y se mantiene el aislamiento entre features. No se agrega CSS
nuevo: se reutilizan clases base existentes (`page-stack`, `field`, `button`, `loading`,
`message--error`) y el marcado semántico.

## Design Decisions

1. **Composición del inicio sin tocar `crear.tsx`.** En `rutas.tsx` se define un componente local
   `ComposicionInicio` que renderiza `<Outlet />` (el formulario lazy actual) seguido de
   `<ListaProyectos />`. La ruta índice pasa a ser hija de ese layout pathless. Alternativa
   descartada: modificar `crear.tsx`, porque sus pruebas afirman la cantidad exacta de llamadas a
   `fetch` y agregar el `GET /projects` allí las rompería.
2. **Módulo compartido, un solo hook.** `useProyectos` expone una unión discriminada
   `cargando | listo | error`, usa `listarProyectos()` y evita actualizar estado tras el
   desmontaje con una bandera `vigente`. Alternativa descartada: contexto global o store, porque
   selector y lista nunca conviven y sería abstracción innecesaria (YAGNI).
3. **Ruta pura y etiqueta pura.** `rutaDeSeccion(pathname, idNuevo)` toma el primer tramo después
   del id (o `proyecto` si no hay tramo) y arma `/projects/<idNuevo>/<sección>`.
   `etiquetaProyecto({ id, nombre })` devuelve `<nombre> (#<id>)`. Ambas se prueban con Vitest
   antes de implementar.
4. **Selector nativo accesible.** `<select>` con etiqueta accesible "Proyecto", opciones con la
   etiqueta de cada proyecto, `value` igual al `:id` de la URL y una opción sintética
   `Proyecto <id>` cuando el actual no está en la lista. Al elegir otro proyecto se navega con
   `rutaDeSeccion`; elegir el actual no navega. El enlace "Nuevo proyecto" se renderiza siempre,
   incluso en carga o error, con semántica de enlace (clase `.button` como en el panel). No se
   usa el componente `Button` porque esta historia solo navega con enlaces y un `<select>`; se
   reutilizan `LoadingIndicator`, `ErrorMessage` y el cliente, y no se agrega ninguna acción de
   botón en esta historia.
5. **Estados visibles y no bloqueantes.** Selector: `LoadingIndicator` al cargar y `ErrorMessage`
   ante error; en ambos casos el `<select>` se muestra deshabilitado con la única opción
   `Proyecto <id>` del proyecto actual, y al resolverse la carga se habilita con la lista
   completa. El enlace "Nuevo proyecto" y los enlaces de sección siguen visibles y operativos en
   todos los estados. Lista: los mismos estados más la invitación a crear cuando el listado está
   vacío. Los enlaces de sección del encabezado viven fuera del selector y no dependen del
   listado.
6. **Reinicio por `:id`.** Además de conservar las `key` existentes, `LayoutProyecto` monta
   `<Outlet key={id} />` para que cualquier pantalla de proyecto se remonte al cambiar de
   proyecto y no mezcle estado (FR-013). El selector queda fuera del `Outlet`, así que no se
   desmonta ni vuelve a pedir la lista al navegar.
7. **Errores del cliente tal cual.** El hook usa el mensaje de `ApiError` o `ConnectionError` del
   cliente compartido; si la respuesta exitosa no es un arreglo, la trata como error con un
   mensaje genérico. No se inventan códigos ni se reintenta automáticamente.
8. **Sin estado compartido entre selector y lista.** Cada pantalla monta su consumidor; no hay
   caché global ni sincronización. Es suficiente para la historia y evita complejidad (YAGNI).
   Una lista desactualizada después de crear se resuelve al volver a montar la pantalla, como
   documenta la spec.

## Implementation Sequence (test-first, one commit per phase)

1. **Infraestructura de prueba (verde).** Modificar `web/src/rutas.test.tsx` para simular
   `GET /projects` con `[]` en `beforeEach` y limpiar con `afterEach`. El mock debe devolver una
   `Response` nueva por llamada y despachar por URL: `[]` para `/projects`, una respuesta mínima
   coherente para los `GET /projects/:id` y `/projects/:id/members` que ya usan las pantallas
   existentes. No hay cambio de comportamiento todavía; la suite existente queda verde antes y
   después de que el encabezado empiece a pedir la lista. Commit de prueba, no de RED.
2. **Lógica pura.** `RED: rutaDeSeccion y etiquetaProyecto (falla)` con `ruta.test.ts`; después
   `GREEN: implementar rutaDeSeccion y etiquetaProyecto, tests en verde` en `ruta.ts`. REFACTOR
   solo si hay algo que limpiar.
3. **Hook y API.** `RED: estados de useProyectos (falla)` con `useProyectos.test.tsx`; después
   `GREEN: implementar api.ts y useProyectos.ts, tests en verde`. Cubre carga, éxito, error,
   respuesta no interpretable y desmontaje antes de resolver. REFACTOR si corresponde.
4. **Selector.** `RED: selector de proyecto (falla)` con `SelectorProyecto.test.tsx`; después
   `GREEN: implementar SelectorProyecto.tsx, tests en verde`. Cubre etiqueta actual, opciones,
   fallback, navegación a la misma sección, mismo proyecto y enlace permanente. Los estados de
   carga y error (US4) se agregan sobre el mismo componente: `<select>` deshabilitado con la
   única opción `Proyecto <id>` del proyecto actual junto a `LoadingIndicator` o `ErrorMessage`,
   y se habilita con la lista completa al resolverse la carga.
5. **Lista.** `RED: lista de proyectos (falla)` con `ListaProyectos.test.tsx`; después
   `GREEN: implementar ListaProyectos.tsx, tests en verde`. Cubre carga, error, vacío, orden y
   enlaces a `/projects/<id>/proyecto`.
6. **Integración.** `RED: navegación integrada en rutas (falla)` agregando a `rutas.test.tsx` los
   casos de lista en el inicio, cambio de proyecto conservando sección, enlace "Nuevo proyecto",
   falla del listado con navegación de secciones operativa y reinicio de estado por `:id`.
   Después `GREEN: montar selector, componer el inicio y reiniciar por :id` en `LayoutProyecto.tsx`
   y `rutas.tsx`. REFACTOR si corresponde.
7. **Limpieza.** Eliminar `web/src/compartido/router/Inicio.tsx`, verificar que no queden
   referencias, correr `npm test`, `npm run lint` y `npm run build`. Commit
   `chore: eliminar Inicio.tsx sin uso`. El PR cierra #148 y #141 en su descripción.

> El primer commit es infraestructura de prueba fuera del ciclo RED/GREEN porque solo prepara el
> entorno de prueba para que la suite existente siga verde. Los commits de los pasos 2 a 6 siguen el
> ciclo y cada fase tiene su propio commit. No se crea un commit de REFACTOR vacío.

## Testing Strategy

**Lógica pura (`ruta.test.ts`)**

- `/projects/1/backlog` + `2` → `/projects/2/backlog`.
- `/projects/1/proyecto/edit` + `2` → `/projects/2/proyecto` (solo el primer tramo).
- `/projects/1` + `2` → `/projects/2/proyecto` (sin tramo).
- `/projects/1/otra` + `2` → `/projects/2/otra` (sección desconocida).
- `/projects/1/backlog/` + `2` → `/projects/2/backlog` (barra final).
- `etiquetaProyecto({ id: 7, nombre: 'Alfa' })` → `Alfa (#7)`.
- Dos proyectos con el mismo nombre producen etiquetas distintas por el `#id`.

**Hook (`useProyectos.test.tsx`)**

- Estado inicial `cargando`.
- `GET /projects` con arreglo → `listo` con los proyectos en el orden recibido.
- Error del cliente → `error` con el mensaje del cliente.
- Respuesta exitosa que no es arreglo → `error`, sin romper el render.
- Desmontar antes de que resuelva la promesa no actualiza estado ni genera errores de consola.

**Selector (`SelectorProyecto.test.tsx`)**

- `combobox` accesible llamado "Proyecto"; opciones `<nombre> (#<id>)`.
- El proyecto del `:id` aparece seleccionado.
- Si el `:id` no está en la lista, se agrega `Proyecto <id>` y queda seleccionado.
- Elegir otro proyecto navega a la misma sección; elegir el actual no navega.
- `LoadingIndicator` durante la carga y `ErrorMessage` ante error; en ambos estados el `<select>`
  está deshabilitado con la única opción `Proyecto <id>` del proyecto actual.
- Al resolverse la carga, el `<select>` se habilita con la lista completa.
- El enlace "Nuevo proyecto" está visible en todos los estados.

**Lista (`ListaProyectos.test.tsx`)**

- Enlaces `/projects/<id>/proyecto` con etiqueta `<nombre> (#<id>)`, en el orden de la API.
- Estados de carga, error y vacío; en vacío se invita a crear el primero.
- No se muestra error cuando la lista está vacía.

**Integración (`rutas.test.tsx`)**

- El inicio muestra el formulario de creación y la lista de proyectos; un enlace de la lista abre
  el panel.
- Cambiar de proyecto con el selector conserva la sección (`/backlog` → `/backlog`).
- "Nuevo proyecto" lleva a `/` desde una pantalla de proyecto.
- Con el listado en error, el selector muestra el error y los enlaces de sección siguen
  navegando.
- Cambiar el `:id` remonta la pantalla y no conserva estado del proyecto anterior (se conserva y
  amplía la prueba existente de reinicio).
- La suite existente se mantiene verde con el `GET /projects` simulado.

**Comandos de verificación**

```bash
cd web
npm test
npm run lint
npm run build
```

No se requieren cambios ni pruebas nuevas de backend; el CI de backend sigue corriendo sin
cambios.

## Risks and Mitigations

| Riesgo | Mitigación |
|---|---|
| La composición del inicio podría requerir tocar `crear.tsx` | Se usa un layout pathless en `rutas.tsx` con `<Outlet />`; el formulario queda intacto y sus pruebas no cambian. |
| El selector y la lista podrían pedir la lista dos veces en la misma pantalla | Nunca conviven: el selector vive en `LayoutProyecto` y la lista en `/`. Cada uno hace un solo pedido al montarse. |
| `Outlet` sin `key` conservaría el estado de una pantalla al cambiar de proyecto | Se monta `<Outlet key={id} />`; se conservan además las `key` existentes de las pantallas. |
| El proyecto actual no aparece en el listado y el `<select>` quedaría sin opción válida | Se agrega la opción sintética `Proyecto <id>` y se selecciona por `value`. |
| Un `200` sin cuerpo o con un JSON que no es arreglo rompería el render | El hook valida `Array.isArray` y pasa a estado de error. |
| La respuesta llega después de desmontar el componente | Bandera `vigente` en el `useEffect`; no se actualiza estado tras el desmontaje. |
| El nuevo `fetch` en el encabezado rompe pruebas de rutas existentes | Primer commit con el mock de `GET /projects`; se ejecuta la suite antes de cambiar el encabezado. |
| Estilos nuevos fuera de los archivos permitidos | No se agrega CSS; se reutilizan clases y marcado existentes. |
| `web/README.md` menciona el `Inicio.tsx` eliminado | Queda fuera de los archivos permitidos; se registra como deuda documental para el equipo, sin bloquear esta historia. |

## Complexity Tracking

No hay violaciones a la constitución ni complejidad adicional que justificar. Se descartan
contexto global, caché, retry, paginación y librerías de datos por YAGNI.

## Phase Outputs

- Fase 0: [`research.md`](./research.md)
- Fase 1: [`data-model.md`](./data-model.md), [`contracts/ui.md`](./contracts/ui.md),
  [`quickstart.md`](./quickstart.md)
- Fase 2 (posterior, `/speckit.tasks`): `tasks.md` con las fases RED/GREEN/REFACTOR y sus commits.
