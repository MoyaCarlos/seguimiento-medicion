# Implementation Plan: Pantalla de creación de historias de usuario (HU-01, frontend)

**Branch**: `feature/HU-01-frontend` | **Date**: 2026-10-22 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/008-hu-01-frontend/spec.md`

## Summary

Reemplazar la pantalla provisional de `features/backlog` por el formulario que crea
historias de usuario en el Product Backlog del proyecto de la URL. La pantalla valida
título y descripción (obligatorios, recortados, 200/2000 caracteres Unicode) antes de
enviar, muestra todas las advertencias juntas, envía `POST /backlog` con el cliente
compartido y, según la respuesta, muestra una confirmación con la historia creada o el
error junto al campo / en un aviso general. No agrega dependencias ni modifica
`compartido/`: la lógica de validación, armado del cuerpo y mapeo de errores vive en
funciones puras dentro de `features/backlog`, probadas con Vitest; la pantalla se prueba
con Testing Library + jsdom usando los escenarios de aceptación como nombre de cada test.

## Technical Context

**Language/Version**: TypeScript ~6.0, React 19 (ya instalados en `web/`)

**Primary Dependencies**: `react-router` 8 para `useParams` y el `Component` que monta la
ruta; cliente `get/post` y piezas `Field`, `Button`, `ErrorMessage`, `LoadingIndicator`
de `web/src/compartido/`. Sin dependencias nuevas: Vitest 4, Testing Library y jsdom ya
están en `web/package.json`

**Storage**: N/A (la pantalla no persiste nada; la historia la guarda la API)

**Testing**: Vitest. Funciones puras con tests unitarios; pantalla con Testing Library +
jsdom, un test por escenario de aceptación (el escenario en el nombre) y los casos límite
del spec

**Target Platform**: Navegador moderno; desarrollo local con Vite en `:5173` y proxy a la
API en `:8080`

**Project Type**: SPA que consume la API REST existente

**Performance Goals**: N/A (un formulario, una llamada por guardado)

**Constraints**: Sin tocar Go, `package.json` ni `web/src/compartido/`; la ruta
`projects/:id/backlog` y su registro en `rutas.tsx` ya existen y no cambian; ninguna
feature importa a otra; textos en español; una única creación por guardado

**Scale/Scope**: Una pantalla, cuatro campos, cuatro prioridades MoSCoW y siete valores
de negocio posibles

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principio | Estado | Nota |
|---|---|---|
| I. Test-First | Cumple | Toda la lógica (validación, conteo Unicode, armado del cuerpo, mapeo de errores, parseo del id) se desarrolla con RED→GREEN→REFACTOR y commits separados. El componente solo cablea estado y llamadas; sus pruebas cubren los escenarios de aceptación |
| II. SDD | Cumple | Spec → clarify (5 respuestas integradas) → plan → tasks → implement, una historia a la vez |
| III. BDD (Godog) | Cumple por equivalencia | Los escenarios de aceptación de esta spec se automatizan como pruebas de componente (Testing Library + jsdom), que cubren casos normales, alternativos, límite y de error. El backend que consume la pantalla (`POST /backlog`) ya está cubierto por Godog y esta historia no agrega ni cambia escenarios Gherkin, según la restricción de la spec |
| IV. Clean/Hexagonal | No aplica | Solo frontend; no toca `internal/` |
| V. Clean Code/KISS/YAGNI/DRY | Cumple | Una pantalla, estado local con `useState`, sin librería de formularios ni estado global, sin capas nuevas. Las reglas de validación y el armado del cuerpo viven en un único módulo puro reutilizado por la pantalla y sus tests |
| VI. Stack | Cumple | React + Vite (stack obligatorio). Testing Library y jsdom ya fueron aprobados e incorporados por la base del frontend; no se agrega ninguna dependencia |

Sin violaciones; no hay Complexity Tracking.

## Project Structure

### Documentation (this feature)

```text
specs/008-hu-01-frontend/
├── plan.md              # Este archivo
├── research.md          # Fase 0
├── data-model.md        # Fase 1
├── quickstart.md        # Fase 1
├── contracts/
│   └── api-backlog.md   # Fase 1
├── checklists/
│   └── requirements.md
└── tasks.md             # lo genera /speckit.tasks
```

### Source Code (repository root)

```text
web/
└── src/
    ├── rutas.tsx                      # sin cambios: backlog/* ya está registrado
    ├── compartido/                    # sin cambios: se reutiliza tal cual
    │   ├── api/client.ts
    │   └── componentes/{Field,Button,ErrorMessage,LoadingIndicator}.tsx
    └── features/
        └── backlog/
            ├── index.tsx             # Component: pantalla, estado y eventos
            ├── index.test.tsx        # pruebas de componente (escenarios + límites)
            ├── historia.ts           # funciones puras: validar, contar, armar cuerpo, mapear errores
            ├── historia.test.ts      # pruebas unitarias de esas funciones
            ├── opciones.ts           # prioridades MoSCoW y valores de negocio
            ├── api.ts                # crearHistoria(payload) sobre post() compartido
            └── backlog.css           # solo los controles textarea/select, con variables de la base
