# Implementation Plan: Pantalla de Sprints (HU-05 Frontend, Issue #26)

**Branch**: `feature/HU-05-frontend` | **Date**: 2026-10-24 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/011-hu-05-frontend-sprints/spec.md`

## Summary

Reemplazar el placeholder de `web/src/features/sprints/` por la pantalla de Sprints: un botón para crear, un formulario de inicio por Sprint Pendiente, un botón para cerrar por Sprint Activo y tres secciones fijas por estado. Los Sprints de la sesión viven en el estado local de la pantalla (sin listado: #110). Los errores de la API se clasifican en un archivo propio: con campo del formulario van junto al campo; el resto, a un aviso por Sprint (o de pantalla, para crear). Todo sobre la base del frontend (spec 007): cliente compartido, piezas comunes y `Component` exportado, sin tocar `compartido/`, `package.json` ni los registros.

## Technical Context

**Language/Version**: TypeScript ~6.0, React 19 (ya instalados en `web/`)

**Primary Dependencies**: `react-router` (solo `useParams`), cliente y piezas de `web/src/compartido/`. Sin dependencias nuevas

**Storage**: N/A. Estado local de React; se pierde al recargar (spec FR-009)

**Testing**: Vitest + Testing Library + jsdom (ya configurados por la base); `fetch` simulado con `vi.stubGlobal` en los tests de pantalla

**Target Platform**: Navegador moderno; desarrollo con Vite y proxy a `:8080`

**Project Type**: SPA que consume la API REST existente

**Performance Goals**: N/A (un puñado de Sprints en memoria)

**Constraints**: solo `web/src/features/sprints/` y `specs/011-*/`; ninguna importación de otra feature (lo verifica `features/aislamiento.test.ts`); toda llamada pasa por `compartido/api/client.ts`

**Scale/Scope**: una pantalla, 3 llamadas a la API

## Constitution Check

| Principio | Estado | Nota |
|---|---|---|
| I. Test-First | Cumple | RED/GREEN en el mapeo de errores, la agrupación por estado y los estados de botones. Las piezas sin lógica (secciones, CSS) no llevan RED |
| II. SDD | Cumple | specify → clarify → plan → tasks → analyze → implement |
| III. BDD (Godog) | No aplica | Frontend; Godog cubre la API. Los criterios se verifican con [quickstart.md](quickstart.md) |
| IV. Clean/Hexagonal | No aplica | No toca `internal/` |
| V. Clean Code/KISS/YAGNI/DRY | Cumple | Sin estado global ni librería de datos; sin validaciones duplicadas del backend; errores clasificados en un solo lugar |
| VI. Stack | Cumple | Sin dependencias nuevas |

Sin violaciones; no hay Complexity Tracking. Se re-evaluó tras el diseño: sin cambios.

## Project Structure

### Documentation (this feature)

```text
specs/011-hu-05-frontend-sprints/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── pantalla-y-api.md
├── checklists/requirements.md
└── tasks.md          # lo genera /speckit-tasks
```

### Source Code (repository root)

```text
web/src/features/sprints/
├── index.tsx              # exporta Component: estado local, carga inicial, "Crear Sprint", aviso de crear
├── sprints.css            # estilos propios de la feature
├── api.ts                 # crearSprint, iniciarSprint (sobre post del cliente compartido)
├── BotonCerrar.tsx        # botón "Cerrar Sprint" + cerrarSprint (aislado por #117)
├── FormularioInicio.tsx   # Goal + fechas (input type="date"), aviso propio
├── SeccionSprints.tsx     # una sección por estado; tarjeta "Sprint #<id>" con su acción
├── sprint.ts              # tipo Sprint, ESTADOS, agruparPorEstado, guardarSprint (reemplaza por id)
├── sprint.test.ts
├── errores.ts             # clasificarError(error, camposDelFormulario) -> { campo | aviso }
├── errores.test.ts
├── FormularioInicio.test.tsx
├── BotonCerrar.test.tsx
└── index.test.tsx         # flujo de la pantalla con fetch simulado
```

**Structure Decision**: una carpeta, sin subcarpetas. Cada archivo tiene una sola razón de cambio; `BotonCerrar.tsx` agrupa el botón y su llamada para que un cambio de contrato (#117) toque solo ese archivo (FR-007).

## Design decisions

1. **Estado local, un único punto de carga.** `index.tsx` guarda `sprints: Sprint[]` con `useState`, inicializado por `cargarSprintsIniciales()`, que hoy devuelve `[]` (FR-010). Cuando exista #110 se reemplaza el cuerpo de esa función (o se hace un `useEffect`) sin tocar secciones ni acciones (hecho en la spec 015: `useEffect` con `listarSprints`). Cada respuesta de la API reemplaza al Sprint por `id` (`guardarSprint`); la API es la fuente de verdad.
2. **Agrupación y orden.** `agruparPorEstado(sprints)` devuelve `{ Pendiente, Activo, Finalizado }`, cada lista ordenada por `id` ascendente. Las tres secciones se renderizan siempre, con texto "No hay Sprints" cuando están vacías (FR-008).
3. **Clasificación de errores.** `clasificarError(error, campos)` recibe cualquier error lanzado: si es `ApiError` con `field` incluido en `campos` (`sprint_goal`, `fecha_inicio`, `fecha_fin`) devuelve `{ campo, mensaje }`; cualquier otro caso (`ApiError` sin campo o con `estado`, `proyecto_id`, `id`; `ConnectionError`; error desconocido) devuelve `{ aviso: mensaje }`. No valida nada del dominio (FR-011).
4. **Avisos por acción (clarify).** Cada `FormularioInicio` y cada `BotonCerrar` guardan su propio `aviso`/`errores` y los borran al empezar un envío; el de "Crear Sprint" vive en `index.tsx`. Un error de un Sprint no se ve en otro.
5. **Carga y doble envío.** Cada componente de acción tiene su `enviando`: muestra `LoadingIndicator` y deshabilita su botón (FR-014). Formulario y botón de cierre nunca coexisten en el mismo Sprint, así que no hay dos operaciones simultáneas sobre uno.
6. **Fechas.** `<Field type="date">` nativo; el valor ya es `AAAA-MM-DD`. Sin librería ni selector propio.
7. **Proyecto.** `useParams().id` convertido con `Number`; no se valida, el backend responde 400/404 y se muestra su mensaje (spec, caso límite).
8. **Texto de limitación y encabezado de Sprint.** ~~Texto de limitación~~ (retirado en la spec 015); cada tarjeta muestra "Sprint #<id>" (FR-019).
9. **Aislamiento del cierre.** `BotonCerrar` recibe `sprint` y `alActualizar(sprint)`; no sabe de secciones ni de otros componentes. `cerrarSprint` vive en `BotonCerrar.tsx` por FR-007; `api.ts` solo tiene `crearSprint` e `iniciarSprint`. Si #117 cambia la llamada o la respuesta, solo cambia `BotonCerrar.tsx`.
10. **Contrato de la API.** No se redefine: se enlaza al de la spec 003. Ver [contracts/pantalla-y-api.md](contracts/pantalla-y-api.md).
11. **Contrato con la base (FR-020).** `Component` renderiza un `<h1>` con el texto exacto "Sprints" y una línea "Proyecto actual: <strong>{id}</strong>", con el id de la URL en un único `<strong>`; ningún otro elemento muestra solo ese número (los "Sprint #<id>" de las tarjetas llevan el prefijo en el mismo texto). `web/src/rutas.test.tsx` lo verifica y no se modifica.

## Notas para `/speckit-tasks`

- RED/GREEN: `sprint.test.ts` (agrupar, ordenar, reemplazar por id), `errores.test.ts` (con campo del formulario, con campo ajeno, sin campo, conexión, desconocido), `index.test.tsx` (crear, iniciar, cerrar, botones deshabilitados durante el envío, avisos por Sprint, texto de limitación). El componente sin lógica (CSS, `SeccionSprints` presentacional) no lleva RED propio.
- Fases: setup (tipos y `api.ts`) → mapeo de errores y estado (RED/GREEN) → componentes de acción → secciones y pantalla → verificación manual con [quickstart.md](quickstart.md). Antes de cerrar cada fase: `npm test`, `npm run lint`, `npm run build` en `web/`.
- Verificación: correr `web/src/rutas.test.tsx` completo (`npm test` lo incluye) además de los tests de la feature, para confirmar FR-020.
- Commits: `git add` por nombre de archivo; mensajes con `Co-Authored-By`.
