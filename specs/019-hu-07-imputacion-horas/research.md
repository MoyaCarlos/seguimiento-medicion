# Research: Imputación de Horas Trabajadas (HU-07)

**Fecha**: 2026-10-28 | **Spec**: [spec.md](./spec.md)

No quedaron marcas NEEDS CLARIFICATION en el Technical Context: las seis dudas de negocio
se resolvieron en `/speckit.specify` y `/speckit.clarify` (ver Clarifications de la spec).
Este documento registra las decisiones técnicas tomadas para el diseño.

## R1. Rutas de la API

- **Decision**: `POST /backlog/{id}/horas` para imputar y `GET /backlog/{id}/horas` para
  consultar el esfuerzo real, donde `{id}` es el identificador de la historia.
- **Rationale**: la carga es un recurso hijo de la historia (`historias ||--o{ registros_horas`),
  y `/backlog/{id}` ya es la ruta de la historia (`PUT /backlog/{id}`, HU-15). El proyecto y
  el Sprint no van en la URL: se deducen de la historia (el Sprint es el activo, FR-003).
- **Alternatives considered**: `POST /horas` con `historia_id` en el cuerpo (pierde el 404
  natural de la historia en la ruta y obliga a validar el cuerpo antes de saber si existe);
  `/sprints/{id}/horas` (el Sprint no lo elige quien carga: es el de la historia).

## R2. Duración: de horas decimales a minutos enteros sin `float64`

- **Decision**: el dominio recibe la cantidad de horas como **texto** y la convierte con
  `math/big.Rat`: valida el formato decimal con `^-?\d+(\.\d+)?$`, multiplica por 60 y exige
  que el resultado sea entero, mayor que 0 y menor o igual que 1440 (24 h). Función
  `domain.MinutosDesdeHoras(texto string) (int, error)`.
- **Rationale**: D8 eligió minutos enteros justamente para evitar el error de `float64`
  (`0.1 * 60` no da `6` exacto en binario). `big.Rat` es de la biblioteca estándar y
  representa `0.1` como `1/10`, así que `1,5 h → 90` y `0,01 h → 0,6` (se rechaza) son
  exactos. Recibir texto permite además que "texto no numérico" (Escenario 2) sea un error
  sobre `horas` evaluado en el orden de la spec, y no un "JSON inválido".
- **Alternatives considered**: `float64` y redondeo (acepta en silencio 0,01 h como 1 min o
  0 min; contradice la Clarification de precisión); recibir minutos en la API (traslada la
  conversión al frontend y la duplica, rompe DRY; la spec dice que se ingresan horas).

## R3. Formato de `horas` en el JSON

- **Decision**: el campo `horas` se lee como `json.RawMessage`. Si es un número JSON, se pasa
  su literal tal cual al dominio (`4`, `1.5`); si es un string JSON, se pasa su contenido
  (`"abc"` → `abc`, que el dominio rechaza); si falta o es `null`, se pasa vacío (el dominio
  responde "las horas son obligatorias"). No se usa `float64` en ningún paso.
- **Rationale**: con un campo `float64` o `json.Number`, `"abc"` rompe el `Unmarshal`
  completo y se perdería el error por campo. Con `RawMessage` el handler no valida reglas:
  solo extrae el texto (SRP; la regla vive en el dominio).
- **Alternatives considered**: `json.Number` (falla con strings no numéricos, igual que arriba).

## R4. Orden de las validaciones

- **Decision**: el orden de la spec (Edge Cases, FR-011) se reparte así:
  1. **Handler**: id de ruta inválido → 400; historia inexistente → 404, consultando la
     historia **antes** de leer el cuerpo (mismo patrón que `PUT /backlog/{id}`). Cuerpo que
     no es JSON → 400 "JSON inválido".
  2. **Service** `ImputarHoras`: historia sin Sprint, o con Sprint que no está `Activo` →
     `ErrHistoriaFueraDelSprintActivo`; integrante que no está en
     `ListarIntegrantes(proyecto)` → `ErrIntegranteAjeno`.
  3. **Dominio** `NewRegistroHoras`: horas → actividad → fecha (formato, después rango).
- **Rationale**: la existencia domina sobre la validación (HU-04, 05, 06, 15). Las reglas que
  necesitan leer otras entidades (Sprint, integrantes) están en el service; las que dependen
  solo de los datos de la carga y del Sprint ya leído, en el dominio. La fecha llega como
  texto al dominio por la misma razón que las horas: un formato inválido no debe adelantarse a
  la regla del Sprint activo.
- **Alternatives considered**: validar todo el cuerpo en el handler (`fechaDeCampo`) antes del
  service (rompe el orden fijado en la spec).

## R5. Carrera con el cierre del Sprint

- **Decision**: el adaptador inserta con un `INSERT ... SELECT ... WHERE EXISTS` que repite,
  en la misma sentencia, la condición "la historia sigue en ese Sprint y ese Sprint sigue
  `Activo`". Si no inserta ninguna fila, devuelve `domain.ErrHistoriaFueraDelSprintActivo`.
  ```sql
  INSERT INTO registros_horas (proyecto_id, historia_id, usuario_id, sprint_id, fecha, actividad, minutos, creado_en)
  SELECT ?, ?, ?, ?, ?, ?, ?, ?
  WHERE EXISTS (SELECT 1 FROM historias h JOIN sprints s ON s.id = h.sprint_id
                WHERE h.id = ? AND h.sprint_id = ? AND s.estado = 'Activo')
  ```
