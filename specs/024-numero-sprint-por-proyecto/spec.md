# Feature Specification: Numerar los Sprints por proyecto (Issue #129)

**Feature Branch**: `fix/129-numero-sprint-por-proyecto`

**Created**: 2026-11-01

**Status**: Implementada (Issue #129), pendiente de revisión por PR.

**Input**: Issue #129: el número que identifica a un Sprint es su `id`, la clave primaria autoincremental y **global** de la tabla `sprints`. Si el proyecto 1 tuvo los Sprints #1 y #2, el primer Sprint del proyecto 2 aparece como "Sprint #3".

**Número de spec**: el issue pedía `--number 12`, pero `specs/012` ya está ocupada por HU-02. Se usa el siguiente libre (AGENTS.md).

## Objetivo

Que cada proyecto numere sus Sprints desde 1, para que su lista tenga un orden entendible, sin cambiar el `id` que usan las URLs de la API.

## Decisión

**Decisión del Scrum Master (2026-10-25): alternativa B, número derivado.** `numero` se calcula al leer: es la posición del Sprint entre los de su proyecto, ordenados por `id`. No hay migración ni cambios en las tablas ni en `MODELO-DE-DATOS.md`.

**Límite conocido**: `numero` deja de ser exacto si algún día se borran Sprints (hoy no existe ese borrado). Si se agrega, se pasa a una columna guardada en su propio PR con migración.

## User Story 1 - Ver el número del Sprint dentro de su proyecto (Priority: P1)

Como consumidor de la API, quiero recibir el número de cada Sprint dentro de su proyecto, para mostrar "Sprint 1, 2, 3..." sin importar los Sprints de otros proyectos.

**Acceptance Scenarios**:

1. **Given** el proyecto 1 con dos Sprints, **When** se crea el primer Sprint del proyecto 2, **Then** la respuesta trae `numero` 1.
2. **Given** un proyecto con un Sprint, **When** se crea el siguiente, **Then** la respuesta trae `numero` 2.
3. **Given** un Sprint con `numero` N, **When** se inicia o se cierra, **Then** la respuesta sigue trayendo `numero` N.
4. **Given** dos proyectos con Sprints intercalados, **When** se consultan los Sprints de cada uno, **Then** cada lista numera desde 1, por `id` ascendente.

## Entradas y salidas

- **Entradas**: las mismas de hoy. `numero` no se recibe nunca.
- **Salidas**: `POST /sprints`, `POST /sprints/{id}/iniciar`, `POST /sprints/{id}/cerrar` y `GET /projects/{id}/sprints` agregan a cada Sprint el campo entero `numero`. El resto del contrato no cambia.

## Reglas de negocio y casos límite (EARS)

- El sistema **shall** asignar a cada Sprint un `numero` igual a la cantidad de Sprints de su proyecto con `id` menor o igual al suyo.
- **WHEN** se inicia o se cierra un Sprint, **THEN** el sistema **shall** conservar su `numero`.
- **WHEN** se crean Sprints en otro proyecto, **THEN** el sistema **shall** mantener el `numero` de los Sprints de este proyecto.
- **IF** un proyecto no tiene Sprints, **THEN** el sistema **shall** seguir respondiendo un arreglo vacío.
- **IF** la operación falla (proyecto o Sprint inexistente, validación), **THEN** el sistema **shall** responder el mismo error que antes de este cambio.

## Restricciones

- El `id` sigue siendo el identificador de la API y de las referencias entre tablas: `numero` solo sirve para mostrar.
- Sin cambios de esquema.
- El frontend pasa de "Sprint #<id>" a "Sprint <numero>" en un PR aparte, como indica el issue.

## Criterios de aceptación

- [ ] El primer Sprint de cualquier proyecto tiene número 1, y los siguientes 2, 3... dentro de ese proyecto.
- [ ] La API devuelve `numero` en las respuestas que hoy devuelven un Sprint, sin cambiar el resto del contrato.
- [ ] `go build`, `go vet`, `gofmt -l` y `go test ./...` en verde; los escenarios de Godog existentes no cambian de comportamiento.
