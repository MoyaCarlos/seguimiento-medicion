# Implementation Plan: Markdown seguro en el Product Backlog

**Branch**: `feature/147-markdown-historias` | **Date**: 2026-10-31 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/027-markdown-historias/spec.md`

## Summary

Renderizar como Markdown la descripción y los criterios de aceptación en el listado del
Product Backlog, sin habilitar HTML crudo. Se agrega `react-markdown` como única
dependencia nueva y un componente `Markdown` reutilizable; los formularios de edición
siguen enviando el texto original sin cambios.

## Technical Context

**Language/Version**: TypeScript 6 + React 19; Vite 8

**Primary Dependencies**: `react-markdown` (nueva). Sin `rehype-raw`, sin
`dangerouslySetInnerHTML`; la transformación de URLs por defecto de `react-markdown`
cubre los protocolos seguros.

**Storage**: N/A (solo presentación)

**Testing**: Vitest + Testing Library (jsdom). Test unitario del componente `Markdown`
(CommonMark y XSS) y caso de integración en `Listado.test.tsx`.

**Target Platform**: navegador (SPA React servida por Vite/Docker)

**Project Type**: frontend

**Constraints**: no cambia la API, el esquema, las migraciones ni los límites de
longitud; no se habilita HTML crudo; no se agregan otras dependencias.

**Scale/Scope**: un componente nuevo, dos usos en `Listado.tsx`, estilos Markdown y
tests; una dependencia nueva.

## Constitution Check

| Principio | Estado |
|---|---|
| I. Test-First | RED/GREEN por comportamiento: primero los tests de Markdown/XSS, después el componente y la integración |
| II. SDD | Spec → plan → tasks → analyze antes de implementar |
| III. BDD | No hay backend nuevo; los escenarios se automatizan como pruebas de componente con Vitest, con el escenario en el nombre |
| IV. Clean/Hexagonal | Solo capa de presentación; no toca dominio, servicios ni repositorios |
| V. KISS/YAGNI | Un componente reusable y una dependencia; sin preview opcional ni GFM |
| VI. Stack | React + Vite + Vitest; `react-markdown` es compatible con React 19 |

Sin violaciones. Re-chequeo post-diseño: igual.

## Project Structure

### Documentation (this feature)

```text
specs/027-markdown-historias/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/componente-markdown.md
├── checklists/requirements.md
└── tasks.md
```

### Source Code (repository root)

```text
web/
├── package.json                     # + react-markdown
├── package-lock.json                # actualizado por npm
└── src/features/backlog/
    ├── Markdown.tsx                 # nuevo: render seguro
    ├── Markdown.test.tsx            # nuevo: CommonMark + XSS
    ├── Listado.tsx                  # usa Markdown en descripción y criterios
    ├── Listado.test.tsx             # caso de integración Markdown
    └── backlog.css                  # estilos del contenido Markdown
```

No se modifica `Edicion.tsx`, `crear.tsx`, la API ni el esquema.

## Decisiones de diseño

- **Componente único `Markdown`**: encapsula `react-markdown` y la clase de estilos;
  así el listado no repite configuración y el día que haya que cambiar la política de
  seguridad se cambia en un solo lugar.
- **Sin HTML crudo**: no se agrega `rehype-raw`. El HTML escrito en la descripción no
  se interpreta; los tests verifican que no aparezcan `<script>` ni imágenes con
  `onerror`.
- **Sin GFM**: las listas, negritas, cursivas, código y enlaces de CommonMark cubren el
  caso de uso. Evita `remark-gfm` y mantiene chico el árbol de dependencias.
- **CSS acotado**: los estilos Markdown viven bajo `.backlog` y no afectan otras
  pantallas.
- **Formularios intactos**: `Edicion` y `crear` siguen mostrando el texto en
  `textarea`; la vista previa queda fuera de alcance (opcional del issue).

## Testing Strategy

- `Markdown.test.tsx`:
  - renderiza negritas, listas, código y enlaces comunes;
  - no inserta `<script>` ni `<img onerror>`;
  - no deja un enlace `javascript:` ejecutable.
- `Listado.test.tsx`: una historia con Markdown se ve formateada en el listado.
- Los tests existentes de listado y edición deben seguir en verde, incluido el caso de
  criterios `null`.

## Risks

- La interpretación de Markdown puede cambiar el HTML esperado por tests existentes;
  se mitiga manteniendo el texto visible y corriendo la suite completa.
- Una URL de imagen externa podría filtrar datos al navegador; no se agrega nada que
  cargue recursos no presentes en el texto y se mantiene la política de URLs seguras
  de la librería.
