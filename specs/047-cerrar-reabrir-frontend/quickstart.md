# Quickstart: validar «Cerrar y reabrir un proyecto» (frontend)

## Prerrequisitos
- Backend de la spec 044 en main, corriendo (`go run ./cmd/...` según `AGENTS.md`), y `npm install` hecho en `web/`.
- Dos usuarios de un mismo proyecto: uno Scrum Master y otro Product Builder (los datos de prueba v2 los crean).

## Automatizado (en cada fase)
```bash
cd web
npm test        # Vitest: debe incluir permisos.test.tsx sin modificar (SC-006)
npm run lint
npm run build
```

## Fase 1: zona de cierre, aviso y Editar proyecto
1. Como Scrum Master, sin Sprint Activo: Editar proyecto → «Cerrar proyecto» → «Cancelar». No cambia nada (US1.2).
2. «Cerrar proyecto» → «Confirmar». Aparece «Proyecto cerrado el AAAA-MM-DD», estado «Finalizado», «Guardar cambios» deshabilitado con nota, sin recargar (US1.1, FR-005).
3. «Reabrir proyecto»: el aviso desaparece y Guardar se habilita (US3.1, SC-005).
4. Con un Sprint Activo, confirmar el cierre: se ve «Cerrá el Sprint activo antes de cerrar el proyecto» y la confirmación se cierra (US1.3).
5. Con el foco: al abrir la confirmación el foco está en «Cancelar»; Escape cancela; doble clic en «Confirmar» envía un solo pedido.
6. Como Product Builder: la zona se ve deshabilitada con la nota (US4.1).

## Fase 2: Sprints
Con el proyecto cerrado, en Sprints quedan deshabilitados crear, iniciar y cerrar Sprint, cargar horas y cambiar el estado de las historias, con la nota «El proyecto está cerrado» y el aviso arriba. El listado sigue visible.

## Fase 3: Historias, Votación, Defectos y Equipo
Con el proyecto cerrado, quedan deshabilitados crear y editar Historia, asignar a Sprint, las acciones de Votación, registrar Defectos y asignar integrantes. Tras reabrir y volver a entrar, todo se habilita.

## Cobertura contra la spec
Cada escenario de US1–US4 y cada borde (conexión caída, doble clic, Escape, cargando, nota única, cambio de proyecto) tiene una prueba con su nombre en el título (SC-007).
