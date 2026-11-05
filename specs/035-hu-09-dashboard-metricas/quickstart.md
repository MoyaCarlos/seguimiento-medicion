# Quickstart: validación del Dashboard (HU-09)

## Pruebas automáticas

Desde `web/`:

```bash
npm test
npm run lint
npm run build
```

Resultado esperado: todo en verde. Las pruebas de la feature están en `src/features/dashboard/` y tienen el escenario en su nombre.

Los comportamientos de la API se validan con `go test ./...` (`features/calculo_metricas.feature`).

## Escenarios BDD (#70)

No se agrega un archivo `.feature`: la historia no cambia la API, y Godog no ejercita la interfaz. Los escenarios de la spec se automatizan como pruebas de componente con el escenario en su nombre (Principio III), y los de la API ya existen.

**Pruebas de componente** (`web/src/features/dashboard/`, salvo la navegación):

| Categoría | Escenarios | Archivo de prueba |
|---|---|---|
| Normal | US1.1, US1.2, US1.3 | `GraficoVelocidad.test.tsx` |
| Normal | Navegación y ruta de US1.1 | `web/src/rutas.test.tsx` |
| Normal | US2.1 a US2.5 | `Indicadores.test.tsx` |
| Alternativo | US2.6 (sin horas por Story Point) | `Indicadores.test.tsx` |
| Alternativo | US3.1, US3.2, US3.3 (estado vacío) | `index.test.tsx` |
| Límite | US1.4, US1.5 | `GraficoVelocidad.test.tsx` |
| Límite | Desviación porcentual sin datos; horas y defectos en cero | `Indicadores.test.tsx` |
| Límite | Consulta única; volver y ver valores nuevos | `index.test.tsx` |
| Error | US4.1 (carga), US4.2, US4.3 y error interno | `index.test.tsx` |

**API, ya cubierta con Godog** (`features/calculo_metricas.feature`):

| Comportamiento que usa el Dashboard | Escenario |
|---|---|
| Valores por Sprint | Story Points, velocidad y porcentaje de un Sprint finalizado; Horas y desviación de un Sprint finalizado; Defectos detectados y resueltos en un Sprint |
| Valores del proyecto | Valores del proyecto con dos Sprints finalizados |
| Solo Sprints finalizados | Un Sprint activo no aporta a las métricas |
| Estado vacío | Proyecto sin Sprints finalizados y con horas por Story Point; Proyecto sin Sprints finalizados y sin horas por Story Point |
| "Sin datos" | Métricas de un Sprint finalizado sin horas por Story Point |
| Error | Consultar las métricas de un proyecto inexistente |

## Validación manual

Con el backend y el frontend levantados (ver `README.md`):

| # | Preparación | Acción | Resultado esperado |
|---|---|---|---|
| 1 | Proyecto sin Sprints | Entrar a la pestaña "Dashboard" | Estado vacío con acceso a Sprints; sin gráfico ni indicadores |
| 2 | Proyecto cuyo único Sprint está activo | Entrar al Dashboard | El mismo estado vacío |
| 3 | Cerrar ese Sprint con una historia completada y una sin completar | Volver al Dashboard | Barras "Comprometidos" y "Completados" con los valores del Sprint, e indicadores |
| 4 | Proyecto sin horas por Story Point | Ver los indicadores | Horas estimadas y desviación en "sin datos"; aviso con enlace a la edición del proyecto |
| 5 | Definir las horas por Story Point | Volver al Dashboard | Horas estimadas y desviación con valor |
| 6 | Comparar con `GET /projects/{id}/metricas` | Abrir "Ver valores por Sprint" | Todos los valores coinciden |
| 7 | Ninguna | Entrar a `/projects/999999/dashboard` | Mensaje "proyecto no encontrado", sin estado vacío |
| 8 | Backend detenido | Entrar al Dashboard | Mensaje de conexión, sin estado vacío |
| 9 | Ventana de 320 px, proyecto con 12 Sprints finalizados | Recorrer el panel y las cinco pestañas | Sin desplazamiento horizontal de la página: las pestañas pasan a otra línea y el gráfico se desplaza dentro de su contenedor, con los rótulos legibles |
| 10 | Pestaña de red del navegador, recarga en otra sección | Observar las descargas y luego entrar al Dashboard | El código de Recharts se descarga solo al entrar al Dashboard |
| 11 | Proyecto con varios Sprints finalizados; alguien del equipo que no haya visto la pantalla | Pedirle que diga en qué Sprint se completó menos de lo comprometido y si se trabajó más o menos de lo estimado | Responde sin ayuda en menos de 30 segundos (SC-002) |
