package repository

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

// historiaEnSprint inserta una historia asignada al Sprint. La asignación es de
// HU-06; acá se simula con SQL directo.
func historiaEnSprint(t *testing.T, db *sql.DB, proyectoID, sprintID int64, estado domain.Estado, sp *int) domain.BacklogItem {
	t.Helper()
	var finalizacion any
	if estado == domain.EstadoCompletada {
		finalizacion = "2026-10-20"
	}
	h := domain.BacklogItem{ProyectoID: proyectoID, Estado: estado, EstimacionSP: sp, SprintID: &sprintID}
	err := db.QueryRow(
		`INSERT INTO historias (proyecto_id, titulo, descripcion, prioridad, estado, estimacion_sp, sprint_id, fecha_finalizacion)
		 VALUES (?, 't', 'd', 'M', ?, ?, ?, ?) RETURNING id`,
		proyectoID, string(estado), punteroAValorSQL(sp), sprintID, finalizacion,
	).Scan(&h.ID)
	if err != nil {
		t.Fatalf("no se pudo sembrar la historia: %v", err)
	}
	return h
}

type arrastre struct {
	historiaID, proyectoID int64
	estimacionSP           sql.NullInt64
}

func arrastradasDe(t *testing.T, db *sql.DB, sprintID int64) []arrastre {
	t.Helper()
	filas, err := db.Query(
		"SELECT historia_id, proyecto_id, estimacion_sp FROM historias_arrastradas WHERE sprint_id = ? ORDER BY historia_id", sprintID)
	if err != nil {
		t.Fatalf("no se pudieron leer las historias arrastradas: %v", err)
	}
	defer filas.Close()
	var resultado []arrastre
	for filas.Next() {
		var a arrastre
		if err := filas.Scan(&a.historiaID, &a.proyectoID, &a.estimacionSP); err != nil {
			t.Fatalf("no se pudo leer la historia arrastrada: %v", err)
		}
		resultado = append(resultado, a)
	}
	return resultado
}

func sprintYEstadoDeHistoria(t *testing.T, db *sql.DB, id int64) (sql.NullInt64, string) {
	t.Helper()
	var (
		sprintID sql.NullInt64
		estado   string
	)
	if err := db.QueryRow("SELECT sprint_id, estado FROM historias WHERE id = ?", id).Scan(&sprintID, &estado); err != nil {
		t.Fatalf("no se pudo leer la historia %d: %v", id, err)
	}
	return sprintID, estado
}

func cerrado(t *testing.T, s domain.Sprint) domain.Sprint {
	t.Helper()
	if err := s.Cerrar(); err != nil {
		t.Fatalf("no se esperaba error de dominio: %v", err)
	}
	return s
}

func sp(v int) *int { return &v }

func TestSQLiteSprintRepository_CerrarConArrastre_RegistraFotoYLibera(t *testing.T) {
	db := nuevaBDMigrada(t)
	repo := NewSQLiteSprintRepository(db)
	ctx := context.Background()
	proyectoID := guardarProyecto(t, db)
	activo := iniciarSprint(t, repo, guardarSprint(t, repo, proyectoID))
	completada := historiaEnSprint(t, db, proyectoID, activo.ID, domain.EstadoCompletada, sp(3))
	enProgreso := historiaEnSprint(t, db, proyectoID, activo.ID, domain.Estado("En progreso"), sp(8))
	sinEstimar := historiaEnSprint(t, db, proyectoID, activo.ID, domain.EstadoNueva, nil)

	err := repo.CerrarConArrastre(ctx, cerrado(t, activo), []domain.BacklogItem{enProgreso, sinEstimar})
	if err != nil {
		t.Fatalf("no se esperaba error al cerrar: %v", err)
	}

	if s, _ := repo.ObtenerPorID(ctx, activo.ID); s.Estado != domain.SprintFinalizado {
		t.Errorf("el Sprint debía quedar Finalizado, quedó %q", s.Estado)
	}
	esperadas := []arrastre{
		{historiaID: enProgreso.ID, proyectoID: proyectoID, estimacionSP: sql.NullInt64{Int64: 8, Valid: true}},
		{historiaID: sinEstimar.ID, proyectoID: proyectoID},
	}
	if got := arrastradasDe(t, db, activo.ID); !slices.Equal(got, esperadas) {
		t.Errorf("arrastradas: se esperaba %+v, se obtuvo %+v", esperadas, got)
	}
	for _, id := range []int64{enProgreso.ID, sinEstimar.ID} {
		if sprintID, estado := sprintYEstadoDeHistoria(t, db, id); sprintID.Valid || estado != string(domain.EstadoNueva) {
			t.Errorf("la historia %d debía volver al Product Backlog en Nueva, quedó sprint=%v estado=%q", id, sprintID, estado)
		}
	}
	if sprintID, estado := sprintYEstadoDeHistoria(t, db, completada.ID); sprintID.Int64 != activo.ID || estado != string(domain.EstadoCompletada) {
		t.Errorf("la historia completada debía seguir en el Sprint, quedó sprint=%v estado=%q", sprintID, estado)
	}

	mustExec(t, db, "UPDATE historias SET estimacion_sp = 13 WHERE id = ?", enProgreso.ID)
	if got := arrastradasDe(t, db, activo.ID); got[0].estimacionSP.Int64 != 8 {
		t.Errorf("la foto de Story Points no debía cambiar al editar la historia, quedó %v", got[0].estimacionSP)
	}
}

