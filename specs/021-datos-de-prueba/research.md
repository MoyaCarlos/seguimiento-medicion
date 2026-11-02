# Research: Cargar datos de prueba

## R1. ¿Dónde vive la carga?
- **Decisión**: paquete `internal/demo`, con SQL propio dentro de una transacción.
- **Por qué**: orquestar `CrearProyecto`, `CrearSprint`, etc. abre una transacción por caso de uso y no sería atómico; además `CrearProyecto` no permite fijar el id.
- **Descartado**: repositorio genérico `CargaRepository` con puerto separado (una sola implementación: YAGNI); sembrar con un `.sql` embebido (saltea las validaciones del dominio).

## R2. ¿Cómo se garantiza idempotencia y concurrencia?
- **Decisión**: clave primaria del proyecto 25 + transacción cuya primera sentencia es una escritura; clasificación del 409 tras el rollback.
- **Por qué**: con `SetMaxOpenConns(1)` las cargas de un proceso se serializan en el pool; entre procesos, `busy_timeout` espera el bloqueo de escritura y la PK decide. No hay lectura previa que promover a escritura.
- **Descartado**: `BEGIN IMMEDIATE` (requiere tocar el DSN de `AbrirSQLite`); un chequeo `SELECT` previo (carrera entre procesos).
- **Efecto aceptado**: el id 25 hace avanzar la secuencia `AUTOINCREMENT` de `proyectos`.

## R3. ¿Cómo se carga `En progreso` y `Completada`?
- **Decisión**: se escriben directo (con `fecha_finalizacion` para Completada). El dominio solo expone `Nueva` y `Completada`, y la transición es de HU-11; `En progreso` se declara como constante local de `internal/demo`.
- **Por qué**: es el estado final deseado; el CHECK del esquema lo admite.

## R4. ¿Cómo se cargan las horas de Sprints cerrados?
- **Decisión**: `NewRegistroHoras` valida formato, actividad y fecha dentro del Sprint y hasta hoy; no se llama a `ValidarSprintParaImputar` (regla de la API de HU-07: solo Sprint activo), porque son datos históricos.

## R5. ¿Y el Sprint #3 Pendiente?
- **Decisión** (Scrum Master, 2026-10-30): se crea sin Sprint Goal ni fechas, como exige el CHECK (MODELO-DE-DATOS.md §4.4) y como lo crea `CrearSprint`.

## R6. `horas_por_story_point`
- **Decisión**: `Project` del dominio no tiene ese campo; la columna existe en el esquema (§4.1). Se escribe con el `INSERT` del bloque (valor 4), sin tocar el dominio.

## R7. Dónde monta el botón
- **Decisión**: ver *Structure Decision* del plan. Riesgo de conflicto con #148: bajo (una línea en `crear.tsx`; #148 compone en `rutas.tsx`).
