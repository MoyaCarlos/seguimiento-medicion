# Contrato HTTP: /demo

Prefijo `/demo`, agregado a `PrefijosAPI` y al proxy de Vite. Con `DEMO_HABILITADO` distinto de `true`, **todas** las rutas responden `404 page not found` (cuerpo estándar de Go), sin escribir nada.

## GET /demo/estado
- Habilitada: `200 {"habilitado":true}`.
- Deshabilitada: `404`.

## POST /demo/datos
Sin cuerpo.

| Caso | Estado | Cuerpo |
|---|---|---|
| Carga exitosa | 201 | `{"proyecto_id":25,"usuarios":4,"integrantes":3,"sprints":3,"historias":10,"registros_horas":13,"rondas_votacion":1,"votos":3}` |
| Ya cargado | 409 | `{"mensaje":"los datos de prueba ya están cargados","campo":"ya_cargado"}` |
| Id 25 de otro proyecto | 409 | `{"mensaje":"el id 25 ya lo usa otro proyecto","campo":"proyecto_id"}` |
| Error interno (nada queda) | 500 | `{"mensaje":"no se pudieron cargar los datos de prueba"}` |
| Deshabilitada | 404 | `404 page not found` |

El formato de error `{"mensaje","campo"}` es el que ya parsea `compartido/api/client.ts`.

## Frontend
- `GET /demo/estado` 200 ⇒ botón «Cargar datos de prueba»; cualquier otra cosa ⇒ sin botón.
- 201 ⇒ aviso con enlace «Abrir el proyecto de prueba» (`/projects/25`).
- 409 con `campo = ya_cargado` ⇒ aviso (no error), con el mismo enlace.
- Otro error (incluido 409 `proyecto_id`) ⇒ mensaje de error.
