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

	"github.com/MoyaCarlos/seguimiento-medicion/internal/app"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

type usuarioDePrueba struct {
	ID     int64
	Nombre string
}

type pasosGestionUsuarios struct {
	db              *sql.DB
	api             *http.ServeMux
	tokens          map[string]string
	usuarios        map[string]usuarioDePrueba
	proyectos       map[string]int64
	usuarioActual   string
	contrasena      string
	respuesta       *httptest.ResponseRecorder
	respuestaError  respuestaError
	respuestaSesion respuestaSesion
	respuestaYo     respuestaYo
}

type respuestaSesion struct {
	Usuario struct {
		ID     int64  `json:"id"`
		Nombre string `json:"nombre"`
	} `json:"usuario"`
	Token string `json:"token"`
}

type respuestaYo struct {
	ID     int64  `json:"id"`
	Nombre string `json:"nombre"`
}

type respuestaError struct {
	Campo   string `json:"campo"`
	Mensaje string `json:"mensaje"`
}

type respuestaProyecto struct {
	ID int64 `json:"id"`
}

type respuestaIntegrante struct {
	ID     int64  `json:"id"`
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
}

func (p *pasosGestionUsuarios) preparar() error {
	db, err := repository.AbrirSQLite(":memory:")
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(1)
	if err := repository.AplicarMigraciones(context.Background(), db); err != nil {
		return err
	}
	p.db = db
	p.api = app.NuevoMux(db, app.Opciones{
		DemoHabilitado: func() bool { return false },
		Ahora: func() time.Time {
			return time.Date(2026, 11, 15, 12, 0, 0, 0, time.UTC)
		},
		Contrasenas: service.NewServicioContrasenas(1),
	})
	p.tokens = map[string]string{}
	p.usuarios = map[string]usuarioDePrueba{}
	p.proyectos = map[string]int64{}
	p.usuarioActual = ""
	p.contrasena = ""
	p.respuesta = nil
	p.respuestaError = respuestaError{}
	return nil
}

func (p *pasosGestionUsuarios) cerrar() {
	if p.db != nil {
		_ = p.db.Close()
		p.db = nil
	}
}

func (p *pasosGestionUsuarios) pedir(metodo, ruta, cuerpo, token string) error {
	req := httptest.NewRequest(metodo, ruta, strings.NewReader(cuerpo))
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	p.respuesta = httptest.NewRecorder()
	p.respuestaError = respuestaError{}
	p.api.ServeHTTP(p.respuesta, req)
	return nil
}

func (p *pasosGestionUsuarios) registrar(nombre, contrasena string) error {
	cuerpo, err := json.Marshal(map[string]string{"nombre": nombre, "contrasena": contrasena})
	if err != nil {
		return err
	}
	if err := p.pedir(http.MethodPost, "/auth/registro", string(cuerpo), ""); err != nil {
		return err
	}
	if p.respuesta.Code != http.StatusCreated {
		return fmt.Errorf("no se pudo registrar %q: status %d (%s)", nombre, p.respuesta.Code, p.respuesta.Body.String())
	}
	if err := json.Unmarshal(p.respuesta.Body.Bytes(), &p.respuestaSesion); err != nil {
		return err
	}
	p.tokens[nombre] = p.respuestaSesion.Token
	p.usuarios[nombre] = usuarioDePrueba{
		ID:     p.respuestaSesion.Usuario.ID,
		Nombre: p.respuestaSesion.Usuario.Nombre,
	}
	return nil
}

func (p *pasosGestionUsuarios) queTieneUnaCuenta(nombre, contrasena string) error {
	p.usuarioActual = nombre
	p.contrasena = contrasena
	return p.registrar(nombre, contrasena)
}

func (p *pasosGestionUsuarios) queExisteOtraCuenta(nombre, contrasena string) error {
	return p.registrar(nombre, contrasena)
}

func (p *pasosGestionUsuarios) queCreoProyecto(usuario, proyecto string) error {
	token, ok := p.tokens[usuario]
	if !ok {
		return fmt.Errorf("no hay sesión para %q", usuario)
	}
	cuerpo, err := json.Marshal(map[string]string{"nombre": proyecto, "descripcion": "Proyecto de prueba"})
	if err != nil {
		return err
	}
	if err := p.pedir(http.MethodPost, "/projects", string(cuerpo), token); err != nil {
		return err
	}
	if p.respuesta.Code != http.StatusCreated {
		return fmt.Errorf("no se pudo crear el proyecto: status %d (%s)", p.respuesta.Code, p.respuesta.Body.String())
	}
	var creado respuestaProyecto
	if err := json.Unmarshal(p.respuesta.Body.Bytes(), &creado); err != nil {
		return err
	}
	p.proyectos[proyecto] = creado.ID
	return nil
}

