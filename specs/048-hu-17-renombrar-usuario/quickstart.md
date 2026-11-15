# Quickstart: validar HU-17

## Prerrequisitos

- Go 1.26 y Node con las dependencias de `web/` instaladas.
- Base SQLite local (no hay migración nueva; las migraciones existentes se aplican solas).
- Una cuenta registrada (`POST /auth/registro`) y su token.

## Verificación automática

```sh
go test ./...                    # dominio, servicios, repositorios, HTTP y Godog
cd web && npm test               # Vitest
cd web && npm run lint && npm run build
```

Esperado: todo en verde, incluidos `features/gestion_usuarios.feature` y los escenarios
previos sin cambios.

## Verificación manual de la API

Con la API en `:8080` (`go run ./cmd/api`):

1. Registrar `ana` con una contraseña válida y guardar el token.
2. `PUT /auth/yo` con `Authorization: Bearer <token>` y `{"nombre":"ana.perez"}` → `200`
   con `{"id":...,"nombre":"ana.perez"}`.
3. `GET /auth/yo` con el mismo token → `200` con `ana.perez`.
4. `POST /auth/ingreso` con `ana.perez` y la misma contraseña → `200`; con `ana` → `401`
   «Usuario o contraseña incorrectos».
5. Registrar `bruno`; con el token de `ana`, `PUT /auth/yo` con `{"nombre":"BRUNO"}` → `409`
   «Ese nombre de usuario ya existe»; `GET /auth/yo` sigue devolviendo `ana.perez`.
6. `PUT /auth/yo` con un nombre de 2 caracteres o con `@` → `400` con `campo: "nombre"`.
7. `PUT /auth/yo` sin token → `401`.

## Verificación manual de la interfaz

1. Iniciar sesión y abrir el menú de la cuenta; elegir «Cambiar nombre».
2. Cambiar a un nombre válido y guardar: el menú muestra el nombre nuevo sin recargar.
3. Intentar un nombre ya usado: se muestra «Ese nombre de usuario ya existe» y el nombre
   anterior sigue en el menú.
4. Cancelar: el formulario se cierra sin llamar a la API.
5. Recargar la página: `GET /auth/yo` mantiene el nombre nuevo y la sesión sigue vigente.

## Criterios de salida

- Los escenarios 5 y 6 de HU-17 (renombrado y nombre en uso) automatizados en Godog y en
  pruebas de componente.
- Sin migración ni cambio de esquema; el historial y las sesiones siguen asociados por
  `usuario_id`.
