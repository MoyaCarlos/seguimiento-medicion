package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
)

func nuevaBDDePrueba(t *testing.T) *SQLiteBacklogRepository {
	t.Helper()
	db := nuevaBDMigrada(t)
	guardarProyecto(t, db)
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
		"SELECT estado, estimacion_sp, valor_negocio FROM historias WHERE id = ?", guardado.ID,
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
	if err := repo.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM historias WHERE titulo = ?", "Mismo título").Scan(&total); err != nil {
		t.Fatalf("no se pudo contar: %v", err)
	}
	if total != 2 {
		t.Errorf("se esperaban 2 historias con el mismo título, se obtuvieron %d", total)
	}
}

func TestSQLiteBacklogRepository_Guardar_ProyectoInexistenteEsIntegridad(t *testing.T) {
	repo := nuevaBDDePrueba(t)
	item, err := domain.NewBacklogItem(999, "Historia huérfana", "Sin proyecto", domain.PrioridadMust, nil)
	if err != nil {
		t.Fatalf("no se esperaba error de dominio: %v", err)
	}

	_, err = repo.Guardar(context.Background(), item)

	if !errors.Is(err, ErrIntegridadReferencial) {
		t.Fatalf("se esperaba ErrIntegridadReferencial, se obtuvo %v", err)
	}
}

const estadoEnProgreso domain.Estado = "En progreso"

func guardarHistoria(t *testing.T, repo *SQLiteBacklogRepository, proyectoID int64, titulo string, prioridad domain.Prioridad) domain.BacklogItem {
	t.Helper()
	item, err := domain.NewBacklogItem(proyectoID, titulo, "Descripción de "+titulo, prioridad, nil)
	if err != nil {
		t.Fatalf("no se esperaba error de dominio: %v", err)
	}
	guardado, err := repo.Guardar(context.Background(), item)
	if err != nil {
		t.Fatalf("no se esperaba error al guardar: %v", err)
	}
	return guardado
}

func asignarASprint(t *testing.T, repo *SQLiteBacklogRepository, historiaID, sprintID int64) {
	t.Helper()
	if _, err := repo.db.ExecContext(context.Background(),
		"UPDATE historias SET sprint_id = ?, estado = ? WHERE id = ?", sprintID, string(estadoEnProgreso), historiaID,
	); err != nil {
		t.Fatalf("no se pudo asignar la historia al sprint: %v", err)
	}
}

func titulos(historias []domain.BacklogItem) []string {
	resultado := make([]string, 0, len(historias))
	for _, h := range historias {
		resultado = append(resultado, h.Titulo)
	}
	return resultado
}

func TestSQLiteBacklogRepository_ListarSinSprint_ExcluyeConSprintYOtrosProyectos(t *testing.T) {
	repo := nuevaBDDePrueba(t)
	otroProyecto := guardarProyecto(t, repo.db)
	sprint := guardarSprint(t, NewSQLiteSprintRepository(repo.db), 1)

	guardarHistoria(t, repo, 1, "Sin sprint", domain.PrioridadMust)
	conSprint := guardarHistoria(t, repo, 1, "Con sprint", domain.PrioridadMust)
	asignarASprint(t, repo, conSprint.ID, sprint.ID)
	guardarHistoria(t, repo, otroProyecto, "De otro proyecto", domain.PrioridadMust)

	historias, err := repo.ListarSinSprint(context.Background(), 1)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if got := titulos(historias); len(got) != 1 || got[0] != "Sin sprint" {
		t.Errorf("se esperaba solo [Sin sprint], se obtuvo %v", got)
	}
}

func TestSQLiteBacklogRepository_ListarSinSprint_OrdenMoSCoWYCreacion(t *testing.T) {
	repo := nuevaBDDePrueba(t)
	guardarHistoria(t, repo, 1, "S1", domain.PrioridadShould)
	guardarHistoria(t, repo, 1, "M1", domain.PrioridadMust)
	guardarHistoria(t, repo, 1, "W1", domain.PrioridadWont)
	guardarHistoria(t, repo, 1, "C1", domain.PrioridadCould)
	guardarHistoria(t, repo, 1, "M2", domain.PrioridadMust)

	historias, err := repo.ListarSinSprint(context.Background(), 1)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	esperado := []string{"M1", "M2", "S1", "C1", "W1"}
	got := titulos(historias)
	if len(got) != len(esperado) {
		t.Fatalf("se esperaba %v, se obtuvo %v", esperado, got)
	}
	for i := range esperado {
		if got[i] != esperado[i] {
			t.Fatalf("se esperaba el orden %v, se obtuvo %v", esperado, got)
		}
	}
}