func (p *pasosGestionUsuarios) cambioMiNombre(nombre string) error {
	return p.cambiarNombre(p.tokens[p.usuarioActual], nombre)
}

func (p *pasosGestionUsuarios) intentoCambiarMiNombre(nombre string) error {
	return p.cambiarNombre(p.tokens[p.usuarioActual], nombre)
}

func (p *pasosGestionUsuarios) cambiarNombre(token, nombre string) error {
	cuerpo, err := json.Marshal(map[string]string{"nombre": nombre})
	if err != nil {
		return err
	}
	return p.pedir(http.MethodPut, "/auth/yo", string(cuerpo), token)
}

func (p *pasosGestionUsuarios) intentoCambiarSinSesion(nombre string) error {
	cuerpo, err := json.Marshal(map[string]string{"nombre": nombre})
	if err != nil {
		return err
	}
	return p.pedir(http.MethodPut, "/auth/yo", string(cuerpo), "")
}

func (p *pasosGestionUsuarios) respondeConNombre(status int, nombre string) error {
	if p.respuesta.Code != status {
		return fmt.Errorf("status = %d, se esperaba %d (%s)", p.respuesta.Code, status, p.respuesta.Body.String())
	}
	if err := json.Unmarshal(p.respuesta.Body.Bytes(), &p.respuestaYo); err != nil {
		return err
	}
	if p.respuestaYo.Nombre != nombre || p.respuestaYo.ID == 0 {
		return fmt.Errorf("respuesta = %+v, se esperaba id y nombre %q", p.respuestaYo, nombre)
	}
	return nil
}

func (p *pasosGestionUsuarios) equipoMuestra(proyecto, nombre string) error {
	id, ok := p.proyectos[proyecto]
	if !ok {
		return fmt.Errorf("no existe el proyecto %q", proyecto)
	}
	token, ok := p.tokens[p.usuarioActual]
	if !ok {
		return fmt.Errorf("no hay sesión para %q", p.usuarioActual)
	}
	if err := p.pedir(http.MethodGet, fmt.Sprintf("/projects/%d/members", id), "", token); err != nil {
		return err
	}
	if p.respuesta.Code != http.StatusOK {
		return fmt.Errorf("status = %d, se esperaba 200 (%s)", p.respuesta.Code, p.respuesta.Body.String())
	}
	var miembros []respuestaIntegrante
	if err := json.Unmarshal(p.respuesta.Body.Bytes(), &miembros); err != nil {
		return err
	}
	for _, miembro := range miembros {
		if miembro.Nombre != nombre {
			continue
		}
		if miembro.ID != p.usuarios[p.usuarioActual].ID {
			return fmt.Errorf("el integrante %q tiene id %d, se esperaba %d", nombre, miembro.ID, p.usuarios[p.usuarioActual].ID)
		}
		return nil
	}
	return fmt.Errorf("el equipo de %q no muestra a %q", proyecto, nombre)
}

func (p *pasosGestionUsuarios) puedoIngresar(nombre, contrasena string) error {
	cuerpo, err := json.Marshal(map[string]string{"nombre": nombre, "contrasena": contrasena})
	if err != nil {
		return err
	}
	if err := p.pedir(http.MethodPost, "/auth/ingreso", string(cuerpo), ""); err != nil {
		return err
	}
	if p.respuesta.Code != http.StatusOK {
		return fmt.Errorf("status = %d, se esperaba 200 (%s)", p.respuesta.Code, p.respuesta.Body.String())
	}
	var sesion respuestaSesion
	if err := json.Unmarshal(p.respuesta.Body.Bytes(), &sesion); err != nil {
		return err
	}
	if sesion.Token == "" {
		return fmt.Errorf("el ingreso no devolvió token")
	}
	return nil
}

func (p *pasosGestionUsuarios) ingresoRechazado(nombre string) error {
	cuerpo, err := json.Marshal(map[string]string{"nombre": nombre, "contrasena": p.contrasena})
	if err != nil {
		return err
	}
	if err := p.pedir(http.MethodPost, "/auth/ingreso", string(cuerpo), ""); err != nil {
		return err
	}
	if p.respuesta.Code != http.StatusUnauthorized {
		return fmt.Errorf("status = %d, se esperaba 401 (%s)", p.respuesta.Code, p.respuesta.Body.String())
	}
	var rechazo respuestaError
	if err := json.Unmarshal(p.respuesta.Body.Bytes(), &rechazo); err != nil {
		return err
	}
	if rechazo.Mensaje != "Usuario o contraseña incorrectos" {
		return fmt.Errorf("mensaje = %q, se esperaba el de credenciales inválidas", rechazo.Mensaje)
	}
	return nil
}

