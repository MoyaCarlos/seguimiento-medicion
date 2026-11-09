package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// buscadorEspia cuenta las consultas que llegan al repositorio de sesiones.
type buscadorEspia struct {
	*fakeSesiones
	consultas int
}

func (b *buscadorEspia) ObtenerPorHuella(ctx context.Context, huella string) (domain.SesionConUsuario, error) {
	b.consultas++
	return b.fakeSesiones.ObtenerPorHuella(ctx, huella)
}

// Token de 43 caracteres base64url, el largo de 32 bytes sin relleno.
const tokenDePrueba = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"

func armarObtenerSesion(t *testing.T, ahora time.Time) (*ObtenerSesion, *buscadorEspia, domain.Sesion) {
	t.Helper()
	fake := &fakeSesiones{usuarios: map[int64]domain.User{
		7: {ID: 7, Nombre: "Jimena", NombreNormalizado: "jimena"},
	}}
	espia := &buscadorEspia{fakeSesiones: fake}
	sesion, err := fake.Crear(context.Background(), HuellaDeToken(tokenDePrueba), domain.Sesion{
		UsuarioID: 7,
		CreadoEn:  relojDeEmision(),
		VenceEn:   relojDeEmision().Add(domain.DuracionSesion),
	})
	if err != nil {
		t.Fatalf("preparar sesión: %v", err)
	}
	return NewObtenerSesion(espia, func() time.Time { return ahora }), espia, sesion
}

func TestObtenerSesion_TokenVigenteDevuelveSesionYUsuario(t *testing.T) {
	caso, _, sesion := armarObtenerSesion(t, relojDeEmision().Add(time.Hour))

	resultado, err := caso.Ejecutar(context.Background(), tokenDePrueba)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if resultado.ID != sesion.ID {
		t.Errorf("ID de sesión = %d, se esperaba %d", resultado.ID, sesion.ID)
	}
	if resultado.Usuario.ID != 7 {
		t.Errorf("Usuario.ID = %d, se esperaba 7", resultado.Usuario.ID)
	}
	if resultado.Usuario.Nombre != "Jimena" {
		t.Errorf("Usuario.Nombre = %q, se esperaba %q", resultado.Usuario.Nombre, "Jimena")
	}
}

func TestObtenerSesion_VencidaDevuelveNoAutenticado(t *testing.T) {
	vence := relojDeEmision().Add(domain.DuracionSesion)
	casos := []struct {
		nombre string
		ahora  time.Time
	}{
		{"justo en VenceEn", vence},
		{"después de VenceEn", vence.Add(time.Second)},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			caso, _, _ := armarObtenerSesion(t, c.ahora)

			_, err := caso.Ejecutar(context.Background(), tokenDePrueba)

			if !errors.Is(err, domain.ErrNoAutenticado) {
				t.Fatalf("se esperaba ErrNoAutenticado, se obtuvo %v", err)
			}
		})
	}
}

func TestObtenerSesion_CerradaDevuelveNoAutenticado(t *testing.T) {
	ahora := relojDeEmision().Add(time.Hour)
	caso, espia, sesion := armarObtenerSesion(t, ahora)
	if err := espia.Cerrar(context.Background(), sesion.ID, ahora); err != nil {
		t.Fatalf("cerrar sesión: %v", err)
	}

	_, err := caso.Ejecutar(context.Background(), tokenDePrueba)

	if !errors.Is(err, domain.ErrNoAutenticado) {
		t.Fatalf("se esperaba ErrNoAutenticado, se obtuvo %v", err)
	}
}

func TestObtenerSesion_InexistenteDevuelveNoAutenticado(t *testing.T) {
	caso, espia, _ := armarObtenerSesion(t, relojDeEmision().Add(time.Hour))
	otroToken := strings.Repeat("B", 43)

	_, err := caso.Ejecutar(context.Background(), otroToken)

	if !errors.Is(err, domain.ErrNoAutenticado) {
		t.Fatalf("se esperaba ErrNoAutenticado, se obtuvo %v", err)
	}
	if espia.consultas != 1 {
		t.Errorf("consultas = %d, se esperaba 1", espia.consultas)
	}
}

func TestObtenerSesion_TokenMalFormadoNoConsultaElRepositorio(t *testing.T) {
	casos := []struct {
		nombre string
		token  string
	}{
		{"vacío", ""},
		{"más corto de 43", strings.Repeat("A", 42)},
		{"más largo de 43", strings.Repeat("A", 44)},
		{"con signo más", strings.Repeat("A", 42) + "+"},
		{"con barra", strings.Repeat("A", 42) + "/"},
		{"con relleno igual", strings.Repeat("A", 42) + "="},
		{"con espacio", strings.Repeat("A", 42) + " "},
		{"con caracter no ASCII", strings.Repeat("A", 42) + "ñ"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			caso, espia, _ := armarObtenerSesion(t, relojDeEmision().Add(time.Hour))

			_, err := caso.Ejecutar(context.Background(), c.token)

			if !errors.Is(err, domain.ErrNoAutenticado) {
				t.Fatalf("se esperaba ErrNoAutenticado, se obtuvo %v", err)
			}
			if espia.consultas != 0 {
				t.Errorf("consultas = %d, se esperaba 0", espia.consultas)
			}
		})
	}
}
