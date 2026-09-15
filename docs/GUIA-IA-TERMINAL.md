# Guía rápida: programar con IA desde la terminal

Para todo el equipo, sin importar si usan Claude Code, OpenCode+DeepSeek u
OpenCode+NVIDIA. El contexto del proyecto (stack, arquitectura, principios,
convenciones) ya está en [`AGENTS.md`](../AGENTS.md) — se carga solo al abrir
el agente en esta carpeta, no hace falta pegarlo en el prompt.

## Spec Kit: qué es y para qué lo usamos

[GitHub Spec Kit](https://github.com/github/spec-kit) es la herramienta
oficial de GitHub para **Specification-Driven Development (SDD)** — la
práctica que el TP exige usar. Antes de escribir código de una funcionalidad,
se escribe primero una especificación formal (objetivo, reglas de negocio,
casos límite, condiciones de error, criterios de aceptación); Spec Kit guía
ese proceso con comandos en vez de dejarlo librado a que cada uno escriba la
spec como quiera.

**Flujo por cada historia del Product Backlog:**

1. `/speckit.specify` — convierte la historia en una especificación formal.
2. `/speckit.clarify` — el agente pregunta lo que esté ambiguo antes de seguir.
3. `/speckit.plan` — arma el plan técnico usando el stack ya definido en `AGENTS.md`.
4. `/speckit.tasks` — desglosa el plan en tareas concretas.
5. `/speckit.implement` — recién ahí se escribe código, **una historia a la
   vez**, con revisión humana del diff (nunca todo el backlog de una).

Las specs generadas quedan versionadas en `specs/<nombre-historia>/` dentro
del repo — eso es lo que van a mostrar como evidencia de SDD en la defensa
final, y de ahí sale la trazabilidad Historia → Spec → BDD → Tests → Código.

**Ya está instalado y commiteado en el repo** — con hacer `git pull` en
`seguimiento-medicion/` ya tenés los comandos `/speckit.*` disponibles al
abrir el agente ahí adentro. No hace falta que lo reinstales.

- Claude Code: `/speckit-specify`, `/speckit-plan`, `/speckit-tasks`, etc. (con guión).
- OpenCode: `/speckit.specify`, `/speckit.plan`, `/speckit.tasks`, etc. (con punto).

Si alguna vez hiciera falta reinstalarlo desde cero en otra máquina: parate
**adentro de la carpeta del repo** (`cd ruta/a/seguimiento-medicion`) — el
flag `--here` inicializa en la carpeta donde estés parado, así que si lo
corrés desde otro lado te crea los archivos en el lugar equivocado — y ahí sí:
```bash
uv tool install specify-cli
specify init --here --integration claude    # o: opencode
```

### `/speckit-taskstoissues` necesita tu propio GitHub token

Este skill (sube las tareas de `/speckit.tasks` como Issues) usa el GitHub
MCP server, ya configurado en `.mcp.json` (compartido en el repo, sin ningún
token adentro). Para que te funcione en tu máquina:

1. Creá un **Personal Access Token fine-grained** en
   https://github.com/settings/personal-access-tokens/new — Resource owner:
   tu usuario, Repository access: solo `seguimiento-medicion`, Permissions →
   **Issues: Read and write**. Expiración: alcanza con fin de cuatrimestre
   (ej. 31/12/2026), no pongas "No expiration".
2. Guardalo como variable de entorno, **nunca lo pegues en el chat con el
   agente ni lo commitees**:
   ```bash
   echo 'export GITHUB_PAT=tu_token_aca' >> ~/.zshrc
   ```
   (o `~/.bashrc` si usás bash en vez de zsh).
3. Abrí una **terminal nueva** (los procesos ya corriendo no ven variables
   agregadas después) y arrancá tu agente ahí — recién ahí el MCP conecta.
4. Verificar sin exponer el token: `claude mcp list` tiene que mostrar
   `github ... ✔ Connected` (puede tardar unos segundos en la primera conexión).

## Qué instalar en tu agente (y qué NO)

**No hace falta instalar nada vos misma.** Todo lo que el equipo decidió usar
ya está commiteado en el repo — con `git pull` en `seguimiento-medicion/` te
llega automáticamente, tanto en Claude Code como en OpenCode:

- **Spec Kit** (`/speckit.*`) — flujo de SDD, ver sección de arriba.
- **`test-driven-development`** (skill de
  [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills)) —
  refuerza el ciclo RED/GREEN/REFACTOR con Prove-It Pattern para bugs, test
  pyramid y Arrange-Act-Assert. Recomendada por el profesor. Se activa sola,
  sin comando — el agente la usa cuando estás implementando o corrigiendo algo.

**No instales paquetes de skills adicionales por tu cuenta** (ej. el resto del
catálogo de `agent-skills`, o `mattpocock/skills`) sin avisar al equipo antes.
No es por desconfianza: cada skill que se suma **consume contexto/tokens** aunque
no se use — hay una medición real de esto en
[`docs/METRICA-CONSUMO-MCP.md`](METRICA-CONSUMO-MCP.md) — y varias de esas
skills hacen lo mismo que Spec Kit (specs, planning), lo que genera confusión
sobre cuál usar para qué. Si el profesor recomienda algo puntual, lo evaluamos
como equipo antes de instalarlo (así decidimos con la skill de TDD: se
descartó el resto de ambos catálogos por redundancia).

## Reglas de oro

1. **Todo lo que genere la IA lo tenés que entender antes de commitear.**
   No es opcional: la consigna del TP dice explícitamente que el equipo es
   responsable de todo el código, generado o no con IA, y que hay que poder
   justificarlo en la defensa final. Si no entendés una parte del diff,
   preguntale al agente "explicame por qué hiciste esto" antes de aceptar.

2. **Nunca le pidas que resuelva un Sprint entero de una sola pasada.**
   Una historia (o un ítem del roadmap) a la vez, revisás el diff, recién
   ahí seguís con la próxima. Si el agente arrasa todo de una, perdés la
   evidencia real de TDD y de las ceremonias de Scrum que se evalúan.

3. **Para TDD, pedí el ciclo explícito, no "escribime código y tests".**
   Prompt tipo: *"Escribí el test que falla para `CalcularVelocidad`, no
   toques la implementación todavía"* → confirmás que falla → *"ahora el
   código mínimo para que pase"* → confirmás verde → *"¿hay algo para
   refactorizar acá?"*. Cada paso es su propio commit (ver convención en
   `AGENTS.md`).

