# Research: Markdown seguro en el Product Backlog

## Librería de Markdown

- **Decisión**: `react-markdown` sin `rehype-raw`.
- **Rationale**: renderiza CommonMark con componentes de React (sin
  `dangerouslySetInnerHTML`), no interpreta HTML crudo por defecto y transforma las
  URLs inseguras. Es compatible con React 19 y Vite.
- **Alternativas**:
  - `markdown-to-jsx`: también evita HTML crudo, pero su configuración y su API de
    componentes son más limitadas para el caso.
  - `marked` + `DOMPurify`: obliga a mantener un sanitizador aparte y a inyectar HTML,
    justo lo que se quiere evitar.
  - `micromark` a mano: demasiado trabajo para el alcance.
- **Descartado**: `rehype-raw` y cualquier opción que habilite HTML crudo.

## Modelo de amenaza

El texto de una historia lo escribe una persona, por lo que es entrada no confiable.
Los ataques a cubrir son:

- `<script>alert(1)</script>` en descripción o criterios;
- `<img src=x onerror="...">`;
- enlaces con `javascript:` u otro protocolo inseguro.

`react-markdown` sin `rehype-raw` no crea nodos HTML a partir de etiquetas escritas en
el texto. Su `defaultUrlTransform` solo deja pasar protocolos y rutas seguras, así que
un enlace `javascript:` no queda como `href` ejecutable. Los tests automatizan los tres
casos.

## CommonMark

Alcanza con las construcciones estándar: párrafos, listas, énfasis, código en línea,
bloques de código, encabezados y enlaces. No se agrega `remark-gfm`; tablas, tachado y
listas de tareas quedan fuera.

## Integración

`Listado.tsx` ya tiene el texto de `descripcion` y `criterios_aceptacion`. Se reemplaza
el `<p>` por el componente `Markdown`, que devuelve un contenedor con la clase
`markdown`. Los formularios no se tocan: el texto se sigue editando y enviando tal cual.
