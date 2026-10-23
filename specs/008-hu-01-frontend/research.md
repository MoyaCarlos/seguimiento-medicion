# Research: Pantalla de creación de historias (HU-01, frontend)

Fase 0 de `/speckit.plan`. Resuelve las decisiones técnicas que la spec dejó en manos
del diseño. No quedan `NEEDS CLARIFICATION`.

## 1. Organización de la lógica: un módulo puro dentro de la feature

- **Decisión**: `web/src/features/backlog/historia.ts` concentra
  `contarCaracteres`, `validarHistoria`, `armarCuerpo`, `parsearProyectoId` y
  `mapearError`. La pantalla (`index.tsx`) solo mantiene el estado del formulario y
  dibuja lo que esas funciones devuelven.
- **Rationale**: la spec exige lógica en funciones puras probadas con Vitest y prohíbe
  tocar `compartido/`. Un módulo por pantalla mantiene el aislamiento entre features y
  evita duplicar reglas en el componente. Además permite RED→GREEN→REFACTOR sobre
  validación, conteo y cuerpo sin depender del DOM.
- **Alternativas**: lógica dentro del componente (difícil de testear sin DOM y mezcla
  reglas con render); carpeta `compartido/` (prohibido por la spec y la base);
  librería de formularios (dependencia nueva, rechazada).

## 2. Conteo de caracteres Unicode

- **Decisión**: contar con `[...texto].length` sobre el texto recortado.
- **Rationale**: el operador spread recorre la cadena por puntos de código, no por
  unidades UTF-16, y coincide con `len([]rune(texto))` del backend
  (`internal/domain/backlog_item.go`). La spec y las clarificaciones piden máximos
  inclusivos de 200/2000 sobre el texto recortado y por caracteres Unicode.
- **Alternativas**: `texto.length` (cuenta unidades UTF-16: un emoji fuera del BMP vale 2,
  se rechazaría antes que el backend); `Intl.Segmenter` con grafemas (cuenta distinto
  que el backend y agrega una API que la spec no pide); `Array.from` (equivalente a
  spread, sin ventaja).

## 3. Estado del formulario

- **Decisión**: `useState` local en `index.tsx` con un objeto
  `{ titulo, descripcion, prioridad, valorNegocio }`, más estado para `enviando`,
  `erroresCampos`, `errorGeneral` y `confirmacion`.
- **Rationale**: KISS/YAGNI; una sola pantalla, sin estado compartido (la base ya
  resuelve el proyecto por URL). El proyecto no se guarda: se lee de `useParams` en cada
  envío.
- **Alternativas**: Context o store global (innecesario); react-hook-form u otra
  librería (dependencia nueva y prohibida).

## 4. Controles que `Field` no cubre

- **Decisión**: título con `Field`; descripción con `<textarea>`, prioridad y valor de
  negocio con `<select>`, todos dentro de `<div className="field">`, con `label htmlFor`,
  `id`, `aria-invalid` y `aria-describedby`, y el error con `ErrorMessage`.
- **Rationale**: `Field` solo acepta `InputHTMLAttributes<HTMLInputElement>`; extenderlo
  sería tocar `compartido/`. La clase `field` ya da la separación etiqueta/control y
  `ErrorMessage` ya tiene `role="alert"`; así se mantiene la misma experiencia sin
  duplicar piezas.
- **Alternativas**: cambiar `Field` para aceptar textarea/select (prohibido: modifica la
  base y afecta a otras pantallas); un componente nuevo en `compartido/` (mismo
  problema).

## 5. Estilo de los controles no-input

- **Decisión**: `web/src/features/backlog/backlog.css` con reglas acotadas a
  `.backlog-form textarea` y `.backlog-form select`, reutilizando las variables de
  `App.css` (`--line`, `--ink`, `--lilac`). Se importa desde `index.tsx`.
