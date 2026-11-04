package service

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

type contadorDeRegistros struct {
	*fakeCredenciales
	mu       sync.Mutex
	llamadas int
}

func (c *contadorDeRegistros) Registrar(ctx context.Context, u domain.User, hash string) (domain.User, error) {
	c.mu.Lock()
	c.llamadas++
	c.mu.Unlock()
	return c.fakeCredenciales.Registrar(ctx, u, hash)
}

func (c *contadorDeRegistros) registros() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.llamadas
}

type hasheadorQueCuenta struct {
	llamadas int
}

func (h *hasheadorQueCuenta) Hashear(contrasena string) (string, error) {
	h.llamadas++
	return "hash:" + contrasena, nil
}

func (h *hasheadorQueCuenta) Verificar(contrasena, hash string) bool {
	return hash == "hash:"+contrasena
}

func armarRegistrarUsuario() (*RegistrarUsuario, *contadorDeRegistros, *hasheadorQueCuenta, *fakeSesiones) {
	credenciales := &contadorDeRegistros{fakeCredenciales: &fakeCredenciales{}}
	hasheador := &hasheadorQueCuenta{}
	sesiones := &fakeSesiones{}
	caso := NewRegistrarUsuario(credenciales, hasheador, NewEmisorSesiones(sesiones, relojDeEmision))
	return caso, credenciales, hasheador, sesiones
}

func TestRegistrarUsuario_ExitoGuardaElNombreTalCualYNormalizado(t *testing.T) {
	caso, credenciales, _, _ := armarRegistrarUsuario()

	resultado, err := caso.Ejecutar(context.Background(), "Jimena.Mo", "clave-segura-1")
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if resultado.Usuario.Nombre != "Jimena.Mo" {
		t.Errorf("Nombre = %q, se esperaba %q", resultado.Usuario.Nombre, "Jimena.Mo")
	}
	if resultado.Usuario.NombreNormalizado != "jimena.mo" {
		t.Errorf("NombreNormalizado = %q, se esperaba %q", resultado.Usuario.NombreNormalizado, "jimena.mo")
	}
	if _, ok := credenciales.usuarios["jimena.mo"]; !ok {
		t.Errorf("el usuario no quedó guardado en el repositorio")
	}
}

func TestRegistrarUsuario_ExitoGuardaUnHashDistintoDeLaContrasena(t *testing.T) {
	caso, credenciales, _, _ := armarRegistrarUsuario()

	if _, err := caso.Ejecutar(context.Background(), "jimena", "clave-segura-1"); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	hash := credenciales.hashes["jimena"]
	if hash == "" || hash == "clave-segura-1" {
		t.Errorf("se esperaba un hash distinto del texto plano, se obtuvo %q", hash)
	}
}

func TestRegistrarUsuario_ExitoEmiteUnaSesionConToken(t *testing.T) {
	caso, _, _, sesiones := armarRegistrarUsuario()

	resultado, err := caso.Ejecutar(context.Background(), "jimena", "clave-segura-1")
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if resultado.Token == "" {
		t.Errorf("se esperaba un token")
	}
	if sesiones.creadas != 1 {
		t.Errorf("sesiones creadas = %d, se esperaba 1", sesiones.creadas)
	}
	if resultado.VenceEn.IsZero() {
		t.Errorf("se esperaba una fecha de vencimiento")
	}
}

func TestRegistrarUsuario_NombreInvalidoNoLlamaAlRegistroNiAlHash(t *testing.T) {
	caso, credenciales, hasheador, _ := armarRegistrarUsuario()

	_, err := caso.Ejecutar(context.Background(), "ab", "clave-segura-1")

	assertCampoInvalido(t, err, "nombre")
	if credenciales.registros() != 0 {
		t.Errorf("no se esperaba llamar al registro")
	}
	if hasheador.llamadas != 0 {
		t.Errorf("no se esperaba hashear")
	}
}

