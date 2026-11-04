package http

import (
	"bytes"
	nethttp "net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

const datosEdicionProyecto = `"nombre":"Corregido","descripcion":"","fecha_inicio":null,"fecha_fin":null`

func editarProyectoViaHTTP(handler *ProjectHandler, id int64, cuerpo string) *httptest.ResponseRecorder {
	idStr := strconv.FormatInt(id, 10)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(nethttp.MethodPut, "/projects/"+idStr, bytes.NewBufferString(cuerpo))
	req.SetPathValue("id", idStr)
	handler.Editar(rec, req)
	return rec
}

func horasDelProyecto(t *testing.T, handler *ProjectHandler, id int64) any {
	t.Helper()
	idStr := strconv.FormatInt(id, 10)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(nethttp.MethodGet, "/projects/"+idStr, nil)
	req.SetPathValue("id", idStr)
	handler.Obtener(rec, req)
	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	horas, existe := cuerpoJSON(t, rec)["horas_por_story_point"]
	if !existe {
		t.Fatalf("la respuesta no incluye horas_por_story_point: %s", rec.Body.String())
	}
	return horas
}

func TestProjectHandler_Editar_HorasPorStoryPoint(t *testing.T) {
	casos := []struct {
		nombre   string
		valor    string
		esperado any
	}{
		{"un entero", "4", 4.0},
		{"con decimales", "2.5", 2.5},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			handler, _ := nuevoProjectHandlerDePrueba()
			id := crearProyectoViaHTTP(t, handler, "Original")

			rec := editarProyectoViaHTTP(handler, id, `{`+datosEdicionProyecto+`,"horas_por_story_point":`+c.valor+`}`)

			if rec.Code != nethttp.StatusOK {
				t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
			}
			if horas := cuerpoJSON(t, rec)["horas_por_story_point"]; horas != c.esperado {
				t.Errorf("la respuesta: se esperaba %v, se obtuvo %v", c.esperado, horas)
			}
			if horas := horasDelProyecto(t, handler, id); horas != c.esperado {
				t.Errorf("la consulta: se esperaba %v, se obtuvo %v", c.esperado, horas)
			}
		})
	}
}

func TestProjectHandler_Editar_SinMencionarElParametroLoConserva(t *testing.T) {
	handler, _ := nuevoProjectHandlerDePrueba()
	id := crearProyectoViaHTTP(t, handler, "Original")
	editarProyectoViaHTTP(handler, id, `{`+datosEdicionProyecto+`,"horas_por_story_point":4}`)

	rec := editarProyectoViaHTTP(handler, id, `{`+datosEdicionProyecto+`}`)

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	if horas := horasDelProyecto(t, handler, id); horas != 4.0 {
		t.Errorf("se esperaba conservar 4, se obtuvo %v", horas)
	}
}

func TestProjectHandler_Editar_ConNullBorraElParametro(t *testing.T) {
	handler, _ := nuevoProjectHandlerDePrueba()
	id := crearProyectoViaHTTP(t, handler, "Original")
	editarProyectoViaHTTP(handler, id, `{`+datosEdicionProyecto+`,"horas_por_story_point":4}`)

	rec := editarProyectoViaHTTP(handler, id, `{`+datosEdicionProyecto+`,"horas_por_story_point":null}`)

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	if horas := horasDelProyecto(t, handler, id); horas != nil {
		t.Errorf("se esperaba sin definir, se obtuvo %v", horas)
	}
}

func TestProjectHandler_Editar_HorasPorStoryPointInvalidas(t *testing.T) {
	for _, valor := range []string{"0", "-1", `"abc"`, `""`, `"5"`, "true", "[]", "{}"} {
		t.Run(valor, func(t *testing.T) {
			handler, _ := nuevoProjectHandlerDePrueba()
			id := crearProyectoViaHTTP(t, handler, "Original")
			editarProyectoViaHTTP(handler, id, `{`+datosEdicionProyecto+`,"horas_por_story_point":4}`)

			rec := editarProyectoViaHTTP(handler, id, `{`+datosEdicionProyecto+`,"horas_por_story_point":`+valor+`}`)

			esperarError(t, rec, nethttp.StatusBadRequest, "horas_por_story_point", "las horas por Story Point deben ser un número mayor que cero")
			if horas := horasDelProyecto(t, handler, id); horas != 4.0 {
				t.Errorf("se esperaba conservar 4, se obtuvo %v", horas)
			}
		})
	}
}

func TestProjectHandler_Editar_LasHorasPorStoryPointSeValidanAlFinal(t *testing.T) {
	handler, _ := nuevoProjectHandlerDePrueba()
	id := crearProyectoViaHTTP(t, handler, "Original")

	rec := editarProyectoViaHTTP(handler, id, `{"nombre":"  ","descripcion":"","fecha_inicio":null,"fecha_fin":null,"horas_por_story_point":"abc"}`)

	esperarError(t, rec, nethttp.StatusBadRequest, "nombre", "")
}

func TestProjectHandler_ElParametroEsNullSiNoEstaDefinido(t *testing.T) {
	handler, _ := nuevoProjectHandlerDePrueba()

	alta := httptest.NewRecorder()
	handler.Crear(alta, httptest.NewRequest(nethttp.MethodPost, "/projects", bytes.NewBufferString(`{"nombre":"Nuevo","creador":"Ana"}`)))
	if horas, existe := cuerpoJSON(t, alta)["horas_por_story_point"]; !existe || horas != nil {
		t.Errorf("el alta: se esperaba null, se obtuvo %v (existe: %v)", horas, existe)
	}

	listado := httptest.NewRecorder()
	handler.Listar(listado, httptest.NewRequest(nethttp.MethodGet, "/projects", nil))
	if !bytes.Contains(listado.Body.Bytes(), []byte(`"horas_por_story_point":null`)) {
		t.Errorf("el listado no incluye el parámetro: %s", listado.Body.String())
	}
}
