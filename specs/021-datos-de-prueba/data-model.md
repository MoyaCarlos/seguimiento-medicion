# Data Model: datos de prueba (versión 1)

No agrega ni cambia tablas. Esquema: [`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md) §4.1–§4.8. `hoy` = día de la carga (UTC).

## Usuarios
Beta Uno, Beta Dos, Beta Tres, Beta Cuatro (`INSERT OR IGNORE` por `nombre_normalizado`; se leen sus ids después).

## Proyecto 25 y equipo
`id = 25`, nombre «Licencia de conducir (datos de prueba)», descripción breve, `horas_por_story_point = 4`, inicio `hoy-15`, fin `hoy+45`. Integrantes: Beta Uno `scrum_master`; Beta Dos y Beta Tres `product_builder`.

## Sprints
| # | Estado | Goal | Fechas |
|---|---|---|---|
| 1 | Finalizado | «Gestión inicial del trámite» | hoy-14 a hoy-8 |
| 2 | Activo | «Examen teórico y comprobantes» | hoy-7 a hoy+6 |
| 3 | Pendiente | NULL | NULL |

(Se numeran por orden de creación; los ids reales son los que asigna SQLite.)

## Historias (10)
| Clave | Título | Prio | Estado | SP | Sprint | Criterios |
|---|---|---|---|---|---|---|
| turno | Solicitar turno para el examen teórico | M | Completada | 5 | 1 | sí |
| datos | Registrar datos personales del solicitante | M | Completada | 3 | 1 | no |
| apto | Validar el certificado de aptitud psicofísica | S | Completada | 8 | 1 | sí |
| tasa | Pagar la tasa de la licencia en línea | M | Nueva (arrastrada del 1) | 8 | — | no |
| examen | Rendir el examen teórico en línea | M | En progreso | 13 | 2 | sí |
| comprobante | Emitir el comprobante del turno | S | En progreso | 3 | 2 | sí |
| estado | Consultar el estado del trámite | C | Nueva | 2 | 2 | no |
| renovar | Renovar la licencia vencida | S | Nueva | 5 | 3 | no |
| aviso | Avisar por correo el vencimiento | C | Nueva | — | — | no |
| digital | Mostrar la licencia en formato digital | W | Nueva | — | — | no |

Descripciones con Markdown (negritas y listas) en `turno`, `examen` y `renovar`. Completadas con `fecha_finalizacion` dentro del Sprint 1. `tasa` además tiene fila en `historias_arrastradas` (Sprint 1, `estimacion_sp = 8`).

## Horas (minutos enteros)
Total esperado ∈ [50 %, 150 %] de SP × 4 h para las Completadas.

| Historia | Sprint | Registros (integrante, día, h) | Total |
|---|---|---|---|
| turno | 1 | Dos hoy-13 6; Tres hoy-12 8; Uno hoy-11 4 | 18 h (20 h estimadas) |
| datos | 1 | Tres hoy-13 5; Dos hoy-12 6 | 11 h (12 h) |
| apto | 1 | Dos hoy-11 10; Uno hoy-10 8; Tres hoy-9 12 | 30 h (32 h) |
| tasa | 1 | Uno hoy-10 3; Tres hoy-9 4 | 7 h |
| examen | 2 | Dos hoy-6 8; Tres hoy-4 7; Uno hoy-2 5 | 20 h |
| comprobante | 2 | Tres hoy-3 3; Dos hoy-1 2 | 5 h |

## Votación
Historia `aviso`: ronda 1 revelada con Beta Uno = 3, Beta Dos = 5, Beta Tres = 5 (la historia no cambia: no se escribe `estimacion_sp`).

## Resumen de la respuesta
`proyecto_id`, `usuarios` (creados), `integrantes`, `sprints`, `historias`, `registros_horas`, `rondas_votacion`, `votos`.
