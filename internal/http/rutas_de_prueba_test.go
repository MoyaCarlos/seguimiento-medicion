package http

import (
	"context"
	nethttp "net/http"
	"strings"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// anaDePrueba es la persona con sesión iniciada en los tests de handlers.
var anaDePrueba = domain.User{ID: 7, Nombre: "ana", NombreNormalizado: "ana"}

// conSesionDe devuelve la petición con u como usuario de la sesión, igual que
// la deja ConSesion.
func conSesionDe(r *nethttp.Request, u domain.User) *nethttp.Request {
	return r.WithContext(context.WithValue(r.Context(), claveUsuarioDeSesion{}, u))
}

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
