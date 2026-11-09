package http

import (
	"fmt"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func TestEscribirError_AutenticacionYRegistro(t *testing.T) {
	casos := []struct {
		nombre  string
		err     error
		estado  int
		mensaje string
	}{
		{"credenciales inválidas", domain.ErrCredencialesInvalidas, nethttp.StatusUnauthorized, "Usuario o contraseña incorrectos"},
		{"no autenticado", domain.ErrNoAutenticado, nethttp.StatusUnauthorized, "Tu sesión no es válida"},
		{"nombre de usuario en uso", domain.ErrNombreUsuarioEnUso, nethttp.StatusConflict, "Ese nombre de usuario ya existe"},
		{"credenciales inválidas envueltas", fmt.Errorf("ingresar: %w", domain.ErrCredencialesInvalidas), nethttp.StatusUnauthorized, "Usuario o contraseña incorrectos"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			rec := httptest.NewRecorder()

			escribirError(rec, c.err)

			if rec.Code != c.estado {
				t.Fatalf("se esperaba %d, se obtuvo %d (%s)", c.estado, rec.Code, rec.Body.String())
			}
			if resp := decodificarError(t, rec); resp.Mensaje != c.mensaje {
				t.Errorf("mensaje = %q, se esperaba %q", resp.Mensaje, c.mensaje)
			}
		})
	}
}

func TestEscribirError_LosErroresPreviosMantienenSuCodigo(t *testing.T) {
	casos := []struct {
		nombre string
		err    error
		estado int
	}{
		{"sprint no encontrado", domain.ErrSprintNoEncontrado, nethttp.StatusNotFound},
		{"no es scrum master", domain.ErrNoEsScrumMaster, nethttp.StatusForbidden},
		{"votación ya iniciada", domain.ErrVotacionYaIniciada, nethttp.StatusConflict},
		{"error desconocido", fmt.Errorf("boom"), nethttp.StatusInternalServerError},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			rec := httptest.NewRecorder()

			escribirError(rec, c.err)

			if rec.Code != c.estado {
				t.Errorf("se esperaba %d, se obtuvo %d", c.estado, rec.Code)
			}
		})
	}
}
