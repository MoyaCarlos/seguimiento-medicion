# Data Model: Adopción del modelo de datos (Issue #107)

Este archivo **no redefine** tablas, columnas ni restricciones (AGENTS.md, MODELO-DE-DATOS.md
§11). La fuente única es [`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md); acá solo se dice
qué secciones usa esta adopción.

| Qué | Dónde |
|---|---|
| Convenciones (idioma, ids, fechas, `CHECK`, `typeof`, sin cascadas) | [§2](../../MODELO-DE-DATOS.md#2-convenciones-valen-para-todas-las-tablas) |
| Tablas de la migración 001: `proyectos`, `usuarios`, `integrantes`, `historias`, `sprints` | [§4.1–§4.5](../../MODELO-DE-DATOS.md#4-tablas) |
| Auditoría de `CHECK` frente a `NULL` | [§4.10](../../MODELO-DE-DATOS.md#410-auditoría-de-los-check-frente-a-null) |
| Dónde vive cada regla y traducción de violaciones a errores | [§5](../../MODELO-DE-DATOS.md#5-dónde-vive-cada-regla) |
| Decisiones cerradas que aplican (D3, D9, D11, D12, D13) | [§7](../../MODELO-DE-DATOS.md#7-decisiones-pendientes-e-inconsistencias-del-backlog) |
| Migraciones: archivos, tabla de control, números reservados, base vieja, test de esquema | [§9](../../MODELO-DE-DATOS.md#9-migraciones) |
| Plan de adopción, fases y excepción de proceso | [§10](../../MODELO-DE-DATOS.md#10-cambios-sobre-lo-que-ya-está-en-main) |

## Lo que la adopción agrega al documento (FR-016a)

Estas precisiones se decidieron en el clarify y el plan (ver *Clarifications* de la
[spec](spec.md)) y se **escriben en MODELO-DE-DATOS.md** en el mismo PR (fase 0b); no se
definen acá:

- §9 punto 1: la tabla de control la crea el runner, y su DDL.
- §4.5: nombres de los índices parciales (`sprints_un_activo_por_proyecto`, `sprints_un_pendiente_por_proyecto`).
- §9 punto 5 y §10: las cuatro tablas que delatan una base vieja.
- §10: la frase sobre la ventana en rojo, que arranca en el `GREEN:` de la fase 3a (no en la fase 2).
- §5: el texto real de los mensajes del índice parcial.

## Correspondencia con el esquema viejo (solo para el renombre)

Es la diferencia esperada mientras #107 no esté mergeado (AGENTS.md); el destino de cada
columna es el de §4.

| Hoy en `main` | Destino (§4) | Repositorio / fase |
|---|---|---|
| `users` | `usuarios` | `sqlite_user` / 3a |
| `projects` | `proyectos` (fechas a `YYYY-MM-DD`, D3) | `sqlite_project` / 3b |
| `project_members` | `integrantes` | `sqlite_project` / 3b |
| `sprints` | `sprints` (mismo nombre y columnas; gana FK, `CHECK` e índices) | `sqlite_sprint` / 3c |
| `backlog_items` | `historias` | `sqlite_backlog` / 3d |
