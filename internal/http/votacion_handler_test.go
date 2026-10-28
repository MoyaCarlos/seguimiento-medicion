package http

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// equipoDeVotacion es el proyecto 1 con Carlos (scrum_master), Candela y
// Jimena (product_builder), un usuario ajeno (9) y la historia 3.
const equipoDeVotacion = `
INSERT INTO proyectos (id, nombre, creado_en) VALUES (1, 'Software Metrics & Estimation', '2026-10-01T00:00:00Z');
INSERT INTO usuarios (id, nombre, nombre_normalizado, creado_en) VALUES
    (1, 'Carlos', 'carlos', '2026-10-01T00:00:00Z'),
    (2, 'Candela', 'candela', '2026-10-01T00:00:00Z'),
    (3, 'Jimena', 'jimena', '2026-10-01T00:00:00Z'),
    (9, 'Ajeno', 'ajeno', '2026-10-01T00:00:00Z');
INSERT INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES
    (1, 1, 'scrum_master', '2026-10-01T00:00:00Z'),
    (1, 2, 'product_builder', '2026-10-01T00:00:00Z'),
    (1, 3, 'product_builder', '2026-10-01T00:00:00Z');
INSERT INTO historias (id, proyecto_id, titulo, descripcion, prioridad, estado) VALUES (3, 1, 't', 'd', 'M', 'Nueva');
`

const (
	rondaAbierta          = "INSERT INTO rondas_votacion (id, historia_id, numero, revelada_en) VALUES (1, 3, 1, NULL)"
	rondaRevelada         = "INSERT INTO rondas_votacion (id, historia_id, numero, revelada_en) VALUES (1, 3, 1, '2026-10-27T15:04:05Z')"
	votosDeCarlosYCandela = "INSERT INTO votos (ronda_id, usuario_id, valor) VALUES (1, 1, 3), (1, 2, 8)"
)

func nuevoVotacionHandlerDePrueba(t *testing.T, sentencias ...string) *VotacionHandler {
	t.Helper()
	db, err := repository.AbrirSQLite(":memory:")
	if err != nil {
		t.Fatalf("no se pudo abrir sqlite: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := repository.AplicarMigraciones(context.Background(), db); err != nil {
		t.Fatalf("no se pudieron aplicar las migraciones: %v", err)
	}
	for _, s := range append([]string{equipoDeVotacion}, sentencias...) {
		ejecutar(t, db, s)
	}

	historias := repository.NewSQLiteBacklogRepository(db)
	proyectos := repository.NewSQLiteProjectRepository(db)
	votaciones := repository.NewSQLiteVotacionRepository(db)
	return NewVotacionHandler(
		service.NewObtenerHistoria(historias),
		service.NewVotar(historias, proyectos, votaciones),
		service.NewConsultarVotacion(historias, proyectos, votaciones),
	)
}

func ejecutar(t *testing.T, db *sql.DB, sentencia string) {
	t.Helper()
	if _, err := db.Exec(sentencia); err != nil {
		t.Fatalf("no se pudo ejecutar %q: %v", sentencia, err)
	}
}

func votarEn(handler *VotacionHandler, id, cuerpo string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(nethttp.MethodPost, "/backlog/"+id+"/votacion/votos", bytes.NewBufferString(cuerpo))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", id)
	handler.Votar(rec, req)
	return rec
}

func consultarVotacionDe(handler *VotacionHandler, id, consulta string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(nethttp.MethodGet, "/backlog/"+id+"/votacion"+consulta, nil)
	req.SetPathValue("id", id)
	handler.Consultar(rec, req)
	return rec
}

func leerVotacion(t *testing.T, rec *httptest.ResponseRecorder) votacionResponse {
	t.Helper()
	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	var resp votacionResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta JSON inválida: %v", err)
	}
	return resp
}

