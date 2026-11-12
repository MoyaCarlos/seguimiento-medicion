package features

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/cucumber/godog"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/demo"
	apphttp "github.com/MoyaCarlos/seguimiento-medicion/internal/http"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

var tablasDeDatosDePrueba = []string{"usuarios", "proyectos", "integrantes", "sprints", "historias",
	"historias_arrastradas", "registros_horas", "rondas_votacion", "votos", "defectos"}

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

func (p *pasosDatosPrueba) pedir(metodo, ruta string, cuerpo io.Reader) *httptest.ResponseRecorder {
	mux := http.NewServeMux()
	hasheador := service.NewServicioContrasenas(1)
	demoHandler := apphttp.NewDemoHandler(func() bool { return p.habilitado }, demo.NewCargador(p.db, time.Now, hasheador))
	for _, r := range apphttp.Rutas(apphttp.Handlers{Demo: demoHandler}) {
		if strings.Contains(r.Patron, "/demo/") {
			mux.HandleFunc(r.Patron, r.Handler)
		}
	}
	mux.HandleFunc("POST /auth/ingreso", p.autenticacion(hasheador).Ingreso)
	respuesta := httptest.NewRecorder()
	mux.ServeHTTP(respuesta, httptest.NewRequest(metodo, ruta, cuerpo))
	return respuesta
}

