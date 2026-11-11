package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
	"time"

	apihttp "github.com/MoyaCarlos/seguimiento-medicion/internal/http"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// nivelesDeLaSpec es la tabla «Qué se protege» de la spec 042 (HU-21), escrita a
// mano: si una ruta nueva no se clasifica, o se clasifica distinto, falla el test.
var nivelesDeLaSpec = map[string]apihttp.NivelRuta{
	"GET /health":         apihttp.Publica,
	"POST /auth/registro": apihttp.Publica,
	"POST /auth/ingreso":  apihttp.Publica,
	"GET /demo/estado":    apihttp.Publica,
	"POST /demo/datos":    apihttp.Publica,
	"/demo/":              apihttp.Publica,

	"GET /auth/yo":     apihttp.Lectura,
	"POST /auth/salir": apihttp.SoloSesion,
	"POST /projects":   apihttp.SoloSesion,
	"GET /users":       apihttp.SoloSesion,

	"POST /backlog":                     apihttp.SesionIntegrante,
	"PUT /backlog/{id}":                 apihttp.SesionIntegrante,
	"PUT /backlog/{id}/estado":          apihttp.SesionIntegrante,
	"PUT /backlog/{id}/sprint":          apihttp.SesionIntegrante,
	"POST /backlog/{id}/horas":          apihttp.SesionIntegrante,
	"POST /projects/{id}/defectos":      apihttp.SesionIntegrante,
	"POST /defectos/{id}/resolver":      apihttp.SesionIntegrante,
	"POST /backlog/{id}/votacion/votos": apihttp.SesionIntegrante,

	"POST /sprints":                          apihttp.SesionScrumMaster,
	"POST /sprints/{id}/iniciar":             apihttp.SesionScrumMaster,
	"POST /sprints/{id}/cerrar":              apihttp.SesionScrumMaster,
	"POST /backlog/{id}/votacion":            apihttp.SesionScrumMaster,
	"POST /backlog/{id}/votacion/revelar":    apihttp.SesionScrumMaster,
	"POST /backlog/{id}/votacion/rondas":     apihttp.SesionScrumMaster,
	"POST /backlog/{id}/votacion/estimacion": apihttp.SesionScrumMaster,
	"PUT /projects/{id}":                     apihttp.SesionScrumMaster,
	"POST /projects/{id}/members":            apihttp.SesionScrumMaster,

	"GET /projects":                           apihttp.Lectura,
	"GET /projects/{id}":                      apihttp.Lectura,
	"GET /projects/{id}/status":               apihttp.Lectura,
	"GET /projects/{id}/members":              apihttp.Lectura,
	"GET /projects/{id}/backlog":              apihttp.Lectura,
	"GET /projects/{id}/sprints":              apihttp.Lectura,
	"GET /projects/{id}/metricas":             apihttp.Lectura,
	"GET /projects/{id}/historial":            apihttp.Lectura,
	"GET /projects/{id}/historial/{sprintId}": apihttp.Lectura,
	"GET /projects/{id}/avance":               apihttp.Lectura,
	"GET /projects/{id}/defectos":             apihttp.Lectura,
	"GET /projects/{id}/reporte":              apihttp.Lectura,
	"GET /backlog/{id}":                       apihttp.Lectura,
	"GET /backlog/{id}/horas":                 apihttp.Lectura,
	"GET /backlog/{id}/votacion":              apihttp.Lectura,
	"GET /backlog/{id}/votacion/rondas":       apihttp.Lectura,
	"GET /sprints/{id}/backlog":               apihttp.Lectura,
	"GET /sprints/{id}/reporte":               apihttp.Lectura,
}

func TestRutas_CadaRutaTieneElNivelDeLaSpec(t *testing.T) {
	declaradas := map[string]apihttp.NivelRuta{}
	for _, r := range apihttp.Rutas(apihttp.Handlers{}) {
		declaradas[r.Patron] = r.Nivel
	}
	for patron, nivel := range nivelesDeLaSpec {
		got, ok := declaradas[patron]
		if !ok {
			t.Errorf("la spec protege %q pero la ruta no está registrada", patron)
			continue
		}
		if got != nivel {
			t.Errorf("%q tiene nivel %d, la spec pide %d", patron, got, nivel)
		}
	}
	var sinClasificar []string
	for patron := range declaradas {
		if _, ok := nivelesDeLaSpec[patron]; !ok {
			sinClasificar = append(sinClasificar, patron)
		}
	}
	sort.Strings(sinClasificar)
	if len(sinClasificar) > 0 {
		t.Errorf("rutas sin clasificar en la tabla de la spec: %v", sinClasificar)
	}
}

func apiDePrueba(t *testing.T) *http.ServeMux {
	t.Helper()
	db, err := repository.AbrirSQLite(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	db.SetMaxOpenConns(1)
	if err := repository.AplicarMigraciones(context.Background(), db); err != nil {
		t.Fatal(err)
	}
	return NuevoMux(db, Opciones{
		DemoHabilitado: func() bool { return true },
		Ahora:          time.Now,
		Contrasenas:    service.NewServicioContrasenas(1),
	})
}

// pedidoDe arma método y ruta concretos a partir del patrón, con 1 en cada comodín.
func pedidoDe(patron string) (string, string) {
	metodo, ruta, ok := strings.Cut(patron, " ")
	if !ok {
		return http.MethodGet, strings.TrimSuffix(patron, "/") + "/otra"
	}
	ruta = strings.NewReplacer("{id}", "1", "{sprintId}", "1").Replace(ruta)
	return metodo, ruta
}

// respuestaSinSesion es el cuerpo del 401 de sesión (HU-20); un 401 con otro
// cuerpo, como el de credenciales incorrectas en /auth/ingreso, no es de sesión.
const respuestaSinSesion = "{\"mensaje\":\"Tu sesión no es válida\"}\n"

func TestRutas_LasProtegidasResponden401SinSesionYLasPublicasNo(t *testing.T) {
	mux := apiDePrueba(t)
	for patron, nivel := range nivelesDeLaSpec {
		t.Run(patron, func(t *testing.T) {
			metodo, ruta := pedidoDe(patron)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, httptest.NewRequest(metodo, ruta, strings.NewReader("{}")))

			rechazoDeSesion := rec.Code == http.StatusUnauthorized && rec.Body.String() == respuestaSinSesion
			if nivel == apihttp.Publica && rechazoDeSesion {
				t.Errorf("%s %s es pública y la rechazó la capa de sesión", metodo, ruta)
			}
			if nivel != apihttp.Publica && !rechazoDeSesion {
				t.Errorf("%s %s sin sesión respondió %d %q, se esperaba 401 %q", metodo, ruta, rec.Code, rec.Body.String(), respuestaSinSesion)
			}
		})
	}
}
