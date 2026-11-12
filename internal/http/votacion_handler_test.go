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
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
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
		service.NewConsultarVotacion(historias, proyectos, votaciones, domain.CoincidenciaExacta{}),
		service.NewRevelarVotacion(historias, proyectos, votaciones),
		service.NewIniciarVotacion(historias, proyectos, votaciones),
		service.NewRegistrarEstimacion(historias, proyectos, votaciones, repository.NewSQLiteSprintRepository(db)),
		service.NewAbrirNuevaRonda(historias, proyectos, votaciones),
		service.NewConsultarHistorialRondas(historias, proyectos, votaciones),
	)
}

func ejecutar(t *testing.T, db *sql.DB, sentencia string) {
	t.Helper()
	if _, err := db.Exec(sentencia); err != nil {
		t.Fatalf("no se pudo ejecutar %q: %v", sentencia, err)
	}
}

func votarEn(handler *VotacionHandler, id string, actor int64, cuerpo string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(nethttp.MethodPost, "/backlog/"+id+"/votacion/votos", bytes.NewBufferString(cuerpo))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", id)
	handler.Votar(rec, comoActor(req, actor))
	return rec
}

func consultarVotacionDe(handler *VotacionHandler, id string, actor int64, consulta string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(nethttp.MethodGet, "/backlog/"+id+"/votacion"+consulta, nil)
	req.SetPathValue("id", id)
	handler.Consultar(rec, comoActor(req, actor))
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

	rec := votarEn(handler, "3", 2, `{"valor":8}`)

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
		id      string
		actor   int64
		cuerpo  string
		previas []string
		status  int
		campo   string
		mensaje string
	}{
		"id no numérico":                     {"abc", 2, `{"valor":8}`, nil, nethttp.StatusBadRequest, "id", ""},
		"historia inexistente con JSON roto": {"999", 2, `{no-es-json`, nil, nethttp.StatusNotFound, "", "historia no encontrada"},
		"JSON inválido":                      {"3", 2, `{no-es-json`, []string{rondaAbierta}, nethttp.StatusBadRequest, "", "JSON inválido"},
		"sin sesión":                         {"3", 0, `{"valor":8}`, []string{rondaAbierta}, nethttp.StatusUnauthorized, "", "Tu sesión no es válida"},
		"no es integrante":                   {"3", 9, `{"valor":8}`, []string{rondaAbierta}, nethttp.StatusForbidden, "", "No sos integrante de este proyecto"},
		"sin votación":                       {"3", 2, `{"valor":8}`, nil, nethttp.StatusNotFound, "", "la historia no tiene votación iniciada"},
		"carta fuera de la baraja":           {"3", 2, `{"valor":4}`, []string{rondaAbierta}, nethttp.StatusBadRequest, "valor", ""},
		"ronda revelada":                     {"3", 2, `{"valor":8}`, []string{rondaRevelada}, nethttp.StatusBadRequest, "estado", ""},
	}
	for nombre, c := range casos {
		t.Run(nombre, func(t *testing.T) {
			handler := nuevoVotacionHandlerDePrueba(t, c.previas...)

			rec := votarEn(handler, c.id, c.actor, c.cuerpo)

			esperarError(t, rec, c.status, c.campo, c.mensaje)
		})
	}
}

func TestVotacionHandler_Votar_ValorNoEntero(t *testing.T) {
	casos := map[string]struct {
		actor  int64
		cuerpo string
	}{
		"decimal":        {2, `{"valor":2.5}`},
		"entre comillas": {2, `{"valor":"5"}`},
		"null":           {2, `{"valor":null}`},
		"sin el campo":   {2, `{}`},
	}
	for nombre, c := range casos {
		t.Run(nombre, func(t *testing.T) {
			handler := nuevoVotacionHandlerDePrueba(t, rondaAbierta)

			rec := votarEn(handler, "3", c.actor, c.cuerpo)

			esperarError(t, rec, nethttp.StatusBadRequest, "valor", "")
		})
	}
}

