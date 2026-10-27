package http

import (
	"os"
	"regexp"
	"slices"
	"testing"
)

func TestPrefijosAPI_IncluyenLosQueProxeaVite(t *testing.T) {
	contenido, err := os.ReadFile("../../web/vite.config.ts")
	if err != nil {
		t.Fatalf("leer vite.config.ts: %v", err)
	}

	coincidencias := regexp.MustCompile(`['"](/[^'"]*)['"]\s*:\s*apiProxy`).FindAllStringSubmatch(string(contenido), -1)
	if len(coincidencias) == 0 {
		t.Fatal("no se encontraron prefijos en vite.config.ts: ¿cambió el formato?")
	}

	for _, coincidencia := range coincidencias {
		prefijo := coincidencia[1]
		if !slices.Contains(PrefijosAPI, prefijo) {
			t.Errorf("vite.config.ts proxea %q pero no está en PrefijosAPI", prefijo)
		}
	}
}