4. **Corré los tests vos mismo antes de commitear.** No confíes en que el
   agente diga "ya están en verde" — `go test ./...` (o el comando de Godog)
   con tus propios ojos.

5. **Nunca pegues API keys en el chat ni las commitees.** Van en variables
   de entorno (`.env`, con `.gitignore`), referenciadas como `{env:VAR}` en
   `opencode.json`.

## Prompts que funcionan mejor

- **Sé específico con archivos y funciones**: "en `internal/domain/sprint.go`,
  la función `CerrarSprint`..." en vez de "arreglá el cierre de sprint".
- **Decí qué NO tocar**: si estás en medio de otra cosa, aclaralo.
- **Si el pedido es ambiguo, dejá que pregunte** en vez de que asuma y
  tengas que deshacer.
- **Pedí el "por qué" en decisiones no obvias**, no solo el código — te va a
  servir para la presentación final.

## Multi-modelo en OpenCode

Quienes no tienen Claude Code usan OpenCode con dos proveedores ya
configurados en `opencode.json`:

- **DeepSeek** → para codear (`opencode -m deepseek/deepseek-chat`).
- **NVIDIA (build.nvidia.com)** → para planificar/discutir (`opencode -m
  nvidia/<modelo>`).

Se puede tener una ventana de cada uno abierta al mismo tiempo, cada sesión
usa su propio modelo sin pisarse. Cambiar de modelo dentro de una sesión ya
abierta: `/models`.

## Skills nativas de Claude Code (no hace falta instalarlas)

Estas ya vienen incluidas en Claude Code — no son algo que se instale con
`npx skills`, simplemente están disponibles. OpenCode no tiene este sistema
de skills nativas, así que esta lista es solo para quien use Claude Code
(distinto de `test-driven-development`, que sí instalamos para ambas
herramientas — ver arriba). Se invocan con `/nombre-skill`:

- **`/code-review`** — revisa el diff/PR pendiente buscando bugs antes de
  mergear. Usar antes de abrir un PR importante.
- **`/security-review`** — revisión de seguridad del diff pendiente, útil en
  todo lo que toque la API REST o la persistencia.
- **`simplify`** — pasada de limpieza (reuso, simplificación, eficiencia)
  sobre lo que cambiaste.
- **`dataviz`** — guía de diseño para los gráficos del Dashboard (paleta de
  colores, forma de los charts, legibilidad).
- **`artifact-diagramming`** — para armar diagramas (arquitectura, flujo de
  trazabilidad Historia→SDD→BDD→Tests→Código) para la documentación o la
  defensa final.
- **`run`** — levanta el proyecto (backend Go + frontend React) para probar
  una feature en vivo en vez de confiar solo en los tests.
- **`init`** — genera documentación automática del código, útil más
  adelante cuando ya haya bastante escrito.
- **`fewer-permission-prompts`** — reduce las confirmaciones repetidas de
  comandos (git, go test, npm) si te resultan molestas.

## Antes de pedir ayuda al equipo

Si el agente se traba o da vueltas en círculos con lo mismo, no insistas 5
veces con el mismo prompt — cortá, escribí en el grupo qué intentaste, y
seguime con otra tarea mientras alguien lo destraba.
