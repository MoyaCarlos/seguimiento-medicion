package http

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const rutaSalud = "/health"

// PrefijosAPI son las rutas que atiende la API. Todas menos /health comparten
// ruta con pantallas del SPA y deben incluir los prefijos que proxea web/vite.config.ts.
var PrefijosAPI = []string{"/backlog", "/projects", "/sprints", rutaSalud}

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
			http.ServeFile(w, r, index)
		default:
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
	return r.URL.Path != rutaSalud
}
