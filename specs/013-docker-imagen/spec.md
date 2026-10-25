# Feature Specification: Imagen Docker multi-stage sin root para correr la app

**Feature Branch**: `feature/132-docker-imagen`

**Created**: 2026-10-24

**Status**: Implementada, pendiente de revisión por PR (Issue #132)

**Input**: User description: "Imagen Docker multi-stage sin root para correr la app (Issue #132). UNA sola imagen: el binario Go sirve la API y el frontend (web/dist). PORT y DB_PATH configurables por entorno; el servidor replica la regla del proxy de Vite para servir el SPA; WEB_DIR configurable; Dockerfile multi-stage con usuario no root; .dockerignore; sección de README. Fuera de alcance: Azure, Trivy, registro de imágenes, CI de Docker. No se toca el esquema de datos."

> La carpeta es `specs/013-docker-imagen` por indicación del equipo (el issue menciona `specs/013-docker/`).

## Clarifications

### Session 2026-10-24

- Q: ¿Qué responde el servidor a un `GET` con `Accept: text/html` en una ruta que no es de la API ni un archivo (ej. `/`, `/otra-pantalla`)? → A: El `index.html`, igual que el fallback de SPA del servidor de desarrollo de Vite. Un `GET` sin `text/html` a un archivo inexistente responde 404.
- Nota: resuelta por el agente a partir de `web/vite.config.ts` (comportamiento actual en desarrollo), sin consulta interactiva; Carlos la valida al revisar la spec.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Correr toda la app con un solo contenedor (Priority: P1)

Como integrante del equipo, construyo la imagen desde un clon limpio y la corro con un único comando; abro el navegador en el puerto publicado y veo la interfaz, y la API responde en el mismo puerto, sin instalar Go ni Node.

**Why this priority**: Es el objetivo del issue y lo que la cátedra pide para el despliegue posterior.

**Independent Test**: `docker build` y `docker run -p 8080:8080 -v ...`; `GET /health` responde 200, `GET /` y `GET /projects/1/status` con `Accept: text/html` devuelven el `index.html`, y `POST /projects` crea un proyecto (201).

**Acceptance Scenarios**:

1. **Given** un clon limpio del repositorio, **When** construyo la imagen, **Then** el build termina sin errores.
2. **Given** el contenedor en ejecución con el puerto publicado, **When** pido `GET /health`, **Then** responde 200.
3. **Given** el contenedor en ejecución, **When** pido `GET /` o `GET /projects/1/status` con `Accept: text/html`, **Then** devuelve el `index.html`.
4. **Given** el contenedor en ejecución, **When** creo un proyecto con `curl -X POST localhost:8080/projects -H 'Content-Type: application/json' -d '{"nombre":"Demo","creador":"Ana","fecha_inicio":"2026-10-01","fecha_fin":"2026-12-01"}'`, **Then** responde 201.

---

### User Story 2 - Los datos sobreviven al reinicio (Priority: P1)

Como integrante del equipo, si detengo y vuelvo a levantar el contenedor con el mismo volumen, los proyectos creados siguen ahí.

**Why this priority**: Sin persistencia la imagen no sirve para nada más que una demo descartable.

**Independent Test**: crear un proyecto con el `curl` de US1, `docker rm -f`, volver a correr con el mismo volumen y `GET /projects/1`.

**Acceptance Scenarios**:

1. **Given** un proyecto creado con un volumen montado en `/data`, **When** reinicio el contenedor con el mismo volumen, **Then** `GET /projects/1` devuelve el proyecto.
2. **Given** un volumen nuevo y vacío, **When** arranca el contenedor, **Then** la base se crea y migra sin errores de permisos.

---

### User Story 3 - El contenedor no corre como root (Priority: P1)

Como Agile Enabler, necesito que el proceso no corra con privilegios de superusuario.

**Why this priority**: Requisito explícito de la cátedra y del issue.

**Independent Test**: `docker inspect` muestra un usuario no root en la configuración de la imagen.

**Acceptance Scenarios**:

1. **Given** la imagen construida, **When** inspecciono su usuario, **Then** no es root (UID distinto de 0).
2. **Given** el contenedor corriendo, **When** la app escribe su base en `/data`, **Then** lo logra sin ser root.

---

### User Story 4 - Configurar puerto, base y estáticos por entorno (Priority: P2)

Como persona que despliega, configuro `PORT`, `DB_PATH` y `WEB_DIR` por variables de entorno; sin ellas, la app se comporta como hoy.

**Why this priority**: Habilita la imagen, pero el default conserva el comportamiento de desarrollo local.

**Independent Test**: test unitario de la lectura de configuración con un `getenv` falso (con y sin valores).

**Acceptance Scenarios**:

1. **Given** ninguna variable definida, **When** se lee la configuración, **Then** puerto `8080`, base `seguimiento.db`, estáticos `web/dist`.
2. **Given** `PORT`, `DB_PATH` y `WEB_DIR` definidas, **When** se lee la configuración, **Then** se usan esos valores.

---

### User Story 5 - El servidor Go sirve el SPA sin romper la API (Priority: P1)

Como usuario del navegador, entro a una pantalla por su URL (que comparte ruta con un endpoint de la API) y veo la interfaz; las llamadas a la API siguen devolviendo JSON.

**Why this priority**: Sin esto una sola imagen no puede servir ambas cosas; hoy lo resuelve el proxy de Vite en desarrollo.

**Independent Test**: tests HTTP contra el servidor con un directorio de estáticos temporal.

**Acceptance Scenarios**:

1. **Given** estáticos disponibles, **When** `GET /projects/1/status` con `Accept: text/html`, **Then** responde el `index.html`.
2. **Given** estáticos disponibles, **When** `GET /projects/1/status` sin `Accept: text/html`, **Then** responde la API (igual que hoy).
3. **Given** estáticos disponibles, **When** `GET` de un archivo de `assets/`, **Then** responde ese archivo.
4. **Given** estáticos disponibles, **When** `POST` o `PUT` a una ruta de la API (aunque el `Accept` incluya `text/html`), **Then** va a la API.
5. **Given** que el directorio de estáticos no existe, **When** arranca el servidor, **Then** arranca igual y solo sirve la API.

### Edge Cases

- `Accept` con varios tipos (`text/html,application/xhtml+xml,...`) que incluyen `text/html`: se trata como pedido de página.
- `GET /health` con `Accept: text/html` (navegador): `/health` no es ruta proxeada por Vite, así que sigue respondiendo `ok`. La regla de FR-004 aplica solo a `/backlog`, `/projects` y `/sprints`.
- `GET` de un asset inexistente: 404, no `index.html`.
- Volumen montado con permisos de root en el host: queda documentado cómo evitarlo (ver Assumptions).
- `PORT` o `DB_PATH` definidas pero vacías: se usa el default.

## Requirements *(mandatory)*

### Functional Requirements (EARS)

- **FR-001** (Ubicuo): La imagen MUST contener un único binario que sirve la API y el frontend.
- **FR-002** (Ubicuo): La lectura de configuración MUST ser una función que recibe un `getenv` y devuelve puerto, ruta de la base y directorio de estáticos.
- **FR-003** (No deseado): Si `PORT`, `DB_PATH` o `WEB_DIR` no están definidas o están vacías, entonces el sistema MUST usar `8080`, `seguimiento.db` y `web/dist` respectivamente.
- **FR-004** (Basado en evento): Cuando llega un `GET` a una ruta de la API que Vite hoy proxea (`/backlog`, `/projects`, `/sprints` y subrutas) con `Accept` que incluye `text/html`, el sistema MUST responder el `index.html`.
- **FR-005** (Basado en evento): Cuando llega un `GET` o `HEAD` de un archivo existente en el directorio de estáticos, el sistema MUST servirlo.
- **FR-005b** (Basado en evento): Cuando llega un `GET` con `Accept` que incluye `text/html` a una ruta que no es de la API ni un archivo existente (ej. `/`), el sistema MUST responder el `index.html`; sin `text/html`, 404.
- **FR-006** (Ubicuo): Toda otra petición (incluidos `POST`/`PUT` y `GET` sin `text/html`) MUST llegar a la API sin cambios de comportamiento.
- **FR-007** (Basado en estado): Mientras el directorio de estáticos no exista, el sistema MUST arrancar y servir solo la API.
- **FR-008** (Ubicuo): El Dockerfile MUST ser multi-stage: build del frontend, build del binario sin CGO, imagen final mínima con solo el binario y `web/dist`.
- **FR-009** (Ubicuo): La imagen final MUST correr con usuario no root.
- **FR-010** (Ubicuo): La imagen MUST fijar `DB_PATH=/data/seguimiento.db`, `PORT` y `EXPOSE` coherentes, y `/data` MUST ser escribible por el usuario no root.
- **FR-011** (Ubicuo): El repositorio MUST incluir `.dockerignore` que excluya al menos `node_modules` (en cualquier nivel, incluido `web/node_modules`), `.git`, `*.db`, `docs/`, `specs/` y `web/dist`.
- **FR-012** (Ubicuo): El README MUST tener una sección "Correr con Docker" con `docker build` y `docker run` (`-p`, `-v` a `/data`).
- **FR-013** (Ubicuo): El esquema de datos MUST NOT modificarse.
- **FR-014** (Ubicuo): Los prefijos de la API que proxea `web/vite.config.ts` MUST estar incluidos en los que usa `ServirSPA`; un test lo verifica.

### Key Entities

Ninguna nueva. Usa la base existente descrita en [`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md), sin cambios.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Desde un clon limpio, la imagen se construye con un solo comando y sin errores.
- **SC-002**: El contenedor responde 200 en `/health` y `GET /` y `GET /projects/1/status` con `Accept: text/html` devuelven el `index.html`, y `POST /projects` responde 201.
- **SC-003**: Tras reiniciar con el mismo volumen, el 100 % de los proyectos creados siguen disponibles por `GET /projects/{id}`.
- **SC-004**: El usuario efectivo del contenedor no es root.
- **SC-005**: `go test ./...` y `npm test` siguen en verde; los tests nuevos cubren la configuración, la regla de servido y la coherencia entre los prefijos de Vite y los del servidor Go (FR-014).

## Assumptions

- La creación de un proyecto desde la pantalla depende del PR #121 y se valida cuando esté mergeado; no es criterio de este issue.
- Una sola imagen (decisión de Carlos); sin segunda imagen para el frontend.
- Fuera de alcance: despliegue a Azure, Trivy, registro de imágenes y CI de Docker.
- Se usa el volumen o bind mount de `/data` con permisos para el UID no root de la imagen (en bind mounts el usuario debe ajustar el dueño; en volúmenes con nombre Docker lo copia de la imagen).
- La regla de negociación por `Accept` replica la de `web/vite.config.ts`; si cambian los prefijos proxeados allí, hay que cambiarlos también en Go (lo detecta un test, FR-014).
- El `Dockerfile` depende de las versiones de `go.mod` (Go 1.26.x) y de Node 24.
