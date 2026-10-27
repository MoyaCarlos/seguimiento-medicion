package http

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

const (
	cuerpoIndex = "INDEX"
	cuerpoAPI   = "API"
	cuerpoAsset = "console.log('app')"
)

func apiFalsa() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(cuerpoAPI))
	})
}

func dirWebDePrueba(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "assets"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte(cuerpoIndex), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "assets", "app.js"), []byte(cuerpoAsset), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func pedir(handler http.Handler, metodo, ruta, accept string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(metodo, ruta, nil)
	if accept != "" {
		req.Header.Set("Accept", accept)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

const aceptaNavegador = "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"

func TestServirSPA_ReglaDeEnrutado(t *testing.T) {
	handler := ServirSPA(apiFalsa(), dirWebDePrueba(t))

	casos := []struct {
		nombre string
		metodo string
		ruta   string
		accept string
		cuerpo string
	}{
		{"pantalla compartida con la API pide HTML", http.MethodGet, "/projects/1/status", "text/html", cuerpoIndex},
		{"navegador con varios tipos pide HTML", http.MethodGet, "/projects/1/status", aceptaNavegador, cuerpoIndex},
		{"lista de proyectos pide HTML", http.MethodGet, "/projects", "text/html", cuerpoIndex},
		{"backlog pide HTML", http.MethodGet, "/backlog", "text/html", cuerpoIndex},
		{"sprints pide HTML", http.MethodGet, "/sprints", "text/html", cuerpoIndex},
		{"raíz pide HTML", http.MethodGet, "/", "text/html", cuerpoIndex},
		{"ruta desconocida pide HTML", http.MethodGet, "/otra-pantalla", "text/html", cuerpoIndex},
		{"misma ruta sin HTML va a la API", http.MethodGet, "/projects/1/status", "application/json", cuerpoAPI},
		{"misma ruta sin Accept va a la API", http.MethodGet, "/projects/1/status", "", cuerpoAPI},
		{"asset se sirve normal", http.MethodGet, "/assets/app.js", "", cuerpoAsset},
		{"asset con Accept de navegador se sirve normal", http.MethodGet, "/assets/app.js", "*/*", cuerpoAsset},
		{"health con HTML sigue en la API", http.MethodGet, "/health", "text/html", cuerpoAPI},
		{"POST va a la API aunque acepte HTML", http.MethodPost, "/projects", "text/html", cuerpoAPI},
		{"PUT va a la API aunque acepte HTML", http.MethodPut, "/projects/1", "text/html", cuerpoAPI},
		{"prefijo parecido no es de la API", http.MethodGet, "/projectsx", "text/html", cuerpoIndex},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			rec := pedir(handler, caso.metodo, caso.ruta, caso.accept)

			if rec.Code != http.StatusOK {
				t.Fatalf("se esperaba 200, se obtuvo %d", rec.Code)
			}
			if rec.Body.String() != caso.cuerpo {
				t.Errorf("se esperaba %q, se obtuvo %q", caso.cuerpo, rec.Body.String())
			}
		})
	}
}

func TestServirSPA_AssetInexistenteSinHTMLNoEsLaAPI(t *testing.T) {
	handler := ServirSPA(apiFalsa(), dirWebDePrueba(t))

	rec := pedir(handler, http.MethodGet, "/assets/no-existe.js", "")

	if rec.Code == http.StatusOK && rec.Body.String() == cuerpoIndex {
		t.Fatal("un asset inexistente no debe responder el index.html")
	}
}

func TestServirSPA_SinDirectorioTodoVaALaAPI(t *testing.T) {
	handler := ServirSPA(apiFalsa(), filepath.Join(t.TempDir(), "no-existe"))

	for _, ruta := range []string{"/", "/projects/1/status", "/assets/app.js"} {
		rec := pedir(handler, http.MethodGet, ruta, "text/html")

		if rec.Body.String() != cuerpoAPI {
			t.Errorf("%s: se esperaba la API, se obtuvo %q", ruta, rec.Body.String())
		}
	}
}