func TestVotacionHandler_Votar_PermisoAntesQueElValor(t *testing.T) {
	handler := nuevoVotacionHandlerDePrueba(t, rondaAbierta)

	rec := votarEn(handler, "3", 9, `{"valor":2.5}`)

	esperarError(t, rec, nethttp.StatusForbidden, "", "No sos integrante de este proyecto")
}

func TestVotacionHandler_Consultar_SoloElVotoPropioEsVisible(t *testing.T) {
	handler := nuevoVotacionHandlerDePrueba(t, rondaAbierta, votosDeCarlosYCandela)

	resp := leerVotacion(t, consultarVotacionDe(handler, "3", 2, ""))

	esperarValores(t, resp, -1, 8, -1)
	if !resp.Integrantes[0].YaVoto || resp.Integrantes[2].YaVoto {
		t.Errorf("Carlos votó (oculto) y Jimena no: %+v", resp.Integrantes)
	}
	if resp.ReveladaEn != nil || resp.Promedio != nil {
		t.Errorf("una ronda oculta no tiene revelada_en ni promedio: %+v", resp)
	}
}

// Escenario 7 de HU-21: el «mi voto» de la vista es siempre el de la sesión;
// un usuario_id en la query se ignora y no deja ver el voto oculto de otro.
func TestVotacionHandler_Consultar_IgnoraElUsuarioDeLaQuery(t *testing.T) {
	handler := nuevoVotacionHandlerDePrueba(t, rondaAbierta, votosDeCarlosYCandela)

	resp := leerVotacion(t, consultarVotacionDe(handler, "3", 2, "?usuario_id=1"))

	esperarValores(t, resp, -1, 8, -1)
}

// Escenario 7 de HU-21: el voto queda a nombre de la sesión aunque el cuerpo
// traiga el usuario_id de otra persona.
func TestVotacionHandler_Votar_IgnoraElUsuarioDelCuerpo(t *testing.T) {
	handler := nuevoVotacionHandlerDePrueba(t, rondaAbierta)

	if rec := votarEn(handler, "3", 2, `{"usuario_id":1,"valor":8}`); rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}

	esperarValores(t, leerVotacion(t, consultarVotacionDe(handler, "3", 2, "")), -1, 8, -1)
	esperarValores(t, leerVotacion(t, consultarVotacionDe(handler, "3", 1, "")), -1, -1, -1)
}

func TestVotacionHandler_Consultar_Revelada(t *testing.T) {
	handler := nuevoVotacionHandlerDePrueba(t, rondaRevelada,
		"INSERT INTO votos (ronda_id, usuario_id, valor) VALUES (1, 1, 3), (1, 2, 8), (1, 3, 5)")

	resp := leerVotacion(t, consultarVotacionDe(handler, "3", 3, ""))

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
		id       string
		actor    int64
		consulta string
		previas  []string
		status   int
		campo    string
		mensaje  string
	}{
		"id no numérico":       {"abc", 2, "", nil, nethttp.StatusBadRequest, "id", ""},
		"historia inexistente": {"999", 2, "", nil, nethttp.StatusNotFound, "", "historia no encontrada"},
		"sin sesión":           {"3", 0, "", []string{rondaAbierta}, nethttp.StatusUnauthorized, "", "Tu sesión no es válida"},
		"sin votación":         {"3", 2, "", nil, nethttp.StatusNotFound, "", "la historia no tiene votación iniciada"},
	}
	for nombre, c := range casos {
		t.Run(nombre, func(t *testing.T) {
			handler := nuevoVotacionHandlerDePrueba(t, c.previas...)

			rec := consultarVotacionDe(handler, c.id, c.actor, c.consulta)

			esperarError(t, rec, c.status, c.campo, c.mensaje)
		})
	}
}

func revelarEn(handler *VotacionHandler, id string, actor int64, cuerpo string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(nethttp.MethodPost, "/backlog/"+id+"/votacion/revelar", bytes.NewBufferString(cuerpo))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", id)
	handler.Revelar(rec, comoActor(req, actor))
	return rec
}

const votosDeTodos = "INSERT INTO votos (ronda_id, usuario_id, valor) VALUES (1, 1, 3), (1, 2, 8), (1, 3, 5)"

