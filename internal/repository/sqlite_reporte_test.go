package repository

import (
	"context"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func insertarArrastrada(t *testing.T, repo *SQLiteBacklogRepository, sprintID, historiaID, proyectoID int64, sp int) {
	t.Helper()
	if _, err := repo.db.ExecContext(context.Background(),
		`INSERT INTO historias_arrastradas (sprint_id, historia_id, proyecto_id, estimacion_sp) VALUES (?, ?, ?, ?)`,
		sprintID, historiaID, proyectoID, sp,
	); err != nil {
		t.Fatalf("no se pudo sembrar la historia arrastrada: %v", err)
	}
}

func estimaciones(historias []domain.BacklogItem) []int {
	resultado := make([]int, 0, len(historias))
	for _, h := range historias {
		if h.EstimacionSP == nil {
			resultado = append(resultado, -1)
			continue
		}
		resultado = append(resultado, *h.EstimacionSP)
	}
	return resultado
}

func TestSQLiteBacklogRepository_ListarArrastradasPorSprint(t *testing.T) {
	repo := nuevaBDDePrueba(t)
	ctx := context.Background()
	sprints := NewSQLiteSprintRepository(repo.db)
	sprint := iniciarSprint(t, sprints, guardarSprint(t, sprints, 1))
	otroSprint := guardarSprint(t, sprints, 1)
	must := guardarHistoria(t, repo, 1, "M1", domain.PrioridadMust)
	should := guardarHistoria(t, repo, 1, "S1", domain.PrioridadShould)
	could := guardarHistoria(t, repo, 1, "C1", domain.PrioridadCould)
	wont := guardarHistoria(t, repo, 1, "W1", domain.PrioridadWont)
	otraDelMismoProyecto := guardarHistoria(t, repo, 1, "De otro Sprint", domain.PrioridadMust)

	insertarArrastrada(t, repo, sprint.ID, should.ID, 1, 3)
	insertarArrastrada(t, repo, sprint.ID, must.ID, 1, 5)
	insertarArrastrada(t, repo, sprint.ID, wont.ID, 1, 1)
	insertarArrastrada(t, repo, sprint.ID, could.ID, 1, 2)
	insertarArrastrada(t, repo, otroSprint.ID, otraDelMismoProyecto.ID, 1, 8)

	if _, err := repo.db.ExecContext(ctx, "UPDATE historias SET estimacion_sp = 99 WHERE id = ?", must.ID); err != nil {
		t.Fatalf("no se pudo editar la historia después del cierre: %v", err)
	}

	historias, err := repo.ListarArrastradasPorSprint(ctx, sprint.ID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	esperadoTitulos := []string{"M1", "S1", "C1", "W1"}
	if got := titulos(historias); len(got) != len(esperadoTitulos) {
		t.Fatalf("se esperaba %v, se obtuvo %v", esperadoTitulos, got)
	}
	for i, titulo := range esperadoTitulos {
		if historias[i].Titulo != titulo {
			t.Fatalf("se esperaba el orden %v, se obtuvo %v", esperadoTitulos, titulos(historias))
		}
		if historias[i].ProyectoID != 1 {
			t.Errorf("se esperaba el proyecto 1 para %q, se obtuvo %d", titulo, historias[i].ProyectoID)
		}
	}
	esperadoSP := []int{5, 3, 2, 1}
	got := estimaciones(historias)
	if len(got) != len(esperadoSP) {
		t.Fatalf("se esperaban las estimaciones %v, se obtuvo %v", esperadoSP, got)
	}
	for i := range esperadoSP {
		if got[i] != esperadoSP[i] {
			t.Fatalf("se esperaban las estimaciones %v, se obtuvo %v", esperadoSP, got)
		}
	}
}

func TestSQLiteBacklogRepository_ListarArrastradasPorSprint_VacioNoNil(t *testing.T) {
	repo := nuevaBDDePrueba(t)
	sprint := guardarSprint(t, NewSQLiteSprintRepository(repo.db), 1)

	historias, err := repo.ListarArrastradasPorSprint(context.Background(), sprint.ID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if historias == nil || len(historias) != 0 {
		t.Errorf("se esperaba un slice vacío no nil, se obtuvo %#v", historias)
	}
}
