package repository

import (
	"context"
	"testing"
	"time"
)

func TestSQLiteProjectRepository_Listar_Vacio(t *testing.T) {
	repo := NewSQLiteProjectRepository(nuevaBDMigrada(t))

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

func TestSQLiteProjectRepository_Listar_OrdenadoPorID(t *testing.T) {
	db := nuevaBDMigrada(t)
	repo := NewSQLiteProjectRepository(db)
	mustExec(t, db, `INSERT INTO proyectos (id, nombre, descripcion, fecha_inicio, fecha_fin, creado_en) VALUES
		(9, 'Nueve', '', NULL, NULL, '2026-10-23T00:00:00Z'),
		(2, 'Dos', 'Descripción', '2026-03-01', NULL, '2026-10-23T00:00:00Z'),
		(5, 'Cinco', '', NULL, NULL, '2026-10-23T00:00:00Z')`)

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
	if proyectos[0].FechaInicio == nil {
		t.Fatal("se esperaba la fecha de inicio del proyecto 2")
	}
	if !proyectos[0].FechaInicio.Equal(time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("fecha de inicio inesperada: %v", proyectos[0].FechaInicio)
	}
}