func TestVotacionHandler_Revelar_Exitoso(t *testing.T) {
	handler := nuevoVotacionHandlerDePrueba(t, rondaAbierta, votosDeTodos)

	resp := leerVotacion(t, revelarEn(handler, "3", 1, `{}`))

	esperarValores(t, resp, 3, 8, 5)
	if resp.ReveladaEn == nil {
		t.Fatal("se esperaba revelada_en")
	}
	if _, err := time.Parse(time.RFC3339, *resp.ReveladaEn); err != nil {
		t.Errorf("revelada_en debe ser RFC 3339: %v", err)
	}
	if resp.Promedio == nil || *resp.Promedio != (promedioResponse{Exacto: 5.33, CartaSugerida: 5}) {
		t.Errorf("se esperaba el promedio 5.33 / 5, se obtuvo %+v", resp.Promedio)
	}
}

func TestVotacionHandler_Revelar_Errores(t *testing.T) {
	casos := map[string]struct {
		id      string
		actor   int64
		cuerpo  string
		previas []string
		status  int
		campo   string
		mensaje string
	}{
		"historia inexistente con JSON roto": {"999", 1, `{no-es-json`, nil, nethttp.StatusNotFound, "", "historia no encontrada"},
		"sin sesión":                         {"3", 0, `{}`, []string{rondaAbierta, votosDeTodos}, nethttp.StatusUnauthorized, "", "Tu sesión no es válida"},
		"no es Scrum Master":                 {"3", 2, `{}`, []string{rondaAbierta, votosDeTodos}, nethttp.StatusForbidden, "", "Solo el Scrum Master puede realizar esta acción"},
		"sin votación":                       {"3", 1, `{}`, nil, nethttp.StatusNotFound, "", "la historia no tiene votación iniciada"},
		"falta votar Jimena":                 {"3", 1, `{}`, []string{rondaAbierta, votosDeCarlosYCandela}, nethttp.StatusBadRequest, "votos", ""},
		"ya revelada":                        {"3", 1, `{}`, []string{rondaRevelada, votosDeTodos}, nethttp.StatusBadRequest, "estado", ""},
	}
	for nombre, c := range casos {
		t.Run(nombre, func(t *testing.T) {
			handler := nuevoVotacionHandlerDePrueba(t, c.previas...)

			rec := revelarEn(handler, c.id, c.actor, c.cuerpo)

			esperarError(t, rec, c.status, c.campo, c.mensaje)
		})
	}
}

// historiaCompletada pasa la historia 3 a Completada dentro de un Sprint: el
// esquema exige Sprint para todo estado que no sea Nueva.
const historiaCompletada = `
INSERT INTO sprints (id, proyecto_id, sprint_goal, fecha_inicio, fecha_fin, estado) VALUES (10, 1, 'g', '2026-10-01', '2026-10-15', 'Activo');
UPDATE historias SET estado = 'Completada', sprint_id = 10, fecha_finalizacion = '2026-10-10' WHERE id = 3;
`

func iniciarEn(handler *VotacionHandler, id string, actor int64, cuerpo string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(nethttp.MethodPost, "/backlog/"+id+"/votacion", bytes.NewBufferString(cuerpo))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", id)
	handler.Iniciar(rec, comoActor(req, actor))
	return rec
}

func TestVotacionHandler_Iniciar_Exitoso(t *testing.T) {
	handler := nuevoVotacionHandlerDePrueba(t)

	rec := iniciarEn(handler, "3", 1, `{}`)

	if rec.Code != nethttp.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	var resp votacionResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta JSON inválida: %v", err)
	}
	if resp.HistoriaID != 3 || resp.Ronda != 1 || resp.ReveladaEn != nil || resp.Promedio != nil {
		t.Errorf("se esperaba la ronda 1 de la historia 3, oculta y sin promedio, se obtuvo %+v", resp)
	}
	esperarValores(t, resp, -1, -1, -1)
	for _, i := range resp.Integrantes {
		if i.YaVoto {
			t.Errorf("nadie votó todavía: %+v", resp.Integrantes)
		}
	}
}

