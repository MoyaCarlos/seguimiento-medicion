# Quickstart: validar #225

## Prerrequisitos

- Node con las dependencias de `web/` instaladas.
- La API y el frontend levantados (`go run ./cmd/api` y `npm run dev`) o las pruebas de
  componente con Vitest.

## Verificación automática

```sh
cd web
npm test -- src/features/votacion/useVotacion.test.ts src/features/votacion/index.test.tsx
npm test
npm run lint
npm run build
```

Esperado: todo en verde, sin dependencias nuevas.

## Verificación manual

1. Abrir la pantalla de votación de una historia con una ronda abierta, en dos navegadores o
   dos ventanas con usuarios distintos.
2. Votar en la segunda ventana.
3. En la primera, sin pulsar «Actualizar», comprobar que la fila del otro integrante pasa a
   «Ya votó» / «Oculto» en menos de 5 segundos.
4. Cambiar a otra pestaña y dejar pasar varios intervalos; al volver, comprobar que se
   reanuda el sondeo y se actualiza la votación.
5. Revelar la votación desde el Scrum Master; comprobar que el resultado aparece solo y que no
   siguen llegando consultas automáticas.
6. Pulsar «Actualizar» y comprobar que sigue funcionando como antes.

## Criterios de salida

- Los votos ajenos y la revelación aparecen sin pulsar «Actualizar».
- El sondeo se pausa con la pestaña oculta, la ronda revelada o una acción en curso.
- La pantalla no muestra el indicador de carga por las consultas automáticas.
- Sin cambios de API, esquema ni dependencias.
