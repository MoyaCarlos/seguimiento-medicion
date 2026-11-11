package app

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

type proyectoConRol struct {
	ID     int64  `json:"id"`
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
}

// Escenarios 5 y 9 de HU-21 contra la API real: GET /projects lista solo los
// proyectos de la sesión con su rol, y [] (nunca 403) a quien no integra ninguno.
func TestMisProyectos_CadaUnoVeLosSuyosConSuRol(t *testing.T) {
	c := nuevoClienteDePrueba(t)
	c.registrar("ana")
	c.registrar("bruno")
	c.registrar("carla")
	c.esperar("ana", http.MethodPost, "/projects", `{"nombre":"Tienda"}`, http.StatusCreated)
	c.esperar("bruno", http.MethodPost, "/projects", `{"nombre":"Blog"}`, http.StatusCreated)
	c.esperar("ana", http.MethodPost, "/projects/1/members", `{"nombre":"bruno","rol":"product_builder"}`, http.StatusCreated)

	var deAna, deBruno []proyectoConRol
	leer(t, c.esperar("ana", http.MethodGet, "/projects", "", http.StatusOK).Body.Bytes(), &deAna)
	leer(t, c.esperar("bruno", http.MethodGet, "/projects", "", http.StatusOK).Body.Bytes(), &deBruno)

	if len(deAna) != 1 || deAna[0] != (proyectoConRol{1, "Tienda", "scrum_master"}) {
		t.Errorf("Ana debía ver solo Tienda como scrum_master: %+v", deAna)
	}
	if len(deBruno) != 2 || deBruno[0] != (proyectoConRol{1, "Tienda", "product_builder"}) ||
		deBruno[1] != (proyectoConRol{2, "Blog", "scrum_master"}) {
		t.Errorf("Bruno debía ver Tienda como product_builder y Blog como scrum_master: %+v", deBruno)
	}
	if body := strings.TrimSpace(c.esperar("carla", http.MethodGet, "/projects", "", http.StatusOK).Body.String()); body != "[]" {
		t.Errorf("Carla no integra proyectos: se esperaba [], se obtuvo %s", body)
	}

	var panel proyectoConRol
	leer(t, c.esperar("bruno", http.MethodGet, "/projects/1", "", http.StatusOK).Body.Bytes(), &panel)
	if panel.Rol != "product_builder" {
		t.Errorf("el panel debía traer el rol de Bruno (product_builder), se obtuvo %q", panel.Rol)
	}
}

func leer(t *testing.T, datos []byte, destino any) {
	t.Helper()
	if err := json.Unmarshal(datos, destino); err != nil {
		t.Fatalf("respuesta JSON inválida: %v (%s)", err, datos)
	}
}
