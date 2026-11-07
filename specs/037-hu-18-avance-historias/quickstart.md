# Quickstart: Avance de las Historias del Proyecto (HU-18)

Guía para validar la historia de punta a punta. Contrato: [contracts/openapi.yaml](./contracts/openapi.yaml). Modelo: [data-model.md](./data-model.md).

## Tests automáticos

```bash
go test ./...                 # dominio, service (fakes), repositorio (SQLite :memory:), handler
go test ./features/...        # escenarios Godog, incluido features/avance_historias.feature (@HU-18)
gofmt -l . && go vet ./...    # lo mismo que el check `backend` del CI
cd web && npm run lint && npm test && npm run build   # el check `frontend`
```

Esperado: todo en verde y `gofmt -l .` sin salida.

## Prueba manual contra la API

Prerrequisito: el servidor corriendo (`go run ./cmd/api`) sobre una base nueva.

1. Crear un proyecto y consultar `GET /projects/<proyecto>/avance`.
   Esperado: `200`, `total` en cero y los tres grupos con `cantidad: 0` y `historias: []`.
2. Crear tres historias (`POST /backlog`) y estimar dos con `PUT /backlog/<historia>` (el cuerpo exige
   `titulo`, `descripcion` y `criterios_aceptacion`, además de `estimacion_sp`): una de 5 y una de 3 Story Points;
   la tercera queda sin estimar. Consultar de nuevo.
   Esperado: las tres en `Nueva`, `sprint: null`, cantidad 3 y 8 Story Points (la sin estimar suma 0 y tiene
   `estimacion_sp: null`).
3. Crear e iniciar un Sprint (`POST /sprints`, `POST /sprints/<sprint>/iniciar`), asignarle la de 5 y la de 3
   (`PUT /backlog/<historia>/sprint`) y pasar la de 5 a En progreso y luego a Completada (`PUT /backlog/<historia>/estado`).
   Esperado tras En progreso: la de 5 en `En progreso` (cantidad 1, 5 Story Points). Esperado tras Completada: la
   de 5 en `Completada` (cantidad 1, 5 Story Points) con `sprint` `{id, numero: 1}`; la de 3 sigue en `Nueva`, pero con
   `sprint` `{id, numero: 1}` (está asignada a un Sprint, no en el Product Backlog); la sin estimar sigue en `Nueva` con `sprint: null`.
4. Cerrar el Sprint (`POST /sprints/<sprint>/cerrar`) con la de 3 sin completar.
   Esperado: la de 3 vuelve a `Nueva` con `sprint: null`; la de 5 sigue `Completada` en el Sprint 1.
5. `GET /projects/abc/avance` → `400` con `campo: id`; `GET /projects/9999/avance` → `404` "proyecto no encontrado".
6. `GET /projects/<proyecto>/backlog` (HU-15), después del paso 4.
   Esperado: solo las historias sin Sprint (la de 3 y la sin estimar), sin cambios por esta historia.

## Prueba manual en la interfaz

1. Abrir `/projects/<proyecto>/status`: se ven los tres grupos con cantidad y Story Points, y el Sprint de cada historia ("Sprint 1" o "Product Backlog").
2. Proyecto sin historias: mensaje de que no hay historias y totales en cero.
3. `/projects/abc/status`: mensaje de error sin pedir datos al servidor.
4. Con el servidor detenido: mensaje de error de conexión y ningún dato parcial.

## Resultado de la validación (2026-11-07)

Con la app real: backend Go sobre una base temporal (`DEMO_HABILITADO=true`, `POST /demo/datos`, proyecto 25) y el build de producción servido por el mismo backend, en Chromium headless con ventana de 320 px.

| Paso | Resultado |
|---|---|
| API 1 | OK: `200`, `total` en cero y los tres grupos con cantidad 0 e `historias: []` |
| API 2 | OK: tres historias en `Nueva`, `sprint: null`, cantidad 3 y 8 Story Points; la sin estimar con `estimacion_sp: null`. El paso se reescribió: el `PUT` exige `criterios_aceptacion` |
| API 3 | OK, con lo esperado corregido: tras En progreso y Completada la de 5 cambia de grupo con `sprint` `{id, numero: 1}`; la de 3 queda en `Nueva` con su Sprint. El texto anterior ("una historia en cada grupo") no coincidía con los datos |
| API 4 | OK: la de 3 vuelve a `Nueva` con `sprint: null`; la de 5 sigue `Completada` en el Sprint 1 |
| API 5 | OK: `400` `{"campo":"id","mensaje":"el identificador debe ser un entero positivo"}` y `404` `{"mensaje":"proyecto no encontrado"}` |
| API 6 | OK: `GET /projects/<proyecto>/backlog` devuelve solo las dos historias sin Sprint |
| Extra | OK: `GET /projects/<proyecto>/status` con `Accept: text/html` sirve la SPA y `GET /projects/<proyecto>/avance` con `Accept: application/json` sirve la API |
| UI 1 | OK en el proyecto 25 (más una historia de 120 caracteres sin espacios): los tres grupos con cantidad y Story Points, "Sprint N" o "Product Backlog" en cada historia, "sin estimar" donde corresponde y total 11 historias y 47 Story Points. A 320 px la página no desborda (`scrollWidth` 320 = `clientWidth` 320, ningún elemento pasa de los 320 px) y el título largo se parte en líneas. Con la barra de scroll clásica puede aparecer el desplazamiento de 15 px del `min-width: 320px` de `web/src/App.css`, anterior a esta historia (ver la spec 035) |
| UI 2 | OK en un proyecto recién creado: mensaje "El proyecto todavía no tiene historias.", total 0 y 0, sin grupos ni errores |
| UI 3 | OK en `/projects/abc/status`: alerta `El proyecto "abc" no es válido` y ningún pedido a `/avance`. Además, `/projects/9999/status` muestra "El proyecto no existe" |
| UI 4 | OK, verificado a mano por el Scrum Master (2026-11-07): con el servidor detenido y la pantalla ya cargada, aparece el mensaje de conexión y ningún dato parcial. No se pudo reproducir en el navegador headless; además lo cubre la prueba de componente "US3.4" |
