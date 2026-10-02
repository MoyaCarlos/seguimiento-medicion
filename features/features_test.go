package features

import (
	"testing"

	"github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		Name:                "HU-01",
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"crear_historia_backlog.feature"},
			TestingT: t,
		},
	}
	if suite.Run() != 0 {
		t.Fatal("fallaron los escenarios BDD de HU-01")
	}
}
