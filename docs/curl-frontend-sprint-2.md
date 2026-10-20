# Datos de ejemplo para el frontend del Sprint 2

Estos comandos cargan un proyecto, integrantes, un Sprint y una historia
para probar las pantallas contra la API local.

La API debe estar ejecutándose en `http://localhost:8080`.

## Crear un proyecto

```bash
curl -X POST http://localhost:8080/projects ^
  -H "Content-Type: application/json" ^
  -d "{\"nombre\":\"Proyecto Demo\",\"descripcion\":\"Datos para probar el frontend\",\"creador\":\"Ana\",\"fecha_inicio\":\"2026-10-20\",\"fecha_fin\":\"2026-12-20\"}"
```

La respuesta incluye el `id` del proyecto. En los comandos siguientes se usa
`PROJECT_ID=1` como ejemplo; reemplazarlo por el identificador real.

## Agregar integrantes

```bash
curl -X POST http://localhost:8080/projects/1/members ^
  -H "Content-Type: application/json" ^
  -d "{\"nombre\":\"Carlos\",\"rol\":\"product_builder\"}"
```

```bash
curl http://localhost:8080/projects/1/members
```

## Crear un Sprint

```bash
curl -X POST http://localhost:8080/sprints ^
  -H "Content-Type: application/json" ^
  -d "{\"proyecto_id\":1}"
```

La respuesta incluye el `id` del Sprint. Reemplazar `SPRINT_ID` por ese valor.

```bash
curl -X POST http://localhost:8080/sprints/1/iniciar ^
  -H "Content-Type: application/json" ^
  -d "{\"sprint_goal\":\"Validar el seguimiento\",\"fecha_inicio\":\"2026-10-20\",\"fecha_fin\":\"2026-11-03\"}"
```

## Crear una historia

```bash
curl -X POST http://localhost:8080/backlog ^
  -H "Content-Type: application/json" ^
  -d "{\"proyecto_id\":1,\"titulo\":\"Validar dashboard\",\"descripcion\":\"Historia de ejemplo\",\"prioridad\":\"Must have\"}"
```

## Consultar y actualizar el proyecto

```bash
curl http://localhost:8080/projects/1
```

```bash
curl -X PUT http://localhost:8080/projects/1 ^
  -H "Content-Type: application/json" ^
  -d "{\"nombre\":\"Proyecto Demo actualizado\",\"descripcion\":\"Descripción actualizada\",\"fecha_inicio\":\"2026-10-20\",\"fecha_fin\":\"2026-12-20\"}"
```
