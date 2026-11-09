# Data Model: Criterios de aceptación al crear una historia

**No hay cambio de esquema.** La columna `historias.criterios_aceptacion` (`TEXT?`) existe desde la migración 001 y está documentada en [`MODELO-DE-DATOS.md`](../../MODELO-DE-DATOS.md) §4.4. Este documento enlaza a esa sección y no redefine tablas.

## Entidad existente

**Historia** (`historias`)

| Columna | Tipo | Uso en esta feature |
|---|---|---|
| `criterios_aceptacion` | `TEXT?` | Criterios de aceptación opcionales, texto libre. `NULL` = sin criterios. Máximo 2000 caracteres después del trim. |
| `descripcion` | `TEXT NOT NULL` | No cambia; sigue siendo obligatoria. |
| `titulo` | `TEXT NOT NULL` | No cambia. |
| `prioridad` | `TEXT NOT NULL` | No cambia. |
| `estado` | `TEXT NOT NULL` | La creación sigue forzando `Nueva`. |

## Reglas de normalización y validación

| Regla | Dónde vive |
|---|---|
| Trimming de espacios externos | `internal/domain` (helper compartido con la edición) |
| Vacío o solo espacios → `NULL` | `internal/domain` (`ConCriteriosAceptacion`) |
| Máximo 2000 caracteres (runas) | `internal/domain` (`ConCriteriosAceptacion`); error `ValidationError{Campo: "criterios_aceptacion"}` |
| El campo es opcional en el alta | DTO y service (puntero `*string`) |
| La respuesta devuelve el valor persistido o `null` | `historiaResponse` existente |

## Formas de los datos

**Entrada `POST /backlog`** (campo nuevo, opcional):

```json
{
  "proyecto_id": 1,
  "titulo": "Iniciar sesión",
  "descripcion": "Autenticación con usuario y contraseña",
  "prioridad": "M",
  "valor_negocio": 8,
  "criterios_aceptacion": "Dado un usuario registrado, cuando ingresa, entonces accede."
}
```

`criterios_aceptacion` puede omitirse o ser `null`. Si es un string vacío o solo espacios, se guarda `NULL`.

**Respuesta 201** (misma forma que el listado, campo ya existente):

```json
{
  "id": 12,
  "proyecto_id": 1,
  "titulo": "Iniciar sesión",
  "descripcion": "Autenticación con usuario y contraseña",
  "criterios_aceptacion": "Dado un usuario registrado, cuando ingresa, entonces accede.",
  "prioridad": "M",
  "estado": "Nueva",
  "valor_negocio": 8,
  "estimacion_sp": null,
  "sprint_id": null,
  "fecha_finalizacion": null
}
```

Sin criterios, `criterios_aceptacion` es `null`.

## Tipos del frontend

- `FormularioHistoria`: suma `criterios: string` (estado inicial `''`).
- `CuerpoCrearHistoria`: suma `criterios_aceptacion?: string`.
- `HistoriaCreada`: suma `criterios_aceptacion: string | null`.

## Estados de la pantalla

| Estado | Condición | Qué se muestra |
|---|---|---|
| Formulario | Sin envío | Campo opcional de criterios; sin error si está vacío |
| Error de validación | Criterios de más de 2000 | Mensaje bajo el campo; no se envía |
| Creada con criterios | Respuesta 201 con string | Confirmación con los criterios renderizados en Markdown |
| Creada sin criterios | Respuesta 201 con `null` | Confirmación con "Sin criterios de aceptación" |
| Error de API | 4xx/5xx o conexión | Mensaje de error; no se muestra una historia parcial |
