# Seguimiento y Medición de Proyectos de Software

Trabajo Práctico Integrador de la asignatura **Ingeniería y Calidad de Software**
(UTN Facultad Regional San Rafael, 2026).

## Integrantes

-   Moya Carlos Esteban - Scrum Master
-   Iriarte Lopez Ana Valentina - Product Builder
-   Vulcano Candela Nair - Product Builder
-   Martinez Oldani Jimena - Product Builder

## Descripción del proyecto

**Software Metrics & Estimation** es una aplicación para administrar un proyecto
de software gestionado con Scrum, permitiendo estimar, planificar, hacer
seguimiento y medir su calidad. Funciona como un tablero de gestión ágil que
además calcula automáticamente las métricas del proyecto.

Principales funcionalidades:

-   Gestión de proyectos e integrantes.
-   Product Backlog con historias de usuario (prioridad, estado, story points,
    criterios de aceptación).
-   Gestión de Sprints (Sprint Goal, historias asignadas, cierre, historial).
-   Estimación mediante Story Points y **Planning Poker**.
-   Registro de esfuerzo (horas trabajadas) por integrante e historia/tarea.
-   Gestión de defectos (severidad, estado, sprint de detección/resolución).
-   Cálculo de métricas: velocidad del equipo, desviación esfuerzo estimado vs.
    real, porcentaje de historias completadas, defectos detectados/resueltos.
-   Dashboard con representación gráfica de las métricas.
-   Generación de reportes de proyecto/Sprint.

## Stack técnico

-   **Backend:** Go (núcleo de negocio y API REST).
-   **Frontend:** React + Vite (SPA).
-   **Persistencia:** SQLite embebido (`modernc.org/sqlite`).
-   **Testing:** Godog (BDD) + `testing` estándar de Go (TDD).
-   **Reportes:** generación de PDF con Maroto.

## Correr con Docker

Una sola imagen: el binario Go sirve la API y el frontend, y corre sin root.

```sh
docker build -t seguimiento-medicion .
docker run -d --name seg -p 8080:8080 -v seguimiento-data:/data seguimiento-medicion
```

- `GET http://localhost:8080/health` responde `ok`; la interfaz se abre en `http://localhost:8080`.
- La base SQLite vive en `/data` (volumen `seguimiento-data`): sobrevive a `docker rm` y a reiniciar el contenedor.
- Variables de entorno: `PORT` (default `8080`), `DB_PATH` (`/data/seguimiento.db` en la imagen) y `WEB_DIR` (`/app/web/dist`).
- Si el puerto 8080 está ocupado, publicá en otro: `-p 8081:8080`.
- Con un directorio del host en lugar de un volumen (`-v $(pwd)/datos:/data`), tiene que ser escribible por el usuario `65532` (`sudo chown 65532:65532 datos`).

## Despliegue en Render

La aplicación está publicada en https://seguimiento-medicion-sprintometro.onrender.com (plan gratis de Render, imagen Docker del `Dockerfile` de la raíz, región Oregon).

- **Cómo se despliega:** cada merge a `main` que toque código se publica solo, después de que pasen los checks de CI (`backend` y `frontend`). No hay pasos manuales.
- **Qué no redeploya:** los cambios que solo toquen `docs/`, `specs/`, `features/`, `*.md`, `.github/` o `.specify/`. Cualquier otro archivo (incluido `render.yaml`) sí redeploya.
- **La base se vacía en cada deploy:** no hay disco persistente, así que cada deploy o reinicio arranca con la base vacía.
- **El servicio se duerme a los 15 minutos sin tráfico:** el primer pedido después tarda decenas de segundos. Abrí la URL unos minutos antes de una demo.
- **`DEMO_HABILITADO`:** la variable de entorno (`true` en Render) controla el botón de datos de prueba (#144).
- **Regla operativa:** no mergear en la hora previa a la presentación del 16-nov, porque un deploy en curso puede dejar la app caída unos minutos y vacía la base.
- **`render.yaml`:** es solo referencia versionada de la configuración, que se administra a mano en el dashboard. No crear un Blueprint a partir de él: Render lo aplicaría al servicio existente y lo que el archivo no declare volvería a su valor por defecto.

## Especificación de funcionalidades (SDD)

El proyecto aplica **Specification-Driven Development** con
[GitHub Spec Kit](https://github.com/github/spec-kit), la herramienta oficial
de GitHub para esta práctica: antes de implementar cada historia de usuario
se escribe una especificación formal (objetivo, reglas de negocio, casos
límite, condiciones de error, criterios de aceptación), siguiendo el flujo
`/speckit.specify` → `/speckit.clarify` → `/speckit.plan` → `/speckit.tasks` →
`/speckit.implement`. Las especificaciones quedan versionadas en `specs/`, y
son la base de la trazabilidad Historia → Spec → Escenario BDD → Tests → Código
que se muestra en la defensa final.

Ya está configurado en el repo para Claude Code y OpenCode — detalles de uso
en [`docs/GUIA-IA-TERMINAL.md`](./docs/GUIA-IA-TERMINAL.md).

Metodología, arquitectura, patrones de diseño y principios aplicados están
documentados en [`AGENTS.md`](./AGENTS.md).

## Documentación del proyecto

- [`docs/PRODUCT-BACKLOG.md`](./docs/PRODUCT-BACKLOG.md) — las 15 historias de usuario.
- [`docs/GUIA-SCRUM-MASTER.md`](./docs/GUIA-SCRUM-MASTER.md) — checklist de
  ceremonias por Sprint, Definition of Done y entregables del TP.
- [`docs/GUIA-IA-TERMINAL.md`](./docs/GUIA-IA-TERMINAL.md) — cómo usar los
  agentes de IA (Claude Code / OpenCode) en este proyecto.
- [`docs/METRICA-CONSUMO-MCP.md`](./docs/METRICA-CONSUMO-MCP.md) — medición
  de consumo de tokens de las herramientas de IA usadas.
