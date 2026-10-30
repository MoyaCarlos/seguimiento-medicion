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

_Se completa al implementar (fecha, código HTTP y observaciones de cada pedido)._
