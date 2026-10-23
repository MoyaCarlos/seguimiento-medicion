package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

type stubListadorProyectos struct {
	proyectos []domain.Project
	err       error
}

func (s *stubListadorProyectos) Listar(_ context.Context) ([]domain.Project, error) {
	return s.proyectos, s.err
}

func TestListarProyectos_OrdenaPorID(t *testing.T) {
	servicio := NewListarProyectos(&stubListadorProyectos{
		proyectos: []domain.Project{
			{ID: 3, Nombre: "Tres"},
			{ID: 1, Nombre: "Uno"},
			{ID: 2, Nombre: "Dos"},
		},
	})

	proyectos, err := servicio.Ejecutar(context.Background())

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(proyectos) != 3 {
		t.Fatalf("se esperaban 3 proyectos, se obtuvieron %d", len(proyectos))
	}
	for i, esperado := range []int64{1, 2, 3} {
		if proyectos[i].ID != esperado {
			t.Errorf("posición %d: se esperaba el id %d, se obtuvo %d", i, esperado, proyectos[i].ID)
		}
	}
}

func TestListarProyectos_VacioDevuelveArregloNoNulo(t *testing.T) {
	servicio := NewListarProyectos(&stubListadorProyectos{proyectos: []domain.Project{}})

	proyectos, err := servicio.Ejecutar(context.Background())

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if proyectos == nil {
		t.Fatal("se esperaba un arreglo vacío, se obtuvo nil")
	}
	if len(proyectos) != 0 {
		t.Fatalf("se esperaba una lista vacía, se obtuvo %+v", proyectos)
	}
}

func TestListarProyectos_PropagaErrorDePersistencia(t *testing.T) {
	errPersistencia := errors.New("fallo de persistencia")
	servicio := NewListarProyectos(&stubListadorProyectos{err: errPersistencia})

	_, err := servicio.Ejecutar(context.Background())

	if !errors.Is(err, errPersistencia) {
		t.Fatalf("se esperaba el error de persistencia, se obtuvo %v", err)
	}
}