func TestSQLiteBacklogRepository_ListarSinSprint_SinHistoriasDevuelveVacio(t *testing.T) {
	repo := nuevaBDDePrueba(t)

	historias, err := repo.ListarSinSprint(context.Background(), 1)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if historias == nil || len(historias) != 0 {
		t.Errorf("se esperaba un slice vacío no nil, se obtuvo %#v", historias)
	}
}

func TestSQLiteBacklogRepository_ListarPorSprint_OrdenMoSCoWYCreacion(t *testing.T) {
	t.Run("ordena por prioridad y luego por creación, y excluye otros Sprints", func(t *testing.T) {
		repo := nuevaBDDePrueba(t)
		ctx := context.Background()
		sprints := NewSQLiteSprintRepository(repo.db)
		activo := iniciarSprint(t, sprints, guardarSprint(t, sprints, 1))
		pendiente := guardarSprint(t, sprints, 1)

		for _, h := range []struct {
			titulo    string
			prioridad domain.Prioridad
		}{
			{"S1", domain.PrioridadShould},
			{"M1", domain.PrioridadMust},
			{"W1", domain.PrioridadWont},
			{"C1", domain.PrioridadCould},
			{"M2", domain.PrioridadMust},
		} {
			guardada := guardarHistoria(t, repo, 1, h.titulo, h.prioridad)
			if err := repo.AsignarASprint(ctx, guardada.ID, pendiente.ID); err != nil {
				t.Fatalf("no se pudo asignar %q al Sprint Pendiente: %v", h.titulo, err)
			}
		}
		ajena := guardarHistoria(t, repo, 1, "De otro Sprint", domain.PrioridadMust)
		if err := repo.AsignarASprint(ctx, ajena.ID, activo.ID); err != nil {
			t.Fatalf("no se pudo asignar la historia del otro Sprint: %v", err)
		}

		historias, err := repo.ListarPorSprint(ctx, pendiente.ID)
		if err != nil {
			t.Fatalf("no se esperaba error: %v", err)
		}

		esperado := []string{"M1", "M2", "S1", "C1", "W1"}
		got := titulos(historias)
		if len(got) != len(esperado) {
			t.Fatalf("se esperaba el orden %v, se obtuvo %v", esperado, got)
		}
		for i := range esperado {
			if got[i] != esperado[i] {
				t.Fatalf("se esperaba el orden %v, se obtuvo %v", esperado, got)
			}
		}
	})

	t.Run("Sprint sin historias devuelve un slice vacío no nil", func(t *testing.T) {
		repo := nuevaBDDePrueba(t)
		sprints := NewSQLiteSprintRepository(repo.db)
		otroProyecto := guardarProyecto(t, repo.db)
		vacio := guardarSprint(t, sprints, otroProyecto)

		historias, err := repo.ListarPorSprint(context.Background(), vacio.ID)
		if err != nil {
			t.Fatalf("no se esperaba error: %v", err)
		}
		if historias == nil || len(historias) != 0 {
			t.Errorf("se esperaba un slice vacío no nil, se obtuvo %#v", historias)
		}
	})

	t.Run("Sprint cerrado con arrastre solo conserva la historia completada", func(t *testing.T) {
		repo := nuevaBDDePrueba(t)
		ctx := context.Background()
		sprints := NewSQLiteSprintRepository(repo.db)
		activo := iniciarSprint(t, sprints, guardarSprint(t, sprints, 1))
		completada := historiaEnSprint(t, repo.db, 1, activo.ID, domain.EstadoCompletada, nil)
		noCompletada := historiaEnSprint(t, repo.db, 1, activo.ID, domain.EstadoNueva, nil)

		if err := sprints.CerrarConArrastre(ctx, cerrado(t, activo), []domain.BacklogItem{noCompletada}); err != nil {
			t.Fatalf("no se esperaba error al cerrar el Sprint: %v", err)
		}

		historias, err := repo.ListarPorSprint(ctx, activo.ID)
		if err != nil {
			t.Fatalf("no se esperaba error: %v", err)
		}
		if len(historias) != 1 || historias[0].ID != completada.ID {
			t.Errorf("se esperaba solo la historia completada %d, se obtuvo %v", completada.ID, titulos(historias))
		}
	})
}

