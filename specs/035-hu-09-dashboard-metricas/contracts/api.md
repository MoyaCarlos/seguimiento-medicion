# Contrato que consume la pantalla

Esta historia **no agrega ni cambia endpoints**. El contrato vigente es el de HU-03: [`specs/032-hu-03-calculo-metricas/contracts/openapi.yaml`](../../032-hu-03-calculo-metricas/contracts/openapi.yaml).

## `GET /projects/{id}/metricas`

| Respuesta | Qué hace la pantalla |
|---|---|
| `200` con `sprints` con elementos | Muestra el gráfico y los indicadores |
| `200` con `sprints: []` | Muestra el estado vacío |
| `400` identificador inválido | Muestra el mensaje de la API |
| `404` `proyecto no encontrado` | Muestra el mensaje de la API |
| `500` error interno | Muestra el mensaje de la API |
| Sin conexión o gateway caído | Muestra el mensaje de `ConnectionError` |

Campos en [data-model.md](../data-model.md). La pantalla no recalcula ninguno.

## Contrato de la pantalla

| Elemento | Valor |
|---|---|
| Ruta | `/projects/:id/dashboard` |
| Navegación | Entrada "Dashboard" en las pestañas del proyecto |
| Enlace del estado vacío | `/projects/:id/sprints` |
| Enlace del aviso de horas por Story Point | `/projects/:id/proyecto/edit` |
