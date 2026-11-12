package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

type stubListadorUsuarios struct {
	usuarios []domain.User
	err      error
}

func (s stubListadorUsuarios) ListarRegistrados(context.Context) ([]domain.User, error) {
	return s.usuarios, s.err
}

// HU-21, escenario 8: el Scrum Master elige a quién sumar entre los usuarios
// registrados.
func TestListarUsuarios_DevuelveLosRegistrados(t *testing.T) {
	registrados := []domain.User{{ID: 3, Nombre: "Ana"}, {ID: 1, Nombre: "bruno"}}

	usuarios, err := NewListarUsuarios(stubListadorUsuarios{usuarios: registrados}).Ejecutar(context.Background())

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(usuarios) != 2 || usuarios[0].ID != 3 || usuarios[1].ID != 1 {
		t.Errorf("se esperaban los registrados en el orden del repositorio, se obtuvo %+v", usuarios)
	}
}

func TestListarUsuarios_SinUsuariosDevuelveArregloNoNulo(t *testing.T) {
	usuarios, err := NewListarUsuarios(stubListadorUsuarios{}).Ejecutar(context.Background())

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if usuarios == nil || len(usuarios) != 0 {
		t.Fatalf("se esperaba un arreglo vacío (no nil), se obtuvo %#v", usuarios)
	}
}

func TestListarUsuarios_PropagaErrorDePersistencia(t *testing.T) {
	errPersistencia := errors.New("fallo de persistencia")

	_, err := NewListarUsuarios(stubListadorUsuarios{err: errPersistencia}).Ejecutar(context.Background())

	if !errors.Is(err, errPersistencia) {
		t.Fatalf("se esperaba el error de persistencia, se obtuvo %v", err)
	}
}
