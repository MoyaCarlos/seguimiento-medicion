# Data Model: Mejoras de Markdown en el backlog (#215)

## Sin cambios de datos

Esta mejora es solo de presentación en el frontend. No agrega tablas, columnas, campos,
endpoints ni cambios en el contrato de las historias.

- `historias.descripcion` y `historias.criterios_aceptacion` siguen siendo texto libre
  opcional (2000 caracteres como máximo).
- El Markdown se renderiza en el cliente con `react-markdown`, como desde #147/#211.
- No se suben ni se guardan imágenes: las imágenes son URLs externas dentro del texto.
- El límite, la normalización y el saneamiento existentes no cambian.

El componente `CampoMarkdown` solo mantiene estado de UI (vista de edición/preview y
selección del `textarea`); ese estado no se persiste ni se envía a la API.

No se modifica `MODELO-DE-DATOS.md`.
