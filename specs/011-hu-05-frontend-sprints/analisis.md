# Análisis de consistencia: spec 011 (HU-05 Frontend, Issue #26)

**Fecha**: 2026-10-24 · **Artefactos**: spec.md, plan.md, tasks.md, constitución v1.0.0 y el código actual de `web/` · Solo lectura, antes de implementar.

**Resultado**: 0 CRITICAL, 0 HIGH, 3 MEDIUM, 6 LOW. Se puede continuar; la columna "Resolución" la completa el Scrum Master o quien aplique cada cambio.

## Verificaciones pedidas

1. **`web/src/rutas.test.tsx` (FR-020)**: hoy `npm test` pasa (4 archivos, 39 tests). El test de `/projects/42/sprints` hace `findByRole('heading', { name: 'Sprints' })` (exige un único encabezado con ese nombre exacto) y `getByText('42')` (exige un único elemento cuyo texto sea exactamente "42"); el de cambio de id usa `findByText('1', { selector: 'strong' })`. Con la pantalla prevista (un `<h1>` "Sprints", un único `<strong>` con el id, "Sprint #<id>" como un solo texto, sin llamadas a la API al montar) debería seguir en verde, siempre que se cumpla el hallazgo U2. T035 lo comprueba antes del commit GREEN.
2. **Archivos fuera de `web/src/features/sprints/` y `specs/011-*/`**: ninguna tarea modifica `compartido/`, `package.json`, README, `rutas.tsx`, `rutas.test.tsx` ni otras features; `aislamiento.test.ts`, `rutas.test.tsx` y `vitest` solo se ejecutan. Excepción informada: T046 escribe fuera del repositorio (ver S1).
3. **Cobertura de FR**: 20 de 20 tienen tarea; FR-017 solo de forma implícita (ver C1).

## Hallazgos

