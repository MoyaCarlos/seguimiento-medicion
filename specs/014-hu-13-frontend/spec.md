# Feature Specification: Estado del proyecto en el panel principal (HU-13 Frontend)

**Feature Branch**: `feature/HU-13-frontend-panel`

**Fix Branch**: `fix/151-spec-hu-13-frontend` (issue #151)

**Issue**: #41 (parte de #13)

**Created**: 2026-10-23

**Status**: Implementada y mergeada (PR #150). Las mejoras de esta spec y de sus tests se
registraron en el issue #151 y se integraron en el PR #152.

**Dependencia**: el panel principal pertenece a la pantalla de HU-04
(`features/proyecto`, PR #121), ya presente en `main`. Esta rama se apoya en esa pantalla
y no duplica su implementación.

## Objetivo

Mostrar en el panel principal del proyecto el estado derivado por la API:
`Planificado`, `En curso` o `Finalizado`.

## Entrada y salida

- **Entrada**: `GET /projects/{id}/status`, que responde
  `{"estado":"Planificado" | "En curso" | "Finalizado"}`.
- **Salida en pantalla**: una etiqueta visible con el estado dentro del panel principal
  del proyecto.
- La consulta de estado se hace al abrir el panel, junto con la lectura del proyecto.
- Si la consulta falla, el panel muestra el error y no inventa un estado.

## Reglas de comportamiento (EARS)

- **WHEN** se abre el panel de un proyecto, **THEN** el sistema **shall** consultar
  `GET /projects/{id}/status`.
- **WHILE** la consulta está en curso, **THEN** el panel **shall** mostrar un indicador
  de carga para el estado.
- **WHEN** la API devuelve un estado válido, **THEN** el panel **shall** mostrar
  exactamente la etiqueta recibida (`Planificado`, `En curso` o `Finalizado`).
- **IF** la consulta de estado falla, **THEN** el panel **shall** mostrar el mensaje de
  error y conservar visible la información del proyecto que sí se pudo cargar.
- **WHEN** se recarga el panel, **THEN** el sistema **shall** volver a consultar el
  estado, sin guardarlo ni congelarlo.
- El cálculo del estado **no** se replica en el frontend: la única fuente es la API.

## Restricciones

- El estado **shall** provenir siempre de `GET /projects/{id}/status`; el frontend
  **shall not** recalcularlo a partir de fechas u otros datos del proyecto.
- El frontend **shall not** persistir ni cachear el estado entre consultas: cada apertura
  o recarga del panel **shall** consultar la API.
- La etiqueta **shall** mostrarse tal cual la devuelve la API, sin traducirla ni
  normalizarla.
- El panel **shall** reutilizar la pantalla de proyecto de HU-04 y **shall not** duplicar
  su implementación.
- El contrato de `GET /projects/{id}/status` **shall** ser el definido por HU-13; el
  frontend **shall not** modificarlo.

## Casos límite

- **IF** el proyecto no existe y la API responde 404, **THEN** el panel **shall** mostrar
  el mensaje de error de la consulta de estado y **shall not** inventar un estado.
- **IF** la API devuelve un estado fuera de `Planificado`, `En curso` o `Finalizado`,
  **THEN** el panel **shall** mostrar el valor recibido tal cual, sin corregirlo ni
  inferir otro estado.

## Condiciones de error

- **IF** la consulta de estado responde 500 o la red no está disponible, **THEN** el panel
  **shall** mostrar el mensaje de error correspondiente y **shall** conservar visible la
  información del proyecto que sí se pudo cargar.
- **IF** la consulta de estado falla, **THEN** la etiqueta de estado **shall** mostrar
  `Sin datos` y **shall not** reutilizar el estado del proyecto anterior.
- **WHEN** cambia el proyecto mostrado, **THEN** el panel **shall** reiniciar el estado,
  el error y la carga de la consulta anterior antes de mostrar los datos del nuevo
  proyecto.

## Criterios de aceptación

1. **Given** un proyecto cuyo estado es `Planificado`, **When** abro su panel,
   **Then** veo `Planificado`.
2. **Given** un proyecto cuyo estado es `En curso`, **When** abro su panel,
   **Then** veo `En curso`.
3. **Given** un proyecto cuyo estado es `Finalizado`, **When** abro su panel,
   **Then** veo `Finalizado`.
4. **Given** un panel abierto, **When** la consulta de estado falla, **Then** veo el
   error de la consulta y la información del proyecto sigue visible.
5. **Given** el panel de un proyecto, **When** se recarga, **Then** se vuelve a consultar
   el estado a la API.
6. **Given** el panel abierto de un proyecto, **When** cambio a otro proyecto, **Then** se
   vuelve a consultar el estado y, mientras llega, veo `Consultando…` en lugar del estado
   del proyecto anterior.
7. **Given** una consulta de estado que falla, **When** termina la consulta, **Then** veo
   `Sin datos` en la etiqueta de estado, además del mensaje de error.

## Fuera de alcance

- Recalcular el estado en el frontend.
- Persistir o cachear el estado.
- Implementar la pantalla independiente de estado (`/projects/:id/status`); su contenido
  queda para otra tarea.
- Cambiar el contrato de la API de HU-13.

## Trazabilidad

| Criterio | Implementación | Prueba |
|---|---|---|
| Estado visible | `Panel` + `obtenerEstadoProyecto` | test del panel para los tres valores |
| Consulta al abrir | `useEffect` del panel | test que verifica la llamada a `/projects/{id}/status` |
| Consulta al cambiar de proyecto | `useEffect` dependiente de `id` | test que navega a otro proyecto y verifica la nueva consulta |
| Reinicio al cambiar de proyecto | `key={id}` en `PanelLayout` | test que navega por el componente real y no ve el estado anterior |
| Carga del estado | etiqueta con `Consultando…` | test que cambia de proyecto y verifica el indicador |
| Error sin ocultar el proyecto | estado de error del panel | test del panel con fallo de la consulta |
| Estado no inventado | etiqueta con `Sin datos` | test que cambia a un proyecto cuya consulta falla |
| Sin lógica duplicada | la API calcula; el frontend solo muestra | revisión del código |
