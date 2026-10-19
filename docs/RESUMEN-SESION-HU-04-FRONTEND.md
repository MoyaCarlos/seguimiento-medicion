# Resumen de trabajo — HU-04 Frontend

Fecha: 2026-10-14  
Historia: **HU-04 — Creación de Proyecto y Asignación de Equipo**  
Sub-issue: **#36 — Frontend**  
Sprint: **Sprint 2 — Interfaz**  
Rama: `feature/HU-04-frontend`

Este documento resume el trabajo realizado hasta el momento en la interfaz de
HU-04. El maquetado actual es provisorio porque el esqueleto visual definitivo
se está coordinando con el resto del equipo.

## Alcance de la sub-issue

La sub-issue #36 contempla las siguientes vistas:

- Pantalla principal con el formulario para crear un proyecto.
- Panel principal del proyecto después de crearlo.
- Vista de configuración para asignar integrantes y roles.
- Vista de edición del proyecto con sus errores de validación.

En esta primera etapa se trabajó sobre el formulario de creación y el panel
posterior a la creación.

## Trabajo realizado

### Formulario de creación

Se reemplazó la pantalla inicial de Vite por un formulario React con los
siguientes campos:

- Nombre del proyecto.
- Descripción opcional.
- Fecha de inicio.
- Fecha de fin.
- Nombre del creador.

El formulario utiliza validación del lado del cliente antes de realizar la
petición:

- El nombre del proyecto es obligatorio.
- El nombre del creador es obligatorio.
- La fecha de inicio es obligatoria.
- La fecha de fin es obligatoria.
- La fecha de fin debe ser posterior a la fecha de inicio.

Los errores se muestran junto al campo correspondiente y se utilizan atributos
de accesibilidad como `aria-invalid`, `aria-describedby` y `role="alert"`.

### Integración con la API

La creación se conecta con el endpoint existente:

```http
POST /projects
```

El frontend transforma los nombres de sus campos al contrato del backend:

```json
{
  "nombre": "Proyecto Demo",
  "descripcion": "Descripción del proyecto",
  "creador": "Ana",
  "fecha_inicio": "2026-10-15",
  "fecha_fin": "2026-11-15"
}
```

También se contempla el manejo de errores devueltos por la API. Mientras la
petición está en curso, el botón queda deshabilitado y cambia su texto a
`Creando…`.

### Panel posterior a la creación

Cuando la API responde correctamente, la pantalla muestra el panel principal
del proyecto creado con:

- Nombre.
- Descripción.
- Fecha de inicio.
- Fecha de fin.

También se agregó una acción para volver al formulario y crear otro proyecto.
La redirección es actualmente una transición de estado dentro de la aplicación;
se podrá adaptar a un sistema de rutas cuando el equipo defina el maquetado y
la navegación final.

## Archivos involucrados

### Componentes y estilos

- `web/src/App.tsx`: formulario y panel principal provisorio.
- `web/src/App.css`: estilos de la pantalla actual.
- `web/src/index.css`: estilos globales mínimos.

### Lógica y API

- `web/src/projectForm.ts`: tipos y validación del formulario.
- `web/src/api.ts`: función para crear proyectos mediante la API.
- `web/src/projectForm.test.ts`: pruebas unitarias de las validaciones.

### Configuración

- `web/package.json`: script `npm test` y dependencia de Vitest.
- `web/package-lock.json`: lockfile actualizado.
- `web/vite.config.ts`: proxy de `/projects` hacia `http://localhost:8080`.

## Desarrollo guiado por pruebas

Se aplicó el ciclo TDD para la lógica de validación:

1. **RED:** se escribieron pruebas para campos obligatorios, fechas inválidas
   y formulario válido.
2. **GREEN:** se implementó `validateProjectForm` hasta hacer pasar las tres
   pruebas.
3. **REFACTOR:** se tiparon explícitamente los datos del formulario y los
   errores para preservar la seguridad de tipos de TypeScript.

Las pruebas no dependen de la implementación visual y podrán conservarse
cuando el formulario se distribuya en componentes definitivos.

## Verificaciones realizadas

Se ejecutaron correctamente:

```text
npm test
npm run build
npm run lint
go test ./...
```

Resultados:

- 1 archivo de pruebas frontend ejecutado.
- 3 pruebas frontend aprobadas.
- Build de TypeScript y Vite aprobado.
- Oxlint aprobado.
- Suite completa de Go aprobada.

También se verificó manualmente en navegador el flujo completo:

1. Abrir la pantalla de creación.
2. Intentar enviar el formulario vacío.
3. Ver los errores de validación.
4. Completar los datos válidos.
5. Crear el proyecto mediante la API.
6. Visualizar el panel con los datos devueltos por el backend.

## Decisiones provisionales

El CSS y la estructura de `App.tsx` son una primera versión funcional. No se
consideran el diseño visual definitivo porque el equipo todavía está
terminando el esqueleto del frontend.

Cuando se incorpore el maquetado común, será necesario revisar:

- Componentes compartidos de navegación y layout.
- Sistema de rutas.
- Colores, tipografías y espaciados.
- Nombres visibles de los campos.
- Ubicación del formulario y del panel.
- Compatibilidad con la vista de configuración y la vista de edición.

La lógica de validación y la función de integración con la API se mantienen
separadas para facilitar esa integración.

## Pendiente dentro de #36

- Integrar el formulario en el maquetado definitivo.
- Incorporar el panel principal al layout final.
- Crear la vista de configuración de integrantes y roles usando:
  - `GET /projects/{id}/members`
  - `POST /projects/{id}/members`
- Crear la vista de edición usando:
  - `GET /projects/{id}`
  - `PUT /projects/{id}`
- Mostrar los errores de validación de la edición.
- Agregar pruebas para las nuevas vistas y flujos.
- Revisar la experiencia completa con el resto del equipo antes de cerrar la
  sub-issue.

## Estado actual

La primera parte funcional de la sub-issue está implementada y verificada en la
rama `feature/HU-04-frontend`. El trabajo queda preparado para recibir el
maquetado definitivo sin descartar la lógica ni las pruebas ya realizadas.
