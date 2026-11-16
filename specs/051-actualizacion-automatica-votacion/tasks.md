---
description: "Tareas de #225: actualización automática de la votación de Planning Poker (frontend, TDD)"
---

# Tasks: Actualización automática de la votación (#225)

**Input**: spec.md, plan.md, research.md, data-model.md, contracts/polling.md, quickstart.md

**Tests**: OBLIGATORIOS (constitución, Principio I). Cada comportamiento va como
RED → COMMIT → GREEN → COMMIT. No se hace `git push` ni merge desde las tareas.

## Format: `[ID] [P?] Description`

- **[P]**: se puede hacer en paralelo (archivos distintos).

---

## Phase 1: Documentación SDD

- [x] T001 Crear la spec 051 con objetivo, entradas/salidas, escenarios, Edge Cases en EARS,
  requisitos, criterios de éxito y supuestos.
- [x] T002 Crear `plan.md`, `research.md`, `data-model.md`, `contracts/polling.md` y
  `quickstart.md`.
- [x] T003 COMMIT `docs: agregar la spec 051 de actualización automática de la votación (Refs #225)`.

---

## Phase 2: Polling silencioso

- [x] T004 RED: crear `useVotacion.test.ts` con escenarios Given-When-Then de consulta cada 5
  segundos, pausa con la pestaña oculta, pausa con la ronda revelada, pausa con una acción en
  curso, reanudación al volver a visible y consulta sin indicador de carga. Confirmar que FALLA.
- [x] T005 COMMIT `RED: agregar pruebas del polling de la votación (fallan) (Refs #225)`.
- [x] T006 GREEN: agregar a `useVotacion.ts` el efecto de polling silencioso, las pausas por
  visibilidad/revelación/acción, el descarte de respuestas viejas y el procesamiento común de
  la votación. `npm test -- useVotacion` en verde.
- [x] T007 COMMIT `GREEN: implementar el polling de la votación, tests en verde (Refs #225)`.

---

## Phase 3: Documentación y cierre

- [x] T008 Agregar la fila de #225 a las mejoras técnicas de `docs/PRODUCT-BACKLOG.md`.
- [x] T009 COMMIT `docs: registrar #225 en el Product Backlog (Refs #225)`.
- [x] T010 `cd web && npm test`, `npm run lint` y `npm run build` en verde.
- [x] T011 PR creado con «Closes #225».
