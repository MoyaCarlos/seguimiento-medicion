package app

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

const soloScrumMaster = `{"mensaje":"Solo el Scrum Master puede realizar esta acción"}`

// proyectoConEquipo arma el proyecto 1 de ana (Scrum Master) con bruno como
// Product Builder, la historia 1 y el Sprint 1 Pendiente con esa historia.
func proyectoConEquipo(t *testing.T) *clienteDePrueba {
	c := nuevoClienteDePrueba(t)
	c.registrar("ana")
	c.registrar("bruno")
	c.esperar("ana", http.MethodPost, "/projects", `{"nombre":"Tienda"}`, http.StatusCreated)
	c.esperar("ana", http.MethodPost, "/projects/1/members", `{"nombre":"bruno","rol":"product_builder"}`, http.StatusCreated)
	c.esperar("ana", http.MethodPost, "/backlog", `{"proyecto_id":1,"titulo":"Carrito","descripcion":"d","prioridad":"M"}`, http.StatusCreated)
	c.esperar("ana", http.MethodPost, "/sprints", `{"proyecto_id":1}`, http.StatusCreated)
	c.esperar("bruno", http.MethodPut, "/backlog/1/sprint", `{"sprint_id":1}`, http.StatusOK)
	return c
}

func fechas() (string, string) {
	hoy := time.Now()
	return hoy.Format(time.DateOnly), hoy.AddDate(0, 0, 7).Format(time.DateOnly)
}

// Escenario 2 de HU-21: el Product Builder recibe 403 en cada acción de
// gestión, también con datos inválidos, y nada cambia.
func TestGestion_UnProductBuilderRecibe403EnCadaAccionDeGestion(t *testing.T) {
	c := proyectoConEquipo(t)
	inicio, fin := fechas()
	c.esperar("ana", http.MethodPost, "/backlog/1/votacion", ``, http.StatusCreated)

	pedidos := []struct{ nombre, metodo, ruta, cuerpo string }{
		{"crear un Sprint", http.MethodPost, "/sprints", `{"proyecto_id":1}`},
		{"iniciar el Sprint", http.MethodPost, "/sprints/1/iniciar", `{"sprint_goal":"g","fecha_inicio":"` + inicio + `","fecha_fin":"` + fin + `"}`},
		{"iniciar el Sprint con fechas inválidas", http.MethodPost, "/sprints/1/iniciar", `{"sprint_goal":"","fecha_inicio":"no-es-fecha","fecha_fin":"x"}`},
		{"cerrar el Sprint", http.MethodPost, "/sprints/1/cerrar", ``},
		{"iniciar la votación", http.MethodPost, "/backlog/1/votacion", ``},
		{"revelar la votación", http.MethodPost, "/backlog/1/votacion/revelar", ``},
		{"abrir una nueva ronda", http.MethodPost, "/backlog/1/votacion/rondas", ``},
		{"registrar la estimación", http.MethodPost, "/backlog/1/votacion/estimacion", `{"estimacion_sp":5}`},
		{"editar el proyecto", http.MethodPut, "/projects/1", `{"nombre":"Otro","descripcion":"","fecha_inicio":null,"fecha_fin":null}`},
		{"editar el proyecto con datos inválidos", http.MethodPut, "/projects/1", `{"nombre":"","descripcion":"","fecha_inicio":"no-es-fecha","fecha_fin":null}`},
		{"configurar el equipo", http.MethodPost, "/projects/1/members", `{"nombre":"carla","rol":"product_builder"}`},
	}
	for _, p := range pedidos {
		t.Run(p.nombre, func(t *testing.T) {
			rec := c.esperar("bruno", p.metodo, p.ruta, p.cuerpo, http.StatusForbidden)
			if strings.TrimSpace(rec.Body.String()) != soloScrumMaster {
				t.Errorf("cuerpo = %s, se esperaba %s", rec.Body.String(), soloScrumMaster)
			}
		})
	}

	panel := c.esperar("ana", http.MethodGet, "/projects/1", "", http.StatusOK).Body.String()
	if !strings.Contains(panel, `"nombre":"Tienda"`) {
		t.Errorf("el proyecto no debía cambiar: %s", panel)
	}
	sprints := c.esperar("ana", http.MethodGet, "/projects/1/sprints", "", http.StatusOK).Body.String()
	if strings.Count(sprints, `"id":`) != 1 || !strings.Contains(sprints, `"Pendiente"`) {
		t.Errorf("el Sprint debía seguir siendo uno solo y Pendiente: %s", sprints)
	}
}

// Escenario 1 de HU-21: el Scrum Master hace cada acción de gestión.
func TestGestion_ElScrumMasterRealizaCadaAccionDeGestion(t *testing.T) {
	c := proyectoConEquipo(t)
	inicio, fin := fechas()

	c.esperar("ana", http.MethodPut, "/projects/1", `{"nombre":"Tienda 2","descripcion":"","fecha_inicio":null,"fecha_fin":null}`, http.StatusOK)
	c.esperar("ana", http.MethodPost, "/sprints/1/iniciar", `{"sprint_goal":"g","fecha_inicio":"`+inicio+`","fecha_fin":"`+fin+`"}`, http.StatusOK)
	c.esperar("ana", http.MethodPost, "/backlog/1/votacion", ``, http.StatusCreated)
	c.esperar("ana", http.MethodPost, "/backlog/1/votacion/votos", `{"valor":3}`, http.StatusOK)
	c.esperar("bruno", http.MethodPost, "/backlog/1/votacion/votos", `{"valor":5}`, http.StatusOK)
	c.esperar("ana", http.MethodPost, "/backlog/1/votacion/revelar", ``, http.StatusOK)
	c.esperar("ana", http.MethodPost, "/backlog/1/votacion/rondas", ``, http.StatusCreated)
	c.esperar("ana", http.MethodPost, "/backlog/1/votacion/votos", `{"valor":5}`, http.StatusOK)
	c.esperar("bruno", http.MethodPost, "/backlog/1/votacion/votos", `{"valor":5}`, http.StatusOK)
	c.esperar("ana", http.MethodPost, "/backlog/1/votacion/revelar", ``, http.StatusOK)
	c.esperar("ana", http.MethodPost, "/backlog/1/votacion/estimacion", `{"estimacion_sp":5}`, http.StatusOK)
	c.esperar("ana", http.MethodPost, "/sprints/1/cerrar", ``, http.StatusOK)
	c.esperar("ana", http.MethodPost, "/sprints", `{"proyecto_id":1}`, http.StatusCreated)
}
