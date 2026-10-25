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
