# Seguimiento y Medición de Proyectos de Software

Trabajo Práctico Integrador de la asignatura **Ingeniería y Calidad de Software**
(UTN Facultad Regional San Rafael, 2026).

## Integrantes

-   Moya Carlos Esteban - Scrum Master
-   Iriarte Lopez Ana Valentina - Product Builder
-   Vulcano Candela Nair - Product Builder
-   Martinez Oldani Jimera - Product Builder

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

Metodología, arquitectura, patrones de diseño y principios aplicados están
documentados en [`claude.md`](./claude.md).
