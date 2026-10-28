# Research: Navegar entre proyectos desde la interfaz

**Feature**: #148 (cierra #141) | **Date**: 2026-10-27 | **Spec**: [spec.md](./spec.md)

## Scope

No quedan `NEEDS CLARIFICATION` en el contexto técnico. Todas las decisiones se resuelven con la
spec 022, el código actual del frontend y la constitución del proyecto. No se requiere
investigación externa ni comparación de librerías porque la historia prohíbe dependencias nuevas.

## Decisions

### D1. Componer el inicio sin modificar `crear.tsx`

- **Decision**: En `web/src/rutas.tsx`, la ruta índice pasa a ser hija de un layout pathless
  local `ComposicionInicio` que renderiza `<Outlet />` (el formulario lazy actual) y
  `<ListaProyectos />`.
- **Rationale**: El formulario sigue siendo el mismo módulo `features/proyecto/crear`; no se
  toca su código ni su prueba. `rutas.tsx` ya es el registro de rutas y es uno de los archivos
  permitidos por la spec.
- **Alternatives considered**:
  - Modificar `features/proyecto/crear.tsx` para agregar la lista: rechazado porque
    `crear.test.tsx` afirma la cantidad exacta de llamadas a `fetch` y el `GET /projects`
    rompería esas pruebas, además de mezclar una pantalla de feature con navegación transversal.
  - Crear una pantalla de inicio nueva que importe `crear`: rechazado porque una pieza de
    `compartido/` importaría una feature y rompería la separación del proyecto.
  - Renderizar la lista dentro de `App.tsx`: rechazado porque la lista solo corresponde a `/`,
    no a todas las rutas.

### D2. Módulo compartido `compartido/proyectos/`

- **Decision**: Concentrar API, hook, lógica pura, selector y lista en
  `web/src/compartido/proyectos/`.
- **Rationale**: La navegación es transversal al encabezado de proyecto y al inicio; no
  pertenece a una sola feature. Se reutiliza el cliente compartido y se evita duplicar la
  lectura del listado.
- **Alternatives considered**:
  - Colocarlo en `features/proyecto/`: rechazado porque el inicio y el layout compartido
    tendrían que importar una feature, algo que la prueba de aislamiento y la arquitectura
    desaconsejan.
  - Duplicar la lógica en cada pantalla: rechazado por DRY.

### D3. Un hook con unión discriminada de estados

- **Decision**: `useProyectos()` devuelve una unión discriminada
  `{ estado: 'cargando' } | { estado: 'listo'; proyectos: ProyectoResumen[] } |
  { estado: 'error'; mensaje: string }`. Usa `listarProyectos()` en un `useEffect` con bandera
  `vigente` para no actualizar estado tras el desmontaje.
- **Rationale**: El estado solicitado por la spec (cargando, listo, error) se representa sin
  combinaciones imposibles; la bandera evita la actualización tardía y es simple de probar.
- **Alternatives considered**:
  - Tres booleanos/valores separados (`cargando`, `error`, `proyectos`): rechazado porque permite
    estados contradictorios y más ramas de render.
  - Contexto global, store o SWR/React Query: rechazado por YAGNI y por la prohibición de
    dependencias nuevas. Selector y lista nunca se montan juntos.
  - `AbortController`: el cliente compartido no acepta una señal y modificarlo queda fuera del
    alcance; la bandera cubre el requisito funcional.

### D4. Validación defensiva de la respuesta

- **Decision**: Si `GET /projects` responde con éxito pero el resultado no es un arreglo, el hook
  pasa a `error` con un mensaje genérico. Los errores del cliente (`ApiError`,
  `ConnectionError`) muestran su mensaje.
- **Rationale**: La spec exige un estado de error ante respuestas no interpretables y no romper
  la pantalla. El cliente ya cubre JSON inválido, pero no un `200` con cuerpo vacío o un objeto.
- **Alternatives considered**:
  - Validar con una librería de esquemas: rechazado por dependencia nueva y por el tamaño del
    contrato (solo `id` y `nombre`).
  - Confiar ciegamente en el tipo: rechazado porque un cuerpo vacío dejaría `undefined` y
    rompería el render.

### D5. Selector nativo y accesible

- **Decision**: Usar `<select>` con etiqueta accesible "Proyecto"; opciones con
  `<nombre> (#<id>)`; `value` igual al `:id` de la URL; opción sintética `Proyecto <id>` si el
  actual no está en la lista. Al cambiar, navegar con `rutaDeSeccion`; elegir el actual no hace
  nada. El enlace "Nuevo proyecto" se renderiza siempre como enlace real. No se usa el componente
  `Button` porque esta historia solo navega con enlaces y un `<select>`.
- **Rationale**: Un `<select>` nativo es operable por teclado, accesible y no requiere
  dependencias. Mantiene el foco en la navegación y no en un menú propio.
- **Alternatives considered**:
  - Menú desplegable propio: rechazado por más código, manejo de foco y teclado sin beneficio.
  - Lista de enlaces en el encabezado: más ruidosa con muchos proyectos; el issue pide un
    selector.
  - Convertir "Nuevo proyecto" en `<Button>`: rechazado porque es una navegación; un enlace real
    conserva abrir en otra pestaña y la semántica. Se reutiliza la clase visual `.button` como en
    el panel.