```

**Structure Decision**: todo el código de la historia vive en `web/src/features/backlog/`
(reemplaza el placeholder) y no se toca `compartido/` ni `rutas.tsx`. La pantalla usa
`Field` solo para `<input>` (título); descripción (`<textarea>`), prioridad y valor de
negocio (`<select>`) se arman en la feature con `ErrorMessage` y la clase `field`, porque
`Field` no soporta esos controles (clarificación de la spec). `backlog.css` existe solo
para que esos dos controles hereden el mismo borde/espaciado que el `input` de la base,
usando sus variables CSS; no se agregan librerías ni se modifica el CSS global.

## Design decisions

1. **Funciones puras primero.** `historia.ts` expone `contarCaracteres`, `validarHistoria`,
   `armarCuerpo`, `parsearProyectoId` y `mapearError`. La pantalla no repite ninguna regla:
   solo guarda el formulario, llama a esas funciones y dibuja el resultado. Es la única
   lógica con ciclo RED→GREEN→REFACTOR. Ver [research.md](research.md).
2. **Conteo Unicode.** `contarCaracteres` usa `[...texto].length`, que cuenta puntos de
   código y no unidades UTF-16; es lo que hace el backend con `len([]rune(texto))`. La
   validación recorta primero (`trim`), tal como se clarificó, y usa los máximos inclusivos
   200/2000.
3. **Un solo intento, un solo resultado.** Al presionar "Guardar" se limpian la
   confirmación y los errores anteriores, se valida localmente y, solo si no hay errores,
   se llama a la API. Un intento detenido por validación también reemplaza la confirmación
   anterior (clarificación). Todos los errores de campo se muestran juntos.
4. **Cuerpo del POST.** Se envían `proyecto_id` (número), `titulo` y `descripcion`
   recortados, `prioridad` con el código `M/S/C/W`; `valor_negocio` se omite si la opción
   es "Sin valor" (clarificación). El detalle exacto está en
   [contracts/api-backlog.md](contracts/api-backlog.md).
5. **Errores del servidor.** `ApiError` ya trae `status`, `message` y `field`; se mapea un
   `field` conocido (`titulo`, `descripcion`, `prioridad`, `valor_negocio`) junto a su
   campo; cualquier otro (incluido `proyecto_id`), un 404 o un `ConnectionError` van al
   aviso general `ErrorMessage` arriba del formulario. Los mensajes no se borran al
   escribir: permanecen hasta el próximo intento.
6. **Éxito.** La respuesta 201 se muestra en una confirmación con título, prioridad legible,
   "Nueva", "sin estimación" y, si vino, valor de negocio. Luego el formulario vuelve a
   `{ titulo: '', descripcion: '', prioridad: 'M', valorNegocio: '' }` y la confirmación
   queda hasta el próximo intento.
7. **Identificador inválido.** `parsearProyectoId` acepta solo dígitos que representen un
   entero positivo representable sin pérdida de precisión en JavaScript
   (`Number.isSafeInteger`, hasta 9007199254740991). Es más estricta que la regla del
   backend a propósito: un número mayor no se puede serializar con exactitud. Si falla, se
   muestra un aviso general y el formulario queda deshabilitado (`fieldset disabled`), sin
   llamar a la API (clarificación).
8. **Pruebas de la pantalla.** `index.test.tsx` simula `fetch` (el mismo mecanismo del
   cliente compartido) para verificar el cuerpo enviado y cada respuesta; cada escenario de
   aceptación del spec es un `it` cuyo nombre lo reproduce, más los casos límite y de error
   de la sección Edge Cases.

## Post-Design Constitution Check

Vuelto a evaluar después de Fase 1: los artefactos no agregan principio nuevo. La
validación, el conteo y el mapeo quedan en funciones puras con RED/GREEN; no hay
persistencia, auth ni dependencias. Sigue sin violaciones.
