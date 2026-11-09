# Feature Specification: Criterios de aceptación al crear una historia

**Feature Branch**: `feature/174-criterios-aceptacion`

**Created**: 2026-11-08

**Status**: Draft

**Input**: User description: "Issue #174: poder cargar los criterios de aceptación en el formulario de creación de una historia (HU-01), igual que en la edición (HU-15), y mostrarlos en el resultado de la creación y en el listado. Texto libre opcional de hasta 2000 caracteres. Alcance acordado: solo criterios de aceptación; notas y depende_de quedan fuera de este incremento (D17)."

## Clarifications

### Session 2026-11-08

- Q: ¿Qué alcance entra en este incremento? → A: Solo criterios de aceptación. Notas y depende_de (D17) no entran: la opción A (dos columnas nuevas, `depende_de` texto libre) quedó decidida para una historia futura y requiere migración; no se implementa acá.
- Q: ¿Los criterios son obligatorios al crear? → A: No. Son opcionales; vacío o solo espacios equivale a no tener criterios (se guarda `NULL`).
- Q: ¿Qué límite tienen? → A: 2000 caracteres, el mismo que la descripción y que la edición de HU-15.
- Q: ¿Cómo se muestran? → A: Con el mismo formato Markdown que el listado del backlog (HU-15/#147). Sin criterios se muestra el texto "Sin criterios de aceptación".
- Q: ¿Cambia el esquema? → A: No. `historias.criterios_aceptacion` ya existe desde la migración 001; no hay tablas, columnas ni migraciones nuevas.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Cargar los criterios de aceptación al crear la historia (Priority: P1)

Como integrante del equipo, al crear una historia de usuario puedo cargar sus criterios de aceptación en el mismo formulario, sin tener que guardarla primero y editarla después, para dejar la historia completa desde el momento en que la doy de alta.

**Why this priority**: Es el objetivo del issue #174 y completa el contenido mínimo de la historia (título, descripción, prioridad y criterios) en un solo paso. Sin esto, los criterios recién se pueden cargar desde la edición.

**Independent Test**: Se crea una historia enviando criterios válidos y se verifica que la respuesta de creación los devuelve y que quedan persistidos; se repite sin criterios y con criterios demasiado largos.

**Acceptance Scenarios**:

1. **Given** el formulario de una nueva historia, **When** ingreso criterios de aceptación de hasta 2000 caracteres y guardo, **Then** la historia se crea con esos criterios y la respuesta de creación los devuelve.
2. **Given** el formulario de una nueva historia, **When** dejo los criterios vacíos o solo con espacios y guardo, **Then** la historia se crea sin criterios (`null`) y no se muestra ningún error por ese campo.
3. **Given** el formulario de una nueva historia, **When** ingreso criterios de más de 2000 caracteres y guardo, **Then** el sistema rechaza la creación con un error de validación sobre el campo `criterios_aceptacion` y no registra la historia.
4. **Given** criterios con espacios al principio o al final, **When** creo la historia, **Then** se guardan recortados.

---

### User Story 2 - Ver los criterios en la confirmación y en el listado (Priority: P2)

Como integrante del equipo, después de crear una historia veo sus criterios de aceptación en el resultado de la creación y en el listado del Product Backlog, para confirmar que se guardaron bien y consultarlos sin abrir la edición.

**Why this priority**: Cierra el circuito visual del issue: hoy la confirmación de creación no muestra los criterios y el listado ya los muestra desde HU-15.

**Independent Test**: Se crea una historia con criterios y se verifica que aparecen en el panel "Historia creada" y en el listado; se crea otra sin criterios y se verifica el texto "Sin criterios de aceptación".

**Acceptance Scenarios**:

1. **Given** una historia recién creada con criterios, **When** se muestra el resultado de la creación, **Then** los criterios aparecen con formato Markdown.
2. **Given** una historia recién creada sin criterios, **When** se muestra el resultado de la creación, **Then** aparece "Sin criterios de aceptación".
3. **Given** el listado del Product Backlog de un proyecto, **When** una historia tiene criterios, **Then** se muestran con formato Markdown; si no tiene, se indica que no hay criterios (comportamiento ya existente de HU-15 que no debe romperse).

---

### Edge Cases

Redactados en formato EARS.

- **WHILE** los criterios de aceptación son opcionales, el sistema **shall** permitir crear la historia sin informarlos.
- **IF** los criterios están vacíos o contienen solo espacios, **THEN** el sistema **shall** guardarlos como ausentes (`NULL`) y **shall** no mostrar un error.
- **IF** los criterios superan los 2000 caracteres contando espacios externos, **THEN** el sistema **shall** rechazar la creación con un error sobre el campo `criterios_aceptacion` y **shall** no registrar la historia.
- **IF** el proyecto no existe, **THEN** la creación **shall** seguir respondiendo "proyecto no encontrado" y **shall** no registrar la historia, con o sin criterios.
- **IF** la creación falla por cualquier motivo, **THEN** el sistema **shall** no mostrar datos parciales ni una historia creada a medias.
- **WHEN** los criterios contienen Markdown, **THEN** el sistema **shall** mostrarlos con el mismo renderizado que el listado del backlog, sin interpretarlos como HTML inseguro.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: El sistema MUST aceptar, en la creación de una historia, un campo opcional `criterios_aceptacion` de texto libre.
- **FR-002**: El sistema MUST recortar los espacios externos de los criterios y MUST tratar el resultado vacío como ausencia de criterios (`NULL`).
- **FR-003**: El sistema MUST rechazar con un error de validación sobre `criterios_aceptacion` los criterios que superen los 2000 caracteres, sin registrar la historia.
- **FR-004**: El sistema MUST persistir los criterios junto con la historia en `historias.criterios_aceptacion` y MUST devolverlos en la respuesta de creación.
- **FR-005**: El sistema MUST omitir el campo `criterios_aceptacion` o devolverlo como `null` cuando la historia se cree sin criterios.
- **FR-006**: La pantalla de creación MUST ofrecer un campo opcional "Criterios de aceptación" y MUST validar el límite de 2000 caracteres antes de enviar.
- **FR-007**: El resultado de la creación MUST mostrar los criterios de la historia recién creada, o "Sin criterios de aceptación" si no tiene.
- **FR-008**: El listado del Product Backlog MUST seguir mostrando los criterios de cada historia como hasta ahora.
- **FR-009**: Este incremento MUST NOT agregar columnas, tablas ni migraciones: la columna ya existe en el esquema vigente.
- **FR-010**: Este incremento MUST NOT implementar notas ni "depende de": quedan fuera de alcance (D17 decidida para una historia futura).
- **FR-011**: El sistema MUST mantener sin cambios las reglas de creación ya existentes (título, descripción, prioridad, valor de negocio, estado "Nueva", proyecto existente).

### Key Entities *(include if feature involves data)*

- **Historia**: la de HU-01/HU-15. Suma el uso de `criterios_aceptacion`, texto libre opcional de hasta 2000 caracteres que ya existe en el modelo; no se agregan atributos nuevos.
- **Product Backlog**: listado de historias del proyecto; ya muestra los criterios y debe conservar ese comportamiento.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: El 100% de las historias creadas con criterios válidos los devuelve la respuesta de creación y los conserva al consultar el listado.
- **SC-002**: El 100% de las creaciones con criterios de más de 2000 caracteres es rechazado sin registrar la historia.
- **SC-003**: El 100% de las creaciones sin criterios (campo omitido, vacío o solo espacios) se registra con criterios ausentes y sin error.
- **SC-004**: Una persona puede crear una historia completa con criterios sin pasar por la pantalla de edición.
- **SC-005**: El cambio no introduce tablas, columnas ni migraciones nuevas.

## Assumptions

- `historias.criterios_aceptacion` existe desde la migración 001 y admite `NULL`; no hace falta cambio de esquema.
- Los criterios son texto libre con Markdown, igual que en HU-15; el renderizado ya existe (#147).
- El límite de 2000 caracteres se cuenta sobre el texto recortado, igual que en la edición.
- Los criterios no reemplazan ni modifican la descripción de la historia.
- No cambian los permisos ni la autenticación: el sistema sigue sin control de acceso por rol en este incremento.
- Notas y depende_de (D17) quedan fuera; cuando se implementen, su decisión ya tomada es la opción A (columnas nuevas, `depende_de` texto libre), en un PR con migración y `MODELO-DE-DATOS.md`.
