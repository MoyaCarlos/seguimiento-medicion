# Quickstart: criterios de aceptación al crear una historia

## Pruebas automáticas

Desde la raíz:

```bash
go build ./...
go vet ./...
gofmt -l .
go test ./...
```

Desde `web/`:

```bash
npm run lint
npm test
npm run build
```

Resultado esperado: todo en verde, incluidos los escenarios nuevos de `features/criterios_aceptacion_creacion.feature` y las pruebas de `web/src/features/backlog/crear.test.tsx`.

## Validación manual con la API

Con el backend levantado (`go run ./cmd/api`) y un proyecto existente (`<id>`):

1. **Con criterios**:

   ```bash
   curl -i -X POST http://localhost:8080/backlog \
     -H 'Content-Type: application/json' \
     -d '{"proyecto_id":<id>,"titulo":"Iniciar sesión","descripcion":"Autenticación","prioridad":"M","criterios_aceptacion":"Dado un usuario, cuando ingresa, entonces accede."}'
   ```

   Esperado: `201` y `criterios_aceptacion` en la respuesta.

2. **Sin criterios** (campo omitido o `null`):

   ```bash
   curl -i -X POST http://localhost:8080/backlog \
     -H 'Content-Type: application/json' \
     -d '{"proyecto_id":<id>,"titulo":"Sin criterios","descripcion":"Descripción","prioridad":"S"}'
   ```

   Esperado: `201` y `criterios_aceptacion: null`.

3. **Solo espacios**:

   Enviar `"criterios_aceptacion":"   "`. Esperado: `201` y `null`.

4. **Límite**:

   Enviar `"criterios_aceptacion"` con 2001 caracteres. Esperado: `400` con `campo: "criterios_aceptacion"` y sin crear la historia.

5. **Listado**:

   `GET /projects/<id>/backlog` debe devolver el valor persistido (o `null`).

## Validación manual de la pantalla

Con `npm run dev` (o el build servido por el backend):

1. Entrar a `Product Backlog` → `Nueva historia`.
2. Verificar que existe el campo opcional **Criterios de aceptación**.
3. Crear una historia con criterios: en "Historia creada" se ven los criterios con formato Markdown.
4. Crear una historia sin criterios: la confirmación muestra "Sin criterios de aceptación".
5. En el listado, la historia creada muestra sus criterios con el mismo formato.
6. Pegar más de 2000 caracteres: el formulario muestra el error bajo el campo y no envía.
7. Forzar un error 400 de `criterios_aceptacion` desde la API (por ejemplo, con un cliente que omita la validación del navegador): el mensaje aparece bajo el campo.

## Fuera de alcance

Notas y depende_de (D17) no se prueban en este incremento. No hay migración ni cambio en `MODELO-DE-DATOS.md`.
