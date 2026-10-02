package repository

import (
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func TestSQLiteUserRepository_CrearYBuscar(t *testing.T) {
	repo := NewSQLiteUserRepository(abrirBDDePrueba(t))

	u, err := domain.NewUser("Ana Valentina")
	if err != nil {
		t.Fatalf("no se esperaba error de dominio: %v", err)
	}
	if err := repo.Create(&u); err != nil {
		t.Fatalf("no se esperaba error al crear: %v", err)
	}
	if u.ID == "" {
		t.Fatal("se esperaba un ID asignado")
	}

	encontrado, err := repo.FindByNormalizedName("ana valentina")
	if err != nil {
		t.Fatalf("no se esperaba error al buscar: %v", err)
	}
	if encontrado.ID != u.ID || encontrado.Name != "Ana Valentina" {
		t.Errorf("usuario inesperado: %+v", encontrado)
	}
}

func TestSQLiteUserRepository_NoEncuentra(t *testing.T) {
	repo := NewSQLiteUserRepository(abrirBDDePrueba(t))

	_, err := repo.FindByNormalizedName("inexistente")
	if !errors.Is(err, ErrNoEncontrado) {
		t.Fatalf("se esperaba ErrNoEncontrado, se obtuvo %v", err)
	}
}

func TestSQLiteUserRepository_NoDuplicaNombreNormalizado(t *testing.T) {
	repo := NewSQLiteUserRepository(abrirBDDePrueba(t))

	primero, _ := domain.NewUser("Ana")
	if err := repo.Create(&primero); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	// "  ana " normaliza a "ana": debe chocar con el UNIQUE.
	segundo, _ := domain.NewUser("  ana ")
	if err := repo.Create(&segundo); !errors.Is(err, ErrUsuarioDuplicado) {
		t.Fatalf("se esperaba ErrUsuarioDuplicado, se obtuvo %v", err)
	}
}
