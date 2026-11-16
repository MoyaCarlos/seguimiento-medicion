# Feature Specification: Datos de prueba finales (siete proyectos)

**Feature Branch**: `feature/052-datos-de-prueba-final`
**Status**: Borrador, pendiente de revisión.
**Spec base**: [`specs/021-datos-de-prueba`](../021-datos-de-prueba/spec.md) y [`specs/041-datos-de-prueba-v2`](../041-datos-de-prueba-v2/spec.md). Esta spec **no las redefine**: habilitación por `DEMO_HABILITADO`, proyecto 25, transacción única, bloques, fechas relativas, idempotencia y concurrencia siguen vigentes. Sin cambio de esquema ni de contrato HTTP; el `Resumen` solo gana campos al final.

## Objetivo

Para la demo del 16-nov, que `POST /demo/datos` cargue, además del proyecto 25, **siete proyectos** con Ana como Scrum Master y los mismos cuatro usuarios, de modo que cada feature importante se vea con datos reales sin armar nada a mano.

## Entradas / Salidas

- Entrada: la carga de siempre (sin cuerpo nuevo). Un bloque nuevo, `bloqueProyectosAdicionales`, corre al final de `bloquesDemo`, en la misma transacción y con `c.usuarios`.
- Salida: el `Resumen` de siempre (sigue describiendo el proyecto 25) más `proyectos_adicionales` y `adicionales_*` al final, con lo creado por el bloque nuevo.

## Proyectos (ids asignados por la base)

| # | Nombre (todos terminan en «(datos de prueba)») | Estado | Rasgo |
|---|---|---|---|
| 1 | Plataforma de turnos médicos | En curso | 4 Sprints (2 Finalizado, 1 Activo, 1 Pendiente), ~25 historias, Markdown, arrastre, horas, rondas, defectos, `horas_por_story_point` |
| 2 | Sistema de biblioteca | Finalizado (cierre) | Sprints Finalizado + `fecha_cierre` |
| 3 | Tienda de barrio | Finalizado (fecha) | Sprints Finalizado, fin pasado; Jimena no integra |
| 4 | App de reservas | Planificado | Fechas futuras, backlog sin estimar |
| 5 | Portal de alumnos | En curso | Un integrante con `fecha_baja` |
| 6 | Proyecto en blanco | Planificado | Solo el equipo |
| 7 | Campaña de donaciones | En curso | Sprints Finalizado y fecha de fin futura (#315) |

## Requisitos (EARS)

- **FR-001** (ubicuo): El sistema DEBE crear los siete proyectos con Ana como único Scrum Master de cada uno.
- **FR-002** (ubicuo): `CalcularEstadoProyecto` DEBE dar a cada proyecto el estado de la tabla.
- **FR-003** (basado en evento): CUANDO el proyecto 7 se cierra, el sistema DEBE pasarlo a «Finalizado».
- **FR-004** (basado en estado): MIENTRAS el proyecto 2 está cerrado, toda escritura DEBE responder 409 «El proyecto está cerrado».
- **FR-005** (ubicuo): El proyecto 3 NO DEBE incluir a Jimena entre sus integrantes.
- **FR-006** (ubicuo): El integrante dado de baja del proyecto 5 NO DEBE figurar entre los activos, DEBE figurar en el histórico, y sus horas DEBEN seguir existiendo.
- **FR-007** (no deseado): SI el integrante del proyecto 5 fue dado de baja, ENTONCES NO DEBE quedar ningún voto suyo en una ronda sin revelar.
- **FR-008** (ubicuo): El proyecto 1 DEBE tener una ronda abierta donde Carlos no votó, una revelada con consenso y una con divergencia seguida de una segunda ronda.
- **FR-009** (no deseado): SI una historia tiene ronda revelada o está en un Sprint Activo/Finalizado, su `estimacion_sp` ya cargada DEBE ser coherente con D2b: no se carga ninguna historia de un Sprint Finalizado sin estimar.
- **FR-010** (no deseado): SI el bloque falla, ENTONCES NO DEBE quedar ningún dato, ni el proyecto 25 ni los nuevos.
- **FR-011** (ubicuo): El proyecto 25 y su `Resumen` NO DEBEN cambiar.

## Casos límite y errores

Segunda carga → 409 `ya_cargado` (lo da el proyecto 25). Dos cargas simultáneas → una gana. Una historia con título largo y otra sin criterios; Markdown sin imágenes externas.

## Criterios de aceptación

Los 11 tests de `internal/demo/proyectos_adicionales_test.go` (convivencia, estados, idempotencia, atomicidad, concurrencia, baja, cierre, regla #315, no integrante, votación en vivo, resumen).
