# Feature Specification: Rediseño visual del frontend (#230)

**Feature Branch**: `feature/230-rediseno-completo` (rama de integración y PR único del rediseño; la ola 1 salió primero en ramas `feature/230-ola-1-*`)

**Created**: 2026-11-11

**Status**: Implementado en la rama de integración (en revisión). Actualizada el 2026-11-15 para reflejar el alcance final: ola 1 y ola 2 en un solo PR

**Input**: User description: "Rediseño visual del frontend (issue #230): nueva identidad visual (azul marino y teal) y mejor aprovechamiento del ancho de pantalla, con el prototipo «Métrica Remake Sprint 4» (versión final) como referencia visual, sin cambiar la API ni el modelo de datos y sin agregar pantallas ni endpoints. Se entrega por olas."

## Objetivo

Darle a la aplicación una identidad visual nueva y un armazón de navegación que aproveche el ancho de la pantalla, sin cambiar lo que la aplicación hace. El prototipo «Métrica Remake Sprint 4» (versión final) es la guía visual y de estructura: no se copia su código ni sus datos, y lo que el prototipo marca como «Propuesta» o con «*» (datos que la API no entrega) no forma parte de este rediseño.

## Actualización del 15/11 (alcance final)

El rediseño se entrega en **un solo PR de integración** (`feature/230-rediseno-completo`) que incluye las dos olas. Cuando se escribió esta spec la ola 2 quedaba fuera porque esperaba #247, HU-21 y cerrar/reabrir; esas historias ya están en `main` y la rama las integró (#321, #322, #324, #315, #319, #320, #323, #325 y #326, con la regla «queda el markup del rediseño y se reinyecta lo de `main`»). Las decisiones de la sección «Fuera de alcance» marcadas como tomadas el 11 y el 12/11 fueron **revisadas por el equipo** al adoptar el prototipo final (V10); abajo figura el estado vigente. Lo que quedó escrito como alcance de la ola 2 se lee como histórico.

## Alcance por olas

### Ola 1 (esta entrega)

| PR | Parte | Contenido |
|----|-------|-----------|
| 1 | Armazón | Barra lateral, barra superior, estilos globales, selector de proyecto, menú de usuario e íconos |
| 2 | Acceso | Pantallas Ingresar y Registro |
| 3 | Seguimiento | Pantallas Estado, Dashboard y Defectos |

### Ola 2 (incluida en el PR de integración)

