package http

import (
	nethttp "net/http"
	"strings"
)

// registrarRutasDe registra, sin envoltorio de sesión, las rutas de la tabla
// cuyo camino empieza con prefijo. Sirve a los tests de un handler.
func registrarRutasDe(mux *nethttp.ServeMux, h Handlers, prefijo string) {
	for _, r := range Rutas(h) {
		_, camino, ok := strings.Cut(r.Patron, " ")
		if !ok {
			camino = r.Patron
		}
		if strings.HasPrefix(camino, prefijo) {
			mux.HandleFunc(r.Patron, r.Handler)
		}
	}
}
