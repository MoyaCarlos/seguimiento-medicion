# Quickstart: verificar la pantalla de Sprints

Requisitos: Node 24 LTS (`web/.nvmrc`) y el backend en `:8080` (`go run ./cmd/api`).

## 1. Calidad (SC-006)
```bash
cd web && npm test && npm run lint && npm run build
```
Esperado: 0 errores. `git diff --stat main` solo muestra `web/src/features/sprints/` y `specs/011-*/`.

## 2. Preparar datos
Con los `curl` de [`docs/curl-frontend-sprint-2.md`](../../docs/curl-frontend-sprint-2.md), crear un proyecto y anotar su id. No crear el Sprint por `curl`: se crea desde la pantalla.

## 3. Crear, iniciar y cerrar (US1–US3, SC-001, SC-002)
```bash
cd web && npm run dev
```
Abrir `http://localhost:5173/projects/<id>/sprints`.
1. Se ven los Sprints existentes del proyecto en sus secciones (o las tres vacías si no hay), sin texto de limitación (spec 015).
2. "Crear Sprint": aparece "Sprint #N" en Pendiente, sin Goal ni fechas.
3. Completar Goal y fechas, iniciar: el Sprint pasa a Activo con Goal y fechas.
4. "Cerrar Sprint": pasa a Finalizado, sin acciones.

## 4. Errores (US1–US3, SC-003)
- "Crear Sprint" con un Pendiente ya existente: aviso arriba de la pantalla con el mensaje de la API.
- Iniciar con fecha de fin anterior a la de inicio: el mensaje aparece junto a "Fecha de fin".
- Iniciar con otro Sprint Activo en el proyecto: aviso arriba del formulario de ese Sprint (campo `estado`).
- Apagar el backend y pulsar cualquier acción: aviso de conexión en la acción que falló.

## 5. Recarga (US4)
Recargar la página: los Sprints vuelven a cargarse desde la API (spec 015).

## 6. Aislamiento (SC-005)
`git diff --stat` del PR: ningún archivo fuera de la carpeta de la feature y de `specs/011-*/`; `BotonCerrar.tsx` es el único archivo que conoce la llamada de cierre.
