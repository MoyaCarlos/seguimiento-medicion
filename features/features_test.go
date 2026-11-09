package features

import (
	"testing"

	"github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		Name: "BDD",
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			sc := InitializeScenario(ctx)                   // HU-01 (+ HU-05)
			InitializeScenarioProyecto(ctx)                 // HU-04
			InitializeScenarioListado(ctx)                  // HU-04 (listado)
			InitializeScenarioEstado(ctx, sc)               // HU-13
			InitializeScenarioBacklog(ctx, sc)              // HU-15
			InitializeScenarioConsultarSprints(ctx)         // HU-05 (consulta)
			InitializeScenarioAsignacion(ctx, sc)           // HU-06
			votacion := InitializeScenarioVotacion(ctx, sc) // HU-02
			InitializeScenarioConsenso(ctx, votacion)       // HU-14
			InitializeScenarioHoras(ctx, sc)                // HU-07
			InitializeScenarioDatosPrueba(ctx)              // Issue #144
			InitializeScenarioEstadoHistoria(ctx, sc)       // HU-11
			InitializeScenarioDefectos(ctx, sc)             // HU-08
			InitializeScenarioMetricas(ctx)                 // HU-03
			InitializeScenarioHistorial(ctx)                // HU-12
			InitializeScenarioObtenerHistoria(ctx)          // Issue #222
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
