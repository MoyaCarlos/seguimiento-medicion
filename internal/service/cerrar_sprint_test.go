package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// historiasEnMemoria simula las historias asignadas a Sprints.
type historiasEnMemoria struct {
	items []domain.BacklogItem
}

func (h *historiasEnMemoria) ListarPorSprint(_ context.Context, sprintID int64) ([]domain.BacklogItem, error) {
	var resultado []domain.BacklogItem
	for _, item := range h.items {
		if item.SprintID != nil && *item.SprintID == sprintID {
			resultado = append(resultado, item)
		}
	}
	return resultado, nil
}

func (h *historiasEnMemoria) QuitarDeSprint(context.Context, int64) error { return nil }

func historia(id int64, estado domain.Estado, sprintID int64, estimacion *int) domain.BacklogItem {
	return domain.BacklogItem{ID: id, ProyectoID: 1, Estado: estado, EstimacionSP: estimacion, SprintID: &sprintID}
}

func puntos(v int) *int { return &v }

func sembrarSprintActivo(t *testing.T, repo repository.SprintRepository, proyectoID int64) domain.Sprint {
	t.Helper()
	pendiente := sembrarSprintPendiente(t, repo, proyectoID)
	activo, err := NewIniciarSprint(repo).Ejecutar(context.Background(), inputValido(pendiente.ID))
	if err != nil {
		t.Fatalf("no se pudo iniciar el Sprint: %v", err)
	}
	return activo
}

func TestCerrarSprint_FinalizaYArrastraLasNoCompletadas(t *testing.T) {
	sprints := repository.NewSprintRepositoryEnMemoria()
	sprint := sembrarSprintActivo(t, sprints, 1)
	historias := &historiasEnMemoria{items: []domain.BacklogItem{
		historia(1, domain.EstadoCompletada, sprint.ID, puntos(3)),
		historia(2, domain.Estado("En progreso"), sprint.ID, puntos(8)),
		historia(3, domain.EstadoNueva, sprint.ID, nil),
		historia(4, domain.EstadoNueva, sprint.ID+100, puntos(5)),
	}}

	resultado, err := NewCerrarSprint(sprints, historias).Ejecutar(context.Background(), sprint.ID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if resultado.Estado != domain.SprintFinalizado || estadoPersistido(t, sprints, sprint.ID) != domain.SprintFinalizado {
		t.Errorf("el Sprint debía quedar Finalizado")
	}
	arrastradas := sprints.Arrastradas(sprint.ID)
	if len(arrastradas) != 2 || arrastradas[0].ID != 2 || arrastradas[1].ID != 3 {
		t.Fatalf("debían arrastrarse solo las no completadas de este Sprint (2 y 3), se obtuvo %+v", arrastradas)
	}
	if arrastradas[0].EstimacionSP == nil || *arrastradas[0].EstimacionSP != 8 {
		t.Errorf("la historia 2 debía arrastrarse con 8 Story Points, quedó %v", arrastradas[0].EstimacionSP)
	}
	if arrastradas[1].EstimacionSP != nil {
		t.Errorf("la historia 3 debía arrastrarse sin estimar, quedó %d", *arrastradas[1].EstimacionSP)
	}
}

func TestCerrarSprint_SinHistoriasAsignadas(t *testing.T) {
	sprints := repository.NewSprintRepositoryEnMemoria()
	sprint := sembrarSprintActivo(t, sprints, 1)

	resultado, err := NewCerrarSprint(sprints, &historiasEnMemoria{}).Ejecutar(context.Background(), sprint.ID)
	if err != nil {
		t.Fatalf("cerrar un Sprint sin historias no es un error: %v", err)
	}
	if resultado.Estado != domain.SprintFinalizado || estadoPersistido(t, sprints, sprint.ID) != domain.SprintFinalizado {
		t.Errorf("el Sprint debía quedar Finalizado")
	}
}

func TestCerrarSprint_FalloAlPersistirNoCierraElSprint(t *testing.T) {
	sprints := repository.NewSprintRepositoryEnMemoria()
	sprint := sembrarSprintActivo(t, sprints, 1)
	deOtroProyecto := historia(1, domain.EstadoNueva, sprint.ID, nil)
	deOtroProyecto.ProyectoID = 2
	historias := &historiasEnMemoria{items: []domain.BacklogItem{historia(2, domain.EstadoNueva, sprint.ID, nil), deOtroProyecto}}

	_, err := NewCerrarSprint(sprints, historias).Ejecutar(context.Background(), sprint.ID)

	if !errors.Is(err, repository.ErrIntegridadReferencial) {
		t.Fatalf("se esperaba ErrIntegridadReferencial, se obtuvo %v", err)
	}
	if estadoPersistido(t, sprints, sprint.ID) != domain.SprintActivo {
		t.Errorf("el Sprint debía seguir Activo")
	}
	if got := sprints.Arrastradas(sprint.ID); len(got) != 0 {
		t.Errorf("no debía quedar ninguna historia arrastrada, quedaron %+v", got)
	}
}

func TestCerrarSprint_SoloSiEstaActivo(t *testing.T) {
	sprints := repository.NewSprintRepositoryEnMemoria()
	pendiente := sembrarSprintPendiente(t, sprints, 1)
	historias := &historiasEnMemoria{items: []domain.BacklogItem{historia(1, domain.EstadoNueva, pendiente.ID, nil)}}

	_, err := NewCerrarSprint(sprints, historias).Ejecutar(context.Background(), pendiente.ID)

	esperarValidacion(t, err, "estado")
	if got := sprints.Arrastradas(pendiente.ID); len(got) != 0 {
		t.Errorf("no debían arrastrarse historias de un Sprint que no se cerró: %+v", got)
	}
	if estadoPersistido(t, sprints, pendiente.ID) != domain.SprintPendiente {
		t.Errorf("el Sprint debía seguir Pendiente")
	}
}

func TestCerrarSprint_NoEncontrado(t *testing.T) {
	_, err := NewCerrarSprint(repository.NewSprintRepositoryEnMemoria(), &historiasEnMemoria{}).
		Ejecutar(context.Background(), 999)

	if !errors.Is(err, domain.ErrSprintNoEncontrado) {
		t.Fatalf("se esperaba ErrSprintNoEncontrado, se obtuvo %v", err)
	}
}
