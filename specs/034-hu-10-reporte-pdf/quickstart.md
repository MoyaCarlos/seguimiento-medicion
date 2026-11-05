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
3. Descargar el reporte del Sprint:
   ```bash
   curl -s -D cabeceras.txt -o sprint.pdf localhost:8080/sprints/<sprint>/reporte
   ```
   Esperado: `200`, `Content-Type: application/pdf`,
   `Content-Disposition: attachment; filename="reporte-proyecto-<proyecto>-sprint-<numero>.pdf"`,
   y `sprint.pdf` empieza con `%PDF-`.
4. Descargar el reporte del proyecto:
   ```bash
   curl -s -D cabeceras.txt -o proyecto.pdf localhost:8080/projects/<proyecto>/reporte
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
8. Navegar directo a la URL de descarga desde el navegador (con `Accept: text/html`): debe
   bajar el PDF, no devolver `index.html`.
9. Reiniciar el servidor y repetir: el reporte se recalcula y no modifica la base.

## Resultado de la validación manual

Pendiente de completar al terminar la implementación.
