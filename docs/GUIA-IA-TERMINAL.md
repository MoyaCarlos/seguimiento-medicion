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

### Comandos opcionales (quedan disponibles, no son obligatorios en el flujo)

- **`/speckit.constitution`** — se corre **una sola vez por proyecto** (ya
  no hace falta repetirlo), define los principios/restricciones que Spec Kit
  respeta al generar cualquier spec. Ya está configurado.
- **`/speckit.clarify`** — dentro del flujo principal (ver arriba): hace
  hasta 5 preguntas puntuales sobre lo ambiguo de la spec antes de pasar a
  `/speckit.plan`, y guarda las respuestas en la spec misma.
- **`/speckit.checklist`** — genera una checklist de calidad para validar que
  la spec de una historia esté completa/clara, después de `/speckit.plan`.
  Útil si dudás de que la spec quedó bien armada antes de seguir.
- **`/speckit.analyze`** — chequea que spec, plan y tasks de una misma
  historia sean consistentes entre sí (no se contradigan), después de
  `/speckit.tasks` y antes de `/speckit.implement`.
- **`/speckit.converge`** — compara el código ya escrito contra la spec/plan/
  tasks de una historia y agrega como tareas nuevas lo que falte. Sirve para
  retomar una historia que quedó a medias.
- **`/speckit.taskstoissues`** — sube las tareas de `/speckit.tasks` como
  Issues de GitHub (ver la sección de abajo, necesita tu propio token).

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

5. **Nunca pegues API keys en el chat ni las commitees.** En OpenCode se
   conectan con `/connect` (ver más abajo); en Claude Code (MCP de GitHub) van
   como variable de entorno referenciada con `{env:VAR}` en `.mcp.json`. En
   ningún caso el valor real de una key debería aparecer en un archivo
   versionado del repo.

## Ahorrar tokens: una tarea por sesión

Cada mensaje que mandás en una sesión reenvía **toda la conversación
anterior** de esa sesión — cuanto más larga y mezclada esté, más caro (en
tokens, en plata, en velocidad de respuesta) sale cada mensaje nuevo, incluso
si la pregunta actual no tiene nada que ver con lo de antes.

**Quedate en la misma sesión cuando:**
- Seguís trabajando en la **misma historia** o el mismo archivo (ej. escribiste
  el test RED, ahora pedís el GREEN — necesita el contexto de lo anterior).
- Estás iterando sobre algo que el agente ya entendió (corrigiendo un detalle,
  ajustando un mensaje de error, etc.).
- Vas a revisar/discutir el mismo cambio que acaban de hacer.

**Abrí una sesión nueva cuando:**
- Vas a empezar **otra historia distinta**, sin relación con la anterior.
- Cambiás de tipo de tarea (pasaste de codear a redactar la spec de otra
  funcionalidad, o a una pregunta general sin relación con lo que venías haciendo).
- La sesión ya lleva rato larga y no la necesitás más — no hace falta "cerrarla
  prolijo", simplemente arrancá una nueva para lo próximo.

**Cómo pedir sin gastar de más:**
- **Un pedido concreto por mensaje**, no una lista de 5 cosas no relacionadas
  en un mismo prompt — si una parte falla o hay que ajustarla, tenés que
  repetir/corregir todo el bloque.
- **No le pidas que relea archivos que ya leyó en la misma sesión** — si ya
  te mostró el contenido de un archivo hace 2 mensajes, no hace falta
  pedírselo de nuevo salvo que haya cambiado.
- **Cortá una respuesta que se fue por las ramas** en vez de dejar que seas vos
  quien la lea entera — pedile que resuma o vaya al punto.
- **Pedí la versión más simple primero** (esto es literal el modo "ponytail"
  que ya tienen activo en este repo): la solución mínima que funciona, no la
  más completa "por si acaso". Si después hace falta más, se agrega.

### Horarios más baratos (si usás DeepSeek)

DeepSeek tiene precios distintos según la hora — **fuera de horario pico
cuesta la mitad**. En hora argentina (ART, UTC-3), el horario pico (caro) es:

- **Lunes a viernes, 22:00 a 01:00** y **03:00 a 07:00** — evitar programar en esas franjas si se puede.
- **Todo el resto de la semana es horario barato**, esto incluye **sábado y domingo completos**, y el horario normal de trabajo/estudio (mañana, tarde, noche hasta las 22hs) entre semana.

En la práctica, el horario en el que normalmente van a estar laburando ya cae en la franja barata — esto importa sobre todo si a alguna se le ocurre dejar una tarea larga corriendo de madrugada, ahí conviene esperar a la mañana.

Esto **no aplica a Claude Code** (suscripción de Carlos, no cobra por
horario) ni a NVIDIA (es por créditos, no por franja horaria).

