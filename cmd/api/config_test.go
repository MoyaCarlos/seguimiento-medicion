package main

import "testing"

func entornoDe(variables map[string]string) func(string) string {
	return func(clave string) string { return variables[clave] }
}

func TestCargarConfig_SinVariablesUsaLosDefaults(t *testing.T) {
	cfg := cargarConfig(entornoDe(nil))

	if cfg.Puerto != "8080" {
		t.Errorf("puerto: se esperaba 8080, se obtuvo %q", cfg.Puerto)
	}
	if cfg.RutaBase != "seguimiento.db" {
		t.Errorf("ruta de la base: se esperaba seguimiento.db, se obtuvo %q", cfg.RutaBase)
	}
	if cfg.DirWeb != "web/dist" {
		t.Errorf("directorio web: se esperaba web/dist, se obtuvo %q", cfg.DirWeb)
	}
}

func TestCargarConfig_ConVariablesUsaEsosValores(t *testing.T) {
	cfg := cargarConfig(entornoDe(map[string]string{
		"PORT":    "9090",
		"DB_PATH": "/data/seguimiento.db",
		"WEB_DIR": "/app/web/dist",
	}))

	if cfg.Puerto != "9090" || cfg.RutaBase != "/data/seguimiento.db" || cfg.DirWeb != "/app/web/dist" {
		t.Errorf("no se usaron los valores del entorno: %+v", cfg)
	}
}

func TestCargarConfig_VariableVaciaUsaElDefault(t *testing.T) {
	cfg := cargarConfig(entornoDe(map[string]string{"PORT": "", "DB_PATH": "", "WEB_DIR": ""}))

	if cfg.Puerto != "8080" || cfg.RutaBase != "seguimiento.db" || cfg.DirWeb != "web/dist" {
		t.Errorf("las variables vacías deben caer al default: %+v", cfg)
	}
}

func TestCargarConfig_DemoHabilitadoSoloConTrue(t *testing.T) {
	casos := map[string]bool{"": false, "false": false, "TRUE": false, "1": false, "true": true}
	for valor, esperado := range casos {
		cfg := cargarConfig(entornoDe(map[string]string{"DEMO_HABILITADO": valor}))
		if cfg.DemoHabilitado != esperado {
			t.Errorf("DEMO_HABILITADO=%q: se esperaba %v, se obtuvo %v", valor, esperado, cfg.DemoHabilitado)
		}
	}
}
