# Contrato HTTP: /demo (cambios de la v2)

Rutas y habilitación: sin cambios respecto de [`021/contracts/api.md`](../../021-datos-de-prueba/contracts/api.md).

## POST /demo/datos

| Caso | Estado | Cuerpo |
|---|---|---|
| Carga exitosa | 201 | `{"proyecto_id":25,"usuarios":4,"integrantes":4,"sprints":4,"historias":12,"registros_horas":20,"rondas_votacion":2,"votos":6,"defectos":5}` |
| Ya cargado | 409 | `{"mensaje":"los datos de prueba ya están cargados","campo":"ya_cargado"}` (sin cambios) |
| **Usuario de la demo existente** (**nuevo**) | 409 | `{"mensaje":"ya existe el usuario ana.iriarte","campo":"usuario_existente"}` |
| Id 25 de otro proyecto | 409 | `{"mensaje":"el id 25 ya lo usa otro proyecto","campo":"proyecto_id"}` (sin cambios) |
| Error interno | 500 | `{"mensaje":"no se pudieron cargar los datos de prueba"}` |

Prioridad cuando coinciden: `ya_cargado` > `usuario_existente` > `proyecto_id`.

## Frontend (`/ingresar`)
- `GET /demo/estado` 200 ⇒ botón «Cargar datos de prueba»; otra cosa ⇒ sin botón ni tabla.
- 201 o 409 `ya_cargado` ⇒ aviso + tabla (usuario, contraseña, rol) + enlace al proyecto de prueba.
- 409 `usuario_existente` ⇒ aviso con el mensaje del servidor, sin tabla.
- Otro error (incluido `proyecto_id`) ⇒ mensaje de error.
