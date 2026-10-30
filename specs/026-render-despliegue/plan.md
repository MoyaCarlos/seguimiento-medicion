# Implementation Plan: Despliegue continuo en Render

**Branch**: `feature/139-render` | **Date**: 2026-10-30 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/026-render-despliegue/spec.md`

## Summary

El servicio de Render ya existe y funciona. Esta tarea no cambia código: agrega `render.yaml`
(referencia versionada, validada contra el esquema oficial de Render), la sección "Despliegue en
Render" del README y la verificación manual del servicio real, anotada en `quickstart.md`.

## Technical Context

**Language/Version**: N/A (YAML y Markdown); el `Dockerfile` existente no se toca

**Primary Dependencies**: ninguna nueva. La validación usa el esquema público `https://render.com/schema/render.yaml.json` con herramientas temporales fuera del repo

**Storage**: N/A; sin disco persistente (la base se vacía en cada deploy); esquema sin cambios

**Testing**: validación del YAML contra el esquema (JSON Schema) y verificación manual con `curl`; sin tests de Go ni de Vitest nuevos

**Target Platform**: Render, servicio web Docker, plan Free, región Oregon

**Project Type**: infraestructura y documentación

**Constraints**: no tocar el dashboard ni secretos; solo `GET` contra el servicio real; no hacer push ni abrir el PR

**Scale/Scope**: 1 archivo nuevo en la raíz (`render.yaml`), 1 sección de README, 1 spec

## Constitution Check

| Principio | Estado |
|---|---|
| I. Test-First | No aplica el ciclo RED/GREEN: no hay lógica de negocio. La validación contra el esquema hace de chequeo |
| II. SDD | Spec → clarify → plan → tasks → analyze antes de implementar |
| III. BDD | Sin comportamiento de negocio: no aplica `.feature` |
| IV. Clean/Hexagonal | No se toca código ni capas |
| V. KISS/YAGNI | Un solo archivo de configuración y una sección de README; sin Blueprint ni automatización nueva |
| VI. Stack | Sin dependencias nuevas |

Sin violaciones. Re-chequeo post-diseño: igual.

## Project Structure

### Documentation (this feature)

```text
specs/026-render-despliegue/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/render-yaml.md
├── checklists/requirements.md
└── tasks.md
```

### Source Code (repository root)

```text
render.yaml   # nuevo: referencia versionada, no se enlaza a ningún Blueprint
README.md     # sección nueva "Despliegue en Render"
```

No se modifica nada más.

## Decisiones de diseño

- `render.yaml` lleva un comentario al inicio que dice que es solo referencia y por qué no se crea un Blueprint (FR-004, FR-005).
- Las claves y su correspondencia con el dashboard están en [contracts/render-yaml.md](contracts/render-yaml.md).
- Una clave que no se pueda validar contra el esquema se deja fuera y se informa (research.md).
- `render.yaml` no se agrega a `ignoredPaths`: el archivo debe reproducir la configuración actual (ver Clarifications de la spec).
