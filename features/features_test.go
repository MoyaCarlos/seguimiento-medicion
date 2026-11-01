package features

import (
	"testing"

	"github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		Name: "BDD",
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			sc := InitializeScenario(ctx)           // HU-01 (+ HU-05)
			InitializeScenarioProyecto(ctx)         // HU-04
			InitializeScenarioListado(ctx)          // HU-04 (listado)
			InitializeScenarioEstado(ctx, sc)       // HU-13
			InitializeScenarioBacklog(ctx, sc)      // HU-15
			InitializeScenarioConsultarSprints(ctx) // HU-05 (consulta)
			InitializeScenarioAsignacion(ctx, sc)   // HU-06
			InitializeScenarioVotacion(ctx, sc)     // HU-02
			InitializeScenarioHoras(ctx, sc)        // HU-07
			InitializeScenarioDatosPrueba(ctx)      // Issue #144
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"."},
			TestingT: t,
			Strict:   true,
		},
	}
	if suite.Run() != 0 {
		t.Fatal("fallaron los escenarios BDD")
	}
}
