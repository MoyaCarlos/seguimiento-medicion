package app

import (
	"net/http"
	"strings"
	"testing"
)

// Escenario 8 de HU-21 contra la API real: el Scrum Master ve a los usuarios
// registrados y arma el equipo solo con ellos, como Product Builder.
func TestEquipo_SoloConUsuariosRegistradosComoProductBuilder(t *testing.T) {
	c := nuevoClienteDePrueba(t)
	c.registrar("ana")
	c.registrar("bruno")
	c.esperar("ana", http.MethodPost, "/projects", `{"nombre":"Tienda"}`, http.StatusCreated)

	usuarios := strings.TrimSpace(c.esperar("ana", http.MethodGet, "/users", "", http.StatusOK).Body.String())
	if usuarios != `[{"id":1,"nombre":"ana"},{"id":2,"nombre":"bruno"}]` {
		t.Errorf("GET /users inesperado: %s", usuarios)
	}
	c.esperar("", http.MethodGet, "/users", "", http.StatusUnauthorized)

	alta := c.esperar("ana", http.MethodPost, "/projects/1/members", `{"nombre":"bruno"}`, http.StatusCreated).Body.String()
	if !strings.Contains(alta, `"rol":"product_builder"`) {
		t.Errorf("bruno debía entrar como product_builder: %s", alta)
	}
	noEncontrado := c.esperar("ana", http.MethodPost, "/projects/1/members", `{"nombre":"carla"}`, http.StatusNotFound).Body.String()
	if strings.TrimSpace(noEncontrado) != `{"mensaje":"usuario no encontrado"}` {
		t.Errorf("cuerpo inesperado para un usuario sin cuenta: %s", noEncontrado)
	}
	c.esperar("ana", http.MethodPost, "/projects/1/members", `{"nombre":"bruno","rol":"scrum_master"}`, http.StatusBadRequest)
	c.registrar("carla")
	c.esperar("ana", http.MethodPost, "/projects/1/members", `{"nombre":"carla","rol":"scrum_master"}`, http.StatusConflict)

	if usuarios := strings.TrimSpace(c.esperar("ana", http.MethodGet, "/users", "", http.StatusOK).Body.String()); strings.Contains(usuarios, "hash") || strings.Contains(usuarios, "contrasena") {
		t.Errorf("GET /users no debe exponer credenciales: %s", usuarios)
	}
}
