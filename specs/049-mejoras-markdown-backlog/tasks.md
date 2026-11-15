---

description: "Tareas de #215: mejoras de presentación del Markdown en el backlog (frontend, TDD)"
---

# Tasks: Mejoras de Markdown en el backlog (#215)

**Input**: spec.md, plan.md, research.md, data-model.md, contracts/ui.md, quickstart.md

**Tests**: OBLIGATORIOS (constitución, Principio I). Cada comportamiento va como
RED → COMMIT → GREEN → COMMIT. No se hace `git push` ni merge desde las tareas.

## Format: `[ID] [P?] Description`

- **[P]**: se puede hacer en paralelo (archivos distintos).

---

## Phase 1: Editor reutilizable

- [x] T001 RED: crear `web/src/features/backlog/CampoMarkdown.test.tsx` con escenarios
  Given-When-Then: muestra label, textarea y ayuda; Negrita/Cursiva/Código envuelven la
  selección; Título y Lista prefijan líneas; Enlace e Imagen insertan la sintaxis; la pestaña
  Vista previa renderiza el Markdown y conserva el texto; texto vacío informa que no hay
  contenido; error junto al campo; campo deshabilitado deshabilita la barra. Confirmar que
  FALLA (compila con un stub sin lógica).
- [x] T002 COMMIT `RED: agregar pruebas del editor Markdown del backlog (fallan) (Refs #215)`.
- [x] T003 GREEN: crear `CampoMarkdown.tsx` con la barra, la selección del `textarea`, la
  ayuda, las pestañas y la vista previa con `Markdown`; `npm test -- CampoMarkdown` en verde.
- [x] T004 COMMIT `GREEN: implementar el editor Markdown del backlog, test en verde (Refs #215)`.

---

## Phase 2: Integración en los formularios

- [x] T005 RED: en `crear.test.tsx`, comprobar que Descripción y Criterios muestran la ayuda
  de Markdown y la barra de formato. En `Listado.test.tsx`, comprobar lo mismo en el
  formulario de edición. Confirmar que FALLA.
- [x] T006 COMMIT `RED: agregar pruebas de ayuda y barra en crear/editar (fallan) (Refs #215)`.
- [x] T007 GREEN: reemplazar los `textarea` de Descripción y Criterios por `CampoMarkdown` en
  `crear.tsx` y `Edicion.tsx`; conservar labels, errores, `disabled` y el envío actual.
  `npm test -- crear Listado` en verde.
- [x] T008 COMMIT `GREEN: integrar el editor Markdown en crear y editar, tests en verde (Refs #215)`.

---

## Phase 3: Presentación

- [x] T009 Agregar a `backlog.css` los tamaños de `h1`–`h3` dentro de `.backlog .markdown`,
  el `max-height` de las imágenes y los estilos de barra, pestañas, ayuda y vista previa.
- [x] T010 COMMIT `style: acotar títulos e imágenes y estilizar el editor Markdown (Refs #215)`.

---

## Phase 4: Documentación y cierre

- [x] T011 Agregar la spec 049, actualizar el `Status` y marcar estas tareas.
- [x] T012 COMMIT `docs: agregar la spec 049 de mejoras de Markdown (Refs #215)`.
- [x] T013 Agregar la fila de #215 a las mejoras técnicas de `docs/PRODUCT-BACKLOG.md` con
  el estado y el PR.
- [x] T014 COMMIT `docs: registrar #215 en el Product Backlog (Refs #215)`.
- [x] T015 `cd web && npm test`, `npm run lint` y `npm run build` en verde.
- [x] T016 PR #326 creado con «Closes #215».
