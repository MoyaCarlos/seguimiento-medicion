package features

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/cucumber/godog"

	apihttp "github.com/MoyaCarlos/seguimiento-medicion/internal/http"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

const iteracionesDePruebaAutenticacion = 1

// pasosAutenticacion guarda el estado de los escenarios de HU-20. Arma su propia
// base en memoria: no depende del proyecto del Background de HU-01.
type pasosAutenticacion struct {
	db              *sql.DB
	sesiones        *repository.SQLiteSesionRepository
	handler         *apihttp.AuthHandler
	ahora           time.Time
	contrasena      string
	cuentasAntes    int
	respuesta       *httptest.ResponseRecorder
	respuestaSesion struct {
		Usuario struct {
			ID     int64  `json:"id"`
			Nombre string `json:"nombre"`
		} `json:"usuario"`
		Token   string `json:"token"`
		VenceEn string `json:"vence_en"`
	}
	respuestaError struct {
		Campo   string `json:"campo"`
		Mensaje string `json:"mensaje"`
	}
}

func (p *pasosAutenticacion) preparar() error {
	db, err := repository.AbrirSQLite(":memory:")
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(1)
	if err := repository.AplicarMigraciones(context.Background(), db); err != nil {
		return err
	}
	p.db = db
	p.ahora = time.Date(2026, 11, 4, 12, 0, 0, 0, time.UTC)
	p.sesiones = repository.NewSQLiteSesionRepository(db)
	emisor := service.NewEmisorSesiones(p.sesiones, func() time.Time { return p.ahora })
	credenciales := repository.NewSQLiteCredencialRepository(db)
	hasheador := service.NewServicioContrasenas(iteracionesDePruebaAutenticacion)
	p.handler = apihttp.NewAuthHandler(
		service.NewRegistrarUsuario(credenciales, hasheador, emisor),
		service.NewIngresar(credenciales, hasheador, emisor),
	)
	p.contrasena, p.cuentasAntes, p.respuesta = "", 0, nil
	p.respuestaSesion.Token = ""
	p.respuestaError.Campo, p.respuestaError.Mensaje = "", ""
	return nil
}

func (p *pasosAutenticacion) cerrar() {
	if p.db != nil {
		_ = p.db.Close()
		p.db = nil
	}
}

func (p *pasosAutenticacion) cuentas() (int, error) {
	var total int
	err := p.db.QueryRow(`SELECT COUNT(*) FROM usuarios`).Scan(&total)
	return total, err
}

func (p *pasosAutenticacion) registrar(nombre, contrasena string) (*httptest.ResponseRecorder, error) {
	cuerpo, err := json.Marshal(map[string]string{"nombre": nombre, "contrasena": contrasena})
	if err != nil {
		return nil, err
	}
	req := httptest.NewRequest(http.MethodPost, "/auth/registro", strings.NewReader(string(cuerpo)))
	rec := httptest.NewRecorder()
	p.handler.Registro(rec, req)
	return rec, nil
}

func (p *pasosAutenticacion) yaExisteLaCuenta(nombre string) error {
	rec, err := p.registrar(nombre, "clave-previa-1")
	if err != nil {
		return err
	}
	if rec.Code != http.StatusCreated {
		return fmt.Errorf("no se pudo crear la cuenta previa: status %d", rec.Code)
	}
	return nil
}

func (p *pasosAutenticacion) meRegistro(nombre, contrasena string) error {
	antes, err := p.cuentas()
	if err != nil {
		return err
	}
	p.cuentasAntes, p.contrasena = antes, contrasena
	if p.respuesta, err = p.registrar(nombre, contrasena); err != nil {
		return err
	}
	cuerpo := p.respuesta.Body.Bytes()
	if p.respuesta.Code == http.StatusCreated {
		return json.Unmarshal(cuerpo, &p.respuestaSesion)
	}
	return json.Unmarshal(cuerpo, &p.respuestaError)
}

