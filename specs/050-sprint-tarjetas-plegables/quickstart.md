# Quickstart: validar #219

## Prerrequisitos

- Node con las dependencias de `web/` instaladas.
- La API y el frontend levantados (`go run ./cmd/api` y `npm run dev`) o las pruebas de
  componente con Vitest.

## Verificación automática

```sh
cd web
npm test -- src/features/sprints/EstadoHistoria.test.tsx src/features/sprints/SprintBacklog.test.tsx
npm test
npm run lint
npm run build
```

Esperado: todo en verde, sin dependencias nuevas.

## Verificación manual

1. Abrir la pantalla de Sprints de un proyecto con un Sprint activo y varias historias.
2. Comprobar que cada historia aparece plegada con título, estado y Story Points.
3. Activar una tarjeta con el ratón y con el teclado; se ve la prioridad, la fecha si está
   completada, el control de estado y el esfuerzo.
4. Pulsar «Iniciar»: aparece la confirmación y la historia todavía no cambia.
5. Pulsar «Cancelar»: no cambia nada. Repetir y cancelar con Escape y con un clic fuera.
6. Pulsar «Iniciar», confirmar y comprobar que la historia pasa a «En progreso».
7. Pulsar «Completar», confirmar y comprobar el estado «Completada» y la fecha.
8. Con el proyecto cerrado, comprobar que la acción queda deshabilitada con la nota.

## Criterios de salida

- Tarjetas plegadas con resumen legible.
- Confirmación en línea antes de iniciar o completar.
- Cancelar, Escape y clic fuera no envían nada.
- Doble envío bloqueado mientras se manda.
- Sin cambios de API, esquema ni dependencias.
