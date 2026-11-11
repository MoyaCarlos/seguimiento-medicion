package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func TestSQLiteUserRepository_CrearYBuscar(t *testing.T) {
	ctx := context.Background()
	repo := NewSQLiteUserRepository(nuevaBDMigrada(t))

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
	repo := NewSQLiteUserRepository(nuevaBDMigrada(t))

	_, err := repo.ObtenerPorNombreNormalizado(context.Background(), "inexistente")
	if !errors.Is(err, domain.ErrNoEncontrado) {
		t.Fatalf("se esperaba domain.ErrNoEncontrado, se obtuvo %v", err)
	}
}

func TestSQLiteUserRepository_NoDuplicaNombreNormalizado(t *testing.T) {
	ctx := context.Background()
	repo := NewSQLiteUserRepository(nuevaBDMigrada(t))

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

// HU-21, escenario 8: el equipo se arma solo con usuarios registrados (con
// contraseña), ordenados por nombre normalizado.
func TestSQLiteUserRepository_ListarRegistrados(t *testing.T) {
	db := nuevaBDMigrada(t)
	repo := NewSQLiteUserRepository(db)
	mustExec(t, db, `INSERT INTO usuarios (id, nombre, nombre_normalizado, contrasena_hash, creado_en) VALUES
		(1, 'carla', 'carla', 'hash-c', '2026-10-23T00:00:00Z'),
		(2, 'Sin Cuenta', 'sin cuenta', NULL, '2026-10-23T00:00:00Z'),
		(3, 'Ana', 'ana', 'hash-a', '2026-10-23T00:00:00Z'),
		(4, 'bruno', 'bruno', 'hash-b', '2026-10-23T00:00:00Z')`)

	usuarios, err := repo.ListarRegistrados(context.Background())
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	var obtenidos []string
	for _, u := range usuarios {
		obtenidos = append(obtenidos, u.Nombre)
	}
	if len(usuarios) != 3 || usuarios[0].ID != 3 || usuarios[1].ID != 4 || usuarios[2].ID != 1 {
		t.Fatalf("se esperaban Ana, bruno y carla (sin el usuario sin contraseña), se obtuvo %v", obtenidos)
	}
}

func TestSQLiteUserRepository_ListarRegistrados_Vacio(t *testing.T) {
	usuarios, err := NewSQLiteUserRepository(nuevaBDMigrada(t)).ListarRegistrados(context.Background())
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if usuarios == nil || len(usuarios) != 0 {
		t.Fatalf("se esperaba un arreglo vacío (no nil), se obtuvo %#v", usuarios)
	}
}
