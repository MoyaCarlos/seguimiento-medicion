package app

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"testing"

	apihttp "github.com/MoyaCarlos/seguimiento-medicion/internal/http"
)

// getsFueraDelProyecto son las únicas rutas GET que no leen datos de un
// proyecto: no exigen pertenencia.
var getsFueraDelProyecto = map[string]bool{
	"GET /health":      true,
	"GET /demo/estado": true,
	"GET /auth/yo":     true,
	"GET /projects":    true,
	"GET /users":       true,
}

// clienteConDatosDeEjemplo carga los datos de prueba (proyecto 25, con Sprints
// finalizados) y deja a «intruso» con cuenta pero sin integrar ningún proyecto.
func clienteConDatosDeEjemplo(t *testing.T) *clienteDePrueba {
	t.Helper()
	c := nuevoClienteDePrueba(t)
	c.esperar("", http.MethodPost, "/demo/datos", "", http.StatusCreated)
	c.registrar("intruso")
	rec := c.esperar("", http.MethodPost, "/auth/ingreso", `{"nombre":"carlos.moya","contrasena":"demo1234"}`, http.StatusOK)
	var sesion struct{ Token string }
	if err := json.Unmarshal(rec.Body.Bytes(), &sesion); err != nil || sesion.Token == "" {
		t.Fatalf("no se pudo ingresar como carlos.moya: %s", rec.Body.String())
	}
	c.tokens["carlos.moya"] = sesion.Token
	return c
}

// idsDelProyecto devuelve el proyecto de los datos de ejemplo, una de sus
// historias y uno de sus Sprints finalizados.
func idsDelProyecto(t *testing.T, c *clienteDePrueba) (proyecto, historia, sprint int64) {
	t.Helper()
	proyecto = 25
	var historias []struct{ ID int64 }
	rec := c.esperar("carlos.moya", http.MethodGet, "/projects/25/backlog", "", http.StatusOK)
	if err := json.Unmarshal(rec.Body.Bytes(), &historias); err != nil || len(historias) == 0 {
		t.Fatalf("el proyecto 25 no tiene historias en el backlog: %s", rec.Body.String())
	}
	var sprints []struct {
		ID     int64
		Estado string
	}
	rec = c.esperar("carlos.moya", http.MethodGet, "/projects/25/sprints", "", http.StatusOK)
	if err := json.Unmarshal(rec.Body.Bytes(), &sprints); err != nil {
		t.Fatal(err)
	}
	for _, s := range sprints {
		if s.Estado == "Finalizado" {
			return proyecto, historiaConVotacion(t, c, idsDe(historias)), s.ID
		}
	}
	t.Fatalf("el proyecto 25 no tiene Sprints finalizados: %s", rec.Body.String())
	return 0, 0, 0
}

func idsDe(historias []struct{ ID int64 }) []int64 {
	ids := make([]int64, 0, len(historias))
	for _, h := range historias {
		ids = append(ids, h.ID)
	}
	return ids
}

// historiaConVotacion devuelve la primera historia del Product Backlog del
// proyecto con votación iniciada, para que la lectura de un integrante dé 200.
func historiaConVotacion(t *testing.T, c *clienteDePrueba, historias []int64) int64 {
	t.Helper()
	for _, id := range historias {
		if c.pedir("carlos.moya", http.MethodGet, "/backlog/"+strconv.FormatInt(id, 10)+"/votacion", "").Code == http.StatusOK {
			return id
		}
	}
	t.Fatal("el Product Backlog de los datos de ejemplo no tiene una historia con votación iniciada")
	return 0
}

// rutaConcreta reemplaza los comodines según la familia de la ruta; una ruta GET
// nueva de una familia desconocida hace fallar el test en vez de quedar sin probar.
func rutaConcreta(t *testing.T, patron string, proyecto, historia, sprint int64) string {
	t.Helper()
	_, ruta, _ := strings.Cut(patron, " ")
	id := map[string]int64{"/projects/": proyecto, "/backlog/": historia, "/sprints/": sprint}
	for prefijo, valor := range id {
		if strings.HasPrefix(ruta, prefijo) {
			ruta = strings.Replace(ruta, "{id}", strconv.FormatInt(valor, 10), 1)
			return strings.NewReplacer("{sprintId}", strconv.FormatInt(sprint, 10)).Replace(ruta)
		}
	}
	t.Fatalf("la ruta %q no es de un proyecto, una historia ni un Sprint: clasificala en getsFueraDelProyecto o en rutaConcreta", patron)
	return ""
}

// Toda lectura con alcance de proyecto exige ser integrante (403 «No sos
// integrante de este proyecto»); quien integra el proyecto lee con 200.
func TestLecturaDeProyecto_SoloLosIntegrantesLeen(t *testing.T) {
	c := clienteConDatosDeEjemplo(t)
	proyecto, historia, sprint := idsDelProyecto(t, c)

	probadas := 0
	for _, r := range apihttp.Rutas(apihttp.Handlers{}) {
		if !strings.HasPrefix(r.Patron, "GET ") || getsFueraDelProyecto[r.Patron] {
			continue
		}
		probadas++
		t.Run(r.Patron, func(t *testing.T) {
			if r.Nivel == apihttp.Lectura {
				t.Errorf("%s lee datos de un proyecto y sigue en el nivel Lectura (sin pertenencia)", r.Patron)
			}
			ruta := rutaConcreta(t, r.Patron, proyecto, historia, sprint)
			ajeno := c.pedir("intruso", http.MethodGet, ruta, "")
			if ajeno.Code != http.StatusForbidden || strings.TrimSpace(ajeno.Body.String()) != noSosIntegrante {
				t.Errorf("intruso GET %s: se esperaba 403 %s, se obtuvo %d %.80s", ruta, noSosIntegrante, ajeno.Code, ajeno.Body.String())
			}
			if integrante := c.pedir("carlos.moya", http.MethodGet, ruta, ""); integrante.Code != http.StatusOK {
				t.Errorf("carlos.moya GET %s: se esperaba 200, se obtuvo %d %.80s", ruta, integrante.Code, integrante.Body.String())
			}
		})
	}
	if probadas < 17 {
		t.Errorf("se probaron %d rutas GET de proyecto, se esperaban al menos 17", probadas)
	}
}
