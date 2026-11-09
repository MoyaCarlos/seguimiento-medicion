# Implementation Plan: Autenticación de usuarios (HU-20)

**Branch**: `feature/HU-20-autenticacion` | **Date**: 2026-11-03 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/034-hu-20-autenticacion/spec.md`

## Summary

Registro, ingreso, «quién soy» y cierre de sesión bajo `/auth`, sin exigir sesión en
ninguna ruta existente. La contraseña se guarda con PBKDF2-SHA-256 y sal propia
(`crypto/pbkdf2` de la biblioteca estándar); la sesión es un token aleatorio de 32 bytes
(base64url) del que solo se persiste su SHA-256, con vencimiento fijo de 7 días. El
frontend suma una feature `autenticacion` con un contexto de sesión (token en
`localStorage`), pantallas de ingreso y registro, menú de usuario y ruta protegida.
Backend primero; frontend después, en el mismo PR. Sin dependencias nuevas.

## Technical Context

**Language/Version**: Go 1.26 (backend); TypeScript 6 + React 19, Vite 8 (frontend)

**Primary Dependencies**: ninguna nueva. `crypto/pbkdf2`, `crypto/rand`, `crypto/sha256`,
`crypto/subtle`, `encoding/base64` de la biblioteca estándar.

**Storage**: SQLite (`modernc.org/sqlite`). Migración `006_autenticacion.sql`: columna
`usuarios.contrasena_hash` (NULL) y tabla `sesiones`.

**Testing**: `go test` (tablas de casos en dominio, fakes en servicios, SQLite real en
repositorios y esquema); Godog para la API; Vitest + Testing Library para la interfaz.

**Target Platform**: servidor Linux/Docker + navegador.

**Project Type**: web (API REST + SPA).

**Performance Goals**: sin meta de tiempo automatizada para el ingreso y el registro (el costo
de PBKDF2 es deliberado y sus parámetros se verifican en el test del servicio de contraseñas); «quién soy» es una consulta por
índice único.

**Constraints**: no cambia ninguna ruta, DTO ni escenario existente; `/health` y el resto
siguen sin sesión; contraseña y token nunca en respuestas de error ni en logs.

**Scale/Scope**: 4 rutas nuevas, 4 casos de uso, 1 servicio de contraseñas, 1 emisor de
sesiones, 5 puertos chicos, 1 migración, 1 feature de frontend.

## Constitution Check

| Principio | Estado |
|---|---|
| I. Test-First | RED/GREEN/REFACTOR por comportamiento; orden en [quickstart](quickstart.md) y en `tasks.md` |
| II. SDD | Spec → clarify → plan → tasks → analyze antes de implementar |
| III. BDD | `features/autenticacion.feature` con los 7 escenarios del issue (Godog); en el frontend, pruebas de componente con el escenario en el nombre |
| IV. Clean/Hexagonal | Reglas del nombre y la contraseña en `internal/domain` (puro); puertos en `internal/repository`; hash y token en `internal/service`; HTTP solo traduce |
| V. SOLID/KISS/YAGNI | Puertos mínimos (ISP), reloj inyectado, sin middleware ni limpieza programada de sesiones (nadie los necesita todavía) |
| VI. Stack | Sin dependencias nuevas; `crypto/pbkdf2` es de la biblioteca estándar desde Go 1.24 |

Sin violaciones. Re-chequeo post-diseño: igual. Dos desvíos de las convenciones de
`MODELO-DE-DATOS.md` §2 (instante con regla de negocio, cierre sin borrado) quedan
justificados en [research.md](research.md) D5 y D6 y se documentan en el mismo PR.

## Project Structure

### Documentation (this feature)

```text
specs/034-hu-20-autenticacion/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── api-auth.md
│   └── frontend-sesion.md
├── checklists/requirements.md
└── tasks.md            # lo genera /speckit-tasks
```

### Source Code (repository root)

```text
internal/
├── domain/
│   ├── usuario_credenciales.go      # ValidarNombreUsuario, ValidarContrasena (+ _test, tabla de casos)
│   ├── sesion.go                    # Sesion, SesionConUsuario, DuracionSesion, Vigente (+ _test)
│   └── errors.go                    # + ErrCredencialesInvalidas, ErrNoAutenticado, ErrNombreUsuarioEnUso
├── repository/
│   ├── migraciones/006_autenticacion.sql
│   ├── credencial_repository.go     # RegistroCredenciales, BuscadorCredenciales (puertos)
│   ├── sesion_repository.go         # CreadorSesiones, BuscadorSesiones, CerradorSesiones (puertos)
│   ├── sqlite_credencial.go         # (+ _test, SQLite real)
│   ├── sqlite_sesion.go             # (+ _test, SQLite real)
│   └── esquema_test.go              # + columnas e índices de 006
├── service/
│   ├── contrasenas.go               # ServicioContrasenas: Hashear, Verificar (PBKDF2 + ConstantTimeCompare)
│   ├── emisor_sesiones.go           # token aleatorio + huella + persistir; reloj inyectado
│   ├── registrar_usuario.go         # RegistrarUsuario
│   ├── ingresar.go                  # Ingresar
│   ├── obtener_sesion.go            # ObtenerSesion («quién soy»)
│   ├── cerrar_sesion.go             # CerrarSesion
│   └── fakes_test.go                # + fakes en memoria de los puertos nuevos
└── http/
    ├── auth_dto.go
    ├── auth_handler.go              # Registro, Ingreso, Yo, Salir; lectura de Authorization: Bearer
    ├── sprint_handler.go            # escribirError: + 401 y 409 nuevos
    └── spa.go                       # PrefijosAPI: + "/auth"
