# Research: Dashboard de Métricas (HU-09)

Las decisiones R1 y R2 son las de Clarifications, confirmadas en la reunión de definición.

## R1. Sin endpoint nuevo

- **Decision**: la pantalla consume `GET /projects/{id}/metricas` tal como está.
- **Rationale**: la respuesta ya trae cada Sprint con su número, planificados y completados, y los valores del proyecto. Un endpoint "listo para graficar" duplicaría esa respuesta (Principio V). El contrato de HU-03 ya prevé este consumo. Godog ya cubre la consulta en `features/calculo_metricas.feature`.
- **Consecuencias en las sub-issues**:
  - #71 (Backend): se cierra sin código, con una nota que enlaza esta decisión.
  - #72 (Tests TDD): el TDD pasa a las funciones puras de `metricas.ts`; los casos de handler ya están en HU-03.
  - #70 (BDD): los escenarios se automatizan como pruebas de componente; no hace falta `dashboard_metricas.feature`.
- **Alternatives considered**: endpoint `/dashboard` propio, descartado por duplicar la respuesta sin aportar datos nuevos.

## R2. Recharts para el gráfico de velocidad

- **Decision**: `recharts` 3.x, importada solo desde la feature `dashboard`. La ruta diferida (`lazy`) que ya usan todas las features deja la librería fuera de la carga inicial.
- **Rationale**: resuelve ejes, leyenda y adaptación al ancho, y es lo que supuso la estimación de 8 SP. La versión 3.10.1 declara compatibilidad con React 19 y pide `react-is` como peer.
- **Alcance de la librería**: solo el gráfico de velocidad. Los indicadores son HTML y CSS (R5), sin depender de ella.
- **Alternatives considered**: SVG propio: sin dependencias, pero hay que escribir y mantener ejes, escalas y adaptación al ancho. Queda como alternativa si Recharts no se puede integrar; solo cambiaría `GraficoVelocidad`.

## R3. Pruebas del gráfico

- **Decision**: las pruebas verifican los valores en la tabla de valores por Sprint (R4), que son el mismo dato que recibe el gráfico, y que el gráfico se monte sin fallar. No se verifican coordenadas ni detalles internos del SVG de la librería.
- **Rationale**: jsdom no calcula dimensiones, y probar el SVG de Recharts probaría la librería.
- **Verificado con Recharts 3.10.1** (T004): jsdom no trae `ResizeObserver`, y aun así `ResponsiveContainer` no falla. Sin dimensiones no dibuja nada; con `initialDimension` dibuja el gráfico y la leyenda, sin simular nada. El componente usa `initialDimension`.
- **Barras en cero**: Recharts no dibuja rectángulos de altura cero, así que US1.4 se verifica en la tabla, no contando barras.

## R4. Series distinguibles y valores como texto (FR-006)

- **Decision**: "Comprometidos" con barra clara y borde; "Completados" con barra sólida, más leyenda. Debajo, una tabla plegada "Ver valores por Sprint", con el mismo patrón que "Ver detalle" de HU-07. El ancho mínimo del gráfico crece con la cantidad de Sprints y, si no entra, se desplaza dentro de su propio contenedor, no la página (FR-018).
- **Rationale**: relleno y borde distinguen las series sin depender del color; la tabla da los valores sin puntero y a lectores de pantalla, y plegada no carga el panel con datos crudos.
- **Colores**: se definen con la skill `dataviz` al implementar, a partir de los tokens de `App.css` (`--lilac`, `--complement`).
- **Alternatives considered**: etiquetas de valor sobre cada barra, descartadas porque con muchos Sprints se superponen en 320 px.

## R5. Indicadores en HTML y CSS

- **Decision**: tarjetas con valor y una barra hecha con CSS: proporción para el porcentaje, dos barras en la misma escala para horas y defectos. Velocidad y desviación llevan valor y texto aclaratorio.
- **Rationale**: es la solución más simple, se prueba en jsdom y no agrega código de la librería (KISS).
- **Detalle**: el ancho de cada barra sale de `proporcion(valor, maximo)`, que devuelve 0 si el máximo es 0, para evitar proporciones indefinidas. Es escala de dibujo, no un cálculo de métricas (FR-017).

## R6. Formato y errores

- **Decision**: números con `toLocaleString('es-AR', { maximumFractionDigits: 2 })`, como en HU-07 y HU-02. `null` se muestra como "sin datos". Los errores muestran el `message` de `ApiError` o de `ConnectionError`.
- **Rationale**: `aislamiento.test.ts` impide importar `clasificarError` de Sprints; además, sin campos de formulario alcanza con `error.message`. El 404 ya trae "proyecto no encontrado".
- **Sin reintento**: ninguna pantalla actual ofrece un botón de reintento, así que no se agrega.

## R7. Enlaces

- **Decision**: el estado vacío enlaza a `/projects/{id}/sprints`, y el aviso de horas por Story Point a `/projects/{id}/proyecto/edit`. Son enlaces de ruta, sin importar otras features.
