# Research: Navegar entre proyectos desde la interfaz

**Feature**: #148 | **Date**: 2026-10-27 | **Spec**: [spec.md](./spec.md)

## Scope

No quedan `NEEDS CLARIFICATION` en el contexto técnico. Todas las decisiones se resuelven con la
spec 022, el código actual del frontend y la constitución del proyecto. No se requiere
investigación externa ni comparación de librerías porque la historia prohíbe dependencias nuevas.

## Decisions

### D1. Componer el inicio sin modificar `crear.tsx`

- **Decision**: En `web/src/rutas.tsx`, la ruta índice pasa a ser hija de un layout pathless
  `ComposicionInicio` (en `compartido/proyectos/ComposicionInicio.tsx`, para no exportar un
  componente junto con la configuración de rutas) que renderiza `<Outlet />` (el formulario lazy actual) y
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
    después de `projects` e `id`); si no existe o no está en `navegacion` (`web/src/navegacion.ts`)
    usa `proyecto`; arma `/projects/<idNuevo>/<sección>`.
  - `etiquetaProyecto({ id, nombre })`: devuelve `${nombre} (#${id})`.
- **Rationale**: La regla del issue es conservar solo el primer tramo después del id. Decisión del
  Scrum Master (2026-11-02): si ese tramo no es una sección del menú (p. ej. `votacion/:historiaId`)
  se cae en `proyecto`, sin tabla de rutas, para no llegar a una pantalla inexistente ni trasladar
  el id de una historia de otro proyecto.
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
  lista e integración; al final el nombre accesible del `<nav>` y la verificación. Un commit por fase,
  REFACTOR solo si hay algo que limpiar.
- **Rationale**: La suite existente queda verde antes y después de que el encabezado empiece a
  pedir la lista, y cada comportamiento tiene evidencia RED/GREEN. Es la secuencia pedida por el
  equipo y compatible con la constitución.
- **Alternatives considered**:
  - Cambiar el encabezado y las pruebas al mismo tiempo: rechazado porque oculta si el mock
    mantiene la suite existente verde.
  - Un solo commit con todo: rechazado por la constitución y por la evidencia TDD.

### D10. Estilo en `App.css`

- **Decision**: Reutilizar las clases y el marcado existentes. Primero se agregó a `App.css` solo la
  regla `.field select` equivalente a la de `.field input` (borde, radio, padding y foco); después,
  con el OK del Scrum Master (2026-11-03), el estilo del selector del encabezado y del inicio
  (ver D13), siempre en `App.css` y sin dependencias nuevas. Decisión del
  Scrum Master (2026-11-03): la restricción original de "sin CSS nuevo" no estaba en `spec.md`
  sino en esta decisión y en el plan; la lista de archivos permitidos ahora incluye `App.css`.
- **Rationale**: Sin la regla el `<select>` nativo se ve distinto al resto de los campos.
- **Alternatives considered**:
  - Agregar `proyectos.css`: un archivo más para una sola regla compartida con los campos.
  - No agregar CSS: deja el selector sin estilo.

### D13. Estilo del selector y del inicio: alcance, atajos de CSS y un defecto

- **Decision (selector)**: píldora compacta con flecha (la etiqueta "Proyecto" queda solo para
  lectores de pantalla con `.visually-hidden`), en la fila del menú de secciones y a su izquierda,
  con el botón secundario "+ Nuevo proyecto" a su lado; deshabilitado en carga y error, con el
  aviso debajo de la fila; en celular (390 px) ocupa todo el ancho y "Nuevo proyecto" va debajo.
  Se usan solo las variables de `App.css`.
- **Decision (inicio)**: solo `/` se ensancha (hasta 1320 px desde 1200 px de viewport), sin cambiar
  `.app-main` ni las otras pantallas; formulario de 480 a 520 px y la lista en el resto; fechas
  lado a lado y botones en fila, que se apilan en celular; las dos columnas se apilan por debajo de
  1100 px.