// esperarValores compara el valor visible de cada integrante; -1 = null.
func esperarValores(t *testing.T, resp votacionResponse, esperados ...int) {
	t.Helper()
	vistos := make([]int, len(resp.Integrantes))
	for i, v := range resp.Integrantes {
		vistos[i] = -1
		if v.Valor != nil {
			vistos[i] = *v.Valor
		}
	}
	if len(vistos) != len(esperados) {
		t.Fatalf("se esperaban los valores %v, se obtuvo %v", esperados, vistos)
	}
	for i := range esperados {
		if vistos[i] != esperados[i] {
			t.Fatalf("se esperaban los valores %v, se obtuvo %v", esperados, vistos)
		}
	}
}

func esperarError(t *testing.T, rec *httptest.ResponseRecorder, status int, campo, mensaje string) {
	t.Helper()
	if rec.Code != status {
		t.Fatalf("se esperaba %d, se obtuvo %d (%s)", status, rec.Code, rec.Body.String())
	}
	resp := leerError(t, rec)
	if resp.Campo != campo {
		t.Errorf("se esperaba campo %q, se obtuvo %q", campo, resp.Campo)
	}
	if mensaje != "" && resp.Mensaje != mensaje {
		t.Errorf("se esperaba mensaje %q, se obtuvo %q", mensaje, resp.Mensaje)
	}
}

func TestVotacionHandler_Votar_Exitoso(t *testing.T) {
	handler := nuevoVotacionHandlerDePrueba(t, rondaAbierta,
		"INSERT INTO votos (ronda_id, usuario_id, valor) VALUES (1, 1, 3)")

	rec := votarEn(handler, "3", `{"usuario_id":2,"valor":8}`)

	resp := leerVotacion(t, rec)
	if resp.HistoriaID != 3 || resp.Ronda != 1 {
		t.Errorf("se esperaba la ronda 1 de la historia 3, se obtuvo %+v", resp)
	}
	esperarValores(t, resp, -1, 8, -1)
	carlos := resp.Integrantes[0]
	if carlos.UsuarioID != 1 || carlos.Nombre != "Carlos" || carlos.Rol != "scrum_master" || !carlos.YaVoto {
		t.Errorf("se esperaba a Carlos, scrum_master, con el voto oculto, se obtuvo %+v", carlos)
	}
	if resp.Integrantes[2].YaVoto {
		t.Errorf("Jimena todavía no votó: %+v", resp.Integrantes[2])
	}
	cuerpo := rec.Body.String()
	if !strings.Contains(cuerpo, `"revelada_en":null`) || !strings.Contains(cuerpo, `"promedio":null`) {
		t.Errorf("revelada_en y promedio van en null mientras la ronda está oculta: %s", cuerpo)
	}
}

func TestVotacionHandler_Votar_Errores(t *testing.T) {
	casos := map[string]struct {
		id, cuerpo string
		previas    []string
		status     int
		campo      string
		mensaje    string
	}{
		"id no numérico":                     {"abc", `{"usuario_id":2,"valor":8}`, nil, nethttp.StatusBadRequest, "id", ""},
		"historia inexistente con JSON roto": {"999", `{no-es-json`, nil, nethttp.StatusNotFound, "", "historia no encontrada"},
		"JSON inválido":                      {"3", `{no-es-json`, []string{rondaAbierta}, nethttp.StatusBadRequest, "", "JSON inválido"},
		"sin usuario_id":                     {"3", `{"valor":8}`, []string{rondaAbierta}, nethttp.StatusBadRequest, "usuario_id", ""},
		"usuario_id no positivo":             {"3", `{"usuario_id":0,"valor":8}`, []string{rondaAbierta}, nethttp.StatusBadRequest, "usuario_id", ""},
		"no es integrante":                   {"3", `{"usuario_id":9,"valor":8}`, []string{rondaAbierta}, nethttp.StatusForbidden, "", "el usuario no es integrante del proyecto"},
		"sin votación":                       {"3", `{"usuario_id":2,"valor":8}`, nil, nethttp.StatusNotFound, "", "la historia no tiene votación iniciada"},
		"carta fuera de la baraja":           {"3", `{"usuario_id":2,"valor":4}`, []string{rondaAbierta}, nethttp.StatusBadRequest, "valor", ""},
		"ronda revelada":                     {"3", `{"usuario_id":2,"valor":8}`, []string{rondaRevelada}, nethttp.StatusBadRequest, "estado", ""},
	}
	for nombre, c := range casos {
		t.Run(nombre, func(t *testing.T) {
			handler := nuevoVotacionHandlerDePrueba(t, c.previas...)

			rec := votarEn(handler, c.id, c.cuerpo)

			esperarError(t, rec, c.status, c.campo, c.mensaje)
		})
	}
}

