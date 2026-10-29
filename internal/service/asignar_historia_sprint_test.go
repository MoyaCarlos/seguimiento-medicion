package service

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/MoyaCarlos/seguimiento-medicion/internal/domain"
	"github.com/MoyaCarlos/seguimiento-medicion/internal/repository"
)

func historiaSinSprint(id, proyectoID int64, estado domain.Estado) domain.BacklogItem {
	return domain.BacklogItem{ID: id, ProyectoID: proyectoID, Estado: estado, EstimacionSP: puntos(5)}
}

func fakeConHistorias(items ...domain.BacklogItem) *fakeAsignacionRepository {
	f := &fakeAsignacionRepository{historias: map[int64]domain.BacklogItem{}}
	for _, h := range items {
		f.historias[h.ID] = h
	}
	return f
}

func sembrarSprintFinalizado(t *testing.T, repo repository.SprintRepository, proyectoID int64) domain.Sprint {
	t.Helper()
	activo := sembrarSprintActivo(t, repo, proyectoID)
	if err := activo.Cerrar(); err != nil {
		t.Fatalf("no se pudo cerrar el Sprint: %v", err)
	}
	if err := repo.Actualizar(context.Background(), activo); err != nil {
		t.Fatalf("no se pudo persistir el Sprint cerrado: %v", err)
	}
	return activo
}

func TestAsignarHistoriaASprint_AsignacionValida(t *testing.T) {
	sprints := repository.NewSprintRepositoryEnMemoria()
	pendiente := sembrarSprintPendiente(t, sprints, 1)
	activo := sembrarSprintActivo(t, sprints, 2)
	historias := fakeConHistorias(historiaSinSprint(1, 1, domain.EstadoNueva), historiaSinSprint(2, 2, domain.EstadoNueva))
	caso := NewAsignarHistoriaASprint(historias, sprints)

	for _, c := range []struct{ historiaID, sprintID int64 }{{1, pendiente.ID}, {2, activo.ID}} {
		h, err := caso.Ejecutar(context.Background(), c.historiaID, c.sprintID)
		if err != nil {
			t.Fatalf("no se esperaba error: %v", err)
		}
		if h.SprintID == nil || *h.SprintID != c.sprintID {
			t.Errorf("SprintID = %v, quería %d", h.SprintID, c.sprintID)
		}
		if h.Estado != domain.EstadoNueva {
			t.Errorf("Estado = %q, debía conservar Nueva", h.Estado)
		}
		if h.EstimacionSP == nil || *h.EstimacionSP != 5 {
			t.Errorf("EstimacionSP = %v, debía conservar 5", h.EstimacionSP)
		}
	}
}

func TestAsignarHistoriaASprint_HistoriaInexistenteGanaSobreSprintInexistente(t *testing.T) {
	caso := NewAsignarHistoriaASprint(fakeConHistorias(), repository.NewSprintRepositoryEnMemoria())

	_, err := caso.Ejecutar(context.Background(), 99, 99)

	if !errors.Is(err, domain.ErrHistoriaNoEncontrada) {
		t.Errorf("err = %v, quería ErrHistoriaNoEncontrada", err)
	}
}

func TestAsignarHistoriaASprint_SprintInexistente(t *testing.T) {
	historias := fakeConHistorias(historiaSinSprint(1, 1, domain.EstadoNueva))
	caso := NewAsignarHistoriaASprint(historias, repository.NewSprintRepositoryEnMemoria())

	_, err := caso.Ejecutar(context.Background(), 1, 99)

	if !errors.Is(err, domain.ErrSprintNoEncontrado) {
		t.Errorf("err = %v, quería ErrSprintNoEncontrado", err)
	}
	if n := historias.llamadasAAsignar(); n != 0 {
		t.Errorf("AsignarASprint se llamó %d veces, no debía llamarse", n)
	}
}

func TestAsignarHistoriaASprint_ErroresDeValidacionSeDevuelvenTalCual(t *testing.T) {
	sprints := repository.NewSprintRepositoryEnMemoria()
	delProyecto := sembrarSprintPendiente(t, sprints, 1)
	finalizado := sembrarSprintFinalizado(t, sprints, 3)
	yaAsignada := historiaSinSprint(4, 1, domain.EstadoNueva)
	yaAsignada.SprintID = &delProyecto.ID

	casos := []struct {
		nombre   string
		historia domain.BacklogItem
		sprintID int64
		esperado error
	}{
		{"otro proyecto", historiaSinSprint(1, 2, domain.EstadoNueva), delProyecto.ID, domain.ErrSprintDeOtroProyecto},
		{"Sprint finalizado", historiaSinSprint(2, 3, domain.EstadoNueva), finalizado.ID, domain.ErrSprintFinalizado},
		{"historia completada", historiaSinSprint(3, 1, domain.EstadoCompletada), delProyecto.ID, domain.ErrHistoriaNoPendiente},
		{"ya asignada", yaAsignada, delProyecto.ID, domain.ErrHistoriaYaAsignada},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			historias := fakeConHistorias(c.historia)

			_, err := NewAsignarHistoriaASprint(historias, sprints).Ejecutar(context.Background(), c.historia.ID, c.sprintID)

			if !errors.Is(err, c.esperado) {
				t.Errorf("err = %v, quería %v", err, c.esperado)
			}
			if n := historias.llamadasAAsignar(); n != 0 {
				t.Errorf("AsignarASprint se llamó %d veces, no debía llamarse", n)
			}
		})
	}
}

// ganaOtraAsignacion lee la historia sin Sprint, pero rechaza la escritura:
// simula que otra asignación ganó entre la lectura y el UPDATE.
type ganaOtraAsignacion struct {
	*fakeAsignacionRepository
}

func (ganaOtraAsignacion) AsignarASprint(context.Context, int64, int64) error {
	return domain.ErrHistoriaYaAsignada
}

func TestAsignarHistoriaASprint_PropagaLaCarreraDelRepositorio(t *testing.T) {
	sprints := repository.NewSprintRepositoryEnMemoria()
	sprint := sembrarSprintPendiente(t, sprints, 1)
	historias := ganaOtraAsignacion{fakeConHistorias(historiaSinSprint(1, 1, domain.EstadoNueva))}

	_, err := NewAsignarHistoriaASprint(historias, sprints).Ejecutar(context.Background(), 1, sprint.ID)

	if !errors.Is(err, domain.ErrHistoriaYaAsignada) {
		t.Errorf("err = %v, quería ErrHistoriaYaAsignada", err)
	}
}

func TestAsignarHistoriaASprint_DosAsignacionesSimultaneas(t *testing.T) {
	sprints := repository.NewSprintRepositoryEnMemoria()
	sprint := sembrarSprintPendiente(t, sprints, 1)
	historias := fakeConHistorias(historiaSinSprint(1, 1, domain.EstadoNueva))
	caso := NewAsignarHistoriaASprint(historias, sprints)

	inicio := make(chan struct{})
	errores := make([]error, 2)
	var wg sync.WaitGroup
	for i := range errores {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-inicio
			_, errores[i] = caso.Ejecutar(context.Background(), 1, sprint.ID)
		}()
	}
	close(inicio)
	wg.Wait()

	var exitos, rechazos int
	for _, err := range errores {
		switch {
		case err == nil:
			exitos++
		case errors.Is(err, domain.ErrHistoriaYaAsignada):
			rechazos++
		default:
			t.Errorf("error inesperado: %v", err)
		}
	}
	if exitos != 1 || rechazos != 1 {
		t.Errorf("éxitos = %d, rechazos = %d; quería 1 y 1", exitos, rechazos)
	}
}
