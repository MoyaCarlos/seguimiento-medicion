# Datos de ejemplo para el frontend del Sprint 2

Estos comandos cargan un proyecto, integrantes, un Sprint y una historia
para probar las pantallas contra la API local.

La API debe estar ejecutándose en `http://localhost:8080`.

Hay dos versiones con los mismos datos: la de **bash** (Linux, macOS y Git Bash
en Windows) y la de **cmd** (Windows). Usá la que corresponda a tu terminal;
solo cambia la forma de escribir los comandos.

Los ids que devuelve la API se pierden al recargar una pantalla, porque todavía
no existe un listado de proyectos ni de Sprints. **Anotá el `id` de cada
respuesta**: lo necesitás para abrir la pantalla en
`http://localhost:5173/projects/<id>` y para los comandos siguientes, donde el
ejemplo usa `1`.

## Versión bash (Linux, macOS y Git Bash)

### Crear un proyecto

```bash
curl -X POST http://localhost:8080/projects \
  -H 'Content-Type: application/json' \
  -d '{"nombre":"Proyecto Demo","descripcion":"Datos para probar el frontend","creador":"Ana","fecha_inicio":"2026-10-20","fecha_fin":"2026-12-20"}'
```

La respuesta incluye el `id` del proyecto. En los comandos siguientes se usa
`1` como ejemplo; reemplazalo por el id real.

### Agregar integrantes

```bash
curl -X POST http://localhost:8080/projects/1/members \
  -H 'Content-Type: application/json' \
  -d '{"nombre":"Carlos","rol":"product_builder"}'
```

```bash
curl http://localhost:8080/projects/1/members
```

### Crear un Sprint

```bash
curl -X POST http://localhost:8080/sprints \
  -H 'Content-Type: application/json' \
  -d '{"proyecto_id":1}'
```

La respuesta incluye el `id` del Sprint. En los comandos siguientes se usa
`1` como ejemplo; reemplazalo por el id real.

```bash
curl -X POST http://localhost:8080/sprints/1/iniciar \
  -H 'Content-Type: application/json' \
  -d '{"sprint_goal":"Validar el seguimiento","fecha_inicio":"2026-10-20","fecha_fin":"2026-11-03"}'
```

### Cerrar un Sprint

Cerralo cuando termines de probarlo: el comando no lleva cuerpo. Se usa `1` como
ejemplo; reemplazalo por el id real.

```bash
curl -X POST http://localhost:8080/sprints/1/cerrar
```

La respuesta es el Sprint con `estado` en `Finalizado`. Las historias sin completar
vuelven al Product Backlog en estado `Nueva`. Cerrar un Sprint que no está Activo
responde 400 con `{"campo":"estado","mensaje":"solo se puede cerrar un Sprint activo"}`.
Si el Sprint no existe, responde 404 con `{"mensaje":"sprint no encontrado"}`.

### Crear una historia

```bash
curl -X POST http://localhost:8080/backlog \
  -H 'Content-Type: application/json' \
  -d '{"proyecto_id":1,"titulo":"Validar dashboard","descripcion":"Historia de ejemplo","prioridad":"M"}'
```

### Asignar una historia a un Sprint

```bash
curl -X PUT http://localhost:8080/backlog/1/sprint \
  -H 'Content-Type: application/json' \
  -d '{"sprint_id":1}'
```

La respuesta es la historia con su `sprint_id`. Repetir la asignación responde 400 con
`{"campo":"sprint_id","mensaje":"la historia ya pertenece a un Sprint"}`.

### Consultar el Sprint Backlog

```bash
curl http://localhost:8080/sprints/1/backlog
```

Devuelve las historias del Sprint en orden M, S, C, W y luego por creación; `[]` si
no tiene historias. Si el Sprint no existe, responde 404 con
`{"mensaje":"sprint no encontrado"}`.

### Cambiar el estado de una historia

La historia debe estar asignada a un Sprint activo. Primero se pasa a
`En progreso` y después a `Completada`; el estado solo avanza.

```bash
curl -X PUT http://localhost:8080/backlog/1/estado \
  -H 'Content-Type: application/json' \
  -d '{"estado":"En progreso"}'
```

```bash
curl -X PUT http://localhost:8080/backlog/1/estado \
  -H 'Content-Type: application/json' \
  -d '{"estado":"Completada"}'
```

La respuesta es la historia con su `estado` y, al completarla, su
`fecha_finalizacion`. Los rechazos son `400` con `{"campo":"estado","mensaje":...}`,
`404` si la historia no existe y `409` si otro cambio la modificó antes.

### Consultar y actualizar el proyecto

```bash
curl http://localhost:8080/projects/1
```

