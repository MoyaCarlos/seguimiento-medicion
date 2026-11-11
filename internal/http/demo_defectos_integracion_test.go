package http

import (
	"context"
	"encoding/json"
	nethttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/demo"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

func TestDemo_LosDefectosCargadosSeListanEnElProyecto25(t *testing.T) {
	db, err := repository.AbrirSQLite(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := repository.AplicarMigraciones(context.Background(), db); err != nil {
		t.Fatal(err)
	}
	cargador := demo.NewCargador(db, time.Now, service.NewServicioContrasenas(1))
	if _, err := cargador.Cargar(context.Background()); err != nil {
		t.Fatal(err)
	}
	listar := service.NewListarDefectos(repository.NewSQLiteProjectRepository(db), repository.NewSQLiteDefectoRepository(db))
	mux := nethttp.NewServeMux()
	mux.HandleFunc("GET /projects/{id}/defectos", NewDefectoHandler(nil, nil, nil, listar).Listar)

	respuesta := httptest.NewRecorder()
	mux.ServeHTTP(respuesta, httptest.NewRequest("GET", "/projects/25/defectos", nil))

	if respuesta.Code != nethttp.StatusOK {
		t.Fatalf("estado %d, se esperaba 200", respuesta.Code)
	}
	var defectos []map[string]any
	if err := json.Unmarshal(respuesta.Body.Bytes(), &defectos); err != nil {
		t.Fatal(err)
	}
	if len(defectos) != 5 {
		t.Fatalf("se esperaban 5 defectos, hay %d", len(defectos))
	}
	estados := map[string]int{}
	for _, d := range defectos {
		estado := d["estado"].(string)
		estados[estado]++
		if (estado == "Resuelto") != (d["sprint_resolucion_id"] != nil) {
			t.Errorf("defecto %v: estado %s y Sprint de resolución %v incoherentes", d["id"], estado, d["sprint_resolucion_id"])
		}
		if d["historia"] == "" || d["descripcion"] == "" {
			t.Errorf("defecto %v sin historia o descripción", d["id"])
		}
	}
	if estados["Resuelto"] != 3 || estados["Pendiente"] != 2 {
		t.Errorf("estados %v, se esperaban 3 Resueltos y 2 Pendientes", estados)
	}
}
