# Proyecto: Software Metrics & Estimation

## Objetivo
Desarrollar una aplicación que permita realizar la estimación, seguimiento y medición
de proyectos de software, aplicando los conceptos y prácticas de la asignatura
Ingeniería y Calidad de Software (UTN FRSR, 2026).

El proyecto se desarrolla obligatoriamente utilizando Go, Scrum, SDD, BDD, TDD, Git
y herramientas de Inteligencia Artificial como soporte al proceso.

> Este archivo lo cargan automáticamente tanto Claude Code (vía el symlink
> `CLAUDE.md`) como OpenCode. Consejos prácticos de cómo trabajar con estos
> agentes desde la terminal: [`docs/GUIA-IA-TERMINAL.md`](docs/GUIA-IA-TERMINAL.md).

## Equipo y roles
- **Product Architect:** profesores de la cátedra.
- **Agile Enabler:** Moya Carlos Esteban.
- **Product Builders:** Iriarte Lopez Ana Valentina, Vulcano Candela Nair, Martinez Oldani Jimena.
- **Cliente:** rol compartido entre el equipo y el profesor durante validaciones y reviews.

## Metodología de desarrollo
- **Scrum**, con tablero en **GitHub Projects**.
- Cada Historia de Usuario del Product Backlog es un **Issue padre**, con **sub-issues**
  para desglosar el trabajo (especificación SDD, escenarios BDD, implementación backend,
  tests, implementación frontend), reforzando la trazabilidad exigida por la cátedra.
- Sprints: Sprint 0 (preparación) → Sprint 1 (MVP) → Sprint 2 (interfaz) →
  Sprint 3 (funcionalidad y calidad) → Sprint 4 (cierre y entrega final).
  Duración/calendario a confirmar según fecha de entrega final.

### Convención de Git
- Branches: `feature/HU-XX-nombre-corto`, `fix/nombre-corto`.
- Nada de push directo a `main`: todo cambio entra por Pull Request.
- Mínimo 1 aprobación de otro integrante antes de mergear.
- El PR debe pasar los tests (Godog + `go test`) antes de mergear.

### Definition of Done
Una historia se considera terminada cuando:
- Código implementado y mergeado a `main`.
- Tests unitarios (TDD) y escenario BDD correspondiente, ambos en verde.
- Especificación SDD de la historia escrita en `specs/`.
- Revisada y aprobada por al menos otro integrante (PR).
- Sin defectos abiertos bloqueantes para esa historia.