```bash
curl -X PUT http://localhost:8080/projects/1 \
  -H 'Content-Type: application/json' \
  -d '{"nombre":"Proyecto Demo actualizado","descripcion":"Descripción actualizada","fecha_inicio":"2026-10-20","fecha_fin":"2026-12-20"}'
```

## Versión cmd (Windows)

Mismos pasos y mismos datos que la versión bash. En cmd la continuación de
línea es `^` y las comillas dobles del JSON se escriben con `\"`.

### Crear un proyecto

```bat
curl -X POST http://localhost:8080/projects ^
  -H "Content-Type: application/json" ^
  -d "{\"nombre\":\"Proyecto Demo\",\"descripcion\":\"Datos para probar el frontend\",\"creador\":\"Ana\",\"fecha_inicio\":\"2026-10-20\",\"fecha_fin\":\"2026-12-20\"}"
```

La respuesta incluye el `id` del proyecto. En los comandos siguientes se usa
`1` como ejemplo; reemplazalo por el id real.

### Agregar integrantes

```bat
curl -X POST http://localhost:8080/projects/1/members ^
  -H "Content-Type: application/json" ^
  -d "{\"nombre\":\"Carlos\",\"rol\":\"product_builder\"}"
```

```bat
curl http://localhost:8080/projects/1/members
```

### Crear un Sprint

```bat
curl -X POST http://localhost:8080/sprints ^
  -H "Content-Type: application/json" ^
  -d "{\"proyecto_id\":1}"
```

La respuesta incluye el `id` del Sprint. En los comandos siguientes se usa
`1` como ejemplo; reemplazalo por el id real.

```bat
curl -X POST http://localhost:8080/sprints/1/iniciar ^
  -H "Content-Type: application/json" ^
  -d "{\"sprint_goal\":\"Validar el seguimiento\",\"fecha_inicio\":\"2026-10-20\",\"fecha_fin\":\"2026-11-03\"}"
```

### Cerrar un Sprint

Cerralo cuando termines de probarlo: el comando no lleva cuerpo. Se usa `1` como
ejemplo; reemplazalo por el id real.

```bat
curl -X POST http://localhost:8080/sprints/1/cerrar
```

La respuesta es el Sprint con `estado` en `Finalizado`. Las historias sin completar
vuelven al Product Backlog en estado `Nueva`. Cerrar un Sprint que no está Activo
responde 400 con `{"campo":"estado","mensaje":"solo se puede cerrar un Sprint activo"}`.
Si el Sprint no existe, responde 404 con `{"mensaje":"sprint no encontrado"}`.

### Crear una historia

```bat
curl -X POST http://localhost:8080/backlog ^
  -H "Content-Type: application/json" ^
  -d "{\"proyecto_id\":1,\"titulo\":\"Validar dashboard\",\"descripcion\":\"Historia de ejemplo\",\"prioridad\":\"M\"}"
```

### Asignar una historia a un Sprint

```bat
curl -X PUT http://localhost:8080/backlog/1/sprint ^
  -H "Content-Type: application/json" ^
  -d "{\"sprint_id\":1}"
```

La respuesta es la historia con su `sprint_id`. Repetir la asignación responde 400 con
`{"campo":"sprint_id","mensaje":"la historia ya pertenece a un Sprint"}`.

### Consultar el Sprint Backlog

```bat
curl http://localhost:8080/sprints/1/backlog
```

Devuelve las historias del Sprint en orden M, S, C, W y luego por creación; `[]` si
no tiene historias. Si el Sprint no existe, responde 404 con
`{"mensaje":"sprint no encontrado"}`.

### Cambiar el estado de una historia

La historia debe estar asignada a un Sprint activo. Primero se pasa a
`En progreso` y después a `Completada`; el estado solo avanza.

```bat
curl -X PUT http://localhost:8080/backlog/1/estado ^
  -H "Content-Type: application/json" ^
  -d "{\"estado\":\"En progreso\"}"
```

```bat
curl -X PUT http://localhost:8080/backlog/1/estado ^
  -H "Content-Type: application/json" ^
  -d "{\"estado\":\"Completada\"}"
```

La respuesta es la historia con su `estado` y, al completarla, su
`fecha_finalizacion`. Los rechazos son `400` con `{"campo":"estado","mensaje":...}`,
`404` si la historia no existe y `409` si otro cambio la modificó antes.

### Consultar y actualizar el proyecto

```bat
curl http://localhost:8080/projects/1
```

```bat
curl -X PUT http://localhost:8080/projects/1 ^
  -H "Content-Type: application/json" ^
  -d "{\"nombre\":\"Proyecto Demo actualizado\",\"descripcion\":\"Descripción actualizada\",\"fecha_inicio\":\"2026-10-20\",\"fecha_fin\":\"2026-12-20\"}"
```
