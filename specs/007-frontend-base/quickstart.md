# Quickstart: verificar la base

Requisitos: Node 24 LTS (ver `web/.nvmrc`). Para los pasos 3 y 4, el backend corriendo en `:8080` (`go run ./cmd/api`).

## 1. Compila y lint (SC-001, FR-001/002/003)
```bash
cd web && npm install && npm run build && npm run lint
```
Esperado: 0 errores. No queda contenido de la plantilla de Vite (contador, logos).

## 2. Pruebas (FR-008, FR-010, FR-017, FR-018, SC-005)
```bash
cd web && npm test
```
Esperado: pasan los casos del cliente de FR-010, las pruebas de las piezas, las de rutas/navegación y la que impide importar entre features.

## 3. Conexión con la API (US1, FR-004)
1. Ejecutar los `curl` de `docs/curl-frontend-sprint-2.md` (versión bash o `cmd`) y anotar el id del proyecto.
2. `cd web && npm run dev` y, en otra terminal: `curl -H "Accept: application/json" http://localhost:5173/projects/<id>`.

Esperado: JSON del proyecto, sin error de permisos.

## 4. Recarga con id en la URL (SC-006, FR-013)
Abrir `http://localhost:5173/projects/<id>/backlog` en el navegador y recargar. Esperado: aparece el placeholder de Historias con el mismo proyecto, no un 404 del backend; el enlace "Historias" del menú aparece marcado como activo.

## 5. Sumar una pantalla (US2, SC-003)
Seguir `web/README.md` con una pantalla de prueba. Esperado: `git diff --stat` fuera de su carpeta muestra solo `rutas.tsx` y `navegacion.ts`, una línea cada uno. Descartar la pantalla de prueba.

## 6. Independencia (FR-020)
Agregar en una feature un import de otra feature. Esperado: `npm test` falla (lo detecta la prueba de aislamiento; `npm run lint` no lo detecta). Revertir.

## 7. Página no encontrada (FR-014)
Abrir `http://localhost:5173/no-existe`. Esperado: aviso de "página no encontrada".