func (p *pasosGestionUsuarios) miNombreSigueSiendo(nombre string) error {
	token, ok := p.tokens[p.usuarioActual]
	if !ok {
		return fmt.Errorf("no hay sesión para %q", p.usuarioActual)
	}
	if err := p.pedir(http.MethodGet, "/auth/yo", "", token); err != nil {
		return err
	}
	if p.respuesta.Code != http.StatusOK {
		return fmt.Errorf("status = %d, se esperaba 200 (%s)", p.respuesta.Code, p.respuesta.Body.String())
	}
	var actual respuestaYo
	if err := json.Unmarshal(p.respuesta.Body.Bytes(), &actual); err != nil {
		return err
	}
	if actual.Nombre != nombre {
		return fmt.Errorf("el nombre es %q, se esperaba %q", actual.Nombre, nombre)
	}
	return nil
}

func (p *pasosGestionUsuarios) cambioRechazadoMensaje(status int, mensaje string) error {
	if p.respuesta.Code != status {
		return fmt.Errorf("status = %d, se esperaba %d (%s)", p.respuesta.Code, status, p.respuesta.Body.String())
	}
	if err := json.Unmarshal(p.respuesta.Body.Bytes(), &p.respuestaError); err != nil {
		return err
	}
	if p.respuestaError.Mensaje != mensaje {
		return fmt.Errorf("mensaje = %q, se esperaba %q", p.respuestaError.Mensaje, mensaje)
	}
	return nil
}

func (p *pasosGestionUsuarios) peticionRechazadaCampo(status int, campo string) error {
	if p.respuesta.Code != status {
		return fmt.Errorf("status = %d, se esperaba %d (%s)", p.respuesta.Code, status, p.respuesta.Body.String())
	}
	if err := json.Unmarshal(p.respuesta.Body.Bytes(), &p.respuestaError); err != nil {
		return err
	}
	if p.respuestaError.Campo != campo {
		return fmt.Errorf("campo = %q, se esperaba %q", p.respuestaError.Campo, campo)
	}
	return nil
}

func InitializeScenarioGestionUsuarios(ctx *godog.ScenarioContext) {
	p := &pasosGestionUsuarios{}
	ctx.Before(func(goctx context.Context, _ *godog.Scenario) (context.Context, error) {
		return goctx, p.preparar()
	})
	ctx.After(func(goctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		p.cerrar()
		return goctx, nil
	})

	ctx.Step(`^que "([^"]*)" tiene una cuenta con la contraseña "([^"]*)"$`, p.queTieneUnaCuenta)
	ctx.Step(`^que existe otra cuenta "([^"]*)" con la contraseña "([^"]*)"$`, p.queExisteOtraCuenta)
	ctx.Step(`^que "([^"]*)" creó el proyecto de prueba "([^"]*)"$`, p.queCreoProyecto)
	ctx.Step(`^cambio mi nombre de usuario a "([^"]*)"$`, p.cambioMiNombre)
	ctx.Step(`^intento cambiar mi nombre a "([^"]*)"$`, p.intentoCambiarMiNombre)
	ctx.Step(`^intento cambiar el nombre sin sesión a "([^"]*)"$`, p.intentoCambiarSinSesion)
	ctx.Step(`^el sistema responde (\d+) con el nombre "([^"]*)"$`, p.respondeConNombre)
	ctx.Step(`^el equipo de "([^"]*)" muestra a "([^"]*)"$`, p.equipoMuestra)
	ctx.Step(`^puedo ingresar con "([^"]*)" y la contraseña "([^"]*)"$`, p.puedoIngresar)
	ctx.Step(`^el ingreso con "([^"]*)" es rechazado$`, p.ingresoRechazado)
	ctx.Step(`^mi nombre sigue siendo "([^"]*)"$`, p.miNombreSigueSiendo)
	ctx.Step(`^el cambio se rechaza con (\d+) y el mensaje "([^"]*)"$`, p.cambioRechazadoMensaje)
	ctx.Step(`^la petición se rechaza con (\d+) y el mensaje "([^"]*)"$`, p.cambioRechazadoMensaje)
	ctx.Step(`^la petición se rechaza con (\d+) en el campo "([^"]*)"$`, p.peticionRechazadaCampo)
}
