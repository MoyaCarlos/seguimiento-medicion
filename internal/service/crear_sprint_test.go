package service

import (
	"context"
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