func (p *pasosAutenticacion) cuentaCreada(nombre string) error {
	if p.respuesta.Code != http.StatusCreated {
		return fmt.Errorf("se esperaba 201, se obtuvo %d: %s", p.respuesta.Code, p.respuesta.Body.String())
	}
	var encontrada int
	err := p.db.QueryRow(`SELECT COUNT(*) FROM usuarios WHERE nombre = ?`, nombre).Scan(&encontrada)
	if err != nil {
		return err
	}
	if encontrada != 1 {
		return fmt.Errorf("se esperaba una cuenta %q, hay %d", nombre, encontrada)
	}
	return nil
}

func (p *pasosAutenticacion) sesionValida(nombre string) error {
	if p.respuestaSesion.Token == "" {
		return fmt.Errorf("la respuesta no trae token")
	}
	encontrada, err := p.sesiones.ObtenerPorHuella(context.Background(), service.HuellaDeToken(p.respuestaSesion.Token))
	if err != nil {
		return fmt.Errorf("el token no corresponde a una sesión: %w", err)
	}
	if encontrada.Usuario.Nombre != nombre {
		return fmt.Errorf("la sesión es de %q, se esperaba %q", encontrada.Usuario.Nombre, nombre)
	}
	if encontrada.CerradaEn != nil || !encontrada.VenceEn.After(p.ahora) {
		return fmt.Errorf("la sesión no está vigente")
	}
	return nil
}

func (p *pasosAutenticacion) respuestaSinContrasena() error {
	cuerpo := p.respuesta.Body.String()
	if strings.Contains(cuerpo, p.contrasena) {
		return fmt.Errorf("la respuesta contiene la contraseña")
	}
	return nil
}

func (p *pasosAutenticacion) rechazadoConMensaje(mensaje string) error {
	if p.respuesta.Code != http.StatusConflict || p.respuestaError.Mensaje != mensaje {
		return fmt.Errorf("se esperaba 409 con %q, se obtuvo %d con %q", mensaje, p.respuesta.Code, p.respuestaError.Mensaje)
	}
	return nil
}

func (p *pasosAutenticacion) rechazadoEnElCampo(campo string) error {
	if p.respuesta.Code != http.StatusBadRequest || p.respuestaError.Campo != campo {
		return fmt.Errorf("se esperaba 400 en el campo %q, se obtuvo %d en %q", campo, p.respuesta.Code, p.respuestaError.Campo)
	}
	return nil
}

func (p *pasosAutenticacion) sinCuentaNueva() error {
	ahora, err := p.cuentas()
	if err != nil {
		return err
	}
	if ahora != p.cuentasAntes {
		return fmt.Errorf("había %d cuentas y ahora hay %d", p.cuentasAntes, ahora)
	}
	return nil
}

func InitializeScenarioAutenticacion(ctx *godog.ScenarioContext) {
	p := &pasosAutenticacion{}
	ctx.Before(func(goctx context.Context, _ *godog.Scenario) (context.Context, error) {
		return goctx, p.preparar()
	})
	ctx.After(func(goctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		p.cerrar()
		return goctx, nil
	})

	ctx.Step(`^que ya existe la cuenta "([^"]*)"$`, p.yaExisteLaCuenta)
	ctx.Step(`^me registro con el nombre de usuario "([^"]*)" y la contraseña "([^"]*)"$`, p.meRegistro)
	ctx.Step(`^la cuenta "([^"]*)" queda creada$`, p.cuentaCreada)
	ctx.Step(`^la persona queda con una sesión válida a nombre de "([^"]*)"$`, p.sesionValida)
	ctx.Step(`^la respuesta del registro no contiene la contraseña$`, p.respuestaSinContrasena)
	ctx.Step(`^el registro de la cuenta se rechaza con el mensaje "([^"]*)"$`, p.rechazadoConMensaje)
	ctx.Step(`^el registro de la cuenta se rechaza indicando el campo "([^"]*)"$`, p.rechazadoEnElCampo)
	ctx.Step(`^no se crea ninguna cuenta nueva$`, p.sinCuentaNueva)
}
