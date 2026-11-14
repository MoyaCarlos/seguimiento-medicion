# Implementation Plan: Cerrar y reabrir un proyecto (frontend)

**Branch**: `feature/312-cerrar-reabrir-frontend` | **Date**: 2026-11-13 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/047-cerrar-reabrir-frontend/spec.md`

## Summary

Solo frontend. Se amplían los permisos compartidos de HU-21 (`web/src/compartido/permisos/`) con el estado «cerrado» leído del mismo `GET /projects/{id}` (que ya devuelve `rol` y `fecha_cierre`). Un hook nuevo, `usePermisosProyecto(id)`, hace un solo pedido y entrega `{ rol, fechaCierre, cerrado, puedeGestionar, puedeEscribir }`; `useEsScrumMaster` pasa a ser un envoltorio fino sobre él. Se agregan el aviso «Proyecto cerrado el AAAA-MM-DD», una nota «El proyecto está cerrado» y la zona «Cerrar proyecto» (confirmación en línea) montada desde Editar proyecto. Luego cada pantalla que escribe adopta el hook en tres fases independientes. El contrato de la API es el de la spec 044 y no cambia.

## Technical Context

**Language/Version**: TypeScript, React 19, react-router 8 (el `web/` existente)

**Primary Dependencies**: las actuales. Ninguna nueva (Principio VI)

**Storage**: N/A (el estado de cierre vive en el servidor)

**Testing**: Vitest + Testing Library (`npm test` en `web/`), `fetch` simulado como en `permisos.test.tsx`

**Target Platform**: SPA en navegador

**Project Type**: web, solo el directorio `web/`

**Performance Goals**: sin pedidos nuevos: cada pantalla sigue haciendo el mismo `GET /projects/{id}` de HU-21 (FR-007a)

**Constraints**: no tocar `App.tsx`, `App.css`, `rutas.tsx`, `LayoutProyecto.tsx` ni el backend; estilos solo en el CSS de cada feature; los usos y las pruebas de HU-21 no cambian de comportamiento (SC-006; `permisos.test.tsx` solo suma `fechaCierre: null` a tres aserciones `toEqual`)

**Scale/Scope**: 3 fases; el Panel y 7 pantallas que escriben (Editar proyecto, Sprints, Historias —crear y listado—, Votación, Defectos, Equipo) y los componentes hijos que reciben las props

## Constitution Check

| Principio | Estado |
|---|---|
| I. Test-First | Cumple: cada tarea va en RED → GREEN, con commits `RED:` / `GREEN:` |
| II. SDD | Cumple: specify → clarify → plan. Las 3 fases son cortes de una misma historia (#312), con revisión entre fases; no se mezclan historias |
| III. BDD | Cumple por equivalencia (constitución 1.1.0): cada escenario de la spec es una prueba de componente con el escenario en el nombre |
| IV. Clean/Hexagonal | No aplica (frontend) |
| V. Clean Code / KISS / YAGNI | Cumple: sin store global ni Context nuevos; un hook, dos componentes chicos, props con default |
| VI. Stack | Cumple: sin dependencias nuevas |

Sin violaciones; no hay Complexity Tracking. Se re-evaluó después del diseño: sigue igual.

## Decisiones técnicas

Las cinco decisiones del pedido se adoptan tal cual, con **un ajuste** y dos precisiones (detalle en [research.md](research.md)):

- **Ajuste (aviso)**: para que, tras cerrar o reabrir, Editar se actualice «sin recargar» (FR-005, SC-005) y como el hook es la única fuente del estado, `usePermisosProyecto` devuelve además `aplicarCierre(fechaCierre | null)`. Es un agregado; la forma `{ rol, fechaCierre, cerrado, puedeGestionar, puedeEscribir }` no cambia. Sin esto, la zona y el resto de Editar quedarían con el dato viejo hasta recargar.
- **Precisión**: `puedeGestionar` y `puedeEscribir` valen `true` mientras el rol o el cierre no se conocen (cargando o error); decide el servidor.
- **Precisión**: la zona de cierre usa su propia regla, no `puedeGestionar`: solo el Scrum Master cierra y reabre, y reabrir debe poder usarse justamente con el proyecto cerrado (donde `puedeGestionar` es `false`). Usa `rol` (Product Builder → zona deshabilitada con nota).

## Project Structure

### Documentation (this feature)

```text
specs/047-cerrar-reabrir-frontend/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── ui-contracts.md
└── tasks.md             # lo genera /speckit-tasks
```

### Source Code (repository root)

```text
web/src/
├── compartido/permisos/
│   ├── useRolEnProyecto.ts          # se amplía: suma fechaCierre al estado 'listo'
│   ├── usePermisosProyecto.ts       # nuevo
│   ├── useEsScrumMaster.ts          # envoltorio fino sobre usePermisosProyecto
│   ├── AvisoProyectoCerrado.tsx     # nuevo
│   ├── NotaCerrado.tsx              # nuevo (o NotaGestion con otro texto, ver research)
│   └── permisos.test.tsx            # solo suma fechaCierre: null a 3 toEqual; pruebas nuevas en archivos nuevos
├── features/proyecto/
│   ├── api.ts                       # + cerrarProyecto, reabrirProyecto, fecha_cierre en Proyecto
│   ├── ZonaCierreProyecto.tsx       # nuevo
│   ├── Editar.tsx                   # cambio mínimo: monta la zona, aviso y deshabilita Guardar
│   ├── Panel.tsx                    # cambio mínimo: muestra AvisoProyectoCerrado
│   └── estilos.css                  # estilos de la zona y el aviso
├── features/sprints/                # Fase 2: index, FormularioInicio, BotonCerrar, FormularioHoras, EstadoHistoria
├── features/backlog/                # Fase 3: crear, Edicion, AsignarASprint
├── features/votacion/               # Fase 3: index y AccionesScrumMaster
├── features/defectos/               # Fase 3: PantallaDefectos y FormularioDefecto
└── features/proyecto/Configuracion.tsx  # Fase 3 («Equipo»: asignar integrante)
```

**Structure Decision**: se mantiene la estructura por feature de `web/src`. Lo compartido va en `compartido/permisos/`; la zona de cierre, en `features/proyecto/`. Cada pantalla de las Fases 2 y 3 llama una vez a `usePermisosProyecto` en su raíz y pasa `puedeGestionar` y `puedeEscribir` a sus hijos como props con valor por defecto `true`, así los usos y las pruebas actuales no se rompen.

## Fases (cada una con su RED/GREEN y su validación)

1. **Fase 1**: `useRolEnProyecto` (+`fechaCierre`), `usePermisosProyecto`, `useEsScrumMaster` como envoltorio, `AvisoProyectoCerrado`, `NotaCerrado`, `cerrarProyecto`/`reabrirProyecto`, `ZonaCierreProyecto`, Editar proyecto y el aviso en el Panel. Entrega US1, US3 y US4 completas, y US2 en Editar y Panel. Las pantallas de solo lectura (Estado, Dashboard, Historial) quedan fuera: pueden mostrar el aviso cuando se rediseñen.
2. **Fase 2**: Sprints (crear, iniciar, cerrar Sprint, horas y estado de las historias).
3. **Fase 3**: Historias (crear, editar, asignar a Sprint), Votación, Defectos y Equipo.

Cada fase termina con `npm test`, `npm run lint` y `npm run build` en verde y la validación manual de [quickstart.md](quickstart.md). Se puede cortar después de cualquier fase.
