package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

type fakeCambiadorNombre struct {
	llamado     bool
	id          int64
	nombre      string
	normalizado string
	err         error
}

func (f *fakeCambiadorNombre) CambiarNombre(_ context.Context, id int64, nombre, normalizado string) error {
	f.llamado = true
	f.id = id
	f.nombre = nombre
	f.normalizado = normalizado
	return f.err
}

func TestCambiarNombreUsuario_Valido(t *testing.T) {
	repo := &fakeCambiadorNombre{}
	servicio := NewCambiarNombreUsuario(repo)

	actualizado, err := servicio.Ejecutar(context.Background(), 7, "Ana.Perez")

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if actualizado.ID != 7 || actualizado.Nombre != "Ana.Perez" || actualizado.NombreNormalizado != "ana.perez" {
		t.Errorf("usuario inesperado: %+v", actualizado)
	}
	if !repo.llamado || repo.id != 7 || repo.nombre != "Ana.Perez" || repo.normalizado != "ana.perez" {
		t.Errorf("llamada inesperada al repositorio: %+v", repo)
	}
}

func TestCambiarNombreUsuario_NombresInvalidos(t *testing.T) {
	casos := []string{
		"",
		"ab",
		strings.Repeat("a", 31),
		"ana perez",
		"ana@perez",
		"josé",
	}
	for _, nombre := range casos {
		t.Run(nombre, func(t *testing.T) {
			repo := &fakeCambiadorNombre{}
			servicio := NewCambiarNombreUsuario(repo)

			_, err := servicio.Ejecutar(context.Background(), 7, nombre)

			var verr domain.ValidationError
			if !errors.As(err, &verr) || verr.Campo != "nombre" {
				t.Fatalf("se esperaba ValidationError sobre nombre, se obtuvo %v", err)
			}
			if repo.llamado {
				t.Errorf("no se debía llamar al repositorio con un nombre inválido")
			}
		})
	}
}

func TestCambiarNombreUsuario_Duplicado(t *testing.T) {
	repo := &fakeCambiadorNombre{err: repository.ErrUsuarioDuplicado}
	servicio := NewCambiarNombreUsuario(repo)

	_, err := servicio.Ejecutar(context.Background(), 7, "ana.perez")

	if !errors.Is(err, domain.ErrNombreUsuarioEnUso) {
		t.Fatalf("se esperaba ErrNombreUsuarioEnUso, se obtuvo %v", err)
	}
}

func TestCambiarNombreUsuario_Inexistente(t *testing.T) {
	repo := &fakeCambiadorNombre{err: domain.ErrNoEncontrado}
	servicio := NewCambiarNombreUsuario(repo)

	_, err := servicio.Ejecutar(context.Background(), 999, "ana.perez")

	if !errors.Is(err, domain.ErrNoEncontrado) {
		t.Fatalf("se esperaba domain.ErrNoEncontrado, se obtuvo %v", err)
	}
}

func TestCambiarNombreUsuario_PropagaOtroError(t *testing.T) {
	repo := &fakeCambiadorNombre{err: errors.New("falla de persistencia")}
	servicio := NewCambiarNombreUsuario(repo)

	_, err := servicio.Ejecutar(context.Background(), 7, "ana.perez")

	if err == nil || !strings.Contains(err.Error(), "falla de persistencia") {
		t.Fatalf("se esperaba el error de persistencia, se obtuvo %v", err)
	}
}
