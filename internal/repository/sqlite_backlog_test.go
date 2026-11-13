package repository

import (
	"context"
	"errors"
	"testing"
	"time"

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

func fechaDeFinalizacion(h domain.BacklogItem) string {
	if h.FechaFinalizacion == nil {
		return "<nil>"
	}
	return h.FechaFinalizacion.Format("2006-01-02")
}

func TestSQLiteBacklogRepository_CambiarEstado_NuevaAEnProgreso(t *testing.T) {
	repo := nuevaBDDePrueba(t)
	ctx := context.Background()
	activo, _ := sprintsDePrueba(t, repo, 1)
	historia := historiaEnSprint(t, repo.db, 1, activo.ID, domain.EstadoNueva, nil)

	if err := repo.CambiarEstado(ctx, historia.ID, domain.EstadoNueva, domain.EstadoEnProgreso, nil); err != nil {
		t.Fatalf("no se esperaba error al cambiar de estado: %v", err)
	}

	leida, err := repo.ObtenerPorID(ctx, historia.ID)
	if err != nil {
		t.Fatalf("no se esperaba error al leer: %v", err)
	}
	if leida.Estado != domain.EstadoEnProgreso {
		t.Errorf("se esperaba estado En progreso, se obtuvo %q", leida.Estado)
	}
	if leida.FechaFinalizacion != nil {
		t.Errorf("se esperaba fecha de finalización nil, se obtuvo %s", fechaDeFinalizacion(leida))
	}
}

func TestSQLiteBacklogRepository_CambiarEstado_EnProgresoACompletadaGuardaLaFecha(t *testing.T) {
	repo := nuevaBDDePrueba(t)
	ctx := context.Background()
	activo, _ := sprintsDePrueba(t, repo, 1)
	historia := historiaEnSprint(t, repo.db, 1, activo.ID, domain.EstadoEnProgreso, nil)
	// La hora de entrada se descarta: solo se persiste el día.
	fecha := time.Date(2026, time.October, 29, 15, 4, 5, 0, time.UTC)

	if err := repo.CambiarEstado(ctx, historia.ID, domain.EstadoEnProgreso, domain.EstadoCompletada, &fecha); err != nil {
		t.Fatalf("no se esperaba error al completar: %v", err)
	}

	leida, err := repo.ObtenerPorID(ctx, historia.ID)
	if err != nil {
		t.Fatalf("no se esperaba error al leer: %v", err)
	}
	if leida.Estado != domain.EstadoCompletada {
		t.Errorf("se esperaba estado Completada, se obtuvo %q", leida.Estado)
	}
	if got := fechaDeFinalizacion(leida); got != "2026-10-29" {
		t.Errorf("se esperaba fecha de finalización 2026-10-29, se obtuvo %s", got)
	}
}

func TestSQLiteBacklogRepository_CambiarEstado_SegundaLlamadaConMismoDesdeEsConflicto(t *testing.T) {
	repo := nuevaBDDePrueba(t)
	ctx := context.Background()
	activo, _ := sprintsDePrueba(t, repo, 1)
	historia := historiaEnSprint(t, repo.db, 1, activo.ID, domain.EstadoNueva, nil)
	fecha := time.Date(2026, time.October, 29, 0, 0, 0, 0, time.UTC)

	if err := repo.CambiarEstado(ctx, historia.ID, domain.EstadoNueva, domain.EstadoEnProgreso, nil); err != nil {
		t.Fatalf("no se esperaba error en la primera llamada: %v", err)
	}
	err := repo.CambiarEstado(ctx, historia.ID, domain.EstadoNueva, domain.EstadoCompletada, &fecha)

	if !errors.Is(err, domain.ErrConflictoEstadoHistoria) {
		t.Fatalf("se esperaba ErrConflictoEstadoHistoria, se obtuvo %v", err)
	}
	if _, estado := sprintYEstadoDeHistoria(t, repo.db, historia.ID); estado != string(domain.EstadoEnProgreso) {
		t.Errorf("la historia debía quedar En progreso, quedó %q", estado)
	}
}

func TestSQLiteBacklogRepository_CambiarEstado_HistoriaInexistenteEsConflicto(t *testing.T) {
	repo := nuevaBDDePrueba(t)

	err := repo.CambiarEstado(context.Background(), 999, domain.EstadoNueva, domain.EstadoEnProgreso, nil)

	if !errors.Is(err, domain.ErrConflictoEstadoHistoria) {
		t.Fatalf("se esperaba ErrConflictoEstadoHistoria, se obtuvo %v", err)
	}
	var total int
	if err := repo.db.QueryRow("SELECT COUNT(*) FROM historias").Scan(&total); err != nil {
		t.Fatalf("no se pudo contar: %v", err)
	}
	if total != 0 {
		t.Errorf("no debía existir ninguna historia, hay %d", total)
	}
}

func TestSQLiteBacklogRepository_ObtenerPorID_SinFechaDevuelveNil(t *testing.T) {
	repo := nuevaBDDePrueba(t)
	ctx := context.Background()
	activo, _ := sprintsDePrueba(t, repo, 1)
	sp := 5
	historia := historiaEnSprint(t, repo.db, 1, activo.ID, domain.EstadoEnProgreso, &sp)

	leida, err := repo.ObtenerPorID(ctx, historia.ID)
	if err != nil {
		t.Fatalf("no se esperaba error al leer: %v", err)
	}

	if leida.FechaFinalizacion != nil {
		t.Errorf("se esperaba fecha de finalización nil, se obtuvo %s", fechaDeFinalizacion(leida))
	}
	if leida.Estado != domain.EstadoEnProgreso || leida.EstimacionSP == nil || *leida.EstimacionSP != 5 ||
		leida.SprintID == nil || *leida.SprintID != activo.ID {
		t.Errorf("las demás columnas no debían cambiar: %+v", leida)
	}
}

func punteroInt(v int) *int { return &v }

func rondaDePrueba(t *testing.T, repo *SQLiteBacklogRepository, historiaID int64, numero int, revelada bool) int64 {
	t.Helper()
	var reveladaEn any
	if revelada {
		reveladaEn = "2026-10-27T15:04:05Z"
	}
	resultado, err := repo.db.Exec(
		"INSERT INTO rondas_votacion (historia_id, numero, revelada_en) VALUES (?, ?, ?)", historiaID, numero, reveladaEn)
	if err != nil {
		t.Fatalf("no se pudo sembrar la ronda: %v", err)
	}
	id, err := resultado.LastInsertId()
	if err != nil {
		t.Fatalf("no se pudo obtener el id de la ronda: %v", err)
	}
	return id
}

func TestSQLiteBacklogRepository_GuardarEstimacion_EscribeSoloLosStoryPoints(t *testing.T) {
	ctx := context.Background()
	repo := nuevaBDDePrueba(t)
	historia := guardarHistoria(t, repo, 1, "Historia", domain.PrioridadMust)
	rondaID := rondaDePrueba(t, repo, historia.ID, 1, true)

	if err := repo.GuardarEstimacion(ctx, historia.ID, 8, rondaID); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	guardada, err := repo.ObtenerPorID(ctx, historia.ID)
	if err != nil {
		t.Fatalf("no se esperaba error al leer la historia: %v", err)
	}
	if guardada.EstimacionSP == nil || *guardada.EstimacionSP != 8 {
		t.Errorf("se esperaban 8 Story Points, se obtuvo %v", guardada.EstimacionSP)
	}
	if guardada.Titulo != historia.Titulo || guardada.Descripcion != historia.Descripcion || guardada.Estado != historia.Estado {
		t.Errorf("no debía cambiar nada más: antes %+v, ahora %+v", historia, guardada)
	}
}

func TestSQLiteBacklogRepository_GuardarEstimacion_HistoriaInexistente(t *testing.T) {
	repo := nuevaBDDePrueba(t)

	err := repo.GuardarEstimacion(context.Background(), 999, 8, 0)

	if !errors.Is(err, domain.ErrHistoriaNoEncontrada) {
		t.Errorf("se esperaba ErrHistoriaNoEncontrada, se obtuvo %v", err)
	}
}

func TestSQLiteBacklogRepository_GuardarEstimacion_NoTocaOtrasHistorias(t *testing.T) {
	ctx := context.Background()
	repo := nuevaBDDePrueba(t)
	historia := guardarHistoria(t, repo, 1, "Historia", domain.PrioridadMust)
	otra := guardarHistoria(t, repo, 1, "Otra", domain.PrioridadMust)
	rondaOtra := rondaDePrueba(t, repo, otra.ID, 1, true)
	rondaHistoria := rondaDePrueba(t, repo, historia.ID, 1, true)
	if err := repo.GuardarEstimacion(ctx, otra.ID, 3, rondaOtra); err != nil {
		t.Fatalf("no se esperaba error al cargar la otra historia: %v", err)
	}

	if err := repo.GuardarEstimacion(ctx, historia.ID, 8, rondaHistoria); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	guardada, err := repo.ObtenerPorID(ctx, otra.ID)
	if err != nil {
		t.Fatalf("no se esperaba error al leer la otra historia: %v", err)
	}
	if guardada.EstimacionSP == nil || *guardada.EstimacionSP != 3 {
		t.Errorf("la otra historia debía seguir con 3 Story Points, tiene %v", guardada.EstimacionSP)
	}
}

func TestSQLiteBacklogRepository_Actualizar_RechazaSiElSprintSeIniciaAntesDeEscribir(t *testing.T) {
	ctx := context.Background()
	repo := nuevaBDDePrueba(t)
	sprints := NewSQLiteSprintRepository(repo.db)
	pendiente := guardarSprint(t, sprints, 1)
	historia := guardarHistoria(t, repo, 1, "Historia", domain.PrioridadMust)
	mustExec(t, repo.db, "UPDATE historias SET sprint_id = ?, estimacion_sp = 3 WHERE id = ?", pendiente.ID, historia.ID)
	mustExec(t, repo.db,
		"UPDATE sprints SET estado = 'Activo', sprint_goal = 'g', fecha_inicio = '2026-10-01', fecha_fin = '2026-10-15' WHERE id = ?",
		pendiente.ID)
	leida := 3
	editada := historia
	editada.Titulo = "Cambiada"
	editada.EstimacionSP = punteroInt(8)

	err := repo.Actualizar(ctx, editada)

	if !errors.Is(err, domain.ErrEdicionDesactualizada) {
		t.Fatalf("se esperaba ErrEdicionDesactualizada, se obtuvo %v", err)
	}
	guardada, _ := repo.ObtenerPorID(ctx, historia.ID)
	if guardada.Titulo == "Cambiada" || guardada.EstimacionSP == nil || *guardada.EstimacionSP != leida {
		t.Errorf("la edición no debía persistirse: título %q, estimación %v", guardada.Titulo, guardada.EstimacionSP)
	}
}

func TestSQLiteBacklogRepository_Actualizar_RechazaLaCargaInicialSiElSprintSeFinalizaAntesDeEscribir(t *testing.T) {
	ctx := context.Background()
	repo := nuevaBDDePrueba(t)
	sprints := NewSQLiteSprintRepository(repo.db)
	pendiente := guardarSprint(t, sprints, 1)
	historia := guardarHistoria(t, repo, 1, "Historia", domain.PrioridadMust)
	mustExec(t, repo.db, "UPDATE historias SET sprint_id = ? WHERE id = ?", pendiente.ID, historia.ID)
	mustExec(t, repo.db,
		"UPDATE sprints SET estado = 'Finalizado', sprint_goal = 'g', fecha_inicio = '2026-10-01', fecha_fin = '2026-10-15' WHERE id = ?",
		pendiente.ID)
	editada := historia
	editada.EstimacionSP = punteroInt(5)

	err := repo.Actualizar(ctx, editada)

	if !errors.Is(err, domain.ErrEdicionDesactualizada) {
		t.Fatalf("se esperaba ErrEdicionDesactualizada, se obtuvo %v", err)
	}
	guardada, _ := repo.ObtenerPorID(ctx, historia.ID)
	if guardada.EstimacionSP != nil {
		t.Errorf("la historia debía quedar sin estimar, se obtuvo %v", guardada.EstimacionSP)
	}
}

func TestSQLiteBacklogRepository_Actualizar_PermiteReenviarLaMismaEstimacionConSprintActivo(t *testing.T) {
	ctx := context.Background()
	repo := nuevaBDDePrueba(t)
	activo, _ := sprintsDePrueba(t, repo, 1)
	historia := historiaEnSprint(t, repo.db, 1, activo.ID, domain.EstadoNueva, punteroInt(3))
	editada := historia
	editada.Titulo = "Cambiada"

	if err := repo.Actualizar(ctx, editada); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	guardada, _ := repo.ObtenerPorID(ctx, historia.ID)
	if guardada.Titulo != "Cambiada" || guardada.EstimacionSP == nil || *guardada.EstimacionSP != 3 {
		t.Errorf("se esperaba el título nuevo y 3 Story Points, se obtuvo %+v", guardada)
	}
}

func TestSQLiteBacklogRepository_GuardarEstimacion_RechazaSiElSprintSeIniciaAntesDeEscribir(t *testing.T) {
	ctx := context.Background()
	repo := nuevaBDDePrueba(t)
	sprints := NewSQLiteSprintRepository(repo.db)
	pendiente := guardarSprint(t, sprints, 1)
	historia := guardarHistoria(t, repo, 1, "Historia", domain.PrioridadMust)
	mustExec(t, repo.db, "UPDATE historias SET sprint_id = ?, estimacion_sp = 3 WHERE id = ?", pendiente.ID, historia.ID)
	rondaID := rondaDePrueba(t, repo, historia.ID, 1, true)
	mustExec(t, repo.db,
		"UPDATE sprints SET estado = 'Activo', sprint_goal = 'g', fecha_inicio = '2026-10-01', fecha_fin = '2026-10-15' WHERE id = ?",
		pendiente.ID)

	err := repo.GuardarEstimacion(ctx, historia.ID, 8, rondaID)

	if !errors.Is(err, domain.ErrRegistroEstimacionDesactualizado) {
		t.Fatalf("se esperaba ErrRegistroEstimacionDesactualizado, se obtuvo %v", err)
	}
	guardada, _ := repo.ObtenerPorID(ctx, historia.ID)
	if guardada.EstimacionSP == nil || *guardada.EstimacionSP != 3 {
		t.Errorf("la historia debía conservar 3 Story Points, se obtuvo %v", guardada.EstimacionSP)
	}
}

func TestSQLiteBacklogRepository_GuardarEstimacion_RechazaLaCargaInicialSiElSprintSeFinalizaAntesDeEscribir(t *testing.T) {
	ctx := context.Background()
	repo := nuevaBDDePrueba(t)
	sprints := NewSQLiteSprintRepository(repo.db)
	pendiente := guardarSprint(t, sprints, 1)
	historia := guardarHistoria(t, repo, 1, "Historia", domain.PrioridadMust)
	mustExec(t, repo.db, "UPDATE historias SET sprint_id = ? WHERE id = ?", pendiente.ID, historia.ID)
	rondaID := rondaDePrueba(t, repo, historia.ID, 1, true)
	mustExec(t, repo.db,
		"UPDATE sprints SET estado = 'Finalizado', sprint_goal = 'g', fecha_inicio = '2026-10-01', fecha_fin = '2026-10-15' WHERE id = ?",
		pendiente.ID)

	err := repo.GuardarEstimacion(ctx, historia.ID, 5, rondaID)

	if !errors.Is(err, domain.ErrRegistroEstimacionDesactualizado) {
		t.Fatalf("se esperaba ErrRegistroEstimacionDesactualizado, se obtuvo %v", err)
	}
	guardada, _ := repo.ObtenerPorID(ctx, historia.ID)
	if guardada.EstimacionSP != nil {
		t.Errorf("la historia debía quedar sin estimar, se obtuvo %v", guardada.EstimacionSP)
	}
}

func TestSQLiteBacklogRepository_GuardarEstimacion_RechazaSiSeAbreOtraRonda(t *testing.T) {
	ctx := context.Background()
	repo := nuevaBDDePrueba(t)
	historia := guardarHistoria(t, repo, 1, "Historia", domain.PrioridadMust)
	mustExec(t, repo.db, "UPDATE historias SET estimacion_sp = 3 WHERE id = ?", historia.ID)
	rondaUno := rondaDePrueba(t, repo, historia.ID, 1, true)
	rondaDePrueba(t, repo, historia.ID, 2, false)

	err := repo.GuardarEstimacion(ctx, historia.ID, 8, rondaUno)

	if !errors.Is(err, domain.ErrRegistroEstimacionDesactualizado) {
		t.Fatalf("se esperaba ErrRegistroEstimacionDesactualizado, se obtuvo %v", err)
	}
	guardada, _ := repo.ObtenerPorID(ctx, historia.ID)
	if guardada.EstimacionSP == nil || *guardada.EstimacionSP != 3 {
		t.Errorf("la historia debía conservar 3 Story Points, se obtuvo %v", guardada.EstimacionSP)
	}
}

func TestSQLiteBacklogRepository_GuardarEstimacion_PermiteLaCargaInicialConSprintActivo(t *testing.T) {
	ctx := context.Background()
	repo := nuevaBDDePrueba(t)
	activo, _ := sprintsDePrueba(t, repo, 1)
	historia := historiaEnSprint(t, repo.db, 1, activo.ID, domain.EstadoNueva, nil)
	rondaID := rondaDePrueba(t, repo, historia.ID, 1, true)

	if err := repo.GuardarEstimacion(ctx, historia.ID, 5, rondaID); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	guardada, _ := repo.ObtenerPorID(ctx, historia.ID)
	if guardada.EstimacionSP == nil || *guardada.EstimacionSP != 5 {
		t.Errorf("se esperaban 5 Story Points, se obtuvo %v", guardada.EstimacionSP)
	}
}

func TestSQLiteBacklogRepository_ListarPorProyecto(t *testing.T) {
	t.Run("devuelve todas las historias del proyecto, con y sin Sprint, ordenadas por id", func(t *testing.T) {
		repo := nuevaBDDePrueba(t)
		activo, pendiente := sprintsDePrueba(t, repo, 1)
		cinco, ocho := 5, 8

		sinSprint := guardarHistoria(t, repo, 1, "Sin Sprint", domain.PrioridadWont)
		nuevaAsignada := historiaEnSprint(t, repo.db, 1, pendiente.ID, domain.EstadoNueva, nil)
		enProgreso := historiaEnSprint(t, repo.db, 1, activo.ID, domain.EstadoEnProgreso, &cinco)
		completada := historiaEnSprint(t, repo.db, 1, activo.ID, domain.EstadoCompletada, &ocho)

		historias, err := repo.ListarPorProyecto(context.Background(), 1)
		if err != nil {
			t.Fatalf("no se esperaba error: %v", err)
		}

		esperadas := []domain.BacklogItem{sinSprint, nuevaAsignada, enProgreso, completada}
		if len(historias) != len(esperadas) {
			t.Fatalf("se esperaban %d historias, se obtuvieron %d", len(esperadas), len(historias))
		}
		for i, esperada := range esperadas {
			if historias[i].ID != esperada.ID {
				t.Fatalf("posición %d: se esperaba la historia %d, se obtuvo %d", i, esperada.ID, historias[i].ID)
			}
		}
		if historias[0].SprintID != nil {
			t.Errorf("la historia sin Sprint debe conservar SprintID nil, se obtuvo %d", *historias[0].SprintID)
		}
		if historias[1].SprintID == nil || *historias[1].SprintID != pendiente.ID || historias[1].Estado != domain.EstadoNueva {
			t.Errorf("la Nueva asignada debe conservar su Sprint pendiente %d, se obtuvo %+v", pendiente.ID, historias[1])
		}
		if historias[2].Estado != domain.EstadoEnProgreso || historias[3].Estado != domain.EstadoCompletada {
			t.Errorf("se esperaban los estados En progreso y Completada, se obtuvo %q y %q", historias[2].Estado, historias[3].Estado)
		}
	})

	t.Run("excluye las historias de otros proyectos", func(t *testing.T) {
		repo := nuevaBDDePrueba(t)
		otroProyecto := guardarProyecto(t, repo.db)
		propia := guardarHistoria(t, repo, 1, "Propia", domain.PrioridadMust)
		guardarHistoria(t, repo, otroProyecto, "Ajena", domain.PrioridadMust)

		historias, err := repo.ListarPorProyecto(context.Background(), 1)
		if err != nil {
			t.Fatalf("no se esperaba error: %v", err)
		}
		if len(historias) != 1 || historias[0].ID != propia.ID {
			t.Errorf("se esperaba solo la historia propia %d, se obtuvo %+v", propia.ID, historias)
		}
	})

	t.Run("conserva la estimación en nil cuando la historia no está estimada", func(t *testing.T) {
		repo := nuevaBDDePrueba(t)
		activo, _ := sprintsDePrueba(t, repo, 1)
		tres := 3
		historiaEnSprint(t, repo.db, 1, activo.ID, domain.EstadoEnProgreso, nil)
		historiaEnSprint(t, repo.db, 1, activo.ID, domain.EstadoEnProgreso, &tres)

		historias, err := repo.ListarPorProyecto(context.Background(), 1)
		if err != nil {
			t.Fatalf("no se esperaba error: %v", err)
		}
		if len(historias) != 2 {
			t.Fatalf("se esperaban 2 historias, se obtuvieron %d", len(historias))
		}
		if historias[0].EstimacionSP != nil {
			t.Errorf("la historia sin estimar debe tener EstimacionSP nil, se obtuvo %d", *historias[0].EstimacionSP)
		}
		if historias[1].EstimacionSP == nil || *historias[1].EstimacionSP != 3 {
			t.Errorf("se esperaban 3 Story Points, se obtuvo %v", historias[1].EstimacionSP)
		}
	})

	t.Run("un proyecto sin historias devuelve un slice vacío, nunca nil", func(t *testing.T) {
		repo := nuevaBDDePrueba(t)

		historias, err := repo.ListarPorProyecto(context.Background(), 1)
		if err != nil {
			t.Fatalf("no se esperaba error: %v", err)
		}
		if historias == nil || len(historias) != 0 {
			t.Errorf("se esperaba un slice vacío no nil, se obtuvo %#v", historias)
		}
	})
}
