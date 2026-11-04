# Quickstart: Cálculo de Métricas del Proyecto (HU-03)

Guía para validar la historia de punta a punta. El contrato está en
[openapi.yaml](./contracts/openapi.yaml) y las reglas en [data-model.md](./data-model.md).

## Pruebas automatizadas

```bash
go test ./...          # incluye dominio, services, repositorio, handlers y Godog
gofmt -l .             # no debe listar archivos
cd web && npm run lint && npm test && npm run build
```

## Validación manual con los datos de prueba

**Prerrequisito**: la aplicación levantada, con la función de datos de prueba habilitada y los datos
cargados (spec 021). El proyecto de prueba es el 25, con 4 horas por Story Point. Abajo, `$API` es la
URL base del backend.

### 1. Métricas del Sprint cerrado

```bash
curl -s $API/projects/25/metricas
```

Esperado: `sprints` tiene solo el Sprint 1; el Sprint 2 (`Activo`) no aparece ni aporta horas. Tanto el
Sprint como `proyecto` informan:

| Campo | Valor |
|---|---|
| `story_points_planificados` | 24 |
| `story_points_completados` | 16 |
| `velocidad` | 16 |
| `porcentaje_historias_completadas` | 75 |
| `horas_estimadas` | 96 |
| `horas_reales` | 66 |
| `desviacion_horas` | -30 |
| `desviacion_porcentaje` | -31.25 |
| `defectos_detectados` y `defectos_resueltos` | 0 (los datos de prueba no tienen defectos) |

### 2. Recalcular al cambiar el parámetro

Enviar a `PUT /projects/25` los datos actuales del proyecto más `"horas_por_story_point": 5`.

Esperado: `horas_estimadas` 120, `desviacion_horas` -54, `desviacion_porcentaje` -45. Los demás valores
no cambian.

### 3. Conservar y borrar el parámetro

- `PUT` sin mencionar el campo: `GET /projects/25` sigue devolviendo 5.
- `PUT` con `"horas_por_story_point": null`: `horas_estimadas`, `desviacion_horas` y
  `desviacion_porcentaje` pasan a `null`, y los demás valores no cambian.

### 4. Valores inválidos

`PUT` con `0`, `-1`, `"abc"` y `""`: cada uno responde 400 con `campo: horas_por_story_point`, y el
proyecto queda como estaba.

### 5. Proyecto sin Sprints cerrados y proyecto inexistente

- Un proyecto recién creado responde 200, con `sprints: []`, valores en cero y `horas_estimadas`,
  `desviacion_horas` y `desviacion_porcentaje` en `null`.
- Con el parámetro definido, solo `desviacion_porcentaje` queda en `null`.
- Un proyecto inexistente responde 404 con "proyecto no encontrado".

### 6. Pantalla de edición

En la edición del proyecto: cargar un valor, guardar y volver a abrir (el valor sigue ahí); vaciar el
campo y guardar (queda sin definir); cargar `0` (se ve el mensaje de error junto al campo).

## Regla de Story Points con Sprint finalizado

Los datos de prueba no tienen una historia sin estimar en un Sprint finalizado, así que User Story 5
se valida con los escenarios de `features/calculo_metricas.feature` y las pruebas de los services.
