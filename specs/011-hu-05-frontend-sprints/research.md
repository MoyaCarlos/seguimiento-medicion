# Research: Pantalla de Sprints

Sin incógnitas abiertas; las decisiones del Scrum Master para el plan se registran acá.

## 1. Fechas
- **Decision**: `<input type="date">` nativo, vía `Field`.
- **Rationale**: entrega `AAAA-MM-DD`, que es lo que pide la API; cero dependencias.
- **Alternatives**: librería de fechas o selector propio (fuera de alcance, requeriría tocar `package.json`).

## 2. Estado de los Sprints
- **Decision**: `useState` local en la pantalla; cada respuesta de la API reemplaza al Sprint por `id`.
- **Rationale**: la API no lista Sprints (#110); solo hay que recordar los de la sesión. Cumple KISS/YAGNI.
- **Alternatives**: Context o librería de datos (estado compartido innecesario); `localStorage` (contradice FR-009, el estado "de la sesión" se pierde al recargar).

## 3. Orden dentro de cada sección
- **Decision**: por `id` ascendente.
- **Rationale**: estable, coincide con el orden de creación y permite distinguir varios Finalizados.

## 4. Clasificación de errores
- **Decision**: función pura `clasificarError(error, campos)` en `errores.ts`.
- **Rationale**: es la única lógica de la pantalla que decide qué se ve y dónde; pura y testeable con RED/GREEN. Un `field` que el formulario no tiene (p. ej. `estado`) se trata como aviso.
- **Alternatives**: decidir en cada componente (duplicaría la regla en formulario, cierre y creación).

## 5. Aislamiento del cierre
- **Decision**: `BotonCerrar.tsx` contiene el botón y la llamada `cerrarSprint`.
- **Rationale**: #117 (cierre con arrastre) puede cambiar el contrato; el cambio queda en un archivo. Ver FR-007.
- **Nota**: `api.ts` reúne las llamadas de crear e iniciar; la de cerrar vive junto a su botón por esta razón.

## 6. Estilos
- **Decision**: `sprints.css` propio, importado desde `index.tsx` (como la pantalla de historias de HU-01); reutiliza las clases de las piezas comunes.
- **Rationale**: no se edita `index.css` ni `compartido/`.

## 7. Pruebas
- **Decision**: Vitest + Testing Library; `fetch` simulado con `vi.stubGlobal`, igual que los tests del cliente. Las pruebas de pantalla montan `Component` dentro de un `MemoryRouter`/ruta `projects/:id/sprints`.
