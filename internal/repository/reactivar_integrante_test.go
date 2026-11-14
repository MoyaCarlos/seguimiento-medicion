package repository

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// HU-16 (Escenario 3, D15): agregar de nuevo a quien fue dado de baja lo
// reactiva en la misma fila, con su creado_en original; un activo sigue
// siendo duplicado.

func TestProjectRepositoryEnMemoria_ReactivarIntegrante(t *testing.T) {
	ctx := context.Background()
	repo := NewProjectRepositoryEnMemoria()
	ana, _ := repo.GuardarUsuario(ctx, domain.User{Nombre: "Ana"})
	beto, _ := repo.GuardarUsuario(ctx, domain.User{Nombre: "Beto"})
	p, _ := repo.GuardarConScrumMaster(ctx, domain.Project{Nombre: "Uno"}, ana.ID)
	alta := time.Date(2026, 10, 2, 10, 0, 0, 0, time.UTC)
	if _, err := repo.AgregarIntegrante(ctx, domain.Membership{ProjectID: p.ID, UserID: beto.ID, Role: domain.RolProductBuilder, CreadoEn: alta}); err != nil {
		t.Fatalf("alta: %v", err)
	}
	if _, err := repo.AgregarIntegrante(ctx, domain.Membership{ProjectID: p.ID, UserID: beto.ID, Role: domain.RolProductBuilder}); !errors.Is(err, ErrMiembroDuplicado) {
		t.Fatalf("activo: se esperaba ErrMiembroDuplicado, se obtuvo %v", err)
	}
	if err := repo.DarDeBajaIntegrante(ctx, p.ID, beto.ID, fechaDeBaja); err != nil {
		t.Fatalf("baja: %v", err)
	}

	reactivado, err := repo.AgregarIntegrante(ctx, domain.Membership{ProjectID: p.ID, UserID: beto.ID, Role: domain.RolProductBuilder})
	if err != nil {
		t.Fatalf("reactivar: no se esperaba error: %v", err)
	}
	if !reactivado.CreadoEn.Equal(alta) {
		t.Errorf("se esperaba el creado_en original %v, se obtuvo %v", alta, reactivado.CreadoEn)
	}
	activos, _ := repo.ListarIntegrantes(ctx, p.ID)
	if got := nombresDe(activos); !reflect.DeepEqual(got, []string{"Ana:scrum_master", "Beto:product_builder"}) {
		t.Errorf("se esperaba a Beto una sola vez como product_builder, se obtuvo %v", got)
	}
	historicos, _ := repo.ListarIntegrantesHistoricos(ctx, p.ID)
	if len(historicos) != 2 {
		t.Errorf("la reactivación no debe duplicar el vínculo, se obtuvo %v", nombresDe(historicos))
	}
}

func TestSQLiteProjectRepository_ReactivarIntegrante(t *testing.T) {
	ctx := context.Background()
	db := nuevaBDMigrada(t)
	mustExec(t, db, equipoConVotosYHoras)
	repo := NewSQLiteProjectRepository(db)
	beto := domain.Membership{ProjectID: 1, UserID: 2, Role: domain.RolProductBuilder}

	if _, err := repo.AgregarIntegrante(ctx, beto); !errors.Is(err, ErrMiembroDuplicado) {
		t.Fatalf("activo: se esperaba ErrMiembroDuplicado, se obtuvo %v", err)
	}
	if err := repo.DarDeBajaIntegrante(ctx, 1, 2, fechaDeBaja); err != nil {
		t.Fatalf("baja: %v", err)
	}

	reactivado, err := repo.AgregarIntegrante(ctx, beto)
	if err != nil {
		t.Fatalf("reactivar: no se esperaba error: %v", err)
	}
	original := time.Date(2026, 10, 2, 0, 0, 0, 0, time.UTC)
	if !reactivado.CreadoEn.Equal(original) {
		t.Errorf("se esperaba el creado_en original %v, se obtuvo %v", original, reactivado.CreadoEn)
	}
	var filas int
	var baja *string
	if err := db.QueryRow("SELECT COUNT(*) FROM integrantes WHERE proyecto_id = 1 AND usuario_id = 2").Scan(&filas); err != nil || filas != 1 {
		t.Errorf("se esperaba una sola fila del vínculo, hay %d (%v)", filas, err)
	}
	if err := db.QueryRow("SELECT fecha_baja FROM integrantes WHERE proyecto_id = 1 AND usuario_id = 2").Scan(&baja); err != nil || baja != nil {
		t.Errorf("se esperaba fecha_baja NULL, se obtuvo %v (%v)", baja, err)
	}
	activos, _ := repo.ListarIntegrantes(ctx, 1)
	if got := nombresDe(activos); !reflect.DeepEqual(got, []string{"Ana:scrum_master", "Beto:product_builder", "Carla:product_builder"}) {
		t.Errorf("se esperaba a Beto una sola vez como product_builder, se obtuvo %v", got)
	}
	// Sus aportes previos siguen a su nombre: el voto revelado y las horas.
	if got := votosGuardados(t, db); got[0] != (votoGuardado{1000, 2, 5}) {
		t.Errorf("el voto revelado de Beto debe seguir, se obtuvo %v", got)
	}
	var horas int
	if err := db.QueryRow("SELECT COUNT(*) FROM registros_horas WHERE proyecto_id = 1 AND usuario_id = 2").Scan(&horas); err != nil || horas != 1 {
		t.Errorf("las horas de Beto deben seguir, se obtuvo %d (%v)", horas, err)
	}
	if _, err := repo.AgregarIntegrante(ctx, beto); !errors.Is(err, ErrMiembroDuplicado) {
		t.Errorf("reactivado: se esperaba ErrMiembroDuplicado, se obtuvo %v", err)
	}
}
