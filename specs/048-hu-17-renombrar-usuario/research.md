# Research: Renombrar el propio usuario (HU-17)

## R1. Alcance después de la redefinición

- **Decisión**: la historia queda solo en el renombrado del propio usuario (escenarios 5 y 6
  de la redacción original). Se retiran los escenarios 1, 2, 3, 4 y 7 porque ya los cubren
  HU-20, HU-21 y #235, o contradicen HU-21.
- **Rationale**: con el login, cada persona se registra sola; el nombre es su dato de ingreso.
  Los roles se asignan por proyecto y el Scrum Master es quien crea el proyecto.
- **Alternativas**: mantener el alta/edición de usuarios por el Scrum Master (contradice
  HU-20 y obliga a definir cuál Scrum Master puede renombrar a una persona que está en varios
  proyectos).

## R2. Endpoint y autorización (opción A)

- **Decisión**: `PUT /auth/yo`, nivel `SoloSesion`. El handler toma el `domain.User` que dejó
  `ConSesion` en el contexto y no acepta `usuario_id` en el cuerpo.
- **Rationale**: cada persona solo cambia su propio nombre; no hay que definir permisos por
  proyecto para un dato global del usuario. `PUT` reemplaza el nombre y el cliente ya lo
  soporta.
- **Alternativas**: `PATCH /auth/yo` (obligaría a sumar un método al cliente compartido sin
  ganancia); una ruta de administración `PUT /users/{id}` (permitiría cambiar el ingreso de
  otra persona y requiere reglas de rol que el issue no pide).

## R3. Actualización atómica y unicidad

- **Decisión**: `UPDATE usuarios SET nombre = ?, nombre_normalizado = ? WHERE id = ?` en una
  sola sentencia. El índice `usuarios.nombre_normalizado UNIQUE` cubre la carrera; la
  violación `SQLITE_CONSTRAINT_UNIQUE` se traduce a `ErrUsuarioDuplicado`, que el servicio
  convierte en `domain.ErrNombreUsuarioEnUso` (409).
- **Rationale**: el nombre y su identidad normalizada no pueden quedar desincronizados; no
  hace falta consultar antes ni una transacción explícita.
- **Alternativas**: `SELECT` + `UPDATE` en dos pasos (ventana de carrera); transacción en el
  servicio (el servicio no conoce transacciones).

## R4. Validación del nombre

- **Decisión**: reutilizar `domain.ValidarNombreUsuario` (3 a 30 caracteres, ASCII, `.`, `-`,
  `_`) y `domain.NormalizarNombre` (minúsculas y sin espacios externos). El nombre se guarda
  tal como se recibió, igual que en el registro; los espacios externos se rechazan.
- **Rationale**: el nombre de usuario es el dato de ingreso y debe tener las mismas reglas
  que el registro; la unicidad es case-insensitive por `nombre_normalizado`.
- **Alternativas**: reglas propias para el renombrado (permitiría nombres que el registro
  rechaza); recortar en silencio (cambiaría el nombre mostrado sin avisar).

## R5. Sesiones e historial

- **Decisión**: no tocar `sesiones`, `votos`, `registros_horas`, `integrantes` ni `defectos`.
  Todas referencian `usuario_id`; las lecturas que muestran el nombre hacen JOIN con
  `usuarios` y ven el nombre nuevo. La sesión actual sigue vigente.
- **Rationale**: el renombrado es una actualización de un atributo del usuario, no una
  migración de identidad. No se cierra la sesión de nadie ni se pierde historial.
- **Alternativas**: reemitir el token (innecesario: el token no guarda el nombre); migrar
  filas (rompería la normalización y no hace falta).

## R6. Frontend

- **Decisión**: agregar `cambiarNombre` al contexto de sesión y un formulario en
  `MenuUsuario` que se abre con «Cambiar nombre». Al guardar, llama a `PUT /auth/yo` y
  actualiza `usuario` con la respuesta; el menú muestra el nombre nuevo. Los errores del
  servidor se muestran en el formulario; en `409` se conserva el nombre anterior.
- **Rationale**: el contexto ya concentra la identidad y el token; el menú es el lugar que
  pide el issue. No hace falta una ruta nueva.
- **Alternativas**: una pantalla `/cuenta` (más superficie para un formulario chico); manejar
  el `PUT` directamente en el menú sin el contexto (el nombre del menú quedaría
  desactualizado).

## R7. Documentación

- **Decisión**: actualizar `MODELO-DE-DATOS.md` §4.2 (regla de renombrado), §5 (escritura
  atómica) y §8 (fila HU-17), sin migración; redefinir la sección HU-17 de
  `docs/PRODUCT-BACKLOG.md`.
- **Rationale**: el modelo no cambia, pero la regla y la trazabilidad de la historia deben
  quedar documentadas. El issue pide actualizar el backlog en el mismo PR.
