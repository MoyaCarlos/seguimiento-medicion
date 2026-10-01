<!--
Sync Impact Report
- Version change: (template, no prior content) → 1.0.0
- Modified principles: none (all 6 are new; template had 5 placeholder slots, one added)
- Added sections: Core Principles I-VI, Contexto Académico y Criterios de Evaluación,
  Flujo de Trabajo y Git, Governance
- Removed sections: none
- Deferred items: none — all placeholders resolved from AGENTS.md and team decisions
  recorded in this project's working sessions.
-->

# Software Metrics & Estimation Constitution

Trabajo Práctico Integrador de Ingeniería y Calidad de Software (UTN Facultad
Regional San Rafael, 2026). Esta constitución formaliza, para que todo comando
de Spec Kit (en particular `/speckit.implement`) la respete automáticamente,
las reglas ya acordadas y documentadas en `AGENTS.md`.

## Core Principles

### I. Test-First estricto (NON-NEGOTIABLE)
Toda regla de negocio, cálculo o validación se desarrolla con el ciclo
**RED → GREEN → REFACTOR**: primero el test que falla, luego el código
mínimo que lo hace pasar, luego la limpieza. Cada fase DEBE quedar como su
propio commit, con prefijo obligatorio en el mensaje: `RED: <caso> (falla)`,
`GREEN: <caso>, test en verde`, `REFACTOR: <qué se limpió>`. No se fuerza un
commit de REFACTOR vacío si no hay nada que limpiar. No aplica a código
trivial (getters, DTOs sin lógica). El merge a `main` se hace siempre con
"Create a merge commit", nunca squash — aplastar el historial destruye la
evidencia de este ciclo, que es requisito de evaluación de la cátedra.

### II. Specification-Driven Development (SDD)
Ninguna funcionalidad se implementa sin haber pasado antes por el flujo de
Spec Kit: `/speckit.specify` → `/speckit.clarify` → `/speckit.plan` →
`/speckit.tasks` → `/speckit.implement`. Las especificaciones quedan
versionadas en `specs/<historia>/`. Se trabaja **una Historia de Usuario a
la vez**: `/speckit.implement` no se usa para completar varias historias en
una sola pasada sin revisión humana entre medio.

### III. Behavior-Driven Development (BDD)
Las funcionalidades con comportamiento observable cuentan con escenarios
Given-When-Then (Gherkin) automatizados con Godog, en `/features`. Los
escenarios cubren casos normales, alternativos, límite y de error — no solo
el camino feliz.

### IV. Arquitectura Clean/Hexagonal
La dependencia siempre apunta hacia adentro: `internal/domain` no conoce
SQLite ni HTTP ni ningún detalle de infraestructura. Patrones aplicados:
**Repository** (persistencia detrás de interfaces, con fakes en memoria para
tests), **Strategy** (reglas variables como el consenso de Planning Poker),
**Factory** (constructores que validan invariantes, p.ej. `NewProject`,
`NewBacklogItem`). No se agregan capas ni interfaces que no estén al
servicio de esta separación.

### V. Clean Code, SOLID, KISS, YAGNI, DRY
Nombres claros, funciones pequeñas, sin comentarios que expliquen el "qué".
SRP: un service por caso de uso. Dependency Inversion: los services
dependen de interfaces de repositorio, nunca de una implementación
concreta. KISS: la solución más simple que cumpla el requerimiento. YAGNI:
no se construyen abstracciones para necesidades hipotéticas. DRY: la
lógica de negocio y de cálculo de métricas vive en un único lugar.

### VI. Stack tecnológico obligatorio
Backend en Go (API REST). Frontend en React + Vite (SPA). Persistencia en
SQLite embebido vía `modernc.org/sqlite` (puro Go, sin CGO). BDD con Godog.
TDD con el paquete `testing` estándar de Go. Reportes en PDF con Maroto
(puro Go, sin dependencias externas). No se introduce un dependencia nueva
fuera de este stack sin que el equipo lo decida explícitamente.

## Contexto Académico y Criterios de Evaluación

Este es un Trabajo Práctico universitario, no un producto comercial — los
principios de esta constitución no son sugerencias de estilo, son
requisito de nota. La cátedra evalúa: Producto funcional (25%), SDD/BDD/TDD
(25%), Calidad del software (20%), Gestión del proyecto (20%), Trabajo en
equipo y presentación (10%). El cumplimiento del Principio I (Test-First) y
el Principio II (SDD) es, en los hechos, la mitad de la nota del TP — no
son negociables bajo presión de tiempo.

## Flujo de Trabajo y Git

Branches: `feature/HU-XX-nombre-corto`, `fix/nombre-corto`. Ningún push
directo a `main`: todo cambio entra por Pull Request, con mínimo 1
aprobación de otro integrante, y el PR debe pasar los tests (Godog +
`go test`) antes de mergear. Una historia se considera terminada (Definition
of Done) cuando: el código está mergeado a `main`; los tests unitarios y el
escenario BDD correspondiente están en verde; la especificación SDD de la
historia está escrita en `specs/`; fue revisada y aprobada por al menos otro
integrante; y no tiene defectos abiertos bloqueantes.

## Governance

Esta constitución tiene precedencia sobre cualquier otra práctica informal
del equipo. Se enmienda editando este archivo junto con `AGENTS.md` (deben
mantenerse consistentes entre sí) y commiteando el cambio con un mensaje que
explique el motivo — no se edita en silencio. Versionado semántico:
MAJOR para quitar o redefinir un principio de forma incompatible, MINOR
para agregar un principio o ampliar sustancialmente una guía existente,
PATCH para aclaraciones o correcciones de redacción. El Scrum Master
verifica cumplimiento de esta constitución antes de aprobar cada Pull
Request.

**Version**: 1.0.0 | **Ratified**: 2026-09-07 | **Last Amended**: 2026-10-01
