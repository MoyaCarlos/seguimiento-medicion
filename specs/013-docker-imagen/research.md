# Research: Imagen Docker

## Imágenes base (tags verificados con `docker manifest inspect` el 2026-10-24)

| Uso | Imagen | Existe |
|---|---|---|
| Build frontend | `node:24-alpine` | sí |
| Build backend | `golang:1.26.5-alpine` (`go.mod` pide 1.26.5) | sí |
| Final | `gcr.io/distroless/static-debian13:nonroot` | sí (también `-debian12`) |

- **Decisión**: alpine en los stages de build (más chico, `npm ci` y `go build` no necesitan glibc: sin CGO). Final distroless `static` `:nonroot` (UID 65532).
- **Alternativa descartada**: `scratch` + usuario a mano (más pasos, sin CA certs ni tzdata); `golang` sin alpine (más pesado, sin beneficio).

## Una imagen vs. dos
- **Decisión** (Carlos): una sola; Go sirve `web/dist`. Evita CORS y un segundo servidor.

## Regla de servido del SPA
- Fuente: `web/vite.config.ts` → `bypass`: si `Accept` incluye `text/html` devuelve `/index.html`, si no va a la API; solo para `/backlog`, `/projects`, `/sprints`.
- El servidor de desarrollo además hace fallback de SPA para el resto (ej. `/`); se replica con FR-005b.
- **Alternativa descartada**: `http.FileServer` con fallback por 404 (rompería los 404 JSON de la API y el 405 de métodos).

## `/data` sin shell
- La imagen final no tiene `mkdir`/`chown`. Un volumen con nombre vacío copia contenido y dueño del directorio de la imagen; por eso `/data` se crea en un stage de build y se copia con `--chown=65532:65532` (nonroot).
- Bind mount del host: el dueño lo define el host (documentado en README).

## Sin HEALTHCHECK
- Distroless no trae shell ni curl; se valida con el smoke test de `quickstart.md`.
