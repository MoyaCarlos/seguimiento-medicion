# Quickstart: validar datos de prueba v2

Prerrequisito: `DEMO_HABILITADO=true`, base vacía.

1. Pruebas del backend: `go test ./internal/demo/... ./internal/http/... ./features/...`.
2. Pruebas del frontend: `cd web && npm test -- compartido/demo features/autenticacion`.
3. Levantar la API: `DEMO_HABILITADO=true go run ./cmd/api`; en otra terminal `curl -X POST localhost:8080/demo/datos` (puerto según la configuración). Esperado: 201 con el resumen de [contracts/api.md](contracts/api.md).
4. Repetir el `POST`: 409 `ya_cargado`.
5. Iniciar sesión: `POST /auth/ingreso` con `carlos.moya` / `demo1234` (y los otros tres). Esperado: sesión.
6. Con la base vacía, registrar a mano `ana.iriarte` y cargar: 409 `ya existe el usuario ana.iriarte` y las tablas de la carga vacías.
7. Abrir `/ingresar` en el navegador: aparece el botón; tras cargar, la tabla de credenciales.
8. Abrir el proyecto 25: 2 Sprints finalizados con velocidad 13 y 8, uno activo, uno pendiente; 5 defectos.

Detalles de datos en [data-model.md](data-model.md).
