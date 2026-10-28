# Feature Specification: Estado del proyecto en el panel principal (HU-13 Frontend)

**Feature Branch**: `feature/HU-13-frontend-panel`

**Issue**: #41 (parte de #13)

**Created**: 2026-10-23

**Status**: Implementada, pendiente de revisión por PR.

**Dependencia**: el panel principal pertenece a la pantalla de HU-04
(`features/proyecto`, PR #121). Esta rama se apoya en esa pantalla y no duplica su
implementación.

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
| Error sin ocultar el proyecto | estado de error del panel | test del panel con fallo de la consulta |
| Sin lógica duplicada | la API calcula; el frontend solo muestra | revisión del código |
