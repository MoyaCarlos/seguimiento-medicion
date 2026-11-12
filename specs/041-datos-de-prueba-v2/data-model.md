# Data Model: Datos de prueba v2

No agrega tablas ni columnas. Las tablas están en [`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md) (§4.2 `usuarios`, §4.3 `integrantes`, §4.5 `sprints`, §4.6 `historias_arrastradas`, §4.7 votación, §4.8 horas, §4.9 `defectos`). Esto es solo el contenido que carga la demo. Fechas en días relativos a hoy.

## Usuarios (contraseña `demo1234` para todos)

| Usuario | Rol en el proyecto 25 |
|---|---|
| `carlos.moya` | Scrum Master (único) |
| `ana.iriarte` | Product Builder |
| `candela.vulcano` | Product Builder |
| `jimena.martinez` | Product Builder |

## Sprints

| # | Estado | Inicio | Fin | Sprint Goal |
|---|---|---|---|---|
| 1 | Finalizado | −28 | −15 | Alta del trámite y validación médica |
| 2 | Finalizado | −14 | −8 | Turno y comprobante |
| 3 | Activo | −7 | +6 | Examen teórico y reprogramación de turnos |
| 4 | Pendiente | — | — | — (sin goal, por el CHECK de §4.5) |

Proyecto: inicio −30, fin +45.

## Historias (12)

| Clave | Sprint | Estado | SP | Final. | Origen |
|---|---|---|---|---|---|
| `tipos` | 1 | Completada | 2 | −18 | **nueva** («Consultar los tipos de licencia y sus requisitos») |
| `datos` | 1 | Completada | 3 | −17 | v1 |
| `apto` | 1 | Completada | 8 | −16 | v1 |
| `turno` | 2 | Completada | 5 | −10 | v1 |
| `comprobante` | 2 | Completada | 3 | −9 | v1 (antes En progreso, Sprint 2) |
| `examen` | 3 | En progreso | 13 | — | v1 |
| `reprogramar` | 3 | En progreso | 3 | — | **nueva** («Reprogramar un turno reservado») |
| `estado` | 3 | Nueva | 2 | — | v1 |
| `renovar` | 4 | Nueva | 5 | — | v1 |
| `tasa` | backlog | Nueva | 8 | — | v1, **arrastrada del Sprint 2** (`historias_arrastradas`, 8 SP) |
| `aviso` | backlog | Nueva | sin estimar | — | v1, ronda revelada |
| `digital` | backlog | Nueva | sin estimar | — | v1, ronda abierta |

Velocidad: Sprint 1 = 13 SP, Sprint 2 = 8 SP.

## Horas (20 registros, `horas_por_story_point` = 4)

| Sprint | Historia | Registros (usuario, día, horas) |
|---|---|---|
| 1 | `tipos` | ana −27 5 |
| 1 | `datos` | candela −26 5; ana −25 6; jimena −25 3 |
| 1 | `apto` | ana −24 10; carlos −23 8; candela −22 12 |
| 2 | `turno` | ana −13 6; candela −12 8; carlos −11 4 |
| 2 | `comprobante` | candela −9 4; ana −8 3 |
| 2 | `tasa` | carlos −10 3; candela −9 4; jimena −9 3 |
| 3 | `examen` | ana −6 8; candela −4 7; carlos −2 5 |
| 3 | `reprogramar` | jimena −3 4; candela −1 3 |

Completadas: `tipos` 5 h (8 h estimadas, 62 %), `datos` 14 (12), `apto` 30 (32), `turno` 18 (20), `comprobante` 7 (12). Todas dentro del 50 %–150 % de FR-011.

## Votación

| Historia | Ronda | Votos |
|---|---|---|
| `aviso` | revelada | carlos 3, ana 5, candela 5, jimena 8 |
| `digital` | abierta (`revelada_en NULL`) | ana 3, jimena 5 |

## Defectos (5)

| # | Historia | Severidad | Estado | Detectado en | Resuelto en |
|---|---|---|---|---|---|
| 1 | `apto` | Alta | Resuelto | Sprint 1 | Sprint 2 |
| 2 | `turno` | Media | Resuelto | Sprint 2 | Sprint 2 |
| 3 | `comprobante` | Baja | Pendiente | Sprint 2 | — |
| 4 | `examen` | Alta | Pendiente | Sprint 3 | — |
| 5 | `reprogramar` | Baja | Resuelto | Sprint 3 | Sprint 3 |

Las descripciones concretas se fijan en la tarea de implementación (texto de dominio de licencias).

## Resumen de la respuesta 201

`{"proyecto_id":25,"usuarios":4,"integrantes":4,"sprints":4,"historias":12,"registros_horas":20,"rondas_votacion":2,"votos":6,"defectos":5}`
