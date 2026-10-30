package features

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/cucumber/godog"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/demo"
	apphttp "github.com/MoyaCarlos/seguimiento-medicion/internal/http"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

var tablasDeDatosDePrueba = []string{"usuarios", "proyectos", "integrantes", "sprints", "historias",
	"historias_arrastradas", "registros_horas", "rondas_votacion", "votos"}

// pasosDatosPrueba ejercita POST /demo/datos sobre una base real en memoria.
type pasosDatosPrueba struct {
	db         *sql.DB
	habilitado bool
	respuesta  *httptest.ResponseRecorder
	antes      map[string]int
}

func (p *pasosDatosPrueba) iniciar() error {
	db, err := repository.AbrirSQLite(":memory:")
	if err != nil {
		return err
	}
	if err := repository.AplicarMigraciones(context.Background(), db); err != nil {
		return err
	}
	*p = pasosDatosPrueba{db: db}
	return nil
}

func (p *pasosDatosPrueba) cerrar() {
	if p.db != nil {
		_ = p.db.Close()
		p.db = nil
	}
}

func (p *pasosDatosPrueba) pedir(metodo, ruta string) *httptest.ResponseRecorder {
	mux := http.NewServeMux()
	cargador := demo.NewCargador(p.db, time.Now)
	apphttp.NewDemoHandler(func() bool { return p.habilitado }, cargador).Registrar(mux)
	respuesta := httptest.NewRecorder()
	mux.ServeHTTP(respuesta, httptest.NewRequest(metodo, ruta, nil))
	return respuesta
}

func (p *pasosDatosPrueba) contar(tabla string) (int, error) {
	var n int
	err := p.db.QueryRow(`SELECT COUNT(*) FROM ` + tabla).Scan(&n)
	return n, err
}

func (p *pasosDatosPrueba) foto() (map[string]int, error) {
	foto := map[string]int{}
	for _, tabla := range tablasDeDatosDePrueba {
		n, err := p.contar(tabla)
		if err != nil {
			return nil, err
		}
		foto[tabla] = n
	}
	return foto, nil
}

func (p *pasosDatosPrueba) habilitada() error    { p.habilitado = true; return nil }
func (p *pasosDatosPrueba) deshabilitada() error { p.habilitado = false; return nil }

func (p *pasosDatosPrueba) baseSinDatosDePrueba() error {
	n, err := p.contar("proyectos")
	if err != nil || n != 0 {
		return fmt.Errorf("la base debía estar vacía (proyectos: %d, error: %v)", n, err)
	}
	return nil
}

func (p *pasosDatosPrueba) yaCargados() error {
	p.habilitado = true
	if rec := p.pedir("POST", "/demo/datos"); rec.Code != http.StatusCreated {
		return fmt.Errorf("la carga previa respondió %d", rec.Code)
	}
	return nil
}

func (p *pasosDatosPrueba) otroProyectoUsaElId25() error {
	_, err := p.db.Exec(`INSERT INTO proyectos (id, nombre, descripcion, creado_en) VALUES (25, 'Proyecto real', '', '2026-01-01T00:00:00Z')`)
	return err
}

func (p *pasosDatosPrueba) pidoCargar() error {
	antes, err := p.foto()
	if err != nil {
		return err
	}
	p.antes = antes
	p.respuesta = p.pedir("POST", "/demo/datos")
	return nil
}

func (p *pasosDatosPrueba) respondeConResumen() error {
	if p.respuesta.Code != http.StatusCreated {
		return fmt.Errorf("estado %d, se esperaba 201: %s", p.respuesta.Code, p.respuesta.Body)
	}
	if !strings.Contains(p.respuesta.Body.String(), `"proyecto_id":25`) {
		return fmt.Errorf("el resumen no trae el proyecto 25: %s", p.respuesta.Body)
	}
	return nil
}

func (p *pasosDatosPrueba) proyectoTiene(historias, sprints, integrantes int) error {
	esperadas := map[string]int{"historias": historias, "sprints": sprints, "integrantes": integrantes}
	for tabla, esperado := range esperadas {
		var n int
		if err := p.db.QueryRow(`SELECT COUNT(*) FROM ` + tabla + ` WHERE proyecto_id = 25`).Scan(&n); err != nil {
			return err
		}
		if n != esperado {
			return fmt.Errorf("%s del proyecto 25: %d, se esperaban %d", tabla, n, esperado)
		}
	}
	return nil
}

func (p *pasosDatosPrueba) respondeConflicto(mensaje string) error {
	if p.respuesta.Code != http.StatusConflict {
		return fmt.Errorf("estado %d, se esperaba 409", p.respuesta.Code)
	}
	if !strings.Contains(p.respuesta.Body.String(), mensaje) {
		return fmt.Errorf("el cuerpo %q no contiene %q", p.respuesta.Body, mensaje)
	}
	return nil
}

func (p *pasosDatosPrueba) respondeNoEncontrado() error {
	if p.respuesta.Code != http.StatusNotFound {
		return fmt.Errorf("estado %d, se esperaba 404", p.respuesta.Code)
	}
	return nil
}

func (p *pasosDatosPrueba) nadaCambio() error {
	despues, err := p.foto()
	if err != nil {
		return err
	}
	for tabla, filas := range despues {
		if filas != p.antes[tabla] {
			return fmt.Errorf("%s cambió: %d -> %d", tabla, p.antes[tabla], filas)
		}
	}
	return nil
}

// InitializeScenarioDatosPrueba registra los escenarios del issue #144.
func InitializeScenarioDatosPrueba(ctx *godog.ScenarioContext) {
	p := &pasosDatosPrueba{}
	ctx.Before(func(goctx context.Context, _ *godog.Scenario) (context.Context, error) {
		p.cerrar()
		return goctx, p.iniciar()
	})
	ctx.After(func(goctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		p.cerrar()
		return goctx, nil
	})

	ctx.Step(`^que la carga de datos de prueba está habilitada$`, p.habilitada)
	ctx.Step(`^que la carga de datos de prueba está deshabilitada$`, p.deshabilitada)
	ctx.Step(`^que la base no tiene datos de prueba$`, p.baseSinDatosDePrueba)
	ctx.Step(`^que los datos de prueba ya fueron cargados$`, p.yaCargados)
	ctx.Step(`^que otro proyecto ya usa el identificador 25$`, p.otroProyectoUsaElId25)
	ctx.Step(`^pido cargar los datos de prueba$`, p.pidoCargar)
	ctx.Step(`^el sistema responde 201 con el resumen del proyecto 25$`, p.respondeConResumen)
	ctx.Step(`^el proyecto 25 tiene (\d+) historias, (\d+) Sprints y (\d+) integrantes$`, p.proyectoTiene)
	ctx.Step(`^el sistema responde 409 con el mensaje "([^"]*)"$`, p.respondeConflicto)
	ctx.Step(`^el sistema responde 404 como si la ruta no existiera$`, p.respondeNoEncontrado)
	ctx.Step(`^la carga no agregó ni cambió ningún dato$`, p.nadaCambio)
}