cmd/api/main.go                      # cableado y 4 rutas /auth/*
features/
├── autenticacion.feature
└── steps_autenticacion.go           # + registro en features_test.go
MODELO-DE-DATOS.md                   # §2, §3, §4.2, §4.11 nueva, §8, §9

web/
├── vite.config.ts                   # + '/auth': apiProxy
└── src/
    ├── compartido/api/
    │   ├── client.ts                # agrega Authorization si hay token; avisa ante 401
    │   └── sesionHttp.ts            # punto de enganche (obtenerToken, alNoAutenticado)
    ├── features/autenticacion/
    │   ├── api.ts                   # registrar, ingresar, yo, salir
    │   ├── almacenamientoToken.ts   # localStorage dentro de try/catch
    │   ├── SesionContext.tsx        # proveedor y useSesion
    │   ├── Ingresar.tsx / Registro.tsx
    │   ├── RutaProtegida.tsx / MenuUsuario.tsx
    │   ├── estilos.css
    │   └── *.test.tsx / *.test.ts
    ├── App.tsx                      # SesionProvider + MenuUsuario en el encabezado
    └── rutas.tsx                    # /ingresar, /registro y el resto bajo RutaProtegida
```

**Structure Decision**: se respeta la estructura hexagonal existente. Los puertos nuevos
viven junto a los repositorios actuales, un caso de uso por archivo en `service`, y el
frontend sigue la regla de aislamiento por feature (`features/aislamiento.test.ts`): el
cliente compartido no importa la feature, la feature se engancha a él.

## Decisiones técnicas del equipo y ajustes del plan

Se usan tal como se indicaron. Lo que el diseño agrega o precisa, con su motivo:

| # | Ajuste | Motivo |
|---|---|---|
| 1 | El 401 de «sesión vencida» en el cliente solo actúa si la petición llevó token y no es `/auth/ingreso` ni `/auth/registro` | un 401 de «Usuario o contraseña incorrectos» no es una sesión vencida; sin esto el ingreso fallido redirigiría y mostraría el aviso equivocado |
| 2 | `/auth` se agrega a `PrefijosAPI` y al proxy de Vite | un test ya exige que ambos coincidan |
| 3 | El cierre de sesión marca `cerrada_en`; no borra la fila | `MODELO-DE-DATOS.md` §2: «ninguna historia borra datos» |
| 4 | `vence_en` se guarda como columna y se usa en la regla de vigencia | §2 dice que los instantes son solo auditoría; es una excepción documentada (D5) |
| 5 | Longitud máxima de contraseña: 128 caracteres | cierra el pendiente del spec; evita que una contraseña enorme encarezca el hash |
| 6 | Parámetros de PBKDF2 inyectables; en los tests, pocas iteraciones | los tests de dominio y servicio no deben tardar segundos por ingreso |
| 7 | El hash guardado es autodescriptivo (`pbkdf2-sha256$iter$sal$hash`) | permite subir las iteraciones más adelante sin migrar a nadie |
| 8 | El cliente compartido recibe `obtenerToken` y `alNoAutenticado` por enganche | respeta la regla de aislamiento entre `compartido` y las features |
| 9 | Rutas protegidas del frontend: todas salvo `/ingresar` y `/registro`, incluidas la raíz y «no encontrada»; al redirigir a ingreso se conserva la ruta a la que volver | decidido por el Scrum Master (ver [research.md](research.md) D12), escenario 6; se aplica en un solo lugar (`rutas.tsx`) y ningún test del backend cambia |

## Complexity Tracking

Sin violaciones de la constitución que justificar.
