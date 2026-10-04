package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

func sembrarProyecto(t *testing.T, repo repository.ProjectRepository) domain.Project {
	t.Helper()
	p, err := domain.NewProject("Software Metrics & Estimation", "", nil, nil)
	if err != nil {
		t.Fatalf("no se esperaba error de dominio: %v", err)
	}
	guardado, err := repo.Guardar(context.Background(), p)
	if err != nil {
		t.Fatalf("no se esperaba error al sembrar el proyecto: %v", err)
	}
	return guardado
}

func TestCrearSprint_Exitoso(t *testing.T) {
	proyectos := repository.NewProjectRepositoryEnMemoria()
	sprints := repository.NewSprintRepositoryEnMemoria()
	proyecto := sembrarProyecto(t, proyectos)

	creado, err := NewCrearSprint(proyectos, sprints).Ejecutar(context.Background(), proyecto.ID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if creado.ID == 0 || creado.ProyectoID != proyecto.ID || creado.Estado != domain.SprintPendiente {
		t.Errorf("Sprint creado inesperado: %+v", creado)
	}
	if estadoPersistido(t, sprints, creado.ID) != domain.SprintPendiente {
		t.Error("el Sprint no quedó persistido en estado Pendiente")
	}
}

func TestCrearSprint_ProyectoInexistente(t *testing.T) {
	sprints := repository.NewSprintRepositoryEnMemoria()

	_, err := NewCrearSprint(repository.NewProjectRepositoryEnMemoria(), sprints).Ejecutar(context.Background(), 999)
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("se esperaba ErrProyectoNoEncontrado, se obtuvo %v", err)
	}
	guardados, err := sprints.ListarPorProyecto(context.Background(), 999)
	if err != nil {
		t.Fatalf("no se esperaba error al listar: %v", err)
	}
	if len(guardados) != 0 {
		t.Errorf("no debía guardarse ningún Sprint, hay %d", len(guardados))
	}
}

func TestCrearSprint_RechazaSegundoPendiente(t *testing.T) {
	proyectos := repository.NewProjectRepositoryEnMemoria()
	sprints := repository.NewSprintRepositoryEnMemoria()
	proyecto := sembrarProyecto(t, proyectos)
	sembrarSprintPendiente(t, sprints, proyecto.ID)

	_, err := NewCrearSprint(proyectos, sprints).Ejecutar(context.Background(), proyecto.ID)
	esperarValidacion(t, err, "estado")

	delProyecto, err := sprints.ListarPorProyecto(context.Background(), proyecto.ID)
	if err != nil {
		t.Fatalf("no se esperaba error al listar: %v", err)
	}
	if len(delProyecto) != 1 {
		t.Errorf("el proyecto debía seguir con un solo Sprint, tiene %d", len(delProyecto))
	}
}

func TestCrearSprint_PermitePendienteSiElOtroEstaActivo(t *testing.T) {
	proyectos := repository.NewProjectRepositoryEnMemoria()
	sprints := repository.NewSprintRepositoryEnMemoria()
	proyecto := sembrarProyecto(t, proyectos)
	activo := sembrarSprintPendiente(t, sprints, proyecto.ID)
	if _, err := NewIniciarSprint(sprints).Ejecutar(context.Background(), inputValido(activo.ID)); err != nil {
		t.Fatalf("no se esperaba error al iniciar: %v", err)
	}

	if _, err := NewCrearSprint(proyectos, sprints).Ejecutar(context.Background(), proyecto.ID); err != nil {
		t.Fatalf("con un Sprint Activo se debe poder preparar el siguiente: %v", err)
	}
}
