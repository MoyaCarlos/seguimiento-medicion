# Contrato de UI: editor Markdown del backlog (#215)

## `CampoMarkdown`

Componente de frontend reutilizable en creación y edición de historias.

| Prop | Tipo | Obligatorio | Descripción |
|---|---|---|---|
| `id` | string | sí | id del `textarea` y base de los ids de ayuda/error |
| `name` | string | sí | name del `textarea` |
| `label` | string | sí | texto del `label` asociado |
| `value` | string | sí | texto Markdown actual |
| `onChange` | `(valor: string) => void` | sí | actualiza el valor del formulario |
| `error` | string | no | mensaje de error; se muestra con `role="alert"` |
| `disabled` | boolean | no | deshabilita el `textarea` y los botones de formato |
| `rows` | number | no | filas del `textarea` (4 por defecto) |

### Ayuda

Texto fijo debajo del campo:

> Admite Markdown: negrita, listas, enlaces e imágenes por URL con ![texto](https://...).
> No se pueden pegar imágenes.

Se asocia con `aria-describedby` cuando no hay error; con error, `aria-describedby` apunta
al mensaje de error.

### Barra de formato

`role="toolbar"`, `aria-label="Formato Markdown"`. Botones `type="button"`:

| Botón | Resultado sobre la selección |
|---|---|
| Título | prefija cada línea con `# ` (sin duplicar) |
| Negrita | `**selección**` |
| Cursiva | `*selección*` |
| Lista | prefija cada línea con `- ` (sin duplicar) |
| Enlace | `[selección](https://...)` |
| Código | `` `selección` `` |
| Imagen | `![selección](https://...)` |

Sin selección se inserta un texto de ejemplo y el cursor queda dentro. Con el campo
deshabilitado, los botones también se deshabilitan.

### Pestañas

- **Escribir**: muestra la barra y el `textarea`.
- **Vista previa**: renderiza `value` con `Markdown`; no envía ni guarda nada. Si el texto
  está vacío, muestra «Nada para previsualizar».
- Volver a Escribir conserva el texto.

## Estilos

- `.backlog .markdown h1`: 1.25 rem; `h2`: 1.15 rem; `h3`: 1.1 rem.
- `.backlog .markdown img`: `max-width: 100%`, `height: auto`, `max-height: 20rem`,
  `object-fit: contain`.
- La barra, las pestañas, la ayuda y la vista previa viven en `backlog.css`, sin estilos
  globales nuevos.
