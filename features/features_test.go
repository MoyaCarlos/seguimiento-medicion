package features

import (
	"testing"

	"github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		Name: "BDD",
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			InitializeScenario(ctx)         // HU-01 (+ HU-05)
			InitializeScenarioProyecto(ctx) // HU-04
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"."},
			TestingT: t,
		},
	}
	if suite.Run() != 0 {
		t.Fatal("fallaron los escenarios BDD")
	}
}
