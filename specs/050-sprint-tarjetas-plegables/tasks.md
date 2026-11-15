---
description: "Tareas de #219: tarjetas plegables y confirmaciÃ³n de estado en el Sprint (frontend, TDD)"
---

# Tasks: Tarjetas plegables y confirmaciÃ³n de estado (#219)

**Input**: spec.md, plan.md, research.md, data-model.md, contracts/ui.md, quickstart.md

**Tests**: OBLIGATORIOS (constituciÃ³n, Principio I). Cada comportamiento va como
RED â†’ COMMIT â†’ GREEN â†’ COMMIT. No se hace `git push` ni merge desde las tareas.

## Format: `[ID] [P?] Description`

- **[P]**: se puede hacer en paralelo (archivos distintos).

---

## Phase 1: DocumentaciÃ³n SDD

- [x] T001 Crear la spec 050 con objetivo, entradas/salidas, escenarios, Edge Cases en EARS,
  requisitos, criterios de Ã©xito y supuestos.
- [x] T002 Crear `plan.md`, `research.md`, `data-model.md`, `contracts/ui.md` y `quickstart.md`.
- [x] T003 COMMIT `docs: agregar la spec 050 de tarjetas plegables y confirmaciÃ³n (Refs #219)`.

---

## Phase 2: ConfirmaciÃ³n del cambio de estado

- [x] T004 RED: en `EstadoHistoria.test.tsx`, actualizar los escenarios para pulsar la acciÃ³n y
  luego Â«ConfirmarÂ»; agregar Â«CancelarÂ», Escape, clic fuera, foco en Â«CancelarÂ» y doble envÃ­o.
  Confirmar que FALLA.
- [x] T005 COMMIT `RED: agregar pruebas de confirmaciÃ³n del estado (fallan) (Refs #219)`.
- [x] T006 GREEN: implementar las fases `reposo`/`confirmando`/`enviando` en
  `EstadoHistoria.tsx` con Confirmar, Cancelar, Escape y la protecciÃ³n de doble envÃ­o.
  `npm test -- EstadoHistoria` en verde.
- [x] T007 COMMIT `GREEN: implementar la confirmaciÃ³n del estado, tests en verde (Refs #219)`.

---

## Phase 3: Tarjetas plegables del Sprint Backlog

- [x] T008 RED: en `SprintBacklog.test.tsx`, comprobar que la tarjeta estÃ¡ plegada con tÃ­tulo,
  estado y Story Points, que el resumen no tiene la prioridad y que al abrir aparecen los
  datos adicionales y las acciones. Ajustar los escenarios existentes para abrir la tarjeta
  antes de pulsar Â«IniciarÂ» o Â«CompletarÂ». Confirmar que FALLA.
- [x] T009 COMMIT `RED: agregar pruebas de tarjetas plegables en el Sprint (fallan) (Refs #219)`.
- [x] T010 GREEN: usar `<details>`/`<summary>` en `SprintBacklog.tsx` y conservar el detalle,
  las acciones y el esfuerzo. `npm test -- SprintBacklog` en verde.
- [x] T011 COMMIT `GREEN: implementar las tarjetas plegables, tests en verde (Refs #219)`.

---

## Phase 4: PresentaciÃ³n y cierre

- [x] T012 Agregar a `sprints.css` los estilos de lista, tarjeta, resumen, marca, detalle y
  confirmaciÃ³n.
- [x] T013 COMMIT `style: estilizar las tarjetas plegables y la confirmaciÃ³n (Refs #219)`.
- [x] T014 Agregar la fila de #219 a las mejoras tÃ©cnicas de `docs/PRODUCT-BACKLOG.md`.
- [x] T015 COMMIT `docs: registrar #219 en el Product Backlog (Refs #219)`.
- [x] T016 `cd web && npm test`, `npm run lint` y `npm run build` en verde.
- [x] T017 PR creado con Â«Closes #219Â».
