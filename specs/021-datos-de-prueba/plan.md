# Implementation Plan: Cargar datos de prueba (issue #144, versión 1)

**Branch**: `feature/144-datos-de-prueba` | **Date**: 2026-10-30 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/021-datos-de-prueba/spec.md`

## Summary

Un paquete nuevo `internal/demo` carga el proyecto de demostración (id 25) en **una sola transacción**, armada con bloques independientes (usuarios, proyecto y equipo, Sprints, historias, horas, votación). Cada bloque usa los constructores del dominio y escribe con SQL dentro de la transacción. La idempotencia sale de la clave primaria del proyecto 25. Un handler HTTP expone `GET /demo/estado` y `POST /demo/datos` detrás de un único punto de habilitación (`DEMO_HABILITADO`). El frontend suma un componente propio con el botón. Sin cambio de esquema ni de `MODELO-DE-DATOS.md`.

## Technical Context

**Language/Version**: Go 1.x del `go.mod`; TypeScript + React (Vite) en `web/`

**Primary Dependencies**: las existentes (`modernc.org/sqlite`, Godog, Vitest, Testing Library). Ninguna nueva.

**Storage**: SQLite embebido, esquema de las migraciones 001–004 sin tocar.

**Testing**: `testing` de Go (bloques, integración con migraciones reales, idempotencia, fallo inyectado, concurrencia), Godog en `features/`, Vitest en el frontend.

**Target Platform**: servidor Linux (imagen Docker, Render).

**Project Type**: web-service + SPA.

**Performance Goals**: la carga completa tarda menos de 1 s sobre SQLite local.

**Constraints**: una sola conexión (`SetMaxOpenConns(1)`); `busy_timeout` de 5 s; sin tablas nuevas.

**Scale/Scope**: 1 proyecto, 4 usuarios, 3 Sprints, 10 historias, ~15 registros de horas, 1 ronda con 3 votos.

## Constitution Check

| Principio | Estado | Nota |
|---|---|---|
| I. Test-First (RED/GREEN/REFACTOR) | PASS | Cada bloque y cada regla de idempotencia con su RED y su GREEN. |
| II. SDD | PASS | Spec 021 numerada a mano. |
| III. BDD | PASS | Godog para la API; Vitest con Given-When-Then para el botón. |
| IV. Clean/Hexagonal | PASS | El handler depende de una interfaz `Cargador`; el dominio no conoce SQLite ni HTTP; `internal/demo` es un adaptador de carga que usa solo constructores del dominio. |
| V. Clean Code, KISS, YAGNI, DRY | PASS | Lista de funciones-bloque, sin registro dinámico ni framework de seeds. |
| VI. Stack | PASS | Sin dependencias nuevas. |

Se reevalúa tras el diseño: sin cambios.

## Project Structure

### Documentation

```text
specs/021-datos-de-prueba/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/api.md
└── tasks.md          # lo genera /speckit-tasks
```

### Source Code

```text
internal/demo/
├── carga.go            # Cargador, Resumen, errores, lista de bloques, transacción
├── bloque_usuarios.go
├── bloque_proyecto.go  # proyecto y equipo
├── bloque_sprints.go
├── bloque_historias.go # historias y arrastrada
├── bloque_horas.go
├── bloque_votacion.go
└── *_test.go
internal/http/demo_handler.go (+ test)   # habilitación + rutas
internal/http/spa.go                      # "/demo" en PrefijosAPI
cmd/api/config.go (+ test)                # DemoHabilitado
cmd/api/main.go                           # wiring
features/datos_prueba.feature + steps_datos_prueba.go
web/vite.config.ts                        # proxy '/demo'
web/src/compartido/demo/                  # api.ts, BotonDatosDePrueba.tsx (+ test)
web/src/features/proyecto/crear.tsx       # una línea de montaje
```

**Structure Decision**: el botón vive en `compartido/demo/` y no en una feature, porque el test de aislamiento (`aislamiento.test.ts`) prohíbe que una feature importe a otra. En `main` no existe aún la «lista de proyectos» (la compone #148 en `rutas.tsx`): el botón se monta en la pantalla de inicio actual (`crear.tsx`), que #148 no modifica según su plan, y seguirá en el inicio cuando llegue la lista.

## Decisiones de diseño

1. **Un paquete `internal/demo` que persiste con SQL en la transacción** (en vez de orquestar casos de uso con transacciones separadas, que rompería la atomicidad). Los valores se validan con `NewUser`, `NewProject`, `NewMembership`, `NewSprint`+`Iniciar`+`Cerrar`, `NewBacklogItem`+`ConDatosEditados`, `NewRegistroHoras`, `NuevaRonda`+`Votar`+`Revelar`. El estado `En progreso` y la fecha de finalización no tienen transición en el dominio (HU-11): se escriben directo.
2. **Bloques**: `type bloque func(ctx, tx *sql.Tx, c *carga) error`; `carga` guarda el día base, los ids ya creados (por nombre de usuario y por clave de historia) y los contadores. Lista ordenada `bloques`; una historia futura agrega una función y una entrada.
3. **Idempotencia y concurrencia**: la primera sentencia de la transacción es una escritura (usuarios), así nunca hay una lectura que luego se promueve a escritura (evita `SQLITE_BUSY` entre procesos). El `INSERT` del proyecto con id 25 falla por clave primaria; se hace rollback y se clasifica leyendo el nombre del proyecto 25: nombre reservado ⇒ `ErrYaCargado`, otro ⇒ `ErrIdOcupado`. En un solo proceso la conexión única serializa las cargas. Sin `BEGIN IMMEDIATE`.
4. **Fallo inyectado**: el `Cargador` recibe la lista de bloques; el test agrega un bloque que falla al final y verifica que no queda nada.
5. **Habilitación en un punto**: `soloSiHabilitado(next)` en `demo_handler.go` consulta una `func() bool`; hoy lee `cfg.DemoHabilitado`. Las rutas se registran siempre y el envoltorio responde 404 con `http.NotFound`.
6. **Día base**: el `Cargador` recibe `func() time.Time` (reloj inyectable); las fechas se truncan al día UTC.
7. **Errores del frontend**: el 409 de «ya cargado» lleva `campo: "ya_cargado"`; el de id ocupado, `campo: "proyecto_id"` (ver `contracts/api.md`).