- **Rationale**: el CSS base estiliza `.field input` pero no `textarea`/`select`; sin esa
  hoja los controles quedarían con el estilo por defecto del navegador. La hoja es local a
  la feature, no incorpora librerías y no cambia ninguna otra pantalla.
- **Alternativas**: modificar `App.css` o `compartido/` (cambia código compartido);
  dejar los controles sin estilo (inconsistente con el `input` de la base); agregar una
  librería de estilos (prohibida).

## 6. Armado del cuerpo y valor de negocio

- **Decisión**: `armarCuerpo` devuelve
  `{ proyecto_id, titulo, descripcion, prioridad, valor_negocio? }`; título y descripción
  se envían con `trim`; `valor_negocio` se omite cuando el formulario está en "Sin valor"
  y, en caso contrario, se convierte a número. `proyecto_id` se envía como número.
- **Rationale**: coincide con `crearHistoriaRequest` del backend
  (`internal/http/backlog_dto.go`) y con las clarificaciones de la spec. Omitir el campo
  opcional es la forma más simple de representar "Sin valor".
- **Alternativas**: enviar `null` (el backend lo admite, pero la spec decidió omitirlo);
  enviar strings para los números (el DTO de Go espera `int64`/`int` y fallaría el
  decode).

## 7. Errores de la API

- **Decisión**: usar `ApiError` y `ConnectionError` del cliente compartido.
  `mapearError` devuelve `{ campo?, mensaje }`: los campos conocidos del formulario se
  muestran junto al control; todo lo demás (404 de proyecto, campo desconocido,
  `ConnectionError`, error no interpretable) se muestra en el aviso general.
- **Rationale**: el cliente ya normaliza 400/404/409, el `campo` opcional y la falta de
  conexión (`web/src/compartido/api/client.ts`). La pantalla no reimplementa la lectura
  de respuestas. El 404 de un proyecto inexistente llega como `ApiError` sin `campo` y
  con el mensaje "proyecto no encontrado" (`internal/domain/errors.go`,
  `internal/http/sprint_handler.go`), que se muestra en general.
- **Alternativas**: leer `fetch` en la pantalla (duplica la base); distinguir 404 por
  estado y mostrar un texto fijo (la spec acepta el mensaje del servidor y no agrega esa
  regla).

## 8. Estrategia de pruebas

- **Decisión**: `historia.test.ts` cubre las funciones puras (recorte, Unicode, límites
  exactos, cuerpo con/sin valor, parseo del id, mapeo de errores). `index.test.tsx`
  simula `fetch` y tiene un `it` por escenario de aceptación, con el escenario en el
  nombre, más los casos límite (id inválido, backend caído, doble envío, título repetido).
- **Rationale**: cumple la restricción de la spec (Testing Library + jsdom, escenario en
  el nombre) y prueba el flujo completo —cuerpo enviado, confirmación, errores— sin
  duplicar los tests del cliente compartido. Simular `fetch` verifica el contrato real
  del POST, que es lo que el backend recibe.
- **Alternativas**: inyectar la función de API como prop (agrega un seam solo para
  tests); mockear el módulo `api.ts` (no verifica el cuerpo real); tests de integración
  con backend levantado (más lentos y frágiles para la suite de `npm test`).

## 9. Identificador de proyecto de la URL

- **Decisión**: `parsearProyectoId(valor)` acepta solo `/^\d+$/` y un resultado mayor que
  cero; si no, la pantalla muestra un aviso general y deshabilita el formulario con
  `fieldset disabled`, sin llamar a la API.
- **Rationale**: replica la regla del backend (`parsearID`: entero positivo) y evita
  aceptar valores como `1.5`, `1e3` o ` 1` que el backend rechazaría. Deshabilitar el
  formulario es la decisión registrada en las clarificaciones.
- **Alternativas**: `Number(valor)` (acepta decimales y notación científica); dejar el
  formulario habilitado y avisar al guardar (descartado en clarify).
