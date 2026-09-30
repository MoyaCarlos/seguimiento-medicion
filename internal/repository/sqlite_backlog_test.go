package repository

import (
	"context"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func nuevaBDDePrueba(t *testing.T) *SQLiteBacklogRepository {
	t.Helper()
	db, err := AbrirSQLite(":memory:")
	if err != nil {
		t.Fatalf("no se pudo abrir sqlite: %v", err)
	}
	db.SetMaxOpenConns(1)
	if err := Migrar(context.Background(), db); err != nil {
		t.Fatalf("no se pudo migrar: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return NewSQLiteBacklogRepository(db)
}

func TestSQLiteBacklogRepository_Guardar(t *testing.T) {
	repo := nuevaBDDePrueba(t)
	ctx := context.Background()

	item, err := domain.NewBacklogItem(1, "Historia uno", "Descripción uno", domain.PrioridadMust, nil)
	if err != nil {
		t.Fatalf("no se esperaba error de dominio: %v", err)
	}

	guardado, err := repo.Guardar(ctx, item)
	if err != nil {
		t.Fatalf("no se esperaba error al guardar: %v", err)
	}
	if guardado.ID <= 0 {
		t.Fatalf("se esperaba un ID asignado, se obtuvo %d", guardado.ID)
	}

	var (
		estado       string
		estimacion   *int
		valorNegocio *int
	)
	err = repo.db.QueryRowContext(ctx,
		"SELECT estado, estimacion_sp, valor_negocio FROM backlog_items WHERE id = ?", guardado.ID,
	).Scan(&estado, &estimacion, &valorNegocio)
	if err != nil {
		t.Fatalf("no se pudo releer la historia: %v", err)
	}
	if estado != "Nueva" {
		t.Errorf("se esperaba estado %q, se obtuvo %q", "Nueva", estado)
	}
	if estimacion != nil {
		t.Errorf("se esperaba estimación NULL, se obtuvo %v", *estimacion)
	}
	if valorNegocio != nil {
		t.Errorf("se esperaba valor de negocio NULL, se obtuvo %v", *valorNegocio)
	}
}

func TestSQLiteBacklogRepository_OrdenDeCreacion(t *testing.T) {
	repo := nuevaBDDePrueba(t)
	ctx := context.Background()

	primero, _ := domain.NewBacklogItem(1, "Primera", "Descripción primera", domain.PrioridadMust, nil)
	segundo, _ := domain.NewBacklogItem(1, "Segunda", "Descripción segunda", domain.PrioridadShould, nil)

	g1, err := repo.Guardar(ctx, primero)
	if err != nil {
		t.Fatalf("error al guardar la primera historia: %v", err)
	}
	g2, err := repo.Guardar(ctx, segundo)
	if err != nil {
		t.Fatalf("error al guardar la segunda historia: %v", err)
	}
	if g2.ID <= g1.ID {
		t.Errorf("se esperaba que la segunda historia tuviera un ID posterior: %d <= %d", g2.ID, g1.ID)
	}
}

func TestSQLiteBacklogRepository_TitulosRepetidos(t *testing.T) {
	repo := nuevaBDDePrueba(t)
	ctx := context.Background()

	for i := 0; i < 2; i++ {
		item, err := domain.NewBacklogItem(1, "Mismo título", "Misma descripción", domain.PrioridadCould, nil)
		if err != nil {
			t.Fatalf("no se esperaba error de dominio: %v", err)
		}
		if _, err := repo.Guardar(ctx, item); err != nil {
			t.Fatalf("no se esperaba error al guardar títulos repetidos: %v", err)
		}
	}

	var total int
	if err := repo.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM backlog_items WHERE titulo = ?", "Mismo título").Scan(&total); err != nil {
		t.Fatalf("no se pudo contar: %v", err)
	}
	if total != 2 {
		t.Errorf("se esperaban 2 historias con el mismo título, se obtuvieron %d", total)
	}
}
