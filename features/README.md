# Escenarios BDD (Godog)

Acá van los archivos `.feature` (Gherkin, Given-When-Then) por historia, más
los step definitions en Go que los automatizan. Se agregan cuando arranca la
implementación de cada historia (vía `/speckit.plan` → `/speckit.tasks`), no
antes.

## Runner único

`features_test.go` registra un único runner Godog (`Name: "BDD"`) con
`Paths: ["."]`, que ejecuta **todos** los `.feature` del directorio. Cada
historia expone su propia función de inicialización de pasos (p.ej.
`InitializeScenario` para HU-01, `InitializeScenarioProyecto` para HU-04), y el
runner las invoca dentro de `ScenarioInitializer`.

## Cómo agregar una historia nueva

1. Crear `features/<historia>.feature` con tag `@HU-XX` y los escenarios
   Given-When-Then derivados 1:1 de los Criterios de Aceptación del spec.
2. Crear `features/steps_<historia>.go` con un `*Context` y una función
   `InitializeScenario<Historia>(ctx *godog.ScenarioContext)` que registre sus
   pasos con `ctx.Step(...)`, más los `ctx.Before`/`ctx.After` para inicializar
   y cerrar la base `:memory:`.
3. Invocar esa función dentro del `ScenarioInitializer` de `features_test.go`,
   al lado de las historias existentes.
4. Revisar que los regex de los pasos **no colisionen** con los de otras
   historias (Godog falla si hay dos step definitions con el mismo regex).
   Si colisionan, desambiguar el texto del paso de la historia nueva (y su
   `.feature`) sin salirse del spec.
5. Verificar con `go test ./features/ -v` que todos los escenarios (propios y
   ajenos) siguen en verde.
