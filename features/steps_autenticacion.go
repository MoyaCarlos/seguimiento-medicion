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
	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
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
	api             *http.ServeMux
	ahora           time.Time
	tokenSegunda    string
	contrasena      string
	cuentasAntes    int
	sesionesAntes   int
	respuesta       *httptest.ResponseRecorder
	respuestaSesion struct {
		Usuario struct {
			ID     int64  `json:"id"`
			Nombre string `json:"nombre"`
		} `json:"usuario"`
		Token   string `json:"token"`
		VenceEn string `json:"vence_en"`
	}
	respuestaYo struct {
		ID     int64  `json:"id"`
		Nombre string `json:"nombre"`
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
	reloj := func() time.Time { return p.ahora }
	obtener := service.NewObtenerSesion(p.sesiones, reloj)
	p.handler = apihttp.NewAuthHandler(
		service.NewRegistrarUsuario(credenciales, hasheador, emisor),
		service.NewIngresar(credenciales, hasheador, emisor),
		obtener,
		service.NewCerrarSesion(obtener, p.sesiones, reloj),
	)
	p.api = app.NuevoMux(db, app.Opciones{
		DemoHabilitado: func() bool { return false },
		Ahora:          reloj,
		Contrasenas:    hasheador,
	})
	p.contrasena, p.cuentasAntes, p.sesionesAntes, p.respuesta = "", 0, 0, nil
	p.respuestaSesion.Token, p.tokenSegunda = "", ""
	p.respuestaYo.ID, p.respuestaYo.Nombre = 0, ""
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

func (p *pasosAutenticacion) sesionesAbiertas() (int, error) {
	var total int
	err := p.db.QueryRow(`SELECT COUNT(*) FROM sesiones`).Scan(&total)
	return total, err
}

// existeUsuarioPrevio crea un usuario por el camino de HU-04: solo con nombre, sin credenciales.
func (p *pasosAutenticacion) existeUsuarioPrevio(nombre string) error {
	_, err := repository.NewSQLiteUserRepository(p.db).Guardar(context.Background(), domain.User{
		Nombre:            nombre,
		NombreNormalizado: domain.NormalizarNombre(nombre),
	})
	return err
}

func (p *pasosAutenticacion) ingreso(nombre, contrasena string) error {
	var err error
	if p.sesionesAntes, err = p.sesionesAbiertas(); err != nil {
		return err
	}
	p.contrasena = contrasena
	cuerpo, err := json.Marshal(map[string]string{"nombre": nombre, "contrasena": contrasena})
	if err != nil {
		return err
	}
	req := httptest.NewRequest(http.MethodPost, "/auth/ingreso", strings.NewReader(string(cuerpo)))
	p.respuesta = httptest.NewRecorder()
	p.handler.Ingreso(p.respuesta, req)
	if p.respuesta.Code == http.StatusOK {
		return json.Unmarshal(p.respuesta.Body.Bytes(), &p.respuestaSesion)
	}
	return json.Unmarshal(p.respuesta.Body.Bytes(), &p.respuestaError)
}

func (p *pasosAutenticacion) ingresoRechazado(mensaje string) error {
	if p.respuesta.Code != http.StatusUnauthorized || p.respuestaError.Mensaje != mensaje {
		return fmt.Errorf("se esperaba 401 con %q, se obtuvo %d con %q", mensaje, p.respuesta.Code, p.respuestaError.Mensaje)
	}
	return nil
}

func (p *pasosAutenticacion) sinSesionNueva() error {
	ahora, err := p.sesionesAbiertas()
	if err != nil {
		return err
	}
	if ahora != p.sesionesAntes {
		return fmt.Errorf("había %d sesiones y ahora hay %d", p.sesionesAntes, ahora)
	}
	return nil
}

// inicioSesion registra la cuenta y se queda con el token de la sesión que abre.
func (p *pasosAutenticacion) inicioSesion(nombre string) error {
	p.contrasena = "clave-previa-1"
	rec, err := p.registrar(nombre, p.contrasena)
	if err != nil {
		return err
	}
	if rec.Code != http.StatusCreated {
		return fmt.Errorf("no se pudo iniciar la sesión: status %d", rec.Code)
	}
	return json.Unmarshal(rec.Body.Bytes(), &p.respuestaSesion)
}

// pasanDias adelanta el reloj de los pasos; 7 días es lo que dura una sesión.
func (p *pasosAutenticacion) pasanDias(dias int) error {
	p.ahora = p.ahora.Add(time.Duration(dias) * 24 * time.Hour)
	return nil
}

func (p *pasosAutenticacion) consultarQuienSoy(autorizacion string) error {
	req := httptest.NewRequest(http.MethodGet, "/auth/yo", nil)
	if autorizacion != "" {
		req.Header.Set("Authorization", autorizacion)
	}
	p.respuesta = httptest.NewRecorder()
	p.handler.Yo(p.respuesta, req)
	if p.respuesta.Code == http.StatusOK {
		return json.Unmarshal(p.respuesta.Body.Bytes(), &p.respuestaYo)
	}
	return json.Unmarshal(p.respuesta.Body.Bytes(), &p.respuestaError)
}

func (p *pasosAutenticacion) consultoQuienSoyConMiToken() error {
	return p.consultarQuienSoy("Bearer " + p.respuestaSesion.Token)
}

func (p *pasosAutenticacion) consultoQuienSoyConElToken(token string) error {
	return p.consultarQuienSoy("Bearer " + token)
}

func (p *pasosAutenticacion) consultoQuienSoySinToken() error {
	return p.consultarQuienSoy("")
}

func (p *pasosAutenticacion) respondeQueSoy(nombre string) error {
	if p.respuesta.Code != http.StatusOK || p.respuestaYo.Nombre != nombre || p.respuestaYo.ID == 0 {
		return fmt.Errorf("se esperaba 200 con %q, se obtuvo %d: %s", nombre, p.respuesta.Code, p.respuesta.Body.String())
	}
	return nil
}

func (p *pasosAutenticacion) quienSoyRechazado(mensaje string) error {
	if p.respuesta.Code != http.StatusUnauthorized || p.respuestaError.Mensaje != mensaje {
		return fmt.Errorf("se esperaba 401 con %q, se obtuvo %d con %q", mensaje, p.respuesta.Code, p.respuestaError.Mensaje)
	}
	return nil
}

// abreUnaSegundaSesion ingresa de nuevo con la misma cuenta y guarda ese token
// aparte, sin pisar el de la primera sesión.
func (p *pasosAutenticacion) abreUnaSegundaSesion(nombre string) error {
	cuerpo, err := json.Marshal(map[string]string{"nombre": nombre, "contrasena": p.contrasena})
	if err != nil {
		return err
	}
	req := httptest.NewRequest(http.MethodPost, "/auth/ingreso", strings.NewReader(string(cuerpo)))
	rec := httptest.NewRecorder()
	p.handler.Ingreso(rec, req)
	if rec.Code != http.StatusOK {
		return fmt.Errorf("no se pudo abrir la segunda sesión: status %d", rec.Code)
	}
	var segunda struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &segunda); err != nil {
		return err
	}
	p.tokenSegunda = segunda.Token
	return nil
}