- **Atajos de CSS** (para no tocar `crear.tsx`, `BotonDatosDePrueba.tsx` ni el shell; se retiran
  en el rediseño):
  - Contenedor más ancho que `.app-main` con un margen izquierdo negativo calculado
    (`calc((100% - min(1320px, 100vw - 3rem)) / 2)`), solo desde 1200 px.
  - Encabezado superior alineado con `.app-shell:has(.inicio-composicion) .app-header`; si el
    navegador no soporta `:has()` el encabezado conserva su ancho de siempre.
  - La tarjeta del formulario pasa a ser una grilla y tanto el `<form>` como el bloque de "Cargar
    datos de prueba" usan `display: contents`, para que "Crear proyecto" (que está dentro del
    formulario) y "Cargar datos de prueba" (que está fuera) compartan una fila. Riesgo conocido:
    `display: contents` en un formulario puede hacer que lectores de pantalla antiguos pierdan su
    rol; no se probó con un lector de pantalla.
  - Las reglas dependen del marcado actual: el bloque de datos de prueba es el último hijo de la
    tarjeta (`> div:last-child`). Una prueba de componente en `rutas.test.tsx` fija ese contrato.
- **Defecto encontrado y arreglo (2026-11-03)**: con el aviso de "Cargar datos de prueba" ("los
  datos de prueba ya están cargados" más el enlace "Abrir el proyecto de prueba", o un error), el
  bloque de datos de prueba quedaba en una sola celda de la grilla junto a "Crear proyecto": el
  aviso se mostraba en la misma fila que los botones y "Crear proyecto" se estiraba a toda la
  altura (125 px en lugar de 42 px). Arreglo: el bloque usa `display: contents`, los dos botones
  van con `order: 1` y todo lo que no es botón (aviso, error o indicador de carga) ocupa su propia
  fila debajo (`grid-column: 1 / -1`, `order: 2`), y la tarjeta alinea sus celdas arriba
  (`align-items: start`) y el botón de datos de prueba no se estira (`justify-self: start`).
  Verificado en escritorio (1280 px) y celular (390 px) sin aviso, con "ya cargados" (respuesta real
  409 del backend) y con error (respuesta 500 simulada de `POST /demo/datos`).
- **Deuda para el rediseño general (issue #230)**: los tres atajos de arriba (margen negativo
  calculado de `.inicio-composicion`, `display: contents` en el formulario y en el bloque de datos
  de prueba, y la regla `:has()` del encabezado) llevan un comentario en `App.css` que los marca
  como deuda. Además, el texto del aviso empieza en minúscula ("los datos de prueba ya están
  cargados") y su enlace "Abrir el proyecto de prueba" se parte en dos líneas en la tarjeta; ambas
  cosas quedan para #230.
- **Prueba del aviso**: se pudo agregar sin tocar `crear.tsx` (en `rutas.test.tsx`, fija el
  marcado del que dependen las reglas). No tiene ciclo RED/GREEN porque no cambia código de
  producción: el arreglo es solo CSS y el marcado ya cumplía el contrato; se comprobó que falla
  si el aviso sale del bloque de los botones (mutación manual) y se commiteó junto con el estilo
  de `/`.
- **Fuera de este arreglo (a propósito)**: el texto del aviso ("los datos de prueba ya están
  cargados", sin mayúscula inicial) y el color gris del mensaje de error dentro de la tarjeta
  (`.hu04-tarjeta p` pisa a `.message--error`) ya eran así antes de esta historia; se resuelven en
  el rediseño (issue #230).
- **Alternatives considered**:
  - Modificar `crear.tsx` o `BotonDatosDePrueba.tsx` para agrupar los botones: rechazado, fuera de
    los archivos permitidos de esta historia.
  - Ensanchar `.app-main` para todas las pantallas: rechazado por el Scrum Master (opción A).

### D12. Tarjetas, buscador y dos columnas en el inicio

- **Decision**: La lista del inicio sigue la referencia visual del prototipo de Scrum Master
  (2026-11-03): dos columnas (formulario a la izquierda, lista a la derecha; apiladas en celular),
  encabezado con contador, tarjeta enlace con nombre, id, descripción y fechas, y un buscador por
  nombre del lado del cliente. Se usan los colores y variables reales de `App.css`.
  No se toma del prototipo el estado, el Scrum Master, la barra de avance ni el orden inverso.
- **Rationale**: Es un cambio de presentación y de filtrado local sobre datos que `GET /projects`
  ya devuelve; el estado lo calcula el backend y duplicarlo en el frontend rompería DRY.
- **Alternatives considered**:
  - Buscar en el servidor: rechazado, no hay endpoint ni se permite cambiar la API.
  - Filtrar también por descripción o id: rechazado, la búsqueda es solo por nombre.
  - Fechas con otro formato que el del panel: rechazado por consistencia con `Panel.tsx`.

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
