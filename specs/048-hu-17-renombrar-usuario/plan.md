# Implementation Plan: Renombrar el propio usuario (HU-17)

**Branch**: `feature/HU-17-renombrar-usuario` | **Date**: 2026-11-15 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/048-hu-17-renombrar-usuario/spec.md`

## Summary

Agregar `PUT /auth/yo` para que una persona autenticada cambie solo su propio nombre de
usuario. El handler toma el usuario de la sesión, valida el nombre con las reglas de HU-20 y
llama a un caso de uso que actualiza `nombre` y `nombre_normalizado` en una sola sentencia.
La unicidad la garantiza `usuarios.nombre_normalizado UNIQUE`; el choque se traduce a
`ErrUsuarioEnUso` (409). No hay migración: todas las tablas referencian `usuario_id`, así que
el historial y las sesiones no cambian. El frontend suma «Cambiar nombre» al menú de la
cuenta, con un formulario que actualiza el usuario del contexto de sesión.

## Technical Context

**Language/Version**: Go 1.26 (backend); TypeScript 6 + React 19, Vite 8 (frontend).

**Primary Dependencies**: ninguna nueva. Se reutilizan `domain.ValidarNombreUsuario`,
`domain.NormalizarNombre`, el cliente HTTP compartido y el contexto de sesión de HU-20.

**Storage**: SQLite (`modernc.org/sqlite`). **Sin migración ni cambio de esquema.**
`usuarios` ya tiene `nombre` y `nombre_normalizado UNIQUE` (001/HU-20).

**Testing**: `go test` (dominio/servicio con fakes, SQLite real en el repositorio, handler);
Godog para la API; Vitest + Testing Library para el menú y el contexto.

**Target Platform**: servidor Linux/Docker + navegador.

**Project Type**: web (API REST + SPA).

**Performance Goals**: una sola sentencia `UPDATE` por renombrado; no hay meta de tiempo.

**Constraints**: el usuario sale de la sesión; no se acepta `usuario_id` en el cuerpo; el
nombre se valida igual que en el registro; no se modifican otras tablas.

**Scale/Scope**: 1 ruta, 1 caso de uso, 1 método de repositorio, 1 formulario en el menú y
sus pruebas.

## Constitution Check

| Principio | Estado |
|---|---|
| I. Test-First | RED/GREEN/REFACTOR por comportamiento; commits por fase en [tasks.md](tasks.md) |
| II. SDD | Este plan y la spec se escriben antes de implementar |
| III. BDD | `features/gestion_usuarios.feature` con los escenarios de renombrado; el frontend con `Given-When-Then` en el nombre |
| IV. Clean/Hexagonal | Validación en el dominio, puerto en `repository`, caso de uso en `service`, HTTP solo traduce |
| V. SOLID/KISS/YAGNI | Un puerto de un método, un caso de uso, sin acción de administración de otros usuarios |
| VI. Stack | Sin dependencias nuevas |

Sin violaciones. No hay migración; solo se documenta la regla en `MODELO-DE-DATOS.md` §5 y §8.

## Project Structure

### Documentation (this feature)

```text
specs/048-hu-17-renombrar-usuario/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── api-usuarios.md
└── tasks.md
```

### Source Code (repository root)

```text
internal/
├── domain/
│   └── usuario_credenciales.go          # se reutiliza ValidarNombreUsuario/NormalizarNombre
├── repository/
│   ├── credencial_repository.go         # + ActualizadorNombreUsuario
│   ├── sqlite_credencial.go             # + CambiarNombre (UPDATE atómico)
│   └── sqlite_credencial_test.go        # + casos de éxito, duplicado e inexistente
├── service/
│   ├── cambiar_nombre_usuario.go        # caso de uso
│   ├── cambiar_nombre_usuario_test.go
│   └── fakes_test.go                    # + fake del puerto
├── http/
│   ├── auth_dto.go                      # + cambiarNombreRequest
│   ├── auth_handler.go                  # + CambiarNombre; constructor con el caso de uso
│   ├── auth_handler_test.go             # + 200/400/401/409
│   ├── rutas.go                         # + PUT /auth/yo (SoloSesion)
│   └── sprint_handler.go                # escribirError ya cubre 400/401/409/404
├── app/
│   ├── app.go                           # cableado
│   └── rutas_permisos_test.go           # + nivel de PUT /auth/yo
features/
├── gestion_usuarios.feature             # @HU-17
├── steps_gestion_usuarios.go
└── features_test.go                     # registrar los pasos
web/src/
├── compartido/api/client.ts             # ya tiene put
├── features/autenticacion/
│   ├── api.ts                           # + cambiarNombre
│   ├── SesionContext.tsx                # + cambiarNombre
│   ├── MenuUsuario.tsx                  # + formulario «Cambiar nombre»
│   ├── MenuUsuario.test.tsx             # + escenarios
│   ├── SesionContext.test.tsx           # + actualización del usuario
│   └── estilos.css                      # + estilos del formulario
MODELO-DE-DATOS.md                       # §4.2, §5 y §8: regla y mapa de HU-17 (sin migración)
docs/PRODUCT-BACKLOG.md                  # HU-17 redefinida
```

**Structure Decision**: se respeta la arquitectura existente. El renombrado es una operación
de identidad del usuario, así que el puerto vive junto a los de credenciales y el caso de uso
en `internal/service`; el handler de autenticación es el adaptador natural porque la ruta es
`/auth/yo`. El frontend reutiliza el contexto de sesión y el menú de cuenta de HU-20.

## Decisiones técnicas

| # | Decisión | Motivo |
|---|---|---|
| 1 | `PUT /auth/yo` en lugar de `PATCH` | Es un reemplazo del nombre y el cliente compartido ya expone `put`; consistente con los `PUT` de edición del resto de la API. |
| 2 | El usuario sale de la sesión; el cuerpo solo lleva `nombre` | Evita que alguien cambie el nombre de otra persona (opción A). |
| 3 | `UPDATE usuarios SET nombre = ?, nombre_normalizado = ? WHERE id = ?` | Atómico y cubre la carrera de dos renombrados con el índice único. |
| 4 | `SQLITE_CONSTRAINT_UNIQUE` → `ErrUsuarioDuplicado` → `ErrNombreUsuarioEnUso` | Reutiliza el mapeo y el mensaje del registro. |
| 5 | Sin migración | `usuarios` ya tiene las columnas y las demás tablas apuntan por `usuario_id`. |
| 6 | Actualizar el usuario del contexto con la respuesta | El menú muestra el nombre nuevo sin recargar; la sesión sigue vigente. |

## Complexity Tracking

Sin violaciones de la constitución que justificar.