## Prompts que funcionan mejor

- **Sé específico con archivos y funciones**: "en `internal/domain/sprint.go`, la función `CerrarSprint`..." en vez de "arreglá el cierre de sprint".
- **Decí qué NO tocar**: si estás en medio de otra cosa, aclaralo.
- **Si el pedido es ambiguo, dejá que pregunte** en vez de que asuma y
  tengas que deshacer.
- **Pedí el "por qué" en decisiones no obvias**, no solo el código — te va a
  servir para la presentación final.

## Conectar DeepSeek en OpenCode

El equipo va a trabajar con OpenCode + DeepSeek (Carlos usa Claude Code, pero
solo hasta que se le termine la suscripción). **DeepSeek ya viene como
proveedor nativo de OpenCode** — no hace falta ningún archivo de config en el
repo, se conecta directo.

**1. Instalar OpenCode** (si no lo tenés):
```bash
curl -fsSL https://opencode.ai/install | bash
```
O instalar el CLI y después la extensión de Visual Studio Code.

La API key la pasa Carlos por mensaje privado — no la compartan en ningún
otro lado, ni la peguen en el chat con el agente.

**2. Conectarla con `/connect`** (la key queda guardada solo en tu máquina,
en `~/.local/share/opencode/auth.json`, **nunca en ningún archivo del
repo**):
1. Parada en la carpeta del repo, arrancá OpenCode: `opencode`
2. Adentro, escribí `/connect`
3. Buscá y elegí **DeepSeek** (aparece como proveedor nativo, con todo el
   catálogo de modelos actualizado).
4. Pegá la key cuando te la pida.

**3. Elegir el modelo:**
```bash
opencode -m deepseek/deepseek-flash    # default recomendado: rápido y barato
opencode -m deepseek/deepseek-v4-pro   # solo para algo puntual bien difícil (~3x más caro)
```
**Usen `deepseek-flash` por default** — de sobra para el trabajo normal, con
1M de contexto. Reserven V4 Pro para algo específico que Flash no resuelva bien.

Cambiar de modelo dentro de una sesión ya abierta: `/models`.

*(Nota: si ven `deepseek-chat` o `deepseek-reasoner` en algún lado, son
nombres viejos ya dados de baja por DeepSeek — usen `deepseek-flash` en su lugar.)*

**NVIDIA (build.nvidia.com)** también es proveedor **nativo** de OpenCode —
mismo procedimiento que DeepSeek, sin archivo de config: `/connect` → buscar
**NVIDIA** → pegar la key (se genera en build.nvidia.com, con créditos
gratis). Da acceso a modelos Nemotron y otros modelos abiertos. Cuando esté
lista la key del equipo, se usa igual que arriba, solo cambiando el nombre
del proveedor en `/connect` y en `-m nvidia/<modelo>`.

## Ponytail (opcional, cada una decide)

Si ven mencionado "ponytail": es un plugin que Carlos usa en Claude Code
(activado en `.claude/settings.json`, que está en `.gitignore` — no se
comparte por el repo). Lo que hace: en cada respuesta, empuja al agente a
dar siempre **la solución más simple que funcione** primero — reusar lo que
ya existe, no crear abstracciones que nadie pidió, no meter código de más
"por las dudas". Es básicamente forzar KISS/YAGNI de `AGENTS.md` de forma
automática en vez de tener que pedirlo cada vez.

**También existe para OpenCode** (no es exclusivo de Claude Code) — pero lo
dejamos como decisión de cada una, no lo metemos en el repo compartido para
no imponérselo a nadie. Si lo querés probar, se instala **en tu config
personal** (no toca el repo ni afecta a las demás):

```bash
mkdir -p ~/.config/opencode
```
Y en `~/.config/opencode/opencode.json` (creálo si no existe):
```json
{
  "$schema": "https://opencode.ai/config.json",
  "plugin": ["@dietrichgebert/ponytail"]
}
```

Si no lo instalás, el efecto lo conseguís igual pidiéndolo a mano: *"la
versión más simple que funcione, sin agregar nada que no pedí"* — es la
misma idea que ya está en "Ahorrar tokens" más arriba, solo que ponytail lo
hace automático en cada respuesta.

## Modelos de NVIDIA se dan de baja rápido

Si usando el proveedor NVIDIA les aparece un error tipo `410 Gone` /
`"has reached its end of life"`, no es un problema de configuración —
NVIDIA retira modelos de su catálogo (build.nvidia.com) con fecha de baja
fija, más seguido que DeepSeek con su propia API. Solución: `/models` y
elegir cualquier modelo vigente de la lista, no hay que arreglar nada.

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
