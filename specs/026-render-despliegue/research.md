# Research: Despliegue continuo en Render

## Esquema oficial

- **Decisión**: validar `render.yaml` contra `https://render.com/schema/render.yaml.json` (el mismo que registra SchemaStore para `**/*render.yaml`), descargado el 2026-10-30.
- **Rationale**: el Render CLI no está instalado en la máquina; el esquema es la fuente oficial y permite validar sin instalar nada en el repo.
- **Alternativas**: Render CLI v2.7.0+ (`render blueprints validate`), no disponible acá; queda como chequeo opcional para quien lo tenga.

## Claves confirmadas en el esquema

| Clave en `render.yaml` | Definición en el esquema | Valor |
|---|---|---|
| `type` | `serverService.type` ∈ web, worker, pserv | `web` |
| `name` | `serverService.name` | `seguimiento-medicion-sprintometro` |
| `runtime` | `runtime` | `docker` |
| `plan` | `serverPlan` (incluye `free`) | `free` |
| `region` | `region` ∈ oregon, ohio, frankfurt, singapore, virginia | `oregon` |
| `branch` | `serverService.branch` | `main` |
| `dockerfilePath` | `serverService.dockerfilePath` | `./Dockerfile` |
| `dockerContext` | `serverService.dockerContext` | `.` |
| `healthCheckPath` | `serverService.healthCheckPath` | `/health` |
| `autoDeployTrigger` | `autoDeployTrigger` ∈ off, commit, checksPass | `checksPass` |
| `envVars[].key/value` | `envVarFromKeyValue` (`value` string o número) | `DEMO_HABILITADO` / `"true"` |
| `buildFilter.ignoredPaths` | `buildFilter` (`paths`, `ignoredPaths`; al menos uno) | las 6 rutas |

`autoDeploy` está deprecado: no se usa.

## Pendiente de confirmar al implementar

- Que el archivo valide sin errores contra el esquema (T002).
- Semántica de los globs de `buildFilter.ignoredPaths`: en los filtros de Render `*` no cruza `/`, por lo que `*.md` solo cubre la raíz y `**/*.md` cubre `.md` de cualquier carpeta (review del PR #200).
- Si el esquema marca `name` como requerido: se incluye igual para identificar el servicio.
