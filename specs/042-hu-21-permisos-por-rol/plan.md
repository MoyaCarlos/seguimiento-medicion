# Implementation Plan: Permisos por rol (HU-21)

**Branch**: `feature/HU-21-permisos-por-rol` | **Date**: 2026-11-11 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/042-hu-21-permisos-por-rol/spec.md`

## Summary

Imponer la identidad de HU-20 en toda la API: un envoltorio HTTP `conSesion` resuelve la
sesión (401) y entrega el `usuarioID` a los handlers; los **servicios** comprueban la
pertenencia y el rol por proyecto (403), después de cargar el recurso (404) y antes de
validar (400/409). Se quitan `usuario_id`, `integrante_id` y `creador` de las peticiones;
`GET /projects` pasa a ser «Mis proyectos» con el rol; se suma `GET /users`; configurar el
equipo solo agrega usuarios registrados como Product Builder. El frontend quita los
selectores de identidad, deshabilita las acciones de gestión para el Product Builder y
muestra «No sos integrante» / «No tenés permiso». Sin cambios de esquema. Detalle de
decisiones en [research.md](research.md).

## Technical Context

**Language/Version**: Go (versión del `go.mod`); TypeScript + React (Vite) en `web/`

**Primary Dependencies**: `net/http` estándar (ServeMux con patrones), `modernc.org/sqlite`, Godog; React, Vitest, Testing Library. Sin dependencias nuevas.

**Storage**: SQLite; **sin migración** (lee `usuarios`, `integrantes`, `sesiones`; ver [data-model.md](data-model.md))

**Testing**: `go test` (TDD en `internal/service`, `internal/http`, `internal/repository`), Godog en `features/`, Vitest en `web/`

**Target Platform**: servidor Linux en contenedor (Render) y navegadores actuales

**Project Type**: web-service + SPA

**Performance Goals**: sin metas nuevas; la comprobación de rol agrega a lo sumo una lectura de `integrantes` por petición (ya la hace la votación)

**Constraints**: entrega 16-nov; plan B del backlog si el 13-nov el backend no está completo; check `backend` y `frontend` en verde en el mismo PR

**Scale/Scope**: ~45 rutas, ~15 servicios con actor nuevo, ~20 `.feature` a revisar, ~8 pantallas

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principio | Cumplimiento | Estado |
|---|---|---|
| I. Test-First | Cada regla (sesión, pertenencia, rol, lista filtrada, usuario registrado, orden de errores) entra con RED → GREEN → REFACTOR en commits separados; fakes en memoria para los servicios. | ✅ |
| II. SDD | specify → clarify → plan → tasks → analyze antes de implementar; una historia a la vez. | ✅ |
| III. BDD | `features/permisos_por_rol.feature` con un escenario por criterio (normal, alternativo, límite, error); los `.feature` existentes se ajustan en el mismo PR; frontend con pruebas de componente con el escenario en el nombre. | ✅ |
| IV. Clean/Hexagonal | La sesión se resuelve en `internal/http`; los servicios reciben `usuarioID` y aplican las reglas; el dominio no conoce tokens ni HTTP. Repositorios con interfaces y fakes. | ✅ |
| V. Clean Code / KISS / YAGNI / DRY | Reglas de pertenencia/rol en un solo lugar (`internal/service/integrantes.go`); sin middleware por rol ni capa de permisos genérica; se elimina `resolverUsuario` sin uso. | ✅ |
| VI. Stack | Solo biblioteca estándar y el stack existente. | ✅ |
| Modelo de datos | Sin cambios de esquema; solo texto en `MODELO-DE-DATOS.md` (§4.2, §8). Cualquier necesidad de esquema se frena y se consulta. | ✅ |

**Re-check post-diseño**: sin violaciones; Complexity Tracking vacío.

## Project Structure

### Documentation (this feature)

```text
specs/042-hu-21-permisos-por-rol/
├── spec.md
├── plan.md              # este archivo
├── research.md          # decisiones R1–R11
├── data-model.md        # sin esquema nuevo; lecturas nuevas
├── quickstart.md        # validación automática y manual
├── contracts/
│   └── api-permisos.md
├── checklists/
│   └── requirements.md
└── tasks.md             # /speckit.tasks
```

### Source Code (repository root)

```text
cmd/api/main.go                       # wiring; delega el registro de rutas
internal/
├── domain/errors.go                  # mensajes de ErrNoEsIntegrante / ErrNoEsScrumMaster; ErrUsuarioNoEncontrado
├── service/
│   ├── integrantes.go                # exigirIntegrante; exigirScrumMaster (no integrante → ErrNoEsIntegrante)
│   ├── crear_proyecto.go             # CreadorID desde la sesión
│   ├── listar_proyectos.go           # por integrante, con rol
│   ├── listar_usuarios.go            # NUEVO: GET /users
│   ├── asignar_integrante.go         # solo registrados; siempre Product Builder
│   ├── usuarios.go                   # se elimina resolverUsuario
│   ├── obtener_proyecto.go, editar_proyecto.go
│   ├── crear_sprint.go, iniciar_sprint.go, cerrar_sprint.go
│   ├── crear_historia.go, editar_historia.go, cambiar_estado_historia.go, asignar_historia_sprint.go
│   ├── imputar_horas.go, registrar_defecto.go, resolver_defecto.go
│   └── votar.go, consultar_votacion.go (+ los de votación que ya exigen Scrum Master)
├── repository/
│   ├── project_repository.go (+ sqlite, en memoria)   # ListarPorIntegrante
│   └── user_repository.go (+ sqlite)                  # ListarRegistrados
└── http/
    ├── sesion_middleware.go          # NUEVO: conSesion, usuarioDeSesion
    ├── rutas.go                      # NUEVO: RegistrarRutas (sale de main.go)
    ├── rutas_permisos_test.go        # NUEVO: recorre la tabla «Qué se protege»
    ├── usuarios_handler.go           # NUEVO: GET /users
    ├── project_dto.go, votacion_dto.go, horas_dto.go   # sin campos de identidad; + rol
    └── *_handler.go                  # pasan usuarioDeSesion(r) a los servicios