Proyectos (inicio, con la sección «Proyectos cerrados»), Panel del proyecto, Historias (lista, cajón «Nueva historia», paneles de Editar y Asignar a Sprint, editor de Markdown), Sprints, Votación, Editar proyecto (con la zona «Cerrar proyecto»), Equipo, Historial de Sprints y su detalle. Originalmente esperaban el frontend de HU-21 (#247) y cerrar/reabrir proyecto (#300); al estar en `main`, se rediseñan en este PR conservando su comportamiento.

### Fuera de alcance

- Lo que la API no entrega: estado, avance y Scrum Master en las tarjetas de proyecto; prioridad de cada historia en Estado. (El rol sí se muestra en las tarjetas porque HU-21 ya lo entrega, y la «Cerrado el dd/mm/aaaa» sale de `fecha_cierre`.)
- **Excepción a FR-015, decidida por el equipo:** la barra de avance del cronograma del Panel es solo presentación; muestra el porcentaje de tiempo transcurrido entre la fecha de inicio y la de fin que la API ya entrega (fechas locales, «hoy» inyectable en las pruebas). No es una métrica y no cambia ningún cálculo.
- Lo que el prototipo marca como «Propuesta»: eliminar proyecto / «Zona de peligro» (#232, necesita endpoint nuevo; queda para después de la entrega) y todas las etiquetas «Propuesta» y «Solo del prototipo». Dejaron de ser «Propuesta» y se implementan porque ya están en `main`: «Repetir contraseña», permisos por rol (HU-21), cerrar y reabrir proyecto (#322), quitar integrante (#323) y cambiar nombre (#325).
- Decisiones del equipo: las pantallas de acceso van a pantalla completa, sin cabecera; sin avisos flotantes (el aviso de datos de prueba cargados es en línea).
- Decisiones del 11 y el 12/11 que el equipo **revirtió** al adoptar el prototipo final (V10), y que hoy sí se implementan: ventana modal para «Nuevo proyecto» (ruta `/projects/new`), barra de pestañas inferior en pantallas de 900 px o menos, tarjeta «Equipo» en el Panel y menú de usuario desplegable (con el rol en el proyecto y «Cambiar nombre»).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Navegar con el armazón nuevo (Priority: P1)

Como integrante de un proyecto quiero una barra lateral con las secciones y una barra superior con mi ubicación y mi sesión, para moverme por la aplicación aprovechando el ancho de la pantalla.

**Why this priority**: Todas las pantallas viven dentro del armazón; es la base del resto del rediseño.

**Independent Test**: Con sesión iniciada, abrir cualquier pantalla de un proyecto y recorrer las secciones, cambiar de proyecto y cerrar sesión solo con el teclado.

**Acceptance Scenarios**:

1. **Given** una sesión iniciada y una pantalla de un proyecto, **When** se muestra la aplicación, **Then** hay una barra lateral «Navegación principal» con la marca, el enlace «Proyectos», el selector «Proyecto», el enlace «Nuevo proyecto» y la navegación «Secciones del proyecto» con la sección actual marcada.
2. **Given** una sesión iniciada y una pantalla de un proyecto, **When** se muestra la barra superior, **Then** la «Ruta de navegación» muestra «Proyectos» como enlace al inicio y la sección actual, y el menú de usuario muestra el nombre y el botón «Cerrar sesión».
3. **Given** una sesión iniciada en el inicio, **When** se muestra la barra lateral, **Then** no aparecen el selector ni las secciones del proyecto.
4. **Given** una pantalla angosta, **When** la persona presiona «Abrir menú», **Then** la barra lateral se abre como cajón; **When** elige una sección, presiona Escape o presiona «Cerrar menú», **Then** el cajón se cierra.
5. **Given** una persona sin sesión, **When** abre Ingresar o Registro, **Then** la pantalla ocupa todo el ancho, sin barra lateral ni cabecera, con el panel de presentación (con la marca) a la izquierda y el formulario a la derecha, y el formulario ofrece el enlace a la otra pantalla («Registrate» o «Ingresá»).

---

### User Story 2 - Ingresar y registrarse con la identidad nueva (Priority: P2)

Como persona que usa la aplicación quiero pantallas de acceso claras, con la opción de ver la contraseña que escribo, para entrar sin errores de tipeo.

**Why this priority**: Es la primera impresión de la aplicación y no depende de otras historias en curso.

**Independent Test**: Abrir Ingresar y Registro sin sesión, alternar la visibilidad de la contraseña e ingresar.

**Acceptance Scenarios**:

1. **Given** la pantalla Ingresar o Registro, **When** la persona presiona «Mostrar contraseña», **Then** la contraseña se ve como texto y el botón queda marcado como presionado (el nombre del botón no cambia); **When** lo presiona de nuevo, **Then** vuelve a ocultarse.
2. **Given** la pantalla Ingresar, **When** se muestra, **Then** conserva el formulario, el aviso «Tu sesión venció» cuando corresponde, el enlace a Registro y los datos de prueba, con el comportamiento actual.
3. **Given** la pantalla Registro, **When** se muestra, **Then** cada campo indica la regla que ya aplica el servidor (usuario: de 3 a 30 caracteres, letras sin tildes, números, punto, guion y guion bajo; contraseña: de 8 a 128 caracteres); **When** el servidor rechaza un campo, **Then** el mensaje de error reemplaza esa ayuda.

---

### User Story 3 - Ver Estado, Dashboard y Defectos con la identidad nueva (Priority: P3)

Como integrante quiero que las pantallas de seguimiento usen el ancho disponible y la identidad nueva, para leer el avance de un vistazo.

**Why this priority**: Son pantallas de solo lectura o estables, que ninguna historia en curso modifica.

**Independent Test**: Abrir las tres pantallas con los datos de prueba y compararlas con el prototipo.

**Acceptance Scenarios**:

1. **Given** un proyecto con historias, **When** se abre Estado, **Then** se ve el total del proyecto con una barra que reparte las historias por estado y los tres grupos (Nueva, En progreso, Completada) en columnas, siempre los tres aunque estén en cero.
2. **Given** un proyecto con Sprints cerrados, **When** se abre Dashboard, **Then** se ven el gráfico y los indicadores actuales con el estilo nuevo.
3. **Given** un proyecto con defectos, **When** se abre Defectos, **Then** el formulario y el listado se ven lado a lado en pantallas anchas y la fecha de registro se muestra como dd/mm/aaaa; sin Sprint activo el formulario sigue oculto.

---

### Edge Cases

- Pantalla de 320 px de ancho: ninguna pantalla genera scroll horizontal.
- Nombre de proyecto o de usuario muy largo: se recorta sin romper la barra lateral ni la barra superior.
- El listado de proyectos falla o está cargando: el selector queda deshabilitado con «Proyecto <id>», se ve el aviso y las secciones siguen navegando.
- Sin conexión con el servidor o sesión en verificación: no se muestra la barra lateral.
- Una ruta de proyecto que no es del menú (votación, edición): la ruta de navegación muestra la sección del menú a la que pertenece.

## Requirements *(mandatory)*

### Functional Requirements

Reglas en formato EARS.

- **FR-001** (Ubicuo): El rediseño no debe cambiar la API, el modelo de datos, las rutas ni los textos y nombres accesibles existentes, salvo los que esta spec enumera.
- **FR-002** (Basado en estado): Mientras haya una sesión iniciada, la aplicación debe mostrar la barra lateral «Navegación principal» y la barra superior.
- **FR-003** (Basado en estado): Mientras la ruta sea de un proyecto, la barra lateral debe mostrar el selector «Proyecto» (botón con lista de opciones (listbox personalizado, `aria-haspopup="listbox"`) con «nombre (#id)» y la opción actual marcada), el botón «Nuevo proyecto» fuera de la lista y la navegación «Secciones del proyecto» con la sección actual marcada como página actual.
- **FR-004** (Ubicuo): La barra superior debe mostrar una navegación «Ruta de navegación» y el botón «Menú de usuario» con el nombre de la persona; al abrirlo muestra el panel «Cuenta» con el nombre, el rol en el proyecto (solo dentro de un proyecto), «Cambiar nombre» y «Cerrar sesión». El pie de la barra lateral repite el nombre, el rol y un botón de salida.
- **FR-005** (Basado en evento): Cuando la persona presiona «Abrir menú» en una pantalla angosta, la aplicación debe abrir la barra lateral como cajón; cuando navega, presiona Escape o presiona «Cerrar menú», debe cerrarlo.
- **FR-006** (Basado en estado): Mientras no haya sesión iniciada, la aplicación no debe mostrar la barra lateral ni la cabecera, y las pantallas de acceso deben ofrecer en el formulario el enlace a la otra pantalla («Registrate» o «Ingresá»).
- **FR-007** (Ubicuo): El listado de proyectos se debe pedir una sola vez al montar una pantalla de proyecto y no al cambiar de sección.
- **FR-008** (Basado en evento): Cuando la persona presiona «Mostrar contraseña», el campo debe mostrar el texto y el botón debe quedar con `aria-pressed` verdadero, sin cambiar su nombre accesible; al presionarlo de nuevo debe volver al estado inicial.
- **FR-009** (Ubicuo): Las fechas visibles en las pantallas rediseñadas (ola 1 y ola 2) se deben mostrar como dd/mm/aaaa, incluida la fecha de cierre del proyecto.
- **FR-010** (Ubicuo): Los íconos deben ser decorativos (ocultos para lectores de pantalla) y no deben agregar dependencias.
- **FR-011** (Ubicuo): Todo control interactivo debe tener foco visible, ser operable con teclado y cumplir contraste AA.
- **FR-012** (No deseado): Si el ancho de la pantalla es de 320 px, entonces ninguna pantalla debe generar scroll horizontal.
- **FR-013** (Basado en estado): Mientras la pantalla mida 900 px o menos, la barra lateral no debe ocupar ancho del contenido.
- **FR-014** (Ubicuo): Las pantallas de la ola 2 se rediseñan en el PR de integración conservando su comportamiento y los textos y nombres accesibles existentes, salvo los que cada commit RED declara con su antes y después.
- **FR-015** (No deseado): Si un dato del prototipo no lo entrega la API, entonces no se debe mostrar ni calcular en el cliente.

## Success Criteria *(mandatory)*

### Measurable Outcomes

**Visuales**

- **SC-001**: A 320 px de ancho ninguna pantalla rediseñada tiene scroll horizontal.
- **SC-002**: En pantallas de 1280 px o más, el contenido usa hasta 1180 px de ancho además de la barra lateral.
- **SC-003**: Todas las fechas de las pantallas rediseñadas se leen como dd/mm/aaaa.
- **SC-004**: Todo el recorrido (navegar, cambiar de proyecto, abrir y cerrar el cajón, cerrar sesión, mostrar la contraseña) se completa solo con teclado y con foco visible en cada paso.
- **SC-005**: Texto y controles cumplen contraste AA (4,5:1 texto normal, 3:1 texto grande y bordes de controles).

**No regresión**

- **SC-006**: Las pruebas de componente existentes no se debilitan. Se modifican las que están atadas al diseño anterior (textos, estructura, formato de fechas, orden, nombre accesible del selector, etc.); **cada una está declarada con su antes → después en el mensaje del commit RED (o del merge) que la toca y resumida en la descripción del PR**. Las pruebas de las historias integradas desde `main` (#321 a #326) solo se adaptan al DOM nuevo, con la misma declaración. Lo único que pierde cobertura de texto es la ayuda de «Criterios de aceptación», que ahora es la del editor de Markdown de #326.
- **SC-007**: Se conservan los roles y nombres accesibles: lista «Proyecto», navegación «Secciones del proyecto», botón «Cerrar sesión», enlace «Nuevo proyecto», y el título de cada pantalla.
- **SC-008**: Los cambios de otras historias que tocan pantallas de la ola 2 (#321 a #326) se integran en la rama conservando su comportamiento; cada conflicto se resolvió con la regla «queda el markup del rediseño y se reinyecta lo de `main`» y está anotado en el mensaje del merge.

## Assumptions

- No hay cambios de backend ni de modelo de datos: no se toca `MODELO-DE-DATOS.md`.
- Los tokens de color anteriores quedan como alias de los nuevos hasta que la ola 2 migre sus pantallas; quitarlos es parte de la ola 2.
- La tipografía es **Inter autoalojada** (woff2 propio con licencia OFL-1.1, incluido en el repositorio); no se piden fuentes a servidores externos ni se agregan dependencias.
- El Dashboard cambia su presentación: eje con marcas 0 / mitad / techo, valores sobre las barras, indicadores de Horas y Defectos sin barras y un decimal. Los cálculos de las métricas no cambian.
- BDD: al ser frontend, los escenarios se automatizan como pruebas de componente (Principio III de la constitución). Los cambios solo de estilos no llevan ciclo RED/GREEN; los de estructura, etiquetas o comportamiento sí.
- Cada PR usa «Refs #230»; el issue queda abierto hasta la ola 2.
