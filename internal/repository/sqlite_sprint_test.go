package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func nuevaBD(t *testing.T) *sql.DB {
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
	return db
}

func guardarSprint(t *testing.T, repo *SQLiteSprintRepository, proyectoID int64) domain.Sprint {
	t.Helper()
	s, err := domain.NewSprint(proyectoID)
	if err != nil {
		t.Fatalf("no se esperaba error de dominio: %v", err)
	}
	guardado, err := repo.Guardar(context.Background(), s)
	if err != nil {
		t.Fatalf("no se esperaba error al guardar: %v", err)
	}
	return guardado
}

func guardarProyecto(t *testing.T, db *sql.DB) int64 {
	t.Helper()
	p, err := domain.NewProject("Proyecto", "", nil, nil)
	if err != nil {
		t.Fatalf("no se esperaba error de dominio: %v", err)
	}
	guardado, err := NewSQLiteProjectRepository(db).Guardar(context.Background(), p)
	if err != nil {
		t.Fatalf("no se esperaba error al guardar el proyecto: %v", err)
	}
	return guardado.ID
}

func iniciarSprint(t *testing.T, repo *SQLiteSprintRepository, s domain.Sprint) domain.Sprint {
	t.Helper()
	inicio := time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC)
	if err := s.Iniciar("Objetivo", inicio, inicio.AddDate(0, 0, 14)); err != nil {
		t.Fatalf("no se esperaba error de dominio: %v", err)
	}
	if err := repo.Actualizar(context.Background(), s); err != nil {
		t.Fatalf("no se esperaba error al iniciar el Sprint: %v", err)
	}
	return s
}

func TestMigrar_EsIdempotente(t *testing.T) {
	db := nuevaBD(t)
	if err := Migrar(context.Background(), db); err != nil {
		t.Fatalf("migrar dos veces no debería fallar: %v", err)
	}
}

func TestSQLiteSprintRepository_GuardarYObtenerPendiente(t *testing.T) {
	db := nuevaBDMigrada(t)
	repo := NewSQLiteSprintRepository(db)
	proyectoID := guardarProyecto(t, db)

	guardado := guardarSprint(t, repo, proyectoID)
	if guardado.ID <= 0 {
		t.Fatalf("se esperaba un ID asignado, se obtuvo %d", guardado.ID)
	}

	leido, err := repo.ObtenerPorID(context.Background(), guardado.ID)
	if err != nil {
		t.Fatalf("no se esperaba error al leer: %v", err)
	}
	if leido.Estado != domain.SprintPendiente || leido.ProyectoID != proyectoID {
		t.Errorf("se leyó un Sprint inesperado: %+v", leido)
	}
	if leido.SprintGoal != "" || !leido.FechaInicio.IsZero() || !leido.FechaFin.IsZero() {
		t.Errorf("un Sprint pendiente no debe tener Goal ni fechas: %+v", leido)
	}
}

func TestSQLiteSprintRepository_ObtenerInexistente(t *testing.T) {
	repo := NewSQLiteSprintRepository(nuevaBDMigrada(t))

	_, err := repo.ObtenerPorID(context.Background(), 999)
	if !errors.Is(err, domain.ErrSprintNoEncontrado) {
		t.Fatalf("se esperaba ErrSprintNoEncontrado, se obtuvo %v", err)
	}
}

func TestSQLiteSprintRepository_ActualizarPersisteLaTransicion(t *testing.T) {
	db := nuevaBDMigrada(t)
	repo := NewSQLiteSprintRepository(db)
	ctx := context.Background()
	sprint := guardarSprint(t, repo, guardarProyecto(t, db))

	inicio := time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC)
	fin := inicio.AddDate(0, 0, 14)
	if err := sprint.Iniciar("Entregar el MVP", inicio, fin); err != nil {
		t.Fatalf("no se esperaba error de dominio: %v", err)
	}
	if err := repo.Actualizar(ctx, sprint); err != nil {
		t.Fatalf("no se esperaba error al actualizar: %v", err)
	}

	leido, err := repo.ObtenerPorID(ctx, sprint.ID)
	if err != nil {
		t.Fatalf("no se esperaba error al leer: %v", err)
	}
	if leido.Estado != domain.SprintActivo || leido.SprintGoal != "Entregar el MVP" {
		t.Errorf("no se persistió la transición: %+v", leido)
	}
	if !leido.FechaInicio.Equal(inicio) || !leido.FechaFin.Equal(fin) {
		t.Errorf("no se persistieron las fechas: %v - %v", leido.FechaInicio, leido.FechaFin)
	}
}

