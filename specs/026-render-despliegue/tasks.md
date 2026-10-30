---

description: "Tareas del despliegue en Render (Issue #139): infraestructura y documentación, sin ciclo RED/GREEN"

---

# Tasks: Despliegue continuo en Render

**Input**: `/specs/026-render-despliegue/` (spec.md, plan.md, research.md, data-model.md, contracts/render-yaml.md, quickstart.md)

**Tests**: no hay tests de código. El chequeo es la validación de `render.yaml` contra el esquema oficial y la verificación manual del quickstart.

**Reglas**: cada fase es su propio commit (los COMMIT son tareas); archivos agregados por nombre (nunca `git add -A`); sin `git push` ni PR; no se toca el dashboard de Render ni código, `Dockerfile` o dependencias. No se commitean `docs/reviews/`, `INFORME-*.md`, `Proyecto - 2026.pdf`, `docs/prototipo-sprint2/` ni `specs/007-frontend-base/analisis.md`.

## Phase 1: User Story 2 — `render.yaml` como referencia versionada (P1)

**Goal**: archivo en la raíz que reproduce la configuración actual, validado, con la advertencia del Blueprint.

**Independent Test**: validación contra `https://render.com/schema/render.yaml.json` sin errores y comparación valor por valor con la tabla de la spec.

- [ ] T001 [US2] Crear render.yaml en la raíz con el contenido de contracts/render-yaml.md, más un comentario inicial que diga que es solo referencia y que no se crea un Blueprint porque Render lo aplicaría al servicio existente, sincronizaría con cada push y lo que falte volvería a su valor por defecto (FR-001, FR-003, FR-004, FR-005)
- [ ] T002 [US2] Validar render.yaml contra el esquema oficial (descargado a una carpeta temporal fuera del repo, JSON Schema 2020-12); si alguna clave falla o no se puede validar, quitarla del archivo y anotarla (FR-002). Confirmar 0 errores
- [ ] T003 [US2] Comparar render.yaml con la tabla de configuración de spec.md y la columna "Valor" de research.md: cada valor coincide, sin secretos (SC-001)
- [ ] T004 [US2] COMMIT `Agregar render.yaml como referencia versionada del servicio en Render` (solo render.yaml)

---

## Phase 2: User Story 1 y 3 — README "Despliegue en Render" (P1)

**Goal**: una sola sección que responde cómo se despliega y qué esperar en una demo.

**Independent Test**: revisión del README contra FR-006 a FR-012.

- [ ] T005 [US1] [US3] Agregar la sección "Despliegue en Render" en README.md, después de "Correr con Docker": URL, despliegue automático tras los checks, rutas que no redeployan y aclaración de que cualquier otro archivo (incluido `render.yaml`) sí redeploya, base vacía en cada deploy o reinicio, sueño a los 15 minutos con la recomendación de abrir la URL antes de una demo, `DEMO_HABILITADO` y el botón de datos de prueba (#144), no mergear en la hora previa a la presentación del 16-nov, y que `render.yaml` es solo referencia (FR-006..FR-012)
- [ ] T006 [US1] [US3] Revisar que la sección responde las 7 preguntas de SC-002 y que no contradice a "Correr con Docker"
- [ ] T007 [US1] [US3] COMMIT `Documentar el despliegue en Render en el README` (solo README.md)

---

## Phase 3: User Story 4 — verificación manual contra el servicio real (P2)

**Goal**: comprobar lo documentado sin modificar el servicio.

**Independent Test**: los tres pedidos de quickstart.md.

- [ ] T008 [US4] Despertar el servicio y ejecutar los pedidos `GET` de quickstart.md (`/health`, `/` con `Accept: text/html`, `/projects`); solo lectura, sin tocar el dashboard (FR-014)
- [ ] T009 [US4] Anotar en la sección "Resultado" de specs/026-render-despliegue/quickstart.md la fecha, el código HTTP y lo observado de cada pedido, y el resultado de la validación del esquema (FR-013, SC-003). Lo que no se pueda comprobar (por ejemplo, abrir la raíz en un navegador real) se anota como no verificado
- [ ] T010 [US4] COMMIT `Registrar la verificación manual del servicio en Render` (solo quickstart.md)

---

## Phase 4: Polish

- [ ] T011 Confirmar con `git diff --stat origin/main` que el diff solo toca render.yaml, README.md y specs/026-render-despliegue/ (SC-004, FR-015)
- [ ] T012 Marcar las tareas hechas en este archivo y commitear `Marcar las tareas de la spec 026 como hechas` (solo tasks.md)
