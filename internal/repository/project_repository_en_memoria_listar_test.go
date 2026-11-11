package repository

import (
	"context"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// HU-21, escenario 5: el repositorio en memoria responde igual que el de SQLite.
func TestProjectRepositoryEnMemoria_ListarPorIntegrante(t *testing.T) {
	ctx := context.Background()
	repo := NewProjectRepositoryEnMemoria()
	ana, _ := repo.GuardarUsuario(ctx, domain.User{Nombre: "Ana"})
	bruno, _ := repo.GuardarUsuario(ctx, domain.User{Nombre: "Bruno"})
	carla, _ := repo.GuardarUsuario(ctx, domain.User{Nombre: "Carla"})
	primero, _ := repo.GuardarConScrumMaster(ctx, domain.Project{Nombre: "Uno"}, bruno.ID)
	segundo, _ := repo.GuardarConScrumMaster(ctx, domain.Project{Nombre: "Dos"}, ana.ID)
	if _, err := repo.AgregarIntegrante(ctx, domain.Membership{ProjectID: primero.ID, UserID: ana.ID, Role: domain.RolProductBuilder}); err != nil {
		t.Fatalf("no se pudo preparar el integrante: %v", err)
	}

	deAna, err := repo.ListarPorIntegrante(ctx, ana.ID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(deAna) != 2 {
		t.Fatalf("se esperaban 2 proyectos de Ana, se obtuvieron %+v", deAna)
	}
	if deAna[0].Proyecto.ID != primero.ID || deAna[0].Rol != domain.RolProductBuilder {
		t.Errorf("primero: se esperaba %d como product_builder, se obtuvo %+v", primero.ID, deAna[0])
	}
	if deAna[1].Proyecto.ID != segundo.ID || deAna[1].Rol != domain.RolScrumMaster {
		t.Errorf("segundo: se esperaba %d como scrum_master, se obtuvo %+v", segundo.ID, deAna[1])
	}

	deCarla, err := repo.ListarPorIntegrante(ctx, carla.ID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if deCarla == nil || len(deCarla) != 0 {
		t.Fatalf("se esperaba un arreglo vacío (no nil), se obtuvo %#v", deCarla)
	}
}
