# Feature Specification: Despliegue continuo en Render

**Feature Branch**: `feature/139-render`

**Created**: 2026-10-30

**Status**: Draft (Issue #139)

**Input**: User description: "Despliegue continuo en Render (Issue #139). Infraestructura y documentación, sin tocar código de la aplicación. El servicio ya existe, creado a mano en el dashboard, y funciona en https://seguimiento-medicion-sprintometro.onrender.com. Entregables: `render.yaml` como referencia versionada, sección de README y verificación manual anotada en `quickstart.md`."

> La carpeta es `specs/026-render-despliegue` con `--number 26` por indicación del equipo (del 021 al 025 están reservados para otras tareas).

## Clarifications

### Session 2026-10-30

- Q: ¿Qué cuenta como "toca código" para que Render redeploye? → A: Cualquier cambio en un archivo que no coincida con las rutas ignoradas (por ejemplo `cmd/`, `internal/`, `web/`, `Dockerfile`, `go.mod`, `.dockerignore` y `render.yaml`). Si el merge mezcla rutas ignoradas con una no ignorada, redeploya.
- Q: ¿`render.yaml` es una ruta ignorada? → A: No: la configuración actual del dashboard no la ignora y esta spec no la cambia. Mergear el PR que agrega `render.yaml` redeploya el servicio y vacía la base una vez.
- Nota: resueltas por el agente a partir de la configuración dada, sin consulta interactiva; Carlos las valida al revisar la spec.

## Configuración actual del servicio (fuente de verdad)

Creada a mano en el dashboard de Render; esta spec la documenta, no la cambia.

| Ajuste | Valor |
|--------|-------|
| Servicio | `seguimiento-medicion-sprintometro` |
| URL | https://seguimiento-medicion-sprintometro.onrender.com |
| Runtime | Docker, `Dockerfile` de la raíz, contexto: la raíz |
| Plan / región | Free / Oregon (US West) |
| Rama | `main` |
| Auto-Deploy | "After CI Checks Pass" |
| Health Check Path | `/health` |
| Variable de entorno | `DEMO_HABILITADO=true` |
| Included Paths | vacío |
| Ignored Paths | `docs/**`, `specs/**`, `features/**`, `*.md`, `.github/**`, `.specify/**` |

El `Dockerfile` ya fija `PORT=8080`, `DB_PATH=/data/seguimiento.db` y `WEB_DIR=/app/web/dist` y corre sin root (ver [`specs/013-docker-imagen`](../013-docker-imagen/spec.md)). No hay disco persistente: la base se vacía en cada reinicio o deploy, y el servicio se duerme tras 15 minutos sin tráfico.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Cada merge a `main` llega solo a la URL pública (Priority: P1)

Como integrante del equipo, quiero que lo que se mergea a `main` y pasa los checks se publique solo en la URL de Render, sin pasos manuales.

**Why this priority**: Es el objetivo del issue: despliegue continuo para las demos y la entrega.

**Independent Test**: Documentación revisada contra el dashboard (solo lectura) y verificación manual de la URL (quickstart).

**Acceptance Scenarios**:

1. **Given** un merge a `main` que toca código y con los checks `backend` y `frontend` en verde, **When** termina el CI, **Then** Render construye y publica la nueva versión.
2. **Given** un merge a `main` que solo toca rutas ignoradas, **When** termina el CI, **Then** Render no redeploya.
3. **Given** el README, **When** una persona lee "Despliegue en Render", **Then** encuentra la URL, el modo de despliegue y las rutas que no redeployan.

---

### User Story 2 - La configuración queda versionada como referencia (Priority: P1)

Como Agile Enabler, quiero un `render.yaml` en la raíz que describa la configuración del servicio, para poder revisarla en un PR y recrearla a mano si hace falta.

**Why this priority**: La configuración hoy solo vive en el dashboard y no deja rastro en el historial.

**Independent Test**: Validar el archivo contra el esquema oficial de Render (`https://render.com/schema/render.yaml.json`) y comparar cada clave con la tabla de arriba.

**Acceptance Scenarios**:

1. **Given** `render.yaml`, **When** se valida contra el esquema oficial, **Then** no hay errores.
2. **Given** `render.yaml`, **When** se compara con la configuración actual, **Then** cada valor coincide.
3. **Given** `render.yaml` y la spec, **When** se leen, **Then** ambos advierten que no se debe crear un Blueprint en el dashboard (ver Riesgos).

---

### User Story 3 - El equipo sabe qué esperar en una demo (Priority: P1)

Como integrante del equipo, quiero saber que el servicio se duerme, que la base se vacía en cada deploy y cuándo no mergear, para no llegar sorprendido a una presentación.

**Why this priority**: Son límites del plan gratis que afectan directamente la entrega del 16-nov.

**Independent Test**: Revisión del README.

**Acceptance Scenarios**:

1. **Given** el README, **When** se lee, **Then** dice que el servicio se duerme a los 15 minutos y que conviene abrir la URL unos minutos antes de una demo.
2. **Given** el README, **When** se lee, **Then** dice que la base se vacía en cada deploy o reinicio.
3. **Given** el README, **When** se lee, **Then** dice que `DEMO_HABILITADO` controla el botón de datos de prueba (#144).
4. **Given** el README, **When** se lee, **Then** incluye la regla de no mergear en la hora previa a la presentación del 16-nov.

---

### User Story 4 - El servicio real responde como se documenta (Priority: P2)

Como Agile Enabler, quiero comprobar contra el servicio real, sin modificarlo, que lo documentado es cierto.

**Why this priority**: Evita documentar algo que no se verificó.

**Independent Test**: Los tres pedidos del quickstart.

**Acceptance Scenarios**:

1. **Given** el servicio despierto, **When** se pide `GET /health`, **Then** responde 200 con cuerpo `ok`.
2. **Given** el servicio despierto, **When** se pide la raíz con `Accept: text/html` (o se abre en un navegador), **Then** devuelve el `index.html` de la aplicación. Si no se pudo abrir en un navegador real, se anota como no verificado.
3. **Given** el servicio despierto, **When** se pide `GET /projects`, **Then** responde JSON.

### Edge Cases

- Primer pedido con el servicio dormido: tarda decenas de segundos (arranque en frío) y no es un error.
- Merge que mezcla rutas ignoradas y código: redeploya, porque basta un archivo no ignorado.
- Merge a `main` con un check en rojo: no debería llegar a `main` por la protección de rama; si llegara, Render no despliega ("After CI Checks Pass").
- Reinicio por inactividad o deploy: se pierden los datos cargados, incluidos los de prueba.
- `DEMO_HABILITADO` ausente o distinto de `true`: el botón de datos de prueba no se habilita (comportamiento del #144, fuera de esta spec).
- Cambiar solo `render.yaml` (que no está ignorado) redeploya y vacía la base, aunque no cambie nada de la aplicación.
- Un deploy en curso durante la presentación: la app puede quedar no disponible unos minutos.

## Requirements *(mandatory)*

### Functional Requirements (EARS)

- **FR-001** (Ubicuo): El repositorio MUST incluir `render.yaml` en la raíz, que describa un único servicio web con la configuración de la tabla de arriba.
- **FR-002** (Ubicuo): `render.yaml` MUST usar solo claves definidas en el esquema oficial de Render y MUST validar contra él sin errores.
- **FR-003** (Ubicuo): `render.yaml` MUST NOT contener secretos ni valores sensibles.
- **FR-004** (Ubicuo): `render.yaml` MUST llevar un comentario que diga que es solo referencia y que no se crea un Blueprint a partir de él.
- **FR-005** (No deseado): Si alguien crea un Blueprint en el dashboard a partir de `render.yaml`, entonces Render lo aplicaría al servicio existente y lo que el archivo no declare volvería a su valor por defecto; la spec MUST documentar este riesgo.
- **FR-006** (Ubicuo): El README MUST tener una sección "Despliegue en Render" con la URL pública del servicio.
- **FR-007** (Basado en evento): Cuando un merge a `main` toca código y los checks pasan, el README MUST indicar que Render redeploya solo.
- **FR-008** (Ubicuo): El README MUST listar las rutas que no redeployan: `docs/**`, `specs/**`, `features/**`, `*.md`, `.github/**` y `.specify/**`, y aclarar que cualquier otro archivo (incluido `render.yaml`) sí redeploya.
- **FR-009** (Ubicuo): El README MUST indicar que la base se vacía en cada deploy o reinicio.
- **FR-010** (Basado en estado): Mientras el servicio lleve 15 minutos sin tráfico, el README MUST indicar que se duerme y que conviene abrir la URL unos minutos antes de una demo.
- **FR-011** (Ubicuo): El README MUST indicar que `DEMO_HABILITADO` controla el botón de datos de prueba (#144).
- **FR-012** (Ubicuo): El README MUST incluir la regla operativa: no mergear en la hora previa a la presentación del 16-nov.
- **FR-013** (Ubicuo): `quickstart.md` MUST registrar el resultado de la verificación manual: `GET /health` responde 200 `ok`, la raíz abre la aplicación y `GET /projects` devuelve JSON.
- **FR-014** (Ubicuo): La verificación MUST usar solo pedidos de lectura (`GET`) y MUST NOT modificar el servicio ni su configuración.
- **FR-015** (Ubicuo): Esta tarea MUST NOT modificar código de la aplicación, el `Dockerfile`, el esquema de datos ni agregar dependencias.

### Riesgos

- **Blueprint sobre el servicio existente**: crear un Blueprint a partir de `render.yaml` haría que Render sincronice el archivo con cada push; cualquier ajuste del dashboard que el archivo no declare (por ejemplo la región o las variables) volvería a su valor por defecto. Por eso el archivo es solo referencia y no se crea ningún Blueprint.
- **Deriva**: si alguien cambia el dashboard y no el archivo (o al revés), el archivo queda desactualizado. Un cambio de configuración se registra en el mismo PR que actualiza `render.yaml` y el README.

### Key Entities

Ninguna. No se toca el esquema ([`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md)).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: `render.yaml` valida contra el esquema oficial con 0 errores y el 100 % de sus valores coincide con la tabla de configuración.
- **SC-002**: El README responde, en una sola sección, las 7 preguntas de FR-006 a FR-012 (URL, cómo se despliega, rutas ignoradas, base vacía, sueño, `DEMO_HABILITADO`, regla de la hora previa).
- **SC-003**: Los 3 pedidos de verificación manual dan el resultado esperado y quedan anotados con fecha.
- **SC-004**: El diff del PR no toca código de la aplicación, el `Dockerfile` ni dependencias.

## Assumptions

- El servicio ya existe y funciona; esta tarea no lo crea ni lo cambia, no toca el dashboard y no maneja secretos.
- La región es Oregon (US West), dato confirmado por el equipo.
- Fuera de alcance: disco persistente, otros proveedores, cambios de código o de esquema, crear un Blueprint.
- Los nombres exactos de las claves de `render.yaml` se validan contra `https://render.com/schema/render.yaml.json`; lo que no se pueda validar se deja fuera y se informa.
- No aplica el ciclo RED/GREEN/REFACTOR: es infraestructura y documentación.
