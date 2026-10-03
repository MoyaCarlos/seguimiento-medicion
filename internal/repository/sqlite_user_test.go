package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func TestSQLiteUserRepository_CrearYBuscar(t *testing.T) {
	ctx := context.Background()
	repo := NewSQLiteUserRepository(abrirBDDePrueba(t))

	u, err := domain.NewUser("Ana Valentina")
	if err != nil {
		t.Fatalf("no se esperaba error de dominio: %v", err)
	}
	guardado, err := repo.Guardar(ctx, u)
	if err != nil {
		t.Fatalf("no se esperaba error al crear: %v", err)
	}
	if guardado.ID == 0 {
		t.Fatal("se esperaba un ID asignado")
	}

	encontrado, err := repo.ObtenerPorNombreNormalizado(ctx, "ana valentina")
	if err != nil {
		t.Fatalf("no se esperaba error al buscar: %v", err)
	}
	if encontrado.ID != guardado.ID || encontrado.Nombre != "Ana Valentina" {
		t.Errorf("usuario inesperado: %+v", encontrado)
	}
}

func TestSQLiteUserRepository_NoEncuentra(t *testing.T) {
	repo := NewSQLiteUserRepository(abrirBDDePrueba(t))

	_, err := repo.ObtenerPorNombreNormalizado(context.Background(), "inexistente")
	if !errors.Is(err, domain.ErrNoEncontrado) {
		t.Fatalf("se esperaba domain.ErrNoEncontrado, se obtuvo %v", err)
	}
}

func TestSQLiteUserRepository_NoDuplicaNombreNormalizado(t *testing.T) {
	ctx := context.Background()
	repo := NewSQLiteUserRepository(abrirBDDePrueba(t))

	primero, _ := domain.NewUser("Ana")
	if _, err := repo.Guardar(ctx, primero); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	// "  ana " normaliza a "ana": debe chocar con el UNIQUE.
	segundo, _ := domain.NewUser("  ana ")
	if _, err := repo.Guardar(ctx, segundo); !errors.Is(err, ErrUsuarioDuplicado) {
		t.Fatalf("se esperaba ErrUsuarioDuplicado, se obtuvo %v", err)
	}
}
