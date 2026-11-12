package app

import (
	"net/http"
	"slices"
	"strings"
	"testing"
	"time"

	apihttp "github.com/MoyaCarlos/seguimiento-medicion/internal/http"
)

// Spec 044 (#300): un proyecto cerrado es de solo lectura en toda la API.

const proyectoEstaCerrado = `{"mensaje":"El proyecto está cerrado"}`

// escrituraDeProyecto es una ruta de escritura con un pedido concreto y válido
// sobre el proyecto 1 que arma proyectoCerradoDeAna.
type escrituraDeProyecto struct {
	patron, ruta, cuerpo string
}

// escriturasDeProyecto tiene cada ruta de nivel «Sesión e integrante» o
// «Sesión y Scrum Master» de HU-21, salvo cerrar y reabrir el proyecto.
var escriturasDeProyecto = []escrituraDeProyecto{
	{"POST /backlog", "/backlog", `{"proyecto_id":1,"titulo":"Otra","descripcion":"d","prioridad":"M"}`},
	{"PUT /backlog/{id}", "/backlog/1", `{"titulo":"Cambiada","descripcion":"d","criterios_aceptacion":null,"estimacion_sp":null}`},
	{"PUT /backlog/{id}/estado", "/backlog/1/estado", `{"estado":"En progreso"}`},
	{"PUT /backlog/{id}/sprint", "/backlog/1/sprint", `{"sprint_id":2}`},
	{"POST /backlog/{id}/horas", "/backlog/1/horas", `{"fecha":"2026-11-11","actividad":"x","horas":1}`},
	{"POST /projects/{id}/defectos", "/projects/1/defectos", `{"historia_id":1,"descripcion":"d","severidad":"Alta"}`},
	{"POST /defectos/{id}/resolver", "/defectos/1/resolver", ``},
	{"POST /backlog/{id}/votacion/votos", "/backlog/2/votacion/votos", `{"valor":5}`},
	{"POST /sprints", "/sprints", `{"proyecto_id":1}`},
	{"POST /sprints/{id}/iniciar", "/sprints/2/iniciar", `{"sprint_goal":"g","fecha_inicio":"2026-11-12","fecha_fin":"2026-11-20"}`},
	{"POST /sprints/{id}/cerrar", "/sprints/1/cerrar", ``},
	{"POST /backlog/{id}/votacion", "/backlog/1/votacion", ``},
	{"POST /backlog/{id}/votacion/revelar", "/backlog/2/votacion/revelar", ``},
	{"POST /backlog/{id}/votacion/rondas", "/backlog/2/votacion/rondas", ``},
	{"POST /backlog/{id}/votacion/estimacion", "/backlog/2/votacion/estimacion", `{"estimacion_sp":5}`},
	{"PUT /projects/{id}", "/projects/1", `{"nombre":"Otro","descripcion":"","fecha_inicio":null,"fecha_fin":null}`},
	{"POST /projects/{id}/members", "/projects/1/members", `{"nombre":"diego"}`},
}

// rutasDelCierre no son escrituras sobre los datos del proyecto: lo cierran o reabren.
var rutasDelCierre = []string{"POST /projects/{id}/cerrar", "POST /projects/{id}/reabrir"}

// proyectoConDatosDeAna arma el proyecto 1 de ana con datos en cada tabla:
// historias 1 y 2, Sprint 1 finalizado (con la historia 1 y un defecto),
// Sprint 2 pendiente y una votación abierta en la historia 2. diego tiene
// cuenta y no integra el proyecto.
func proyectoConDatosDeAna(t *testing.T) *clienteDePrueba {
	t.Helper()
	c := nuevoClienteDePrueba(t)
	c.registrar("ana")
	c.registrar("diego")
	c.esperar("ana", http.MethodPost, "/projects", `{"nombre":"Tienda"}`, http.StatusCreated)
	c.esperar("ana", http.MethodPost, "/backlog", `{"proyecto_id":1,"titulo":"Carrito","descripcion":"d","prioridad":"M"}`, http.StatusCreated)
	c.esperar("ana", http.MethodPost, "/backlog", `{"proyecto_id":1,"titulo":"Pago","descripcion":"d","prioridad":"M"}`, http.StatusCreated)
	c.esperar("ana", http.MethodPost, "/sprints", `{"proyecto_id":1}`, http.StatusCreated)
	c.esperar("ana", http.MethodPost, "/sprints/1/iniciar", `{"sprint_goal":"g","fecha_inicio":"2026-11-01","fecha_fin":"2026-11-10"}`, http.StatusOK)
	c.esperar("ana", http.MethodPut, "/backlog/1/sprint", `{"sprint_id":1}`, http.StatusOK)
	c.esperar("ana", http.MethodPost, "/projects/1/defectos", `{"historia_id":1,"descripcion":"d","severidad":"Alta"}`, http.StatusCreated)
	c.esperar("ana", http.MethodPost, "/sprints/1/cerrar", ``, http.StatusOK)
	c.esperar("ana", http.MethodPost, "/sprints", `{"proyecto_id":1}`, http.StatusCreated)
	c.esperar("ana", http.MethodPost, "/backlog/2/votacion", ``, http.StatusCreated)
	return c
}

// fotoDelProyecto junta las lecturas que muestran todos los datos del proyecto 1.
func fotoDelProyecto(c *clienteDePrueba) string {
	var foto strings.Builder
	for _, ruta := range []string{"/projects/1", "/projects/1/backlog", "/projects/1/sprints", "/projects/1/defectos", "/projects/1/members", "/backlog/2/votacion"} {
		foto.WriteString(ruta + " " + c.pedir("ana", http.MethodGet, ruta, "").Body.String())
	}
	return foto.String()
}

