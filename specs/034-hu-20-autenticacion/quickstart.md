# Quickstart: validar HU-20

Guía de validación; no repite el código ni los contratos
([api-auth](contracts/api-auth.md), [frontend-sesion](contracts/frontend-sesion.md),
[data-model](data-model.md)).

## Prerrequisitos

- Go 1.26 y Node con las dependencias de `web/` instaladas.
- Base SQLite local (la migración 006 se aplica sola al arrancar la API).

## Orden de trabajo (Test-First, un commit por fase)

1. Migración 006 y test de esquema (columnas, índice único, base con usuarios previos).
2. Dominio: `ValidarNombreUsuario`, `ValidarContrasena`, `Sesion.Vigente` (tablas de casos).
3. Repositorios SQLite: credenciales y sesiones contra SQLite real.
4. Servicios: contraseñas, emisor de sesiones, registrar, ingresar, obtener sesión, cerrar.
5. HTTP + cableado en `main.go`, `PrefijosAPI` y proxy de Vite; escenarios de Godog.
6. Frontend: enganche del cliente, almacenamiento, contexto, pantallas, menú, ruta protegida.
7. `MODELO-DE-DATOS.md` junto con la migración (mismo PR).

## Verificación automática

```sh
go test ./...                    # dominio, servicios, repositorios, HTTP y Godog
cd web && npm test               # Vitest
cd web && npm run lint && npm run build
```

Esperado: todo en verde, incluidos los escenarios de Godog anteriores sin modificarlos
(SC-007) y `TestPrefijosAPI_IncluyenLosQueProxeaVite`.

## Verificación manual de la API

Con la API en `:8080` (`go run ./cmd/api`):

1. `POST /auth/registro` con un nombre y una contraseña de 8 o más caracteres → 201 con
   `token`. Repetirlo con el nombre en otras mayúsculas → 409.
2. `GET /auth/yo` con `Authorization: Bearer <token>` → 200 con id y nombre.
3. `POST /auth/salir` → 204; repetir el `GET /auth/yo` → 401 con «Tu sesión no es válida».
4. `POST /auth/ingreso` con contraseña incorrecta y con un usuario inexistente → 401 con el
   mismo cuerpo.
5. Un endpoint viejo (`GET /projects`) sin token → 200 como antes.

## Verificación manual de la interfaz

1. Abrir una ruta protegida sin sesión → lleva a `/ingresar`.
2. Registrarse → queda dentro y vuelve a la ruta que pidió; el menú muestra el nombre.
3. Borrar o alterar el token en `localStorage` y navegar → aviso «Tu sesión venció» y vuelta a
   `/ingresar`; al ingresar, regresa a la pantalla en la que estaba.
4. «Cerrar sesión» → lleva a `/ingresar` y el token anterior ya no sirve.

## Criterios de salida

- Los 7 escenarios del issue #233 automatizados (Godog y Vitest).
- Migración aplicada sobre una base con usuarios previos sin error.
- Revisión de logs y respuestas sin contraseña ni token en claro (SC-006).