| ID | Categoría | Severidad | Ubicación | Resumen | Recomendación | Resolución |
|----|-----------|-----------|-----------|---------|---------------|------------|
| L1 | Calidad / lint | MEDIUM | tasks.md T022, T024; spec FR-007; `web/.oxlintrc.json` | `BotonCerrar.tsx` exportaría un componente y la función `cerrarSprint`. La regla `react/only-export-components` (warn) lo marca (comprobado con oxlint): `npm run lint` mostraría una advertencia, y la spec 007 pide lint sin errores ni advertencias | Dejar `cerrarSprint` sin exportar dentro de `BotonCerrar.tsx` y probarla desde el componente (T022 sin la prueba directa (a)), o decidir explícitamente tolerar la advertencia | Aplicado en tasks.md: `cerrarSprint` queda sin exportar dentro de `BotonCerrar.tsx` y se prueba a través del componente (clic → `POST /sprints/{id}/cerrar` sin cuerpo); se quitó la prueba directa (a) de T022; T024 aclara el motivo (`react/only-export-components` de oxlint) y que FR-007 se mantiene. |
| U1 | Subespecificación | MEDIUM | tasks.md T034; spec FR-013 | T034 no dice que el error de "Crear Sprint" pase por `clasificarError(error, [])`. Sin eso, un `ApiError` con campo (`proyecto_id`) o un error de conexión podrían no mostrarse | Agregar a T034: errores de crear vía `clasificarError(error, [])`, siempre al aviso de la pantalla | Aplicado en T034: los errores de "Crear Sprint" pasan por `clasificarError(error, [])` y se muestran siempre en el aviso de la pantalla. |
| U2 | Subespecificación | MEDIUM | tasks.md T020, T028, T034; `web/src/rutas.test.tsx` | (a) T020 no fija `name` de los campos: deben ser `sprint_goal`, `fecha_inicio` y `fecha_fin` para que `error.field` coincida, ni un `id` por Sprint. (b) Un segundo encabezado llamado exactamente "Sprints" rompería `findByRole` en `rutas.test.tsx` | (a) Fijar `name` e `id` (con el id del Sprint) en T020. (b) Aclarar en T028/T034 que las secciones se titulan "Pendiente", "Activo" y "Finalizado" y que ningún otro encabezado se llama "Sprints" | Aplicado: T020 fija `name` `sprint_goal`, `fecha_inicio`, `fecha_fin` e `id` propio con el id del Sprint; T028 y T034 aclaran que las secciones se titulan "Pendiente", "Activo" y "Finalizado" y que ningún otro encabezado se llama "Sprints". |
| S1 | Alcance | LOW | tasks.md T046 | T046 escribe `~/Escritorio/PR-26-frontend-sprints.md`, fuera del repositorio. Es lo que pidió el Scrum Master, pero es la única tarea fuera de `web/src/features/sprints/` y `specs/011-*/` | Dejar como excepción documentada; no se agrega al repositorio | Aceptado sin cambios (excepción documentada: el archivo vive fuera del repositorio). |
| C1 | Cobertura | LOW | spec FR-017 | "Interfaz en español" no tiene tarea ni test propio; queda cubierto por los textos de T032 | Aceptable; revisar en T040 | Aceptado sin cambios. |
| C2 | Cobertura | LOW | spec, Casos límite (id de proyecto no válido) | Ninguna tarea prueba un id no numérico en la URL (`Number('abc')` = `NaN`) | Opcional: un caso en `index.test.tsx` que muestre el error de la API sin validar el id | Aplicado en T032: nuevo caso (h) con un id no numérico en la URL; se muestra el error de la API sin que la pantalla valide el id. |
| I1 | Inconsistencia | LOW | plan.md "Source Code" vs tasks.md T018, T022 | `FormularioInicio.test.tsx` y `BotonCerrar.test.tsx` figuran en tasks.md y no en el árbol del plan (la nota de tasks lo aclara) | Sumarlos al árbol del plan si se quiere coherencia | Aplicado en plan.md: `FormularioInicio.test.tsx` y `BotonCerrar.test.tsx` suman al árbol de archivos. |
| U3 | Subespecificación | LOW | tasks.md T034 | "Quitar las clases/estilos del placeholder" es ambiguo: esos estilos están en `web/src/App.css`, fuera de la carpeta de la feature | Aclarar: solo se dejan de usar las clases en `index.tsx`; no se edita `App.css` ni `index.css` | Aplicado en T034: solo se dejan de usar las clases del placeholder en `index.tsx`; no se edita `App.css` ni `index.css`. |
| A1 | Ambigüedad | LOW | spec SC-001, SC-002 | Los tiempos (2 segundos, 2 minutos) dependen de la máquina y solo se miden a mano | Aceptable: se registran en T040 | Aceptado sin cambios. |

## Cobertura de requisitos

| Requisito | ¿Tiene tarea? | Tareas |
|---|---|---|
| FR-001, FR-002 | Sí | T005, T032, T034 |
| FR-003, FR-004 | Sí | T005, T018, T020, T028 |
| FR-005, FR-006, FR-007 | Sí | T022, T024, T028 (ver L1) |
| FR-008, FR-009 | Sí | T004, T012, T014, T028 |
| FR-010 | Sí | T012, T014, T034 |
| FR-011, FR-012, FR-013 | Sí | T008, T010, T020, T024 (ver U1) |
| FR-014 | Sí | T018, T022, T032 |
| FR-015 | Sí | T038, T041 |
| FR-016, FR-020 | Sí | T032, T034, T035 |
| FR-017 | Implícita | T032 (ver C1) |
| FR-018, FR-019 | Sí | T028, T032, T034 |
| SC-001 a SC-006 | Sí | T012, T024, T038, T040, T041 |

## Alineación con la constitución

Sin conflictos. Principio I: RED/GREEN con commit propio en la lógica; tipos y `api.ts` son envoltorios sin lógica. Principio II: análisis hecho antes de implementar. Principio III: no aplica. Principio VI: sin dependencias nuevas. Flujo de Git: el merge con "Create a merge commit" figura en la descripción del PR (T046).

## Tareas sin requisito

Ninguna. Las tareas de commit, verificación manual y descripción del PR responden al flujo del proyecto.

## Métricas

- Requisitos funcionales: 20 · Criterios de éxito: 6
- Tareas: 46
- Cobertura de FR con al menos una tarea: 100 % (20/20; FR-017 implícita)
- Ambigüedades: 1 · Duplicaciones: 0 · CRITICAL: 0 · HIGH: 0