func TestProyectoCerrado_TodaEscrituraResponde409SinCambios(t *testing.T) {
	c := proyectoConDatosDeAna(t)
	c.esperar("ana", http.MethodPost, "/projects/1/cerrar", ``, http.StatusOK)
	antes := fotoDelProyecto(c)

	for _, e := range escriturasDeProyecto {
		t.Run(e.patron, func(t *testing.T) {
			metodo, _, _ := strings.Cut(e.patron, " ")
			rec := c.pedir("ana", metodo, e.ruta, e.cuerpo)
			if rec.Code != http.StatusConflict || strings.TrimSpace(rec.Body.String()) != proyectoEstaCerrado {
				t.Errorf("se esperaba 409 %s, se obtuvo %d %s", proyectoEstaCerrado, rec.Code, rec.Body.String())
			}
		})
	}

	if despues := fotoDelProyecto(c); despues != antes {
		t.Errorf("las escrituras rechazadas cambiaron datos:\nantes:   %s\ndespués: %s", antes, despues)
	}
}

// Si se agrega una ruta de escritura, tiene que sumarse a escriturasDeProyecto.
func TestProyectoCerrado_LaTablaCubreTodasLasEscrituras(t *testing.T) {
	var deLaAPI, deLaTabla []string
	for _, r := range apihttp.Rutas(apihttp.Handlers{}) {
		if (r.Nivel == apihttp.SesionIntegrante || r.Nivel == apihttp.SesionScrumMaster) && !slices.Contains(rutasDelCierre, r.Patron) {
			deLaAPI = append(deLaAPI, r.Patron)
		}
	}
	for _, e := range escriturasDeProyecto {
		deLaTabla = append(deLaTabla, e.patron)
	}
	slices.Sort(deLaAPI)
	slices.Sort(deLaTabla)
	if !slices.Equal(deLaAPI, deLaTabla) {
		t.Errorf("rutas de escritura de la API:\n %v\nrutas probadas con el proyecto cerrado:\n %v", deLaAPI, deLaTabla)
	}
}

// Las lecturas responden igual con el proyecto abierto y cerrado, salvo el
// estado y la fecha de cierre.
func TestProyectoCerrado_LasLecturasSiguen(t *testing.T) {
	c := proyectoConDatosDeAna(t)
	codigos := func() map[string]int {
		resultado := map[string]int{}
		for patron, nivel := range nivelesDeLaSpec {
			if nivel != apihttp.Lectura {
				continue
			}
			metodo, ruta := pedidoDe(patron)
			resultado[patron] = c.pedir("ana", metodo, ruta, "").Code
		}
		return resultado
	}
	abierto := codigos()

	c.esperar("ana", http.MethodPost, "/projects/1/cerrar", ``, http.StatusOK)

	for patron, codigo := range codigos() {
		if codigo != abierto[patron] {
			t.Errorf("%s: con el proyecto abierto respondía %d y cerrado responde %d", patron, abierto[patron], codigo)
		}
	}
	if rec := c.esperar("ana", http.MethodGet, "/projects/1/status", "", http.StatusOK); strings.TrimSpace(rec.Body.String()) != `{"estado":"Finalizado"}` {
		t.Errorf("estado = %s, se esperaba Finalizado", rec.Body.String())
	}
	if rec := c.esperar("ana", http.MethodGet, "/projects/1/backlog", "", http.StatusOK); !strings.Contains(rec.Body.String(), "Carrito") {
		t.Errorf("el backlog debe seguir legible: %s", rec.Body.String())
	}
}

func TestProyectoCerrado_AlReabrirVuelveAAceptarEscrituras(t *testing.T) {
	c := proyectoConDatosDeAna(t)
	c.esperar("ana", http.MethodPost, "/projects/1/cerrar", ``, http.StatusOK)
	c.esperar("ana", http.MethodPost, "/backlog", `{"proyecto_id":1,"titulo":"Otra","descripcion":"d","prioridad":"M"}`, http.StatusConflict)

	c.esperar("ana", http.MethodPost, "/projects/1/reabrir", ``, http.StatusOK)

	c.esperar("ana", http.MethodPost, "/backlog", `{"proyecto_id":1,"titulo":"Otra","descripcion":"d","prioridad":"M"}`, http.StatusCreated)
}

// SC-003: tras cerrar y reabrir, el estado es el mismo que antes del cierre.
func TestProyectoCerrado_ReabrirDevuelveElEstadoPrevio(t *testing.T) {
	c := nuevoClienteDePrueba(t)
	c.registrar("ana")
	hoy := time.Now()
	inicio, fin := hoy.AddDate(0, 0, -7).Format(time.DateOnly), hoy.AddDate(0, 0, 30).Format(time.DateOnly)
	c.esperar("ana", http.MethodPost, "/projects", `{"nombre":"Tienda","fecha_inicio":"`+inicio+`","fecha_fin":"`+fin+`"}`, http.StatusCreated)
	estado := func() string {
		return strings.TrimSpace(c.esperar("ana", http.MethodGet, "/projects/1/status", "", http.StatusOK).Body.String())
	}
	antes := estado()
	if antes != `{"estado":"En curso"}` {
		t.Fatalf("el proyecto debía empezar En curso, se obtuvo %s", antes)
	}

	c.esperar("ana", http.MethodPost, "/projects/1/cerrar", ``, http.StatusOK)
	if cerrado := estado(); cerrado != `{"estado":"Finalizado"}` {
		t.Fatalf("cerrado debía estar Finalizado, se obtuvo %s", cerrado)
	}
	c.esperar("ana", http.MethodPost, "/projects/1/reabrir", ``, http.StatusOK)

	if despues := estado(); despues != antes {
		t.Errorf("tras reabrir se esperaba %s, se obtuvo %s", antes, despues)
	}
}