func TestSQLiteSprintRepository_CerrarConArrastre_SinNoCompletadas(t *testing.T) {
	db := nuevaBDMigrada(t)
	repo := NewSQLiteSprintRepository(db)
	activo := iniciarSprint(t, repo, guardarSprint(t, repo, guardarProyecto(t, db)))

	if err := repo.CerrarConArrastre(context.Background(), cerrado(t, activo), nil); err != nil {
		t.Fatalf("cerrar sin historias no completadas no es un error: %v", err)
	}
	if s, _ := repo.ObtenerPorID(context.Background(), activo.ID); s.Estado != domain.SprintFinalizado {
		t.Errorf("el Sprint debía quedar Finalizado, quedó %q", s.Estado)
	}
}

// esperarCierreSinAplicar verifica que un cierre fallido no dejó nada: el
// Sprint sigue Activo, sin arrastradas, y la historia sigue en el Sprint.
func esperarCierreSinAplicar(t *testing.T, db *sql.DB, sprintID int64, historia domain.BacklogItem) {
	t.Helper()
	s, err := NewSQLiteSprintRepository(db).ObtenerPorID(context.Background(), sprintID)
	if err != nil || s.Estado != domain.SprintActivo {
		t.Errorf("el Sprint debía seguir Activo, quedó %q (%v)", s.Estado, err)
	}
	if got := arrastradasDe(t, db, sprintID); len(got) != 0 {
		t.Errorf("no debía quedar ninguna historia arrastrada, quedaron %+v", got)
	}
	if actual, estado := sprintYEstadoDeHistoria(t, db, historia.ID); actual.Int64 != sprintID || estado != string(historia.Estado) {
		t.Errorf("la historia %d no debía cambiar, quedó sprint=%v estado=%q", historia.ID, actual, estado)
	}
}

func TestSQLiteSprintRepository_CerrarConArrastre_HistoriaDeOtroProyectoNoAplicaNada(t *testing.T) {
	db := nuevaBDMigrada(t)
	repo := NewSQLiteSprintRepository(db)
	proyectoID, otroProyectoID := guardarProyecto(t, db), guardarProyecto(t, db)
	activo := iniciarSprint(t, repo, guardarSprint(t, repo, proyectoID))
	propia := historiaEnSprint(t, db, proyectoID, activo.ID, domain.Estado("En progreso"), sp(5))
	ajena := historiaEnSprint(t, db, otroProyectoID, guardarSprint(t, repo, otroProyectoID).ID, domain.EstadoNueva, nil)

	err := repo.CerrarConArrastre(context.Background(), cerrado(t, activo), []domain.BacklogItem{propia, ajena})

	if !errors.Is(err, ErrIntegridadReferencial) {
		t.Fatalf("se esperaba ErrIntegridadReferencial, se obtuvo %v", err)
	}
	esperarCierreSinAplicar(t, db, activo.ID, propia)
}

func TestSQLiteSprintRepository_CerrarConArrastre_HistoriaDeOtroSprintNoAplicaNada(t *testing.T) {
	db := nuevaBDMigrada(t)
	repo := NewSQLiteSprintRepository(db)
	proyectoID := guardarProyecto(t, db)
	activo := iniciarSprint(t, repo, guardarSprint(t, repo, proyectoID))
	propia := historiaEnSprint(t, db, proyectoID, activo.ID, domain.Estado("En progreso"), sp(5))
	deOtroSprint := historiaEnSprint(t, db, proyectoID, guardarSprint(t, repo, proyectoID).ID, domain.EstadoNueva, sp(2))

	err := repo.CerrarConArrastre(context.Background(), cerrado(t, activo), []domain.BacklogItem{propia, deOtroSprint})

	if !errors.Is(err, ErrHistoriaFueraDelSprint) {
		t.Fatalf("se esperaba ErrHistoriaFueraDelSprint, se obtuvo %v", err)
	}
	esperarCierreSinAplicar(t, db, activo.ID, propia)
}

func TestSQLiteSprintRepository_CerrarConArrastre_SprintInexistente(t *testing.T) {
	repo := NewSQLiteSprintRepository(nuevaBDMigrada(t))
	s := domain.Sprint{ID: 999, ProyectoID: 1, Estado: domain.SprintFinalizado}

	err := repo.CerrarConArrastre(context.Background(), s, nil)

	if !errors.Is(err, domain.ErrSprintNoEncontrado) {
		t.Fatalf("se esperaba ErrSprintNoEncontrado, se obtuvo %v", err)
	}
}
