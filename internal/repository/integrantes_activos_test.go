package repository

import (
	"context"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// Integrante = integrante activo (HU-16, FR-001): ListarIntegrantes y
// ListarPorIntegrante ignoran a quien tiene fecha de baja, y
// ListarIntegrantesHistoricos lo sigue devolviendo con su nombre y rol.

// marcarBajaEnMemoria da de baja directo sobre el estado interno, sin pasar
// por DarDeBajaIntegrante (se prueba aparte).
func marcarBajaEnMemoria(repo *ProjectRepositoryEnMemoria, proyectoID, usuarioID int64) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.bajas[claveIntegrante{proyecto: proyectoID, usuario: usuarioID}] = time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)
}

func nombresDe(miembros []domain.Member) []string {
	nombres := make([]string, len(miembros))
	for i, m := range miembros {
		nombres[i] = m.Nombre + ":" + string(m.Role)
	}
	return nombres
}

func verificarActivosEHistoricos(t *testing.T, activos, historicos []domain.Member) {
	t.Helper()
	if got := nombresDe(activos); len(got) != 1 || got[0] != "Ana:scrum_master" {
		t.Errorf("ListarIntegrantes: se esperaba solo a Ana (scrum_master), se obtuvo %v", got)
	}
	if got := nombresDe(historicos); len(got) != 2 || got[0] != "Ana:scrum_master" || got[1] != "Beto:product_builder" {
		t.Errorf("ListarIntegrantesHistoricos: se esperaba [Ana:scrum_master Beto:product_builder], se obtuvo %v", got)
	}
}

func TestProjectRepositoryEnMemoria_IntegrantesActivosEHistoricos(t *testing.T) {
	ctx := context.Background()
	repo := NewProjectRepositoryEnMemoria()
	ana, _ := repo.GuardarUsuario(ctx, domain.User{Nombre: "Ana"})
	beto, _ := repo.GuardarUsuario(ctx, domain.User{Nombre: "Beto"})
	p, _ := repo.GuardarConScrumMaster(ctx, domain.Project{Nombre: "Uno"}, ana.ID)
	if _, err := repo.AgregarIntegrante(ctx, domain.Membership{ProjectID: p.ID, UserID: beto.ID, Role: domain.RolProductBuilder}); err != nil {
		t.Fatalf("no se pudo preparar el integrante: %v", err)
	}
	marcarBajaEnMemoria(repo, p.ID, beto.ID)

	activos, err := repo.ListarIntegrantes(ctx, p.ID)
	if err != nil {
		t.Fatalf("ListarIntegrantes: %v", err)
	}
	historicos, err := repo.ListarIntegrantesHistoricos(ctx, p.ID)
	if err != nil {
		t.Fatalf("ListarIntegrantesHistoricos: %v", err)
	}
	verificarActivosEHistoricos(t, activos, historicos)

	deBeto, err := repo.ListarPorIntegrante(ctx, beto.ID)
	if err != nil {
		t.Fatalf("ListarPorIntegrante: %v", err)
	}
	if len(deBeto) != 0 {
		t.Errorf("ListarPorIntegrante: el proyecto donde Beto está dado de baja no debe aparecer, se obtuvo %+v", deBeto)
	}
}

func TestSQLiteProjectRepository_IntegrantesActivosEHistoricos(t *testing.T) {
	ctx := context.Background()
	db := nuevaBDMigrada(t)
	repo := NewSQLiteProjectRepository(db)
	mustExec(t, db, `INSERT INTO usuarios (id, nombre, nombre_normalizado, creado_en) VALUES
		(1, 'Ana', 'ana', '2026-10-23T00:00:00Z'),
		(2, 'Beto', 'beto', '2026-10-23T00:00:00Z')`)
	mustExec(t, db, `INSERT INTO proyectos (id, nombre, creado_en) VALUES (1, 'Uno', '2026-10-23T00:00:00Z')`)
	mustExec(t, db, `INSERT INTO integrantes (proyecto_id, usuario_id, rol, creado_en) VALUES
		(1, 1, 'scrum_master', '2026-10-23T00:00:00Z'),
		(1, 2, 'product_builder', '2026-10-24T00:00:00Z')`)
	mustExec(t, db, `UPDATE integrantes SET fecha_baja = '2026-11-12' WHERE proyecto_id = 1 AND usuario_id = 2`)

	activos, err := repo.ListarIntegrantes(ctx, 1)
	if err != nil {
		t.Fatalf("ListarIntegrantes: %v", err)
	}
	historicos, err := repo.ListarIntegrantesHistoricos(ctx, 1)
	if err != nil {
		t.Fatalf("ListarIntegrantesHistoricos: %v", err)
	}
	verificarActivosEHistoricos(t, activos, historicos)

	deBeto, err := repo.ListarPorIntegrante(ctx, 2)
	if err != nil {
		t.Fatalf("ListarPorIntegrante: %v", err)
	}
	if len(deBeto) != 0 {
		t.Errorf("ListarPorIntegrante: el proyecto donde Beto está dado de baja no debe aparecer, se obtuvo %+v", deBeto)
	}
}