func TestVotacionHandler_Iniciar_Errores(t *testing.T) {
	casos := map[string]struct {
		id      string
		actor   int64
		cuerpo  string
		previas []string
		status  int
		campo   string
		mensaje string
	}{
		"historia inexistente con cuerpo vacío": {"999", 1, `{}`, nil, nethttp.StatusNotFound, "", "historia no encontrada"},
		"sin sesión":                            {"3", 0, `{}`, nil, nethttp.StatusUnauthorized, "", "Tu sesión no es válida"},
		"no es Scrum Master":                    {"3", 2, `{}`, nil, nethttp.StatusForbidden, "", "Solo el Scrum Master puede realizar esta acción"},
		"historia Completada":                   {"3", 1, `{}`, []string{historiaCompletada}, nethttp.StatusBadRequest, "estado", ""},
		"ya tiene votación":                     {"3", 1, `{}`, []string{rondaAbierta}, nethttp.StatusConflict, "", "la historia ya tiene una votación"},
	}
	for nombre, c := range casos {
		t.Run(nombre, func(t *testing.T) {
			handler := nuevoVotacionHandlerDePrueba(t, c.previas...)

			rec := iniciarEn(handler, c.id, c.actor, c.cuerpo)

			esperarError(t, rec, c.status, c.campo, c.mensaje)
		})
	}
}

func TestVotacionHandler_Consultar_Diferencias(t *testing.T) {
	casos := map[string]struct {
		previas  []string
		esperado string
	}{
		"ronda oculta": {[]string{rondaAbierta, votosDeCarlosYCandela}, `"diferencias":null`},
		"revelada con diferencias": {[]string{rondaRevelada, votosDeCarlosYCandela},
			`"diferencias":{"hay_diferencias":true,"minimo":3,"maximo":8}`},
		"revelada sin diferencias": {[]string{rondaRevelada, "INSERT INTO votos (ronda_id, usuario_id, valor) VALUES (1, 1, 5), (1, 2, 5)"},
			`"diferencias":{"hay_diferencias":false,"minimo":null,"maximo":null}`},
	}
	for nombre, c := range casos {
		t.Run(nombre, func(t *testing.T) {
			handler := nuevoVotacionHandlerDePrueba(t, c.previas...)

			rec := consultarVotacionDe(handler, "3", 1, "")

			if rec.Code != nethttp.StatusOK {
				t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
			}
			if !strings.Contains(rec.Body.String(), c.esperado) {
				t.Errorf("se esperaba %s en %s", c.esperado, rec.Body.String())
			}
		})
	}
}

func registrarEstimacionEn(handler *VotacionHandler, id string, actor int64, cuerpo string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(nethttp.MethodPost, "/backlog/"+id+"/votacion/estimacion", bytes.NewBufferString(cuerpo))
	req.SetPathValue("id", id)
	handler.RegistrarEstimacion(rec, comoActor(req, actor))
	return rec
}

