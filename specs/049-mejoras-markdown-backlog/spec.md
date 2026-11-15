# Feature Specification: Mejorar la presentación del Markdown en el backlog (#215)

**Feature Branch**: `feature/215-mejoras-markdown-backlog`

**Issue**: #215

**Created**: 2026-11-15

**Status**: Implementada (PR #326, pendiente de review)

**Dependencia**: #147 (PR #211, Markdown en la descripción y los criterios) y #174
(criterios al crear una historia), ya en `main`. Solo frontend; no cambia la API, el
esquema ni el contrato de las historias.

**Input**: User description: "Mejorar la presentación del Markdown en el backlog: títulos
más chicos, ayuda de sintaxis y barra de formato. Todo es estético o de usabilidad; no
cambia el comportamiento ni la API. Las imágenes pegadas del portapapeles quedan fuera de
alcance."

## Objetivo

Que el Markdown de la descripción y los criterios se vea integrado en la jerarquía de la
pantalla, que quien escribe sepa qué sintaxis se admite y que pueda insertar formato sin
memorizarla, sin agregar dependencias nuevas.

## Entradas y salidas

- **Entrada**: el texto Markdown de la descripción o los criterios, en el formulario de
  creación (HU-01/#174) o de edición (HU-15), y el texto ya guardado en las tarjetas del
  backlog.
- **Salida**: el mismo Markdown renderizado con títulos acotados e imágenes de altura
  limitada; una ayuda visible bajo cada campo; una barra de formato que inserta sintaxis en
  la selección y una pestaña «Vista previa» que reutiliza `Markdown`.
- **Sin cambios**: la API, el esquema, el límite de 2000 caracteres, el envío del texto
  recortado y el saneamiento de enlaces/imágenes que ya hace `react-markdown`.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Leer Markdown sin romper la jerarquía (Priority: P1)

Quien consulta el backlog ve los títulos `#`, `##` y `###` de una historia con un tamaño
menor que el título de la página, y las imágenes grandes no ocupan media pantalla.

**Independent Test**: renderizar una historia con títulos e imagen; comprobar que los
títulos quedan dentro de `.backlog .markdown` con el tamaño acotado por CSS y que la imagen
tiene `max-width: 100%` y `max-height`.

**Acceptance Scenarios**:

1. **Given** una descripción con `# Título`, **When** se muestra en una tarjeta del
   backlog, **Then** el `h1` renderizado no supera 1.25 rem y mantiene la jerarquía debajo
   del título de la página.
2. **Given** una descripción con `##` o `###`, **When** se muestra, **Then** los títulos
   quedan entre 1.1 y 1.25 rem y con márgenes compactos.
3. **Given** una imagen externa por URL, **When** se muestra, **Then** se conserva su `src`
   https, respeta `max-width: 100%` y no supera la altura máxima definida.

---

### User Story 2 - Saber que se admite Markdown (Priority: P1)

Quien crea o edita una historia ve, debajo de Descripción y de Criterios de aceptación, una
ayuda que enumera la sintaxis admitida y aclara que las imágenes son solo por URL.

**Independent Test**: abrir los formularios de creación y edición y comprobar la ayuda en
ambos campos.

**Acceptance Scenarios**:

1. **Given** el formulario de creación, **When** se muestra, **Then** Descripción y
   Criterios de aceptación tienen la ayuda «Admite Markdown: negrita, listas, enlaces e
   imágenes por URL con ![texto](https://...). No se pueden pegar imágenes.».
2. **Given** el formulario de edición, **When** se muestra, **Then** ambos campos tienen la
   misma ayuda.
3. **Given** un error de validación en un campo, **When** se muestra, **Then** el mensaje de
   error sigue junto al campo y la ayuda no se duplica.

---

### User Story 3 - Insertar formato y previsualizar (Priority: P2)

Quien no conoce la sintaxis usa una barra de formato sobre el campo: título, negrita,
cursiva, lista, enlace, código e imagen por URL. Los botones envuelven o prefijan la
selección. Una pestaña «Vista previa» muestra el resultado con el mismo componente Markdown.

**Independent Test**: seleccionar texto, pulsar Negrita y comprobar la sintaxis; escribir un
título y abrir Vista previa para verlo renderizado.

**Acceptance Scenarios**:

1. **Given** un texto seleccionado, **When** se pulsa Negrita, Cursiva o Código, **Then** la
   selección queda envuelta con `**`, `*` o `` ` ``.
2. **Given** una línea o selección, **When** se pulsa Título o Lista, **Then** cada línea
   afectada recibe el prefijo `# ` o `- `.
3. **Given** una selección, **When** se pulsa Enlace o Imagen, **Then** se inserta
   `[texto](https://...)` o `![texto](https://...)` usando la selección como texto.
4. **Given** el formulario con Markdown escrito, **When** se elige «Vista previa», **Then**
   se muestra el Markdown renderizado con el componente existente y el texto no se pierde.
5. **Given** la vista previa vacía, **When** se abre, **Then** se informa que no hay
   contenido para previsualizar.
6. **Given** un campo deshabilitado (proyecto cerrado o envío en curso), **When** se muestra
   la barra, **Then** sus botones de edición quedan deshabilitados.

### Edge Cases

- Sin selección, los botones de énfasis insertan la sintaxis con un texto de ejemplo y dejan
  el cursor dentro; Título y Lista actúan sobre la línea actual.
- Selección que abarca varias líneas: Título y Lista aplican el prefijo a cada línea sin
  duplicarlo si ya estaba.
- El texto vacío se previsualiza como «Nada para previsualizar», no como error.
- El límite de 2000 caracteres y la validación siguen siendo los del formulario; la barra no
  los relaja.
- Pegar una imagen del portapapeles no se convierte en archivo ni en Markdown: solo se admite
  una URL escrita a mano.
- Los enlaces e imágenes con protocolos inseguros siguen saneados por `Markdown`.

## Requirements *(mandatory)*

### Reglas de comportamiento (EARS)

**Ubicuas**

- **FR-001**: Los `h1`, `h2` y `h3` dentro de `.backlog .markdown` MUST tener un tamaño
  entre 1.1 y 1.25 rem y una jerarquía menor que el título de la página.
- **FR-002**: Las imágenes dentro de `.backlog .markdown` MUST conservar `max-width: 100%`,
  `height: auto` y MUST tener una altura máxima razonable.
- **FR-003**: La ayuda de sintaxis MUST mostrarse debajo de Descripción y de Criterios de
  aceptación, en creación y edición, con el texto acordado.
- **FR-004**: La barra de formato MUST ofrecer título, negrita, cursiva, lista, enlace,
  código e imagen por URL, y MUST insertar la sintaxis sobre la selección o la línea actual.
- **FR-005**: La pestaña Vista previa MUST renderizar el valor actual con el componente
  `Markdown` existente, sin enviar ni guardar nada.
- **FR-006**: La funcionalidad MUST ser solo de frontend, sin dependencias nuevas ni cambios
  de API, esquema o contrato.

**Basadas en evento**

- **FR-007**: WHEN se pulsa un botón con texto seleccionado, el sistema SHALL envolver o
  prefijar solo esa selección y conservar el resto del texto.
- **FR-008**: WHEN se cambia a Vista previa, el sistema SHALL mostrar el Markdown renderizado
  y SHALL conservar el texto para volver a Escribir.
- **FR-009**: WHEN el campo está deshabilitado, los botones de edición SHALL quedar
  deshabilitados y SHALL NOT modificar el texto.

**No deseadas**

- **IF** el campo está vacío, THEN la vista previa SHALL mostrar un aviso de que no hay
  contenido, sin error.
- **IF** se intenta pegar una imagen del portapapeles, THEN no SHALL subirse ni guardarse
  ningún archivo; el comportamiento queda fuera de alcance.

### Requisitos de la interfaz

- **FR-010**: El editor MUST reutilizar el componente `Markdown` para la vista previa y el
  mismo estilo `.backlog .markdown` para las tarjetas y la previsualización.
- **FR-011**: Los controles MUST ser accesibles: labels asociados, botones con nombre,
  `aria-describedby` para la ayuda o el error, y `role="toolbar"` en la barra.
- **FR-012**: La ayuda y la barra MUST aparecer en el formulario de creación de HU-01/#174 y
  en el de edición de HU-15.

## Fuera de alcance

- Pegar o subir imágenes del portapapeles: requiere almacenamiento persistente y el servicio
  de Render no lo tiene.
- Un editor Markdown de terceros o una dependencia nueva.
- Cambios en la API, el esquema, el saneamiento o el límite de 2000 caracteres.
- Editar la descripción o los criterios fuera de los formularios actuales.

## Success Criteria *(mandatory)*

- **SC-001**: Los títulos Markdown del backlog quedan visualmente por debajo del título de la
  página y las imágenes no superan la altura máxima.
- **SC-002**: El 100 % de los formularios de creación y edición muestran la ayuda de sintaxis
  en Descripción y Criterios.
- **SC-003**: Los siete botones de formato insertan la sintaxis esperada y la vista previa
  renderiza el mismo Markdown que las tarjetas.
- **SC-004**: Las pruebas de componente de la barra, la ayuda y la vista previa pasan en
  Vitest, junto con las existentes de crear/editar y `Markdown`.
- **SC-005**: `npm test`, `npm run lint` y `npm run build` pasan; no se agregan dependencias.

## Assumptions

- #211 y #174 ya están en `main`; esta mejora no reabre su comportamiento.
- El componente `Markdown` existente ya sanea HTML y protocolos inseguros; se reutiliza.
- La ayuda es la misma para descripción y criterios.
- La barra es una mejora de usabilidad, no un editor completo: no hay atajos de teclado
  globales ni subida de archivos.
- Solo se tocan archivos del frontend (`web/src/features/backlog/` y `backlog.css`).
