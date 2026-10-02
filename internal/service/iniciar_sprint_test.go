package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

var (
	inicioSprint = time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC)
	finSprint    = inicioSprint.AddDate(0, 0, 14)
)

func sembrarSprintPendiente(t *testing.T, repo repository.SprintRepository, proyectoID int64) domain.Sprint {
	t.Helper()
	s, err := domain.NewSprint(proyectoID)
	if err != nil {
		t.Fatalf("no se esperaba error de dominio: %v", err)
	}
	guardado, err := repo.Guardar(context.Background(), s)
	if err != nil {
		t.Fatalf("no se esperaba error al sembrar: %v", err)
	}
	return guardado
}

func inputValido(sprintID int64) IniciarSprintInput {
	return IniciarSprintInput{SprintID: sprintID, SprintGoal: "Entregar el MVP", FechaInicio: inicioSprint, FechaFin: finSprint}
}

func estadoPersistido(t *testing.T, repo repository.SprintRepository, id int64) domain.EstadoSprint {
	t.Helper()
	s, err := repo.ObtenerPorID(context.Background(), id)
	if err != nil {
		t.Fatalf("no se pudo releer el Sprint: %v", err)
	}
	return s.Estado
}

func esperarValidacion(t *testing.T, err error, campo string) {
	t.Helper()
	var verr domain.ValidationError
	if !errors.As(err, &verr) {
		t.Fatalf("se esperaba ValidationError en %q, se obtuvo %v", campo, err)
	}
	if verr.Campo != campo {
		t.Errorf("se esperaba campo %q, se obtuvo %q", campo, verr.Campo)
	}
}

func TestIniciarSprint_Exitoso(t *testing.T) {
	repo := repository.NewSprintRepositoryEnMemoria()
	sprint := sembrarSprintPendiente(t, repo, 1)

	resultado, err := NewIniciarSprint(repo).Ejecutar(context.Background(), inputValido(sprint.ID))
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if resultado.Estado != domain.SprintActivo || resultado.SprintGoal != "Entregar el MVP" {
		t.Errorf("resultado inesperado: %+v", resultado)
	}
	if estadoPersistido(t, repo, sprint.ID) != domain.SprintActivo {
		t.Errorf("el Sprint no quedó Activo en el repositorio")
	}
}

func TestIniciarSprint_RechazaSegundoActivoEnElProyecto(t *testing.T) {
	repo := repository.NewSprintRepositoryEnMemoria()
	servicio := NewIniciarSprint(repo)
	primero := sembrarSprintPendiente(t, repo, 1)
	segundo := sembrarSprintPendiente(t, repo, 1)
	if _, err := servicio.Ejecutar(context.Background(), inputValido(primero.ID)); err != nil {
		t.Fatalf("no se esperaba error al iniciar el primero: %v", err)
	}

	_, err := servicio.Ejecutar(context.Background(), inputValido(segundo.ID))

	esperarValidacion(t, err, "estado")
	if estadoPersistido(t, repo, segundo.ID) != domain.SprintPendiente {
		t.Errorf("el segundo Sprint debía permanecer Pendiente")
	}
}

func TestIniciarSprint_UnActivoEnOtroProyectoNoBloquea(t *testing.T) {
	repo := repository.NewSprintRepositoryEnMemoria()
	servicio := NewIniciarSprint(repo)
	otroProyecto := sembrarSprintPendiente(t, repo, 2)
	if _, err := servicio.Ejecutar(context.Background(), inputValido(otroProyecto.ID)); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	sprint := sembrarSprintPendiente(t, repo, 1)

	if _, err := servicio.Ejecutar(context.Background(), inputValido(sprint.ID)); err != nil {
		t.Fatalf("un Sprint activo de otro proyecto no debía bloquear: %v", err)
	}
}

func TestIniciarSprint_FechasInvalidasNoPersiste(t *testing.T) {
	repo := repository.NewSprintRepositoryEnMemoria()
	sprint := sembrarSprintPendiente(t, repo, 1)
	input := inputValido(sprint.ID)
	input.FechaFin = input.FechaInicio

	_, err := NewIniciarSprint(repo).Ejecutar(context.Background(), input)

	esperarValidacion(t, err, "fecha_fin")
	if estadoPersistido(t, repo, sprint.ID) != domain.SprintPendiente {
		t.Errorf("el Sprint debía permanecer Pendiente")
	}
}

func TestIniciarSprint_NoEncontrado(t *testing.T) {
	repo := repository.NewSprintRepositoryEnMemoria()

	_, err := NewIniciarSprint(repo).Ejecutar(context.Background(), inputValido(999))

	if !errors.Is(err, domain.ErrSprintNoEncontrado) {
		t.Fatalf("se esperaba ErrSprintNoEncontrado, se obtuvo %v", err)
	}
}
