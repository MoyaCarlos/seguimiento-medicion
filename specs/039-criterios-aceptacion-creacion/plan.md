# Implementation Plan: Criterios de aceptación al crear una historia

**Branch**: `feature/174-criterios-aceptacion` | **Date**: 2026-11-08 | **Spec**: [spec.md](./spec.md)

**Issue**: #174 (parte de HU-01 #1; relacionado con HU-15 #43/#98)

**Input**: Feature specification from `specs/039-criterios-aceptacion-creacion/spec.md`

## Summary

Permitir cargar los criterios de aceptación en la creación de una historia y mostrarlos en el resultado de la creación. El backend ya persiste `historias.criterios_aceptacion` y la edición de HU-15 ya los valida; falta aceptarlos en `POST /backlog`, pasarlos por el dominio y devolverlos (la respuesta ya incluye el campo). El frontend suma el campo opcional al formulario de creación y los muestra con Markdown en la confirmación; el listado ya los muestra. No hay cambio de esquema ni migración. Notas y depende_de (D17) quedan fuera.

## Technical Context

**Language/Version**: Go 1.26.5 (backend), TypeScript 6 + React 19 + Vite 8 (frontend).

**Primary Dependencies**: biblioteca estándar de Go, `modernc.org/sqlite` (ya existente), React, react-router 8, Vitest 4 + Testing Library. No se agregan dependencias.

**Storage**: SQLite. Se reutiliza `historias.criterios_aceptacion` (`TEXT?`, migración 001, `MODELO-DE-DATOS.md` §4.4). No hay tablas, columnas ni migraciones nuevas.

**Testing**: Go `testing` para dominio, service, handler y repositorio; Godog para los escenarios BDD de la API/servicio; Vitest + Testing Library para la pantalla de creación y la confirmación (escenario en el nombre).

**Target Platform**: aplicación web (API Go + SPA React servida por el backend).

**Project Type**: web application; el cambio toca `internal/` y `web/`.

**Performance Goals**: sin cambio; una sola creación por envío y una sola consulta de listado como hoy.

**Constraints**: criterios opcionales, máximo 2000 caracteres, recorte de espacios externos y vacío = `NULL`; sin cambio de esquema; sin permisos nuevos; sin notas/depende_de.

**Scale/Scope**: un campo por historia; sin paginación ni filtros nuevos.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principio | Cumplimiento |
|---|---|
| I. Test-First | Cumple. Dominio, service, handler y frontend van con RED → GREEN en commits propios. El repositorio ya persiste el campo; se agrega una prueba de caracterización, sin RED artificial. |
| II. SDD | Cumple. Spec 039 → clarify → plan → tasks → analyze antes de implementar. |
| III. BDD | Cumple. Escenarios Godog para la creación con/sin criterios y el límite; pruebas de componente para el formulario y la confirmación. |
| IV. Clean/Hexagonal | Cumple. La validación vive en `internal/domain`; el handler solo traduce HTTP y el repositorio solo persiste. |
| V. Clean Code, KISS, YAGNI, DRY | Cumple. Se reutiliza `normalizarCriterios` para el alta y la edición; no se agregan capas ni campos que el issue no pida. |
| VI. Stack | Cumple. Sin dependencias nuevas ni cambios de esquema. |

Resultado: sin violaciones, antes y después del diseño.

## Project Structure

### Documentation (this feature)

```text
specs/039-criterios-aceptacion-creacion/
├── spec.md
├── plan.md              # este archivo
├── research.md          # decisiones técnicas
├── data-model.md        # campos y reglas; enlaza a MODELO-DE-DATOS.md
├── contracts/api-backlog.md
├── quickstart.md        # validación manual y automática
├── checklists/requirements.md
└── tasks.md             # lo genera /speckit-tasks
```

### Source Code (repository root)

```text
internal/
├── domain/
│   ├── backlog_item.go            # ConCriteriosAceptacion + normalización compartida con la edición
│   └── backlog_item_test.go       # casos nil/vacío/válido/límite/error
├── service/
│   ├── crear_historia.go          # CrearHistoriaInput.CriteriosAceptacion
│   └── crear_historia_test.go
└── http/
    ├── backlog_dto.go             # crearHistoriaRequest.CriteriosAceptacion
    └── backlog_handler_test.go    # 201 con criterios/null y 400 con campo

features/
├── criterios_aceptacion_creacion.feature   # escenarios @HU-01
└── steps_creacion.go                       # pasos nuevos reutilizando scenarioContext

web/src/features/backlog/
├── historia.ts / historia.test.ts          # tipo, estado inicial, cuerpo y validación
├── crear.tsx / crear.test.tsx              # campo, confirmación y errores
└── api.test.ts                             # si el fixture de creación necesita el campo

docs/PRODUCT-BACKLOG.md                     # nota de #174 en HU-01/HU-15
```

**Structure Decision**: se extiende el flujo existente de HU-01/HU-15 sin crear features ni endpoints nuevos. El repositorio SQLite no cambia porque `Guardar` ya incluye `criterios_aceptacion`.

## Risks

- **Validación en el service**: el alta debe encadenar `ConCriteriosAceptacion` después de `NewBacklogItem` y antes de verificar el proyecto y persistir. Si se olvida, el campo se perdería en silencio; lo cubren los tests de service y del handler. No se cambia la firma de la fábrica, así no hay que tocar los llamados existentes.
- **Confirmación y accesibilidad**: los criterios se muestran con `Markdown` dentro del bloque de resultado; hay que mantener una estructura semántica válida y el texto alternativo "Sin criterios de aceptación".
- **PRs abiertos**: #279 (HU-10) agrega métodos en `internal/repository/sqlite_backlog.go`; este cambio no modifica ese archivo (solo, si hace falta, su test). #276 y #269 no tocan el alta de historias. Verificar `merge-tree` antes de mergear.
- **Límite con Markdown**: el límite se cuenta sobre el texto recortado (runas), igual que la edición; los caracteres de Markdown cuentan como contenido.

## Complexity Tracking

Sin violaciones que justificar.
