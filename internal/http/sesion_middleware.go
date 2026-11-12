package http

import (
	"context"
	nethttp "net/http"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

type claveUsuarioDeSesion struct{}

// ConSesion devuelve un envoltorio que exige una sesión vigente (HU-21): sin
// ella responde 401 sin ejecutar el handler; con ella guarda el usuario en el
// contexto de la petición, de donde lo lee usuarioDeSesion.
func ConSesion(obtener *service.ObtenerSesion) func(nethttp.HandlerFunc) nethttp.HandlerFunc {
	return func(siguiente nethttp.HandlerFunc) nethttp.HandlerFunc {
		return func(w nethttp.ResponseWriter, r *nethttp.Request) {
			token, ok := tokenBearer(r)
			if !ok {
				escribirError(w, domain.ErrNoAutenticado)
				return
			}
			sesion, err := obtener.Ejecutar(r.Context(), token)
			if err != nil {
				escribirError(w, err)
				return
			}
			ctx := context.WithValue(r.Context(), claveUsuarioDeSesion{}, sesion.Usuario)
			siguiente(w, r.WithContext(ctx))
		}
	}
}

// actorDeSesion devuelve el id de quien actúa; sin usuario de sesión responde
// 401 y devuelve false. Protege al handler aunque la ruta no esté envuelta.
func actorDeSesion(w nethttp.ResponseWriter, r *nethttp.Request) (int64, bool) {
	id := usuarioDeSesion(r).ID
	if id <= 0 {
		escribirError(w, domain.ErrNoAutenticado)
		return 0, false
	}
	return id, true
}

// usuarioDeSesion devuelve el usuario que dejó ConSesion en el contexto; el
// valor cero si la ruta no está envuelta.
func usuarioDeSesion(r *nethttp.Request) domain.User {
	usuario, _ := r.Context().Value(claveUsuarioDeSesion{}).(domain.User)
	return usuario
}
