# Quickstart: validar #215

## Prerrequisitos

- Node con las dependencias de `web/` instaladas.
- La API y el frontend levantados (`go run ./cmd/api` y `npm run dev`) o las pruebas de
  componente con Vitest.

## Verificación automática

```sh
cd web
npm test -- src/features/backlog/CampoMarkdown.test.tsx
npm test -- src/features/backlog/crear.test.tsx src/features/backlog/Listado.test.tsx
npm test
npm run lint
npm run build
```

Esperado: todo en verde, sin dependencias nuevas.

## Verificación manual

1. Abrir el backlog de un proyecto con una historia cuya descripción tenga `#`, `##`, `###`,
   una lista y una imagen https.
2. Comprobar que los títulos se ven por debajo del título de la página y que la imagen no
   supera la altura máxima.
3. Abrir «Nueva historia»: Descripción y Criterios muestran la ayuda de Markdown.
4. Seleccionar texto en Descripción y probar Negrita, Cursiva, Código, Enlace e Imagen; el
   texto queda con la sintaxis alrededor de la selección.
5. Probar Título y Lista sobre una línea o varias líneas; se agrega el prefijo sin duplicarlo.
6. Escribir un título, cambiar a «Vista previa» y comprobar que se renderiza; volver a
   «Escribir» y comprobar que el texto sigue.
7. Repetir la ayuda y la barra en el formulario de edición de una historia.
8. Con el proyecto cerrado, comprobar que la barra no permite editar.

## Criterios de salida

- Títulos e imágenes acotados en el backlog.
- Ayuda visible en los dos formularios y en los dos campos.
- Barra de formato y vista previa funcionando con el componente `Markdown` existente.
- Sin cambios de API, esquema ni dependencias.
