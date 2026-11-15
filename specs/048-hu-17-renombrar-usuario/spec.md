# Feature Specification: Renombrar el propio usuario (HU-17)

**Feature Branch**: `feature/HU-17-renombrar-usuario`

**Issue**: #176 (HU-17 Gestión de Usuarios)

**Created**: 2026-11-15

**Status**: Implementada (PR #325, pendiente de review)

**Dependencia**: HU-20 (autenticación, #233) y HU-21 (permisos por rol, #234), ya en `main`.
No hay cambio de esquema: `usuarios` ya existe y las demás tablas referencian `usuario_id`.

**Input**: User description: "HU-17: Gestión de Usuarios. Después de la redefinición acordada
con el Scrum Master (2026-11-15), la historia queda solo en el renombrado del propio usuario.
Cada persona cambia únicamente su nombre de usuario y conserva su historial. Los escenarios
de alta de usuarios, alta de integrantes, elección del Scrum Master y usuarios de la demo se
retiran porque ya los cubren HU-20, HU-21 y #235."

## Clarifications

### Session 2026-11-15

- Q: ¿Quién puede renombrar a un usuario? → A: **Opción A**: cada persona solo cambia su
  propio nombre, desde su sesión. Nadie cambia el dato de ingreso de otra persona. El
  endpoint es `PUT /auth/yo`, con el usuario tomado de la sesión.
- Q: ¿Qué pasa con los escenarios 1, 2, 3, 4 y 7 de la historia original? → A: Se retiran:
  el alta y el nombre repetido los cubre HU-20; elegir integrantes y roles los cubre HU-21;
  elegir al Scrum Master contradice HU-21 (lo es quien crea el proyecto) y los usuarios de
  la demo los cubre #235.
- Q: ¿Cambia el esquema? → A: No. `usuarios(id, nombre, nombre_normalizado, contrasena_hash,
  creado_en)` ya existe. El renombrado actualiza `nombre` y `nombre_normalizado` en una sola
  sentencia; votos, horas, integrantes y sesiones apuntan por `usuario_id`.

## Objetivo

Que una persona autenticada pueda corregir el nombre de usuario con el que ingresó, sin
perder su historial (votos, horas, equipos y sesiones) ni cerrar su sesión actual.

## Entradas y salidas

- **Entrada**: `PUT /auth/yo` con `{ "nombre": "nuevo.nombre" }` y la sesión vigente.
- **Salida exitosa**: `200` con `{ "id": <id>, "nombre": "nuevo.nombre" }`. La sesión sigue
  vigente y `GET /auth/yo` devuelve el nombre nuevo. El nombre nuevo es el dato de ingreso.
- **Errores**: `401` sin sesión; `400` con `campo: "nombre"` si el nombre no cumple las
  reglas de HU-20; `409` con «Ese nombre de usuario ya existe» si choca con otro usuario.
- **Historial**: el renombrado no cambia ninguna otra tabla; todo lo asociado por
  `usuario_id` sigue mostrando el nombre nuevo en las lecturas.

## Fuera de alcance

- Crear usuarios: lo cubre HU-20 (`POST /auth/registro`).
- Crear, editar o eliminar otros usuarios; cambiar roles de proyecto (HU-21/HU-16).
- Cambiar la contraseña, recuperarla o administrar cuentas ajenas.
- Elegir al Scrum Master al crear un proyecto: lo es quien crea el proyecto (HU-21).
- Volver a cargar los usuarios predeterminados de la demo: lo cubre #235.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Cambiar el propio nombre de usuario (Priority: P1)

Una persona con sesión inicia el cambio de su nombre desde el menú de la cuenta. Si el
nombre es válido y no está en uso, el sistema lo guarda y la persona sigue con la misma
sesión y su historial.

**Why this priority**: es el único comportamiento que queda de HU-17 y resuelve el error de
tipeo sin perder los aportes.

**Independent Test**: con una cuenta y una sesión, cambiar el nombre; comprobar el `200`, que
`GET /auth/yo` devuelve el nombre nuevo y que el ingreso con el nombre viejo ya no funciona y
con el nuevo sí.

**Acceptance Scenarios**:

1. **Given** una cuenta `ana` con sesión y una contraseña válida, **When** cambia su nombre a
   `ana.perez`, **Then** el sistema responde `200` con el nombre nuevo, la sesión sigue
   vigente y `GET /auth/yo` devuelve `ana.perez`.
2. **Given** la misma persona después del cambio, **When** ingresa con `ana.perez` y su
   contraseña, **Then** ingresa correctamente; **When** intenta con `ana`, **Then** recibe
   «Usuario o contraseña incorrectos».
3. **Given** un usuario que integra un proyecto, **When** cambia su nombre, **Then** el
   equipo, sus votos, sus horas y sus defectos siguen asociados a la misma persona y se
   muestran con el nombre nuevo.
4. **Given** un usuario con una sesión y un nombre nuevo inválido (menos de 3 caracteres, más
   de 30, con espacios, `@` o tildes), **When** lo envía, **Then** recibe `400` con
   `campo: "nombre"` y su nombre no cambia.

---

### User Story 2 - Rechazar un nombre ya en uso (Priority: P1)

Si otra cuenta ya usa el nombre, sin distinguir mayúsculas, el cambio se rechaza y la
persona conserva su nombre anterior.

**Why this priority**: la unicidad es la identidad del usuario; permitir duplicados rompería
el ingreso y la asignación de integrantes.

**Independent Test**: con dos cuentas, intentar que una tome el nombre de la otra y
comprobar el `409` y que el nombre anterior sigue vigente.

**Acceptance Scenarios**:

1. **Given** que existe el usuario `bruno` y la cuenta `ana` tiene sesión, **When** `ana`
   intenta cambiar su nombre a `BRUNO` o `bruno`, **Then** el sistema responde `409` con
   «Ese nombre de usuario ya existe» y `GET /auth/yo` sigue devolviendo `ana`.
2. **Given** dos pedidos simultáneos que intentan tomar el mismo nombre, **When** se
   procesan, **Then** solo uno cambia el nombre y el otro recibe `409`.

---

### User Story 3 - Sesión obligatoria (Priority: P2)

El cambio de nombre solo lo puede hacer la persona autenticada sobre su propia cuenta.

**Why this priority**: sin sesión no hay a quién renombrar; protege el dato de ingreso.

**Acceptance Scenarios**:

1. **Given** una petición sin token, con token inválido o vencido, **When** llama a
   `PUT /auth/yo`, **Then** recibe `401` y no se modifica ningún usuario.
2. **Given** una sesión vigente, **When** se envía el cambio, **Then** el handler usa el
   usuario de la sesión y no acepta un identificador de usuario en el cuerpo.

### Edge Cases

- Nombre con espacios externos: se rechaza con `400` por caracteres no permitidos, igual que
  en el registro; no se recorta en silencio. La comparación de duplicados sí es sin
  distinguir mayúsculas porque usa `nombre_normalizado`.
- Nombre de exactamente 3 o 30 caracteres: válido.
- Nombre igual al actual, con las mismas mayúsculas: la actualización es válida y no cambia
  nada observable.
- Renombrar no invalida la sesión actual ni las demás sesiones del usuario: las sesiones
  referencian `usuario_id`, no el nombre.
- Un usuario sin contraseña (previo a HU-20) puede ser renombrado si tiene una sesión; su
  nombre normalizado sigue siendo único.
- Un proyecto cerrado no impide cambiar el nombre del usuario: el renombrado no es una
  escritura sobre el proyecto.

## Requirements *(mandatory)*

### Reglas de negocio (EARS)

**Ubicuas**

- **FR-001**: El sistema SHALL permitir que una persona autenticada cambie únicamente su
  propio `usuarios.nombre`; el usuario SHALL tomarse de la sesión.
- **FR-002**: El nombre SHALL validarse con las mismas reglas de HU-20: de 3 a 30
  caracteres, solo letras ASCII, números, punto, guion y guion bajo.
- **FR-003**: El sistema SHALL guardar `nombre` con el valor recibido (sin recortar) y
  `nombre_normalizado` como `strings.ToLower(strings.TrimSpace(nombre))`.
- **FR-004**: La unicidad SHALL ser sin distinguir mayúsculas, apoyada en
  `usuarios.nombre_normalizado UNIQUE`.
- **FR-005**: El renombrado SHALL ser atómico: `nombre` y `nombre_normalizado` se actualizan
  en una sola sentencia `UPDATE`.
- **FR-006**: El renombrado SHALL conservar el `usuario_id`; votos, horas, integrantes,
  defectos y sesiones SHALL seguir asociados a la misma persona.
- **FR-007**: La sesión actual y las demás sesiones del usuario SHALL seguir vigentes después
  del renombrado.

**Basadas en evento**

- **FR-008**: WHEN una persona autenticada envía un nombre válido y no usado, el sistema
  SHALL actualizar el usuario y responder `200` con `{id, nombre}`.
- **FR-009**: WHEN el nombre normalizado choca con otro usuario, el sistema SHALL responder
  `409` «Ese nombre de usuario ya existe» sin modificar el usuario.
- **FR-010**: WHEN el nombre no cumple las reglas, el sistema SHALL responder `400` con
  `campo: "nombre"` sin modificar el usuario.

**Basadas en estado**

- **WHILE** la persona tiene una sesión vigente, el sistema SHALL permitir el cambio solo
  sobre el usuario de esa sesión.
- **WHILE** el nombre nuevo coincide con el actual, el sistema SHALL aceptar la operación.

**No deseadas**

- **IF** la petición no tiene sesión válida, THEN el sistema SHALL responder `401` sin
  modificar nada.
- **IF** el cuerpo no es JSON o falta `nombre`, THEN el sistema SHALL responder `400`.
- **IF** el usuario de la sesión no existe (caso no alcanzable desde la API), THEN el sistema
  SHALL responder `404` sin modificar nada.
- **IF** dos renombrados simultáneos toman el mismo nombre, THEN solo uno SHALL aplicarse y
  el otro SHALL responder `409`.

### Requisitos de la interfaz *(frontend)*

- **FR-011**: El menú de la cuenta MUST ofrecer «Cambiar nombre» cuando hay sesión.
- **FR-012**: El formulario MUST mostrar el nombre actual y permitir guardar o cancelar; al
  guardar, MUST llamar a `PUT /auth/yo` y reflejar el nombre nuevo en el menú.
- **FR-013**: Si el servidor rechaza el cambio, la interfaz MUST mostrar el mensaje (junto al
  campo si trae `campo: "nombre"`) y conservar lo escrito; en `409`, MUST mostrar «Ese nombre
  de usuario ya existe» y mantener el nombre anterior.
- **FR-014**: Cancelar MUST cerrar el formulario sin llamar a la API.

### Key Entities

- **Usuario**: `id`, `nombre`, `nombre_normalizado` (único), `contrasena_hash` opcional,
  `creado_en`. El renombrado cambia solo `nombre` y `nombre_normalizado`.
- **Sesión**: referencia `usuario_id`; no guarda el nombre. Sigue vigente después del
  renombrado y `GET /auth/yo` lee el nombre actual del usuario.

## Success Criteria *(mandatory)*

- **SC-001**: El 100 % de los renombrados válidos y no duplicados responde `200`, deja el
  nombre nuevo en `GET /auth/yo` y mantiene la sesión vigente.
- **SC-002**: El 100 % de los intentos con un nombre en uso responde `409` y deja el nombre
  anterior intacto.
- **SC-003**: El 100 % de los nombres inválidos responde `400` con `campo: "nombre"` y no
  modifica al usuario.
- **SC-004**: Después del renombrado, el equipo, los votos, las horas y los defectos siguen
  asociados al mismo `usuario_id` y muestran el nombre nuevo.
- **SC-005**: Los escenarios de Godog y las pruebas de componente de esta historia pasan en
  el CI, junto con los existentes.

## Assumptions

- Se elige la opción A de la redefinición: cada persona solo su propio nombre.
- El endpoint es `PUT /auth/yo` (consistente con los `PUT` de edición existentes) y toma el
  usuario de la sesión; no recibe `usuario_id`.
- El nombre de usuario es el dato de ingreso; después del cambio, el nombre anterior deja de
  servir para ingresar y el nuevo funciona con la misma contraseña.
- No hay migración ni cambio de esquema.
- El frontend usa el contexto de sesión existente y actualiza el usuario en memoria con la
  respuesta del `PUT`.