func (p *pasosDatosPrueba) autenticacion(hasheador *service.ServicioContrasenas) *apphttp.AuthHandler {
	sesiones := repository.NewSQLiteSesionRepository(p.db)
	emisor := service.NewEmisorSesiones(sesiones, time.Now)
	credenciales := repository.NewSQLiteCredencialRepository(p.db)
	obtener := service.NewObtenerSesion(sesiones, time.Now)
	return apphttp.NewAuthHandler(
		service.NewRegistrarUsuario(credenciales, hasheador, emisor),
		service.NewIngresar(credenciales, hasheador, emisor),
		obtener,
		service.NewCerrarSesion(obtener, sesiones, time.Now),
	)
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
	if rec := p.pedir("POST", "/demo/datos", nil); rec.Code != http.StatusCreated {
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
	p.respuesta = p.pedir("POST", "/demo/datos", nil)
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

func (p *pasosDatosPrueba) usuarioYaExiste(nombre string) error {
	_, err := p.db.Exec(`INSERT INTO usuarios (nombre, nombre_normalizado, creado_en) VALUES (?, ?, '2026-01-01T00:00:00Z')`,
		nombre, strings.ToLower(nombre))
	return err
}

func (p *pasosDatosPrueba) baseRechazaDefectosBaja() error {
	_, err := p.db.Exec(`CREATE TRIGGER falla_en_defecto_baja BEFORE INSERT ON defectos WHEN NEW.severidad = 'Baja'
		BEGIN SELECT RAISE(ABORT, 'fallo inyectado'); END`)
	return err
}

func (p *pasosDatosPrueba) ingresoConUsuarioDeLaDemo(nombre, contrasena string) error {
	cuerpo, err := json.Marshal(map[string]string{"nombre": nombre, "contrasena": contrasena})
	if err != nil {
		return err
	}
	p.respuesta = p.pedir("POST", "/auth/ingreso", bytes.NewReader(cuerpo))
	return nil
}

func (p *pasosDatosPrueba) seAbreLaSesionDe(nombre string) error {
	if p.respuesta.Code != http.StatusOK {
		return fmt.Errorf("estado %d, se esperaba 200: %s", p.respuesta.Code, p.respuesta.Body)
	}
	var sesion struct {
		Usuario struct {
			Nombre string `json:"nombre"`
		} `json:"usuario"`
		Token string `json:"token"`
	}
	if err := json.Unmarshal(p.respuesta.Body.Bytes(), &sesion); err != nil {
		return err
	}
	if sesion.Usuario.Nombre != nombre || sesion.Token == "" {
		return fmt.Errorf("sesión de %q con token %q, se esperaba la de %q con token", sesion.Usuario.Nombre, sesion.Token, nombre)
	}
	return nil
}

func (p *pasosDatosPrueba) ningunaContrasenaEnClaro() error {
	filas, err := p.db.Query(`SELECT nombre, COALESCE(contrasena_hash, '') FROM usuarios`)
	if err != nil {
		return err
	}
	defer filas.Close()
	usuarios := 0
	for filas.Next() {
		var nombre, hash string
		if err := filas.Scan(&nombre, &hash); err != nil {
			return err
		}
		usuarios++
		if hash == "" || strings.Contains(hash, "demo1234") {
			return fmt.Errorf("%s guarda %q: debe ser un hash, no la contraseña", nombre, hash)
		}
	}
	if usuarios != 4 {
		return fmt.Errorf("se esperaban 4 usuarios, hay %d", usuarios)
	}
	return filas.Err()
}

func (p *pasosDatosPrueba) noHayUsuariosBeta() error {
	var n int
	if err := p.db.QueryRow(`SELECT COUNT(*) FROM usuarios WHERE nombre LIKE 'Beta %'`).Scan(&n); err != nil {
		return err
	}
	if n != 0 {
		return fmt.Errorf("hay %d usuarios Beta", n)
	}
	return nil
}

func (p *pasosDatosPrueba) resumenInforma(usuarios, integrantes, sprints, historias, horas, rondas, votos, defectos int) error {
	var resumen demo.Resumen
	if err := json.Unmarshal(p.respuesta.Body.Bytes(), &resumen); err != nil {
		return err
	}
	esperado := demo.Resumen{ProyectoID: 25, Usuarios: usuarios, Integrantes: integrantes, Sprints: sprints,
		Historias: historias, RegistrosHoras: horas, Rondas: rondas, Votos: votos, Defectos: defectos}
	if resumen != esperado {
		return fmt.Errorf("resumen %+v, se esperaba %+v", resumen, esperado)
	}
	return nil
}

func (p *pasosDatosPrueba) proyectoTieneDefectos(total, resueltos, pendientes int) error {
	esperados := map[string]int{"": total, "Resuelto": resueltos, "Pendiente": pendientes}
	for estado, esperado := range esperados {
		consulta, args := `SELECT COUNT(*) FROM defectos WHERE proyecto_id = 25`, []any{}
		if estado != "" {
			consulta, args = consulta+` AND estado = ?`, []any{estado}
		}
		var n int
		if err := p.db.QueryRow(consulta, args...).Scan(&n); err != nil {
			return err
		}
		if n != esperado {
			return fmt.Errorf("defectos %q del proyecto 25: %d, se esperaban %d", estado, n, esperado)
		}
	}
	return nil
}

func (p *pasosDatosPrueba) esElUnicoScrumMaster(nombre string) error {
	var scrumMasters, siendoEl int
	if err := p.db.QueryRow(`SELECT COUNT(*) FROM integrantes WHERE proyecto_id = 25 AND rol = 'scrum_master'`).Scan(&scrumMasters); err != nil {
		return err
	}
	if err := p.db.QueryRow(`SELECT COUNT(*) FROM integrantes i JOIN usuarios u ON u.id = i.usuario_id
		WHERE i.proyecto_id = 25 AND i.rol = 'scrum_master' AND u.nombre_normalizado = ?`, nombre).Scan(&siendoEl); err != nil {
		return err
	}
	if scrumMasters != 1 || siendoEl != 1 {
		return fmt.Errorf("Scrum Masters del proyecto 25: %d (%s: %d), se esperaba uno solo y era %s", scrumMasters, nombre, siendoEl, nombre)
	}
	return nil
}

func (p *pasosDatosPrueba) respondeErrorInterno() error {
	if p.respuesta.Code != http.StatusInternalServerError {
		return fmt.Errorf("estado %d, se esperaba 500", p.respuesta.Code)
	}
	if strings.Contains(p.respuesta.Body.String(), "fallo inyectado") {
		return fmt.Errorf("el detalle del error llegó al cliente: %s", p.respuesta.Body)
	}
	return nil
}

func (p *pasosDatosPrueba) mensajeNoNombraUsuarios() error {
	if strings.Contains(p.respuesta.Body.String(), "ya existe el usuario") {
		return fmt.Errorf("el mensaje nombra a un usuario: %s", p.respuesta.Body)
	}
	return nil
}

// InitializeScenarioDatosPrueba registra los escenarios de los issues #144 y #235.
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
	ctx.Step(`^que el usuario "([^"]*)" ya existe$`, p.usuarioYaExiste)
	ctx.Step(`^que la base rechaza los defectos de severidad Baja$`, p.baseRechazaDefectosBaja)
	ctx.Step(`^ingreso con el usuario de la demo "([^"]*)" y la contraseña "([^"]*)"$`, p.ingresoConUsuarioDeLaDemo)
	ctx.Step(`^el sistema abre la sesión de "([^"]*)"$`, p.seAbreLaSesionDe)
	ctx.Step(`^ningún usuario de la demo guarda su contraseña en claro$`, p.ningunaContrasenaEnClaro)
	ctx.Step(`^no existen usuarios «Beta»$`, p.noHayUsuariosBeta)
	ctx.Step(`^el resumen informa (\d+) usuarios, (\d+) integrantes, (\d+) Sprints, (\d+) historias, (\d+) registros de horas, (\d+) rondas de votación, (\d+) votos y (\d+) defectos$`, p.resumenInforma)
	ctx.Step(`^el proyecto 25 tiene (\d+) defectos, (\d+) Resueltos y (\d+) Pendientes$`, p.proyectoTieneDefectos)
	ctx.Step(`^"([^"]*)" es el único Scrum Master del proyecto 25$`, p.esElUnicoScrumMaster)
	ctx.Step(`^el sistema responde 500 sin filtrar el detalle del error$`, p.respondeErrorInterno)
	ctx.Step(`^el mensaje no nombra a ningún usuario$`, p.mensajeNoNombraUsuarios)
}