func TestSQLiteBacklogRepository_ObtenerPorID_Inexistente(t *testing.T) {
	repo := nuevaBDDePrueba(t)

	_, err := repo.ObtenerPorID(context.Background(), 999)

	if !errors.Is(err, domain.ErrHistoriaNoEncontrada) {
		t.Fatalf("se esperaba ErrHistoriaNoEncontrada, se obtuvo %v", err)
	}
}

func TestSQLiteBacklogRepository_ObtenerPorID_LeeTodosLosCampos(t *testing.T) {
	repo := nuevaBDDePrueba(t)
	ctx := context.Background()
	sprint := guardarSprint(t, NewSQLiteSprintRepository(repo.db), 1)
	valor := 8
	item, _ := domain.NewBacklogItem(1, "Historia", "Descripción", domain.PrioridadShould, &valor)
	guardado, err := repo.Guardar(ctx, item)
	if err != nil {
		t.Fatalf("no se esperaba error al guardar: %v", err)
	}
	asignarASprint(t, repo, guardado.ID, sprint.ID)
	if _, err := repo.db.ExecContext(ctx,
		"UPDATE historias SET criterios_aceptacion = ?, estimacion_sp = 5 WHERE id = ?", "Dado, cuando, entonces", guardado.ID,
	); err != nil {
		t.Fatalf("no se pudo sembrar la historia: %v", err)
	}

	h, err := repo.ObtenerPorID(ctx, guardado.ID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if h.ID != guardado.ID || h.ProyectoID != 1 || h.Titulo != "Historia" || h.Descripcion != "Descripción" {
		t.Errorf("datos básicos inesperados: %+v", h)
	}
	if h.Prioridad != domain.PrioridadShould || h.Estado != estadoEnProgreso {
		t.Errorf("prioridad o estado inesperados: %q %q", h.Prioridad, h.Estado)
	}
	if h.CriteriosAceptacion == nil || *h.CriteriosAceptacion != "Dado, cuando, entonces" {
		t.Errorf("criterios inesperados: %v", h.CriteriosAceptacion)
	}
	if h.ValorNegocio == nil || *h.ValorNegocio != 8 {
		t.Errorf("valor de negocio inesperado: %v", h.ValorNegocio)
	}
	if h.EstimacionSP == nil || *h.EstimacionSP != 5 {
		t.Errorf("estimación inesperada: %v", h.EstimacionSP)
	}
	if h.SprintID == nil || *h.SprintID != sprint.ID {
		t.Errorf("sprint inesperado: %v", h.SprintID)
	}
}

func TestSQLiteBacklogRepository_Actualizar(t *testing.T) {
	repo := nuevaBDDePrueba(t)
	ctx := context.Background()
	sprint := guardarSprint(t, NewSQLiteSprintRepository(repo.db), 1)
	valor := 13
	item, _ := domain.NewBacklogItem(1, "Original", "Descripción original", domain.PrioridadCould, &valor)
	guardado, err := repo.Guardar(ctx, item)
	if err != nil {
		t.Fatalf("no se esperaba error al guardar: %v", err)
	}
	asignarASprint(t, repo, guardado.ID, sprint.ID)

	criterios, sp := "Criterios nuevos", 3
	editado := guardado
	editado.Titulo, editado.Descripcion = "Nuevo", "Descripción nueva"
	editado.CriteriosAceptacion, editado.EstimacionSP = &criterios, &sp
	if err := repo.Actualizar(ctx, editado); err != nil {
		t.Fatalf("no se esperaba error al actualizar: %v", err)
	}

	h, err := repo.ObtenerPorID(ctx, guardado.ID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if h.Titulo != "Nuevo" || h.Descripcion != "Descripción nueva" {
		t.Errorf("título o descripción sin actualizar: %+v", h)
	}
	if h.CriteriosAceptacion == nil || *h.CriteriosAceptacion != "Criterios nuevos" {
		t.Errorf("criterios sin actualizar: %v", h.CriteriosAceptacion)
	}
	if h.EstimacionSP == nil || *h.EstimacionSP != 3 {
		t.Errorf("estimación sin actualizar: %v", h.EstimacionSP)
	}
	if h.Prioridad != domain.PrioridadCould || h.Estado != estadoEnProgreso ||
		h.ValorNegocio == nil || *h.ValorNegocio != 13 || h.SprintID == nil || *h.SprintID != sprint.ID {
		t.Errorf("se modificaron campos que no se editan: %+v", h)
	}
}

func TestSQLiteBacklogRepository_Actualizar_VuelveANull(t *testing.T) {
	repo := nuevaBDDePrueba(t)
	ctx := context.Background()
	guardado := guardarHistoria(t, repo, 1, "Historia", domain.PrioridadMust)
	if _, err := repo.db.ExecContext(ctx,
		"UPDATE historias SET criterios_aceptacion = ?, estimacion_sp = 5 WHERE id = ?", "Algo", guardado.ID,
	); err != nil {
		t.Fatalf("no se pudo sembrar la historia: %v", err)
	}

	guardado.CriteriosAceptacion, guardado.EstimacionSP = nil, nil
	if err := repo.Actualizar(ctx, guardado); err != nil {
		t.Fatalf("no se esperaba error al actualizar: %v", err)
	}

	var (
		criterios  *string
		estimacion *int
	)
	if err := repo.db.QueryRowContext(ctx,
		"SELECT criterios_aceptacion, estimacion_sp FROM historias WHERE id = ?", guardado.ID,
	).Scan(&criterios, &estimacion); err != nil {
		t.Fatalf("no se pudo releer la historia: %v", err)
	}
	if criterios != nil || estimacion != nil {
		t.Errorf("se esperaban criterios y estimación NULL, se obtuvo %v %v", criterios, estimacion)
	}
}

func TestSQLiteBacklogRepository_Guardar_PersisteCriterios(t *testing.T) {
	repo := nuevaBDDePrueba(t)
	ctx := context.Background()
	item, _ := domain.NewBacklogItem(1, "Historia", "Descripción", domain.PrioridadMust, nil)
	criterios := "Dado, cuando, entonces"
	item.CriteriosAceptacion = &criterios

	guardado, err := repo.Guardar(ctx, item)
	if err != nil {
		t.Fatalf("no se esperaba error al guardar: %v", err)
	}

	h, err := repo.ObtenerPorID(ctx, guardado.ID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if h.CriteriosAceptacion == nil || *h.CriteriosAceptacion != criterios {
		t.Errorf("se esperaban los criterios guardados, se obtuvo %v", h.CriteriosAceptacion)
	}
}

func TestSQLiteBacklogRepository_Guardar_CriteriosNull(t *testing.T) {
	repo := nuevaBDDePrueba(t)
	guardado := guardarHistoria(t, repo, 1, "Historia", domain.PrioridadMust)

	var criterios *string
	if err := repo.db.QueryRowContext(context.Background(),
		"SELECT criterios_aceptacion FROM historias WHERE id = ?", guardado.ID,
	).Scan(&criterios); err != nil {
		t.Fatalf("no se pudo releer la historia: %v", err)
	}
	if criterios != nil {
		t.Errorf("se esperaban criterios NULL, se obtuvo %q", *criterios)
	}
}

func sprintsDePrueba(t *testing.T, repo *SQLiteBacklogRepository, proyectoID int64) (activo, pendiente domain.Sprint) {
	t.Helper()
	sprints := NewSQLiteSprintRepository(repo.db)
	activo = iniciarSprint(t, sprints, guardarSprint(t, sprints, proyectoID))
	pendiente = guardarSprint(t, sprints, proyectoID)
	return activo, pendiente
}

func TestSQLiteBacklogRepository_AsignarASprint_AsignaYConservaDatos(t *testing.T) {
	repo := nuevaBDDePrueba(t)
	ctx := context.Background()
	_, pendiente := sprintsDePrueba(t, repo, 1)
	historia := guardarHistoria(t, repo, 1, "Iniciar sesión", domain.PrioridadMust)
	mustExec(t, repo.db, "UPDATE historias SET estimacion_sp = 5 WHERE id = ?", historia.ID)

	if err := repo.AsignarASprint(ctx, historia.ID, pendiente.ID); err != nil {
		t.Fatalf("no se esperaba error al asignar: %v", err)
	}

	leida, err := repo.ObtenerPorID(ctx, historia.ID)
	if err != nil {
		t.Fatalf("no se esperaba error al leer: %v", err)
	}
	if leida.SprintID == nil || *leida.SprintID != pendiente.ID {
		t.Errorf("se esperaba el Sprint %d, se obtuvo %v", pendiente.ID, leida.SprintID)
	}
	if leida.Estado != domain.EstadoNueva {
		t.Errorf("se esperaba estado Nueva, se obtuvo %q", leida.Estado)
	}
	if leida.EstimacionSP == nil || *leida.EstimacionSP != 5 {
		t.Errorf("se esperaban 5 Story Points, se obtuvo %v", leida.EstimacionSP)
	}
}

func TestSQLiteBacklogRepository_AsignarASprint_SegundaAsignacionYaAsignada(t *testing.T) {
	repo := nuevaBDDePrueba(t)
	ctx := context.Background()
	activo, pendiente := sprintsDePrueba(t, repo, 1)
	historia := guardarHistoria(t, repo, 1, "Iniciar sesión", domain.PrioridadMust)

	if err := repo.AsignarASprint(ctx, historia.ID, activo.ID); err != nil {
		t.Fatalf("no se esperaba error en la primera asignación: %v", err)
	}
	err := repo.AsignarASprint(ctx, historia.ID, pendiente.ID)

	if !errors.Is(err, domain.ErrHistoriaYaAsignada) {
		t.Fatalf("se esperaba ErrHistoriaYaAsignada, se obtuvo %v", err)
	}
	if sprintID, _ := sprintYEstadoDeHistoria(t, repo.db, historia.ID); sprintID.Int64 != activo.ID {
		t.Errorf("la historia debía seguir en el Sprint %d, quedó %v", activo.ID, sprintID)
	}
}

func TestSQLiteBacklogRepository_AsignarASprint_SprintDeOtroProyectoEsIntegridad(t *testing.T) {
	repo := nuevaBDDePrueba(t)
	otroProyecto := guardarProyecto(t, repo.db)
	_, ajeno := sprintsDePrueba(t, repo, otroProyecto)
	historia := guardarHistoria(t, repo, 1, "Iniciar sesión", domain.PrioridadMust)

	err := repo.AsignarASprint(context.Background(), historia.ID, ajeno.ID)

	if !errors.Is(err, ErrIntegridadReferencial) {
		t.Fatalf("se esperaba ErrIntegridadReferencial, se obtuvo %v", err)
	}
	if sprintID, _ := sprintYEstadoDeHistoria(t, repo.db, historia.ID); sprintID.Valid {
		t.Errorf("la historia debía seguir sin Sprint, quedó %v", sprintID)
	}
}

func TestSQLiteBacklogRepository_AsignarASprint_SprintInexistenteEsIntegridad(t *testing.T) {
	repo := nuevaBDDePrueba(t)
	historia := guardarHistoria(t, repo, 1, "Iniciar sesión", domain.PrioridadMust)

	err := repo.AsignarASprint(context.Background(), historia.ID, 999)

	if !errors.Is(err, ErrIntegridadReferencial) {
		t.Fatalf("se esperaba ErrIntegridadReferencial, se obtuvo %v", err)
	}
}

func TestSQLiteBacklogRepository_AsignarASprint_HistoriaArrastradaSeReasigna(t *testing.T) {
	repo := nuevaBDDePrueba(t)
	ctx := context.Background()
	sprints := NewSQLiteSprintRepository(repo.db)
	activo, pendiente := sprintsDePrueba(t, repo, 1)
	arrastrada := historiaEnSprint(t, repo.db, 1, activo.ID, domain.EstadoNueva, nil)
	if err := sprints.CerrarConArrastre(ctx, cerrado(t, activo), []domain.BacklogItem{arrastrada}); err != nil {
		t.Fatalf("no se esperaba error al cerrar: %v", err)
	}

	if err := repo.AsignarASprint(ctx, arrastrada.ID, pendiente.ID); err != nil {
		t.Fatalf("no se esperaba error al reasignar: %v", err)
	}

	if sprintID, _ := sprintYEstadoDeHistoria(t, repo.db, arrastrada.ID); sprintID.Int64 != pendiente.ID {
		t.Errorf("la historia debía quedar en el Sprint %d, quedó %v", pendiente.ID, sprintID)
	}
}