- **Rationale**: entre que el service lee el Sprint y que inserta, otro pedido puede cerrar el
  Sprint (`CerrarConArrastre` devuelve las no completadas al backlog y finaliza el Sprint).
  La sentencia única es atómica en SQLite, así que la carga queda en el Sprint correcto o se
  rechaza con el mismo mensaje que la validación previa (edge case de la spec). Es el mismo
  criterio que el `UPDATE ... WHERE estado = 'Activo'` del fix #145 para el cierre doble.
- **Alternatives considered**: transacción con `SELECT` previo (con una sola conexión
  funciona, pero es más código para el mismo resultado); confiar solo en el service (deja la
  carrera abierta).

## R6. Integrante del proyecto

- **Decision**: el service consulta `ListarIntegrantes(proyectoID)` (método que ya existe en
  `ProjectRepository`) a través de un puerto angosto `IntegrantesRepository`, y busca el
  `usuario_id`. Un usuario inexistente o de otro proyecto da el mismo
  `ValidationError{Campo: "integrante_id"}` (400), como pide la spec ("error sobre el
  integrante", no "no encontrado"). La FK compuesta `(proyecto_id, usuario_id) →
  integrantes` queda como red de seguridad (409, `MODELO-DE-DATOS.md` §5).
- **Rationale**: reutiliza un método probado; el volumen de integrantes es chico. ISP: el
  service no depende de todo `ProjectRepository`.
- **Alternatives considered**: un método nuevo `EsIntegrante` (YAGNI con el volumen actual).

## R7. "Hoy" y comparación de fechas

- **Decision**: el handler pasa `time.Now()` al service y el dominio lo reduce al día civil
  local a las 00:00 UTC, con el mismo helper `inicioDelDiaCivil` que usa HU-13. La fecha de
  la carga y las del Sprint ya están a las 00:00 UTC (`time.DateOnly`), así que se comparan
  por día. Rango aceptado: `inicio ≤ fecha ≤ min(hoy, fin)`.
- **Rationale**: misma convención que el estado del proyecto; los tests fijan "hoy" sin leer
  el reloj (Assumptions de la spec).
- **Alternatives considered**: `time.Now().UTC()` (en Argentina, después de las 21 h cambiaría
  de día antes de tiempo).

## R8. Total acumulado y detalle

- **Decision**: `ListarPorHistoria` devuelve las cargas con el nombre del integrante (`JOIN
  usuarios`), ordenadas por `fecha, id`. El total lo calcula el dominio sumando los minutos
  (`domain.TotalMinutos`). La respuesta expone `minutos` (exacto) y `horas`
  (`minutos / 60`, solo para mostrar) en cada carga y en el total.
- **Rationale**: ya se leen todas las cargas para el detalle, así que la suma en Go es exacta
  (enteros) y evita una segunda consulta. HU-03 sumará por Sprint con `SUM(minutos)` cuando
  lo necesite.
- **Alternatives considered**: `SUM(minutos)` aparte (dos consultas para el mismo dato);
  devolver solo horas (pierde la exactitud que D8 quiere preservar).

## R9. Errores y códigos HTTP

| Caso | Error | HTTP |
|---|---|---|
| Historia inexistente | `domain.ErrHistoriaNoEncontrada` | 404 |
| Historia fuera del Sprint activo (incluida la carrera de R5) | `ValidationError{Campo: "historia_id"}` | 400 |
| Integrante ajeno o inexistente | `ValidationError{Campo: "integrante_id"}` | 400 |
| Horas, actividad o fecha inválidas | `ValidationError{Campo: "horas" / "actividad" / "fecha"}` | 400 |
| FK violada (no debería pasar: el service valida antes) | `repository.ErrIntegridadReferencial` | 409 |
| `CHECK` violado (defecto: el dominio dejó pasar algo) | genérico | 500 (D13) |

- **Rationale**: es la tabla de `MODELO-DE-DATOS.md` §5. "Fuera del Sprint activo" es una
  regla de estado, como "solo se puede cerrar un Sprint activo", que el proyecto ya responde
  con `ValidationError` y 400. `escribirError` no cambia.

## R10. Migración 004 y tests de esquema

- **Decision**: `internal/repository/migraciones/004_registros_horas.sql` con la tabla de
  `MODELO-DE-DATOS.md` §4.8 tal cual (FK compuestas, `CHECK (typeof(minutos) = 'integer' AND
  minutos > 0)`, `CHECK (date(fecha) IS fecha)`). La 003 la tiene reservada HU-02; el runner
  aplica las que faltan aunque haya huecos (§9.1). Se agregan a `esquema_test.go` los casos
  mínimos de §9.6 para `registros_horas`.
- **Rationale**: el número 004 está reservado para HU-07 en §9.3; el documento no cambia
  (el esquema ya está definido), así que no hace falta tocar `MODELO-DE-DATOS.md`.
- **Alternatives considered**: ninguna; el esquema es la fuente de verdad.

## R11. BDD sin HU-06

- **Decision**: los pasos Godog siembran el Sprint `Activo` y asignan la historia con
  `UPDATE historias SET sprint_id = ?`, como ya hacen `steps_sprints.go` y
  `steps_backlog.go`. Los pasos llaman al service directamente, y las fechas del Sprint y de
  la carga se arman relativas al día actual con `hoyCivil()`, como en
  `steps_estado_proyecto.go` (HU-13). En los tests unitarios de dominio y de service, "hoy"
  se pasa fijo.
- **Rationale**: HU-06 todavía no está en `main` (Assumptions de la spec).