func TestRegistrarUsuario_ValidaLargoDeContrasena(t *testing.T) {
	casos := []struct {
		nombre     string
		contrasena string
		campo      string
	}{
		{"7 caracteres es inválida", strings.Repeat("a", 7), "contrasena"},
		{"128 caracteres es válida", strings.Repeat("a", 128), ""},
		{"129 caracteres es inválida", strings.Repeat("a", 129), "contrasena"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			caso, _, _, _ := armarRegistrarUsuario()

			_, err := caso.Ejecutar(context.Background(), "jimena", c.contrasena)

			if c.campo == "" {
				if err != nil {
					t.Fatalf("no se esperaba error: %v", err)
				}
				return
			}
			assertCampoInvalido(t, err, c.campo)
		})
	}
}

func TestRegistrarUsuario_ValidaPrimeroElNombre(t *testing.T) {
	caso, _, _, _ := armarRegistrarUsuario()

	_, err := caso.Ejecutar(context.Background(), "ab", "corta")

	assertCampoInvalido(t, err, "nombre")
}

func TestRegistrarUsuario_NombreRepetidoConOtrasMayusculasEsNombreEnUso(t *testing.T) {
	caso, _, _, sesiones := armarRegistrarUsuario()
	if _, err := caso.Ejecutar(context.Background(), "Jimena", "clave-segura-1"); err != nil {
		t.Fatalf("no se esperaba error en el primer registro: %v", err)
	}

	_, err := caso.Ejecutar(context.Background(), "JIMENA", "otra-clave-22")

	if !errors.Is(err, domain.ErrNombreUsuarioEnUso) {
		t.Fatalf("se esperaba ErrNombreUsuarioEnUso, se obtuvo %v", err)
	}
	if sesiones.creadas != 1 {
		t.Errorf("sesiones creadas = %d, se esperaba 1 (solo la del primer registro)", sesiones.creadas)
	}
}

func TestRegistrarUsuario_ErrorDelRepositorioSePropaga(t *testing.T) {
	caso, credenciales, _, sesiones := armarRegistrarUsuario()
	falla := errors.New("base caída")
	credenciales.err = falla

	_, err := caso.Ejecutar(context.Background(), "jimena", "clave-segura-1")

	if !errors.Is(err, falla) {
		t.Fatalf("se esperaba el error del repositorio, se obtuvo %v", err)
	}
	if errors.Is(err, domain.ErrNombreUsuarioEnUso) {
		t.Errorf("un error que no es duplicado no debe convertirse en nombre en uso")
	}
	if sesiones.creadas != 0 {
		t.Errorf("no se esperaba emitir sesión")
	}
}

func TestRegistrarUsuario_DosRegistrosSimultaneosConElMismoNombreDanUnSoloExito(t *testing.T) {
	caso, _, _, sesiones := armarRegistrarUsuario()
	errores := make([]error, 2)
	var espera sync.WaitGroup
	arranque := make(chan struct{})

	for i := range errores {
		espera.Add(1)
		go func() {
			defer espera.Done()
			<-arranque
			_, errores[i] = caso.Ejecutar(context.Background(), "jimena", "clave-segura-1")
		}()
	}
	close(arranque)
	espera.Wait()

	exitos, enUso := 0, 0
	for _, err := range errores {
		switch {
		case err == nil:
			exitos++
		case errors.Is(err, domain.ErrNombreUsuarioEnUso):
			enUso++
		default:
			t.Errorf("error inesperado: %v", err)
		}
	}
	if exitos != 1 || enUso != 1 {
		t.Errorf("éxitos = %d, nombre en uso = %d; se esperaba 1 y 1", exitos, enUso)
	}
	if sesiones.creadas != 1 {
		t.Errorf("sesiones creadas = %d, se esperaba 1", sesiones.creadas)
	}
}

func assertCampoInvalido(t *testing.T, err error, campo string) {
	t.Helper()
	var validacion domain.ValidationError
	if !errors.As(err, &validacion) {
		t.Fatalf("se esperaba ValidationError, se obtuvo %v", err)
	}
	if validacion.Campo != campo {
		t.Errorf("Campo = %q, se esperaba %q", validacion.Campo, campo)
	}
}