func TestSQLiteSprintRepository_ListarPorProyecto(t *testing.T) {
	db := nuevaBDMigrada(t)
	repo := NewSQLiteSprintRepository(db)
	proyecto1, proyecto2 := guardarProyecto(t, db), guardarProyecto(t, db)
	primero := iniciarSprint(t, repo, guardarSprint(t, repo, proyecto1))
	guardarSprint(t, repo, proyecto2)
	segundo := guardarSprint(t, repo, proyecto1)

	sprints, err := repo.ListarPorProyecto(context.Background(), proyecto1)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(sprints) != 2 || sprints[0].ID != primero.ID || sprints[1].ID != segundo.ID {
		t.Errorf("se esperaban los 2 Sprints del proyecto 1 en orden, se obtuvo %+v", sprints)
	}
}

func TestSQLiteBacklogRepository_HistoriasDeSprint(t *testing.T) {
	db := nuevaBDMigrada(t)
	repo := NewSQLiteBacklogRepository(db)
	sprints := NewSQLiteSprintRepository(db)
	ctx := context.Background()
	proyectoID := guardarProyecto(t, db)
	activo := iniciarSprint(t, sprints, guardarSprint(t, sprints, proyectoID))
	pendiente := guardarSprint(t, sprints, proyectoID)

	var ids []int64
	for _, titulo := range []string{"Completada", "Pendiente", "Otro sprint"} {
		item, _ := domain.NewBacklogItem(proyectoID, titulo, "Descripción", domain.PrioridadMust, nil)
		g, err := repo.Guardar(ctx, item)
		if err != nil {
			t.Fatalf("no se esperaba error al guardar: %v", err)
		}
		ids = append(ids, g.ID)
	}
	// La asignación de historias a un Sprint es de HU-06; acá se simula con SQL directo.
	mustExec(t, db, "UPDATE historias SET sprint_id = ? WHERE id IN (?, ?)", activo.ID, ids[0], ids[1])
	mustExec(t, db, "UPDATE historias SET sprint_id = ? WHERE id = ?", pendiente.ID, ids[2])
	mustExec(t, db, "UPDATE historias SET estado = ?, fecha_finalizacion = '2026-10-20' WHERE id = ?", string(domain.EstadoCompletada), ids[0])
	mustExec(t, db, "UPDATE historias SET estado = ? WHERE id = ?", "En progreso", ids[1])

	historias, err := repo.ListarPorSprint(ctx, activo.ID)
	if err != nil {
		t.Fatalf("no se esperaba error al listar: %v", err)
	}
	if len(historias) != 2 {
		t.Fatalf("se esperaban 2 historias del Sprint activo, se obtuvieron %d", len(historias))
	}
	if historias[0].Estado != domain.EstadoCompletada || historias[0].SprintID == nil || *historias[0].SprintID != activo.ID {
		t.Errorf("no se leyeron bien estado/sprint de la historia: %+v", historias[0])
	}

	if err := repo.QuitarDeSprint(ctx, ids[1]); err != nil {
		t.Fatalf("no se esperaba error al quitar: %v", err)
	}
	historias, _ = repo.ListarPorSprint(ctx, activo.ID)
	if len(historias) != 1 || historias[0].ID != ids[0] {
		t.Errorf("se esperaba que solo quede la historia completada en el Sprint activo: %+v", historias)
	}
	var estado string
	if err := db.QueryRow("SELECT estado FROM historias WHERE id = ?", ids[1]).Scan(&estado); err != nil {
		t.Fatalf("no se pudo consultar el estado de la historia devuelta al backlog: %v", err)
	}
	if estado != string(domain.EstadoNueva) {
		t.Errorf("la historia devuelta al backlog debía quedar en estado Nueva, quedó %q", estado)
	}
}

func mustExec(t *testing.T, db *sql.DB, query string, args ...any) {
	t.Helper()
	if _, err := db.Exec(query, args...); err != nil {
		t.Fatalf("no se pudo ejecutar %q: %v", query, err)
	}
}
