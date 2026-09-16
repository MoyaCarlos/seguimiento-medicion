# Seguimiento y Medición de Proyectos de Software

Trabajo Práctico Integrador de la asignatura **Ingeniería y Calidad de Software**
(UTN Facultad Regional San Rafael, 2026).

## Integrantes

-   Moya Carlos Esteban - Scrum Master
-   Iriarte Lopez Ana Valentina - Product Builder
-   Vulcano Candela Nair - Product Builder
-   Martinez Oldani Jimena - Product Builder

## Descripción del proyecto

**Software Metrics & Estimation** es una aplicación para administrar un proyecto
de software gestionado con Scrum, permitiendo estimar, planificar, hacer
seguimiento y medir su calidad. Funciona como un tablero de gestión ágil que
además calcula automáticamente las métricas del proyecto.

Principales funcionalidades:

-   Gestión de proyectos e integrantes.
-   Product Backlog con historias de usuario (prioridad, estado, story points,
    criterios de aceptación).
-   Gestión de Sprints (Sprint Goal, historias asignadas, cierre, historial).
-   Estimación mediante Story Points y **Planning Poker**.
-   Registro de esfuerzo (horas trabajadas) por integrante e historia/tarea.
-   Gestión de defectos (severidad, estado, sprint de detección/resolución).
-   Cálculo de métricas: velocidad del equipo, desviación esfuerzo estimado vs.
    real, porcentaje de historias completadas, defectos detectados/resueltos.
-   Dashboard con representación gráfica de las métricas.
-   Generación de reportes de proyecto/Sprint.

## Stack técnico

-   **Backend:** Go (núcleo de negocio y API REST).
-   **Frontend:** React + Vite (SPA).
-   **Persistencia:** SQLite embebido (`modernc.org/sqlite`).
-   **Testing:** Godog (BDD) + `testing` estándar de Go (TDD).
-   **Reportes:** generación de PDF con Maroto.

## Especificación de funcionalidades (SDD)

El proyecto aplica **Specification-Driven Development** con
[GitHub Spec Kit](https://github.com/github/spec-kit), la herramienta oficial
de GitHub para esta práctica: antes de implementar cada historia de usuario
se escribe una especificación formal (objetivo, reglas de negocio, casos
límite, condiciones de error, criterios de aceptación), siguiendo el flujo
`/speckit.specify` → `/speckit.clarify` → `/speckit.plan` → `/speckit.tasks` →
`/speckit.implement`. Las especificaciones quedan versionadas en `specs/`, y
son la base de la trazabilidad Historia → Spec → Escenario BDD → Tests → Código
que se muestra en la defensa final.

Ya está configurado en el repo para Claude Code y OpenCode — detalles de uso
en [`docs/GUIA-IA-TERMINAL.md`](./docs/GUIA-IA-TERMINAL.md).

Metodología, arquitectura, patrones de diseño y principios aplicados están
documentados en [`AGENTS.md`](./AGENTS.md).
