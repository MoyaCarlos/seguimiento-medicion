# Resumen de sesión — HU-01 (Creación de Historias de Usuario)

Fecha: 2026-09-30
Historia: **HU-01 — Creación de Historias de Usuario** (Product Backlog)
Método: SDD con Spec Kit (`/speckit.specify` → `/speckit.clarify` → `/speckit.plan` → `/speckit.tasks` → `/speckit.analyze`)

Este resumen documenta lo realizado **antes** de pasar a la implementación (`/speckit.implement`).

## Qué se hizo

1. **Especificación** (`/speckit.specify`): se generó la spec de HU-01 en
   `specs/001-hu-01-crear-historias-usuario/spec.md`, con 2 historias de usuario
   (US1 creación P1, US2 valor de negocio P2), requisitos funcionales (FR-001..FR-014),
   criterios de éxito medibles, casos borde y supuestos.

Modo de uso:
```
/speckit.specify

Historia HU-N: Nombre de la historia

descripcion de la historia

Criterios de aceptación (BDD):


Prioridad: . Story Points:.

```

EJmeplo: (Se copian del archivo PRODUCT-BACKLOG.md, esta en la carpeta docs)

```
/speckit.specify

Historia HU-01: Creación de Historias de Usuario.

Como Product Builder quiero crear historias de usuario con prioridad, estado
y estimación para poder alimentar y organizar el Product Backlog del proyecto.

Criterios de aceptación (BDD):
Escenario 1: Creación exitosa. Dado que me encuentro logueado en el panel
del Product Backlog, Cuando ingreso los datos obligatorios (título,
descripción, prioridad) y presiono "Guardar", Entonces la historia debe
aparecer al final de la lista con estado "Nuevo" y guardarse en la base de
datos SQLite.
Escenario 2: Faltan campos. Dado que intento crear una historia, Cuando
dejo el campo "Título" en blanco y presiono "Guardar", Entonces el sistema
debe mostrar una advertencia y no debe registrar la HU.

Prioridad: Must have. Story Points: 5.

```

2. **Clarificación** (`/speckit.clarify`): 5 preguntas resueltas con el equipo
   (estimación, valor de negocio, roles, longitudes, escala Fibonacci).
3. **Plan** (`/speckit.plan`): diseño técnico hexagonal en
   `plan.md`, `research.md`, `data-model.md`, `contracts/openapi.yaml` y `quickstart.md`.
4. **Tareas** (`/speckit.tasks`): `tasks.md` con 30 tareas organizadas por historia,
   con TDD (RED→GREEN→REFACTOR) y BDD exigidos por `AGENTS.md`.
5. **Análisis** (`/speckit.analyze`): dos pasadas de consistencia spec/plan/tasks; se
   resolvieron todas las inconsistencias detectadas (0 CRITICAL al cierre).
6. **Escenario BDD**: `features/crear_historia_backlog.feature` con 3 escenarios
   (creación exitosa, falta de título, prioridad inválida).

## Decisiones acordadas

- **Estado inicial canónico**: `"Nueva"` (normalizado en spec, plan, tasks, data-model,
  contrato y `docs/PRODUCT-BACKLOG.md`).
- **Estimación**: no se ingresa al crear; queda vacía y se define vía Planning Poker (HU-02).
- **Valor de Negocio**: opcional, escala Fibonacci `{1, 2, 3, 5, 8, 13, 21}`.
- **Longitudes**: título ≤ 200, descripción ≤ 2000 caracteres.
- **Roles**: Product Builder y Scrum Master; la autorización por rol (FR-014) se difiere a HU-04.
- **Proyecto**: `proyecto_id` obligatorio y `> 0`; la existencia del proyecto se difiere a HU-04.
- **Listado/UI**: fuera de alcance de HU-01; el orden de creación se preserva con `id`
  incremental (la vista de backlog se cubre en un incremento de frontend).
- **Edición y borrado**: fuera de alcance (YAGNI).

## Diseño técnico (resumen)

Arquitectura hexagonal liviana:

- `internal/domain`: entidad `BacklogItem` + Factory `NewBacklogItem` (invariantes).
- `internal/service`: caso de uso `CrearHistoriaBacklog`.
- `internal/repository`: puerto `BacklogRepository` + adaptador SQLite
  (`modernc.org/sqlite`, sin CGO) y migración de `backlog_items`.
- `internal/http`: handler `POST /backlog` (201 con la historia creada, 400 con validación).

## Artefactos generados

- `specs/001-hu-01-crear-historias-usuario/`:
  `spec.md`, `plan.md`, `tasks.md`, `research.md`, `data-model.md`,
  `contracts/openapi.yaml`, `quickstart.md`, `checklists/requirements.md`.
- `features/crear_historia_backlog.feature`.
- `docs/PRODUCT-BACKLOG.md` (estado normalizado a `"Nueva"`).

## Pendiente

- Implementar con `/speckit.implement` (una historia a la vez, commits RED/GREEN/REFACTOR).
- Opcional (gobernanza): ratificar la constitución con `/speckit.constitution`.
