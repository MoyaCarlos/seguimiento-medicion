# Research: Mejoras de Markdown en el backlog (#215)

## R1. Títulos más chicos

- **Decisión**: agregar en `backlog.css` `.backlog .markdown h1 { font-size: 1.25rem; }`,
  `h2 { 1.15rem; }` y `h3 { 1.1rem; }`, con `line-height` y márgenes compactos.
- **Rationale**: el `h1` global de `App.css` se aplica dentro del Markdown y compite con el
  título de la página. El selector acotado a `.backlog .markdown` no afecta otras pantallas.
- **Alternativas**: bajar el `h1` global (afecta toda la app) o dejar que el autor use otro
  nivel (no controla el Markdown que ya existe).

## R2. Ayuda de sintaxis

- **Decisión**: un texto único bajo Descripción y Criterios: «Admite Markdown: negrita,
  listas, enlaces e imágenes por URL con ![texto](https://...). No se pueden pegar
  imágenes.».
- **Rationale**: es corto, menciona el caso de las imágenes y evita que se intente pegar un
  archivo. Se muestra en creación y edición mediante el mismo componente.
- **Alternativas**: ayuda expandible o enlace a una guía externa (más superficie para algo
  que se resuelve con una línea).

## R3. Barra de formato

- **Decisión**: botones de Título, Negrita, Cursiva, Lista, Enlace, Código e Imagen. La
  barra opera sobre `selectionStart`/`selectionEnd` del `textarea`:
  - énfasis: envuelve la selección con `**`, `*` o backticks;
  - título/lista: prefija cada línea seleccionada con `# ` o `- `, sin duplicar el prefijo;
  - enlace/imagen: usa la selección como texto de `[texto](https://...)` o
    `![texto](https://...)`;
  - sin selección, inserta un texto de ejemplo y deja el cursor dentro.
- **Rationale**: no agrega dependencias y cubre la sintaxis que el issue enumera. Se usa
  `onMouseDown.preventDefault()` para no perder la selección al pulsar un botón.
- **Alternativas**: un editor de terceros (peso y dependencias, descartado por el issue);
  insertar solo texto fijo (no aprovecha la selección).

## R4. Vista previa

- **Decisión**: pestañas «Escribir» y «Vista previa». La vista previa renderiza el valor
  actual con el componente `Markdown` existente y no envía nada. Si el texto está vacío,
  muestra «Nada para previsualizar».
- **Rationale**: reutiliza el mismo renderizador y saneamiento que las tarjetas; el estado
  del formulario no se toca al cambiar de pestaña.
- **Alternativas**: un `iframe` o un segundo renderizador (duplicaría reglas y sanitización).

## R5. Imágenes

- **Decisión**: `.backlog .markdown img` conserva `max-width: 100%` y `height: auto`, y suma
  `max-height` (20 rem) con `object-fit: contain`. Las imágenes solo se admiten por URL.
- **Rationale**: una imagen externa grande no rompe la tarjeta; el navegador la pide por su
  URL, aceptado mientras el contenido lo escriba el equipo. No hay carga de archivos.
- **Alternativas**: ocultar imágenes (`disallowedElements={['img']}`) o subirlas a un
  servicio externo; ambas están fuera de alcance.

## R6. Accesibilidad y pruebas

- **Decisión**: `label` asociado, `role="toolbar"` con `aria-label`, botones con nombre,
  `aria-describedby` a la ayuda cuando no hay error y al error cuando lo hay, y pruebas de
  componente con Given-When-Then. Se mantienen las pruebas de `Markdown`, `crear` y
  `Listado`.
- **Rationale**: la barra no debe romper el formulario ni los mensajes de error existentes;
  el helper `mensajeDe` de las pruebas lee un único id en `aria-describedby`.
- **Alternativas**: combinar ayuda y error en `aria-describedby` (rompería el helper actual
  y duplicaría el anuncio).
