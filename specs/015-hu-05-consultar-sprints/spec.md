# Feature Specification: Consultar los Sprints de un proyecto (GET /projects/{id}/sprints) — ampliación de HU-05

**Feature Branch**: `feature/HU-05-consultar-sprints`

**Issue**: #110 (parte de #5)

**Created**: 2026-10-26

**Status**: En implementación.

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

## Fuera de alcance

- El campo `numero` del Sprint (issue #129, PR aparte).
- Paginación, filtros por estado, frontend.
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
| Escenarios BDD | [`features/consultar_sprints.feature`](../../features/consultar_sprints.feature) | Godog (`features/steps_consultar_sprints.go`) |