features/
├── permisos_por_rol.feature          # NUEVO (@HU-21)
├── steps_permisos.go                 # NUEVO: «Dado que inicié sesión como "X"» + llamadas con token
└── *.feature / steps_*.go            # ajustes: sin usuario_id / integrante_id / creador
web/src/
├── compartido/api/client.ts          # traducción del 403
├── compartido/proyectos/ListaProyectos.tsx   # rol y estado vacío
├── compartido/permisos/              # NUEVO: NoIntegrante, AvisoSinPermiso, useEsScrumMaster
├── features/votacion/                # quitar SelectorActor; AccionesScrumMaster deshabilitadas
├── features/sprints/                 # FormularioHoras sin integrante; crear/iniciar/cerrar deshabilitados
└── features/proyecto/                # crear sin creador; Editar y Configuracion (GET /users, sin rol)
MODELO-DE-DATOS.md                    # solo texto: §4.2 y §8
```

**Structure Decision**: estructura web existente (backend Go hexagonal + `web/`); se agregan
tres archivos en `internal/http`, un servicio, una feature de Godog y una carpeta
`compartido/permisos` en el frontend.

## Fases sugeridas para `/speckit.tasks`

1. **Base**: mensajes de error, `conSesion`/`usuarioDeSesion`, `RegistrarRutas` y el test
   de rutas (US3, escenario 3).
2. **Identidad** (US1, escenario 7): votar, consultar votación, horas, crear proyecto.
3. **Pertenencia** (US4, escenario 4): `exigirIntegrante` en escrituras, panel y votación.
4. **Scrum Master** (US2, escenarios 1 y 2): Sprints, votación, editar proyecto, equipo.
5. **Mis proyectos y usuarios** (US5, US6; escenarios 5, 6, 8 y 9).
6. **Godog**: `permisos_por_rol.feature` y ajuste de los `.feature` existentes.
7. **Frontend** (US7) y documentación (`MODELO-DE-DATOS.md`, límite conocido de lectura).

Las fases 1 a 6 son las que pide el plan B del backlog antes del 13-nov.

## Complexity Tracking

Sin violaciones de la constitución.
