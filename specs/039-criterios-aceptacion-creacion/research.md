# Research: Criterios de aceptación al crear una historia

## R1 — Reutilizar `historias.criterios_aceptacion`

- **Decision**: No se agrega esquema. La columna `historias.criterios_aceptacion` (`TEXT?`) existe desde la migración 001 y el repositorio SQLite ya la escribe en `Guardar` y la lee en `listarHistorias`.
- **Rationale**: Es exactamente el dato que pide el issue; agregar otra columna o una tabla duplicaría información. `MODELO-DE-DATOS.md` §4.4 ya define el campo.
- **Alternatives considered**: Crear una tabla de criterios (innecesaria; el modelo decidió texto libre en una sola columna, D10); guardarlos dentro de la descripción (rompe la edición y el listado actuales).

## R2 — Validar los criterios en el dominio sin cambiar la fábrica

- **Decision**: se agrega `ConCriteriosAceptacion(criterios *string) (BacklogItem, error)` a la entidad, con la misma forma que `ConEstimacion`. La normalización (trim, vacío → `nil`, máximo 2000) se extrae a un helper compartido y lo usan tanto el alta como `ConDatosEditados`. El service encadena `NewBacklogItem` → `ConCriteriosAceptacion` antes de verificar el proyecto y persistir.
- **Rationale**: la validación sigue en el dominio (Constitución IV/V) y no cambia la firma de `NewBacklogItem`, así no hay que tocar los llamados existentes de producción y tests. Compartir el helper evita que el alta y la edición acepten reglas distintas (DRY).
- **Alternatives considered**: agregar un parámetro a `NewBacklogItem` (obliga a actualizar todos los llamados y no aporta nada si los criterios son opcionales); validar en el service (deja la entidad sin la regla); crear una segunda fábrica (dos caminos de creación).

## R3 — Contrato HTTP

- **Decision**: `POST /backlog` acepta `criterios_aceptacion` como string opcional o `null`. El campo viaja también en la respuesta 201, que ya usa `historiaResponse`. El error de longitud es 400 con `campo: "criterios_aceptacion"`.
- **Rationale**: es el endpoint real del alta (el issue lo menciona como `POST /projects/{id}/backlog`, pero el contrato vigente tiene `proyecto_id` en el cuerpo y la ruta es `POST /backlog`). Mantener el contrato actual evita romper el frontend y los tests de HU-01.
- **Alternatives considered**: agregar una ruta nueva por proyecto (duplicaría el endpoint sin necesidad); devolver 422 (la API usa 400 para validaciones).

## R4 — Formulario y confirmación

- **Decision**: el formulario de creación suma un `<textarea>` opcional "Criterios de aceptación" con validación de 2000 caracteres; el cuerpo lo envía solo si, después del trim, no está vacío. La confirmación muestra los criterios con `Markdown` (el mismo componente del listado) o "Sin criterios de aceptación".
- **Rationale**: replica la edición de HU-15 y el renderizado de #147; el listado ya muestra el campo. Enviar el campo solo cuando tiene contenido mantiene el cuerpo mínimo y el backend interpreta la ausencia como `NULL`.
- **Alternatives considered**: usar un editor Markdown (fuera de alcance y sin dependencias nuevas); mostrar texto plano en la confirmación (inconsistente con el listado).

## R5 — Pruebas

- **Decision**: BDD con Godog sobre el service (éxito con criterios, sin criterios y criterios demasiado largos) y pruebas de handler para el HTTP (201 con el campo, 201 con `null`, 400 con `campo`). En el frontend, pruebas de componente con el escenario en el nombre y pruebas unitarias de las funciones puras.
- **Rationale**: Godog ya cubre HU-01 a nivel de service y el handler se prueba con `httptest`; el frontend no es ejercitable con Godog (Constitución III).
- **Alternatives considered**: probar el alta solo desde el handler (deja el service sin cobertura directa); omitir BDD por ser un cambio chico (el proyecto exige escenarios para el comportamiento observable).

## R6 — Notas y depende_de

- **Decision**: fuera de alcance. D17 quedó decidida por el Scrum Master como opción A (dos columnas nuevas, `depende_de` texto libre) para una historia futura; este PR no toca `MODELO-DE-DATOS.md` ni agrega migraciones.
- **Rationale**: la planificación del 2026-11-03 acotó este issue a criterios de aceptación y una migración requiere coordinar el número con #236.
- **Alternatives considered**: implementar D17 ahora (mayor alcance, migración y decisión de UI no pedidas en este incremento).
