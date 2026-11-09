# Research: Generación de Reporte en PDF (HU-10)

## R1. Rutas y descarga

- **Decision**: `GET /projects/{id}/reporte` para el reporte del proyecto y
  `GET /sprints/{id}/reporte` para el de un Sprint finalizado. Responder `200` con
  `Content-Type: application/pdf` y `Content-Disposition: attachment` con un nombre ASCII
  (`reporte-proyecto-<id>.pdf` y `reporte-proyecto-<id>-sprint-<numero>.pdf`).
- **Rationale**: respeta el estilo REST en español del resto de la API (`/metricas`,
  `/defectos`) y evita una ruta `/reports` nueva. El nombre identifica el proyecto y, en el
  reporte de Sprint, también el Sprint (FR-024).
- **Alternatives considered**: `/projects/{id}/reports/pdf` (más largo, sin precedente);
  devolver base64 en JSON (obliga al frontend a decodificar y no es una descarga directa).

## R2. Maroto y UTF-8

- **Decision**: Maroto v2.4.4, generación secuencial en memoria con `document.GetBytes()` y
  fuente incorporada Arial (cp1252). No se registran fuentes externas. Compresión desactivada
  para mantener el contenido inspeccionable.
- **Rationale**: es la versión que exige Go 1.26.1 (el proyecto usa 1.26.5); cp1252 cubre
  tildes, eñes y signos de apertura del español. Embeber un `.ttf` agrandaría el PDF y el
  repositorio sin necesidad para el alcance. Con compresión apagada, el texto generado puede
  verificarse en un smoke test sin descomprimir.
- **Alternatives considered**: Maroto v1 (sin soporte ni API actual); fuente Noto embebida
  (solo necesaria para scripts fuera de cp1252, fuera de alcance).

## R3. Reutilización de datos y DRY

- **Decision**: el servicio `GenerarReporte` compone el agregado con los puertos existentes:
  `ProjectRepository` (proyecto e integrantes), `SprintRepository` (Sprints),
  `ProductBacklogRepository` (historias sin Sprint), `HistoriasDeSprintRepository`
  (historias del Sprint), `DefectoRepository` (defectos) y `DatosMetricasRepository`.
  Las métricas se calculan con `domain.CalcularMetricas`, la misma función de HU-03.
- **Rationale**: evita duplicar reglas y no crea una consulta agregada nueva para un volumen
  académico. El reporte es solo lectura.
- **Alternatives considered**: una consulta SQL única con todas las secciones (más rápida,
  pero duplicaría el cálculo de métricas y las reglas del dominio).

## R4. Historias planificadas y arrastradas

- **Decision**: un Sprint finalizado muestra como planificadas sus historias completadas
  (`historias.sprint_id = sprint`) y las arrastradas al cerrarlo
  (`historias_arrastradas`), con los Story Points congelados en esa tabla. Las arrastradas se
  leen con un puerto nuevo `HistoriasArrastradasRepository` implementado en
  `SQLiteBacklogRepository`. Una historia arrastrada y completada después aparece completada
  en el reporte del Sprint donde se completó, no en el anterior.
- **Rationale**: `CerrarConArrastre` devuelve las no completadas al Product Backlog y guarda
  su foto en `historias_arrastradas`; sin esa tabla el reporte perdería lo planificado.
- **Alternatives considered**: recalcular desde el estado actual (incorrecto: cambia el
  histórico, prohibido por FR-014).

## R5. Secciones y "sin datos"

- **Decision**: el reporte de proyecto siempre incluye datos del proyecto, integrantes,
  Product Backlog y, si existe, Sprint Activo. Los Sprints finalizados se listan en orden con
  su detalle. Los totales de métricas, esfuerzo y defectos se toman solo de Sprints
  finalizados. Si no hay ninguno, se escriben como "sin datos" y el PDF se genera igual.
- **Rationale**: reproduce la decisión del Scrum Master del 2026-11-01 y FR-015/FR-017.
- **Alternatives considered**: incluir el Sprint Activo en los totales (lo prohíbe la spec).

## R6. Errores y archivo parcial

- **Decision**: el service valida existencia y estado antes de generar. El generador produce
  `[]byte` en memoria; el handler recién escribe cabeceras y cuerpo cuando la generación
  terminó bien. Un Sprint no finalizado responde `400` con
  `ErrSprintNoFinalizadoParaReporte`; un Sprint o proyecto inexistente, `404`; un fallo de generación,
  `500` con "no se pudo generar el reporte".
- **Rationale**: cumple US4 y FR-022 sin dejar archivos corruptos. Mantiene el orden de
  precedencia de la spec (existencia antes que estado).
- **Alternatives considered**: escribir directo al `ResponseWriter` (dejaría archivos
  parciales si Maroto falla a mitad de la escritura).

## R7. Excepción en `ServirSPA`

- **Decision**: `esPedidoDePagina` deja de tratar como página los `GET` cuyo camino termina en
  `/reporte`, para que la descarga directa desde el navegador (`Accept: text/html`) llegue a la
  API. Se agrega un test en `spa_test.go`.
- **Rationale**: hoy todo GET con `Accept: text/html` devuelve `index.html`; una navegación
  directa a la ruta de descarga rompería.
- **Alternatives considered**: exigir `fetch` desde el frontend (el frontend lo hará, pero la
  URL debe funcionar igual).

## R8. Nombre y contenido del archivo

- **Decision**: el nombre del PDF es ASCII y estable por identificadores:
  `reporte-proyecto-<id>.pdf` y `reporte-proyecto-<id>-sprint-<numero>.pdf`. La fecha de
  generación se imprime en el documento, no se usa en el nombre.
- **Rationale**: evita problemas de codificación en `Content-Disposition` y cumple FR-024.
- **Alternatives considered**: incluir el nombre del proyecto en el archivo (requiere sanear
  tildes y espacios; sin valor para el alcance).

## R9. Tamaño y paginación

- **Decision**: Maroto pagina automáticamente cuando una fila no entra; las tablas se arman
  fila por fila para no truncar. No se usan listas con header repetido global.
- **Rationale**: cumple SC-006 con el volumen del proyecto y mantiene KISS.
- **Alternatives considered**: header de tabla global repetido en cada página (impediría el
  título y los datos del proyecto en la primera página).
