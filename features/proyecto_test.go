package features

import (
	"testing"

	"github.com/cucumber/godog"
)

func TestProyectoFeatures(t *testing.T) {
	suite := godog.TestSuite{
		Name:                "HU-04",
		ScenarioInitializer: InitializeScenarioProyecto,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"creacion_proyecto_equipo.feature"},
			TestingT: t,
		},
	}
	if suite.Run() != 0 {
		t.Fatal("fallaron los escenarios BDD de HU-04")
	}
}
