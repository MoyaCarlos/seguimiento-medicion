package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func TestProjectRepositoryEnMemoria_AgregarIntegranteProyectoInexistente(t *testing.T) {
	repo := NewProjectRepositoryEnMemoria()
	_, err := repo.AgregarIntegrante(context.Background(), domain.Membership{
		ProjectID: 1, UserID: 1, Role: domain.RolScrumMaster,
	})
	if !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("se esperaba ErrProyectoNoEncontrado, se obtuvo %v", err)
	}
}

func TestProjectRepositoryEnMemoria_AgregarIntegranteNoValidaUsuario(t *testing.T) {
	repo := NewProjectRepositoryEnMemoria()
	guardado, err := repo.Guardar(context.Background(), domain.Project{Nombre: "Proyecto"})
	if err != nil {
		t.Fatalf("no se esperaba error al guardar: %v", err)
	}
	_, err = repo.AgregarIntegrante(context.Background(), domain.Membership{
		ProjectID: guardado.ID, UserID: 42, Role: domain.RolScrumMaster,
	})
	if err != nil {
		t.Fatalf("el fake no debe validar usuarios en otro repositorio: %v", err)
	}
}

func TestProjectRepositoryEnMemoria_GuardarConScrumMaster_RollbackSiFallaMembresia(t *testing.T) {
	repo := NewProjectRepositoryEnMemoria()
	repo.miembros = []domain.Membership{{ProjectID: 1, UserID: 42, Role: domain.RolScrumMaster}}

	_, err := repo.GuardarConScrumMaster(context.Background(), domain.Project{Nombre: "Proyecto"}, 42)
	if !errors.Is(err, ErrMiembroDuplicado) {
		t.Fatalf("se esperaba ErrMiembroDuplicado, se obtuvo %v", err)
	}
	if _, err := repo.ObtenerPorID(context.Background(), 1); !errors.Is(err, domain.ErrProyectoNoEncontrado) {
		t.Fatalf("el proyecto debía revertirse, se obtuvo %v", err)
	}
}

func TestProjectRepositoryEnMemoria_AgregarIntegranteDuplicado(t *testing.T) {
	repo := NewProjectRepositoryEnMemoria()
	guardado, _ := repo.Guardar(context.Background(), domain.Project{Nombre: "Proyecto"})
	usuario, _ := repo.GuardarUsuario(context.Background(), domain.User{Nombre: "Ana", NombreNormalizado: "ana"})
	m := domain.Membership{ProjectID: guardado.ID, UserID: usuario.ID, Role: domain.RolScrumMaster}

	if _, err := repo.AgregarIntegrante(context.Background(), m); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if _, err := repo.AgregarIntegrante(context.Background(), m); !errors.Is(err, ErrMiembroDuplicado) {
		t.Fatalf("se esperaba ErrMiembroDuplicado, se obtuvo %v", err)
	}
}

func TestProjectRepositoryEnMemoria_ListarIntegrantesDevuelveNombre(t *testing.T) {
	repo := NewProjectRepositoryEnMemoria()
	guardado, _ := repo.Guardar(context.Background(), domain.Project{Nombre: "Proyecto"})
	usuario, _ := repo.GuardarUsuario(context.Background(), domain.User{Nombre: "Candela", NombreNormalizado: "candela"})

	if _, err := repo.AgregarIntegrante(context.Background(), domain.Membership{
		ProjectID: guardado.ID, UserID: usuario.ID, Role: domain.RolProductBuilder,
	}); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	miembros, err := repo.ListarIntegrantes(context.Background(), guardado.ID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(miembros) != 1 || miembros[0].Nombre != "Candela" {
		t.Fatalf("se esperaba el integrante con su nombre, se obtuvo %+v", miembros)
	}
}