### SDD (Specification-Driven Development)
- Herramienta: **[GitHub Spec Kit](https://github.com/github/spec-kit)**, instalado
  con `specify init --here --integration claude` (o `--integration opencode`).
  Genera las specs versionadas en `specs/<feature>/` — reemplaza lo que se había
  pensado como `/docs/specs`.
- Flujo por historia: `/speckit.constitution` (una vez) → `/speckit.specify` →
  `/speckit.clarify` → `/speckit.plan` → `/speckit.tasks` → `/speckit.implement`.
  En Claude Code los comandos usan guión (`/speckit-specify`); en OpenCode, punto
  (`/speckit.specify`) — mismo comando, distinta sintaxis según la herramienta.
- Cada spec debe cubrir como mínimo: objetivo, entradas, salidas esperadas, reglas
  de negocio, restricciones, casos límite, condiciones de error y criterios de
  aceptación. Para reglas de negocio/casos límite/errores, redactar en formato
  **EARS** (Ubicuo / Basado en evento / Basado en estado / No deseado / Opcional)
  en vez de prosa libre — más preciso y se traduce directo a escenarios BDD.
- `/speckit.implement` se corre **una historia a la vez**, con revisión humana del
  diff antes de seguir — no se usa para completar varias historias de una sola vez
  (ver convención de commits RED/GREEN/REFACTOR más abajo).

### BDD (Behavior-Driven Development)
- **Godog** (Gherkin) para automatizar escenarios Given-When-Then.
- Archivos `.feature` en `/features`, con casos normales, alternativos, límite y de error.

### TDD (Test-Driven Development)
- Ciclo **RED → GREEN → REFACTOR** para reglas de negocio y cálculos.
- Paquete `testing` estándar de Go, foco en la capa de dominio (`/internal/domain`)
  y de casos de uso (`/internal/service`).
- Evidencia del proceso mediante historial de commits: para cálculos de métricas,
  estimación, reglas de negocio y validaciones, cada fase del ciclo es su propio
  commit, con prefijo en el mensaje:
  ```
  RED: agregar test de <caso> (falla)
  GREEN: implementar <caso>, test en verde
  REFACTOR: <qué se limpió>
  ```
  No forzar un commit de REFACTOR vacío si no hay nada que limpiar. No aplica a
  código trivial (getters, DTOs).
- **Merge a `main` siempre con "Create a merge commit"**, nunca squash — si se
  aplasta el PR en un solo commit se pierde la evidencia del ciclo RED/GREEN/REFACTOR
  en el historial de `main`.

## Principios de diseño
- **Clean Code:** nombres claros, funciones pequeñas, sin comentarios que expliquen el "qué".
- **SOLID:** especialmente SRP (un service por caso de uso) y Dependency Inversion
  (los services dependen de interfaces de repositorio, no de SQLite concreto).
- **KISS:** preferir siempre la solución más simple que cumpla el requerimiento.
- **YAGNI:** no construir abstracciones para necesidades hipotéticas (por eso se
  descartó, por ejemplo, un patrón Observer para el Dashboard: las métricas se
  recalculan on-demand, no en tiempo real).
- **DRY:** la lógica de negocio y de cálculo de métricas vive en un único lugar
  (capa de dominio), sin duplicarse entre backend y frontend.

## Arquitectura
Clean/Hexagonal liviana (Ports & Adapters), sin ceremonia de más capas de las necesarias:

```
/cmd/api             -> main.go: wiring de dependencias, arranque del servidor
/internal/domain     -> entidades (Project, Sprint, Story, Defect, Effort)
                        y reglas de negocio / cálculo de métricas
/internal/service    -> casos de uso (CrearProyecto, CerrarSprint, RegistrarEsfuerzo,
                        CalcularMetricas, EjecutarRondaPlanningPoker, ...)
/internal/repository -> interfaces de persistencia (puertos) + implementación SQLite
/internal/http       -> handlers REST, rutas, DTOs
/internal/report     -> generación de informes PDF (Maroto)
/features            -> escenarios BDD (Godog)
/web                 -> frontend React (Vite)
/specs               -> especificaciones SDD versionadas (generadas con Spec Kit)
```

La dependencia siempre apunta hacia adentro: el dominio no conoce SQLite ni HTTP.

## Patrones de diseño aplicados
- **Repository:** abstrae la persistencia detrás de interfaces (`ProjectRepository`,
  `SprintRepository`, etc.), permite tests unitarios con fakes en memoria.
- **Strategy:** para reglas variables como el cálculo de consenso en Planning Poker
  o variantes en el cálculo de métricas.
- **Factory (function):** constructores (`NewProject`, `NewBacklogItem`) que validan
  invariantes al crear la entidad.

## Stack tecnológico
- **Backend:** Go (API REST, núcleo de negocio).
- **Frontend:** React + Vite (SPA), consume la API vía REST/JSON.
- **Persistencia:** SQLite embebido, vía `modernc.org/sqlite` (puro Go, sin CGO,
  evita depender de un compilador C en las máquinas del equipo).
- **Testing BDD:** Godog.
- **Testing TDD:** paquete `testing` estándar de Go.
- **Reportes:** generación de PDF con Maroto (puro Go, sin dependencias externas).
- **Control de versiones:** Git + GitHub, tablero en GitHub Projects con sub-issues.

## Uso de Inteligencia Artificial
- Herramientas usadas como **asistencia al desarrollo** (no forman parte de la
  arquitectura del producto): **Claude**, **DeepSeek API** y modelos disponibles en
  **build.nvidia.com**. Se usan para análisis de requisitos, especificaciones,
  generación de código/tests, refactorización y revisión. Cualquier incorporación
  adicional se documentará en este archivo.
- Todo resultado generado con IA debe ser comprendido, revisado y validado por el
  equipo antes de incorporarse al proyecto. El equipo es responsable de todo el
  código, sea manual o asistido por IA.
