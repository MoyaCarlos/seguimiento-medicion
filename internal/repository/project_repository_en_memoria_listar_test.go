package repository

import (
	"context"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func TestProjectRepositoryEnMemoria_Listar_Vacio(t *testing.T) {
	repo := NewProjectRepositoryEnMemoria()

	proyectos, err := repo.Listar(context.Background())
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

func TestProjectRepositoryEnMemoria_Listar_OrdenadoPorID(t *testing.T) {
	repo := NewProjectRepositoryEnMemoria()
	repo.proyectos[9] = domain.Project{ID: 9, Nombre: "Nueve"}
	repo.proyectos[2] = domain.Project{ID: 2, Nombre: "Dos"}
	repo.proyectos[5] = domain.Project{ID: 5, Nombre: "Cinco"}

	proyectos, err := repo.Listar(context.Background())
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(proyectos) != 3 {
		t.Fatalf("se esperaban 3 proyectos, se obtuvieron %d", len(proyectos))
	}
	for i, esperado := range []int64{2, 5, 9} {
		if proyectos[i].ID != esperado {
			t.Errorf("posición %d: se esperaba el id %d, se obtuvo %d", i, esperado, proyectos[i].ID)
		}
	}
}

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
