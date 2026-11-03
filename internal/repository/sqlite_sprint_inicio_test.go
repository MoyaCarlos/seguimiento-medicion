package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// leidoEIniciado lee el Sprint y lo inicia en memoria, sin guardarlo. Dos
// llamadas seguidas simulan dos inicios simultáneos: los dos lo leyeron
// Pendiente antes de que el primero llegara a la base.
func leidoEIniciado(t *testing.T, repo *SQLiteSprintRepository, id int64, goal string) domain.Sprint {
	t.Helper()
	leido, err := repo.ObtenerPorID(context.Background(), id)
	if err != nil {
		t.Fatalf("no se pudo leer el Sprint: %v", err)
	}
	inicio := time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC)
	if err := leido.Iniciar(goal, inicio, inicio.AddDate(0, 0, 14)); err != nil {
		t.Fatalf("no se esperaba error de dominio: %v", err)
	}
	return leido
}

func TestSQLiteSprintRepository_Actualizar_SegundoInicioSimultaneo(t *testing.T) {
	db := nuevaBDMigrada(t)
	repo := NewSQLiteSprintRepository(db)
	ctx := context.Background()
	pendiente := guardarSprint(t, repo, guardarProyecto(t, db))
	primero := leidoEIniciado(t, repo, pendiente.ID, "Goal del primero")
	segundo := leidoEIniciado(t, repo, pendiente.ID, "Goal del segundo")

	if err := repo.Actualizar(ctx, primero); err != nil {
		t.Fatalf("el primer inicio no debía fallar: %v", err)
	}
	err := repo.Actualizar(ctx, segundo)

	var verr domain.ValidationError
	if !errors.As(err, &verr) || verr.Campo != "estado" {
		t.Fatalf("se esperaba el error de estado de un Sprint que no está Pendiente, se obtuvo %v", err)
	}
	guardado, err := repo.ObtenerPorID(ctx, pendiente.ID)
	if err != nil {
		t.Fatalf("no se pudo releer el Sprint: %v", err)
	}
	if guardado.SprintGoal != "Goal del primero" {
		t.Errorf("el segundo inicio no debía pisar el Goal del primero, quedó %q", guardado.SprintGoal)
	}
}