func (p *pasosAutenticacion) cierroMiSesion() error {
	req := httptest.NewRequest(http.MethodPost, "/auth/salir", nil)
	req.Header.Set("Authorization", "Bearer "+p.respuestaSesion.Token)
	p.respuesta = httptest.NewRecorder()
	p.handler.Salir(p.respuesta, req)
	return nil
}

func (p *pasosAutenticacion) cierreSinContenido() error {
	if p.respuesta.Code != http.StatusNoContent || p.respuesta.Body.Len() != 0 {
		return fmt.Errorf("se esperaba 204 sin cuerpo, se obtuvo %d con %q", p.respuesta.Code, p.respuesta.Body.String())
	}
	return nil
}

func (p *pasosAutenticacion) consultoQuienSoyConElTokenDeMiSegundaSesion() error {
	return p.consultarQuienSoy("Bearer " + p.tokenSegunda)
}

func (p *pasosAutenticacion) consultoLaListaDeProyectosSinToken() error {
	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	p.respuesta = httptest.NewRecorder()
	p.api.ServeHTTP(p.respuesta, req)
	return nil
}

func (p *pasosAutenticacion) listaDeProyectosRechazada(mensaje string) error {
	if p.respuesta.Code != http.StatusUnauthorized {
		return fmt.Errorf("se esperaba 401, se obtuvo %d: %s", p.respuesta.Code, p.respuesta.Body.String())
	}
	if err := json.Unmarshal(p.respuesta.Body.Bytes(), &p.respuestaError); err != nil {
		return err
	}
	if p.respuestaError.Mensaje != mensaje {
		return fmt.Errorf("se esperaba el mensaje %q, se obtuvo %q", mensaje, p.respuestaError.Mensaje)
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
	ctx.Step(`^que existe el usuario previo "([^"]*)" creado solo con nombre$`, p.existeUsuarioPrevio)
	ctx.Step(`^ingreso con el nombre de usuario "([^"]*)" y la contraseña "([^"]*)"$`, p.ingreso)
	ctx.Step(`^la respuesta del ingreso no contiene la contraseña$`, p.respuestaSinContrasena)
	ctx.Step(`^el ingreso se rechaza con el mensaje "([^"]*)"$`, p.ingresoRechazado)
	ctx.Step(`^no hay ninguna sesión nueva$`, p.sinSesionNueva)
	ctx.Step(`^que "([^"]*)" inició sesión$`, p.inicioSesion)
	ctx.Step(`^pasan (\d+) días$`, p.pasanDias)
	ctx.Step(`^consulto quién soy con mi token$`, p.consultoQuienSoyConMiToken)
	ctx.Step(`^consulto quién soy con el token "([^"]*)"$`, p.consultoQuienSoyConElToken)
	ctx.Step(`^consulto quién soy sin token$`, p.consultoQuienSoySinToken)
	ctx.Step(`^el sistema responde que soy "([^"]*)"$`, p.respondeQueSoy)
	ctx.Step(`^la respuesta de quién soy no contiene la contraseña$`, p.respuestaSinContrasena)
	ctx.Step(`^la consulta de quién soy se rechaza con el mensaje "([^"]*)"$`, p.quienSoyRechazado)
	ctx.Step(`^"([^"]*)" abre una segunda sesión$`, p.abreUnaSegundaSesion)
	ctx.Step(`^cierro mi sesión$`, p.cierroMiSesion)
	ctx.Step(`^el cierre de sesión responde sin contenido$`, p.cierreSinContenido)
	ctx.Step(`^consulto quién soy con el token de mi segunda sesión$`, p.consultoQuienSoyConElTokenDeMiSegundaSesion)
	ctx.Step(`^consulto la lista de proyectos sin token$`, p.consultoLaListaDeProyectosSinToken)
	ctx.Step(`^la lista de proyectos se rechaza con el mensaje "([^"]*)"$`, p.listaDeProyectosRechazada)
}
