# Quickstart: Navegar entre proyectos desde la interfaz

**Feature**: #148 (cierra #141) | **Date**: 2026-10-27 | **Spec**: [spec.md](./spec.md)

## Prerequisites

- Node.js `>=24.15.0` (ver `web/.nvmrc` y `web/package.json`).
- Backend Go disponible en `http://localhost:8080` para la validación manual.
- Al menos dos proyectos creados para probar el cambio de proyecto. Se pueden cargar con los
  comandos de [`docs/curl-frontend-sprint-2.md`](../../docs/curl-frontend-sprint-2.md) o desde
  la pantalla de inicio.
- Dependencias del frontend instaladas (`npm install` dentro de `web/`).

## Automated Validation

Desde `web/`:

```bash
npm test
npm run lint
npm run build
```

Esperado:

- `npm test`: pasa la suite existente (incluida `rutas.test.tsx` con `GET /projects` simulado) y
  las nuevas pruebas de `ruta`, `useProyectos`, `SelectorProyecto` y `ListaProyectos`.
- `npm run lint`: termina sin errores.
- `npm run build`: `tsc -b` y `vite build` terminan sin errores.
- No se modifica ni se ejecuta una suite nueva de backend; el CI de backend sigue igual.

## Manual Validation

### 1. Lista de proyectos en el inicio

1. Levantar el backend: `go run ./cmd/api` desde la raíz.
2. Levantar el frontend: `npm run dev` desde `web/`.
3. Abrir `http://localhost:5173/`.
4. Verificar que debajo del formulario "Crear proyecto" aparece la lista de proyectos existentes,
   cada uno como `<nombre> (#<id>)`.
5. Activar el enlace de un proyecto.

**Esperado**: se abre `/projects/<id>/proyecto` con el panel de ese proyecto.

### 2. Estado vacío del inicio

1. Abrir `/` con una base sin proyectos.

**Esperado**: se ve la invitación a crear el primer proyecto y el formulario sigue disponible; no
se muestra un error.

### 3. Selector en el encabezado

1. Abrir `/projects/<id>/backlog` con dos proyectos cargados.
2. Abrir el selector "Proyecto" del encabezado.
3. Elegir el otro proyecto.

**Esperado**:

- La dirección pasa a `/projects/<otroId>/backlog` (se conserva la sección).
- El selector marca el proyecto nuevo.
- La pantalla del backlog se reinicia y no muestra estado del proyecto anterior.

### 4. Conservar solo la sección

1. Abrir `/projects/<id>/proyecto/edit`.
2. Elegir otro proyecto en el selector.

**Esperado**: se navega a `/projects/<otroId>/proyecto`, no a la sub-ruta de edición.

### 5. Enlace "Nuevo proyecto"

1. Abrir cualquier sección de un proyecto.
2. Pulsar "Nuevo proyecto" en el encabezado.

**Esperado**: se abre `/` con el formulario de creación y la lista de proyectos.

### 6. Estados de carga y error

1. Abrir una pantalla de proyecto y observar el selector mientras carga: se ve el indicador de
   carga y el enlace "Nuevo proyecto" sigue visible.
2. Detener el backend y recargar la pantalla de proyecto.

**Esperado**:

- El selector muestra el mensaje de error de la API.
- Los enlaces de sección (Proyecto, Historias, Sprints, Estado) siguen funcionando.
- "Nuevo proyecto" sigue visible y lleva a `/`.
- En `/`, el formulario de creación sigue disponible aunque la lista muestre el error.

### 7. Proyecto actual fuera del listado

1. Abrir `/projects/<id>/backlog` con un id que no exista en el listado (por ejemplo, un proyecto
   borrado o una lista incompleta).

**Esperado**: el selector muestra una opción `Proyecto <id>` seleccionada y permite cambiar a los
proyectos que sí están en el listado.

## References

- Contrato de UI y consumo de `GET /projects`: [`contracts/ui.md`](./contracts/ui.md)
- Entidades y estados de presentación: [`data-model.md`](./data-model.md)
- Criterios de aceptación: [`spec.md`](./spec.md)
