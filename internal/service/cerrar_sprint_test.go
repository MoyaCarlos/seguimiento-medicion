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

func (h *historiasEnMemoria) estadoDe(id int64) domain.Estado {
	for _, item := range h.items {
		if item.ID == id {
			return item.Estado
		}
	}
	return ""
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

func (h *historiasEnMemoria) QuitarDeSprint(_ context.Context, historiaID int64) error {
	for i := range h.items {
		if h.items[i].ID == historiaID {
			h.items[i].SprintID = nil
			h.items[i].Estado = domain.EstadoNueva
		}
	}
	return nil
}

func (h *historiasEnMemoria) sprintDe(id int64) *int64 {
	for _, item := range h.items {
		if item.ID == id {
			return item.SprintID
		}
	}
	return nil
}

func historia(id int64, estado domain.Estado, sprintID int64) domain.BacklogItem {
	return domain.BacklogItem{ID: id, Estado: estado, SprintID: &sprintID}
}

func sembrarSprintActivo(t *testing.T, repo repository.SprintRepository, proyectoID int64) domain.Sprint {
	t.Helper()
	pendiente := sembrarSprintPendiente(t, repo, proyectoID)
	activo, err := NewIniciarSprint(repo).Ejecutar(context.Background(), inputValido(pendiente.ID))
	if err != nil {
		t.Fatalf("no se pudo iniciar el Sprint: %v", err)
	}
	return activo
}

func TestCerrarSprint_FinalizaYDevuelveAlBacklogLasNoCompletadas(t *testing.T) {
	sprints := repository.NewSprintRepositoryEnMemoria()
	sprint := sembrarSprintActivo(t, sprints, 1)
	historias := &historiasEnMemoria{items: []domain.BacklogItem{
		historia(1, domain.EstadoCompletada, sprint.ID),
		historia(2, domain.Estado("En progreso"), sprint.ID),
		historia(3, domain.EstadoNueva, sprint.ID+100),
	}}

	resultado, err := NewCerrarSprint(sprints, historias).Ejecutar(context.Background(), sprint.ID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if resultado.Estado != domain.SprintFinalizado || estadoPersistido(t, sprints, sprint.ID) != domain.SprintFinalizado {
		t.Errorf("el Sprint debía quedar Finalizado")
	}
	if s := historias.sprintDe(1); s == nil || *s != sprint.ID {
		t.Errorf("la historia completada debía seguir vinculada al Sprint (FR-009)")
	}
	if historias.sprintDe(2) != nil {
		t.Errorf("la historia no completada debía volver al Product Backlog (FR-008)")
	}
	if estado := historias.estadoDe(2); estado != domain.EstadoNueva {
		t.Errorf("la historia no completada debía volver a estado Nueva, quedó %q", estado)
	}
	if s := historias.sprintDe(3); s == nil || *s != sprint.ID+100 {
		t.Errorf("no debían tocarse historias de otro Sprint")
	}
}

func TestCerrarSprint_SinHistoriasAsignadas(t *testing.T) {
	sprints := repository.NewSprintRepositoryEnMemoria()
	sprint := sembrarSprintActivo(t, sprints, 1)

	resultado, err := NewCerrarSprint(sprints, &historiasEnMemoria{}).Ejecutar(context.Background(), sprint.ID)
	if err != nil {
		t.Fatalf("cerrar un Sprint sin historias no es un error: %v", err)
	}
	if resultado.Estado != domain.SprintFinalizado {
		t.Errorf("el Sprint debía quedar Finalizado")
	}
}

func TestCerrarSprint_SoloSiEstaActivo(t *testing.T) {
	sprints := repository.NewSprintRepositoryEnMemoria()
	pendiente := sembrarSprintPendiente(t, sprints, 1)
	historias := &historiasEnMemoria{items: []domain.BacklogItem{historia(1, domain.EstadoNueva, pendiente.ID)}}

	_, err := NewCerrarSprint(sprints, historias).Ejecutar(context.Background(), pendiente.ID)

	esperarValidacion(t, err, "estado")
	if historias.sprintDe(1) == nil {
		t.Errorf("no debían liberarse historias de un Sprint que no se cerró")
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
