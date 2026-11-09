package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// cerradorEspia cuenta las llamadas a Cerrar que llegan al repositorio.
type cerradorEspia struct {
	*fakeSesiones
	llamadas int
}

func (c *cerradorEspia) Cerrar(ctx context.Context, sesionID int64, ahora time.Time) error {
	c.llamadas++
	return c.fakeSesiones.Cerrar(ctx, sesionID, ahora)
}

// Segundo token de 43 caracteres, para una segunda sesión del mismo usuario.
var otroTokenDePrueba = strings.Repeat("B", 43)

// armarCerrarSesion usa el ObtenerSesion real sobre el fake, con dos sesiones del usuario 7.
func armarCerrarSesion(t *testing.T, ahora time.Time) (*CerrarSesion, *ObtenerSesion, *cerradorEspia) {
	t.Helper()
	fake := &fakeSesiones{usuarios: map[int64]domain.User{
		7: {ID: 7, Nombre: "Jimena", NombreNormalizado: "jimena"},
	}}
	for _, token := range []string{tokenDePrueba, otroTokenDePrueba} {
		_, err := fake.Crear(context.Background(), HuellaDeToken(token), domain.Sesion{
			UsuarioID: 7,
			CreadoEn:  relojDeEmision(),
			VenceEn:   relojDeEmision().Add(domain.DuracionSesion),
		})
		if err != nil {
			t.Fatalf("preparar sesión: %v", err)
		}
	}
	reloj := func() time.Time { return ahora }
	obtener := NewObtenerSesion(fake, reloj)
	espia := &cerradorEspia{fakeSesiones: fake}
	return NewCerrarSesion(obtener, espia, reloj), obtener, espia
}

func TestCerrarSesion_CierraSesionVigente(t *testing.T) {
	caso, obtener, espia := armarCerrarSesion(t, relojDeEmision().Add(time.Hour))

	if err := caso.Ejecutar(context.Background(), tokenDePrueba); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if espia.llamadas != 1 {
		t.Errorf("llamadas a Cerrar = %d, se esperaba 1", espia.llamadas)
	}
	_, err := obtener.Ejecutar(context.Background(), tokenDePrueba)
	if !errors.Is(err, domain.ErrNoAutenticado) {
		t.Fatalf("tras cerrar se esperaba ErrNoAutenticado, se obtuvo %v", err)
	}
}

func TestCerrarSesion_OtraSesionDelMismoUsuarioSigueVigente(t *testing.T) {
	caso, obtener, _ := armarCerrarSesion(t, relojDeEmision().Add(time.Hour))

	if err := caso.Ejecutar(context.Background(), tokenDePrueba); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	resultado, err := obtener.Ejecutar(context.Background(), otroTokenDePrueba)
	if err != nil {
		t.Fatalf("la otra sesión debía seguir vigente: %v", err)
	}
	if resultado.Usuario.ID != 7 {
		t.Errorf("Usuario.ID = %d, se esperaba 7", resultado.Usuario.ID)
	}
}

func TestCerrarSesion_TokenNoVigenteDevuelveNoAutenticadoSinCerrar(t *testing.T) {
	vence := relojDeEmision().Add(domain.DuracionSesion)
	casos := []struct {
		nombre string
		token  string
		ahora  time.Time
	}{
		{"token mal formado", "corto", relojDeEmision().Add(time.Hour)},
		{"token inexistente", strings.Repeat("C", 43), relojDeEmision().Add(time.Hour)},
		{"sesión vencida", tokenDePrueba, vence.Add(time.Second)},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			caso, _, espia := armarCerrarSesion(t, c.ahora)

			err := caso.Ejecutar(context.Background(), c.token)

			if !errors.Is(err, domain.ErrNoAutenticado) {
				t.Fatalf("se esperaba ErrNoAutenticado, se obtuvo %v", err)
			}
			if espia.llamadas != 0 {
				t.Errorf("llamadas a Cerrar = %d, se esperaba 0", espia.llamadas)
			}
		})
	}
}

func TestCerrarSesion_YaCerradaDevuelveNoAutenticadoSinCerrar(t *testing.T) {
	ahora := relojDeEmision().Add(time.Hour)
	caso, obtener, espia := armarCerrarSesion(t, ahora)
	sesion, err := obtener.Ejecutar(context.Background(), tokenDePrueba)
	if err != nil {
		t.Fatalf("obtener sesión: %v", err)
	}
	if err := espia.fakeSesiones.Cerrar(context.Background(), sesion.ID, ahora); err != nil {
		t.Fatalf("cerrar sesión: %v", err)
	}

	err = caso.Ejecutar(context.Background(), tokenDePrueba)

	if !errors.Is(err, domain.ErrNoAutenticado) {
		t.Fatalf("se esperaba ErrNoAutenticado, se obtuvo %v", err)
	}
	if espia.llamadas != 0 {
		t.Errorf("llamadas a Cerrar = %d, se esperaba 0", espia.llamadas)
	}
}

func TestCerrarSesion_CerrarDosVecesDaErrorLaSegunda(t *testing.T) {
	caso, _, espia := armarCerrarSesion(t, relojDeEmision().Add(time.Hour))

	if err := caso.Ejecutar(context.Background(), tokenDePrueba); err != nil {
		t.Fatalf("la primera vez no se esperaba error: %v", err)
	}
	err := caso.Ejecutar(context.Background(), tokenDePrueba)

	if !errors.Is(err, domain.ErrNoAutenticado) {
		t.Fatalf("la segunda vez se esperaba ErrNoAutenticado, se obtuvo %v", err)
	}
	if espia.llamadas != 1 {
		t.Errorf("llamadas a Cerrar = %d, se esperaba 1", espia.llamadas)
	}
}
