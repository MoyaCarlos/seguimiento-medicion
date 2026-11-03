# Quickstart: Consulta de Sprints Anteriores (HU-12)

Guía para validar la historia de punta a punta. Contrato: [contracts/openapi.yaml](./contracts/openapi.yaml).
Modelo: [data-model.md](./data-model.md).

**Prerrequisito**: el cálculo de métricas de HU-03 mergeado en `main` y el contrato revisado (research.md R2).

## Tests automáticos

```bash
go test ./...                 # dominio, services (fakes), repositorio (SQLite :memory:), handlers
go test ./features/...        # escenarios Godog, incluido features/historial_sprints.feature (@HU-12)
gofmt -l . && go vet ./...    # lo mismo que el check `backend` del CI
cd web && npm run lint && npm test && npm run build   # el check `frontend`
```

Esperado: todo en verde y `gofmt -l .` sin salida.

## Prueba manual contra la API

Prerrequisito: el servidor corriendo (`go run ./cmd/api`) sobre una base nueva.

1. Crear un proyecto (`POST /projects`) y consultar `GET /projects/<proyecto>/historial`.
   Esperado: `200` con `[]`.
2. Crear e iniciar un Sprint, y asignarle dos historias estimadas, por ejemplo de 5 y 3 Story Points.
   Pedir su detalle:
   ```bash
   curl -s localhost:8080/projects/<proyecto>/historial/<sprint>
   ```
   Esperado: `400` con `campo: estado` y el mensaje "el Sprint todavía no está finalizado". El
   historial sigue vacío.
3. Imputar horas a ambas historias, registrar dos defectos y resolver uno. Completar la historia de
   5 Story Points y cerrar el Sprint (`POST /sprints/<sprint>/cerrar`).
4. Consultar el historial. Esperado: un Sprint con `numero`, `sprint_goal`, fechas y `metricas`, con
   8 Story Points planificados, 5 completados y 50% de historias completadas.
5. Consultar el detalle. Esperado:
   - La historia de 5 con `completada: true` y la de 3 con `completada: false`.
   - La suma de `minutos` de las historias equivale a `metricas.horas_reales`.
   - Dos defectos: uno detectado y resuelto en el Sprint, y otro solo detectado.
6. Editar los Story Points de la historia arrastrada, que volvió al Product Backlog, y repetir el
   paso 5. Esperado: sigue figurando con 3 Story Points.
7. Iniciar un segundo Sprint con esa historia, completarla, resolver el defecto pendiente y cerrar.
   Esperado: el historial lista dos Sprints en orden de número. En el detalle del primero la historia
   sigue como no completada; en el del segundo, como completada. El defecto aparece en ambos, con
   `resuelto_en_sprint: true` solo en el segundo.
8. Casos de error:

   | Pedido | Esperado |
   |---|---|
   | `GET /projects/9999/historial` | `404` "proyecto no encontrado" |
   | `GET /projects/abc/historial` | `400` con `campo: id` |
   | `GET /projects/<proyecto>/historial/9999` | `404` "sprint no encontrado" |
   | `GET /projects/<proyecto>/historial/abc` | `400` con `campo: sprint_id` |
   | Detalle con un Sprint de otro proyecto | `404` "sprint no encontrado" |

## Prueba manual de la interfaz

Con la aplicación levantada, abrir un proyecto y entrar a "Historial de Sprints":

1. Sin Sprints cerrados se ve un mensaje que indica que todavía no hay ninguno.
2. Con Sprints cerrados se ve cada uno como "Sprint <número>", con su Sprint Goal, fechas y métricas.
3. Al seleccionar uno se ven sus métricas, historias, horas y defectos, y se puede volver a la lista.
4. Sin horas por Story Point, las horas estimadas y la desviación aparecen como "sin datos"; una
   historia sin estimación aparece como "sin estimar".
5. Recargar el detalle (F5) vuelve a mostrar la misma pantalla.
