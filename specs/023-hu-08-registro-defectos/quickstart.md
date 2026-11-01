# Quickstart: Registro y Seguimiento de Defectos (HU-08)

Guía para validar la historia de punta a punta. Contrato: [contracts/openapi.yaml](./contracts/openapi.yaml).
Modelo: [data-model.md](./data-model.md).

## Tests automáticos

```bash
go test ./...                 # dominio, services (fakes), repositorio (SQLite :memory:), handlers
go test ./features/...        # escenarios Godog, incluido features/registro_defectos.feature (@HU-08)
gofmt -l . && go vet ./...    # lo mismo que el check `backend` del CI
cd web && npm test            # el check `frontend`: solo cambia el proxy de vite.config.ts
```

Esperado: todo en verde, `gofmt -l .` sin salida y `TestPrefijosAPI_IncluyenLosQueProxeaVite` en
verde con el prefijo `/defectos`.

## Prueba manual contra la API

Prerrequisito: el servidor corriendo (`go run ./cmd/api`) sobre una base nueva (la migración 005
se aplica sola al arrancar).

1. Crear un proyecto (`POST /projects`) y anotar su `id`.
2. Crear dos historias (`POST /backlog`). Sin iniciar ningún Sprint, intentar registrar un defecto:
   ```bash
   curl -s -X POST localhost:8080/projects/<proyecto>/defectos \
     -H 'Content-Type: application/json' \
     --data-binary @defecto.json
   ```
   con `defecto.json` = `{"historia_id": <historia1>, "descripcion": "El total no se actualiza", "severidad": "Alta"}`.
   Esperado: `400` con `campo: sprint_deteccion_id`.
3. Crear un Sprint (`POST /sprints`) e iniciarlo (`POST /sprints/{id}/iniciar`). Repetir el paso 2.
   Esperado: `201`, `estado: Pendiente`, `sprint_deteccion_id` = el Sprint activo y
   `sprint_resolucion_id: null`. La historia no hace falta que esté en el Sprint.
4. Casos de error (cada uno `400` con el `campo` indicado; el listado no cambia):

   | Cuerpo | `campo` |
   |---|---|
   | sin `historia_id` | `historia_id` |
   | `"historia_id": 9999` | `historia_id` |
   | `"historia_id"` de una historia de otro proyecto | `historia_id` |
   | `"descripcion": "   "` | `descripcion` |
   | `"descripcion"` de 2001 caracteres | `descripcion` |
   | `"severidad": "alta"` | `severidad` |
   | `"severidad": "Crítica"` | `severidad` |

5. Listar: `curl -s localhost:8080/projects/<proyecto>/defectos`. Esperado: `200` con el defecto
   del paso 3 y el título de su historia.
6. Resolver: `curl -s -X POST localhost:8080/defectos/<defecto>/resolver`. Esperado: `200`,
   `estado: Resuelto` y `sprint_resolucion_id` = el Sprint activo. Resolverlo otra vez: `400` con
   `campo: estado`.
7. Registrar un segundo defecto, cerrar el Sprint (`POST /sprints/{id}/cerrar`) e intentar
   resolverlo: `400` con `campo: sprint_resolucion_id`. Crear e iniciar otro Sprint, y resolverlo:
   `sprint_resolucion_id` es el Sprint nuevo y `sprint_deteccion_id` sigue siendo el anterior.
8. `POST /projects/9999/defectos` y `GET /projects/9999/defectos`: `404` "proyecto no
   encontrado". `POST /defectos/9999/resolver`: `404` "defecto no encontrado".
9. Reiniciar el servidor y listar otra vez: los defectos siguen con sus Sprints (FR-015).

En Windows, mandar los cuerpos con acentos desde un archivo (`--data-binary @archivo.json`), como
se anotó en HU-07: pasados como argumento, la consola los recodifica.

## Resultado de la validación manual (2026-10-31)

Contra `go run ./cmd/api` sobre una base nueva, con los cuerpos enviados en UTF-8:

| Paso | Resultado |
|---|---|
| 1–2 | Sin Sprint activo: `400` con `campo: sprint_deteccion_id` |
| 3 | `201`, `estado: Pendiente`, `sprint_deteccion_id` del Sprint activo, `sprint_resolucion_id: null`, título de la historia en `historia` |
| 4 | Los 7 cuerpos inválidos dan `400` con el `campo` esperado; el listado no cambia |
| 5 | `200` con el defecto del paso 3 |
| 6 | `200` `Resuelto` en el Sprint activo; la segunda vez `400` con `campo: estado` |
| 7 | Sin Sprint activo: `400` con `campo: sprint_resolucion_id`; con un Sprint nuevo, resuelto en él y detectado en el anterior |
| 8 | `404` "proyecto no encontrado" (POST y GET) y `404` "defecto no encontrado" |
| 9 | Después de reiniciar, los dos defectos siguen con sus Sprints de detección y resolución |
