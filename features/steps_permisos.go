package features

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/cucumber/godog"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/app"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/service"
)

// contrasenaDeEscenario es la contraseña con la que el paso de sesión registra
// a cada persona en la API de prueba.
const contrasenaDeEscenario = "una-clave-de-prueba"

// sesionDeEscenario es el contexto compartido de HU-21: quién actúa en el
// escenario. Los pasos que llaman a los servicios leen el nombre con actor() y
// lo resuelven en su propia base; los que usan la API real piden con peticion(),
// que levanta la API armada como en producción (app.NuevoMux) solo si hace
// falta y registra o hace ingresar a la persona con /auth para obtener su token.
type sesionDeEscenario struct {
	nombre string
	db     *sql.DB
	mux    *http.ServeMux
	tokens map[string]string
	ahora  time.Time
}

func (s *sesionDeEscenario) reiniciar() {
	s.cerrar()
	s.nombre = ""
	s.tokens = map[string]string{}
	s.ahora = time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
}

func (s *sesionDeEscenario) cerrar() {
	if s.db != nil {
		_ = s.db.Close()
	}
	s.db, s.mux = nil, nil
}

// inicieSesionComo fija quién actúa en el escenario.
func (s *sesionDeEscenario) inicieSesionComo(nombre string) error {
	s.nombre = nombre
	return nil
}

// actor devuelve el nombre de quien actúa; vacío si el escenario no inició sesión.
func (s *sesionDeEscenario) actor() string {
	return s.nombre
}

// api devuelve la API completa sobre una base en memoria, creada al primer uso.
func (s *sesionDeEscenario) api() (*http.ServeMux, error) {
	if s.mux != nil {
		return s.mux, nil
	}
	db, err := repository.AbrirSQLite(":memory:")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if err := repository.AplicarMigraciones(context.Background(), db); err != nil {
		_ = db.Close()
		return nil, err
	}
	s.db = db
	s.mux = app.NuevoMux(db, app.Opciones{
		DemoHabilitado: func() bool { return false },
		Ahora:          func() time.Time { return s.ahora },
		Contrasenas:    service.NewServicioContrasenas(iteracionesDePruebaAutenticacion),
	})
	return s.mux, nil
}

// token registra a la persona en la API (o la hace ingresar si ya tiene
// cuenta) y devuelve su token de sesión.
func (s *sesionDeEscenario) token(nombre string) (string, error) {
	if t, ok := s.tokens[nombre]; ok {
		return t, nil
	}
	credenciales := map[string]string{"nombre": nombre, "contrasena": contrasenaDeEscenario}
	rec, err := s.enviar(http.MethodPost, "/auth/registro", credenciales, "")
	if err != nil {
		return "", err
	}
	if rec.Code == http.StatusConflict {
		if rec, err = s.enviar(http.MethodPost, "/auth/ingreso", credenciales, ""); err != nil {
			return "", err
		}
	}
	if rec.Code != http.StatusCreated && rec.Code != http.StatusOK {
		return "", fmt.Errorf("no se pudo iniciar la sesión de %q: %d %s", nombre, rec.Code, rec.Body.String())
	}
	var resp struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil || resp.Token == "" {
		return "", fmt.Errorf("respuesta de sesión inválida para %q: %s", nombre, rec.Body.String())
	}
	s.tokens[nombre] = resp.Token
	return resp.Token, nil
}

// peticion llama a la API con la sesión de quien actúa; sin sesión iniciada
// la petición viaja sin cabecera Authorization.
func (s *sesionDeEscenario) peticion(metodo, ruta string, cuerpo any) (*httptest.ResponseRecorder, error) {
	autorizacion := ""
	if s.nombre != "" {
		t, err := s.token(s.nombre)
		if err != nil {
			return nil, err
		}
		autorizacion = "Bearer " + t
	}
	return s.enviar(metodo, ruta, cuerpo, autorizacion)
}

func (s *sesionDeEscenario) enviar(metodo, ruta string, cuerpo any, autorizacion string) (*httptest.ResponseRecorder, error) {
	mux, err := s.api()
	if err != nil {
		return nil, err
	}
	var lector *bytes.Reader
	if cuerpo == nil {
		lector = bytes.NewReader(nil)
	} else {
		datos, err := json.Marshal(cuerpo)
		if err != nil {
			return nil, err
		}
		lector = bytes.NewReader(datos)
	}
	req := httptest.NewRequest(metodo, ruta, lector)
	if autorizacion != "" {
		req.Header.Set("Authorization", autorizacion)
	}
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec, nil
}

// InitializeScenarioSesion registra el paso de contexto de HU-21 y devuelve la
// sesión del escenario para los contextos que la necesitan.
func InitializeScenarioSesion(ctx *godog.ScenarioContext) *sesionDeEscenario {
	s := &sesionDeEscenario{}
	ctx.Before(func(c context.Context, _ *godog.Scenario) (context.Context, error) {
		s.reiniciar()
		return c, nil
	})
	ctx.After(func(c context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		s.cerrar()
		return c, nil
	})
	ctx.Step(`^que inicié sesión como "([^"]*)"$`, s.inicieSesionComo)
	return s
}
