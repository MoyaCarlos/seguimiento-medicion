# Guía del Scrum Master

Referencia práctica para vos (Carlos, Agile Enabler) para seguir manejando el
proyecto sin depender de tener un agente de IA disponible todo el tiempo.
Todo lo que menciona acá está versionado en el repo — si dudás de algo,
la fuente completa está en los links de la sección [Dónde está todo](#dónde-está-todo).

## Dónde estamos (al 28/09/2026)

- **Sprint 0 (Preparación): cerrado.** Repo, tablero, Spec Kit, esqueleto de
  carpetas, todo commiteado.
- **Product Backlog: 13 historias** en `docs/PRODUCT-BACKLOG.md`, cargadas
  como Issues #1-#13 en el tablero.
- **Sprint 1 (MVP): en definición** — reunión pendiente/en curso para repartir
  entre los 4: HU-04, HU-01, HU-13, HU-05 (propuesta, a confirmar en la reunión).
- **Pendiente sin resolver todavía**: la fecha de entrega final del
  cuatrimestre. Sin eso, el calendario de Sprints 2-4 sigue siendo tentativo.

## Tu rol como Scrum Master / Agile Enabler

No sos quien construye el producto (eso son las Product Builders) — sos quien:
- Organiza las ceremonias (Planning, Daily, Review, Retrospective).
- Cuida que el proceso (Scrum + SDD + BDD + TDD) se siga de verdad, no solo
  de nombre.
- Mantiene el tablero honesto (que refleje lo que realmente está pasando).
- Es el punto de contacto con el profesor/Cliente.

## Checklist por Sprint (repetir en cada uno: 1, 2, 3, 4)

### Al abrir el Sprint (Planning)

- [ ] Confirmar qué historias entran (campo **Iteration** en el tablero → el
      nombre del Sprint correspondiente).
- [ ] Cada historia tiene seteados: **Priority**, **Estimate**, **Tipo**
      (`Historia de Usuario`), **Status** en `Backlog` o `Ready`.
- [ ] Definir el **Sprint Goal** (una frase: qué se puede demostrar al cerrar
      este Sprint) y dejarlo escrito — puede ir como comentario en el Sprint
      Review posterior o en un doc de actas.
- [ ] Repartir una historia por persona (o coordinar si dos trabajan juntas
      en una historia grande).

### Durante el Sprint

- [ ] **Daily corta** (no tiene que ser reunión formal, puede ser un mensaje
      de grupo): qué hice, qué voy a hacer, si hay algún bloqueo.
- [ ] Por cada historia, el flujo es: `/speckit.specify` → `/speckit.clarify`
      → `/speckit.plan` → `/speckit.tasks` → `/speckit.implement` (ver
      `docs/GUIA-IA-TERMINAL.md` si alguien olvida los comandos exactos).
- [ ] TDD real: commits `RED:` / `GREEN:` / `REFACTOR:` separados (ver
      convención completa en `AGENTS.md`).
- [ ] Todo cambio entra por **Pull Request**, con **mínimo 1 aprobación** de
      otra persona, y se mergea con **"Create a merge commit"** (nunca squash
      — si alguien squashea sin querer, se pierde la evidencia de TDD).
- [ ] Mover el **Status** de cada historia en el tablero a medida que avanza:
      `Backlog` → `Ready` → `In progress` → `In review` → `Done`.

### Al cerrar el Sprint (Review + Retrospectiva)

- [ ] **Sprint Review**: mostrar lo que quedó funcionando (demo real, no solo
      código). Comparar contra el Sprint Goal que se definió al abrir.
- [ ] Anotar métricas simples: Story Points planificados vs. completados —
      esto va a hacer falta para justificar "Gestión del proyecto" al final.
- [ ] **Retrospectiva**: qué salió bien, qué no, una acción concreta para el
      próximo Sprint. **Dejar un acta escrita** (la consigna la pide como
      entregable explícito — no alcanza con haberla hablado).
- [ ] Cerrar el campo Iteration de las historias completadas; las que no se
      terminaron, vuelven al Product Backlog sin iteración asignada (mismo
      criterio que define HU-05 para el sistema real).

## Definition of Done (cuándo una historia está realmente terminada)

Copiado de `AGENTS.md` para tenerlo a mano sin tener que buscarlo:

- Código implementado y mergeado a `main`.
- Tests unitarios (TDD) y escenario BDD correspondiente, ambos en verde.
- Especificación SDD de la historia escrita en `specs/`.
- Revisada y aprobada por al menos otro integrante (PR).
- Sin defectos abiertos bloqueantes para esa historia.

Si una historia no cumple esto, no está "Done" — puede estar "In review",
pero no cerrada.

## Entregables del TP (no perder de vista ninguno)

La consigna pide estos 13 al final. Marcá los que ya tengan avance real:

- [ ] Repositorio Git con historial de contribuciones.
- [ ] Tablero Scrum (GitHub Projects).
- [ ] Product Backlog y Sprint Backlogs (`docs/PRODUCT-BACKLOG.md` + tablero).
- [ ] Especificaciones SDD (`specs/`, generadas con Spec Kit por historia).
- [ ] Escenarios BDD (`.feature` en `/features`).
- [ ] Pruebas automatizadas.
- [ ] Código fuente en Go.
- [ ] Evidencias de aplicación de TDD (historial de commits RED/GREEN/REFACTOR).
- [ ] Software funcional.
- [ ] Informe de métricas y cobertura de pruebas.
- [ ] **Actas de retrospectivas** (fácil de olvidar — no es solo la reunión, hay que escribirla).
- [ ] Documentación técnica y manual breve de usuario.
- [ ] Presentación y demostración final.

## Cómo se reparte la nota (para priorizar esfuerzo bajo presión de tiempo)

Si en algún momento hay que elegir en qué invertir el tiempo que queda:

| Criterio | % |
|---|---|
| Producto funcional | 25% |
| SDD, BDD y TDD | 25% |
| Calidad del software | 20% |
| Gestión del proyecto | 20% |
| Trabajo en equipo y presentación | 10% |

Es decir: **el proceso (SDD+BDD+TDD) pesa lo mismo que el producto
funcionando.** Si alguna vez hay que elegir entre "meter una feature más
rápido saltando specs/tests" o "hacerlo bien pero más lento", la segunda
opción es la que conviene según cómo puntúa la cátedra.

## Cuando se termine tu Claude Code

El resto del equipo sigue con OpenCode + DeepSeek (nativo, sin instalar
nada — con `git pull` ya tienen todo). Vos podés seguir exactamente igual
que ellas: instalar OpenCode, `/connect` con tu propia key de DeepSeek
(pasos completos en `docs/GUIA-IA-TERMINAL.md`). El contexto del proyecto
(`AGENTS.md`) se carga solo, no perdés nada de lo que ya armamos.

## Trazabilidad (para la defensa final)

Van a tener que mostrar el recorrido completo de **al menos una historia**:
Historia de Usuario → Especificación SDD → Criterios de Aceptación →
Escenario BDD → Tests → Código Go. Conviene elegir de antemano cuál van a
usar de ejemplo (una que esté bien completa, no la más simple) y tenerla
mapeada: número de HU → carpeta en `specs/` → archivo `.feature` → archivo
de test → archivo de implementación.

## Dónde está todo

- [`AGENTS.md`](../AGENTS.md) — metodología completa, arquitectura,
  principios, stack, convenciones de Git/TDD.
- [`docs/PRODUCT-BACKLOG.md`](PRODUCT-BACKLOG.md) — las 13 historias con
  todos sus campos.
- [`docs/GUIA-IA-TERMINAL.md`](GUIA-IA-TERMINAL.md) — cómo usar Claude Code /
  OpenCode, Spec Kit, buenas prácticas, ahorro de tokens.
- [`docs/METRICA-CONSUMO-MCP.md`](METRICA-CONSUMO-MCP.md) — para mostrarle
  al profesor el costo real de las herramientas de IA usadas.
- `specs/` — se va llenando historia por historia a medida que avanza la
  implementación (todavía vacío al 28/09/2026).
- Tablero: Project "Software Metrics & Estimation" en GitHub, del repo
  `MoyaCarlos/seguimiento-medicion`.

## Si algo no sale como estaba planeado

- **Si el profesor cambia un requisito o fecha**: actualizá `AGENTS.md` (o
  el archivo que corresponda) el mismo día, con un commit que lo explique —
  así queda evidencia de que el equipo reaccionó a tiempo, no que se enteró
  tarde.
- **Si una historia no llega a terminarse en su Sprint**: no se fuerza — se
  documenta en la retro por qué, y vuelve al backlog para el próximo Sprint.
  Es más valioso mostrar un proceso honesto que un backlog "completo" a costa
  de saltear SDD/BDD/TDD.
- **Si alguien no tiene claro cómo seguir un paso técnico puntual** (cómo se
  usa tal comando, cómo se configura tal cosa): primero buscarlo en los docs
  de arriba, y si no está, preguntarle directo al agente (DeepSeek/OpenCode)
  — tiene el mismo contexto del proyecto que tuvo Claude.