func TestVotacionHandler_Consultar_SoloElVotoPropioEsVisible(t *testing.T) {
	handler := nuevoVotacionHandlerDePrueba(t, rondaAbierta, votosDeCarlosYCandela)

	resp := leerVotacion(t, consultarVotacionDe(handler, "3", "?usuario_id=2"))

	esperarValores(t, resp, -1, 8, -1)
	if !resp.Integrantes[0].YaVoto || resp.Integrantes[2].YaVoto {
		t.Errorf("Carlos votó (oculto) y Jimena no: %+v", resp.Integrantes)
	}
	if resp.ReveladaEn != nil || resp.Promedio != nil {
		t.Errorf("una ronda oculta no tiene revelada_en ni promedio: %+v", resp)
	}
}

func TestVotacionHandler_Consultar_SinUsuarioTodoOculto(t *testing.T) {
	handler := nuevoVotacionHandlerDePrueba(t, rondaAbierta, votosDeCarlosYCandela)

	resp := leerVotacion(t, consultarVotacionDe(handler, "3", ""))

	esperarValores(t, resp, -1, -1, -1)
}

func TestVotacionHandler_Consultar_Revelada(t *testing.T) {
	handler := nuevoVotacionHandlerDePrueba(t, rondaRevelada,
		"INSERT INTO votos (ronda_id, usuario_id, valor) VALUES (1, 1, 3), (1, 2, 8), (1, 3, 5)")

	resp := leerVotacion(t, consultarVotacionDe(handler, "3", ""))

	esperarValores(t, resp, 3, 8, 5)
	if resp.ReveladaEn == nil || *resp.ReveladaEn != "2026-10-27T15:04:05Z" {
		t.Errorf("se esperaba revelada_en en RFC 3339, se obtuvo %v", resp.ReveladaEn)
	}
	if resp.Promedio == nil || *resp.Promedio != (promedioResponse{Exacto: 5.33, CartaSugerida: 5}) {
		t.Errorf("se esperaba el promedio 5.33 / 5, se obtuvo %+v", resp.Promedio)
	}
}

func TestVotacionHandler_Consultar_Errores(t *testing.T) {
	casos := map[string]struct {
		id, consulta string
		previas      []string
		status       int
		campo        string
		mensaje      string
	}{
		"id no numérico":         {"abc", "", nil, nethttp.StatusBadRequest, "id", ""},
		"historia inexistente":   {"999", "?usuario_id=abc", nil, nethttp.StatusNotFound, "", "historia no encontrada"},
		"usuario_id no numérico": {"3", "?usuario_id=abc", []string{rondaAbierta}, nethttp.StatusBadRequest, "usuario_id", ""},
		"sin votación":           {"3", "?usuario_id=2", nil, nethttp.StatusNotFound, "", "la historia no tiene votación iniciada"},
	}
	for nombre, c := range casos {
		t.Run(nombre, func(t *testing.T) {
			handler := nuevoVotacionHandlerDePrueba(t, c.previas...)

			rec := consultarVotacionDe(handler, c.id, c.consulta)

			esperarError(t, rec, c.status, c.campo, c.mensaje)
		})
	}
}
