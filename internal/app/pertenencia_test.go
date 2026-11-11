package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// clienteDePrueba pide a la API real con la sesión de cada persona.
type clienteDePrueba struct {
	t      *testing.T
	mux    *http.ServeMux
	tokens map[string]string
}

func nuevoClienteDePrueba(t *testing.T) *clienteDePrueba {
	return &clienteDePrueba{t: t, mux: apiDePrueba(t), tokens: map[string]string{}}
}

func (c *clienteDePrueba) registrar(nombre string) {
	c.t.Helper()
	rec := c.pedir("", http.MethodPost, "/auth/registro", `{"nombre":"`+nombre+`","contrasena":"una-clave-larga"}`)
	var resp struct{ Token string }
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil || resp.Token == "" {
		c.t.Fatalf("no se pudo registrar a %s: %d %s", nombre, rec.Code, rec.Body.String())
	}
	c.tokens[nombre] = resp.Token
}

func (c *clienteDePrueba) pedir(quien, metodo, ruta, cuerpo string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(metodo, ruta, strings.NewReader(cuerpo))
	if quien != "" {
		req.Header.Set("Authorization", "Bearer "+c.tokens[quien])
	}
	rec := httptest.NewRecorder()
	c.mux.ServeHTTP(rec, req)
	return rec
}

func (c *clienteDePrueba) esperar(quien, metodo, ruta, cuerpo string, status int) *httptest.ResponseRecorder {
	c.t.Helper()
	rec := c.pedir(quien, metodo, ruta, cuerpo)
	if rec.Code != status {
		c.t.Fatalf("%s %s %s: se esperaba %d, se obtuvo %d (%s)", quien, metodo, ruta, status, rec.Code, rec.Body.String())
	}
	return rec
}

// proyectoDeAna arma el proyecto 1 de ana con la historia 1 en votación y deja
// a diego con cuenta pero fuera del proyecto.
func proyectoDeAna(t *testing.T) *clienteDePrueba {
	c := nuevoClienteDePrueba(t)
	c.registrar("ana")
	c.registrar("diego")
	c.esperar("ana", http.MethodPost, "/projects", `{"nombre":"Tienda"}`, http.StatusCreated)
	c.esperar("ana", http.MethodPost, "/backlog", `{"proyecto_id":1,"titulo":"Carrito","descripcion":"d","prioridad":"M"}`, http.StatusCreated)
	c.esperar("ana", http.MethodPost, "/backlog/1/votacion", ``, http.StatusCreated)
	return c
}

const noSosIntegrante = `{"mensaje":"No sos integrante de este proyecto"}`

// Escenario 4 de HU-21.
func TestPertenencia_UnNoIntegranteRecibe403EnElPanelLaVotacionYLasEscrituras(t *testing.T) {
	c := proyectoDeAna(t)
	pedidos := []struct{ metodo, ruta, cuerpo string }{
		{http.MethodGet, "/projects/1", ""},
		{http.MethodGet, "/backlog/1/votacion", ""},
		{http.MethodPost, "/backlog", `{"proyecto_id":1,"titulo":"Otra","descripcion":"d","prioridad":"M"}`},
		{http.MethodPut, "/backlog/1", `{"titulo":"Cambiada","descripcion":"d","criterios_aceptacion":null,"estimacion_sp":null}`},
		{http.MethodPut, "/backlog/1/estado", `{"estado":"En Progreso"}`},
		{http.MethodPut, "/backlog/1/sprint", `{"sprint_id":1}`},
		{http.MethodPost, "/backlog/1/horas", `{"fecha":"2026-11-11","actividad":"x","horas":1}`},
		{http.MethodPost, "/projects/1/defectos", `{"historia_id":1,"descripcion":"d","severidad":"Alta"}`},
		{http.MethodPost, "/backlog/1/votacion/votos", `{"valor":5}`},
	}
	for _, p := range pedidos {
		t.Run(p.metodo+" "+p.ruta, func(t *testing.T) {
			rec := c.esperar("diego", p.metodo, p.ruta, p.cuerpo, http.StatusForbidden)
			if strings.TrimSpace(rec.Body.String()) != noSosIntegrante {
				t.Errorf("cuerpo = %s, se esperaba %s", rec.Body.String(), noSosIntegrante)
			}
		})
	}
	if rec := c.esperar("ana", http.MethodGet, "/backlog/1", "", http.StatusOK); strings.Contains(rec.Body.String(), "Cambiada") {
		t.Errorf("la historia no debía cambiar: %s", rec.Body.String())
	}
}

// Límite conocido de HU-21: con sesión, un no integrante puede leer los demás
// datos del proyecto.
func TestPertenencia_UnNoIntegrantePuedeLeerLosDemasDatos(t *testing.T) {
	c := proyectoDeAna(t)
	for _, ruta := range []string{"/projects/1/backlog", "/projects/1/sprints", "/projects/1/defectos", "/backlog/1", "/backlog/1/horas"} {
		c.esperar("diego", http.MethodGet, ruta, "", http.StatusOK)
	}
}

func TestPertenencia_ElPanelDeUnProyectoInexistenteEs404(t *testing.T) {
	c := proyectoDeAna(t)
	c.esperar("diego", http.MethodGet, "/projects/999", "", http.StatusNotFound)
}
