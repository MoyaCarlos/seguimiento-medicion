package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

type stubBuscadorProyectos struct{ err error }

func (s stubBuscadorProyectos) ObtenerPorID(_ context.Context, id int64) (domain.Project, error) {
	return domain.Project{ID: id}, s.err
}

type stubListadorSprints struct {
	sprints []domain.Sprint
	err     error
}

func (s stubListadorSprints) ListarPorProyecto(_ context.Context, _ int64) ([]domain.Sprint, error) {
	return s.sprints, s.err
}

func TestConsultarSprints_OrdenaPorID(t *testing.T) {
	servicio := NewConsultarSprints(stubBuscadorProyectos{}, stubListadorSprints{
		sprints: []domain.Sprint{{ID: 3}, {ID: 1}, {ID: 2}},
	})

	sprints, err := servicio.Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(sprints) != 3 {
		t.Fatalf("se esperaban 3 Sprints, se obtuvieron %d", len(sprints))
	}
	for i, esperado := range []int64{1, 2, 3} {
		if sprints[i].ID != esperado {
			t.Errorf("posición %d: se esperaba el id %d, se obtuvo %d", i, esperado, sprints[i].ID)
		}
	}
}

func TestConsultarSprints_SinSprintsDevuelveArregloNoNulo(t *testing.T) {
	servicio := NewConsultarSprints(stubBuscadorProyectos{}, stubListadorSprints{})

	sprints, err := servicio.Ejecutar(context.Background(), 1)

	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if sprints == nil || len(sprints) != 0 {
		t.Fatalf("se esperaba un arreglo vacío no nulo, se obtuvo %#v", sprints)
	}
}

func TestConsultarSprints_ProyectoInexistente(t *testing.T) {
	servicio := NewConsultarSprints(
		stubBuscadorProyectos{err: domain.ErrProyectoNoEncontrado},
		stubListadorSprints{sprints: []domain.Sprint{{ID: 1}}},
	)

	_, err := servicio.Ejecutar(context.Background(), 999)

	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("se esperaba ErrProyectoNoEncontrado, se obtuvo %v", err)
	}
}

func TestConsultarSprints_PropagaErrorDePersistencia(t *testing.T) {
	errPersistencia := errors.New("fallo de persistencia")
	servicio := NewConsultarSprints(stubBuscadorProyectos{}, stubListadorSprints{err: errPersistencia})

	_, err := servicio.Ejecutar(context.Background(), 1)

	if !errors.Is(err, errPersistencia) {
		t.Fatalf("se esperaba el error de persistencia, se obtuvo %v", err)
	}
}
