# Quickstart: Generación de Reporte en PDF (HU-10)

## Prerrequisitos

- Go 1.26.5.
- Backend con la base migrada (la 005 ya incluye defectos; no hay migración nueva).
- Datos de un proyecto con al menos un Sprint finalizado para ver todas las secciones.

## Pruebas automatizadas

```bash
go build ./...
go vet ./...
gofmt -l .
go test ./...
go test ./features/...        # escenarios Godog, incluido features/reporte_pdf.feature (@HU-10)
```

## Prueba manual contra la API

1. Levantar el servidor sobre una base nueva:
   ```bash
   go run ./cmd/api
   ```
2. Crear un proyecto, un Sprint finalizado con historias completadas, horas y defectos
   (se puede usar `POST /demo` si está habilitado, o los endpoints de las historias
   anteriores). Anotar el `id` del proyecto y del Sprint.
   Desde HU-21 la descarga pide sesión: ingresar con `POST /auth/ingreso` y guardar el
   token en `TOKEN` (con la demo, `carlos.moya` / `demo1234`).
3. Descargar el reporte del Sprint:
   ```bash
   curl -s -D cabeceras.txt -o sprint.pdf -H "Authorization: Bearer $TOKEN" localhost:8080/sprints/<sprint>/reporte
   ```
   Esperado: `200`, `Content-Type: application/pdf`,
   `Content-Disposition: attachment; filename="reporte-proyecto-<proyecto>-sprint-<numero>.pdf"`,
   y `sprint.pdf` empieza con `%PDF-`.
4. Descargar el reporte del proyecto:
   ```bash
   curl -s -D cabeceras.txt -o proyecto.pdf -H "Authorization: Bearer $TOKEN" localhost:8080/projects/<proyecto>/reporte
   ```
   Esperado: `200` y un PDF con datos del proyecto, integrantes, Product Backlog, Sprints
   finalizados y Sprint Activo si existe.
5. Casos de error:
   - Sprint `Pendiente` o `Activo`: `400` con `campo: estado` y sin archivo.
   - Sprint inexistente: `404` "sprint no encontrado".
   - Proyecto inexistente: `404` "proyecto no encontrado".
   - Id no numérico o `0`: `400` con `campo: id`.
6. Proyecto sin Sprints finalizados (con o sin Sprint Activo): el PDF se genera y muestra
   métricas, esfuerzo y defectos como "sin datos".
7. Abrir los PDF en un visor: verificar tildes y eñes, que las secciones estén completas y
   que el documento pagine si hay muchas historias.
8. Navegar directo a la URL de descarga desde el navegador (con `Accept: text/html`): desde
   HU-21 responde `401` `{"mensaje":"Tu sesión no es válida"}`, porque el navegador no manda
   el token (nunca devuelve `index.html`). Con la cabecera `Authorization` baja el PDF. Por
   eso la pantalla descarga el reporte con el cliente de la API (`fetch` + blob con
   `Authorization`), no con un enlace común (issue #78).
9. Reiniciar el servidor y repetir: el reporte se recalcula y no modifica la base.

## Resultado de la validación manual

Validado el 2026-11-06 contra el servidor local (`DB_PATH` en un archivo temporal,
`DEMO_HABILITADO=true`, `WEB_DIR` inexistente) y los datos de `POST /demo/datos`
(proyecto 25, con el Sprint 1 finalizado y el Sprint 2 activo):

- `GET /sprints/1/reporte` → `200`, `Content-Type: application/pdf`,
  `Content-Disposition: attachment; filename="reporte-proyecto-25-sprint-1.pdf"`,
  archivo de 17.316 bytes que empieza con `%PDF-`.
- `GET /projects/25/reporte` → `200`, `Content-Type: application/pdf`,
  `Content-Disposition: attachment; filename="reporte-proyecto-25.pdf"`,
  archivo de 39.618 bytes que empieza con `%PDF-`.
- Errores: Sprint activo → `400` con `campo: estado`; Sprint inexistente → `404`
  "sprint no encontrado"; proyecto inexistente → `404` "proyecto no encontrado";
  id `abc` o `0` → `400` con `campo: id`. En todos los casos responde JSON y no
  entrega archivo.
- Proyecto recién creado sin Sprints finalizados (id 26) → `200` y PDF de 6.871
  bytes que empieza con `%PDF-` (secciones de métricas, esfuerzo y defectos como
  "sin datos").
- La descarga directa con `Accept: text/html` se cubre con el test de
  `internal/http/spa_test.go` (`TestServirSPA_ReportesLleganALaAPIaunqueAceptenHTML`),
  porque en esta corrida manual `WEB_DIR` no existía y la SPA no estaba montada.
- Los reportes son de solo lectura: no modifican la base.
