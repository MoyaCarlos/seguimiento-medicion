package service

import (
	"context"
	"errors"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

// fakeHistoriasDeSprintRepository devuelve las historias sembradas tal cual,
// en el orden en que se cargaron.
type fakeHistoriasDeSprintRepository struct {
	historias []domain.BacklogItem
	llamadas  int
}

func (f *fakeHistoriasDeSprintRepository) ListarPorSprint(_ context.Context, _ int64) ([]domain.BacklogItem, error) {
	f.llamadas++
	return f.historias, nil
}

func TestConsultarSprintBacklog_DevuelveLasHistoriasEnElOrdenDelRepositorio(t *testing.T) {
	sprints := repository.NewSprintRepositoryEnMemoria()
	sprint := sembrarSprintPendiente(t, sprints, 1)
	esperadas := []domain.BacklogItem{
		{ID: 2, ProyectoID: 1, Titulo: "Iniciar sesión", Prioridad: domain.PrioridadMust},
		{ID: 1, ProyectoID: 1, Titulo: "Cerrar sesión", Prioridad: domain.PrioridadShould},
	}
	servicio := NewConsultarSprintBacklog(sprints, &fakeHistoriasDeSprintRepository{historias: esperadas})

	historias, err := servicio.Ejecutar(context.Background(), sprint.ID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if len(historias) != len(esperadas) {
		t.Fatalf("se esperaban %d historias, se obtuvieron %d", len(esperadas), len(historias))
	}
	for i := range esperadas {
		if historias[i].ID != esperadas[i].ID {
			t.Fatalf("se esperaba el orden del repositorio %v, se obtuvo %v", esperadas, historias)
		}
	}
}

func TestConsultarSprintBacklog_SinHistoriasDevuelveVacio(t *testing.T) {
	sprints := repository.NewSprintRepositoryEnMemoria()
	sprint := sembrarSprintPendiente(t, sprints, 1)
	servicio := NewConsultarSprintBacklog(sprints, &fakeHistoriasDeSprintRepository{})

	historias, err := servicio.Ejecutar(context.Background(), sprint.ID)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if historias == nil || len(historias) != 0 {
		t.Errorf("se esperaba un slice vacío no nil, se obtuvo %#v", historias)
	}
}

func TestConsultarSprintBacklog_SprintInexistenteNoConsultaElListado(t *testing.T) {
	historias := &fakeHistoriasDeSprintRepository{}
	servicio := NewConsultarSprintBacklog(repository.NewSprintRepositoryEnMemoria(), historias)

	_, err := servicio.Ejecutar(context.Background(), 99)

	if !errors.Is(err, domain.ErrSprintNoEncontrado) {
		t.Fatalf("se esperaba domain.ErrSprintNoEncontrado, se obtuvo %v", err)
	}
	if historias.llamadas != 0 {
		t.Errorf("no se esperaba consultar el listado de un Sprint inexistente")
	}
}