### D6. Estados del selector y de la lista

- **Decision**:
  - Selector: `LoadingIndicator` mientras carga y `ErrorMessage` si falla; en ambos estados el
    `<select>` se muestra deshabilitado con la única opción `Proyecto <id>` del proyecto actual.
    Cuando la carga se resuelve, el `<select>` se habilita con la lista completa. El enlace
    "Nuevo proyecto" y los enlaces de sección siguen visibles y operativos.
  - Lista: `LoadingIndicator` mientras carga; `ErrorMessage` si falla; invitación a crear si está
    vacía; enlaces si hay proyectos.
- **Rationale**: Son los tres estados pedidos por la spec y no bloquean la navegación ni el
  formulario. Mantener el selector deshabilitado con el proyecto actual conserva el contexto sin
  ofrecer opciones que todavía no se pueden elegir. El encabezado mantiene los enlaces de
  sección fuera del selector.
- **Alternatives considered**:
  - Ocultar el encabezado completo ante error: rechazado; la spec exige que la navegación siga
    funcionando.
  - Ocultar el selector ante carga o error: rechazado porque la spec pide que el proyecto actual
    siga visible y que el contexto no se pierda.
  - Mostrar el `<select>` habilitado sin datos o con la lista parcial: rechazado porque ofrecería
    opciones no disponibles mientras la carga no terminó.

### D7. Ruta y etiqueta puras

- **Decision**:
  - `rutaDeSeccion(pathname, idNuevo)`: divide el pathname, toma el tercer segmento (índice 2
    después de `projects` e `id`) o `proyecto` si no existe, y arma `/projects/<idNuevo>/<sección>`.
  - `etiquetaProyecto({ id, nombre })`: devuelve `${nombre} (#${id})`.
- **Rationale**: La regla del issue es explícita: conservar solo el primer tramo después del id.
  Las funciones son deterministas, no dependen de React y se prueban con Vitest.
- **Alternatives considered**:
  - Usar `useParams`/`useLocation` para armar la ruta dentro del selector: rechazado porque
    mezcla lógica pura con React y dificulta las pruebas.
  - Conservar la ruta completa: rechazado por la spec (las sub-rutas no se trasladan).

### D8. Reinicio por `:id`

- **Decision**: Montar `<Outlet key={id} />` en `LayoutProyecto` y conservar las `key` existentes
  de las pantallas.
- **Rationale**: Garantiza que cualquier pantalla de proyecto, incluidas las que hoy no tienen
  una `key` propia, se remonte al cambiar de proyecto y no mezcle estado. El selector queda fuera
  del `Outlet`, por lo que no se remonta.
- **Alternatives considered**:
  - Confiar solo en las `key` actuales de panel y Sprints: rechazado porque la spec pide el
    reinicio por `:id` para las pantallas en general.
  - Navegar con recarga completa: rechazado; rompe la SPA y la experiencia.

### D9. Estrategia de pruebas y orden de commits

- **Decision**: Primero un commit de infraestructura que simula `GET /projects` en
  `rutas.test.tsx` (respuesta nueva por llamada y despacho por URL, para no alterar los `GET
  /projects/:id` y `/members` existentes); después RED/GREEN de lógica pura, hook, selector,
  lista e integración; al final eliminación de `Inicio.tsx` y verificación. Un commit por fase,
  REFACTOR solo si hay algo que limpiar.
- **Rationale**: La suite existente queda verde antes y después de que el encabezado empiece a
  pedir la lista, y cada comportamiento tiene evidencia RED/GREEN. Es la secuencia pedida por el
  equipo y compatible con la constitución.
- **Alternatives considered**:
  - Cambiar el encabezado y las pruebas al mismo tiempo: rechazado porque oculta si el mock
    mantiene la suite existente verde.
  - Un solo commit con todo: rechazado por la constitución y por la evidencia TDD.

### D10. Sin CSS nuevo

- **Decision**: Reutilizar las clases y el marcado existentes; no agregar un archivo de estilos ni
  modificar `App.css`.
- **Rationale**: La spec limita los archivos permitidos y no pide un rediseño. El navegador
  aplica estilos base suficientes; el foco es la navegación y los estados.
- **Alternatives considered**:
  - Agregar `proyectos.css`: fuera de los archivos permitidos por la spec.
  - Modificar `App.css`: tampoco está permitido.

### D11. Sin caché compartida ni actualización en vivo

- **Decision**: Cada consumidor pide la lista al montarse. No hay caché, polling ni refresco
  automático después de crear un proyecto.
- **Rationale**: La spec lo declara fuera de alcance y el costo de un pedido único es bajo.
- **Alternatives considered**:
  - Contexto con caché e invalidación: rechazado por YAGNI y por agregar estado global sin
    necesidad.
  - Polling: rechazado por carga innecesaria y porque no se pide tiempo real.

## Open Questions

Ninguna. Las decisiones de producto y de UI ya están en la spec 022; las técnicas se resuelven en
este documento.
