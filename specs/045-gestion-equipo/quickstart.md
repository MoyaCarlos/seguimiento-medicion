# Quickstart: validar la gestión del equipo (HU-16)

## Pruebas automatizadas

```sh
go test ./...                          # dominio, servicios, repositorio SQLite, HTTP y Godog
go test ./features -run TestFeatures   # solo los escenarios BDD (incluye gestion_equipo.feature)
cd web && npm test                     # pruebas de componente de la pantalla del equipo
```

Esperado: todo en verde, incluidos los escenarios existentes de `creacion_proyecto_equipo`,
`permisos_por_rol`, `scrum_master_unico`, `votacion_planning_poker` e `imputacion_horas`.

## Prueba manual con curl

Prerrequisito: servidor levantado (`go run ./cmd/api`), dos usuarios registrados (`ana` como
Scrum Master que crea el proyecto, `beto` como Product Builder) y los tokens de cada uno
(ver `docs/` guía de curl).

1. `ana` agrega a `beto` (`POST /projects/{id}/members`) → 201.
2. `beto` vota en una ronda y la ronda se revela; `beto` imputa horas.
3. `ana` quita a `beto`: `DELETE /projects/{id}/members/{id_beto}` → 204.
4. `GET /projects/{id}/members` → `beto` ya no figura.
5. Historial de rondas y horas de la historia → `beto` sigue con su nombre y sus valores;
   métricas sin cambios.
6. Con el token de `beto`: `GET /projects/{id}` → 403 «No sos integrante de este proyecto».
7. `ana` se intenta quitar a sí misma → 409 «No se puede quitar al Scrum Master del
   proyecto».
8. `ana` repite el paso 3 → 404 «integrante no encontrado».
9. `ana` vuelve a agregar a `beto` → 201, aparece una sola vez como Product Builder.
10. `DELETE /projects/9999/members/1` → 404 «proyecto no encontrado».

Contrato completo: [contracts/api-equipo.md](contracts/api-equipo.md).