func TestVotacionHandler_RegistrarEstimacion_Exitoso(t *testing.T) {
	handler := nuevoVotacionHandlerDePrueba(t, rondaRevelada, votosDeCarlosYCandela)

	rec := registrarEstimacionEn(handler, "3", 1, `{"estimacion_sp":5}`)

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"estimacion_sp":5`) {
		t.Errorf("se esperaba la historia con 5 Story Points: %s", rec.Body.String())
	}
}

const sprintActivo = "INSERT INTO sprints (id, proyecto_id, estado, sprint_goal, fecha_inicio, fecha_fin) VALUES (1, 1, 'Activo', 'Objetivo', '2026-10-01', '2026-10-15')"

func TestVotacionHandler_RegistrarEstimacion_Congelada(t *testing.T) {
	handler := nuevoVotacionHandlerDePrueba(t, rondaRevelada, sprintActivo, "UPDATE historias SET sprint_id = 1, estimacion_sp = 3 WHERE id = 3")

	rec := registrarEstimacionEn(handler, "3", 1, `{"estimacion_sp":5}`)

	esperarError(t, rec, nethttp.StatusBadRequest, "estimacion_sp", "los Story Points ya cargados no se pueden cambiar en un Sprint Activo o Finalizado")
}

func TestVotacionHandler_RegistrarEstimacion_ValorNoEntero(t *testing.T) {
	casos := map[string]struct {
		actor  int64
		cuerpo string
	}{
		"decimal":        {1, `{"estimacion_sp":2.5}`},
		"entre comillas": {1, `{"estimacion_sp":"5"}`},
		"null":           {1, `{"estimacion_sp":null}`},
		"sin el campo":   {1, `{}`},
	}
	for nombre, c := range casos {
		t.Run(nombre, func(t *testing.T) {
			handler := nuevoVotacionHandlerDePrueba(t, rondaRevelada)

			rec := registrarEstimacionEn(handler, "3", c.actor, c.cuerpo)

			esperarError(t, rec, nethttp.StatusBadRequest, "estimacion_sp", "")
		})
	}
}

func TestVotacionHandler_RegistrarEstimacion_PermisoAntesQueElValor(t *testing.T) {
	handler := nuevoVotacionHandlerDePrueba(t, rondaRevelada)

	rec := registrarEstimacionEn(handler, "3", 2, `{"estimacion_sp":2.5}`)

	esperarError(t, rec, nethttp.StatusForbidden, "", "Solo el Scrum Master puede realizar esta acción")
}

func TestVotacionHandler_RegistrarEstimacion_Errores(t *testing.T) {
	casos := map[string]struct {
		id      string
		actor   int64
		cuerpo  string
		previas []string
		status  int
		campo   string
		mensaje string
	}{
		"id no numérico":                        {"abc", 1, `{"estimacion_sp":5}`, nil, nethttp.StatusBadRequest, "id", ""},
		"historia inexistente con JSON roto":    {"999", 1, `{no-es-json`, nil, nethttp.StatusNotFound, "", "historia no encontrada"},
		"historia inexistente con cuerpo vacío": {"999", 1, `{}`, nil, nethttp.StatusNotFound, "", "historia no encontrada"},
		"historia inexistente con pedido bueno": {"999", 1, `{"estimacion_sp":5}`, nil, nethttp.StatusNotFound, "", "historia no encontrada"},
		"JSON inválido":                         {"3", 1, `{no-es-json`, []string{rondaRevelada}, nethttp.StatusBadRequest, "", "JSON inválido"},
		"sin sesión":                            {"3", 0, `{"estimacion_sp":5}`, []string{rondaRevelada}, nethttp.StatusUnauthorized, "", "Tu sesión no es válida"},
		"carta fuera de la baraja":              {"3", 1, `{"estimacion_sp":4}`, []string{rondaRevelada}, nethttp.StatusBadRequest, "estimacion_sp", ""},
		"ronda oculta":                          {"3", 1, `{"estimacion_sp":5}`, []string{rondaAbierta}, nethttp.StatusBadRequest, "estado", "la ronda actual todavía no fue revelada"},
		"Product Builder":                       {"3", 2, `{"estimacion_sp":5}`, []string{rondaRevelada}, nethttp.StatusForbidden, "", "Solo el Scrum Master puede realizar esta acción"},
		"sin votación":                          {"3", 1, `{"estimacion_sp":5}`, nil, nethttp.StatusNotFound, "", "la historia no tiene votación iniciada"},
	}
	for nombre, c := range casos {
		t.Run(nombre, func(t *testing.T) {
			handler := nuevoVotacionHandlerDePrueba(t, c.previas...)

			rec := registrarEstimacionEn(handler, c.id, c.actor, c.cuerpo)

			esperarError(t, rec, c.status, c.campo, c.mensaje)
		})
	}
}

func abrirRondaEn(handler *VotacionHandler, id string, actor int64, cuerpo string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(nethttp.MethodPost, "/backlog/"+id+"/votacion/rondas", bytes.NewBufferString(cuerpo))
	req.SetPathValue("id", id)
	handler.AbrirRonda(rec, comoActor(req, actor))
	return rec
}

func TestVotacionHandler_AbrirRonda_Exitoso(t *testing.T) {
	handler := nuevoVotacionHandlerDePrueba(t, rondaRevelada, votosDeCarlosYCandela)

	rec := abrirRondaEn(handler, "3", 1, `{}`)

	if rec.Code != nethttp.StatusCreated {
		t.Fatalf("se esperaba 201, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	var resp votacionResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("respuesta JSON inválida: %v", err)
	}
	if resp.Ronda != 2 || resp.ReveladaEn != nil || resp.Diferencias != nil || resp.Promedio != nil {
		t.Errorf("se esperaba la ronda 2 oculta, sin diferencias ni promedio, se obtuvo %+v", resp)
	}
	for _, i := range resp.Integrantes {
		if i.YaVoto {
			t.Errorf("nadie votó todavía en la ronda 2: %+v", i)
		}
	}
	if !strings.Contains(rec.Body.String(), `"diferencias":null`) {
		t.Errorf("diferencias va en null en una ronda oculta: %s", rec.Body.String())
	}
}

func TestVotacionHandler_AbrirRonda_HistoriaCompletada(t *testing.T) {
	handler := nuevoVotacionHandlerDePrueba(t, rondaRevelada, sprintActivo,
		"UPDATE historias SET estado = 'Completada', sprint_id = 1, fecha_finalizacion = '2026-10-10' WHERE id = 3")

	rec := abrirRondaEn(handler, "3", 1, `{}`)

	esperarError(t, rec, nethttp.StatusBadRequest, "estado", "no se puede abrir una ronda en una historia completada")
}

func TestVotacionHandler_AbrirRonda_Errores(t *testing.T) {
	casos := map[string]struct {
		id      string
		actor   int64
		cuerpo  string
		previas []string
		status  int
		campo   string
		mensaje string
	}{
		"id no numérico":                        {"abc", 1, `{}`, nil, nethttp.StatusBadRequest, "id", ""},
		"historia inexistente con cuerpo vacío": {"999", 1, `{}`, nil, nethttp.StatusNotFound, "", "historia no encontrada"},
		"sin sesión":                            {"3", 0, `{}`, []string{rondaRevelada}, nethttp.StatusUnauthorized, "", "Tu sesión no es válida"},
		"ronda oculta":                          {"3", 1, `{}`, []string{rondaAbierta}, nethttp.StatusBadRequest, "estado", "la ronda actual todavía no fue revelada"},
		"Product Builder":                       {"3", 2, `{}`, []string{rondaRevelada}, nethttp.StatusForbidden, "", "Solo el Scrum Master puede realizar esta acción"},
		"sin votación":                          {"3", 1, `{}`, nil, nethttp.StatusNotFound, "", "la historia no tiene votación iniciada"},
	}
	for nombre, c := range casos {
		t.Run(nombre, func(t *testing.T) {
			handler := nuevoVotacionHandlerDePrueba(t, c.previas...)

			rec := abrirRondaEn(handler, c.id, c.actor, c.cuerpo)

			esperarError(t, rec, c.status, c.campo, c.mensaje)
		})
	}
}

func historialDe(handler *VotacionHandler, id string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(nethttp.MethodGet, "/backlog/"+id+"/votacion/rondas", nil)
	req.SetPathValue("id", id)
	handler.Historial(rec, req)
	return rec
}

func TestVotacionHandler_Historial_SoloRondasReveladas(t *testing.T) {
	handler := nuevoVotacionHandlerDePrueba(t, rondaRevelada, votosDeCarlosYCandela,
		"INSERT INTO rondas_votacion (id, historia_id, numero, revelada_en) VALUES (2, 3, 2, NULL)",
		"INSERT INTO votos (ronda_id, usuario_id, valor) VALUES (2, 1, 13)")

	rec := historialDe(handler, "3")

	esperado := `{"historia_id":3,"rondas":[{"ronda":1,"revelada_en":"2026-10-27T15:04:05Z","votos":[` +
		`{"usuario_id":1,"nombre":"Carlos","rol":"scrum_master","valor":3},` +
		`{"usuario_id":2,"nombre":"Candela","rol":"product_builder","valor":8}]}]}`
	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	if strings.TrimSpace(rec.Body.String()) != esperado {
		t.Errorf("se esperaba %s, se obtuvo %s", esperado, rec.Body.String())
	}
}

func TestVotacionHandler_Historial_SinRondasReveladas(t *testing.T) {
	handler := nuevoVotacionHandlerDePrueba(t, rondaAbierta)

	rec := historialDe(handler, "3")

	if rec.Code != nethttp.StatusOK {
		t.Fatalf("se esperaba 200, se obtuvo %d (%s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"rondas":[]`) {
		t.Errorf("rondas va como lista vacía, no null: %s", rec.Body.String())
	}
}

func TestVotacionHandler_Historial_Errores(t *testing.T) {
	handler := nuevoVotacionHandlerDePrueba(t)

	esperarError(t, historialDe(handler, "abc"), nethttp.StatusBadRequest, "id", "")
	esperarError(t, historialDe(handler, "999"), nethttp.StatusNotFound, "", "historia no encontrada")
}

func TestEscribirError_RondaYaAbierta(t *testing.T) {
	rec := httptest.NewRecorder()

	escribirError(rec, domain.ErrRondaYaAbierta)

	esperarError(t, rec, nethttp.StatusConflict, "", domain.ErrRondaYaAbierta.Error())
}

func TestVotacionHandler_Registrar_AtiendeCadaRuta(t *testing.T) {
	casos := []struct {
		metodo, ruta string
		actor        int64
		cuerpo       string
		status       int
		esperado     string
	}{
		{nethttp.MethodGet, "/backlog/3/votacion", 1, "", nethttp.StatusOK, `"diferencias":{`},
		{nethttp.MethodPost, "/backlog/3/votacion", 1, `{}`, nethttp.StatusConflict, "la historia ya tiene una votación"},
		{nethttp.MethodPost, "/backlog/3/votacion/votos", 2, `{"valor":8}`, nethttp.StatusBadRequest, "la ronda ya fue revelada"},
		{nethttp.MethodPost, "/backlog/3/votacion/revelar", 1, `{}`, nethttp.StatusBadRequest, "la ronda ya fue revelada"},
		{nethttp.MethodPost, "/backlog/3/votacion/estimacion", 1, `{"estimacion_sp":5}`, nethttp.StatusOK, `"estimacion_sp":5`},
		{nethttp.MethodPost, "/backlog/3/votacion/rondas", 1, `{}`, nethttp.StatusCreated, `"ronda":2`},
		{nethttp.MethodGet, "/backlog/3/votacion/rondas", 1, "", nethttp.StatusOK, `"rondas":[{"ronda":1`},
	}
	for _, c := range casos {
		t.Run(c.metodo+" "+c.ruta, func(t *testing.T) {
			mux := nethttp.NewServeMux()
			registrarRutasDe(mux, Handlers{Votacion: nuevoVotacionHandlerDePrueba(t, rondaRevelada, votosDeCarlosYCandela)}, "/backlog/{id}/votacion")
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, comoActor(httptest.NewRequest(c.metodo, c.ruta, strings.NewReader(c.cuerpo)), c.actor))

			if rec.Code != c.status {
				t.Fatalf("se esperaba %d, se obtuvo %d (%s)", c.status, rec.Code, rec.Body.String())
			}
			if !strings.Contains(rec.Body.String(), c.esperado) {
				t.Errorf("se esperaba %s en %s", c.esperado, rec.Body.String())
			}
		})
	}
}

func TestVotacionHandler_RegistrarEstimacion_CargaConSprintFinalizado(t *testing.T) {
	sprintFinalizado := strings.Replace(sprintActivo, "'Activo'", "'Finalizado'", 1)
	handler := nuevoVotacionHandlerDePrueba(t, rondaRevelada, sprintFinalizado, "UPDATE historias SET sprint_id = 1, estimacion_sp = NULL WHERE id = 3")

	rec := registrarEstimacionEn(handler, "3", 1, `{"estimacion_sp":5}`)

	esperarError(t, rec, nethttp.StatusBadRequest, "estimacion_sp", "no se pueden cargar Story Points en un Sprint Finalizado")
}
