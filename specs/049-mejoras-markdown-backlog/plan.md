# Implementation Plan: Mejoras de Markdown en el backlog (#215)

**Branch**: `feature/215-mejoras-markdown-backlog` | **Date**: 2026-11-15 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/049-mejoras-markdown-backlog/spec.md`

## Summary

Mejora solo de frontend sobre el Markdown ya existente (#211/#174). Se acotan los títulos
`h1`–`h3` y las imágenes dentro de `.backlog .markdown`, se agrega una ayuda de sintaxis
debajo de Descripción y Criterios en creación y edición, y se incorpora un editor liviano
con barra de formato y pestaña «Vista previa». La barra opera sobre la selección del
`textarea`; la vista previa reutiliza el componente `Markdown`. Sin dependencias nuevas, sin
cambios de API ni de esquema.

## Technical Context

**Language/Version**: TypeScript 6 + React 19, Vite 8, Vitest + Testing Library (jsdom).

**Primary Dependencies**: las existentes (`react-markdown` ya está); no se agregan
dependencias.

**Storage**: no aplica (solo frontend; no se guardan imágenes).

**Testing**: pruebas de componente Vitest con el escenario Given-When-Then en el nombre;
se mantienen en verde las pruebas de `Markdown`, `crear` y `Listado`.

**Target Platform**: navegador.

**Project Type**: SPA React del monorepo.

**Constraints**: no tocar la API, el esquema ni el contrato de las historias; respetar el
límite de 2000 caracteres y el saneamiento de `react-markdown`; no subir archivos.

**Scale/Scope**: un componente reutilizable (`CampoMarkdown`), su integración en dos
formularios, estilos en `backlog.css` y pruebas de componente.

## Constitution Check

| Principio | Estado |
|---|---|
| I. Test-First | RED/GREEN por componente y por integración; commits por fase en [tasks.md](tasks.md) |
| II. SDD | Spec, plan, research, contrato, quickstart y tasks antes de implementar |
| III. BDD | Escenarios de componente con Given-When-Then en el nombre; no hay cambios de API |
| IV. Clean/Hexagonal | Solo frontend; la lógica de Markdown sigue en el componente existente |
| V. SOLID/KISS/YAGNI | Un componente reutilizable, sin editor de terceros ni subida de archivos |
| VI. Stack | Sin dependencias nuevas |

Sin violaciones.

## Project Structure

### Documentation (this feature)

```text
specs/049-mejoras-markdown-backlog/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/ui.md
└── tasks.md
```

### Source Code (repository root)

```text
web/src/features/backlog/
├── CampoMarkdown.tsx        # NUEVO: label + barra + textarea + ayuda + vista previa
├── CampoMarkdown.test.tsx   # NUEVO: pruebas de componente
├── Markdown.tsx             # sin cambios; se reutiliza en la vista previa
├── crear.tsx                # usa CampoMarkdown en descripción y criterios
├── Edicion.tsx              # usa CampoMarkdown en descripción y criterios
├── crear.test.tsx           # + ayuda y barra en el formulario de creación
├── Listado.test.tsx         # + ayuda y barra en el formulario de edición
└── backlog.css              # títulos acotados, max-height de imágenes, estilos del editor
```

**Structure Decision**: el editor vive en un componente del feature `backlog` para
reutilizarlo en creación y edición sin duplicar la lógica de selección. La vista previa
reutiliza `Markdown`; el renderizado de las tarjetas no cambia.

## Decisiones técnicas

| # | Decisión | Motivo |
|---|---|---|
| 1 | Acotar `h1`–`h3` con `.backlog .markdown` | El CSS global de `h1` rompe la jerarquía; el selector queda acotado al backlog. |
| 2 | `max-height` en imágenes | Evita que una imagen externa ocupe media pantalla, sin bloquear su carga por URL. |
| 3 | Componente `CampoMarkdown` | Reutiliza la misma ayuda, barra y vista previa en creación y edición. |
| 4 | Barra sobre `textarea` con `selectionStart/End` | No agrega dependencias y funciona con el teclado y el ratón. |
| 5 | Vista previa con `Markdown` | Un solo renderizador, el mismo saneamiento que las tarjetas. |
| 6 | No subir imágenes | Render no tiene disco persistente; solo se admite URL escrita. |
| 7 | `aria-describedby` a la ayuda o al error (uno solo) | Mantiene los helpers de pruebas existentes y no duplica el anuncio. |

## Complexity Tracking

Sin violaciones de la constitución que justificar.
