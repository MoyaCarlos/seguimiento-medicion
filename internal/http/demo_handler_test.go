package http

import (
	"context"
	"errors"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/demo"
)

type cargadorFalso struct {
	llamadas int
	resumen  demo.Resumen
	err      error
}

func (c *cargadorFalso) Cargar(context.Context) (demo.Resumen, error) {
	c.llamadas++
	return c.resumen, c.err
}

func pedirDemo(habilitado bool, cargador *cargadorFalso, metodo, ruta string) *httptest.ResponseRecorder {
	mux := nethttp.NewServeMux()
	NewDemoHandler(func() bool { return habilitado }, cargador).Registrar(mux)
	respuesta := httptest.NewRecorder()
	mux.ServeHTTP(respuesta, httptest.NewRequest(metodo, ruta, nil))
	return respuesta
}

func TestDemo_DeshabilitadaTodasLasRutasResponden404SinCargar(t *testing.T) {
	for _, caso := range []struct{ metodo, ruta string }{{"GET", "/demo/estado"}, {"POST", "/demo/datos"}, {"POST", "/demo/estado"}, {"DELETE", "/demo/datos"}, {"GET", "/demo/otra"}} {
		cargador := &cargadorFalso{}

		respuesta := pedirDemo(false, cargador, caso.metodo, caso.ruta)

		if respuesta.Code != nethttp.StatusNotFound {
			t.Errorf("%s %s: estado %d, se esperaba 404", caso.metodo, caso.ruta, respuesta.Code)
		}
		if strings.TrimSpace(respuesta.Body.String()) != "404 page not found" {
			t.Errorf("%s %s: cuerpo %q, se esperaba el 404 estándar", caso.metodo, caso.ruta, respuesta.Body.String())
		}
		if cargador.llamadas != 0 {
			t.Errorf("%s %s: no debía cargar nada", caso.metodo, caso.ruta)
		}
	}
}

func TestDemo_HabilitadaElEstadoResponde200(t *testing.T) {
	respuesta := pedirDemo(true, &cargadorFalso{}, "GET", "/demo/estado")

	if respuesta.Code != nethttp.StatusOK || strings.TrimSpace(respuesta.Body.String()) != `{"habilitado":true}` {
		t.Errorf("estado %d, cuerpo %q", respuesta.Code, respuesta.Body.String())
	}
}

func TestDemo_CargaExitosaResponde201ConElResumen(t *testing.T) {
	cargador := &cargadorFalso{resumen: demo.Resumen{ProyectoID: 25, Usuarios: 4, Historias: 10}}

	respuesta := pedirDemo(true, cargador, "POST", "/demo/datos")

	if respuesta.Code != nethttp.StatusCreated {
		t.Fatalf("estado %d, se esperaba 201", respuesta.Code)
	}
	if !strings.Contains(respuesta.Body.String(), `"proyecto_id":25`) || !strings.Contains(respuesta.Body.String(), `"historias":10`) {
		t.Errorf("cuerpo %q", respuesta.Body.String())
	}
}

func TestDemo_ErroresDeLaCarga(t *testing.T) {
	casos := []struct {
		nombre   string
		err      error
		estado   int
		contiene string
	}{
		{"ya cargado", demo.ErrYaCargado, 409, `"campo":"ya_cargado"`},
		{"id ocupado", demo.ErrIdOcupado, 409, `"campo":"proyecto_id"`},
		{"error interno", errors.New("disco lleno"), 500, `"mensaje":"no se pudieron cargar los datos de prueba"`},
	}
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			respuesta := pedirDemo(true, &cargadorFalso{err: caso.err}, "POST", "/demo/datos")

			if respuesta.Code != caso.estado || !strings.Contains(respuesta.Body.String(), caso.contiene) {
				t.Errorf("estado %d, cuerpo %q", respuesta.Code, respuesta.Body.String())
			}
			if strings.Contains(respuesta.Body.String(), "disco lleno") {
				t.Error("el error interno no debe filtrarse al cliente")
			}
		})
	}
}

func TestDemo_LosMensajesDel409ViajanEnElCuerpo(t *testing.T) {
	respuesta := pedirDemo(true, &cargadorFalso{err: demo.ErrYaCargado}, "POST", "/demo/datos")

	if !strings.Contains(respuesta.Body.String(), "los datos de prueba ya están cargados") {
		t.Errorf("cuerpo %q", respuesta.Body.String())
	}
}
