package features

import (
	"testing"

	"github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		Name: "BDD",
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			_ = InitializeScenarioSesion(ctx)                 // HU-21 (quién actúa)
			sc := InitializeScenario(ctx)                     // HU-01 (+ HU-05)
			proyecto := InitializeScenarioProyecto(ctx)       // HU-04
			InitializeScenarioScrumMasterUnico(ctx, proyecto) // #236
			InitializeScenarioListado(ctx)                    // HU-04 (listado)
			InitializeScenarioEstado(ctx, sc)                 // HU-13
			InitializeScenarioBacklog(ctx, sc)                // HU-15
			InitializeScenarioConsultarSprints(ctx)           // HU-05 (consulta)
			InitializeScenarioAsignacion(ctx, sc)             // HU-06
			votacion := InitializeScenarioVotacion(ctx, sc)   // HU-02
			InitializeScenarioConsenso(ctx, votacion)         // HU-14
			InitializeScenarioHoras(ctx, sc)                  // HU-07
			InitializeScenarioDatosPrueba(ctx)                // Issue #144
			InitializeScenarioEstadoHistoria(ctx, sc)         // HU-11
			InitializeScenarioDefectos(ctx, sc)               // HU-08
			InitializeScenarioMetricas(ctx)                   // HU-03
			InitializeScenarioReporte(ctx)                    // HU-10
			InitializeScenarioAutenticacion(ctx)              // HU-20
			InitializeScenarioHistorial(ctx)                  // HU-12
			InitializeScenarioAvance(ctx, sc)                 // HU-18
			InitializeScenarioObtenerHistoria(ctx)            // Issue #222
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
