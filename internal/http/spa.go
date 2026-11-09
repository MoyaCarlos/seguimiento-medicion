package http

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const rutaSalud = "/health"

// PrefijosAPI documenta las rutas que atiende la API. No decide el enrutado:
// ServirSPA responde index.html a todo GET con text/html salvo /health. Un test
// verifica que incluya los prefijos que proxea web/vite.config.ts. Las descargas
// que terminan en /reporte (PDF de HU-10) también llegan a la API aunque acepten
// text/html; cualquier otra excepción se agrega junto a /health.
var PrefijosAPI = []string{"/backlog", "/defectos", "/projects", "/sprints", "/demo", "/auth", rutaSalud}

func ServirSPA(api http.Handler, dirWeb string) http.Handler {
	if info, err := os.Stat(dirWeb); err != nil || !info.IsDir() {
		return api
	}
	archivos := http.Dir(dirWeb)
	servidorArchivos := http.FileServer(archivos)
	index := filepath.Join(dirWeb, "index.html")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case esLecturaDeArchivo(r, archivos):
			servidorArchivos.ServeHTTP(w, r)
		case esPedidoDePagina(r):
			w.Header().Add("Vary", "Accept")
			w.Header().Set("Cache-Control", "no-cache")
			http.ServeFile(w, r, index)
		default:
			w.Header().Add("Vary", "Accept")
			api.ServeHTTP(w, r)
		}
	})
}

func esLecturaDeArchivo(r *http.Request, archivos http.Dir) bool {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return false
	}
	archivo, err := archivos.Open(r.URL.Path)
	if err != nil {
		return false
	}
	defer archivo.Close()
	info, err := archivo.Stat()
	return err == nil && !info.IsDir()
}

func esPedidoDePagina(r *http.Request) bool {
	if r.Method != http.MethodGet || !strings.Contains(r.Header.Get("Accept"), "text/html") {
		return false
	}
	if r.URL.Path == rutaSalud || strings.HasSuffix(r.URL.Path, "/reporte") {
		return false
	}
	return true
}
