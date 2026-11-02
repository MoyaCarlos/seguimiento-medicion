# Quickstart: Imputación de horas en pantalla (HU-07 Frontend)

## Pruebas automáticas

```sh
cd web
npm run lint
npm test
npm run build
```

Todo tiene que quedar en verde, incluido `aislamiento.test.ts`.

## Validación manual

**Preparación**: desde la raíz del repo, `go run ./cmd/api`; en otra terminal, `cd web` y `npm run dev`. Crear un proyecto con al menos dos integrantes y una historia, crear un Sprint, asignarle la historia e iniciarlo con la fecha de hoy dentro de su rango. Abrir la pantalla **Sprints**.

| # | Paso | Resultado esperado | Escenario |
|---|---|---|---|
| 1 | Ver la historia en el Sprint activo | Muestra `0 h`, sin "Ver detalle" | US3.2, US3.6 |
| 2 | Elegir "Registrar horas" | Se abre el formulario, con la lista sin nadie seleccionado y la fecha de hoy | US1.1 |
| 3 | Guardar sin elegir integrante | "Elegí un integrante" junto a la lista; no se envía nada | US2.1 |
| 4 | Elegir integrante, actividad y horas `abc` | Error de la API junto a horas; lo escrito se conserva | US2.2, US2.7 |
| 5 | Cambiar horas a `-2`, después a `25` | Error de la API junto a horas | US2.2 |
| 6 | Horas `1,5` y actividad solo con espacios | Error de la API junto a actividad | US2.3 |
| 7 | Actividad válida y fecha de mañana | Error de la API junto a fecha | US2.4 |
| 8 | Fecha de hoy y guardar | Confirmación; total `1,5 h`; se vacían actividad y horas; siguen integrante y fecha | US1.2, US1.3, US3.3 |
| 9 | Guardar otra carga de `4` y elegir "Ver detalle" | Total `5,5 h`; se listan las dos cargas | US3.1, US3.5, US3.7 |
| 10 | Completar la historia (Iniciar y Completar) y cerrar el Sprint | Sin "Registrar horas"; el total sigue visible. Una historia no completada volvería al backlog y saldría de la tarjeta | US1.4 |
| 11 | Dejar el formulario abierto en otra pestaña, cerrar el Sprint desde esta y guardar allá | Aviso general: solo se imputan horas a historias del Sprint activo | US2.6 |
| 12 | Recargar la página | Los totales siguen igual | SC-005 de la spec 019 |

## Resultado de la validación (2026-11-02)

Los 12 pasos pasaron en Chrome con un script automatizado, contra el backend real (`go run ./cmd/api`, sirviendo `web/dist`) y una base temporal, sin errores de consola. Queda pendiente cronometrar una carga hecha por una persona (SC-001).
