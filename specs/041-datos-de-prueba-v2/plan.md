# Implementation Plan: Datos de prueba v2 (issue #235)

**Branch**: `feature/235-datos-de-prueba-v2` | **Date**: 2026-11-09 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/041-datos-de-prueba-v2/spec.md`. Parte de la [021](../021-datos-de-prueba/plan.md) y la modifica; lo que no se menciona aquí queda como en la 021.

## Summary

Se ajusta el paquete `internal/demo` (sin paquetes nuevos): los usuarios pasan a ser reales con contraseña hasheada con **un solo hash** de `demo1234` calculado antes de abrir la transacción, la lista de Sprints pasa de tres a cuatro, las 10 historias se reubican (más 2 nuevas), se suma un bloque de defectos y el clasificador de conflictos distingue el usuario existente. En el frontend, `/ingresar` monta `BotonDatosDePrueba`, que ahora también muestra la tabla de credenciales y un aviso para el 409 de usuario existente. Sin cambio de esquema ni de `MODELO-DE-DATOS.md`.

## Technical Context

**Language/Version**: Go 1.x del `go.mod`; TypeScript + React (Vite) en `web/`

**Primary Dependencies**: las existentes. Ninguna nueva (`internal/service.ServicioContrasenas` ya existe).

**Storage**: SQLite embebido, migraciones 001–007 sin tocar (usa `usuarios.contrasena_hash`, `defectos`, `integrantes_un_scrum_master`).

**Testing**: `testing` de Go (bloques sobre migraciones reales, idempotencia, rollback, concurrencia), Godog en `features/datos_prueba.feature`, Vitest + Testing Library.

**Target Platform**: Linux (Docker, Render plan gratis).

**Project Type**: web-service + SPA.

**Performance Goals**: una sola derivación PBKDF2 (600 000 iteraciones) por carga, no cuatro (SC-007, verificable con un hasheador falso que cuenta llamadas). En pruebas, hasheador de 1 iteración.

**Constraints**: una conexión (`SetMaxOpenConns(1)`), `busy_timeout` de 5 s; el hash se calcula **fuera** de la transacción, para no tener tomada la escritura durante el PBKDF2.

**Scale/Scope**: 4 usuarios, 1 proyecto, 4 Sprints, 12 historias, 20 registros de horas, 2 rondas (6 votos), 5 defectos.

## Constitution Check

| Principio | Estado | Nota |
|---|---|---|
| I. Test-First | PASS | Cada fase del `tasks.md` empieza con RED. Conteos, FK y CHECK se prueban contra las migraciones reales. |
| II. SDD | PASS | Spec 041 numerada a mano; plan, tasks y analyze antes de implementar. |
| III. BDD | PASS | Godog (`datos_prueba.feature`) para la API; Vitest con Given-When-Then para `/ingresar`. |
| IV. Clean/Hexagonal | PASS | `internal/demo` define su propia interfaz `Hasheador { Hashear(string) (string, error) }`; `service.ServicioContrasenas` la satisface sin que `demo` importe `service`. El dominio no cambia. |
| V. Clean Code, KISS, YAGNI, DRY | PASS | Un hash para los cuatro usuarios; los bloques de la v1 se editan, no se reemplazan; el hash lo calcula `Cargar`, no cada bloque. |
| VI. Stack | PASS | Sin dependencias nuevas. |

Se reevalúa tras el diseño: sin cambios.

## Project Structure

### Documentation

```text
specs/041-datos-de-prueba-v2/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/api.md
├── checklists/requirements.md
└── tasks.md          # lo genera /speckit-tasks
```

### Source Code

```text
internal/demo/
├── carga.go             # Hasheador, errores, Resumen.Defectos, orden de bloques, clasificador, hash único
├── bloque_usuarios.go   # lista de 4 usuarios con rol; INSERT estricto con el hash
├── bloque_proyecto.go   # equipo de 4; inicio hoy−30
├── bloque_sprints.go    # 4 Sprints
├── bloque_historias.go  # 12 historias, arrastrada en el Sprint 2
├── bloque_horas.go      # 20 registros, Sprints 1 a 3
├── bloque_votacion.go   # ronda revelada (4 votos) + ronda abierta (2 votos)
├── bloque_defectos.go   # NUEVO: 5 defectos
└── *_test.go
internal/http/demo_handler.go (+ test)   # 409 de usuario existente
cmd/api/main.go                          # pasa `hasheador` a NewCargador
features/datos_prueba.feature + steps_datos_prueba.go
web/src/compartido/demo/                 # BotonDatosDePrueba.tsx, CredencialesDemo.tsx (NUEVO) (+ tests); api.ts sin cambios
web/src/features/autenticacion/Ingresar.tsx (+ test)
README.md
```

**Structure Decision**: lo nuevo de backend cabe en `internal/demo` (un archivo, `bloque_defectos.go`). En el frontend, un componente chico `CredencialesDemo` dentro de `compartido/demo/`: la regla de aislamiento impide que una feature importe a otra, pero `autenticacion` sí puede usar `compartido`. `crear.tsx` no cambia.

## Decisiones de diseño

1. **Orden de bloques**: `usuarios → proyecto → sprints → historias → horas → votación → defectos`. Los usuarios van primero: es la primera escritura de la transacción (regla heredada de la 021: nunca una lectura que se promueva a escritura) y el choque con un usuario existente falla lo antes posible. Defectos va al final porque depende de historias y Sprints. La lista pasa de `bloquesV1` a `bloquesDemo` (único cambio en `atomicidad_test.go`).
2. **Hash único fuera de la transacción**: `Cargador` recibe un `Hasheador`. `Cargar` llama a `Hashear("demo1234")` **una vez** antes de `ejecutarBloques` y lo deja en `carga.hashContrasena`; `bloqueUsuarios` usa ese valor para los cuatro. Misma sal para los cuatro: aceptable porque la contraseña es pública y conocida (decisión del Scrum Master). Costo asumido: una carga que termina en 409 también paga ese hash.
3. **Usuarios**: lista ordenada `usuariosDemo` con `{nombre, rol}`. Cada nombre pasa por `domain.ValidarNombreUsuario` (HU-20) y por `NewUser` para el normalizado; el `INSERT` es **estricto** (sin `OR IGNORE`): un usuario existente viola `UNIQUE` y aborta la transacción. El rol vive en la lista y `bloqueProyecto` arma el equipo de ahí, con `carlos.moya` como único Scrum Master (lo respalda `integrantes_un_scrum_master`).
4. **Conflicto doble** (FR-005/006): el clasificador corre con la transacción ya deshecha y **no depende de qué restricción disparó**. Orden: (1) el proyecto 25 existe con el nombre reservado ⇒ `ErrYaCargado`; (2) si no, algún usuario de la demo existe (normalizado, en el orden de la lista) ⇒ `ErrUsuarioExistente{Nombre}`; (3) si no, el proyecto 25 existe con otro nombre ⇒ `ErrIdOcupado`; (4) si no, el error original. Así gana `ya_cargado` y el mensaje del usuario nombra al primero en conflicto. `ErrUsuarioExistente` es un tipo con el nombre (`errors.As`) y mensaje «ya existe el usuario X».
5. **Cuatro Sprints reutilizando los bloques de la v1**: `bloqueSprints` pasa a una lista de 4 definiciones (goal, inicio, fin, estado) recorridas por el mismo bucle que insertaba tres: #1 Finalizado (−28/−15), #2 Finalizado (−14/−8), #3 Activo (−7/+6), #4 Pendiente (sin goal ni fechas). `sprintIniciado`, `Cerrar` e `insertarSprint` no cambian. `c.sprints` pasa a tener 4 ids; el proyecto empieza en hoy−30.
6. **Reubicación de las historias**: se editan las filas de `historiasDemo` (campo `sprint` de 0 a 4) y se agregan las 2 nuevas. Reparto y fechas de finalización en [data-model.md](data-model.md). La arrastrada (`tasa`) se registra en `c.sprints[1]`; `bloqueHoras` y `bloqueVotacion` solo cambian sus tablas de datos.
7. **Votación**: dos rondas con el mismo código de construcción (`NuevaRonda`, `Votar`, `Revelar` solo en la primera); la abierta se inserta con `revelada_en NULL`. El resumen cuenta 2 rondas y 6 votos.
8. **Defectos**: `bloqueDefectos` usa `domain.NewDefecto` con el Sprint de detección como `activo` (solo se usa su `ID`) y `Defecto.Resolver` con el de resolución, e inserta con SQL las columnas de §4.9. Las FK compuestas quedan satisfechas porque todo pertenece al proyecto 25 y los ids vienen de `c.historias` y `c.sprints`. `creado_en` es relativo al día de la carga. Los «Resueltos» en un Sprint ya finalizado son datos históricos, igual que las horas de la v1: no se exige Sprint activo.
9. **Handler**: `errors.As(err, &demo.ErrUsuarioExistente{})` ⇒ 409 `{"campo":"usuario_existente","mensaje":"ya existe el usuario X"}`.
10. **Frontend**: `CredencialesDemo` (tabla usuario / contraseña / rol, datos constantes) se muestra dentro de `BotonDatosDePrueba` tras el 201 y tras el 409 `ya_cargado`. El 409 `usuario_existente` es un aviso (no error) sin tabla. `Ingresar` monta `<BotonDatosDePrueba />` bajo el formulario; el botón de `crear.tsx` queda como está.
11. **Riesgo conocido**: `Ingresar.test.tsx` simula un único `fetch` para todas las llamadas; con el botón montado habrá un `GET /demo/estado` más, así que esas pruebas deberán distinguir por URL (fase E).
12. **Sin migración**: no hay `CREATE`/`ALTER`; `MODELO-DE-DATOS.md` no se toca.
