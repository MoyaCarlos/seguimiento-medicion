# Quickstart: Despliegue continuo en Render

Verificación manual contra el servicio real. Solo pedidos de lectura (`GET`); no se toca el dashboard ni se modifica el servicio.

URL: https://seguimiento-medicion-sprintometro.onrender.com

## 1. Validar `render.yaml`

Contra el esquema oficial `https://render.com/schema/render.yaml.json` (JSON Schema 2020-12), con herramientas temporales fuera del repo. Resultado esperado: sin errores.

## 2. Despertar el servicio

El plan Free se duerme tras 15 minutos sin tráfico; el primer pedido puede tardar decenas de segundos.

```sh
curl -s -o /dev/null -w '%{http_code}\n' --max-time 120 https://seguimiento-medicion-sprintometro.onrender.com/health
```

## 3. Pedidos de verificación

```sh
curl -si https://seguimiento-medicion-sprintometro.onrender.com/health
curl -s -H 'Accept: text/html' https://seguimiento-medicion-sprintometro.onrender.com/ | head -c 300
curl -si https://seguimiento-medicion-sprintometro.onrender.com/projects
```

Esperado:

1. `/health` → `200` y cuerpo `ok`.
2. `/` con `Accept: text/html` → el `index.html` de la aplicación; en un navegador, la aplicación abre.
3. `/projects` → `200` con JSON.

## Resultado

Verificado el 2026-10-30 (11:36, hora local) contra el servicio real, solo con `GET`; no se tocó el dashboard.

| Chequeo | Resultado |
|---|---|
| Validación de `render.yaml` contra `https://render.com/schema/render.yaml.json` | 0 errores. Control negativo: agregar una clave inventada produce errores, así que el validador no es permisivo. El Render CLI no está instalado, no se usó |
| Despertar el servicio (`GET /health`) | `200` en 12,6 s (arranque en frío) |
| `GET /health` | `200`, cuerpo `ok` |
| `GET /` con `Accept: text/html` | `200`, `content-type: text/html`, devuelve el `index.html` (`<title>Métrica · Seguimiento de proyectos</title>`, `<div id="root">`, bundle `/assets/index-*.js`) |
| `GET /projects` | `200`, `content-type: application/json`, cuerpo `[]` (base vacía, coherente con que se vacía en cada deploy) |

No verificado: abrir la raíz en un navegador real (se comprobó que el servidor entrega el `index.html`, no que la interfaz renderice); y que `ignoredPaths` de `render.yaml` use la misma semántica de globs que "Ignored Paths" del dashboard, porque comprobarlo exige tocar el dashboard o hacer merges de prueba.
