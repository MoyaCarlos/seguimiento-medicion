# Data Model: Dashboard de Métricas (HU-09)

Esta historia **no agrega ni cambia tablas, columnas ni migraciones**. El cálculo y sus fuentes están en [`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md) §6 y en la [spec de HU-03](../032-hu-03-calculo-metricas/spec.md).

## Tipos que consume la pantalla

Reflejan la respuesta de `GET /projects/{id}/metricas` ([contrato](./contracts/api.md)). `null` significa "sin datos".

**`Metricas`**

| Campo | Tipo | Uso en el panel |
|---|---|---|
| `story_points_planificados` | number | Serie "Comprometidos" del gráfico |
| `story_points_completados` | number | Serie "Completados" del gráfico |
| `velocidad` | number | Indicador de velocidad (promedio por Sprint en el proyecto) |
| `porcentaje_historias_completadas` | number | Indicador con proporción |
| `horas_estimadas` | number \| null | Barras de horas comparadas |
| `horas_reales` | number | Barras de horas comparadas |
| `desviacion_horas` | number \| null | Indicador de desviación |
| `desviacion_porcentaje` | number \| null | Indicador de desviación |
| `defectos_detectados` | number | Barras de defectos comparadas |
| `defectos_resueltos` | number | Barras de defectos comparadas |

**`MetricasProyecto`**: `proyecto_id`, `horas_por_story_point` (number \| null), `proyecto` (`Metricas`) y `sprints`, una lista de `{ sprint_id, numero, metricas }` ordenada por número.

## Estados de la pantalla

| Estado | Condición | Se muestra |
|---|---|---|
| Cargando | Consulta en curso | Indicador de carga |
| Con datos | `sprints` tiene al menos un elemento | Gráfico e indicadores |
| Vacío | `sprints` está vacío (no se miran los valores en cero) | Estado vacío con acceso a Sprints |
| Error | La consulta falla | Mensaje de error, nada más |

## Funciones puras (`metricas.ts`)

| Función | Resultado |
|---|---|
| `datosDelGrafico(sprints)` | Una fila por Sprint: número, comprometidos, completados |
| `formatoNumero(valor)` | Número con coma decimal, o "sin datos" si es `null` |
| `textoDesviacion(horas)` | "más", "menos" o "sin desviación" respecto de lo estimado; `null` si no hay datos |
| `proporcion(valor, maximo)` | Porcentaje de ancho entre 0 y 100; 0 si el máximo es 0 o el valor es `null` |
