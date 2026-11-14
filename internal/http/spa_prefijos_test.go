package http

import (
	"os"
	"regexp"
	"slices"
	"strings"
	"testing"
)

func prefijosQueProxeaVite(t *testing.T) []string {
	t.Helper()
	contenido, err := os.ReadFile("../../web/vite.config.ts")
	if err != nil {
		t.Fatalf("leer vite.config.ts: %v", err)
	}

	coincidencias := regexp.MustCompile(`['"](/[^'"]*)['"]\s*:\s*apiProxy`).FindAllStringSubmatch(string(contenido), -1)
	if len(coincidencias) == 0 {
		t.Fatal("no se encontraron prefijos en vite.config.ts: ¿cambió el formato?")
	}

	prefijos := make([]string, 0, len(coincidencias))
	for _, coincidencia := range coincidencias {
		prefijos = append(prefijos, coincidencia[1])
	}
	return prefijos
}

func TestPrefijosAPI_IncluyenLosQueProxeaVite(t *testing.T) {
	for _, prefijo := range prefijosQueProxeaVite(t) {
		if !slices.Contains(PrefijosAPI, prefijo) {
			t.Errorf("vite.config.ts proxea %q pero no está en PrefijosAPI", prefijo)
		}
	}
}

func TestPrefijosAPI_LosProxeaVite(t *testing.T) {
	proxeados := prefijosQueProxeaVite(t)
	for _, prefijo := range PrefijosAPI {
		if prefijo == rutaSalud {
			continue
		}
		if !slices.Contains(proxeados, prefijo) {
			t.Errorf("%q está en PrefijosAPI pero vite.config.ts no lo proxea", prefijo)
		}
	}
}

func TestPrefijosAPI_CubrenLaTablaDeRutas(t *testing.T) {
	for _, ruta := range Rutas(Handlers{}) {
		camino := ruta.Patron
		if _, despues, tieneMetodo := strings.Cut(camino, " "); tieneMetodo {
			camino = despues
		}
		segmento, _, _ := strings.Cut(strings.TrimPrefix(camino, "/"), "/")
		prefijo := "/" + segmento
		if !slices.Contains(PrefijosAPI, prefijo) {
			t.Errorf("la ruta %q usa el prefijo %q, que no está en PrefijosAPI", ruta.Patron, prefijo)
		}
	}
}
