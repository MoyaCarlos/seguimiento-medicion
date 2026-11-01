# Feature Specification: Markdown seguro en el Product Backlog

**Feature Branch**: `feature/147-markdown-historias`

**Issue**: #147 (extiende la pantalla de #98, HU-15 Frontend)

**Created**: 2026-10-31

**Status**: Implementada (pendiente de PR y review)

**Dependencia**: el listado y la edición del Product Backlog ya están en `main`
(spec 020). Esta spec solo cambia cómo se muestra la descripción y los criterios de
aceptación; no cambia la API, el esquema ni los límites de longitud.

**Input**: User description: "Mostrar la descripción y los criterios de aceptación de
las historias con Markdown (Issue #147). Hoy el backend guarda el texto tal cual y la
pantalla lo muestra como texto plano. Hay que renderizarlo con una librería sin
habilitar HTML crudo, para no abrir un hueco de XSS; un test debe comprobar que un
`<script>` o un `<img onerror=...>` escrito en la descripción no se ejecuta ni se
inserta como HTML. Opcional: vista previa en los formularios de creación y edición."

## Objetivo

Que un Product Builder vea la descripción y los criterios de aceptación de cada
historia con su formato Markdown (listas, negritas, bloques de código, enlaces) en el
listado del Product Backlog, sin que el contenido escrito por las personas pueda
inyectar HTML ni ejecutar código en el navegador.

## Entrada y salida

- **Entrada**: `GET /projects/{id}/backlog` sigue devolviendo `descripcion` y
  `criterios_aceptacion` como texto plano, sin cambios (spec 020).
- **Salida en pantalla**: la descripción y los criterios se muestran interpretando
  Markdown dentro del listado. El texto plano actual se sigue viendo igual.
- **Edición**: los formularios siguen mostrando el texto original en sus campos
  editables y `PUT /backlog/{id}` sigue enviando exactamente el mismo texto, sin
  normalizarlo ni convertirlo.
- **Sin cambio de contrato**: no se agregan campos, endpoints, tablas, migraciones ni
  límites nuevos.

## Decisión de dependencia

Esta es la primera dependencia nueva del frontend desde la base. Se agrega
`react-markdown` (y su árbol de dependencias) porque:

- Interpreta CommonMark con componentes de React, sin usar `dangerouslySetInnerHTML`.
- **No renderiza HTML crudo** por defecto: basta con no agregar `rehype-raw`. Eso
  cubre el riesgo de XSS pedido por el issue.
- Sanitiza las URLs de los enlaces y las imágenes con su transformación por defecto
  (solo protocolos seguros), sin necesidad de otra librería.
- Es compatible con React 19 y con Vite; no requiere configuración especial.

Las alternativas (`markdown-to-jsx`, `marked` + sanitizador, `micromark` a mano)
obligan a sumar o mantener más piezas para el mismo resultado. `react-markdown` es la
opción más directa para el stack actual.

## Reglas de comportamiento (EARS)

- **WHEN** se muestra la descripción de una historia, **THEN** la pantalla **shall**
  interpretarla como Markdown y **shall not** mostrarla como texto plano.
- **WHEN** una historia tiene criterios de aceptación, **THEN** la pantalla **shall**
  interpretarlos como Markdown.
- **WHEN** una historia no tiene criterios de aceptación, **THEN** la pantalla
  **shall** seguir mostrando el texto "Sin criterios de aceptación", sin intentar
  interpretar un valor nulo.
- **WHEN** el texto contiene listas, negritas, cursivas, enlaces o bloques de código,
  **THEN** la pantalla **shall** mostrar el resultado formateado correspondiente.
- **WHEN** el texto es plano, sin marcas de Markdown, **THEN** la pantalla **shall**
  mostrar el mismo contenido visible que antes.
- **IF** el texto contiene HTML crudo, incluidos `<script>` o `<img onerror="...">`,
  **THEN** la pantalla **shall not** insertarlo como HTML, **shall not** ejecutar
  scripts ni manejadores de eventos, y **shall not** permitir URLs con protocolos
  inseguros.
- **WHEN** se abre la edición de una historia, **THEN** los campos **shall** seguir
  mostrando el texto original (Markdown sin interpretar) para poder editarlo.
- **WHEN** se guarda una edición, **THEN** el cuerpo enviado a `PUT /backlog/{id}`
  **shall** conservar el mismo texto que escribió la persona, sin convertirlo ni
  recortarlo más allá de lo que ya hace la pantalla.
- **WHEN** cambia el proyecto o se actualiza una historia, **THEN** el listado **shall**
  volver a mostrar la descripción y los criterios con el formato correspondiente.

## Casos límite

- Descripción con un bloque de código que contiene etiquetas HTML: se muestra como
  código, no como HTML.
- Criterios `null`: se mantiene el mensaje "Sin criterios de aceptación".
- Texto de hasta 2000 caracteres: se sigue mostrando completo; el límite no cambia.
- Enlaces con `javascript:` u otro protocolo inseguro: no se convierten en enlaces
  ejecutables.
- Imágenes Markdown: se muestran con una URL segura; los atributos `onerror` escritos
  como HTML crudo no se conservan.
- Markdown mal formado: se muestra como texto razonable, sin romper el listado ni
  lanzar un error.
- Texto que solo contiene espacios: no debe romper la pantalla.

## Condiciones de error

- **IF** la interpretación de Markdown falla para un texto, **THEN** la pantalla
  **shall** seguir mostrando el listado y no **shall** dejar la historia inaccesible.
- **IF** el texto contiene un intento de XSS, **THEN** la pantalla **shall** mostrarlo
  como contenido inerte o descartarlo, sin ejecutar código.
- **IF** la API falla al cargar o guardar, **THEN** se mantiene el manejo de errores ya
  definido en la spec 020; esta historia no lo modifica.

## Criterios de aceptación

1. **Given** una historia cuya descripción usa negritas, listas o código, **When** se
   abre el Product Backlog, **Then** la descripción se ve formateada y no como texto
   plano.
2. **Given** una historia con criterios de aceptación en Markdown, **When** se abre el
   listado, **Then** los criterios se ven formateados.
3. **Given** una descripción que contiene `<script>alert(1)</script>` y
   `<img src=x onerror="alert(2)">`, **When** se muestra en el listado, **Then** no se
   inserta un `<script>`, no se inserta una imagen con `onerror`, no se ejecuta código
   y el resto de la descripción sigue visible.
4. **Given** una historia sin criterios de aceptación, **When** se abre el listado,
   **Then** se ve "Sin criterios de aceptación".
5. **Given** una historia con Markdown, **When** se abre su edición y se guarda sin
   cambios, **Then** los campos muestran el texto original y el cuerpo de
   `PUT /backlog/{id}` conserva ese mismo texto.
6. **Given** un enlace Markdown con protocolo `javascript:`, **When** se muestra,
   **Then** no queda un enlace ejecutable.

## Fuera de alcance

- Vista previa del Markdown en los formularios de creación y edición (parte opcional
  del issue; se puede agregar en otra historia).
- Extensiones GFM como tablas, tachado o listas de tareas; alcanza con CommonMark.
- Cambios en la API, el esquema, las migraciones o los límites de longitud.
- Sanitización de HTML crudo para permitirlo de forma controlada: el HTML crudo no se
  habilita.
- Cambios en el orden, los filtros o la edición del backlog.

## Trazabilidad

| Criterio | Implementación prevista | Prueba |
|---|---|---|
| Descripción y criterios formateados | `web/src/features/backlog/Markdown.tsx` + `Listado.tsx` | `Markdown.test.tsx` y caso de listado |
| HTML crudo inerte | `react-markdown` sin `rehype-raw` | test de XSS de `Markdown.test.tsx` |
| Sin criterios | `Listado.tsx` | caso existente de listado |
| Edición sin cambios de contrato | `Edicion.tsx` sin cambios | tests existentes de edición |
| Enlaces seguros | transformación por defecto de `react-markdown` | test de enlace `javascript:` |

## Verificación

- `Markdown.test.tsx`: 4/4 en verde (texto plano, CommonMark, XSS y enlace
  `javascript:`).
- `Listado.test.tsx`: 17/17 en verde, incluido el caso de Markdown en descripción y
  criterios.
- `npm test -- --run`: 213/213 en verde en la corrida final. Una corrida previa tuvo un
  timeout aislado de carga en `rutas.test.tsx`; ese archivo pasó solo.
- `npm run lint` y `npm run build`: en verde. `npm audit --omit=dev`: 0 vulnerabilidades
  de producción (las 4 reportadas son de Vite/Vitest en desarrollo y ya venían con el
  stack).

## Assumptions

- El texto llega desde la API tal como se guardó; el frontend no lo normaliza.
- CommonMark cubre los casos reales del equipo (listas, negritas, código y enlaces).
- La dependencia se agrega al `package.json` de `web/` y se justifica en esta spec.
- No hay HTML crudo legítimo que mostrar; si algún día se necesita, será una decisión
  explícita con sanitización, fuera de esta historia.
