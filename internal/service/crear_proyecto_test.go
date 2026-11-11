package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

const creadorDePrueba int64 = 7

func TestCrearProyecto_Ejecutar_ElCreadorDeLaSesionQuedaComoScrumMaster(t *testing.T) {
	proyectos := &fakeProjectRepository{}
	servicio := NewCrearProyecto(proyectos)

	proyecto, err := servicio.Ejecutar(context.Background(), CrearProyectoInput{
		Nombre:      "Software Metrics & Estimation",
		Descripcion: "Sistema de estimación",
		CreadorID:   creadorDePrueba,
	})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if proyecto.ID == 0 {
		t.Error("se esperaba ID asignado")
	}
	if proyectos.creates != 1 {
		t.Errorf("se esperaba 1 proyecto persistido, se crearon %d", proyectos.creates)
	}
	if len(proyectos.miembros) != 1 {
		t.Fatalf("se esperaba 1 vinculación, se obtuvieron %d", len(proyectos.miembros))
	}
	m := proyectos.miembros[0]
	if m.UserID != creadorDePrueba || m.Role != domain.RolScrumMaster || m.ProjectID != proyecto.ID {
		t.Errorf("se esperaba al creador %d como scrum_master del proyecto %d, se obtuvo %+v", creadorDePrueba, proyecto.ID, m)
	}
}

func TestCrearProyecto_Ejecutar_NoPersisteSiInvalido(t *testing.T) {
	proyectos := &fakeProjectRepository{}
	servicio := NewCrearProyecto(proyectos)

	_, err := servicio.Ejecutar(context.Background(), CrearProyectoInput{
		Nombre:    "   ",
		CreadorID: creadorDePrueba,
	})
	var verr domain.ValidationError
	if !errors.As(err, &verr) || verr.Campo != "nombre" {
		t.Fatalf("se esperaba ValidationError de nombre, se obtuvo %v", err)
	}
	if proyectos.creates != 0 || len(proyectos.miembros) != 0 {
		t.Error("no se debía persistir nada con datos inválidos")
	}
}

func TestCrearProyecto_Ejecutar_SinCreadorNoPersisteNada(t *testing.T) {
	proyectos := &fakeProjectRepository{}
	servicio := NewCrearProyecto(proyectos)

	_, err := servicio.Ejecutar(context.Background(), CrearProyectoInput{Nombre: "Proyecto"})

	if !errors.Is(err, domain.ErrNoAutenticado) {
		t.Fatalf("se esperaba ErrNoAutenticado, se obtuvo %v", err)
	}
	if proyectos.creates != 0 || len(proyectos.miembros) != 0 {
		t.Error("no se debía persistir el proyecto sin creador")
	}
}

func TestCrearProyecto_Ejecutar_RepoEnMemoria(t *testing.T) {
	proyectos := repository.NewProjectRepositoryEnMemoria()
	servicio := NewCrearProyecto(proyectos)

	proyecto, err := servicio.Ejecutar(context.Background(), CrearProyectoInput{
		Nombre:    "Proyecto",
		CreadorID: creadorDePrueba,
	})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	miembros, err := proyectos.ListarIntegrantes(context.Background(), proyecto.ID)
	if err != nil {
		t.Fatalf("no se esperaba error al listar integrantes: %v", err)
	}
	if len(miembros) != 1 || miembros[0].UserID != creadorDePrueba || miembros[0].Role != domain.RolScrumMaster {
		t.Fatalf("se esperaba el Scrum Master vinculado, se obtuvo %+v", miembros)
	}
}
