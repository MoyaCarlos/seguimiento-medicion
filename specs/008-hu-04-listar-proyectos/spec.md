# Feature Specification: Listado de proyectos (GET /projects) — ampliación de HU-04

**Feature Branch**: `feature/HU-04-listar-proyectos`

**Issue**: #109 (parte de #4)

**Created**: 2026-10-23

**Status**: Draft

**Fuente de diseño**: [`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md) §4.1. Esta spec no
redefine tablas ni columnas: el proyecto se lee de `proyectos` con los mismos campos que
`GET /projects/{id}`.

## Contexto

El frontend ya tiene el panel y la edición de un proyecto, pero solo puede abrirlos
conociendo el `id`. La API no ofrece un listado, así que no hay forma de elegir qué
proyecto editar. Esta ampliación agrega únicamente la lectura del listado.

## Objetivo

Permitir que un cliente autenticado como Scrum Master obtenga todos los proyectos
existentes, ordenados por `id`, con la misma forma que `GET /projects/{id}`.

## Entradas y salidas

- **Entrada**: `GET /projects`, sin parámetros ni cuerpo.
- **Salida exitosa**: `200` con un arreglo JSON de proyectos. Cada elemento tiene
  `id`, `nombre`, `descripcion`, `fecha_inicio`, `fecha_fin` y `creado_en`, igual que
  `GET /projects/{id}`.
- **Sin proyectos**: `200` con `[]` (arreglo vacío, no `null` ni error).
- **Error de persistencia**: `500` con el formato de error vigente.

## Reglas de negocio (EARS)

- **WHEN** se consulta `GET /projects`, **THEN** el sistema **shall** devolver todos los
  proyectos persistidos, ordenados ascendentemente por `id`.
- **WHEN** no hay proyectos persistidos, **THEN** el sistema **shall** devolver `[]`.
- **WHEN** un proyecto tiene fechas sin definir, **THEN** el sistema **shall** devolver
  `null` en `fecha_inicio` o `fecha_fin`, igual que la consulta por `id`.
- **IF** falla la lectura de la base, **THEN** el sistema **shall** responder con el error
  interno vigente y no inventar una lista vacía.
- **WHEN** el cliente necesite paginar o filtrar, **THEN** esta historia no lo cubre: el
  listado devuelve todos los proyectos.

## Casos límite y de error

- Un solo proyecto: se devuelve un arreglo con un elemento.
- Proyectos con `id` no consecutivos: se conserva el orden por `id`, sin rellenar huecos.
- Base vacía: `[]`.
- Falla de persistencia: `500`, sin lista parcial.

## Criterios de aceptación

1. **Given** dos o más proyectos persistidos, **When** se consulta `GET /projects`,
   **Then** responde `200` con todos los proyectos ordenados por `id` y con los mismos
   campos que `GET /projects/{id}`.
2. **Given** una base sin proyectos, **When** se consulta `GET /projects`, **Then**
   responde `200` con `[]`.
3. **Given** un proyecto con fechas, **When** se lista, **Then** las fechas se devuelven
   con el mismo formato `AAAA-MM-DD` que la consulta individual.
4. **Given** un fallo de persistencia, **When** se lista, **Then** responde `500` y no
   devuelve una lista vacía como si no hubiera datos.

## Fuera de alcance

- Paginación, filtros, búsqueda y orden configurable.
- Borrado de proyectos.
- Autenticación y autorización por rol.
- Cambios en los endpoints existentes.

## Trazabilidad

| Criterio | Implementación | Pruebas |
|---|---|---|
| Lista ordenada por `id` | `service.ListarProyectos`, `SQLiteProjectRepository.Listar`, `ProjectRepositoryEnMemoria.Listar` | tests de servicio, repositorio SQLite y fake |
| Lista vacía `[]` | servicio y handler inicializan un arreglo vacío | test de servicio y handler |
| Mismos campos que `GET /projects/{id}` | reutiliza `aProyectoResponse` | test de handler |
| Error de persistencia | servicio propaga el error y `escribirError` responde 500 | test del servicio |
| Contrato HTTP | [`specs/002-hu-04-crear-proyecto-equipo/contracts/openapi.yaml`](../002-hu-04-crear-proyecto-equipo/contracts/openapi.yaml) | revisión del contrato |
