# Feature Specification: Consultar los Sprints de un proyecto (GET /projects/{id}/sprints) — ampliación de HU-05

**Feature Branch**: `feature/HU-05-consultar-sprints`

**Issue**: #110 (parte de #5)

**Created**: 2026-10-26

**Status**: Backend y frontend implementados, pendiente de revisión por PR.

**Fuente de diseño**: [`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md). Esta spec no toca el
esquema: los Sprints se leen con `SprintRepository.ListarPorProyecto`, que ya existe.

## Contexto

La API permite crear, iniciar y cerrar Sprints, pero no consultarlos: el frontend no puede
mostrar los Sprints de un proyecto. Esta ampliación agrega únicamente la lectura.

## Objetivo

Permitir que un cliente obtenga los Sprints de un proyecto, ordenados por `id`, con la misma
forma que devuelven `POST /sprints` y las demás rutas de Sprints.

## Entradas y salidas

- **Entrada**: `GET /projects/{id}/sprints`, sin cuerpo.
- **Salida exitosa**: `200` con un arreglo JSON. Cada elemento tiene `id`, `proyecto_id`,
  `sprint_goal`, `fecha_inicio`, `fecha_fin` y `estado` (`sprintResponse`).
- **Proyecto sin Sprints**: `200` con `[]` (nunca `null`).
- **Proyecto inexistente**: `404` con el mismo formato de error que `GET /projects/{id}/members`.
- **`id` no numérico o no positivo**: `400` con `campo: "id"`, igual que el resto de las rutas con `{id}`.
- **Error de persistencia**: `500` con el formato de error vigente.

## Reglas de negocio (EARS)

- **WHEN** se consulta `GET /projects/{id}/sprints` de un proyecto existente, **THEN** el
  sistema **shall** devolver todos sus Sprints ordenados ascendentemente por `id`.
- **WHEN** el proyecto no tiene Sprints, **THEN** el sistema **shall** devolver `[]`.
- **WHEN** un Sprint está Pendiente, **THEN** el sistema **shall** devolver `null` en
  `sprint_goal`, `fecha_inicio` y `fecha_fin`, igual que `POST /sprints`.
- **IF** el proyecto no existe, **THEN** el sistema **shall** responder `404` y no una lista vacía.
- **IF** `{id}` no es un entero positivo, **THEN** el sistema **shall** responder `400`.
- **IF** falla la lectura de la base, **THEN** el sistema **shall** responder `500` sin lista parcial.

## Casos límite y de error

- Sprints creados en desorden de `id` en el repositorio: se ordena por `id` en el caso de uso.
- Un solo Sprint: arreglo de un elemento.
- Sprints de otro proyecto: no aparecen.
- `id` = `abc`, `0` o `-1`: `400`.

## Criterios de aceptación

1. **Given** un proyecto con varios Sprints, **When** se consulta `GET /projects/{id}/sprints`,
   **Then** responde `200` con sus Sprints ordenados por `id`.
2. **Given** un proyecto sin Sprints, **When** se consulta, **Then** responde `200` con `[]`.
3. **Given** un proyecto inexistente, **When** se consulta, **Then** responde `404`.
4. **Given** un `id` no numérico, **When** se consulta, **Then** responde `400`.

## Frontend: carga de Sprints en la pantalla `/projects/:id/sprints`

Reemplaza el stub `cargarSprintsIniciales()` de la spec [011](../011-hu-05-frontend-sprints/spec.md).

### Requisitos (EARS)

- **FR-F01**: **WHEN** se abre la pantalla de Sprints de un proyecto, **THEN** el sistema **shall** pedir `GET /projects/{id}/sprints` y usar el resultado como estado inicial.
- **FR-F02**: **WHILE** se carga el listado, el sistema **shall** mostrar el indicador de carga.
- **FR-F03**: **IF** el listado falla, **THEN** el sistema **shall** mostrar el mensaje de la API tal cual en el aviso de la pantalla, sin romperla (las secciones y "Crear Sprint" siguen disponibles).
- **FR-F04**: **WHEN** cambia el `:id` del proyecto, **THEN** el sistema **shall** volver a cargar y no dejar a la vista los Sprints del anterior.
- **FR-F05**: el sistema **shall** dejar de mostrar el texto de limitación de la spec 011 (FR-018).
- **FR-F06**: **WHEN** se crea, inicia o cierra un Sprint, **THEN** el sistema **shall** seguir actualizándolo en el estado (`guardarSprint`).
- Se muestra "Sprint #<id>"; el campo `numero` es el issue #129.

### Criterios de aceptación

5. **Given** un proyecto con Sprints, **When** se abre la pantalla, **Then** se ven agrupados por estado tras la carga.
6. **Given** un proyecto sin Sprints, **When** se abre, **Then** las tres secciones muestran "No hay Sprints".
7. **Given** que el listado falla, **When** se abre, **Then** se ve el mensaje de la API y la pantalla sigue usable.
8. **Given** la pantalla abierta, **When** cambia el proyecto de la URL, **Then** se recarga el listado del nuevo proyecto.

## Fuera de alcance

- El campo `numero` del Sprint (issue #129, PR aparte).
- Paginación y filtros por estado.
- Cambios en `SprintRepository` o en el esquema.

## Clarificaciones

- **¿Dónde vive el handler?** En `SprintHandler` (ruta de Sprints), reutilizando `idDeRuta`;
  el `404` sale de `escribirError`, el mismo formato que `members`.
- **¿Quién ordena?** El caso de uso, como `ListarProyectos`, para no depender del orden del repositorio.
- **¿Quién valida el proyecto?** El caso de uso, como `ListarIntegrantes` (`ObtenerPorID`).

## Trazabilidad

| Criterio | Implementación | Pruebas |
|---|---|---|
| 1, 2, 3 | `service.ConsultarSprints` | `consultar_sprints_test.go` |
| 1, 2, 3, 4 | `SprintHandler.Listar`, ruta en `cmd/api/main.go` | `sprint_handler_test.go` |
| FR-F01 a FR-F06 | `web/src/features/sprints/index.tsx`, `api.ts` (`listarSprints`) | `web/src/features/sprints/index.test.tsx` |
| Escenarios BDD | [`features/consultar_sprints.feature`](../../features/consultar_sprints.feature) | Godog (`features/steps_consultar_sprints.go`) |
