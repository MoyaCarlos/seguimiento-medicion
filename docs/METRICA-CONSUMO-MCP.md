# Métrica: consumo de tokens del GitHub MCP Server

Registro del overhead de contexto que implica tener el servidor MCP de
GitHub conectado a Claude Code, medido en la sesión donde se usó por
primera vez para cargar el Product Backlog al repositorio (13 historias
creadas como Issues, 2026-09-14/15).

## Metodología

Claude Code no expone un contador de tokens en tiempo real por servidor,
así que esta es una **estimación razonada** a partir de lo efectivamente
observado en la sesión (tamaño de los bloques de sistema que aparecieron
al conectar el servidor y al cargar cada herramienta), no una medición
instrumentada exacta. Sirve para entender el orden de magnitud del costo,
no como cifra contable precisa.

## Diseño relevante: "deferred tools"

Claude Code no carga el detalle de uso (parámetros, tipos, validaciones)
de las ~48 herramientas que expone el servidor MCP de GitHub apenas se
conecta — solo registra sus **nombres**. El detalle completo de una
herramienta puntual se carga recién cuando hace falta usarla, vía un
buscador interno (`ToolSearch`). Esto evita pagar el costo de las ~45
herramientas que nunca se llegaron a usar en la tarea.

## Desglose de costos observados

| Concepto | Costo aprox. (tokens) | Frecuencia |
|---|---:|---|
| Conexión inicial: nombres de las ~48 herramientas + instrucciones de uso del servidor | ~750 | Una sola vez, al conectar |
| Carga de schema completo (solo de las 3 herramientas usadas: `get_me`, `issue_write`, `list_issue_fields`) | ~600 | Una vez por herramienta, no por llamada |
| Resultado de las llamadas de verificación (`get_me`, `list_issue_fields`) | ~150 | Por llamada |
| **Total infraestructura del MCP** | **~1.500** | |
| Cuerpo de las 13 Issues creadas (Descripción + criterios BDD + metadatos) | ~3.500 | Contenido de negocio — hubiera existido igual, con o sin MCP, ya estaba redactado en `docs/PRODUCT-BACKLOG.md` |

## Comparación con carga completa (sin "deferred tools")

Si el cliente hubiera cargado el schema completo de las 48 herramientas
del servidor al conectar (en vez de solo sus nombres), el costo de
conexión estimado sería de **~5.000–8.000 tokens**, en vez de los ~750
reales. La carga diferida redujo el costo de conexión entre 3 y 5 veces.

## Conclusión

El costo fijo de tener el GitHub MCP conectado ronda los **~1.500 tokens
por sesión** (una sola vez, no se repite por mensaje). Cada llamada
adicional a una herramienta ya cargada solo cuesta lo que devuelve esa
llamada puntual — no hay overhead recurrente por el solo hecho de tener
el servidor conectado y disponible sin usarlo.
